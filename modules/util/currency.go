package util

import (
	"sort"

	"github.com/samber/lo"
	goprices "github.com/site-name/go-prices"
	"golang.org/x/text/currency"
	"golang.org/x/text/language"
)

// GetCurrencyForCountry returns currency code for givent country_code
//
// eg: us -> USD
//
// returns empty string if givent country_code is invalid
func GetCurrencyForCountry(country_code string) string {
	region, err := language.ParseRegion(country_code)
	if err != nil {
		return ""
	}
	unit, ok := currency.FromRegion(region)
	if ok {
		return unit.String()
	}
	return ""
}

// MinMaxMoneyInMoneySlice takes a list of moneys, compare them and returns min, max moneys respectively.
//
// NOTE: moneys must have same currency
func MinMaxMoneyInMoneySlice(moneys []*goprices.Money) (min *goprices.Money, max *goprices.Money) {
	if len(moneys) == 0 {
		return nil, nil
	}
	if len(moneys) == 1 {
		return moneys[0], moneys[0]
	}

	moneys = lo.Filter(moneys, func(v *goprices.Money, _ int) bool {
		return v != nil
	})
	sort.Slice(moneys, func(i, j int) bool {
		return moneys[i].LessThan(moneys[j])
	})

	return moneys[0], moneys[len(moneys)-1]
}
