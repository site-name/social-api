package api

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"
	"unsafe"

	"github.com/samber/lo"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/web"
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

var (
	shopLanguages = lo.MapToSlice(model.Languages, func(code model.LanguageCodeEnum, name string) *LanguageDisplay {
		return &LanguageDisplay{
			Code:     code,
			Language: name,
		}
	})
)

func (r *Resolver) Shop(ctx context.Context) (*Shop, error) {
	embedCtx := GetContextValue[*web.Context](ctx, WebCtx)
	shopSettings := embedCtx.App.Config().ShopSettings
	serviceSettings := embedCtx.App.Config().ServiceSettings

	vats, appErr := embedCtx.App.Srv().DiscountService().FilterVats(&model.VatFilterOptions{})
	if appErr != nil {
		return nil, appErr
	}

	taxes := lo.SliceToMap(vats, func(item *model.Vat) (model.CountryCode, *model.Vat) {
		return item.CountryCode, item
	})

	shopCountries := []*CountryDisplay{}
	for countryCode, name := range model.Countries {
		shopCountries = append(shopCountries, &CountryDisplay{
			Code:    countryCode.String(),
			Country: name,
			Vat:     systemVatToGraphqlVat(taxes[countryCode]),
		})
	}

	res := Shop{
		Countries:                           shopCountries,
		DefaultMailSenderName:               shopSettings.DefaultMailSenderName,
		DefaultMailSenderAddress:            shopSettings.DefaultMailSenderAddress,
		AutomaticFulfillmentDigitalProducts: shopSettings.AutomaticFulfillmentDigitalProducts,
		DefaultDigitalMaxDownloads:          (*int32)(unsafe.Pointer(shopSettings.DefaultDigitalMaxDownloads)),
		DefaultDigitalURLValidDays:          (*int32)(unsafe.Pointer(shopSettings.DefaultDigitalUrlValidDays)),
		IncludeTaxesInPrices:                *shopSettings.IncludeTaxesInPrice,
		FulfillmentAutoApprove:              *shopSettings.FulfillmentAutoApprove,
		DisplayGrossPrices:                  *shopSettings.DisplayGrossPrices,
		ChargeTaxesOnShipping:               *shopSettings.ChargeTaxesOnShipping,
		TrackInventoryByDefault:             shopSettings.TrackInventoryByDefault,
		DefaultWeightUnit:                   shopSettings.DefaultWeightUnit,
		CompanyAddress:                      SystemAddressToGraphqlAddress(shopSettings.Address),
		CustomerSetPasswordURL:              shopSettings.CustomerSetPasswordUrl,
		Version:                             model.CurrentVersion,
		Name:                                *serviceSettings.SiteName,
		Languages:                           shopLanguages,

		Description:   nil, // TODO: fix me
		Domain:        nil, // TODO: fix me
		Permissions:   nil, // TODO: fix me
		PhonePrefixes: nil, // TODO: fix me
		HeaderText:    nil, // TODO: fix me
		Limits:        nil, // TODO: fix me
		Translation:   nil, // TODO: fix me
	}

	return &res, nil
}

// NOTE: Refer to ./schemas/shop.graphqls for details on directive used.
func (r *Resolver) GiftCardSettings(ctx context.Context) (*GiftCardSettings, error) {
	shopSettings := GetContextValue[*web.Context](ctx, WebCtx).App.Config().ShopSettings

	return &GiftCardSettings{
		ExpiryType: *shopSettings.GiftcardExpiryType,
		ExpiryPeriod: &TimePeriod{
			Amount: *(*int32)(unsafe.Pointer(shopSettings.GiftcardExpiryPeriod)),
			Type:   *shopSettings.GiftcardExpiryPeriodType,
		},
	}, nil
}
