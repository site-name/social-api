package api

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"
	"net/http"
	"unsafe"

	"github.com/Masterminds/squirrel"
	"github.com/gosimple/slug"
	"github.com/samber/lo"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/web"
)

// NOTE: Refer to ./schemas/attribute.graphqls for details on directive used
func (r *Resolver) AttributeCreate(ctx context.Context, args struct{ Input AttributeCreateInput }) (*AttributeCreate, error) {
	embedCtx := GetContextValue[*web.Context](ctx, WebCtx)

	// clean input
	inputType := args.Input.InputType
	if inputType != nil && *inputType == model.AttributeInputTypeReference &&
		(args.Input.EntityType == nil || !args.Input.EntityType.IsValid()) {
		return nil, model.NewAppError("AttributeCreate", "api.attribute.entity_type_missing.app_error", nil, "entity type is required when REFERENCE input type is used", http.StatusBadRequest)
	}

	// no need to initialize Slug here, since it will be done in PreSave() method
	attribute := &model.Attribute{
		Name:       args.Input.Name,
		Type:       args.Input.Type,
		Unit:       (*string)(unsafe.Pointer(args.Input.Unit)),
		EntityType: args.Input.EntityType,
	}
	if inputType != nil {
		attribute.InputType = *inputType
	}
	if v := args.Input.ValueRequired; v != nil {
		attribute.ValueRequired = *v
	}
	if v := args.Input.IsVariantOnly; v != nil {
		attribute.IsVariantOnly = *v
	}
	if v := args.Input.VisibleInStorefront; v != nil {
		attribute.VisibleInStoreFront = *v
	}
	if v := args.Input.FilterableInStorefront; v != nil {
		attribute.FilterableInStorefront = *v
	}
	if v := args.Input.FilterableInDashboard; v != nil {
		attribute.FilterableInDashboard = *v
	}
	if v := args.Input.StorefrontSearchPosition; v != nil {
		attribute.StorefrontSearchPosition = int(*v)
	}
	if v := args.Input.AvailableInGrid; v != nil {
		attribute.AvailableInGrid = *v
	}

	// clean attribute
	appErr := cleanAttributeSettings(attribute, &args.Input)
	if appErr != nil {
		return nil, appErr
	}

	// clean values
	validatedAttributeValues, appErr := newAttributeMixin[*AttributeCreateInput](embedCtx.App.Srv(), "values").
		cleanValues(&args.Input, attribute)
	if appErr != nil {
		return nil, appErr
	}

	// save
	savedAttr, appErr := embedCtx.App.Srv().AttributeService().UpsertAttribute(attribute)
	if appErr != nil {
		return nil, appErr
	}

	// save attribute values
	_, appErr = embedCtx.App.Srv().AttributeService().BulkUpsertAttributeValue(nil, validatedAttributeValues)
	if appErr != nil {
		return nil, appErr
	}

	return &AttributeCreate{
		Attribute: SystemAttributeToGraphqlAttribute(savedAttr),
	}, nil
}

// NOTE: Refer to ./schemas/attribute.graphqls for details on directive used
func (r *Resolver) AttributeDelete(ctx context.Context, args struct{ Id UUID }) (*AttributeDelete, error) {
	embedCtx := GetContextValue[*web.Context](ctx, WebCtx)

	_, appErr := embedCtx.App.Srv().AttributeService().DeleteAttributes(args.Id.String())
	if appErr != nil {
		return nil, appErr
	}
	return &AttributeDelete{
		Attribute: SystemAttributeToGraphqlAttribute(&model.Attribute{Id: args.Id.String()}),
	}, nil
}

