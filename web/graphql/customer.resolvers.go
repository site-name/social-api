package graphql

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"
)

func (r *mutationResolver) CustomerCreate(ctx context.Context, input UserCreateInput) (*CustomerCreate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) CustomerUpdate(ctx context.Context, id string, input CustomerInput) (*CustomerUpdate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) CustomerDelete(ctx context.Context, id string) (*CustomerDelete, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) CustomerBulkDelete(ctx context.Context, ids []*string) (*CustomerBulkDelete, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) Customers(ctx context.Context, filter *CustomerFilterInput, sortBy *UserSortingInput, before *string, after *string, first *int, last *int) (*UserCountableConnection, error) {
	panic(fmt.Errorf("not implemented"))
}
