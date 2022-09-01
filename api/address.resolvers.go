package api

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"

	"github.com/sitename/sitename/api/gqlmodel"
)

func (r *Resolver) AddressCreate(ctx context.Context, args struct {
	input  gqlmodel.AddressInput
	userID string
}) (*gqlmodel.AddressCreate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) AddressUpdate(ctx context.Context, args struct {
	id    string
	input gqlmodel.AddressInput
}) (*gqlmodel.AddressUpdate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) AddressDelete(ctx context.Context, args struct{ id string }) (*gqlmodel.AddressDelete, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) AddressSetDefault(ctx context.Context, args struct {
	addressID string
	typeArg   gqlmodel.AddressTypeEnum
	userID    string
}) (*gqlmodel.AddressSetDefault, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) AddressValidationRules(ctx context.Context, args struct {
	countryCode gqlmodel.CountryCode
	countryArea *string
	city        *string
	cityArea    *string
}) (*gqlmodel.AddressValidationData, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) Address(ctx context.Context, args struct{ id string }) (*gqlmodel.Address, error) {
	panic(fmt.Errorf("not implemented"))
}