// NOTE: Refer to ./schemas/attribute.graphqls for details on directive used
func (r *Resolver) AttributeUpdate(ctx context.Context, args struct {
	Id    UUID
	Input AttributeUpdateInput
}) (*AttributeUpdate, error) {
	embedCtx := GetContextValue[*web.Context](ctx, WebCtx)
	// get attribute
	attribute, appErr := embedCtx.App.Srv().AttributeService().AttributeByOption(&model.AttributeFilterOption{
		Conditions: squirrel.Eq{model.AttributeTableName + ".Id": args.Id},
	})
	if appErr != nil {
		return nil, appErr
	}

	appErr = cleanAttributeSettings(attribute, &args.Input)
	if appErr != nil {
		return nil, appErr
	}

	// clean add values, add them to db
	addAttributeValues, appErr := newAttributeMixin[*AttributeUpdateInput](embedCtx.App.Srv(), "add_values").
		cleanValues(&args.Input, attribute)
	if appErr != nil {
		return nil, appErr
	}
	_, appErr = embedCtx.App.Srv().AttributeService().BulkUpsertAttributeValue(nil, addAttributeValues)
	if appErr != nil {
		return nil, appErr
	}

	// clean remove values
	if len(args.Input.RemoveValues) > 0 {
		removeValues, appErr := embedCtx.App.Srv().AttributeService().FilterAttributeValuesByOptions(model.AttributeValueFilterOptions{
			Conditions: squirrel.Eq{model.AttributeValueTableName + ".Id": args.Input.RemoveValues},
		})
		if appErr != nil {
			return nil, appErr
		}
		// validate all found attribute values are children of attribute
		if !lo.EveryBy(removeValues, func(vl *model.AttributeValue) bool { return vl.AttributeID == attribute.Id }) {
			return nil, model.NewAppError("AttributeUpdate", model.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "removeValues"}, "one of attribute values does not belong to given attribute", http.StatusBadRequest)
		}

		// remove attribute values designated:
		removeAttributeValueIds := *(*[]string)(unsafe.Pointer(&args.Input.RemoveValues))
		_, appErr = embedCtx.App.Srv().AttributeService().DeleteAttributeValues(nil, removeAttributeValueIds...)
		if appErr != nil {
			return nil, appErr
		}
	}

	// construct instance
	if v := args.Input.Name; v != nil && *v != attribute.Name {
		attribute.Name = *v
	}
	if v := args.Input.Slug; v != nil && *v != attribute.Slug {
		attribute.Slug = *v
	}
	if v := args.Input.Unit; v != nil && v.IsValid() {
		attribute.Unit = (*string)(unsafe.Pointer(v))
	}
	if v := args.Input.ValueRequired; v != nil {
		attribute.ValueRequired = *v
	}
	if v := args.Input.IsVariantOnly; v != nil {
		attribute.IsVariantOnly = *v
	}
	if v := args.Input.VisibleInStorefront; v != nil {
		attribute.VisibleInStoreFront = *v
	}
	if v := args.Input.FilterableInStorefront; v != nil {
		attribute.FilterableInStorefront = *v
	}
	if v := args.Input.FilterableInDashboard; v != nil {
		attribute.FilterableInDashboard = *v
	}
	if v := args.Input.AvailableInGrid; v != nil {
		attribute.AvailableInGrid = *v
	}
	if v := args.Input.StorefrontSearchPosition; v != nil {
		attribute.StorefrontSearchPosition = int(*v)
	}

	// update attribute in db
	savedAttr, appErr := embedCtx.App.Srv().AttributeService().UpsertAttribute(attribute)
	if appErr != nil {
		return nil, appErr
	}

	return &AttributeUpdate{
		Attribute: SystemAttributeToGraphqlAttribute(savedAttr),
	}, nil
}

func (r *Resolver) AttributeTranslate(ctx context.Context, args struct {
	Id           UUID
	Input        NameTranslationInput
	LanguageCode LanguageCodeEnum
}) (*AttributeTranslate, error) {
	panic(fmt.Errorf("not implemented"))
}

// NOTE: Refer to ./schemas/attribute.graphqls for details on directive used
func (r *Resolver) AttributeBulkDelete(ctx context.Context, args struct{ Ids []UUID }) (*AttributeBulkDelete, error) {
	embedCtx := GetContextValue[*web.Context](ctx, WebCtx)
	strAttributeIds := *(*[]string)(unsafe.Pointer(&args.Ids))
	count, appErr := embedCtx.App.Srv().AttributeService().DeleteAttributes(strAttributeIds...)
	if appErr != nil {
		return nil, appErr
	}
	return &AttributeBulkDelete{
		Count: int32(count),
	}, nil
}

