package api

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"
)

func (r *Resolver) PageTypeCreate(ctx context.Context, args struct{ Input PageTypeCreateInput }) (*PageTypeCreate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) PageTypeUpdate(ctx context.Context, args struct {
	Id    *string
	Input PageTypeUpdateInput
}) (*PageTypeUpdate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) PageTypeDelete(ctx context.Context, args struct{ Id string }) (*PageTypeDelete, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) PageTypeBulkDelete(ctx context.Context, args struct{ Ids []string }) (*PageTypeBulkDelete, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) PageTypeReorderAttributes(ctx context.Context, args struct {
	Moves      []ReorderInput
	PageTypeID string
}) (*PageTypeReorderAttributes, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) PageType(ctx context.Context, args struct{ Id string }) (*PageType, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) PageTypes(ctx context.Context, args struct {
	SortBy *PageTypeSortingInput
	Filter *PageTypeFilterInput
	Before *string
	After  *string
	First  *int
	Last   *int
}) (*PageTypeCountableConnection, error) {
	panic(fmt.Errorf("not implemented"))
}
