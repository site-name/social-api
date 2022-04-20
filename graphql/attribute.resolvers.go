package graphql

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/gosimple/slug"
	"github.com/graph-gophers/dataloader"
	"github.com/pkg/errors"
	"github.com/sitename/sitename/app"
	"github.com/sitename/sitename/graphql/dataloaders"
	graphql1 "github.com/sitename/sitename/graphql/generated"
	"github.com/sitename/sitename/graphql/gqlmodel"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model/attribute"
	"github.com/sitename/sitename/modules/util"
	"github.com/sitename/sitename/store"
)

func (r *attributeResolver) ProductTypes(ctx context.Context, obj *gqlmodel.Attribute, before *string, after *string, first *int, last *int) (*gqlmodel.ProductTypeCountableConnection, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *attributeResolver) ProductVariantTypes(ctx context.Context, obj *gqlmodel.Attribute, before *string, after *string, first *int, last *int) (*gqlmodel.ProductTypeCountableConnection, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *attributeResolver) Choices(ctx context.Context, obj *gqlmodel.Attribute, sortBy *gqlmodel.AttributeChoicesSortingInput, filter *gqlmodel.AttributeValueFilterInput, before *string, after *string, first *int, last *int) (*gqlmodel.AttributeValueCountableConnection, error) {
	if obj.InputType == nil || !attribute.TYPES_WITH_CHOICES.Contains(strings.ToLower(string(*obj.InputType))) {
		return nil, nil
	}

	var (
		key            string
		orderDirection = gqlmodel.OrderDirectionAsc
		// construct attribute value filter options
		attrValueFilterOptions = &attribute.AttributeValueFilterOptions{
			AttributeID:            squirrel.Eq{store.AttributeValueTableName + ".AttributeID": obj.ID},
			SelectRelatedAttribute: true, //
		}
	)

	// parse order by
	if sortBy != nil {
		switch sortBy.Field {
		case gqlmodel.AttributeChoicesSortFieldName:
			attrValueFilterOptions.OrderBy = store.AttributeValueTableName + ".Name"
		case gqlmodel.AttributeChoicesSortFieldSlug:
			attrValueFilterOptions.OrderBy = store.AttributeValueTableName + ".Slug"
		}

		key = attrValueFilterOptions.OrderBy
		orderDirection = sortBy.Direction
		attrValueFilterOptions.OrderBy += " " + string(sortBy.Direction)
	}

	// parse search value
	if filter != nil && filter.Search != nil {
		attrValueFilterOptions.Extra = squirrel.Or{
			squirrel.ILike{store.AttributeValueTableName + ".Name": "%" + *filter.Search + "%"},
			squirrel.ILike{store.AttributeValueTableName + ".Slug": "%" + *filter.Search + "%"},
		}
	}

	parser := &GraphqlArgumentsParser{
		First:          first,
		Last:           last,
		Before:         before,
		After:          after,
		OrderDirection: orderDirection,
	}
	if appErr := parser.IsValid(); appErr != nil {
		return nil, appErr
	}

	expression, appErr := parser.ConstructSqlExpr(key)
	if appErr != nil {
		return nil, appErr
	}

	attrValueFilterOptions.Extra = squirrel.And{
		attrValueFilterOptions.Extra,
		expression,
	}

	limit := parser.Limit()
	if limit != 0 {
		attrValueFilterOptions.Limit = uint64(limit + 1) // + 1 to check if there is next page available
	}

	// find
	attributeValues, appErr := r.Srv().AttributeService().FilterAttributeValuesByOptions(*attrValueFilterOptions)
	if appErr != nil {
		if appErr.StatusCode == http.StatusInternalServerError {
			return nil, appErr
		}

		// if not found error:
		return &gqlmodel.AttributeValueCountableConnection{
			TotalCount: model.NewInt(0),
		}, nil
	}

	// count
	numOfAttrValues, err := r.Srv().Store.AttributeValue().Count(attrValueFilterOptions)
	if err != nil {
		return nil, model.NewAppError("graphql.attribute.Choices", ".error_counting_attribute_values.app_error", nil, err.Error(), http.StatusInternalServerError)
	}

	result := &gqlmodel.AttributeValueCountableConnection{
		TotalCount: model.NewInt(int(numOfAttrValues)),
	}

	for i := 0; i < util.Min(len(attributeValues), limit); i++ {
		var cursor = attributeValues[i].Name

		if sortBy != nil && sortBy.Field == gqlmodel.AttributeChoicesSortFieldSlug {
			cursor = attributeValues[i].Slug
		}

		result.Edges = append(result.Edges, &gqlmodel.AttributeValueCountableEdge{
			Cursor: util.Base64Encode(cursor),
			Node:   gqlmodel.ModelAttributeValueToGraphqlAttributeValue(attributeValues[i]),
		})
	}

	result.PageInfo = &gqlmodel.PageInfo{
		HasNextPage:     len(attributeValues) > limit,
		HasPreviousPage: parser.HasPreviousPage(),
		StartCursor:     &result.Edges[0].Cursor,
		EndCursor:       &result.Edges[len(result.Edges)-1].Cursor,
	}

	return result, nil
}

