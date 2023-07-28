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
	"github.com/sitename/sitename/app"
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
func (r *Resolver) AttributeDelete(ctx context.Context, args struct{ Id string }) (*AttributeDelete, error) {
	// validate argument(s)
	if !model.IsValidId(args.Id) {
		return nil, model.NewAppError("AttributeDelete", app.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "id"}, "id = "+args.Id+" is invalid id", http.StatusBadRequest)
	}
	embedCtx := GetContextValue[*web.Context](ctx, WebCtx)

	_, appErr := embedCtx.App.Srv().AttributeService().DeleteAttributes(args.Id)
	if appErr != nil {
		return nil, appErr
	}
	return &AttributeDelete{
		Attribute: SystemAttributeToGraphqlAttribute(&model.Attribute{Id: args.Id}),
	}, nil
}

// NOTE: Refer to ./schemas/attribute.graphqls for details on directive used
func (r *Resolver) AttributeUpdate(ctx context.Context, args struct {
	Id    string
	Input AttributeUpdateInput
}) (*AttributeUpdate, error) {
	// validate params:

	if !model.IsValidId(args.Id) {
		return nil, model.NewAppError("AttributeUpdate", app.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "id"}, "please provide valid attribute id", http.StatusBadRequest)
	}
	if !lo.EveryBy(args.Input.RemoveValues, model.IsValidId) {
		return nil, model.NewAppError("AttributeUpdate", app.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "removeValues"}, "please provide valid attribute value ids", http.StatusBadRequest)
	}

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
			return nil, model.NewAppError("AttributeUpdate", app.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "removeValues"}, "one of attribute values does not belong to given attribute", http.StatusBadRequest)
		}

		// remove attribute values designated:
		_, appErr = embedCtx.App.Srv().AttributeService().DeleteAttributeValues(args.Input.RemoveValues...)
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
	Id           string
	Input        NameTranslationInput
	LanguageCode LanguageCodeEnum
}) (*AttributeTranslate, error) {
	panic(fmt.Errorf("not implemented"))
}

// NOTE: Refer to ./schemas/attribute.graphqls for details on directive used
func (r *Resolver) AttributeBulkDelete(ctx context.Context, args struct{ Ids []string }) (*AttributeBulkDelete, error) {
	// validate params
	if !lo.EveryBy(args.Ids, model.IsValidId) {
		return nil, model.NewAppError("AttributeBulkDelete", app.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "ids"}, "please provide valid ids", http.StatusBadRequest)
	}

	embedCtx := GetContextValue[*web.Context](ctx, WebCtx)
	count, appErr := embedCtx.App.Srv().AttributeService().DeleteAttributes(args.Ids...)
	if appErr != nil {
		return nil, appErr
	}
	return &AttributeBulkDelete{
		Count: int32(count),
	}, nil
}

// NOTE: Refer to ./schemas/attribute.graphqls for details on directive used
func (r *Resolver) AttributeValueBulkDelete(ctx context.Context, args struct{ Ids []string }) (*AttributeValueBulkDelete, error) {
	// valdate params
	if !lo.EveryBy(args.Ids, model.IsValidId) {
		return nil, model.NewAppError("AttributeValueBulkDelete", app.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "ids"}, "please provide valid ids", http.StatusBadRequest)
	}

	embedCtx := GetContextValue[*web.Context](ctx, WebCtx)
	count, appErr := embedCtx.App.Srv().AttributeService().DeleteAttributeValues(args.Ids...)
	if appErr != nil {
		return nil, appErr
	}
	return &AttributeValueBulkDelete{
		Count: int32(count),
	}, nil
}

