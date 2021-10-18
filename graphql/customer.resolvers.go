package graphql

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"

	"github.com/sitename/sitename/graphql/gqlmodel"
)

func (r *mutationResolver) CustomerCreate(ctx context.Context, input gqlmodel.UserCreateInput) (*gqlmodel.CustomerCreate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) CustomerUpdate(ctx context.Context, id string, input gqlmodel.CustomerInput) (*gqlmodel.CustomerUpdate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) CustomerDelete(ctx context.Context, id string) (*gqlmodel.CustomerDelete, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) CustomerBulkDelete(ctx context.Context, ids []*string) (*gqlmodel.CustomerBulkDelete, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) Customers(ctx context.Context, filter *gqlmodel.CustomerFilterInput, sortBy *gqlmodel.UserSortingInput, before *string, after *string, first *int, last *int) (*gqlmodel.UserCountableConnection, error) {
	panic(fmt.Errorf("not implemented"))
}