func (r *attributeResolver) Translation(ctx context.Context, obj *gqlmodel.Attribute, languageCode gqlmodel.LanguageCodeEnum) (*gqlmodel.AttributeTranslation, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *attributeResolver) WithChoices(ctx context.Context, obj *gqlmodel.Attribute) (bool, error) {
	return obj.InputType != nil && attribute.TYPES_WITH_CHOICES.Contains(strings.ToLower(string(*obj.InputType))), nil
}

func (r *attributeValueResolver) Translation(ctx context.Context, obj *gqlmodel.AttributeValue, languageCode gqlmodel.LanguageCodeEnum) (*gqlmodel.AttributeValueTranslation, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *attributeValueResolver) InputType(ctx context.Context, obj *gqlmodel.AttributeValue) (*gqlmodel.AttributeInputTypeEnum, error) {
	session, appErr := CheckUserAuthenticated("InputType", ctx)
	if appErr != nil {
		return nil, appErr
	}

	// extract data loaders
	thunk := ctx.Value(dataloaders.DataloaderContextKey).(*dataloaders.DataLoaders).AttributeLoader.
		Load(ctx, dataloader.StringKey(obj.AttributeID))

	result, err := thunk()
	if err != nil {
		return nil, errors.Wrap(err, "failed to execute thunk")
	}

	attr := result.(*attribute.Attribute) // type casting
	if attr.Type == attribute.PAGE_TYPE {

		if r.Srv().AccountService().SessionHasPermissionTo(session, model.PermissionManagePages) {
			res := gqlmodel.AttributeInputTypeEnum(strings.ToUpper(attr.InputType))
			return &res, nil
		}

		return nil, model.NewAppError("InputType", PermissionDeniedId, nil, "", http.StatusForbidden)

	} else if r.Srv().AccountService().SessionHasPermissionTo(session, model.PermissionManageProducts) {
		res := gqlmodel.AttributeInputTypeEnum(strings.ToUpper(attr.InputType))
		return &res, nil
	}

	return nil, model.NewAppError("InputType", PermissionDeniedId, nil, "", http.StatusForbidden)
}

