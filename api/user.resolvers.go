package api

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"

	"github.com/99designs/gqlgen/graphql"
	"github.com/sitename/sitename/api/gqlmodel"
)

func (r *Resolver) Login(ctx context.Context, input gqlmodel.LoginInput) (*gqlmodel.LoginResponse, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) UserAvatarUpdate(ctx context.Context, image graphql.Upload) (*gqlmodel.UserAvatarUpdate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) UserAvatarDelete(ctx context.Context) (*gqlmodel.UserAvatarDelete, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) UserBulkSetActive(ctx context.Context, ids []*string, isActive bool) (*gqlmodel.UserBulkSetActive, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) Me(ctx context.Context) (*gqlmodel.User, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) User(ctx context.Context, id *string, email *string) (*gqlmodel.User, error) {
	panic(fmt.Errorf("not implemented"))
}
