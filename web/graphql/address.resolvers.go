package graphql

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"
)

func (r *mutationResolver) AddressCreate(ctx context.Context, input AddressInput, userID string) (*AddressCreate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) AddressUpdate(ctx context.Context, id string, input AddressInput) (*AddressUpdate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) AddressDelete(ctx context.Context, id string) (*AddressDelete, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) AddressSetDefault(ctx context.Context, addressID string, typeArg AddressTypeEnum, userID string) (*AddressSetDefault, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) AddressValidationRules(ctx context.Context, countryCode CountryCode, countryArea *string, city *string, cityArea *string) (*AddressValidationData, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) Address(ctx context.Context, id string) (*Address, error) {
	panic(fmt.Errorf("not implemented"))
}