func (r *attributeValueResolver) Reference(ctx context.Context, obj *gqlmodel.AttributeValue) (*string, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) AttributeCreate(ctx context.Context, input gqlmodel.AttributeCreateInput) (*gqlmodel.AttributeCreate, error) {
	session, appErr := CheckUserAuthenticated("AttributeCreate", ctx)
	if appErr != nil {
		return nil, appErr
	}

	// check if user has permission to proceed
	var permission *model.Permission
	if input.Type == gqlmodel.AttributeTypeEnumPageType {
		permission = model.PermissionManagePageTypesAndAttributes
	} else {
		permission = model.PermissionManageProductTypesAndAttributes
	}

	if !r.Srv().AccountService().SessionHasPermissionTo(session, permission) {
		return nil, r.Srv().AccountService().MakePermissionError(session, permission)
	}

	// validate if input type & entity type are valid
	if input.InputType != nil &&
		*input.InputType == gqlmodel.AttributeInputTypeEnumReference &&
		(input.EntityType == nil || !(*input.EntityType).IsValid()) {

		return nil, model.NewAppError("AttributeCreate", app.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "input.EntityType"}, "Entity type is required when REFERENCE input type is used", http.StatusBadRequest)
	}

	// clean attribute.
	var slugValue string

	if input.Slug != nil {
		slugValue = *input.Slug
	}
	if slugValue == "" {
		slugValue = slug.Make(input.Name)
	}

	// check if slug is unique,
	// if not, generate a new one
	for {
		_, appErr := r.Srv().AttributeService().AttributeByOption(&attribute.AttributeFilterOption{
			Slug: squirrel.Eq{store.AttributeTableName + ".Slug": slugValue},
		})
		if appErr != nil {
			if appErr.StatusCode == http.StatusInternalServerError {
				return nil, appErr
			}

			// this error means this slug is valid
			break
		}

		slugValue = slugValue + "-" + model.NewId()
	}

	for key, value := range gqlmodel.ATTRIBUTE_PROPERTIES_CONFIGURATION {
		if input.InputType != nil {

			vl := input.GetValueByField(key)
			vlValid := vl != nil

			if vlValid {
				switch t := vl.(type) {
				case *bool:
					vlValid = *t // true
				case *int:
					vlValid = *t != 0
				}
			}

			if !gqlmodel.AttributeInputTypeEnumInSlice(*input.InputType, value...) && vlValid {
				return nil, model.NewAppError("AttributeCreate", app.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "input.InputType"}, fmt.Sprintf("Cannot set %s on a %s", key, *input.InputType), http.StatusBadRequest)
			}
		}
	}

	// clean values
	if len(input.Values) != 0 {

		if input.InputType != nil &&
			gqlmodel.AttributeInputTypeEnumInSlice(*input.InputType, gqlmodel.AttributeInputTypeEnumFile, gqlmodel.AttributeInputTypeEnumReference) &&
			len(input.Values) != 0 {
			return nil, model.NewAppError("AttributeCreate", app.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "input.InputType, input.Values"}, fmt.Sprintf("Values cannot be used with input type %s", *input.InputType), http.StatusBadRequest)
		}

		for _, value := range input.Values {
			if value != nil {
				appErr = r.validateValue(value, *input.InputType == gqlmodel.AttributeInputTypeEnumNumeric, *input.InputType == gqlmodel.AttributeInputTypeEnumSwatch)
				if appErr != nil {
					return nil, appErr
				}
			}
		}
	}

	// construct instance
	attr := &attribute.Attribute{
		Name: input.Name,
		Type: strings.ToLower(string(input.Type)),
		Slug: slugValue,
	}
	if v := input.InputType; v != nil {
		attr.InputType = strings.ToLower(string(*v))
	}
	if v := input.EntityType; v != nil {
		attr.EntityType = model.NewString(strings.ToLower(string(*v)))
	}
	if v := input.Unit; v != nil {
		attr.Unit = model.NewString(strings.ToLower(string(*v)))
	}
	if v := input.ValueRequired; v != nil {
		attr.ValueRequired = *v
	}
	if v := input.IsVariantOnly; v != nil {
		attr.IsVariantOnly = *v
	}
	if v := input.VisibleInStorefront; v != nil {
		attr.VisibleInStoreFront = *v
	}
	if v := input.FilterableInStorefront; v != nil {
		attr.FilterableInStorefront = *v
	}
	if v := input.FilterableInDashboard; v != nil {
		attr.FilterableInDashboard = *v
	}
	if v := input.StorefrontSearchPosition; v != nil {
		attr.StorefrontSearchPosition = *v
	}
	if v := input.AvailableInGrid; v != nil {
		attr.AvailableInGrid = *v
	}

	savedAttr, appErr := r.Srv().AttributeService().UpsertAttribute(attr)
	if appErr != nil {
		return nil, appErr
	}

	// create attribute values if input.Values is provided.
	if len(input.Values) > 0 {
		for _, value := range input.Values {

			var aValue string
			if value.Value != nil {
				aValue = *value.Value
			}

			attrValue := &attribute.AttributeValue{
				AttributeID: savedAttr.Id,

				Name:        value.Name,
				RichText:    value.RichText,
				FileUrl:     value.FileURL,
				ContentType: value.ContentType,
				Value:       aValue,
			}

			_, appErr = r.Srv().AttributeService().UpsertAttributeValue(attrValue)
			if appErr != nil {
				return nil, appErr
			}
		}
	}

	return &gqlmodel.AttributeCreate{
		Attribute: gqlmodel.ModelAttributeToGraphqlAttribute(savedAttr),
	}, nil
}

func (r *mutationResolver) AttributeDelete(ctx context.Context, id string) (*gqlmodel.AttributeDelete, error) {
	session, appErr := CheckUserAuthenticated("AttributeDelete", ctx)
	if appErr != nil {
		return nil, appErr
	}

	if !r.Srv().AccountService().SessionHasPermissionTo(session, model.PermissionManageProductTypesAndAttributes) {
		return nil, r.Srv().AccountService().MakePermissionError(session, model.PermissionManageProductTypesAndAttributes)
	}

	_, appErr = r.Srv().AttributeService().DeleteAttributes(id)
	if appErr != nil {
		return nil, appErr
	}

	return &gqlmodel.AttributeDelete{
		Attribute: &gqlmodel.Attribute{
			ID: id,
		},
	}, nil
}

