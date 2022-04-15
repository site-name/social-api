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
	"github.com/sitename/sitename/app"
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

	// construct attribute value filter options
	attrValueFilterOptions := &attribute.AttributeValueFilterOptions{
		AttributeID: squirrel.Eq{store.AttributeValueTableName + ".AttributeID": obj.ID},
	}

	var (
		key            string
		orderDirection = gqlmodel.OrderDirectionAsc
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
			squirrel.ILike{store.AttributeValueTableName + ".Name": *filter.Search},
			squirrel.ILike{store.AttributeValueTableName + ".Slug": *filter.Search},
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
		var cursor string

		switch sortBy.Field {
		case gqlmodel.AttributeChoicesSortFieldName:
			cursor = attributeValues[i].Name
		case gqlmodel.AttributeChoicesSortFieldSlug:
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
	panic(fmt.Errorf("not implemented"))
}

func (r *attributeValueResolver) Reference(ctx context.Context, obj *gqlmodel.AttributeValue) (*string, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *attributeValueResolver) Date(ctx context.Context, obj *gqlmodel.AttributeValue) (*time.Time, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *attributeValueResolver) DateTime(ctx context.Context, obj *gqlmodel.AttributeValue) (*time.Time, error) {
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
				appErr = r.validateValue(*value, *input.InputType == gqlmodel.AttributeInputTypeEnumNumeric, *input.InputType == gqlmodel.AttributeInputTypeEnumSwatch)
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
	if input.InputType != nil {
		attr.InputType = strings.ToLower(string(*input.InputType))
	}
	if input.EntityType != nil {
		attr.EntityType = model.NewString(strings.ToLower(string(*input.EntityType)))
	}
	if input.Unit != nil {
		attr.Unit = model.NewString(strings.ToLower(string(*input.Unit)))
	}
	if input.ValueRequired != nil {
		attr.ValueRequired = *input.ValueRequired
	}
	if input.IsVariantOnly != nil {
		attr.IsVariantOnly = *input.IsVariantOnly
	}
	if input.VisibleInStorefront != nil {
		attr.VisibleInStoreFront = *input.VisibleInStorefront
	}
	if input.FilterableInStorefront != nil {
		attr.FilterableInStorefront = *input.FilterableInStorefront
	}
	if input.FilterableInDashboard != nil {
		attr.FilterableInDashboard = *input.FilterableInDashboard
	}
	if input.StorefrontSearchPosition != nil {
		attr.StorefrontSearchPosition = *input.StorefrontSearchPosition
	}
	if input.AvailableInGrid != nil {
		attr.AvailableInGrid = *input.AvailableInGrid
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
	if input.Slug != nil && strings.TrimSpace(*input.Slug) == "" {
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

	if len(input.AddValues) != 0 {

		if attributeInputType == attribute.FILE ||
			attributeInputType == attribute.REFERENCE {

			return nil, model.NewAppError("AttributeUpdate", app.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "attribute.InputType"}, "Values cannot be used with input type "+attributeInputType, http.StatusBadRequest)
		}

		// for _, valueData := range input.AddValues {

		// }
	}

	panic("not implemented")
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
	panic(fmt.Errorf("not implemented"))
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
		session, _             = CheckUserAuthenticated("Attributes", ctx)
		orderDirection         = gqlmodel.OrderDirectionAsc // default to "ASC"
		key                    string
		attributeFilterOptions = &attribute.AttributeFilterOption{Distinct: true}
	)

	// if user not authenticated or
	// authenticated but does not has specific permission(s),
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

		key = attributeFilterOptions.OrderBy //
		attributeFilterOptions.OrderBy += " " + string(orderDirection)
	}

	if filter != nil {
		attributeFilterOptions.ValueRequired = filter.ValueRequired
		attributeFilterOptions.IsVariantOnly = filter.IsVariantOnly
		attributeFilterOptions.VisibleInStoreFront = filter.VisibleInStorefront
		attributeFilterOptions.FilterableInStorefront = filter.FilterableInStorefront
		attributeFilterOptions.FilterableInDashboard = filter.FilterableInDashboard
		attributeFilterOptions.AvailableInGrid = filter.AvailableInGrid
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
