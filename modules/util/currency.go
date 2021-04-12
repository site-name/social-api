package util

import (
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

func ToLocalCurrency() {

}
