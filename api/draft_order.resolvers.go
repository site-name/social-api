package api

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"
)

func (r *Resolver) DraftOrderComplete(ctx context.Context, args struct{ Id string }) (*DraftOrderComplete, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) DraftOrderCreate(ctx context.Context, args struct {
	Input DraftOrderCreateInput
}) (*DraftOrderCreate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) DraftOrderDelete(ctx context.Context, args struct{ Id string }) (*DraftOrderDelete, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) DraftOrderBulkDelete(ctx context.Context, args struct{ Ids []*string }) (*DraftOrderBulkDelete, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) DraftOrderUpdate(ctx context.Context, args struct {
	Id    string
	Input DraftOrderInput
}) (*DraftOrderUpdate, error) {
	panic(fmt.Errorf("not implemented"))
}