// NOTE: Refer to ./schemas/attribute.graphqls for details on directive used
func (r *Resolver) AttributeValueCreate(ctx context.Context, args struct {
	AttributeID string
	Input       AttributeValueCreateInput
}) (*AttributeValueCreate, error) {
	// validate params
	if !model.IsValidId(args.AttributeID) {
		return nil, model.NewAppError("AttributeValueCreate", app.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "attributeID"}, "id="+args.AttributeID+" is in valid", http.StatusBadRequest)
	}

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
			return nil, model.NewAppError("AttributeValueCreate", app.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "fileUrl or contentType"}, "fieUrl and contentType can only be defined for swatch attribute", http.StatusBadRequest)
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
func (r *Resolver) AttributeValueDelete(ctx context.Context, args struct{ Id string }) (*AttributeValueDelete, error) {
	// validate params
	if !model.IsValidId(args.Id) {
		return nil, model.NewAppError("AttributeValueDelete", app.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "Id"}, "id="+args.Id+" is in valid", http.StatusBadRequest)
	}

	embedCtx := GetContextValue[*web.Context](ctx, WebCtx)
	attrValues, appErr := embedCtx.App.Srv().AttributeService().FilterAttributeValuesByOptions(model.AttributeValueFilterOptions{
		Conditions: squirrel.Eq{model.AttributeValueTableName + ".Id": args.Id},
	})
	if appErr != nil {
		return nil, appErr
	}
	if len(attrValues) == 0 {
		return nil, model.NewAppError("AttributeValueDelete", app.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "Id"}, "id="+args.Id+" is in valid", http.StatusBadRequest)
	}
	attrValue := attrValues[0]

	_, appErr = embedCtx.App.Srv().AttributeService().DeleteAttributeValues(args.Id)
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
	Id    string
	Input AttributeValueUpdateInput
}) (*AttributeValueUpdate, error) {
	// validate params
	if !model.IsValidId(args.Id) {
		return nil, model.NewAppError("AttributeValueUpdate", app.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "Id"}, "id="+args.Id+" is in valid", http.StatusBadRequest)
	}
	embedCtx := GetContextValue[*web.Context](ctx, WebCtx)

	attrValues, appErr := embedCtx.App.Srv().AttributeService().FilterAttributeValuesByOptions(model.AttributeValueFilterOptions{
		Conditions: squirrel.Eq{model.AttributeValueTableName + ".Id": args.Id},
	})
	if appErr != nil {
		return nil, appErr
	}
	if len(attrValues) == 0 {
		return nil, model.NewAppError("AttributeValueUpdate", app.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "Id"}, "id="+args.Id+" is in valid", http.StatusBadRequest)
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
	Id           string
	Input        AttributeValueTranslationInput
	LanguageCode LanguageCodeEnum
}) (*AttributeValueTranslate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) AttributeReorderValues(ctx context.Context, args struct {
	AttributeID string
	Moves       []*ReorderInput
}) (*AttributeReorderValues, error) {
	// validate params
	if !model.IsValidId(args.AttributeID) {
		return nil, model.NewAppError("AttributeReorderValues", app.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "Id"}, "id="+args.AttributeID+" is in valid", http.StatusBadRequest)
	}

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
			return nil, model.NewAppError("AttributeReorderValues", app.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "move.Id"}, fmt.Sprintf("attribute value with id=%s does not belong to the attribute", move.ID), http.StatusBadRequest)
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

func (r *Resolver) Attributes(ctx context.Context, args struct {
	Filter    *AttributeFilterInput
	SortBy    *AttributeSortingInput
	ChannelID *string
	GraphqlParams
}) (*AttributeCountableConnection, error) {
	// validate params
	var channelID string
	if args.ChannelID != nil {
		if !model.IsValidId(channelID) {
			return nil, model.NewAppError("Attributes", app.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "channel id"}, "please provide valid channel id", http.StatusBadRequest)
		}
	}

	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) Attribute(ctx context.Context, args struct {
	Id   *string
	Slug *string
}) (*Attribute, error) {
	if (args.Id == nil && args.Slug == nil) || (args.Id != nil && args.Slug != nil) {
		return nil, model.NewAppError("Attribute", app.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "id/slug"}, "please provide id or slug, not both", http.StatusBadRequest)
	}
	var attrId string
	if args.Id != nil {
		if !model.IsValidId(attrId) {
			return nil, model.NewAppError("Attribute", app.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "id"}, "please provide valid id", http.StatusBadRequest)
		}
	}
	if args.Slug != nil && !slug.IsSlug(*args.Slug) {
		return nil, model.NewAppError("Attribute", app.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "slug"}, "please provide valid slug", http.StatusBadRequest)
	}

	conditions := squirrel.Eq{}
	if attrId != "" {
		conditions[model.AttributeTableName+".Id"] = *args.Id
	} else {
		conditions[model.AttributeTableName+".Slug"] = *args.Slug
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
