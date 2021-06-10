package graphql

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"
)

func (r *mutationResolver) PageTypeCreate(ctx context.Context, input PageTypeCreateInput) (*PageTypeCreate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) PageTypeUpdate(ctx context.Context, id *string, input PageTypeUpdateInput) (*PageTypeUpdate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) PageTypeDelete(ctx context.Context, id string) (*PageTypeDelete, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) PageTypeBulkDelete(ctx context.Context, ids []string) (*PageTypeBulkDelete, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) PageTypeReorderAttributes(ctx context.Context, moves []ReorderInput, pageTypeID string) (*PageTypeReorderAttributes, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) PageType(ctx context.Context, id string) (*PageType, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) PageTypes(ctx context.Context, sortBy *PageTypeSortingInput, filter *PageTypeFilterInput, before *string, after *string, first *int, last *int) (*PageTypeCountableConnection, error) {
	panic(fmt.Errorf("not implemented"))
}
