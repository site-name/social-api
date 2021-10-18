package graphql

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"

	"github.com/sitename/sitename/graphql/gqlmodel"
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

func (r *mutationResolver) DraftOrderUpdate(ctx context.Context, id string, input gqlmodel.DraftOrderInput) (*gqlmodel.DraftOrderUpdate, error) {
	panic(fmt.Errorf("not implemented"))
}

// !!! WARNING !!!
// The code below was going to be deleted when updating resolvers. It has been copied here so you have
// one last chance to move it out of harms way if you want. There are two reasons this happens:
//  - When renaming or deleting a resolver the old code will be put in here. You can safely delete
//    it when you're done.
//  - You have helper methods in this file. Move them out to keep these resolver files clean.
func (r *mutationResolver) DraftOrderLinesBulkDelete(ctx context.Context, ids []*string) (*gqlmodel.DraftOrderLinesBulkDelete, error) {
	panic(fmt.Errorf("not implemented"))
}
