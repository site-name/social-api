package api

import (
	"context"
	"encoding/base64"
	"fmt"
	"strings"
	"unsafe"

	"github.com/Masterminds/squirrel"
	"github.com/graph-gophers/dataloader/v7"
	"github.com/samber/lo"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/web"
)

// attribute value

type AttributeValue struct {
	ID       string     `json:"id"`
	Name     string     `json:"name"`
	Slug     string     `json:"slug"`
	Value    string     `json:"value"`
	RichText JSONString `json:"richText"`
	Boolean  *bool      `json:"boolean"`
	Date     *Date      `json:"date"`
	DateTime *DateTime  `json:"dateTime"`
	File     *File      `json:"file"`

	attributeID string

	// Translation *AttributeValueTranslation `json:"translation"`
	// InputType   *model.AttributeInputType    `json:"inputType"`
	// Reference   *string                    `json:"reference"`
}

func SystemAttributeValueToGraphqlAttributeValue(attrValue *model.AttributeValue) *AttributeValue {
	if attrValue == nil {
		return nil
	}

	res := &AttributeValue{
		ID:          attrValue.Id,
		Name:        attrValue.Name,
		Slug:        attrValue.Slug,
		Value:       attrValue.Value,
		Boolean:     attrValue.Boolean,
		RichText:    JSONString(attrValue.RichText),
		attributeID: attrValue.AttributeID,
	}

	if attr := attrValue.Attribute; attr != nil && attrValue.Datetime != nil {
		switch attr.InputType {
		case model.AttributeInputTypeDate:
			res.Date = &Date{DateTime{*attrValue.Datetime}}

		case model.AttributeInputTypeDateTime:
			res.DateTime = &DateTime{*attrValue.Datetime}
		}
	}

	if attrValue.FileUrl != nil && len(*attrValue.FileUrl) > 0 {
		res.File = &File{
			URL:         *attrValue.FileUrl,
			ContentType: attrValue.ContentType,
		}
	}

	return res
}

func (a *AttributeValue) Translation(ctx context.Context, args struct{ LanguageCode LanguageCodeEnum }) (*AttributeValueTranslation, error) {
	panic("not implemented")
}

func (a *AttributeValue) InputType(ctx context.Context) (*model.AttributeInputType, error) {
	attr, err := AttributesByAttributeIdLoader.Load(ctx, a.attributeID)()
	if err != nil {
		return nil, err
	}

	return &attr.InputType, nil
}

// the result would has format of "EntityType:slug"
func (a *AttributeValue) Reference(ctx context.Context) (*string, error) {
	attribute, err := AttributesByAttributeIdLoader.Load(ctx, a.attributeID)()
	if err != nil {
		return nil, err
	}

	if attribute.InputType != model.AttributeInputTypeReference {
		return nil, nil
	}

	splitSlug := strings.Split(a.Slug, "_")
	if len(splitSlug) < 2 || attribute.EntityType == nil { // we need at least length of 2 to get element at 1st index
		return nil, nil
	}

	referenceID := base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", *attribute.EntityType, splitSlug[1])))
	return &referenceID, nil
}

// --------------------- Attribute --------------------

type Attribute struct {
	ID              string                     `json:"id"`
	PrivateMetadata []*MetadataItem            `json:"privateMetadata"`
	Metadata        []*MetadataItem            `json:"metadata"`
	InputType       *model.AttributeInputType  `json:"inputType"`
	EntityType      *model.AttributeEntityType `json:"entityType"`
	Name            *string                    `json:"name"`
	Slug            *string                    `json:"slug"`
	Type            *model.AttributeType       `json:"type"`
	Unit            *MeasurementUnitsEnum      `json:"unit"`
	WithChoices     bool                       `json:"withChoices"`

	attr *model.Attribute

	// ValueRequired   bool                     `json:"valueRequired"`
	// Choices                  *AttributeValueCountableConnection `json:"choices"`
	// VisibleInStorefront      bool                               `json:"visibleInStorefront"`
	// FilterableInStorefront   bool                               `json:"filterableInStorefront"`
	// FilterableInDashboard    bool                               `json:"filterableInDashboard"`
	// AvailableInGrid          bool                               `json:"availableInGrid"`
	// Translation              *AttributeTranslation              `json:"translation"`
	// StorefrontSearchPosition int32                              `json:"storefrontSearchPosition"`
	// ProductTypes             *ProductTypeCountableConnection    `json:"productTypes"`
	// ProductVariantTypes      *ProductTypeCountableConnection    `json:"productVariantTypes"`
}

