package vatlayer

import (
	"net/http"
	"strings"

	goprices "github.com/site-name/go-prices"
	"github.com/sitename/sitename/app"
	"github.com/sitename/sitename/app/plugin"
	"github.com/sitename/sitename/app/plugin/interfaces"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model/account"
	"github.com/sitename/sitename/model/checkout"
	"github.com/sitename/sitename/model/product_and_discount"
)

var (
	_ interfaces.BasePluginInterface = (*VatlayerPlugin)(nil)

	manifest = &interfaces.PluginManifest{
		PluginID:           "sitename.taxes.vatlayer",
		PluginName:         "Vatlayer",
		MetaCodeKey:        "vatlayer.code",
		MetaDescriptionKey: "vatlayer.description",
		DefaultConfiguration: []model.StringInterface{
			{"name": "Access key", "value": nil},
			{"name": "origin_country", "value": nil},
			{"name": "countries_to_calculate_taxes_from_origin", "value": nil},
			{"name": "excluded_countries", "value": nil},
		},
		ConfigStructure: map[string]model.StringInterface{
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
)

type VatlayerPlugin struct {
	*plugin.BasePlugin

	config      VatlayerConfiguration
	cachedTaxes model.StringInterface
}

func init() {
	plugin.RegisterVatlayerPlugin(func(cfg *plugin.NewPluginConfig) interfaces.BasePluginInterface {

		basePlg := plugin.NewBasePlugin(cfg)

		// override base plugin's manifest
		basePlg.Manifest = manifest

		vp := &VatlayerPlugin{
			BasePlugin: basePlg,
		}

		var configuration = map[string]interface{}{}
		for _, item := range vp.Configuration {
			configuration[item["name"].(string)] = item["value"]
		}
		var originCountry = configuration["origin_country"].(string)
		if upper := strings.ToUpper(originCountry); model.Countries[upper] != "" {
			originCountry = upper
		} else {
			originCountry = ""
		}

		var countriesFromOrigin = configuration["countries_to_calculate_taxes_from_origin"].(string)
		var splitCountriesFromOrigin = []string{}

		for _, str := range strings.Split(countriesFromOrigin, ",") {
			if upper := strings.ToUpper(strings.TrimSpace(str)); model.Countries[upper] != "" {
				splitCountriesFromOrigin = append(splitCountriesFromOrigin, upper)
			}
		}

		var excludedCountries = configuration["excluded_countries"].(string)
		var splitExcludedCountries = []string{}

		for _, str := range strings.Split(excludedCountries, ",") {
			if upper := strings.ToUpper(strings.TrimSpace(str)); model.Countries[upper] != "" {
				splitExcludedCountries = append(splitExcludedCountries, upper)
			}
		}

		vp.config = VatlayerConfiguration{
			AccessKey:           configuration["Access key"].(string),
			OriginCountry:       originCountry,
			ExcludedCountries:   splitExcludedCountries,
			CountriesFromOrigin: splitCountriesFromOrigin,
		}

		vp.cachedTaxes = make(model.StringInterface)
		return vp

	}, manifest)
}

// previousValue must be either: *Decimal, *TaxedMoney, *TaxedMoneyRange
func (vp *VatlayerPlugin) skipPlugin(previousValue interface{}) bool {
	if !vp.Active || vp.config.AccessKey == "" {
		return true
	}

	// The previous plugin already calculated taxes so we can skip our logic
	switch t := previousValue.(type) {
	case *goprices.TaxedMoneyRange:
		equal1, err1 := t.Start.Net.Equal(t.Start.Gross)
		equal2, err2 := t.Stop.Net.Equal(t.Stop.Gross)

		return err1 == nil && err2 == nil && !equal1 && !equal2

	case goprices.TaxedMoneyRange:
		equal1, err1 := t.Start.Net.Equal(t.Start.Gross)
		equal2, err2 := t.Stop.Net.Equal(t.Stop.Gross)

		return err1 == nil && err2 == nil && !equal1 && !equal2

	case *goprices.TaxedMoney:
		equal, err := t.Net.Equal(t.Gross)
		return err == nil && !equal

	case goprices.TaxedMoney:
		equal, err := t.Net.Equal(t.Gross)
		return err == nil && !equal

	default:
		return false
	}
}

// previousValue must be either TaxedMoneyRange or TaxedMoney
func (vp *VatlayerPlugin) CalculateCheckoutTotal(checkoutInfo checkout.CheckoutInfo, lines checkout.CheckoutLineInfos, address *account.Address, discounts []*product_and_discount.DiscountInfo, previousValue goprices.TaxedMoney) (*goprices.TaxedMoney, *model.AppError) {
	if vp.skipPlugin(previousValue) {
		return &previousValue, nil
	}

	checkoutInfo.Checkout.PopulateNonDbFields() // this is needed

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
		return nil, model.NewAppError("CalculateCheckoutTotal", app.ErrorCalculatingMoneyErrorID, nil, err.Error(), http.StatusInternalServerError)
	}

	sub, err := sum.Sub(checkoutInfo.Checkout.Discount)
	if err != nil {
		return nil, model.NewAppError("CalculateCheckoutTotal", app.ErrorCalculatingMoneyErrorID, nil, err.Error(), http.StatusInternalServerError)
	}

	return sub, nil
}

// Try to fetch cached taxes on the plugin level.
//
// If the plugin doesn't have cached taxes for a given country it will fetch it
// from cache or db.
func (vp *VatlayerPlugin) getTaxesForCountry(country string) {
	if country == "" {
		originCountryCode := vp.config.OriginCountry
		if originCountryCode == "" {

		}
	}
}
