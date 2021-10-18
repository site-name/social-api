package graphql

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"
	"net/http"

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
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) AttributeDelete(ctx context.Context, id string) (*gqlmodel.AttributeDelete, error) {
	panic(fmt.Errorf("not implemented"))
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
	// validate if either arguments are provided
	if id == nil && slug == nil {
		return nil, model.NewAppError("Attribute", app.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "'id', 'slug'"}, "", http.StatusBadRequest)
	}

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
