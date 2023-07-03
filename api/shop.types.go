package api

import (
	"context"
	"fmt"
	"net/http"
	"unsafe"

	"github.com/Masterminds/squirrel"
	"github.com/gosimple/slug"
	"github.com/samber/lo"
	"github.com/sitename/sitename/app"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/store"
	"github.com/sitename/sitename/web"
)

var (
	shopLanguages = lo.MapToSlice(model.Languages, func(code model.LanguageCodeEnum, name string) *LanguageDisplay {
		return &LanguageDisplay{
			Code:     code,
			Language: name,
		}
	})
)

type Shop struct {
	Countries []*CountryDisplay `json:"countries"`
	// NOTE: Refer to ./schemas/shop.graphql for details on directives used
	DefaultMailSenderName *string `json:"defaultMailSenderName"`
	// NOTE: Refer to ./schemas/shop.graphql for details on directives used
	DefaultMailSenderAddress *string `json:"defaultMailSenderAddress"`
	// NOTE: Refer to ./schemas/shop.graphql for details on directives used
	AutomaticFulfillmentDigitalProducts *bool `json:"automaticFulfillmentDigitalProducts"`
	// NOTE: Refer to ./schemas/shop.graphql for details on directives used
	DefaultDigitalMaxDownloads *int32 `json:"defaultDigitalMaxDownloads"`
	// NOTE: Refer to ./schemas/shop.graphql for details on directives used
	DefaultDigitalURLValidDays *int32 `json:"defaultDigitalUrlValidDays"`

	Description             *string            `json:"description"`
	Domain                  *Domain            `json:"domain"`
	Languages               []*LanguageDisplay `json:"languages"`
	Name                    string             `json:"name"`
	Permissions             []*Permission      `json:"permissions"`
	PhonePrefixes           []string           `json:"phonePrefixes"`
	HeaderText              *string            `json:"headerText"`
	IncludeTaxesInPrices    bool               `json:"includeTaxesInPrices"`
	FulfillmentAutoApprove  bool               `json:"fulfillmentAutoApprove"`
	FulfillmentAllowUnpaid  bool               `json:"fulfillmentAllowUnpaid"`
	DisplayGrossPrices      bool               `json:"displayGrossPrices"`
	ChargeTaxesOnShipping   bool               `json:"chargeTaxesOnShipping"`
	TrackInventoryByDefault *bool              `json:"trackInventoryByDefault"`
	DefaultWeightUnit       *WeightUnitsEnum   `json:"defaultWeightUnit"`
	Translation             *ShopTranslation   `json:"translation"`
	CompanyAddress          *Address           `json:"companyAddress"`
	CustomerSetPasswordURL  *string            `json:"customerSetPasswordUrl"`
	Limits                  *LimitInfo         `json:"limits"`
	Version                 string             `json:"version"`

	// DefaultCountry                      *CountryDisplay    `json:"defaultCountry"`
	// StaffNotificationRecipients         []*StaffNotificationRecipient `json:"staffNotificationRecipients"`
	// AvailablePaymentGateways            []*PaymentGateway             `json:"availablePaymentGateways"`
	// AvailableExternalAuthentications    []*ExternalAuthentication     `json:"availableExternalAuthentications"`
	// AvailableShippingMethods            []*ShippingMethod             `json:"availableShippingMethods"`
	// ChannelCurrencies                   []string                      `json:"channelCurrencies"`
}

func systemConfigToGraphqlShop(cfg *model.Config, vats []*model.Vat) *Shop {
	taxes := lo.SliceToMap(vats, func(item *model.Vat) (model.CountryCode, *model.Vat) {
		return item.CountryCode, item
	})

	shopCountries := lo.MapToSlice(model.Countries, func(code model.CountryCode, name string) *CountryDisplay {
		return &CountryDisplay{
			Code:    code.String(),
			Country: name,
			Vat:     systemVatToGraphqlVat(taxes[code]),
		}
	})

	return &Shop{
		Countries:                           shopCountries,
		Languages:                           shopLanguages,
		DefaultMailSenderName:               cfg.ShopSettings.DefaultMailSenderName,
		DefaultMailSenderAddress:            cfg.ShopSettings.DefaultMailSenderAddress,
		AutomaticFulfillmentDigitalProducts: cfg.ShopSettings.AutomaticFulfillmentDigitalProducts,
		DefaultDigitalMaxDownloads:          (*int32)(unsafe.Pointer(cfg.ShopSettings.DefaultDigitalMaxDownloads)),
		DefaultDigitalURLValidDays:          (*int32)(unsafe.Pointer(cfg.ShopSettings.DefaultDigitalUrlValidDays)),
		IncludeTaxesInPrices:                *cfg.ShopSettings.IncludeTaxesInPrice,
		FulfillmentAutoApprove:              *cfg.ShopSettings.FulfillmentAutoApprove,
		DisplayGrossPrices:                  *cfg.ShopSettings.DisplayGrossPrices,
		ChargeTaxesOnShipping:               *cfg.ShopSettings.ChargeTaxesOnShipping,
		TrackInventoryByDefault:             cfg.ShopSettings.TrackInventoryByDefault,
		DefaultWeightUnit:                   cfg.ShopSettings.DefaultWeightUnit,
		CompanyAddress:                      SystemAddressToGraphqlAddress(cfg.ShopSettings.Address),
		CustomerSetPasswordURL:              cfg.ShopSettings.CustomerSetPasswordUrl,
		Version:                             model.CurrentVersion,
		Name:                                *cfg.ServiceSettings.SiteName,
		PhonePrefixes:                       []string{"84"},
		Description:                         cfg.ShopSettings.Description,

		Domain:      nil, // TODO: fix me
		Permissions: nil, // TODO: fix me
		HeaderText:  nil, // TODO: fix me
		Limits:      nil, // TODO: fix me
		Translation: nil, // TODO: fix me
	}
}

