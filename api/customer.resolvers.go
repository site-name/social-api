package api

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"
)

func (r *Resolver) CustomerCreate(ctx context.Context, args struct{ Input UserCreateInput }) (*CustomerCreate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) CustomerUpdate(ctx context.Context, args struct {
	Id    string
	Input CustomerInput
}) (*CustomerUpdate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) CustomerDelete(ctx context.Context, args struct{ Id string }) (*CustomerDelete, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) CustomerBulkDelete(ctx context.Context, args struct{ Ids []string }) (*CustomerBulkDelete, error) {
	panic(fmt.Errorf("not implemented"))
}

// NOTE: graphql directive(s) validated before this. Refer to ./schemas/customer.graphqls for details
func (r *Resolver) Customers(ctx context.Context, args struct {
	Filter *CustomerFilterInput
	SortBy *UserSortingInput
	GraphqlParams
}) (*UserCountableConnection, error) {
	panic(fmt.Errorf("not implemented"))
}
