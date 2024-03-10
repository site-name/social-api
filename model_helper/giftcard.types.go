package model_helper

import (
	"github.com/site-name/decimal"
	"github.com/sitename/sitename/model"
)

type BalanceObject struct {
	Giftcard        model.Giftcard
	PreviousBalance *decimal.Decimal
}

type BalanceData []BalanceObject
