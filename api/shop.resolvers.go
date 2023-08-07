package api

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"unsafe"

	"github.com/samber/lo"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/web"
)

func (r *Resolver) ShopDomainUpdate(ctx context.Context, args struct{ Input *SiteDomainInput }) (*ShopDomainUpdate, error) {
	panic(fmt.Errorf("not implemented"))
}

// NOTE: Refer to ./schemas/shop.graphqls for details on directives used.
func (r *Resolver) ShopSettingsUpdate(ctx context.Context, args struct{ Input ShopSettingsInput }) (*ShopSettingsUpdate, error) {
	embedCtx := GetContextValue[*web.Context](ctx, WebCtx)

	// validate argument
	if inputUrl := args.Input.CustomerSetPasswordURL; inputUrl != nil {
		allowedClientHosts := strings.Fields(*embedCtx.App.Config().ServiceSettings.AllowCorsFrom)

		urlParse, err := url.Parse(*inputUrl)
		if err != nil || !lo.Contains(allowedClientHosts, urlParse.Host) {
			return nil, model.NewAppError("ShopSettingsUpdate", model.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "customer_set_password_url"}, err.Error(), http.StatusBadRequest)
		}
	}

	embedCtx.App.UpdateConfig(func(c *model.Config) {
		if vl := args.Input.HeaderText; vl != nil {
			// TODO: implement me later
		}
		if vl := args.Input.Description; vl != nil && *vl != *c.ShopSettings.Description {
			c.ShopSettings.Description = vl
		}
		if vl := args.Input.IncludeTaxesInPrices; vl != nil {
			c.ShopSettings.IncludeTaxesInPrice = vl
		}
		if vl := args.Input.DisplayGrossPrices; vl != nil {
			c.ShopSettings.DisplayGrossPrices = vl
		}
		if vl := args.Input.ChargeTaxesOnShipping; vl != nil {
			c.ShopSettings.ChargeTaxesOnShipping = vl
		}
		if vl := args.Input.TrackInventoryByDefault; vl != nil {
			c.ShopSettings.TrackInventoryByDefault = vl
		}
		if vl := args.Input.DefaultWeightUnit; vl != nil {
			c.ShopSettings.DefaultWeightUnit = vl
		}
		if vl := args.Input.AutomaticFulfillmentDigitalProducts; vl != nil {
			c.ShopSettings.AutomaticFulfillmentDigitalProducts = vl
		}
		if vl := args.Input.FulfillmentAutoApprove; vl != nil {
			c.ShopSettings.FulfillmentAutoApprove = vl
		}
		if vl := args.Input.FulfillmentAllowUnpaid; vl != nil {
			c.ShopSettings.FulfillmentAllowUnPaid = vl
		}
		if vl := args.Input.DefaultDigitalMaxDownloads; vl != nil {
			c.ShopSettings.DefaultDigitalMaxDownloads = (*int)(unsafe.Pointer(vl))
		}
		if vl := args.Input.DefaultDigitalURLValidDays; vl != nil {
			c.ShopSettings.DefaultDigitalUrlValidDays = (*int)(unsafe.Pointer(vl))
		}
		if vl := args.Input.DefaultMailSenderName; vl != nil {
			c.ShopSettings.DefaultMailSenderName = vl
		}
		if vl := args.Input.DefaultMailSenderAddress; vl != nil {
			c.ShopSettings.DefaultMailSenderAddress = vl
		}
		if vl := args.Input.CustomerSetPasswordURL; vl != nil {
			c.ShopSettings.CustomerSetPasswordUrl = vl
		}
	})

	shop, err := r.Shop(ctx)
	if err != nil {
		return nil, err
	}

	return &ShopSettingsUpdate{
		Shop: shop,
	}, nil
}

// NOTE: Refer to ./schemas/shop.graphqls for details on directives used.
func (r *Resolver) ShopFetchTaxRates(ctx context.Context) (*ShopFetchTaxRates, error) {
	embedCtx := GetContextValue[*web.Context](ctx, WebCtx)
	pluginMng := embedCtx.App.Srv().PluginService().GetPluginManager()

	boolValue, appErr := pluginMng.FetchTaxesData()
	if appErr != nil {
		return nil, appErr
	}

	if !boolValue {
		return nil, model.NewAppError("ShopFetchTaxRates", "api.shop.no_credential_for_tax_plugin.app_error", nil, "Please provile a valid credential for your tax plugin", http.StatusNotAcceptable)
	}

	shop, err := r.Shop(ctx)
	if err != nil {
		return nil, err
	}

	return &ShopFetchTaxRates{
		Shop: shop,
	}, nil
}

func (r *Resolver) ShopSettingsTranslate(ctx context.Context, args struct {
	Input        ShopSettingsTranslationInput
	LanguageCode LanguageCodeEnum
}) (*ShopSettingsTranslate, error) {
	panic(fmt.Errorf("not implemented"))
}

// NOTE: Refer to ./schemas/shop.graphqls for details on directives used.
func (r *Resolver) ShopAddressUpdate(ctx context.Context, args struct{ Input AddressInput }) (*ShopAddressUpdate, error) {
	// validate argument
	if err := args.Input.Validate("ShopAddressUpdate"); err != nil {
		return nil, err
	}

	embedCtx := GetContextValue[*web.Context](ctx, WebCtx)
	embedCtx.App.UpdateConfig(func(c *model.Config) {
		args.Input.PatchAddress(c.ShopSettings.Address)
	})

	shop, err := r.Shop(ctx)
	if err != nil {
		return nil, err
	}

	return &ShopAddressUpdate{
		Shop: shop,
	}, nil
}

func (r *Resolver) Shop(ctx context.Context) (*Shop, error) {
	embedCtx := GetContextValue[*web.Context](ctx, WebCtx)

	vats, appErr := embedCtx.App.Srv().DiscountService().FilterVats(&model.VatFilterOptions{})
	if appErr != nil {
		return nil, appErr
	}
	return systemConfigToGraphqlShop(embedCtx.App.Config(), vats), nil
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
