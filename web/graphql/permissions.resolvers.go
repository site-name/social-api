package graphql

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"
)

func (r *mutationResolver) PermissionGroupCreate(ctx context.Context, input PermissionGroupCreateInput) (*PermissionGroupCreate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) PermissionGroupUpdate(ctx context.Context, id string, input PermissionGroupUpdateInput) (*PermissionGroupUpdate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) PermissionGroupDelete(ctx context.Context, id string) (*PermissionGroupDelete, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) PermissionGroups(ctx context.Context, filter *PermissionGroupFilterInput, sortBy *PermissionGroupSortingInput, before *string, after *string, first *int, last *int) (*GroupCountableConnection, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) PermissionGroup(ctx context.Context, id string) (*Group, error) {
	panic(fmt.Errorf("not implemented"))
}
