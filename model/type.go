package model

import (
	"github.com/shopspring/decimal"
)

// var (
// 	ErrNotSameWeightUnit = errors.New("weights need to have same unit")
// 	ErrNotSameCurrency   = errors.New("Moneys need to have same currency")
// )

type StringMap map[string]string

type Money struct {
	Amount   *decimal.Decimal
	Currency string
}

type Weight struct {
	Weight     float32
	WeightUnit string
}
