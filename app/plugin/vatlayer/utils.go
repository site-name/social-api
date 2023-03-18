package vatlayer

import (
	goprices "github.com/site-name/go-prices"
	"github.com/sitename/sitename/modules/util"
)

type VatlayerConfiguration struct {
	AccessKey           string
	ExcludedCountries   util.AnyArray[string]
	CountriesFromOrigin util.AnyArray[string]
	OriginCountry       string
}

// Naively convert Money to TaxedMoney.
//
// It is meant for consistency with price handling logic across the codebase,
// passthrough other money types.
func convertToNaiveTaxedMoney[M goprices.MoneyObject](base M, taxes M, rateName string) any {
	if base == nil {
		return nil
	}

	switch t := any(base).(type) {
	case *goprices.TaxedMoney, *goprices.TaxedMoneyRange:
		return t

	case *goprices.Money:
		return &goprices.TaxedMoney{
			Net:   t,
			Gross: t,
		}

	default:
		return nil
	}
}

// func ApplyTaxToPrice(taxes)
