package graphql

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"

	"github.com/sitename/sitename/web/graphql/gqlmodel"
)

func (r *mutationResolver) AccountAddressCreate(ctx context.Context, input gqlmodel.AddressInput, typeArg *gqlmodel.AddressTypeEnum) (*gqlmodel.AccountAddressCreate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) AccountAddressUpdate(ctx context.Context, id string, input gqlmodel.AddressInput) (*gqlmodel.AccountAddressUpdate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) AccountAddressDelete(ctx context.Context, id string) (*gqlmodel.AccountAddressDelete, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) AccountSetDefaultAddress(ctx context.Context, id string, typeArg gqlmodel.AddressTypeEnum) (*gqlmodel.AccountSetDefaultAddress, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) AccountRegister(ctx context.Context, input gqlmodel.AccountRegisterInput) (*gqlmodel.AccountRegister, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) AccountUpdate(ctx context.Context, input gqlmodel.AccountInput) (*gqlmodel.AccountUpdate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) AccountRequestDeletion(ctx context.Context, channel *string, redirectURL string) (*gqlmodel.AccountRequestDeletion, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) AccountDelete(ctx context.Context, token string) (*gqlmodel.AccountDelete, error) {
	panic(fmt.Errorf("not implemented"))
}
