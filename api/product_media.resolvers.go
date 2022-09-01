package api

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"

	"github.com/sitename/sitename/api/gqlmodel"
)

func (r *Resolver) ProductMediaCreate(ctx context.Context, args struct {
	input gqlmodel.ProductMediaCreateInput
}) (*gqlmodel.ProductMediaCreate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) ProductMediaDelete(ctx context.Context, args struct{ id string }) (*gqlmodel.ProductMediaDelete, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) ProductMediaBulkDelete(ctx context.Context, args struct{ ids []*string }) (*gqlmodel.ProductMediaBulkDelete, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) ProductMediaReorder(ctx context.Context, args struct {
	mediaIds  []*string
	productID string
}) (*gqlmodel.ProductMediaReorder, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) ProductMediaUpdate(ctx context.Context, args struct {
	id    string
	input gqlmodel.ProductMediaUpdateInput
}) (*gqlmodel.ProductMediaUpdate, error) {
	panic(fmt.Errorf("not implemented"))
}
