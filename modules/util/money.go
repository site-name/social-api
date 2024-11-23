package util

import (
	goprices "github.com/site-name/go-prices"
	"github.com/sitename/sitename/model"
)

// ZeroMoney returns zero money with currency unit is given currency
func ZeroMoney(currency model.Currency) (*goprices.Money, error) {
	return goprices.NewMoney(0, currency.String())
}

// ZeroTaxedMoney returns zero-taxed money with currency unit of given currency
func ZeroTaxedMoney(currency model.Currency) (*goprices.TaxedMoney, error) {
	zero, err := goprices.NewMoney(0, currency.String())
	if err != nil {
		return nil, err
	}
	return goprices.NewTaxedMoney(*zero, *zero)
}

func ZeroMoneyRange(currency model.Currency) (*goprices.MoneyRange, error) {
	zero, err := goprices.NewMoney(0, currency.String())
	if err != nil {
		return nil, err
	}
	return goprices.NewMoneyRange(*zero, *zero)
}

func ZeroTaxedMoneyRange(currency model.Currency) (*goprices.TaxedMoneyRange, error) {
	zero, err := ZeroTaxedMoney(currency)
	if err != nil {
		return nil, err
	}
	return goprices.NewTaxedMoneyRange(*zero, *zero)
}
