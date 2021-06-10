package graphql

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"

	"github.com/99designs/gqlgen/graphql"
)

func (r *mutationResolver) Login(ctx context.Context, input LoginInput) (*LoginResponse, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) UserAvatarUpdate(ctx context.Context, image graphql.Upload) (*UserAvatarUpdate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) UserAvatarDelete(ctx context.Context) (*UserAvatarDelete, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) UserBulkSetActive(ctx context.Context, ids []*string, isActive bool) (*UserBulkSetActive, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) Me(ctx context.Context) (*User, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) User(ctx context.Context, id *string, email *string) (*User, error) {
	panic(fmt.Errorf("not implemented"))
}
