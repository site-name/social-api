package model

import (
	"github.com/shopspring/decimal"
)

type StringMap map[string]string

type Money struct {
	Amount   *decimal.Decimal
	Currency string
}
