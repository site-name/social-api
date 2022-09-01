package api

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"

	"github.com/sitename/sitename/api/gqlmodel"
)

func (r *Resolver) PermissionGroupCreate(ctx context.Context, args struct {
	Input gqlmodel.PermissionGroupCreateInput
}) (*gqlmodel.PermissionGroupCreate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) PermissionGroupUpdate(ctx context.Context, args struct {
	Id    string
	Input gqlmodel.PermissionGroupUpdateInput
}) (*gqlmodel.PermissionGroupUpdate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) PermissionGroupDelete(ctx context.Context, args struct{ Id string }) (*gqlmodel.PermissionGroupDelete, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) PermissionGroups(ctx context.Context, args struct {
	Filter *gqlmodel.PermissionGroupFilterInput
	SortBy *gqlmodel.PermissionGroupSortingInput
	Before *string
	After  *string
	First  *int
	Last   *int
}) (*gqlmodel.GroupCountableConnection, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) PermissionGroup(ctx context.Context, args struct{ Id string }) (*gqlmodel.Group, error) {
	panic(fmt.Errorf("not implemented"))
}
