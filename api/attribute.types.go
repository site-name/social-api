package api

import (
	"context"
	"encoding/base64"
	"net/http"
	"unsafe"

	"github.com/Masterminds/squirrel"
	"github.com/graph-gophers/dataloader/v7"
	"github.com/samber/lo"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/store"
	"github.com/sitename/sitename/web"
)

// attribute value

type AttributeValue struct {
	ID       string     `json:"id"`
	Name     *string    `json:"name"`
	Slug     *string    `json:"slug"`
	Value    *string    `json:"value"`
	RichText JSONString `json:"richText"`
	Boolean  *bool      `json:"boolean"`
	Date     *Date      `json:"date"`
	DateTime *DateTime  `json:"dateTime"`
	File     *File      `json:"file"`

	attributeID string

	// Translation *AttributeValueTranslation `json:"translation"`
	// InputType   *AttributeInputTypeEnum    `json:"inputType"`
	// Reference   *string                    `json:"reference"`
}

func SystemAttributeValueToGraphqlAttributeValue(attr *model.AttributeValue) *AttributeValue {
	if attr == nil {
		return nil
	}

	res := &AttributeValue{
		ID:          attr.Id,
		Name:        &attr.Name,
		Slug:        &attr.Slug,
		Value:       &attr.Value,
		Boolean:     attr.Boolean,
		RichText:    JSONString(attr.RichText),
		attributeID: attr.AttributeID,
	}

	if attr.GetAttribute() != nil && attr.Datetime != nil {
		switch attr.GetAttribute().InputType {
		case model.DATE:
			res.Date = &Date{DateTime{*attr.Datetime}}

		case model.DATE_TIME:
			res.DateTime = &DateTime{*attr.Datetime}
		}
	}

	if attr.FileUrl != nil && len(*attr.FileUrl) > 0 {
		res.File = &File{
			URL:         *attr.FileUrl,
			ContentType: attr.ContentType,
		}
	}

	return res
}

func (a *AttributeValue) Translation(ctx context.Context, args struct{ LanguageCode LanguageCodeEnum }) (*AttributeValueTranslation, error) {
	panic("not implemented")
}

func (a *AttributeValue) InputType(ctx context.Context) (*AttributeInputTypeEnum, error) {
	resolveInputType := func(attr Attribute) (*AttributeInputTypeEnum, error) {
		embedCtx, err := GetContextValue[*web.Context](ctx, WebCtx)
		if err != nil {
			return nil, err
		}

		var permToCheck = model.PermissionManageProducts
		if attr.Type != nil && *attr.Type == AttributeTypeEnumPageType {
			permToCheck = model.PermissionManagePages
		}

		if !embedCtx.App.Srv().
			AccountService().
			SessionHasPermissionTo(embedCtx.AppContext.Session(), permToCheck) {
			return attr.InputType, nil
		}

		return nil, model.NewAppError("AttributeValue.InputType", ErrorUnauthorized, nil, "You are not allowed to see this", http.StatusUnauthorized)
	}

	attr, err := AttributesByAttributeIdLoader.Load(ctx, a.attributeID)()
	if err != nil {
		return nil, err
	}

	return resolveInputType(*SystemAttributeToGraphqlAttribute(attr))
}

func (a *AttributeValue) Reference(ctx context.Context) (*string, error) {
	panic("not implemented")
}

// --------------------- Attribute --------------------

