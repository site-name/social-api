package util

import (
	goprices "github.com/site-name/go-prices"
	"golang.org/x/text/currency"
	"golang.org/x/text/language"
)

// GetCurrencyForCountry returns currency code for givent country_code
//
// eg: us -> USD
//
// returns empty string if given countryCode is invalid
func GetCurrencyForCountry(countryCode string) string {
	region, err := language.ParseRegion(countryCode)
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
		return
	}

	for _, money := range moneys {
		if min == nil || money.LessThan(*min) {
			min = money
		}
		if max == nil || !money.LessThanOrEqual(*max) {
			max = money
		}
	}

	return
}
