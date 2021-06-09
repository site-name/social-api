package graphql

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"
)

func (r *mutationResolver) ProductAttributeAssign(ctx context.Context, operations []*ProductAttributeAssignInput, productTypeID string) (*ProductAttributeAssign, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) ProductAttributeUnassign(ctx context.Context, attributeIds []*string, productTypeID string) (*ProductAttributeUnassign, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) ProductCreate(ctx context.Context, input ProductCreateInput) (*ProductCreate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) ProductDelete(ctx context.Context, id string) (*ProductDelete, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) ProductBulkDelete(ctx context.Context, ids []*string) (*ProductBulkDelete, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) ProductUpdate(ctx context.Context, id string, input ProductInput) (*ProductUpdate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) ProductTranslate(ctx context.Context, id string, input TranslationInput, languageCode LanguageCodeEnum) (*ProductTranslate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) Product(ctx context.Context, id *string, slug *string, channel *string) (*Product, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) Products(ctx context.Context, filter *ProductFilterInput, sortBy *ProductOrder, channel *string, before *string, after *string, first *int, last *int) (*ProductCountableConnection, error) {
	panic(fmt.Errorf("not implemented"))
}
