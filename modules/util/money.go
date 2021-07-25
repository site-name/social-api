package util

import (
	"github.com/shopspring/decimal"
	goprices "github.com/site-name/go-prices"
)

// ZeroMoney returns zero money with currency unit is given currency
func ZeroMoney(currency string) (*goprices.Money, error) {
	return goprices.NewMoney(&decimal.Zero, currency)
}

// ZeroTaxedMoney returns zero-taxed money with currency unit of given currency
func ZeroTaxedMoney(currency string) (*goprices.TaxedMoney, error) {
	zero, err := ZeroMoney(currency)
	if err != nil {
		return nil, err
	}
	return goprices.NewTaxedMoney(zero, zero)
}
