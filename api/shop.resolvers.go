package api

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"
)

func (r *Resolver) ShopDomainUpdate(ctx context.Context, args struct{ Input *SiteDomainInput }) (*ShopDomainUpdate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) ShopSettingsUpdate(ctx context.Context, args struct{ Input ShopSettingsInput }) (*ShopSettingsUpdate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) ShopFetchTaxRates(ctx context.Context) (*ShopFetchTaxRates, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) ShopSettingsTranslate(ctx context.Context, args struct {
	Input        ShopSettingsTranslationInput
	LanguageCode LanguageCodeEnum
}) (*ShopSettingsTranslate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) ShopAddressUpdate(ctx context.Context, args struct{ Input *AddressInput }) (*ShopAddressUpdate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) Shop(ctx context.Context) (*Shop, error) {
	panic(fmt.Errorf("not implemented"))
}

// NOTE: Refer to ./schemas/shop.graphqls for details on directive used.
func (r *Resolver) GiftCardSettings(ctx context.Context) (*GiftCardSettings, error) {
	panic(fmt.Errorf("not implemented"))
}
