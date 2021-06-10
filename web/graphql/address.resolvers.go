package graphql

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"

	"github.com/sitename/sitename/web/graphql/gqlmodel"
)

func (r *mutationResolver) AddressCreate(ctx context.Context, input gqlmodel.AddressInput, userID string) (*gqlmodel.AddressCreate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) AddressUpdate(ctx context.Context, id string, input gqlmodel.AddressInput) (*gqlmodel.AddressUpdate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) AddressDelete(ctx context.Context, id string) (*gqlmodel.AddressDelete, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) AddressSetDefault(ctx context.Context, addressID string, typeArg gqlmodel.AddressTypeEnum, userID string) (*gqlmodel.AddressSetDefault, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) AddressValidationRules(ctx context.Context, countryCode gqlmodel.CountryCode, countryArea *string, city *string, cityArea *string) (*gqlmodel.AddressValidationData, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) Address(ctx context.Context, id string) (*gqlmodel.Address, error) {
	panic(fmt.Errorf("not implemented"))
}
