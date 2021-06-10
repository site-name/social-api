package graphql

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"

	"github.com/sitename/sitename/web/graphql/gqlmodel"
)

func (r *mutationResolver) ProductMediaCreate(ctx context.Context, input gqlmodel.ProductMediaCreateInput) (*gqlmodel.ProductMediaCreate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) ProductMediaDelete(ctx context.Context, id string) (*gqlmodel.ProductMediaDelete, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) ProductMediaBulkDelete(ctx context.Context, ids []*string) (*gqlmodel.ProductMediaBulkDelete, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) ProductMediaReorder(ctx context.Context, mediaIds []*string, productID string) (*gqlmodel.ProductMediaReorder, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) ProductMediaUpdate(ctx context.Context, id string, input gqlmodel.ProductMediaUpdateInput) (*gqlmodel.ProductMediaUpdate, error) {
	panic(fmt.Errorf("not implemented"))
}
