package graphql

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"
)

func (r *mutationResolver) ProductTypeCreate(ctx context.Context, input ProductTypeInput) (*ProductTypeCreate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) ProductTypeDelete(ctx context.Context, id string) (*ProductTypeDelete, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) ProductTypeBulkDelete(ctx context.Context, ids []*string) (*ProductTypeBulkDelete, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) ProductTypeUpdate(ctx context.Context, id string, input ProductTypeInput) (*ProductTypeUpdate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) ProductTypeReorderAttributes(ctx context.Context, moves []*ReorderInput, productTypeID string, typeArg ProductAttributeType) (*ProductTypeReorderAttributes, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) ProductType(ctx context.Context, id string) (*ProductType, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) ProductTypes(ctx context.Context, filter *ProductTypeFilterInput, sortBy *ProductTypeSortingInput, before *string, after *string, first *int, last *int) (*ProductTypeCountableConnection, error) {
	panic(fmt.Errorf("not implemented"))
}
