package api

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"
)

func (r *Resolver) PermissionGroupCreate(ctx context.Context, args struct {
	Input PermissionGroupCreateInput
}) (*PermissionGroupCreate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) PermissionGroupUpdate(ctx context.Context, args struct {
	Id    string
	Input PermissionGroupUpdateInput
}) (*PermissionGroupUpdate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) PermissionGroupDelete(ctx context.Context, args struct{ Id string }) (*PermissionGroupDelete, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) PermissionGroups(ctx context.Context, args struct {
	Filter *PermissionGroupFilterInput
	SortBy *PermissionGroupSortingInput
	Before *string
	After  *string
	First  *int
	Last   *int
}) (*GroupCountableConnection, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) PermissionGroup(ctx context.Context, args struct{ Id string }) (*Group, error) {
	panic(fmt.Errorf("not implemented"))
}
