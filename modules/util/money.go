package util

import (
	goprices "github.com/site-name/go-prices"
)

// ZeroMoney returns zero money with currency unit is given currency
func ZeroMoney(currency string) (*goprices.Money, error) {
	return goprices.NewMoney(0, currency)
}

// ZeroTaxedMoney returns zero-taxed money with currency unit of given currency
func ZeroTaxedMoney(currency string) (*goprices.TaxedMoney, error) {
	zero, err := goprices.NewMoney(0, currency)
	if err != nil {
		return nil, err
	}
	return goprices.NewTaxedMoney(*zero, *zero)
}

func ZeroMoneyRange(currency string) (*goprices.MoneyRange, error) {
	zero, err := goprices.NewMoney(0, currency)
	if err != nil {
		return nil, err
	}
	return goprices.NewMoneyRange(*zero, *zero)
}

func ZeroTaxedMoneyRange(currency string) (*goprices.TaxedMoneyRange, error) {
	zero, err := ZeroTaxedMoney(currency)
	if err != nil {
		return nil, err
	}
	return goprices.NewTaxedMoneyRange(*zero, *zero)
}