type Attribute struct {
	ID              string                   `json:"id"`
	PrivateMetadata []*MetadataItem          `json:"privateMetadata"`
	Metadata        []*MetadataItem          `json:"metadata"`
	InputType       *AttributeInputTypeEnum  `json:"inputType"`
	EntityType      *AttributeEntityTypeEnum `json:"entityType"`
	Name            *string                  `json:"name"`
	Slug            *string                  `json:"slug"`
	Type            *AttributeTypeEnum       `json:"type"`
	Unit            *MeasurementUnitsEnum    `json:"unit"`
	WithChoices     bool                     `json:"withChoices"`

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

		attr: attr,
	}
	if graphqlAttributeInputType := AttributeInputTypeEnum(string(attr.InputType)); graphqlAttributeInputType.IsValid() {
		res.InputType = &graphqlAttributeInputType
	}
	if attr.EntityType != nil {
		if graphqlAttributeEntityType := AttributeEntityTypeEnum(*attr.EntityType); graphqlAttributeEntityType.IsValid() {
			res.EntityType = &graphqlAttributeEntityType
		}
	}
	if graphqlAttributeType := AttributeTypeEnum(attr.Type); graphqlAttributeType.IsValid() {
		res.Type = &graphqlAttributeType
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
		search := *args.Filter.Search

		attributeValues = lo.Filter(attributeValues, func(v *model.AttributeValue, _ int) bool {
			return v.Name == search || v.Slug == search
		})
	}
	totalCount := len(attributeValues)

	p := &graphqlPaginator[*model.AttributeValue, string]{
		data:          attributeValues,
		keyFunc:       keyFunc,
		GraphqlParams: args.GraphqlParams,
	}

	data, hasPrev, hasNext, appErr := p.parse("Attribute.Choices")
	if appErr != nil {
		return nil, appErr
	}

	res := &AttributeValueCountableConnection{
		TotalCount: (*int32)(unsafe.Pointer(&totalCount)),
		Edges: lo.Map(data, func(v *model.AttributeValue, _ int) *AttributeValueCountableEdge {
			return &AttributeValueCountableEdge{
				Node:   SystemAttributeValueToGraphqlAttributeValue(v),
				Cursor: base64.StdEncoding.EncodeToString([]byte(keyFunc(v))),
			}
		}),
	}

	res.PageInfo = &PageInfo{
		HasPreviousPage: hasPrev,
		HasNextPage:     hasNext,
		StartCursor:     &res.Edges[0].Cursor,
		EndCursor:       &res.Edges[len(data)-1].Cursor,
	}

	return res, nil
}

func (a *Attribute) ProductTypes(ctx context.Context, args GraphqlParams) (*ProductTypeCountableConnection, error) {
	embedCtx, err := GetContextValue[*web.Context](ctx, WebCtx)
	if err != nil {
		return nil, err
	}

	productTypes, appErr := embedCtx.App.Srv().
		ProductService().
		ProductTypesByOptions(&model.ProductTypeFilterOption{
			AttributeProducts_AttributeID: squirrel.Eq{store.AttributeProductTableName + ".AttributeID": a.ID},
		})
	if appErr != nil {
		return nil, appErr
	}
	totalCount := len(productTypes)

	p := &graphqlPaginator[*model.ProductType, string]{
		data:          productTypes,
		keyFunc:       func(pt *model.ProductType) string { return pt.Slug },
		GraphqlParams: args,
	}

	data, hasPrev, hasNext, appErr := p.parse("Attribute.ProductTypes")
	if appErr != nil {
		return nil, appErr
	}

	res := &ProductTypeCountableConnection{
		TotalCount: (*int32)(unsafe.Pointer(&totalCount)),
		Edges: lo.Map(data, func(pt *model.ProductType, _ int) *ProductTypeCountableEdge {
			return &ProductTypeCountableEdge{
				Node:   SystemProductTypeToGraphqlProductType(pt),
				Cursor: base64.StdEncoding.EncodeToString([]byte(p.keyFunc(pt))),
			}
		}),
	}

	res.PageInfo = &PageInfo{
		HasNextPage:     hasNext,
		HasPreviousPage: hasPrev,
		StartCursor:     &res.Edges[0].Cursor,
		EndCursor:       &res.Edges[len(data)-1].Cursor,
	}

	return res, nil
}

