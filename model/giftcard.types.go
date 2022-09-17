package model

import "github.com/site-name/decimal"

type BalanceObject struct {
	Giftcard        GiftCard
	PreviousBalance *decimal.Decimal
}

type BalanceData []BalanceObject
