package api

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"

	"github.com/sitename/sitename/api/gqlmodel"
)

func (r *Resolver) ProductTypeCreate(ctx context.Context, input gqlmodel.ProductTypeInput) (*gqlmodel.ProductTypeCreate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) ProductTypeDelete(ctx context.Context, id string) (*gqlmodel.ProductTypeDelete, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) ProductTypeBulkDelete(ctx context.Context, ids []*string) (*gqlmodel.ProductTypeBulkDelete, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) ProductTypeUpdate(ctx context.Context, id string, input gqlmodel.ProductTypeInput) (*gqlmodel.ProductTypeUpdate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) ProductTypeReorderAttributes(ctx context.Context, moves []*gqlmodel.ReorderInput, productTypeID string, typeArg gqlmodel.ProductAttributeType) (*gqlmodel.ProductTypeReorderAttributes, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) ProductType(ctx context.Context, id string) (*gqlmodel.ProductType, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) ProductTypes(ctx context.Context, filter *gqlmodel.ProductTypeFilterInput, sortBy *gqlmodel.ProductTypeSortingInput, before *string, after *string, first *int, last *int) (*gqlmodel.ProductTypeCountableConnection, error) {
	panic(fmt.Errorf("not implemented"))
}