func (a *Attribute) ProductVariantTypes(ctx context.Context, args GraphqlParams) (*ProductTypeCountableConnection, error) {
	embedCtx, err := GetContextValue[*web.Context](ctx, WebCtx)
	if err != nil {
		return nil, err
	}

	productTypes, appErr := embedCtx.App.Srv().
		ProductService().
		ProductTypesByOptions(&model.ProductTypeFilterOption{
			AttributeVariants_AttributeID: squirrel.Eq{store.AttributeVariantTableName + ".AttributeID": a.ID},
		})
	if appErr != nil {
		return nil, appErr
	}
	totalCount := len(productTypes)

	p := &graphqlPaginator[*model.ProductType, string]{
		data:          productTypes,
		keyFunc:       func(pt *model.ProductType) string { return pt.Slug },
		GraphqlParams: args,
	}

	data, hasPrev, hasNext, appErr := p.parse("Attribute.ProductVariantTypes")
	if appErr != nil {
		return nil, appErr
	}

	res := &ProductTypeCountableConnection{
		TotalCount: (*int32)(unsafe.Pointer(&totalCount)),
		Edges: lo.Map(data, func(pt *model.ProductType, _ int) *ProductTypeCountableEdge {
			return &ProductTypeCountableEdge{
				Node:   SystemProductTypeToGraphqlProductType(pt),
				Cursor: base64.StdEncoding.EncodeToString([]byte(p.keyFunc(pt))),
			}
		}),
	}

	res.PageInfo = &PageInfo{
		HasNextPage:     hasNext,
		HasPreviousPage: hasPrev,
		StartCursor:     &res.Edges[0].Cursor,
		EndCursor:       &res.Edges[len(data)-1].Cursor,
	}

	return res, nil
}

// If return error is nil, meaning current user can perform action.
// if not, user can't
func (a *Attribute) currentUserHasPermissionToAccess(ctx context.Context) error {
	embedCtx, err := GetContextValue[*web.Context](ctx, WebCtx)
	if err != nil {
		return err
	}

	var permToCheck = model.PermissionManageProducts

	if a.Type != nil && *a.Type == AttributeTypeEnumPageType {
		permToCheck = model.PermissionManagePages
	}

	if !embedCtx.
		App.
		Srv().
		AccountService().
		SessionHasPermissionTo(embedCtx.AppContext.Session(), permToCheck) {
		return model.NewAppError("Attribute.currentUserHasPermissionToAccess", ErrorUnauthorized, nil, "you are not allowed to perform this action", http.StatusUnauthorized)
	}

	return nil
}

func (a *Attribute) VisibleInStorefront(ctx context.Context) (bool, error) {
	if err := a.currentUserHasPermissionToAccess(ctx); err != nil {
		return false, err
	}
	return a.attr.VisibleInStoreFront, nil
}

func (a *Attribute) ValueRequired(ctx context.Context) (bool, error) {
	if err := a.currentUserHasPermissionToAccess(ctx); err != nil {
		return false, err
	}
	return a.attr.ValueRequired, nil
}

func (a *Attribute) StorefrontSearchPosition(ctx context.Context) (int32, error) {
	if err := a.currentUserHasPermissionToAccess(ctx); err != nil {
		return 0, err
	}
	return int32(a.attr.StorefrontSearchPosition), nil
}

func (a *Attribute) FilterableInStorefront(ctx context.Context) (bool, error) {
	if err := a.currentUserHasPermissionToAccess(ctx); err != nil {
		return false, err
	}
	return a.attr.FilterableInStorefront, nil
}

func (a *Attribute) FilterableInDashboard(ctx context.Context) (bool, error) {
	if err := a.currentUserHasPermissionToAccess(ctx); err != nil {
		return false, err
	}
	return a.attr.FilterableInDashboard, nil
}

func (a *Attribute) AvailableInGrid(ctx context.Context) (bool, error) {
	if err := a.currentUserHasPermissionToAccess(ctx); err != nil {
		return false, err
	}
	return a.attr.AvailableInGrid, nil
}

func (a *Attribute) Translation(ctx context.Context, args struct{ LanguageCode LanguageCodeEnum }) (*AttributeTranslation, error) {
	panic("not implemented")
}

