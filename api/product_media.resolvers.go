package api

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"
)

func (r *Resolver) ProductMediaCreate(ctx context.Context, args struct {
	Input ProductMediaCreateInput
}) (*ProductMediaCreate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) ProductMediaDelete(ctx context.Context, args struct{ Id string }) (*ProductMediaDelete, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) ProductMediaBulkDelete(ctx context.Context, args struct{ Ids []*string }) (*ProductMediaBulkDelete, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) ProductMediaReorder(ctx context.Context, args struct {
	MediaIds  []*string
	ProductID string
}) (*ProductMediaReorder, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) ProductMediaUpdate(ctx context.Context, args struct {
	Id    string
	Input ProductMediaUpdateInput
}) (*ProductMediaUpdate, error) {
	panic(fmt.Errorf("not implemented"))
}
