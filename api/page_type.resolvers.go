package api

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"

	"github.com/sitename/sitename/api/gqlmodel"
)

func (r *Resolver) PageTypeCreate(ctx context.Context, input gqlmodel.PageTypeCreateInput) (*gqlmodel.PageTypeCreate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) PageTypeUpdate(ctx context.Context, id *string, input gqlmodel.PageTypeUpdateInput) (*gqlmodel.PageTypeUpdate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) PageTypeDelete(ctx context.Context, id string) (*gqlmodel.PageTypeDelete, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) PageTypeBulkDelete(ctx context.Context, ids []string) (*gqlmodel.PageTypeBulkDelete, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) PageTypeReorderAttributes(ctx context.Context, moves []gqlmodel.ReorderInput, pageTypeID string) (*gqlmodel.PageTypeReorderAttributes, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) PageType(ctx context.Context, id string) (*gqlmodel.PageType, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) PageTypes(ctx context.Context, sortBy *gqlmodel.PageTypeSortingInput, filter *gqlmodel.PageTypeFilterInput, before *string, after *string, first *int, last *int) (*gqlmodel.PageTypeCountableConnection, error) {
	panic(fmt.Errorf("not implemented"))
}