func (r *mutationResolver) AttributeUpdate(ctx context.Context, id string, input gqlmodel.AttributeUpdateInput) (*gqlmodel.AttributeUpdate, error) {
	session, appErr := CheckUserAuthenticated("AttributeUpdate", ctx)
	if appErr != nil {
		return nil, appErr
	}

	if !r.Srv().AccountService().SessionHasPermissionTo(session, model.PermissionManageProductTypesAndAttributes) {
		return nil, r.Srv().AccountService().MakePermissionError(session, model.PermissionManageProductTypesAndAttributes)
	}

	// check if attribute does exist
	attributeToUpdate, appErr := r.Srv().AttributeService().AttributeByOption(&attribute.AttributeFilterOption{
		Id: squirrel.Eq{store.AttributeTableName + ".Id": id},
	})
	if appErr != nil {
		return nil, appErr
	}

	// clean attribute.
	if input.Slug == nil || (input.Slug != nil && strings.TrimSpace(*input.Slug) == "") {
		return nil, model.NewAppError("AttributeUpdate", app.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "input.Slug"}, "Slug value cannot be blank", http.StatusBadRequest)
	}

	attributeInputType := attributeToUpdate.InputType

	for key, value := range gqlmodel.ATTRIBUTE_PROPERTIES_CONFIGURATION {
		vl := input.GetValueByField(key)
		vlValid := vl != nil

		if vlValid {
			switch t := vl.(type) {
			case *bool:
				vlValid = *t
			case *int:
				vlValid = *t != 0
			}
		}

		if !gqlmodel.AttributeInputTypeEnumInSlice(
			gqlmodel.AttributeInputTypeEnum(strings.ToUpper(attributeInputType)),
			value...,
		) && vlValid {
			return nil, model.NewAppError("AttributeUpdate", app.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "attribute.InputType"}, fmt.Sprintf("Cannot set %s on a %s", key, attributeInputType), http.StatusBadRequest)
		}
	}

	// clean values
	if len(input.AddValues) > 0 {
		if attributeInputType == attribute.FILE ||
			attributeInputType == attribute.REFERENCE {

			return nil, model.NewAppError("AttributeUpdate", app.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "attribute.InputType"}, "Values cannot be used with input type "+attributeInputType, http.StatusBadRequest)
		}

		for _, valueData := range input.AddValues {
			if valueData != nil {
				appErr = r.validateValue(valueData, attributeInputType == attribute.NUMERIC, attributeInputType == attribute.SWATCH)
				if appErr != nil {
					return nil, appErr
				}
			}
		}
	}

	// clean remove values
	attributeValueIDs := util.StringPointerSliceToStringSlice(input.RemoveValues)
	for _, valueID := range attributeValueIDs {
		if model.IsValidId(valueID) {
			attrValues, appErr := r.Srv().AttributeService().FilterAttributeValuesByOptions(attribute.AttributeValueFilterOptions{
				Id: squirrel.Eq{store.AttributeValueTableName + ".Id": valueID},
			})
			if appErr != nil {
				return nil, appErr
			}

			// check if attribute parent of this attribute value is current attribute to update
			if attrValues[0].AttributeID != attributeToUpdate.Id {
				return nil, model.NewAppError("AttributeUpdate", app.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "value.AttributeID"}, "Value does not belong to this attribute", http.StatusBadRequest)
			}
		}
	}

	// construct attribute instance
	if input.Name != nil {
		attributeToUpdate.Name = *input.Name
	}
	if input.Slug != nil {
		attributeToUpdate.Slug = *input.Slug
	}
	if input.Unit != nil {
		attributeToUpdate.Unit = model.NewString(strings.ToLower(string(*input.Unit)))
	}
	if v := input.ValueRequired; v != nil {
		attributeToUpdate.ValueRequired = *v
	}
	if v := input.IsVariantOnly; v != nil {
		attributeToUpdate.IsVariantOnly = *v
	}
	if v := input.VisibleInStorefront; v != nil {
		attributeToUpdate.VisibleInStoreFront = *v
	}
	if v := input.FilterableInStorefront; v != nil {
		attributeToUpdate.FilterableInStorefront = *v
	}
	if v := input.FilterableInDashboard; v != nil {
		attributeToUpdate.FilterableInDashboard = *v
	}
	if v := input.StorefrontSearchPosition; v != nil {
		attributeToUpdate.StorefrontSearchPosition = *v
	}
	if v := input.AvailableInGrid; v != nil {
		attributeToUpdate.AvailableInGrid = *v
	}

	attr, appErr := r.Srv().AttributeService().UpsertAttribute(attributeToUpdate)
	if appErr != nil {
		return nil, appErr
	}

	// remove attribute values
	_, appErr = r.Srv().AttributeService().DeleteAttributeValues(attributeValueIDs...)
	if appErr != nil {
		return nil, appErr
	}

	return &gqlmodel.AttributeUpdate{
		Attribute: gqlmodel.ModelAttributeToGraphqlAttribute(attr),
	}, nil
}

