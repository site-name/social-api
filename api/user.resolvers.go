package api

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"

	"github.com/99designs/gqlgen/graphql"
)

func (r *Resolver) Login(ctx context.Context, args struct{ Input LoginInput }) (*LoginResponse, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) UserAvatarUpdate(ctx context.Context, args struct{ Image graphql.Upload }) (*UserAvatarUpdate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) UserAvatarDelete(ctx context.Context) (*UserAvatarDelete, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) UserBulkSetActive(ctx context.Context, args struct {
	Ids      []*string
	IsActive bool
}) (*UserBulkSetActive, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) Me(ctx context.Context) (*User, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) User(ctx context.Context, args struct {
	Id    *string
	Email *string
}) (*User, error) {
	panic(fmt.Errorf("not implemented"))
}
