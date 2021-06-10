package graphql

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"

	"github.com/sitename/sitename/web/graphql/gqlmodel"
)

func (r *mutationResolver) DraftOrderComplete(ctx context.Context, id string) (*gqlmodel.DraftOrderComplete, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) DraftOrderCreate(ctx context.Context, input gqlmodel.DraftOrderCreateInput) (*gqlmodel.DraftOrderCreate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) DraftOrderDelete(ctx context.Context, id string) (*gqlmodel.DraftOrderDelete, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) DraftOrderBulkDelete(ctx context.Context, ids []*string) (*gqlmodel.DraftOrderBulkDelete, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) DraftOrderLinesBulkDelete(ctx context.Context, ids []*string) (*gqlmodel.DraftOrderLinesBulkDelete, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) DraftOrderUpdate(ctx context.Context, id string, input gqlmodel.DraftOrderInput) (*gqlmodel.DraftOrderUpdate, error) {
	panic(fmt.Errorf("not implemented"))
}
