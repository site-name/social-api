package model

import (
	"io"
	"time"

	"github.com/shopspring/decimal"
	"github.com/sitename/sitename/modules/json"
)

const (
	GIFT_CARD_CODE_MAX_LENGTH = 16
)

type GiftCard struct {
	Id                   string           `json:"id"`
	UserID               string           `json:"user_id"`
	Code                 string           `json:"code"`
	CreateAt             int64            `json:"created_at"`
	StartDate            *time.Time       `json:"start_date"`
	EndDate              *time.Time       `json:"end_date"`
	LastUsedOn           int64            `json:"last_used_on"`
	IsActive             bool             `json:"is_active"`
	Currency             string           `json:"currency"`
	InitialBalanceAmount *decimal.Decimal `json:"initial_balance_amount"`
	InitialBalance       *Money           `json:"initial_balance" db:"-"`
	CurrentBalanceAmount *decimal.Decimal `json:"current_balance_amount"`
	CurrentBalance       *Money           `json:"current_balance" db:"-"`
}

func (gc *GiftCard) DisplayCode() string {
	return "****" + gc.Code[len(gc.Code)-4:]
}

func (gc *GiftCard) ToJson() string {
	b, _ := json.JSON.Marshal(gc)
	return string(b)
}

func GiftCardFromJson(data io.Reader) *GiftCard {
	var gc GiftCard
	err := json.JSON.NewDecoder(data).Decode(&gc)
	if err != nil {
		return nil
	}
	return &gc
}
