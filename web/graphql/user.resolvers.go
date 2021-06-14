package graphql

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"

	"github.com/99designs/gqlgen/graphql"
	"github.com/sitename/sitename/web/graphql/gqlmodel"
)

func (r *customerEventResolver) User(ctx context.Context, obj *gqlmodel.CustomerEvent) (*gqlmodel.User, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *customerEventResolver) Order(ctx context.Context, obj *gqlmodel.CustomerEvent) (*gqlmodel.Order, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *customerEventResolver) OrderLine(ctx context.Context, obj *gqlmodel.CustomerEvent) (*gqlmodel.OrderLine, error) {
	panic(fmt.Errorf("not implemented"))
}

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

func (r *userResolver) DefaultShippingAddress(ctx context.Context, obj *gqlmodel.User) (*gqlmodel.Address, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *userResolver) DefaultBillingAddress(ctx context.Context, obj *gqlmodel.User) (*gqlmodel.Address, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *userResolver) Addresses(ctx context.Context, obj *gqlmodel.User) ([]*gqlmodel.Address, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *userResolver) Events(ctx context.Context, obj *gqlmodel.User) ([]*gqlmodel.CustomerEvent, error) {
	panic(fmt.Errorf("not implemented"))
}

// CustomerEvent returns CustomerEventResolver implementation.
func (r *Resolver) CustomerEvent() CustomerEventResolver { return &customerEventResolver{r} }

// User returns UserResolver implementation.
func (r *Resolver) User() UserResolver { return &userResolver{r} }

type customerEventResolver struct{ *Resolver }
type userResolver struct{ *Resolver }
