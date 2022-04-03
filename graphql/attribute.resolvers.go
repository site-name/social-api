package graphql

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/gosimple/slug"
	"github.com/sitename/sitename/app"
	graphql1 "github.com/sitename/sitename/graphql/generated"
	"github.com/sitename/sitename/graphql/gqlmodel"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model/attribute"
)

func (r *attributeResolver) ProductTypes(ctx context.Context, obj *gqlmodel.Attribute, before *string, after *string, first *int, last *int) (*gqlmodel.ProductTypeCountableConnection, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *attributeResolver) ProductVariantTypes(ctx context.Context, obj *gqlmodel.Attribute, before *string, after *string, first *int, last *int) (*gqlmodel.ProductTypeCountableConnection, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *attributeResolver) Choices(ctx context.Context, obj *gqlmodel.Attribute, sortBy *gqlmodel.AttributeChoicesSortingInput, filter *gqlmodel.AttributeValueFilterInput, before *string, after *string, first *int, last *int) (*gqlmodel.AttributeValueCountableConnection, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *attributeResolver) Translation(ctx context.Context, obj *gqlmodel.Attribute, languageCode gqlmodel.LanguageCodeEnum) (*gqlmodel.AttributeTranslation, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *attributeResolver) WithChoices(ctx context.Context, obj *gqlmodel.Attribute) (bool, error) {
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

	// validate input type & entity type are valid
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
		_, appErr := r.Srv().AttributeService().AttributeBySlug(slugValue)
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
					vlValid = *t
				case *int:
					vlValid = *t != 0
				}
			}

			if !gqlmodel.AttributeInputTypeEnumInSlice(*input.InputType, value...) && vlValid {
				return nil, model.NewAppError("AttributeCreate", app.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "input.InputType"}, "", http.StatusBadRequest)
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

	appErr = r.Srv().AttributeService().DeleteAttribute(id)
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
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) AttributeTranslate(ctx context.Context, id string, input gqlmodel.NameTranslationInput, languageCode gqlmodel.LanguageCodeEnum) (*gqlmodel.AttributeTranslate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) AttributeBulkDelete(ctx context.Context, ids []*string) (*gqlmodel.AttributeBulkDelete, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) AttributeValueBulkDelete(ctx context.Context, ids []*string) (*gqlmodel.AttributeValueBulkDelete, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) AttributeValueCreate(ctx context.Context, attribute string, input gqlmodel.AttributeValueCreateInput) (*gqlmodel.AttributeValueCreate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) AttributeValueDelete(ctx context.Context, id string) (*gqlmodel.AttributeValueDelete, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) AttributeValueUpdate(ctx context.Context, id string, input gqlmodel.AttributeValueUpdateInput) (*gqlmodel.AttributeValueUpdate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) AttributeValueTranslate(ctx context.Context, id string, input gqlmodel.AttributeValueTranslationInput, languageCode gqlmodel.LanguageCodeEnum) (*gqlmodel.AttributeValueTranslate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) AttributeReorderValues(ctx context.Context, attributeID string, moves []*gqlmodel.ReorderInput) (*gqlmodel.AttributeReorderValues, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) Attributes(ctx context.Context, filter *gqlmodel.AttributeFilterInput, sortBy *gqlmodel.AttributeSortingInput, before *string, after *string, first *int, last *int) (*gqlmodel.AttributeCountableConnection, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) Attribute(ctx context.Context, id *string, slug *string) (*gqlmodel.Attribute, error) {
	var (
		attr   *attribute.Attribute
		appErr *model.AppError
	)

	// check if `id` is provided correctly:
	if id != nil && model.IsValidId(*id) {
		attr, appErr = r.Srv().AttributeService().AttributeByID(*id)
	} else if slug != nil {
		attr, appErr = r.Srv().AttributeService().AttributeBySlug(*slug)
	}

	if appErr != nil {
		return nil, appErr
	}
	return gqlmodel.ModelAttributeToGraphqlAttribute(attr), nil
}

// Attribute returns graphql1.AttributeResolver implementation.
func (r *Resolver) Attribute() graphql1.AttributeResolver { return &attributeResolver{r} }

type attributeResolver struct{ *Resolver }