func (r *mutationResolver) AttributeTranslate(ctx context.Context, id string, input gqlmodel.NameTranslationInput, languageCode gqlmodel.LanguageCodeEnum) (*gqlmodel.AttributeTranslate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) AttributeBulkDelete(ctx context.Context, ids []*string) (*gqlmodel.AttributeBulkDelete, error) {
	session, appErr := CheckUserAuthenticated("AttributeBulkDelete", ctx)
	if appErr != nil {
		return nil, appErr
	}
	if !r.Srv().AccountService().SessionHasPermissionTo(session, model.PermissionManagePageTypesAndAttributes) {
		return nil, r.Srv().AccountService().MakePermissionError(session, model.PermissionManagePageTypesAndAttributes)
	}

	var (
		deleteIDs  = util.StringPointerSliceToStringSlice(ids)
		numDeleted int64
	)
	if len(deleteIDs) > 0 {
		numDeleted, appErr = r.Srv().AttributeService().DeleteAttributes(deleteIDs...)
		if appErr != nil {
			return nil, appErr
		}
	}

	return &gqlmodel.AttributeBulkDelete{
		Count: int(numDeleted),
	}, nil
}

func (r *mutationResolver) AttributeValueBulkDelete(ctx context.Context, ids []*string) (*gqlmodel.AttributeValueBulkDelete, error) {
	session, appErr := CheckUserAuthenticated("AttributeBulkDelete", ctx)
	if appErr != nil {
		return nil, appErr
	}
	if !r.Srv().AccountService().SessionHasPermissionTo(session, model.PermissionManagePageTypesAndAttributes) {
		return nil, r.Srv().AccountService().MakePermissionError(session, model.PermissionManagePageTypesAndAttributes)
	}

	var (
		deleteIDs  = util.StringPointerSliceToStringSlice(ids)
		numDeleted int64
	)

	if len(deleteIDs) > 0 {
		numDeleted, appErr = r.Srv().AttributeService().DeleteAttributeValues(deleteIDs...)
		if appErr != nil {
			return nil, appErr
		}
	}

	return &gqlmodel.AttributeValueBulkDelete{
		Count: int(numDeleted),
	}, nil
}

func (r *mutationResolver) AttributeValueCreate(ctx context.Context, attributeID string, input gqlmodel.AttributeValueCreateInput) (*gqlmodel.AttributeValueCreate, error) {
	session, appErr := CheckUserAuthenticated("AttributeValueCreate", ctx)
	if appErr != nil {
		return nil, appErr
	}

	if !r.Srv().AccountService().SessionHasPermissionTo(session, model.PermissionManageProducts) {
		return nil, r.Srv().AccountService().MakePermissionError(session, model.PermissionManageProducts)
	}

	// check if parent attribute exist
	parentAttribute, appErr := r.Srv().AttributeService().AttributeByOption(&attribute.AttributeFilterOption{
		Id: squirrel.Eq{store.AttributeTableName + ".Id": attributeID},
	})
	if appErr != nil {
		return nil, appErr
	}

	slug := slug.Make(input.Name)

	if parentAttribute.InputType != attribute.SWATCH {
		if input.FileURL != nil && *input.FileURL != "" {
			return nil, model.NewAppError("AttributeValueCreate", app.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "input.FileURL"}, "The field FileURL can be defined only for swatch attribute", http.StatusBadRequest)
		}
		if input.ContentType != nil && *input.ContentType != "" {
			return nil, model.NewAppError("AttributeValueCreate", app.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "input.ContentType"}, "The field ContentType can be defined only for swatch attribute", http.StatusBadRequest)
		}
	} else {
		if input.Value != nil && *input.Value != "" && input.FileURL != nil && *input.FileURL != "" {
			return nil, model.NewAppError("AttributeValueCreate", app.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "input.FileURL, input.Value"}, "Cannot specify both value and fileURL for swatch attribute", http.StatusBadRequest)
		}
	}

	// construct and save instance
	attributeValue := &attribute.AttributeValue{
		AttributeID: parentAttribute.Id,
		Slug:        slug,
		Name:        input.Name,
		RichText:    input.RichText,
		ContentType: input.ContentType,
		FileUrl:     input.FileURL,
	}
	if input.Value != nil {
		attributeValue.Value = *input.Value
	}

	attributeValue, appErr = r.Srv().AttributeService().UpsertAttributeValue(attributeValue)
	if appErr != nil {
		return nil, appErr
	}

	// this for ModelAttributeValueToGraphqlAttributeValue()
	// know how to handle Date, DateTime
	attributeValue.Attribute = parentAttribute

	return &gqlmodel.AttributeValueCreate{
		Attribute:      gqlmodel.ModelAttributeToGraphqlAttribute(parentAttribute),
		AttributeValue: gqlmodel.ModelAttributeValueToGraphqlAttributeValue(attributeValue),
	}, nil
}

