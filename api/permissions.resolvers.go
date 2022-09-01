package api

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"

	"github.com/sitename/sitename/api/gqlmodel"
)

func (r *Resolver) PermissionGroupCreate(ctx context.Context, args struct {
	input gqlmodel.PermissionGroupCreateInput
}) (*gqlmodel.PermissionGroupCreate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) PermissionGroupUpdate(ctx context.Context, args struct {
	id    string
	input gqlmodel.PermissionGroupUpdateInput
}) (*gqlmodel.PermissionGroupUpdate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) PermissionGroupDelete(ctx context.Context, args struct{ id string }) (*gqlmodel.PermissionGroupDelete, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) PermissionGroups(ctx context.Context, args struct {
	filter *gqlmodel.PermissionGroupFilterInput
	sortBy *gqlmodel.PermissionGroupSortingInput
	before *string
	after  *string
	first  *int
	last   *int
}) (*gqlmodel.GroupCountableConnection, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) PermissionGroup(ctx context.Context, args struct{ id string }) (*gqlmodel.Group, error) {
	panic(fmt.Errorf("not implemented"))
}