// NOTE: Refer to ./schemas/attribute.graphqls for details on directive used
func (r *Resolver) AttributeValueBulkDelete(ctx context.Context, args struct{ Ids []UUID }) (*AttributeValueBulkDelete, error) {
	embedCtx := GetContextValue[*web.Context](ctx, WebCtx)
	strAttributeValueIds := *(*[]string)(unsafe.Pointer(&args.Ids))
	count, appErr := embedCtx.App.Srv().AttributeService().DeleteAttributeValues(nil, strAttributeValueIds...)
	if appErr != nil {
		return nil, appErr
	}
	return &AttributeValueBulkDelete{
		Count: int32(count),
	}, nil
}

// NOTE: Refer to ./schemas/attribute.graphqls for details on directive used
func (r *Resolver) AttributeValueCreate(ctx context.Context, args struct {
	AttributeID UUID
	Input       AttributeValueCreateInput
}) (*AttributeValueCreate, error) {
	embedCtx := GetContextValue[*web.Context](ctx, WebCtx)
	attribute, appErr := embedCtx.App.Srv().AttributeService().AttributeByOption(&model.AttributeFilterOption{
		Conditions: squirrel.Eq{model.AttributeTableName + ".Id": args.AttributeID},
	})
	if appErr != nil {
		return nil, appErr
	}

	// clean input
	if attribute.InputType == model.AttributeInputTypeSwatch {
		appErr = validateSwatchAttributeValue(&args.Input)
		if appErr != nil {
			return nil, appErr
		}
	} else {
		fileUrl := args.Input.FileURL
		contentType := args.Input.ContentType
		if (fileUrl != nil && *fileUrl != "") ||
			(contentType != nil && *contentType != "") {
			return nil, model.NewAppError("AttributeValueCreate", model.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "fileUrl or contentType"}, "fieUrl and contentType can only be defined for swatch attribute", http.StatusBadRequest)
		}
	}

	// construct instance
	attrValue := &model.AttributeValue{
		Name:        args.Input.Name,
		RichText:    model.StringInterface(args.Input.RichText),
		FileUrl:     args.Input.FileURL,
		ContentType: args.Input.ContentType,
		AttributeID: attribute.Id,
	}
	if v := args.Input.getValue(); v != nil && *v != "" {
		attrValue.Value = *v
	}
	// upsert
	savedAttrValue, appErr := embedCtx.App.Srv().AttributeService().UpsertAttributeValue(attrValue)
	if appErr != nil {
		return nil, appErr
	}

	return &AttributeValueCreate{
		Attribute:      SystemAttributeToGraphqlAttribute(attribute),
		AttributeValue: SystemAttributeValueToGraphqlAttributeValue(savedAttrValue),
	}, nil
}

// NOTE: Refer to ./schemas/attribute.graphqls for details on directive used
func (r *Resolver) AttributeValueDelete(ctx context.Context, args struct{ Id UUID }) (*AttributeValueDelete, error) {
	embedCtx := GetContextValue[*web.Context](ctx, WebCtx)
	attrValues, appErr := embedCtx.App.Srv().AttributeService().FilterAttributeValuesByOptions(model.AttributeValueFilterOptions{
		Conditions: squirrel.Eq{model.AttributeValueTableName + ".Id": args.Id},
	})
	if appErr != nil {
		return nil, appErr
	}
	if len(attrValues) == 0 {
		return nil, model.NewAppError("AttributeValueDelete", model.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "Id"}, "id="+args.Id.String()+" is in valid", http.StatusBadRequest)
	}
	attrValue := attrValues[0]

	_, appErr = embedCtx.App.Srv().AttributeService().DeleteAttributeValues(nil, args.Id.String())
	if appErr != nil {
		return nil, appErr
	}

	return &AttributeValueDelete{
		Attribute:      SystemAttributeToGraphqlAttribute(&model.Attribute{Id: attrValue.AttributeID}),
		AttributeValue: SystemAttributeValueToGraphqlAttributeValue(attrValue),
	}, nil
}

