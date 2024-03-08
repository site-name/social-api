package vatlayer

import (
	"net/http"
	"strings"

	goprices "github.com/site-name/go-prices"
	"github.com/sitename/sitename/app/plugin"
	"github.com/sitename/sitename/app/plugin/interfaces"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model_helper"
	"github.com/sitename/sitename/modules/model_types"
)

var manifest = &interfaces.PluginManifest{
	PluginID:           "sitename.taxes.vatlayer",
	PluginName:         "Vatlayer",
	MetaCodeKey:        "vatlayer.code",
	MetaDescriptionKey: "vatlayer.description",
	DefaultConfiguration: []model_types.JSONString{
		{"name": "Access key", "value": nil},
		{"name": "origin_country", "value": nil},
		{"name": "countries_to_calculate_taxes_from_origin", "value": nil},
		{"name": "excluded_countries", "value": nil},
	},
	ConfigStructure: map[string]model_types.JSONString{
		"origin_country": {
			"type":      interfaces.STRING,
			"help_text": "Country code in ISO format, required to calculate taxes for countries from `Countries for which taxes will be calculated from origin country`.",
			"label":     "Origin country",
		},
		"countries_to_calculate_taxes_from_origin": {
			"type":      interfaces.STRING,
			"help_text": "List of destination countries (separated by comma), in ISO format which will use origin country to calculate taxes.",
			"label":     "Countries for which taxes will be calculated from origin country",
		},
		"excluded_countries": {
			"type":      interfaces.STRING,
			"help_text": "List of countries (separated by comma), in ISO format for which no VAT should be added.",
			"label":     "Countries for which no VAT will be added.",
		},
		"Access key": {
			"type":      interfaces.PASSWORD,
			"help_text": "Required to authenticate to Vatlayer API.",
			"label":     "Access key",
		},
	},
}

// type check
var _ interfaces.BasePluginInterface = (*VatlayerPlugin)(nil)

type VatlayerPlugin struct {
	plugin.BasePlugin

	config      VatlayerConfiguration
	cachedTaxes model_types.JSONString
}

func initFunc(cfg *plugin.PluginConfig) interfaces.BasePluginInterface {
	vatPlugin := &VatlayerPlugin{
		BasePlugin: *plugin.NewBasePlugin(cfg),
	}

	var configuration = model_types.JSONString{}
	for _, item := range vatPlugin.Configuration {
		configuration[item.Get("name", "").(string)] = item["value"]
	}
	var originCountry = configuration.Get("origin_country", "").(string)
	if upper := strings.ToUpper(originCountry); model.CountryCode(upper).IsValid() {
		originCountry = upper
	} else {
		originCountry = ""
	}

	var countriesFromOrigin = configuration.Get("countries_to_calculate_taxes_from_origin", "").(string)
	var splitCountriesFromOrigin = []string{}

	for _, str := range strings.FieldsFunc(countriesFromOrigin, func(r rune) bool { return r == ' ' || r == ',' }) {
		if upper := strings.ToUpper(str); model.CountryCode(upper).IsValid() {
			splitCountriesFromOrigin = append(splitCountriesFromOrigin, upper)
		}
	}

	var excludedCountries = configuration.Get("excluded_countries", "").(string)
	var splitExcludedCountries = []string{}

	for _, str := range strings.FieldsFunc(excludedCountries, func(r rune) bool { return r == ' ' || r == ',' }) {
		if upper := strings.ToUpper(str); model.CountryCode(upper).IsValid() {
			splitExcludedCountries = append(splitExcludedCountries, upper)
		}
	}

	vatPlugin.config = VatlayerConfiguration{
		AccessKey:           configuration.Get("Access key", "").(string),
		OriginCountry:       originCountry,
		ExcludedCountries:   splitExcludedCountries,
		CountriesFromOrigin: splitCountriesFromOrigin,
	}

	vatPlugin.cachedTaxes = make(model_types.JSONString)
	return vatPlugin

}

func init() {
	plugin.RegisterPlugin(plugin.PluginInitObjType{
		NewPluginFunc: initFunc,
		Manifest:      manifest,
	})
}

// previousValue must be either *TaxedMoney or *TaxedMoneyRange
func (vp *VatlayerPlugin) skipPlugin(previousValue any) bool {
	if !vp.Active || vp.config.AccessKey == "" {
		return true
	}

	// The previous plugin already calculated taxes so we can skip our logic
	switch t := previousValue.(type) {
	case *goprices.TaxedMoneyRange:
		return !t.Start.Net.Equal(t.Start.Gross) && !t.Stop.Net.Equal(t.Stop.Gross)
	case goprices.TaxedMoneyRange:
		return !t.Start.Net.Equal(t.Start.Gross) && !t.Stop.Net.Equal(t.Stop.Gross)

	case *goprices.TaxedMoney:
		return !t.Net.Equal(t.Gross)
	case goprices.TaxedMoney:
		return !t.Net.Equal(t.Gross)

	default:
		return false
	}
}

// previousValue must be either TaxedMoneyRange or TaxedMoney
func (vp *VatlayerPlugin) CalculateCheckoutTotal(checkoutInfo model_helper.CheckoutInfo, lines model_helper.CheckoutLineInfos, address *model.Address, discounts []*model_helper.DiscountInfo, previousValue goprices.TaxedMoney) (*goprices.TaxedMoney, *model_helper.AppError) {
	if vp.skipPlugin(previousValue) {
		return &previousValue, nil
	}

	checkoutSubTotal, appErr := vp.Manager.Srv.CheckoutService().CheckoutSubTotal(
		vp.Manager,
		checkoutInfo,
		lines,
		address,
		discounts,
	)
	if appErr != nil {
		return nil, appErr
	}

	checkoutShippingPrice, appErr := vp.Manager.Srv.CheckoutService().CheckoutShippingPrice(
		vp.Manager,
		checkoutInfo,
		lines,
		address,
		discounts,
	)
	if appErr != nil {
		return nil, appErr
	}

	sum, err := checkoutSubTotal.Add(checkoutShippingPrice)
	if err != nil {
		return nil, model_helper.NewAppError("CalculateCheckoutTotal", model_helper.ErrorCalculatingMoneyErrorID, nil, err.Error(), http.StatusInternalServerError)
	}

	checkoutInfo.Checkout.PopulateNonDbFields() // this is needed
	if checkoutInfo.Checkout.Discount != nil {
		sum, err = sum.Sub(checkoutInfo.Checkout.Discount)
		if err != nil {
			return nil, model_helper.NewAppError("CalculateCheckoutTotal", model_helper.ErrorCalculatingMoneyErrorID, nil, err.Error(), http.StatusInternalServerError)
		}
	}

	return sum, nil
}

// Try to fetch cached taxes on the plugin level.
//
// If the plugin doesn't have cached taxes for a given country it will fetch it
// from cache or db.
func (vp *VatlayerPlugin) getTaxesForCountry(country string) (any, *model_helper.AppError) {
	if country == "" {
		country = vp.config.OriginCountry
		if country == "" {
			country := vp.Manager.Srv.Config().ShopSettings.Address.Country.String()

			if country == "" {
				country = model.DEFAULT_COUNTRY.String()
			}
		}
	}

	if vp.config.CountriesFromOrigin.Contains(country) {
		country = vp.config.OriginCountry
	}
	if vp.config.ExcludedCountries.Contains(country) {
		return nil, nil
	}

	if tax, ok := vp.cachedTaxes[country]; ok {
		return tax, nil
	}

	panic("not implemented")
}
