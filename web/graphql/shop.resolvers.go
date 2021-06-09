package graphql

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"
)

func (r *mutationResolver) ShopDomainUpdate(ctx context.Context, input *SiteDomainInput) (*ShopDomainUpdate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) ShopSettingsUpdate(ctx context.Context, input ShopSettingsInput) (*ShopSettingsUpdate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) ShopFetchTaxRates(ctx context.Context) (*ShopFetchTaxRates, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) ShopSettingsTranslate(ctx context.Context, input ShopSettingsTranslationInput, languageCode LanguageCodeEnum) (*ShopSettingsTranslate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) ShopAddressUpdate(ctx context.Context, input *AddressInput) (*ShopAddressUpdate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) Shop(ctx context.Context) (*Shop, error) {
	panic(fmt.Errorf("not implemented"))
}
