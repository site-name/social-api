package api

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"

	"github.com/sitename/sitename/api/gqlmodel"
)

func (r *Resolver) ShopDomainUpdate(ctx context.Context, args struct{ Input *gqlmodel.SiteDomainInput }) (*gqlmodel.ShopDomainUpdate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) ShopSettingsUpdate(ctx context.Context, args struct{ Input gqlmodel.ShopSettingsInput }) (*gqlmodel.ShopSettingsUpdate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) ShopFetchTaxRates(ctx context.Context) (*gqlmodel.ShopFetchTaxRates, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) ShopSettingsTranslate(ctx context.Context, args struct {
	Input        gqlmodel.ShopSettingsTranslationInput
	LanguageCode gqlmodel.LanguageCodeEnum
}) (*gqlmodel.ShopSettingsTranslate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) ShopAddressUpdate(ctx context.Context, args struct{ Input *gqlmodel.AddressInput }) (*gqlmodel.ShopAddressUpdate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) Shop(ctx context.Context) (*gqlmodel.Shop, error) {
	panic(fmt.Errorf("not implemented"))
}
