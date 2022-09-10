package api

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"
)

func (r *Resolver) AddressCreate(ctx context.Context, args struct {
	Input  AddressInput
	UserID string
}) (*AddressCreate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) AddressUpdate(ctx context.Context, args struct {
	Id    string
	Input AddressInput
}) (*AddressUpdate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) AddressDelete(ctx context.Context, args struct{ Id string }) (*AddressDelete, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) AddressSetDefault(ctx context.Context, args struct {
	AddressID string
	TypeArg   AddressTypeEnum
	UserID    string
}) (*AddressSetDefault, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) AddressValidationRules(ctx context.Context, args struct {
	CountryCode CountryCode
	CountryArea *string
	City        *string
	CityArea    *string
}) (*AddressValidationData, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) Address(ctx context.Context, args struct{ Id string }) (*Address, error) {
	panic(fmt.Errorf("not implemented"))
}
