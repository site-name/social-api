package api

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"

	"github.com/sitename/sitename/api/gqlmodel"
)

func (r *Resolver) AccountAddressCreate(ctx context.Context, args struct {
	Input   gqlmodel.AddressInput
	TypeArg *gqlmodel.AddressTypeEnum
}) (*gqlmodel.AccountAddressCreate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) AccountAddressUpdate(ctx context.Context, args struct {
	Id    string
	Input gqlmodel.AddressInput
}) (*gqlmodel.AccountAddressUpdate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) AccountAddressDelete(ctx context.Context, args struct{ Id string }) (*gqlmodel.AccountAddressDelete, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) AccountSetDefaultAddress(ctx context.Context, args struct {
	Id      string
	TypeArg gqlmodel.AddressTypeEnum
}) (*gqlmodel.AccountSetDefaultAddress, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) AccountRegister(ctx context.Context, args struct{ Input gqlmodel.AccountRegisterInput }) (*gqlmodel.AccountRegister, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) AccountUpdate(ctx context.Context, args struct{ Input gqlmodel.AccountInput }) (*gqlmodel.AccountUpdate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) AccountRequestDeletion(ctx context.Context, args struct {
	Channel     *string
	RedirectURL string
}) (*gqlmodel.AccountRequestDeletion, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) AccountDelete(ctx context.Context, args struct{ Token string }) (*gqlmodel.AccountDelete, error) {
	panic(fmt.Errorf("not implemented"))
}
