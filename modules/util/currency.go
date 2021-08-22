package util

import (
	"sort"

	goprices "github.com/site-name/go-prices"
	"golang.org/x/text/currency"
	"golang.org/x/text/language"
)

// GetCurrencyForCountry returns currency code for givent country_code
//
//eg: us -> USD
//
// returns empty string if givent country_code is invalid
func GetCurrencyForCountry(country_code string) string {
	rg, err := language.ParseRegion(country_code)
	if err != nil {
		return ""
	}
	unit, ok := currency.FromRegion(rg)
	if ok {
		return unit.String()
	}
	return ""
}

// MinMaxMoneyInMoneySlice takes a list of moneys, compare them and returns min, max moneys respectively.
//
// NOTE: moneys must have same currency
func MinMaxMoneyInMoneySlice(moneys []*goprices.Money) (min *goprices.Money, max *goprices.Money) {

	if len(moneys) == 1 {
		return moneys[0], moneys[0]
	}

	sort.Slice(moneys, func(i, j int) bool {
		less, err := moneys[i].LessThan(moneys[j])
		return less && err == nil
	})

	return moneys[0], moneys[len(moneys)-1]
}
