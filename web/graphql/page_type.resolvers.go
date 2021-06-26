package graphql

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"

	"github.com/sitename/sitename/web/graphql/gqlmodel"
)

func (r *mutationResolver) PageTypeCreate(ctx context.Context, input gqlmodel.PageTypeCreateInput) (*gqlmodel.PageTypeCreate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) PageTypeUpdate(ctx context.Context, id *string, input gqlmodel.PageTypeUpdateInput) (*gqlmodel.PageTypeUpdate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) PageTypeDelete(ctx context.Context, id string) (*gqlmodel.PageTypeDelete, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) PageTypeBulkDelete(ctx context.Context, ids []string) (*gqlmodel.PageTypeBulkDelete, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) PageTypeReorderAttributes(ctx context.Context, moves []gqlmodel.ReorderInput, pageTypeID string) (*gqlmodel.PageTypeReorderAttributes, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) PageType(ctx context.Context, id string) (*gqlmodel.PageType, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) PageTypes(ctx context.Context, sortBy *gqlmodel.PageTypeSortingInput, filter *gqlmodel.PageTypeFilterInput, before *string, after *string, first *int, last *int) (*gqlmodel.PageTypeCountableConnection, error) {
	panic(fmt.Errorf("not implemented"))
}
