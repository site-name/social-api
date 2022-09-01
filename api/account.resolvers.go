package api

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"

	"github.com/sitename/sitename/api/gqlmodel"
)

func (r *Resolver) AccountAddressCreate(ctx context.Context, args struct {
	input   gqlmodel.AddressInput
	typeArg *gqlmodel.AddressTypeEnum
}) (*gqlmodel.AccountAddressCreate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) AccountAddressUpdate(ctx context.Context, args struct {
	id    string
	input gqlmodel.AddressInput
}) (*gqlmodel.AccountAddressUpdate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) AccountAddressDelete(ctx context.Context, args struct{ id string }) (*gqlmodel.AccountAddressDelete, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) AccountSetDefaultAddress(ctx context.Context, args struct {
	id      string
	typeArg gqlmodel.AddressTypeEnum
}) (*gqlmodel.AccountSetDefaultAddress, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) AccountRegister(ctx context.Context, args struct{ input gqlmodel.AccountRegisterInput }) (*gqlmodel.AccountRegister, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) AccountUpdate(ctx context.Context, args struct{ input gqlmodel.AccountInput }) (*gqlmodel.AccountUpdate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) AccountRequestDeletion(ctx context.Context, args struct {
	channel     *string
	redirectURL string
}) (*gqlmodel.AccountRequestDeletion, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) AccountDelete(ctx context.Context, args struct{ token string }) (*gqlmodel.AccountDelete, error) {
	panic(fmt.Errorf("not implemented"))
}
