package vatlayer

import (
	"strings"

	goprices "github.com/site-name/go-prices"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model/account"
	"github.com/sitename/sitename/model/checkout"
	"github.com/sitename/sitename/model/product_and_discount"
	"github.com/sitename/sitename/modules/plugins"
)

const (
	pluginID = "sitename.taxes.vatlayer"
)

var (
	_ plugins.BasePluginInterface = (*VatlayerPlugin)(nil)
)

type VatlayerPlugin struct {
	*plugins.BasePlugin

	config      VatlayerConfiguration
	cachedTaxes model.StringInterface
}

func init() {
	plugins.RegisterVatlayerPlugin(func(cfg plugins.NewPluginConfig) plugins.BasePluginInterface {

		basePlg := plugins.NewBasePlugin(cfg)
		basePlg.Manifest = plugins.PluginManifest{
			ID:                 pluginID,
			Name:               "Vatlayer",
			MetaCodeKey:        "vatlayer.code",
			MetaDescriptionKey: "vatlayer.description",
			ConfigStructure: map[string]model.StringInterface{
				"origin_country": {
					"type":      plugins.STRING,
					"help_test": "Country code in ISO format, required to calculate taxes for countries from `Countries for which taxes will be calculated from origin country`.",
					"label":     "Origin country",
				},
				"countries_to_calculate_taxes_from_origin": {
					"type":      plugins.STRING,
					"help_text": "List of destination countries (separated by comma), in ISO format which will use origin country to calculate taxes.",
					"label":     "Countries for which taxes will be calculated from origin country",
				},
				"excluded_countries": {
					"type":      plugins.STRING,
					"help_text": "List of countries (separated by comma), in ISO format for which no VAT should be added.",
					"label":     "Countries for which no VAT will be added.",
				},
				"Access key": {
					"type":      plugins.PASSWORD,
					"help_text": "Required to authenticate to Vatlayer API.",
					"label":     "Access key",
				},
			},
			DefaultConfiguration: []model.StringInterface{
				{"name": "Access key", "value": nil},
				{"name": "origin_country", "value": nil},
				{"name": "countries_to_calculate_taxes_from_origin", "value": nil},
				{"name": "excluded_countries", "value": nil},
			},
		}

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

	}, pluginID)
}

// previousValue must be either: *Decimal, *TaxedMoney, *TaxedMoneyRange
func (vp *VatlayerPlugin) skipPlugin(previousValue interface{}) bool {
	if !vp.Active || vp.config.AccessKey == "" {
		return true
	}

	// The previous plugin already calculated taxes so we can skip our logic
	if taxedMoneyRange, ok := previousValue.(*goprices.TaxedMoneyRange); ok && taxedMoneyRange != nil {
		equal1, err1 := taxedMoneyRange.Start.Net.Equal(taxedMoneyRange.Start.Gross)
		equal2, err2 := taxedMoneyRange.Stop.Net.Equal(taxedMoneyRange.Stop.Gross)

		return err1 == nil && err2 == nil && !equal1 && !equal2
	}

	if taxedMoney, ok := previousValue.(*goprices.TaxedMoney); ok && taxedMoney != nil {
		equal, err := taxedMoney.Net.Equal(taxedMoney.Gross)
		return err == nil && !equal
	}

	return false
}

func (vp *VatlayerPlugin) CalculateCheckoutTotal(checkoutInfo checkout.CheckoutInfo, lines checkout.CheckoutLineInfos, address *account.Address, discounts []*product_and_discount.DiscountInfo, previousValue goprices.TaxedMoney) (*goprices.TaxedMoney, *plugins.PluginMethodNotImplemented) {
	panic("not implemented")
}