func (r *mutationResolver) AttributeValueDelete(ctx context.Context, id string) (*gqlmodel.AttributeValueDelete, error) {
	session, appErr := CheckUserAuthenticated("AttributeValueDelete", ctx)
	if appErr != nil {
		return nil, appErr
	}

	if !r.Srv().AccountService().SessionHasPermissionTo(session, model.PermissionManageProductTypesAndAttributes) {
		return nil, r.Srv().AccountService().MakePermissionError(session, model.PermissionManageProductTypesAndAttributes)
	}

	_, appErr = r.Srv().AttributeService().DeleteAttributeValues(id)
	if appErr != nil {
		return nil, appErr
	}

	return &gqlmodel.AttributeValueDelete{
		AttributeValue: &gqlmodel.AttributeValue{
			ID: id,
		},
	}, nil
}

func (r *mutationResolver) AttributeValueUpdate(ctx context.Context, id string, input gqlmodel.AttributeValueUpdateInput) (*gqlmodel.AttributeValueUpdate, error) {
	session, appErr := CheckUserAuthenticated("AttributeValueUpdate", ctx)
	if appErr != nil {
		return nil, appErr
	}
	if !r.Srv().AccountService().SessionHasPermissionTo(session, model.PermissionManageProductTypesAndAttributes) {
		return nil, r.Srv().AccountService().MakePermissionError(session, model.PermissionManageProductTypesAndAttributes)
	}

	// check if attribute value with id = given id does exist:
	attributeValues, appErr := r.Srv().AttributeService().FilterAttributeValuesByOptions(attribute.AttributeValueFilterOptions{
		Id:                     squirrel.Eq{store.AttributeValueTableName + ".Id": id},
		SelectRelatedAttribute: true,
	})
	if appErr != nil {
		return nil, appErr
	}

	var (
		attributeValueToUpdate = attributeValues[0]
		parentAttribute        = attributeValueToUpdate.Attribute
	)

	// clean input
	var slugValue string
	if input.Name != nil {
		slugValue = slug.Make(*input.Name)
	}

	if attributeValueToUpdate.Attribute.InputType != attribute.SWATCH {
		if input.FileURL != nil && *input.FileURL != "" {
			return nil, model.NewAppError("AttributeValueUpdate", app.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "input.FileURL"}, "The field FileURL can be defined only for swatch attributes.", http.StatusBadRequest)
		}
		if input.ContentType != nil && *input.ContentType != "" {
			return nil, model.NewAppError("AttributeValueUpdate", app.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "input.ContentType"}, "The field ContentType can be defined only for swatch attributes.", http.StatusBadRequest)
		}
	} else {
		if input.Value != nil && *input.Value != "" && input.FileURL != nil && *input.FileURL != "" {
			return nil, model.NewAppError("AttributeValueUpdate", app.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "input.Value, input.FileURL"}, "Cannot specify both value and file for swatch attribute.", http.StatusBadRequest)
		}
	}

	if input.Value != nil && *input.Value != "" {
		input.FileURL = nil
		input.ContentType = nil
	} else if input.FileURL != nil && *input.FileURL != "" {
		input.Value = nil
	}

	// construct instance
	if v := input.Value; v != nil {
		attributeValueToUpdate.Value = *v
	}
	if v := input.RichText; len(v) > 0 {
		attributeValueToUpdate.RichText = v
	}
	if v := input.FileURL; v != nil {
		attributeValueToUpdate.FileUrl = v
	}
	if v := input.ContentType; v != nil {
		attributeValueToUpdate.ContentType = v
	}
	if v := input.Name; v != nil {
		attributeValueToUpdate.Name = *v
	}
	attributeValueToUpdate.Slug = slugValue

	attrValue, appErr := r.Srv().AttributeService().UpsertAttributeValue(attributeValueToUpdate)
	if appErr != nil {
		return nil, appErr
	}

	// for ModelAttributeValueToGraphqlAttributeValue() to know
	// how to resolve Date and DateTime
	attrValue.Attribute = parentAttribute

	return &gqlmodel.AttributeValueUpdate{
		Attribute:      gqlmodel.ModelAttributeToGraphqlAttribute(attributeValueToUpdate.Attribute),
		AttributeValue: gqlmodel.ModelAttributeValueToGraphqlAttributeValue(attrValue),
	}, nil
}

