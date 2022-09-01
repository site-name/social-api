package api

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"

	"github.com/sitename/sitename/api/gqlmodel"
)

func (r *Resolver) ProductTypeCreate(ctx context.Context, args struct{ Input gqlmodel.ProductTypeInput }) (*gqlmodel.ProductTypeCreate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) ProductTypeDelete(ctx context.Context, args struct{ Id string }) (*gqlmodel.ProductTypeDelete, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) ProductTypeBulkDelete(ctx context.Context, args struct{ Ids []*string }) (*gqlmodel.ProductTypeBulkDelete, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) ProductTypeUpdate(ctx context.Context, args struct {
	Id    string
	Input gqlmodel.ProductTypeInput
}) (*gqlmodel.ProductTypeUpdate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) ProductTypeReorderAttributes(ctx context.Context, args struct {
	Moves         []*gqlmodel.ReorderInput
	ProductTypeID string
	TypeArg       gqlmodel.ProductAttributeType
}) (*gqlmodel.ProductTypeReorderAttributes, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) ProductType(ctx context.Context, args struct{ Id string }) (*gqlmodel.ProductType, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) ProductTypes(ctx context.Context, args struct {
	Filter *gqlmodel.ProductTypeFilterInput
	SortBy *gqlmodel.ProductTypeSortingInput
	Before *string
	After  *string
	First  *int
	Last   *int
}) (*gqlmodel.ProductTypeCountableConnection, error) {
	panic(fmt.Errorf("not implemented"))
}
