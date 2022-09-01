package api

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"

	"github.com/sitename/sitename/api/gqlmodel"
)

func (r *Resolver) AddressCreate(ctx context.Context, args struct {
	Input  gqlmodel.AddressInput
	UserID string
}) (*gqlmodel.AddressCreate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) AddressUpdate(ctx context.Context, args struct {
	Id    string
	Input gqlmodel.AddressInput
}) (*gqlmodel.AddressUpdate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) AddressDelete(ctx context.Context, args struct{ Id string }) (*gqlmodel.AddressDelete, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) AddressSetDefault(ctx context.Context, args struct {
	AddressID string
	TypeArg   gqlmodel.AddressTypeEnum
	UserID    string
}) (*gqlmodel.AddressSetDefault, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) AddressValidationRules(ctx context.Context, args struct {
	CountryCode gqlmodel.CountryCode
	CountryArea *string
	City        *string
	CityArea    *string
}) (*gqlmodel.AddressValidationData, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) Address(ctx context.Context, args struct{ Id string }) (*gqlmodel.Address, error) {
	panic(fmt.Errorf("not implemented"))
}