func (r *mutationResolver) AttributeValueTranslate(ctx context.Context, id string, input gqlmodel.AttributeValueTranslationInput, languageCode gqlmodel.LanguageCodeEnum) (*gqlmodel.AttributeValueTranslate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) AttributeReorderValues(ctx context.Context, attributeID string, moves []*gqlmodel.ReorderInput) (*gqlmodel.AttributeReorderValues, error) {
	session, appErr := CheckUserAuthenticated("AttributeReorderValues", ctx)
	if appErr != nil {
		return nil, appErr
	}

	if !r.Srv().AccountService().SessionHasPermissionTo(session, model.PermissionManageProductTypesAndAttributes) {
		return nil, r.Srv().AccountService().MakePermissionError(session, model.PermissionManageProductTypesAndAttributes)
	}

	// check if an attribute with given id does really exist
	attr, appErr := r.Srv().AttributeService().AttributeByOption(&attribute.AttributeFilterOption{
		Id:                             squirrel.Eq{store.AttributeTableName + ".Id": attributeID},
		PrefetchRelatedAttributeValues: true, //
	})
	if appErr != nil {
		return nil, appErr
	}

	// operations has keys are attribute value ids
	var operations = map[string]*int{}
	for _, value := range attr.AttributeValues {
		if value != nil {
			operations[value.Id] = nil
		}
	}

	for _, moveInfo := range moves {
		_, exist := operations[moveInfo.ID]
		if !exist {
			return nil, model.NewAppError("AttributeReorderValues", "graphql.attribute.error_resolving_an_attribute_value.app_error", nil, "Couldn't resolve to an attribute value: Id="+moveInfo.ID, http.StatusNotFound)
		}

		operations[moveInfo.ID] = moveInfo.SortOrder
	}

	appErr = r.Srv().AttributeService().PerformReordering(attr.AttributeValues, operations)
	if appErr != nil {
		return nil, appErr
	}

	// TODO: consider if we need to refetching attribute value from database
	return &gqlmodel.AttributeReorderValues{
		Attribute: gqlmodel.ModelAttributeToGraphqlAttribute(attr),
	}, nil
}