func attributesByAttributeIdLoader(ctx context.Context, ids []string) []*dataloader.Result[*model.Attribute] {
	var (
		res          = make([]*dataloader.Result[*model.Attribute], len(ids))
		appErr       *model.AppError
		attributes   model.Attributes
		attributeMap = map[string]*model.Attribute{} // keys are attribute ids
	)
	embedCtx, err := GetContextValue[*web.Context](ctx, WebCtx)
	if err != nil {
		goto errorLabel
	}

	attributes, appErr = embedCtx.
		App.
		Srv().
		AttributeService().
		AttributesByOption(&model.AttributeFilterOption{
			Id: squirrel.Eq{store.AttributeTableName + ".Id": ids},
		})
	if appErr != nil {
		err = appErr
		goto errorLabel
	}

	attributeMap = lo.SliceToMap(attributes, func(a *model.Attribute) (string, *model.Attribute) {
		return a.Id, a
	})

	for idx, attrID := range ids {
		res[idx] = &dataloader.Result[*model.Attribute]{Data: attributeMap[attrID]}
	}
	return res

errorLabel:
	for idx := range ids {
		res[idx] = &dataloader.Result[*model.Attribute]{Error: err}
	}
	return res
}

func attributeValuesByAttributeIdLoader(ctx context.Context, attributeIDs []string) []*dataloader.Result[[]*model.AttributeValue] {
	var (
		res             = make([]*dataloader.Result[[]*model.AttributeValue], len(attributeIDs))
		appErr          *model.AppError
		attributeValues model.AttributeValues

		// keys are attribute ids
		attributeValuesMap = map[string][]*model.AttributeValue{}
	)

	embedCtx, err := GetContextValue[*web.Context](ctx, WebCtx)
	if err != nil {
		goto errorLabel
	}

	attributeValues, appErr = embedCtx.App.
		Srv().
		AttributeService().
		FilterAttributeValuesByOptions(model.AttributeValueFilterOptions{
			AttributeID: squirrel.Eq{store.AttributeValueTableName + ".AttributeID": attributeIDs},
		})
	if appErr != nil {
		err = appErr
		goto errorLabel
	}

	for _, value := range attributeValues {
		attributeValuesMap[value.AttributeID] = append(attributeValuesMap[value.AttributeID], value)
	}

	for idx, id := range attributeIDs {
		res[idx] = &dataloader.Result[[]*model.AttributeValue]{Data: attributeValuesMap[id]}
	}
	return res

errorLabel:
	for idx := range attributeIDs {
		res[idx] = &dataloader.Result[[]*model.AttributeValue]{Error: err}
	}
	return res
}

func attributeValueByIdLoader(ctx context.Context, ids []string) []*dataloader.Result[*model.AttributeValue] {
	var (
		res               = make([]*dataloader.Result[*model.AttributeValue], len(ids))
		appErr            *model.AppError
		attributeValues   model.AttributeValues
		attributeValueMap = map[string]*model.AttributeValue{} // keys are attribute value ids
	)

	embedCts, err := GetContextValue[*web.Context](ctx, WebCtx)
	if err != nil {
		goto errorLabel
	}

	attributeValues, appErr = embedCts.App.
		Srv().
		AttributeService().
		FilterAttributeValuesByOptions(model.AttributeValueFilterOptions{
			Id: squirrel.Eq{store.AttributeValueTableName + ".Id": ids},
		})
	if appErr != nil {
		err = appErr
		goto errorLabel
	}

	attributeValueMap = lo.SliceToMap(attributeValues, func(a *model.AttributeValue) (string, *model.AttributeValue) {
		return a.Id, a
	})

	for idx, id := range ids {
		res[idx] = &dataloader.Result[*model.AttributeValue]{Data: attributeValueMap[id]}
	}
	return res

errorLabel:
	for idx := range ids {
		res[idx] = &dataloader.Result[*model.AttributeValue]{Error: err}
	}
	return res
}