// NOTE: Refer to ./schemas/attribute.graphqls for details on directive used
func (r *Resolver) AttributeValueUpdate(ctx context.Context, args struct {
	Id    UUID
	Input AttributeValueUpdateInput
}) (*AttributeValueUpdate, error) {
	embedCtx := GetContextValue[*web.Context](ctx, WebCtx)

	attrValues, appErr := embedCtx.App.Srv().AttributeService().FilterAttributeValuesByOptions(model.AttributeValueFilterOptions{
		Conditions: squirrel.Eq{model.AttributeValueTableName + ".Id": args.Id},
	})
	if appErr != nil {
		return nil, appErr
	}
	if len(attrValues) == 0 {
		return nil, model.NewAppError("AttributeValueUpdate", model.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "Id"}, "id="+args.Id.String()+" is in valid", http.StatusBadRequest)
	}

	attrValue := attrValues[0]

	// clean input
	if v := args.Input.Value; v != nil && *v != "" {
		args.Input.FileURL = nil
		args.Input.ContentType = nil
	} else if v := args.Input.FileURL; v != nil && *v != "" {
		args.Input.Value = nil
	}

	// update value
	attrValue.Name = args.Input.Name
	if v := args.Input.Value; v != nil {
		attrValue.Value = *v
	}
	attrValue.RichText = model.StringInterface(args.Input.RichText)
	attrValue.FileUrl = args.Input.FileURL
	attrValue.ContentType = args.Input.ContentType

	savedAttrValue, appErr := embedCtx.App.Srv().AttributeService().UpsertAttributeValue(attrValue)
	if appErr != nil {
		return nil, appErr
	}

	return &AttributeValueUpdate{
		Attribute:      SystemAttributeToGraphqlAttribute(&model.Attribute{Id: savedAttrValue.AttributeID}),
		AttributeValue: SystemAttributeValueToGraphqlAttributeValue(savedAttrValue),
	}, nil
}

func (r *Resolver) AttributeValueTranslate(ctx context.Context, args struct {
	Id           UUID
	Input        AttributeValueTranslationInput
	LanguageCode LanguageCodeEnum
}) (*AttributeValueTranslate, error) {
	panic(fmt.Errorf("not implemented"))
}

type AttributeReorderValuesArgs struct {
	AttributeID UUID
	Moves       []*ReorderInput
}

func (r *Resolver) AttributeReorderValues(ctx context.Context, args AttributeReorderValuesArgs) (*AttributeReorderValues, error) {
	embedCtx := GetContextValue[*web.Context](ctx, WebCtx)
	// find attribute with given id
	attribute, appErr := embedCtx.App.Srv().AttributeService().AttributeByOption(&model.AttributeFilterOption{
		Conditions:                     squirrel.Eq{model.AttributeTableName + ".Id": args.AttributeID},
		PrefetchRelatedAttributeValues: true,
	})
	if appErr != nil {
		return nil, appErr
	}
	attributeValues := attribute.AttributeValues
	attributeValueMap := lo.SliceToMap(attributeValues, func(a *model.AttributeValue) (string, bool) { return a.Id, true })
	operations := map[string]*int{}

	for _, move := range args.Moves {
		if !attributeValueMap[move.ID] { // not contains
			return nil, model.NewAppError("AttributeReorderValues", model.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "move.Id"}, fmt.Sprintf("attribute value with id=%s does not belong to the attribute", move.ID), http.StatusBadRequest)
		}

		operations[move.ID] = (*int)(unsafe.Pointer(move.SortOrder))
	}

	appErr = embedCtx.App.Srv().AttributeService().PerformReordering(attributeValues, operations)
	if appErr != nil {
		return nil, appErr
	}

	return &AttributeReorderValues{
		Attribute: SystemAttributeToGraphqlAttribute(attribute),
	}, nil
}

type AttributesArgs struct {
	Filter      *AttributeFilterInput
	SortBy      *AttributeSortingInput
	ChannelSlug *string
	GraphqlParams
}