func (r *queryResolver) Attributes(ctx context.Context, filter *gqlmodel.AttributeFilterInput, sortBy *gqlmodel.AttributeSortingInput, chanelSlug *string, before *string, after *string, first *int, last *int) (*gqlmodel.AttributeCountableConnection, error) {
	var (
		session, _     = CheckUserAuthenticated("Attributes", ctx)
		orderDirection = gqlmodel.OrderDirectionAsc // default to "ASC"
		// key                    = ""
		attributeFilterOptions = &attribute.AttributeFilterOption{Distinct: true}
	)

	// if user not authenticated or
	// authenticated but does not have specific permission(s),
	// then show only visible to store attributes
	if session == nil ||
		(session != nil && !r.Srv().AccountService().SessionHasPermissionToAny(
			session,
			model.PermissionManagePageTypesAndAttributes,
			model.PermissionManageProductTypesAndAttributes,
		)) {
		attributeFilterOptions.VisibleInStoreFront = model.NewBool(true)
	}

	if sortBy != nil {
		orderDirection = sortBy.Direction

		switch sortBy.Field {
		case gqlmodel.AttributeSortFieldAvailableInGrid:
			attributeFilterOptions.OrderBy = store.AttributeTableName + ".AvailableInGrid"
		case gqlmodel.AttributeSortFieldFilterableInDashboard:
			attributeFilterOptions.OrderBy = store.AttributeTableName + ".FilterableInDashboard"
		case gqlmodel.AttributeSortFieldFilterableInStorefront:
			attributeFilterOptions.OrderBy = store.AttributeTableName + ".FilterableInDashboard"
		case gqlmodel.AttributeSortFieldIsVariantOnly:
			attributeFilterOptions.OrderBy = store.AttributeTableName + ".IsVariantOnly"
		case gqlmodel.AttributeSortFieldName:
			attributeFilterOptions.OrderBy = store.AttributeTableName + ".Name"
		case gqlmodel.AttributeSortFieldSlug:
			attributeFilterOptions.OrderBy = store.AttributeTableName + ".Slug"
		case gqlmodel.AttributeSortFieldStorefrontSearchPosition:
			attributeFilterOptions.OrderBy = store.AttributeTableName + ".StorefrontSearchPosition"
		case gqlmodel.AttributeSortFieldValueRequired:
			attributeFilterOptions.OrderBy = store.AttributeTableName + ".ValueRequired"
		case gqlmodel.AttributeSortFieldVisibleInStorefront:
			attributeFilterOptions.OrderBy = store.AttributeTableName + ".VisibleInStoreFront"
		}

		// key = attributeFilterOptions.OrderBy //
		attributeFilterOptions.OrderBy += " " + string(orderDirection)
	}

	if filter != nil {
		attributeFilterOptions.ValueRequired = filter.ValueRequired
		attributeFilterOptions.IsVariantOnly = filter.IsVariantOnly
		attributeFilterOptions.VisibleInStoreFront = filter.VisibleInStorefront
		attributeFilterOptions.FilterableInStorefront = filter.FilterableInStorefront
		attributeFilterOptions.FilterableInDashboard = filter.FilterableInDashboard
		attributeFilterOptions.AvailableInGrid = filter.AvailableInGrid

		ids := util.StringPointerSliceToStringSlice(filter.Ids)
		if len(ids) > 0 {
			attributeFilterOptions.Id = squirrel.Eq{store.AttributeTableName + ".Id": ids}
		}

		if filter.Type != nil {
			attributeFilterOptions.Type = squirrel.Eq{store.AttributeTableName + ".Type": strings.ToLower(string(*filter.Type))}
		}

		if filter.Search != nil {
			attributeFilterOptions.Extra = squirrel.Or{
				squirrel.ILike{store.AttributeTableName + ".Slug": "%" + *filter.Search + "%"},
				squirrel.ILike{store.AttributeTableName + ".Name": "%" + *filter.Search + "%"},
			}
		}
	}

	// check graphql arguments
	parser := &GraphqlArgumentsParser{
		First:          first,
		Last:           last,
		Before:         before,
		After:          after,
		OrderDirection: orderDirection,
	}
	if appErr := parser.IsValid(); appErr != nil {
		return nil, appErr
	}

	panic("not implt")
}

func (r *queryResolver) Attribute(ctx context.Context, id *string, slug *string) (*gqlmodel.Attribute, error) {
	option := &attribute.AttributeFilterOption{}

	if id != nil && model.IsValidId(*id) {
		option.Id = squirrel.Eq{store.AttributeTableName + ".Id": *id}
	}
	if slug != nil && *slug != "" {
		option.Slug = squirrel.Eq{store.AttributeTableName + ".Slug": *slug}
	}

	attr, appErr := r.Srv().AttributeService().AttributeByOption(option)
	if appErr != nil {
		return nil, appErr
	}

	return gqlmodel.ModelAttributeToGraphqlAttribute(attr), nil
}

// Attribute returns graphql1.AttributeResolver implementation.
func (r *Resolver) Attribute() graphql1.AttributeResolver { return &attributeResolver{r} }

// AttributeValue returns graphql1.AttributeValueResolver implementation.
func (r *Resolver) AttributeValue() graphql1.AttributeValueResolver {
	return &attributeValueResolver{r}
}

type attributeResolver struct{ *Resolver }
type attributeValueResolver struct{ *Resolver }

// !!! WARNING !!!
// The code below was going to be deleted when updating resolvers. It has been copied here so you have
// one last chance to move it out of harms way if you want. There are two reasons this happens:
//  - When renaming or deleting a resolver the old code will be put in here. You can safely delete
//    it when you're done.
//  - You have helper methods in this file. Move them out to keep these resolver files clean.
func (r *attributeValueResolver) Date(ctx context.Context, obj *gqlmodel.AttributeValue) (*time.Time, error) {
	panic(fmt.Errorf("not implemented"))
}
func (r *attributeValueResolver) DateTime(ctx context.Context, obj *gqlmodel.AttributeValue) (*time.Time, error) {
	panic(fmt.Errorf("not implemented"))
}