func SystemAttributeToGraphqlAttribute(attr *model.Attribute) *Attribute {
	if attr == nil {
		return nil
	}

	res := &Attribute{
		ID:              attr.Id,
		Metadata:        MetadataToSlice(attr.Metadata),
		PrivateMetadata: MetadataToSlice(attr.PrivateMetadata),
		Name:            &attr.Name,
		Slug:            &attr.Slug,
		WithChoices:     model.TYPES_WITH_CHOICES.Contains(attr.InputType),
		InputType:       &attr.InputType,
		EntityType:      attr.EntityType,
		attr:            attr,
		Type:            &attr.Type,
	}

	if attr.Unit != nil {
		if graphqlAttributeUnit := MeasurementUnitsEnum(*attr.Unit); graphqlAttributeUnit.IsValid() {
			res.Unit = &graphqlAttributeUnit
		}
	}

	return res
}

func (a *Attribute) Choices(
	ctx context.Context,
	args struct {
		Filter *AttributeValueFilterInput
		SortBy *AttributeChoicesSortingInput
		GraphqlParams
	},
) (*AttributeValueCountableConnection, error) {
	if !model.TYPES_WITH_CHOICES.Contains(a.attr.InputType) {
		return nil, nil
	}
	appErr := args.GraphqlParams.Validate("Attribute.Choices")
	if appErr != nil {
		return nil, appErr
	}

	attributeValues, err := AttributeValuesByAttributeIdLoader.Load(ctx, a.ID)()
	if err != nil {
		return nil, err
	}

	keyFunc := func(v *model.AttributeValue) string { return v.Name }
	if args.SortBy != nil && args.SortBy.Field == AttributeChoicesSortFieldSlug {
		keyFunc = func(v *model.AttributeValue) string { return v.Slug }
	}

	// parse filter
	if args.Filter != nil && args.Filter.Search != nil {
		search := strings.ToLower(*args.Filter.Search)

		attributeValues = lo.Filter(attributeValues, func(v *model.AttributeValue, _ int) bool {
			lowerName := strings.ToLower(v.Name)
			lowerSlug := strings.ToLower(v.Slug)

			return strings.Contains(lowerName, search) || strings.Contains(lowerSlug, search)
		})
	}

	res, appErr := newGraphqlPaginator(
		attributeValues,
		keyFunc,
		SystemAttributeValueToGraphqlAttributeValue,
		args.GraphqlParams).parse("Attribute.Choices")

	if appErr != nil {
		return nil, appErr
	}

	return (*AttributeValueCountableConnection)(unsafe.Pointer(res)), nil
}

func (a *Attribute) ProductTypes(ctx context.Context, args GraphqlParams) (*ProductTypeCountableConnection, error) {
	embedCtx := GetContextValue[*web.Context](ctx, WebCtx)

	productTypes, appErr := embedCtx.App.Srv().
		ProductService().
		ProductTypesByOptions(&model.ProductTypeFilterOption{
			AttributeProducts_AttributeID: squirrel.Eq{model.AttributeProductTableName + ".AttributeID": a.ID},
		})
	if appErr != nil {
		return nil, appErr
	}

	keyFunc := func(pt *model.ProductType) string { return pt.Slug }
	res, appErr := newGraphqlPaginator(productTypes, keyFunc, SystemProductTypeToGraphqlProductType, args).parse("Attribute.ProductTypes")
	if appErr != nil {
		return nil, appErr
	}

	return (*ProductTypeCountableConnection)(unsafe.Pointer(res)), nil
}

func (a *Attribute) ProductVariantTypes(ctx context.Context, args GraphqlParams) (*ProductTypeCountableConnection, error) {
	embedCtx := GetContextValue[*web.Context](ctx, WebCtx)

	productTypes, appErr := embedCtx.App.Srv().
		ProductService().
		ProductTypesByOptions(&model.ProductTypeFilterOption{
			AttributeVariants_AttributeID: squirrel.Eq{model.AttributeVariantTableName + ".AttributeID": a.ID},
		})
	if appErr != nil {
		return nil, appErr
	}

	keyFunc := func(pt *model.ProductType) string { return pt.Slug }
	res, appErr := newGraphqlPaginator(productTypes, keyFunc, SystemProductTypeToGraphqlProductType, args).parse("Attribute.ProductVariantTypes")
	if appErr != nil {
		return nil, appErr
	}

	return (*ProductTypeCountableConnection)(unsafe.Pointer(res)), nil
}

// If return error is nil, meaning current user can perform action.
// if not, user can't
func (a *Attribute) currentUserHasPermissionToAccess(ctx context.Context, apiName string) error {
	var permsToCheck = model.Permissions{model.PermissionReadProduct, model.PermissionCreateProduct, model.PermissionDeleteProduct, model.PermissionUpdateProduct}
	if a.Type != nil && *a.Type == model.PAGE_TYPE {
		permsToCheck = model.Permissions{model.PermissionReadPage, model.PermissionCreatePage, model.PermissionUpdatePage, model.PermissionUpdatePage}
	}

	embedCtx := GetContextValue[*web.Context](ctx, WebCtx)
	embedCtx.CheckAuthenticatedAndHasPermissionToAll(permsToCheck...)
	return embedCtx.Err
}

