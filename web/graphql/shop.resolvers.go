package graphql

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"

	"github.com/sitename/sitename/web/graphql/gqlmodel"
)

func (r *mutationResolver) ShopDomainUpdate(ctx context.Context, input *gqlmodel.SiteDomainInput) (*gqlmodel.ShopDomainUpdate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) ShopSettingsUpdate(ctx context.Context, input gqlmodel.ShopSettingsInput) (*gqlmodel.ShopSettingsUpdate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) ShopFetchTaxRates(ctx context.Context) (*gqlmodel.ShopFetchTaxRates, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) ShopSettingsTranslate(ctx context.Context, input gqlmodel.ShopSettingsTranslationInput, languageCode gqlmodel.LanguageCodeEnum) (*gqlmodel.ShopSettingsTranslate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) ShopAddressUpdate(ctx context.Context, input *gqlmodel.AddressInput) (*gqlmodel.ShopAddressUpdate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) Shop(ctx context.Context) (*gqlmodel.Shop, error) {
	panic(fmt.Errorf("not implemented"))
}
