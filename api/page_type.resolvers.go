package api

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"

	"github.com/sitename/sitename/api/gqlmodel"
)

func (r *Resolver) PageTypeCreate(ctx context.Context, args struct{ Input gqlmodel.PageTypeCreateInput }) (*gqlmodel.PageTypeCreate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) PageTypeUpdate(ctx context.Context, args struct {
	Id    *string
	Input gqlmodel.PageTypeUpdateInput
}) (*gqlmodel.PageTypeUpdate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) PageTypeDelete(ctx context.Context, args struct{ Id string }) (*gqlmodel.PageTypeDelete, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) PageTypeBulkDelete(ctx context.Context, args struct{ Ids []string }) (*gqlmodel.PageTypeBulkDelete, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) PageTypeReorderAttributes(ctx context.Context, args struct {
	Moves      []gqlmodel.ReorderInput
	PageTypeID string
}) (*gqlmodel.PageTypeReorderAttributes, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) PageType(ctx context.Context, args struct{ Id string }) (*gqlmodel.PageType, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) PageTypes(ctx context.Context, args struct {
	SortBy *gqlmodel.PageTypeSortingInput
	Filter *gqlmodel.PageTypeFilterInput
	Before *string
	After  *string
	First  *int
	Last   *int
}) (*gqlmodel.PageTypeCountableConnection, error) {
	panic(fmt.Errorf("not implemented"))
}