func (a *Attribute) VisibleInStorefront(ctx context.Context) (bool, error) {
	if err := a.currentUserHasPermissionToAccess(ctx, "Attribute.VisibleInStorefront"); err != nil {
		return false, err
	}
	return a.attr.VisibleInStoreFront, nil
}

func (a *Attribute) ValueRequired(ctx context.Context) (bool, error) {
	if err := a.currentUserHasPermissionToAccess(ctx, "Attribute.ValueRequired"); err != nil {
		return false, err
	}
	return a.attr.ValueRequired, nil
}

func (a *Attribute) StorefrontSearchPosition(ctx context.Context) (int32, error) {
	if err := a.currentUserHasPermissionToAccess(ctx, "Attribute.StorefrontSearchPosition"); err != nil {
		return 0, err
	}
	return int32(a.attr.StorefrontSearchPosition), nil
}

func (a *Attribute) FilterableInStorefront(ctx context.Context) (bool, error) {
	if err := a.currentUserHasPermissionToAccess(ctx, "Attribute.FilterableInStorefront"); err != nil {
		return false, err
	}
	return a.attr.FilterableInStorefront, nil
}

func (a *Attribute) FilterableInDashboard(ctx context.Context) (bool, error) {
	if err := a.currentUserHasPermissionToAccess(ctx, "Attribute.FilterableInDashboard"); err != nil {
		return false, err
	}
	return a.attr.FilterableInDashboard, nil
}

func (a *Attribute) AvailableInGrid(ctx context.Context) (bool, error) {
	if err := a.currentUserHasPermissionToAccess(ctx, "Attribute.AvailableInGrid"); err != nil {
		return false, err
	}
	return a.attr.AvailableInGrid, nil
}

func (a *Attribute) Translation(ctx context.Context, args struct{ LanguageCode LanguageCodeEnum }) (*AttributeTranslation, error) {
	panic("not implemented")
}

func attributesByAttributeIdLoader(ctx context.Context, ids []string) []*dataloader.Result[*model.Attribute] {
	res := make([]*dataloader.Result[*model.Attribute], len(ids))
	embedCtx := GetContextValue[*web.Context](ctx, WebCtx)

	attributes, appErr := embedCtx.
		App.
		Srv().
		AttributeService().
		AttributesByOption(&model.AttributeFilterOption{
			Conditions: squirrel.Eq{model.AttributeTableName + ".Id": ids},
		})
	if appErr != nil {
		for idx := range ids {
			res[idx] = &dataloader.Result[*model.Attribute]{Error: appErr}
		}
		return res
	}

	attributeMap := lo.SliceToMap(attributes, func(a *model.Attribute) (string, *model.Attribute) {
		return a.Id, a
	})

	for idx, attrID := range ids {
		res[idx] = &dataloader.Result[*model.Attribute]{Data: attributeMap[attrID]}
	}
	return res
}

func attributeValuesByAttributeIdLoader(ctx context.Context, attributeIDs []string) []*dataloader.Result[[]*model.AttributeValue] {
	res := make([]*dataloader.Result[[]*model.AttributeValue], len(attributeIDs))
	embedCtx := GetContextValue[*web.Context](ctx, WebCtx)

	attributeValues, appErr := embedCtx.App.
		Srv().
		AttributeService().
		FilterAttributeValuesByOptions(model.AttributeValueFilterOptions{
			Conditions: squirrel.Eq{model.AttributeValueTableName + ".AttributeID": attributeIDs},
		})
	if appErr != nil {
		for idx := range attributeIDs {
			res[idx] = &dataloader.Result[[]*model.AttributeValue]{Error: appErr}
		}
		return res
	}

	// keys are attribute ids
	var attributeValuesMap = map[string][]*model.AttributeValue{}
	for _, value := range attributeValues {
		attributeValuesMap[value.AttributeID] = append(attributeValuesMap[value.AttributeID], value)
	}

	for idx, id := range attributeIDs {
		res[idx] = &dataloader.Result[[]*model.AttributeValue]{Data: attributeValuesMap[id]}
	}
	return res
}

func attributeValueByIdLoader(ctx context.Context, ids []string) []*dataloader.Result[*model.AttributeValue] {
	res := make([]*dataloader.Result[*model.AttributeValue], len(ids))
	embedCtx := GetContextValue[*web.Context](ctx, WebCtx)

	attributeValues, appErr := embedCtx.App.
		Srv().
		AttributeService().
		FilterAttributeValuesByOptions(model.AttributeValueFilterOptions{
			Conditions: squirrel.Eq{model.AttributeValueTableName + ".Id": ids},
		})
	if appErr != nil {
		for idx := range ids {
			res[idx] = &dataloader.Result[*model.AttributeValue]{Error: appErr}
		}
		return res
	}

	attributeValueMap := lo.SliceToMap(attributeValues, func(a *model.AttributeValue) (string, *model.AttributeValue) {
		return a.Id, a
	})
	for idx, id := range ids {
		res[idx] = &dataloader.Result[*model.AttributeValue]{Data: attributeValueMap[id]}
	}
	return res
}
