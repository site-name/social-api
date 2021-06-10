package graphql

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"

	"github.com/99designs/gqlgen/graphql"
	"github.com/sitename/sitename/web/graphql/gqlmodel"
)

func (r *mutationResolver) Login(ctx context.Context, input gqlmodel.LoginInput) (*gqlmodel.LoginResponse, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) UserAvatarUpdate(ctx context.Context, image graphql.Upload) (*gqlmodel.UserAvatarUpdate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) UserAvatarDelete(ctx context.Context) (*gqlmodel.UserAvatarDelete, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) UserBulkSetActive(ctx context.Context, ids []*string, isActive bool) (*gqlmodel.UserBulkSetActive, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) Me(ctx context.Context) (*gqlmodel.User, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) User(ctx context.Context, id *string, email *string) (*gqlmodel.User, error) {
	panic(fmt.Errorf("not implemented"))
}
