package graphql

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"

	graphql1 "github.com/sitename/sitename/graphql/generated"
	"github.com/sitename/sitename/graphql/gqlmodel"
)

func (r *groupResolver) Users(ctx context.Context, obj *gqlmodel.Group) ([]*gqlmodel.User, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) PermissionGroupCreate(ctx context.Context, input gqlmodel.PermissionGroupCreateInput) (*gqlmodel.PermissionGroupCreate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) PermissionGroupUpdate(ctx context.Context, id string, input gqlmodel.PermissionGroupUpdateInput) (*gqlmodel.PermissionGroupUpdate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) PermissionGroupDelete(ctx context.Context, id string) (*gqlmodel.PermissionGroupDelete, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) PermissionGroups(ctx context.Context, filter *gqlmodel.PermissionGroupFilterInput, sortBy *gqlmodel.PermissionGroupSortingInput, before *string, after *string, first *int, last *int) (*gqlmodel.GroupCountableConnection, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) PermissionGroup(ctx context.Context, id string) (*gqlmodel.Group, error) {
	panic(fmt.Errorf("not implemented"))
}

// Group returns graphql1.GroupResolver implementation.
func (r *Resolver) Group() graphql1.GroupResolver { return &groupResolver{r} }

type groupResolver struct{ *Resolver }
