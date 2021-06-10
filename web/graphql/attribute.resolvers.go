package graphql

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"

	"github.com/sitename/sitename/web/graphql/gqlmodel"
)

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

func (r *mutationResolver) AttributeValueUpdate(ctx context.Context, id string, input gqlmodel.AttributeValueCreateInput) (*gqlmodel.AttributeValueUpdate, error) {
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
	panic(fmt.Errorf("not implemented"))
}
