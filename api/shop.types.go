package api

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Masterminds/squirrel"
	"github.com/samber/lo"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/store"
	"github.com/sitename/sitename/web"
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

func systemShopSettingsToGraphqlShop(shop *model.ShopSettings) *Shop {
	if shop == nil {
		return nil
	}

	res := &Shop{}

	return res
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
			var key, value string
			for k, v := range cf {
				key = k
				value = fmt.Sprintf("%v", v)
			}

			return &GatewayConfigLine{
				Field: key,
				Value: &value,
			}
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
		var key, value string

		for k, v := range auth {
			key = k
			value = fmt.Sprintf("%v", v)
		}

		return ExternalAuthentication{
			ID:   key,
			Name: &value,
		}
	}), nil
}

func (s *Shop) AvailableShippingMethods(ctx context.Context, args struct {
	Channel string
	Address *AddressInput
}) ([]*PaymentGateway, error) {
	panic("not implemented")
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
