package graphql

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"
)

func (r *mutationResolver) AccountAddressCreate(ctx context.Context, input AddressInput, typeArg *AddressTypeEnum) (*AccountAddressCreate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) AccountAddressUpdate(ctx context.Context, id string, input AddressInput) (*AccountAddressUpdate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) AccountAddressDelete(ctx context.Context, id string) (*AccountAddressDelete, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) AccountSetDefaultAddress(ctx context.Context, id string, typeArg AddressTypeEnum) (*AccountSetDefaultAddress, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) AccountRegister(ctx context.Context, input AccountRegisterInput) (*AccountRegister, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) AccountUpdate(ctx context.Context, input AccountInput) (*AccountUpdate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) AccountRequestDeletion(ctx context.Context, channel *string, redirectURL string) (*AccountRequestDeletion, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) AccountDelete(ctx context.Context, token string) (*AccountDelete, error) {
	panic(fmt.Errorf("not implemented"))
}
