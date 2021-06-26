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

func (r *shopResolver) AvailablePaymentGateways(ctx context.Context, obj *gqlmodel.Shop, currency *string, channel *string) ([]gqlmodel.PaymentGateway, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *shopResolver) AvailableExternalAuthentications(ctx context.Context, obj *gqlmodel.Shop) ([]gqlmodel.ExternalAuthentication, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *shopResolver) AvailableShippingMethods(ctx context.Context, obj *gqlmodel.Shop, channel string, address *gqlmodel.AddressInput) ([]*gqlmodel.ShippingMethod, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *shopResolver) Countries(ctx context.Context, obj *gqlmodel.Shop, languageCode *gqlmodel.LanguageCodeEnum) ([]gqlmodel.CountryDisplay, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *shopResolver) Translation(ctx context.Context, obj *gqlmodel.Shop, languageCode gqlmodel.LanguageCodeEnum) (*gqlmodel.ShopTranslation, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *shopResolver) CompanyAddress(ctx context.Context, obj *gqlmodel.Shop) (*gqlmodel.Address, error) {
	panic(fmt.Errorf("not implemented"))
}

// Shop returns ShopResolver implementation.
func (r *Resolver) Shop() ShopResolver { return &shopResolver{r} }

type shopResolver struct{ *Resolver }