type PaymentGateway struct {
	Name       string               `json:"name"`
	ID         string               `json:"id"`
	Config     []*GatewayConfigLine `json:"config"`
	Currencies []string             `json:"currencies"`
}

func (s *Shop) AvailablePaymentGateways(ctx context.Context, args struct {
	Currency  string
	ChannelId string
}) ([]*PaymentGateway, error) {
	embedCtx := GetContextValue[*web.Context](ctx, WebCtx)
	pluginMng := embedCtx.App.Srv().PluginService().GetPluginManager()

	paymentGateWays := pluginMng.ListPaymentGateways(args.Currency, nil, args.ChannelId, true)

	return lo.Map(paymentGateWays, func(gw *model.PaymentGateway, _ int) *PaymentGateway {
		gw.Config = lo.Filter(gw.Config, func(cf model.StringInterface, _ int) bool { return cf != nil && len(cf) > 0 })

		resConfig := lo.Map(gw.Config, func(cf model.StringInterface, _ int) *GatewayConfigLine {
			var res GatewayConfigLine
			for k, v := range cf {
				vStr := fmt.Sprintf("%v", v)
				res.Field = k
				res.Value = &vStr
			}

			return &res
		})

		return &PaymentGateway{
			ID:         gw.Id,
			Name:       gw.Name,
			Currencies: gw.Currencies,
			Config:     resConfig,
		}
	}), nil
}

func (s *Shop) AvailableExternalAuthentications(ctx context.Context) ([]ExternalAuthentication, error) {
	embedCtx := GetContextValue[*web.Context](ctx, WebCtx)
	pluginMng := embedCtx.App.Srv().PluginService().GetPluginManager()

	auths, appErr := pluginMng.ListExternalAuthentications(true)
	if appErr != nil {
		return nil, appErr
	}
	auths = lo.Filter(auths, func(auth model.StringInterface, _ int) bool { return auth != nil && len(auth) > 0 })

	return lo.Map(auths, func(auth model.StringInterface, _ int) ExternalAuthentication {
		var res ExternalAuthentication
		for k, v := range auth {
			vStr := fmt.Sprintf("%v", v)
			res.ID = k
			res.Name = &vStr
		}

		return res
	}), nil
}

func (s *Shop) AvailableShippingMethods(ctx context.Context, args struct {
	Channel string // NOTE: channel slug
	Address *AddressInput
}) ([]*ShippingMethod, error) {
	// validate argument(s)
	if !slug.IsSlug(args.Channel) {
		return nil, model.NewAppError("Shop.AvailableShippingMethods", app.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "channel"}, args.Channel+" is not a valid channel slug", http.StatusBadRequest)
	}
	if args.Address != nil {
		err := args.Address.Validate()
		if err != nil {
			return nil, err
		}
	}

	var (
		address                  model.Address
		embedCtx                 = GetContextValue[*web.Context](ctx, WebCtx)
		shippingMethodFilterOpts = &model.ShippingMethodFilterOption{
			ShippingZoneChannelSlug:    squirrel.Eq{store.ChannelTableName + ".Slug": args.Channel},
			ChannelListingsChannelSlug: squirrel.Eq{store.ChannelTableName + ".Slug": args.Channel},
		}
	)

	if args.Address != nil && args.Address.Country != nil {
		shippingMethodFilterOpts.ShippingZoneCountries = squirrel.ILike{store.ShippingZoneTableName + ".Countries": "%" + args.Address.Country.String() + "%"}
	}

	availableSippingMethods, appErr := embedCtx.
		App.
		Srv().
		ShippingService().
		ShippingMethodsByOptions(shippingMethodFilterOpts)
	if appErr != nil {
		return nil, appErr
	}

	if args.Address != nil && args.Address.Country != nil {
		args.Address.PatchAddress(&address)
		availableSippingMethods = embedCtx.App.Srv().
			ShippingService().
			FilterShippingMethodsByPostalCodeRules(availableSippingMethods, &address)
	}

	if len(availableSippingMethods) == 0 {
		return []*ShippingMethod{}, nil
	}

	shippingMapping, appErr := embedCtx.App.Srv().ShippingService().GetShippingMethodToShippingPriceMapping(availableSippingMethods, args.Channel)
	if appErr != nil {
		return nil, appErr
	}

	channel, appErr := embedCtx.App.Srv().ChannelService().ChannelByOption(&model.ChannelFilterOption{
		Slug: squirrel.Eq{store.ChannelTableName + ".Slug": args.Channel},
	})
	if appErr != nil {
		return nil, appErr
	}

	pluginMng := embedCtx.App.Srv().PluginService().GetPluginManager()

	for _, shippingMethod := range availableSippingMethods {
		shippingPrice := shippingMapping[shippingMethod.Id]
		taxedPrice, appErr := pluginMng.ApplyTaxesToShipping(*shippingPrice, address, channel.Id)
		if appErr != nil {
			return nil, appErr
		}

		if *embedCtx.App.Config().ShopSettings.DisplayGrossPrices {
			shippingMethod.SetPrice(taxedPrice.Gross)
		} else {
			shippingMethod.SetPrice(taxedPrice.Net)
		}
	}

	return systemRecordsToGraphql(availableSippingMethods, SystemShippingMethodToGraphqlShippingMethod), nil
}

