package vatlayer

import (
	// "github.com/sitename/sitename/model/plugins"
	"strings"

	"github.com/sitename/sitename/app"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/modules/plugins"
)

var (
	_ plugins.BasePluginInterface = (*VatlayerPlugin)(nil)
)

func init() {
	plugins.RegisterVatlayerPlugin(func(srv *app.Server) plugins.BasePluginInterface {
		vp := &VatlayerPlugin{
			BasePlugin: plugins.BasePlugin{
				Manifest: plugins.PluginManifest{
					ID:                 "sitename.taxes.vatlayer",
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
				},
			},
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

		vp.Config = VatlayerConfiguration{
			AccessKey:     configuration["Access key"].(string),
			OriginCountry: originCountry,
			// ExcludedCountries: ,
		}

		return vp
	})
}

type VatlayerPlugin struct {
	plugins.BasePlugin

	Config      VatlayerConfiguration
	cachedTaxes model.StringInterface
}
