package api

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"

	"github.com/sitename/sitename/api/gqlmodel"
)

func (r *Resolver) CustomerCreate(ctx context.Context, args struct{ input gqlmodel.UserCreateInput }) (*gqlmodel.CustomerCreate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) CustomerUpdate(ctx context.Context, args struct {
	id    string
	input gqlmodel.CustomerInput
}) (*gqlmodel.CustomerUpdate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) CustomerDelete(ctx context.Context, args struct{ id string }) (*gqlmodel.CustomerDelete, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) CustomerBulkDelete(ctx context.Context, args struct{ ids []*string }) (*gqlmodel.CustomerBulkDelete, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) Customers(ctx context.Context, args struct {
	filter *gqlmodel.CustomerFilterInput
	sortBy *gqlmodel.UserSortingInput
	before *string
	after  *string
	first  *int
	last   *int
}) (*gqlmodel.UserCountableConnection, error) {
	panic(fmt.Errorf("not implemented"))
}