// NOTE: Refer to ./schemas/shop.graphql for details on directives used
func (s *Shop) ChannelCurrencies(ctx context.Context) ([]string, error) {
	embedCtx := GetContextValue[*web.Context](ctx, WebCtx)

	channels, appErr := embedCtx.App.Srv().ChannelService().ChannelsByOption(&model.ChannelFilterOption{})
	if appErr != nil {
		return nil, appErr
	}

	return channels.Currencies(), nil
}

func (s *Shop) DefaultCountry(ctx context.Context) (*CountryDisplay, error) {
	embedCtx := GetContextValue[*web.Context](ctx, WebCtx)
	defaultCountryName := model.Countries[model.DEFAULT_COUNTRY]

	if defaultCountryName != "" {
		vats, err := embedCtx.App.Srv().Store.Vat().FilterByOptions(&model.VatFilterOptions{
			CountryCode: squirrel.Eq{store.VatTableName + ".CountryCode": model.DEFAULT_COUNTRY},
		})
		if err != nil {
			return nil, model.NewAppError("Shop.DefaultCountry", "app.shop.error_finding_vats.app_error", nil, err.Error(), http.StatusInternalServerError)
		}

		res := &CountryDisplay{
			Code:    model.DEFAULT_COUNTRY.String(),
			Country: defaultCountryName,
		}
		if len(vats) > 0 {
			res.Vat = systemVatToGraphqlVat(vats[0])
		}

		return res, nil
	}

	return nil, nil
}

// NOTE: Refer to ./schemas/shop.graphql for details on directives used
func (s *Shop) StaffNotificationRecipients(ctx context.Context) ([]*StaffNotificationRecipient, error) {
	embedCtx := GetContextValue[*web.Context](ctx, WebCtx)

	notificationRecipients, err := embedCtx.App.
		Srv().
		Store.
		StaffNotificationRecipient().
		FilterByOptions(&model.StaffNotificationRecipientFilterOptions{}) // find all
	if err != nil {
		return nil, model.NewAppError("Shop.StaffNotificationRecipients", "app.account.staff_notification_recipients_by_options.app_error", nil, err.Error(), http.StatusInternalServerError)
	}

	return systemRecordsToGraphql(notificationRecipients, systemStaffNotificationRecipientToGraphqlStaffNotificationRecipient), nil
}

// -------------------- Vat ------------------------

type Vat struct {
	CountryCode  string         `json:"countryCode"`
	StandardRate *float64       `json:"standardRate"`
	ReducedRates []*ReducedRate `json:"reducedRates"`

	vat *model.Vat
}

type ReducedRate struct {
	Rate     float64 `json:"rate"`
	RateType string  `json:"rateType"`
}

func systemVatToGraphqlVat(vat *model.Vat) *Vat {
	if vat == nil {
		return nil
	}

	res := &Vat{
		CountryCode: vat.CountryCode.String(),
		vat:         vat,
	}

	stdRate, ok := vat.Data.Get("standard_rate", 0.0).(float64)
	if ok {
		res.StandardRate = &stdRate
	}
	rdcRate, ok := vat.Data.Get("reduced_rates", map[string]float64{}).(map[string]float64)
	if ok && len(rdcRate) > 0 {
		res.ReducedRates = lo.MapToSlice(rdcRate, func(key string, value float64) *ReducedRate {
			return &ReducedRate{
				Rate:     value,
				RateType: key,
			}
		})
	}

	return res
}
