package api

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"
)

func (r *Resolver) ProductTypeCreate(ctx context.Context, args struct{ Input ProductTypeInput }) (*ProductTypeCreate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) ProductTypeDelete(ctx context.Context, args struct{ Id string }) (*ProductTypeDelete, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) ProductTypeBulkDelete(ctx context.Context, args struct{ Ids []string }) (*ProductTypeBulkDelete, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) ProductTypeUpdate(ctx context.Context, args struct {
	Id    string
	Input ProductTypeInput
}) (*ProductTypeUpdate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) ProductTypeReorderAttributes(ctx context.Context, args struct {
	Moves         []*ReorderInput
	ProductTypeID string
	Type          ProductAttributeType
}) (*ProductTypeReorderAttributes, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) ProductType(ctx context.Context, args struct{ Id string }) (*ProductType, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) ProductTypes(ctx context.Context, args struct {
	Filter *ProductTypeFilterInput
	SortBy *ProductTypeSortingInput
	GraphqlParams
}) (*ProductTypeCountableConnection, error) {
	panic(fmt.Errorf("not implemented"))
}