func (args *AttributesArgs) parse() (*model.AttributeFilterOption, *model.AppError) {
	// validate params
	var attributeFilter = &model.AttributeFilterOption{}

	if args.Filter != nil {
		var appErr *model.AppError
		attributeFilter, appErr = args.Filter.parse("AttributesArgs.parse")
		if appErr != nil {
			return nil, appErr
		}
	}

	if args.ChannelSlug != nil {
		if !slug.IsSlug(*args.ChannelSlug) {
			return nil, model.NewAppError("AttributesArgs.parse", model.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "channel slug"}, "please provide valid channel slug", http.StatusBadRequest)
		}
		attributeFilter.ChannelSlug = args.ChannelSlug
	}

	if args.SortBy != nil && !args.SortBy.Field.IsValid() {
		return nil, model.NewAppError("AttributesArgs.parse", model.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "sort field"}, "please provide valid sort field", http.StatusBadRequest)
	}

	paginValue, appErr := args.GraphqlParams.Parse("AttributesArgs.parse")
	if appErr != nil {
		return nil, appErr
	}

	attributeFilter.GraphqlPaginationValues = *paginValue
	if attributeFilter.GraphqlPaginationValues.OrderBy == "" {

		// sodrt attributes default by slugs
		attributeSortFields := attributeSortFieldMap[AttributeSortFieldSlug].fields
		if args.SortBy != nil {
			attributeSortFields = attributeSortFieldMap[args.SortBy.Field].fields
		}

		orderDirection := args.GraphqlParams.orderDirection()
		attributeFilter.GraphqlPaginationValues.OrderBy = attributeSortFields.
			Map(func(_ int, item string) string { return item + " " + orderDirection }).
			Join(",")
	}

	return attributeFilter, nil
}

func (r *Resolver) Attributes(ctx context.Context, args AttributesArgs) (*AttributeCountableConnection, error) {
	attrFilterOpts, appErr := args.parse()
	if appErr != nil {
		return nil, appErr
	}

	embedCtx := GetContextValue[*web.Context](ctx, WebCtx)

	// find attributes by options
	attributes, appErr := embedCtx.App.Srv().AttributeService().AttributesByOption(attrFilterOpts)
	if appErr != nil {
		return nil, appErr
	}
	// count total number of attributes that satisfy given options
	totalCount, err := embedCtx.App.Srv().Store.Attribute().CountByOptions(attrFilterOpts)
	if err != nil {
		return nil, model.NewAppError("Attributes", "app.attribute.count_by_options.app_error", nil, err.Error(), http.StatusInternalServerError)
	}

	keyFunc := attributeSortFieldMap[AttributeSortFieldSlug].keyFunc
	if args.SortBy != nil {
		keyFunc = attributeSortFieldMap[args.SortBy.Field].keyFunc
	}

	res := constructCountableConnection(attributes, totalCount, args.GraphqlParams, keyFunc, SystemAttributeToGraphqlAttribute)
	return (*AttributeCountableConnection)(unsafe.Pointer(res)), nil
}

func (r *Resolver) Attribute(ctx context.Context, args struct {
	Id   *UUID
	Slug *string
}) (*Attribute, error) {
	if (args.Id == nil && args.Slug == nil) || (args.Id != nil && args.Slug != nil) {
		return nil, model.NewAppError("Attribute", model.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "id/slug"}, "please provide id or slug, not both", http.StatusBadRequest)
	}
	if args.Slug != nil && !slug.IsSlug(*args.Slug) {
		return nil, model.NewAppError("Attribute", model.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "slug"}, "please provide valid slug", http.StatusBadRequest)
	}

	var conditions squirrel.Sqlizer
	if args.Id != nil {
		conditions = squirrel.Expr(model.AttributeTableName+".Id = ?", *args.Id)
	} else {
		conditions = squirrel.Expr(model.AttributeTableName+".Slug = ?", *args.Slug)
	}

	embedCtx := GetContextValue[*web.Context](ctx, WebCtx)
	attribute, appErr := embedCtx.App.Srv().AttributeService().AttributeByOption(&model.AttributeFilterOption{
		Conditions: conditions,
	})
	if appErr != nil {
		return nil, appErr
	}

	return SystemAttributeToGraphqlAttribute(attribute), nil
}
