package vatlayer

import (
	"fmt"

	goprices "github.com/site-name/go-prices"
)

type VatlayerConfiguration struct {
	AccessKey           string
	ExcludedCountries   []string
	CountriesFromOrigin []string
	OriginCountry       string
}

// Naively convert Money to TaxedMoney.
//
// It is meant for consistency with price handling logic across the codebase,
// passthrough other money types.
func convertToNaiveTaxedMoney(base any, taxes any, rateName string) (any, error) {
	if base == nil {
		return nil, nil
	}

	switch t := base.(type) {
	case *goprices.TaxedMoney, *goprices.TaxedMoneyRange:
		return t, nil

	case *goprices.Money:
		return &goprices.TaxedMoney{
			Net:   t,
			Gross: t,
		}, nil

	case *goprices.MoneyRange:
		return nil, nil

	default:
		return nil, fmt.Errorf("unknown base for flat_tax: %T", base)
	}
}

// func ApplyTaxToPrice(taxes)
