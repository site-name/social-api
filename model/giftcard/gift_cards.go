package giftcard

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/shopspring/decimal"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model/checkout"
	"github.com/sitename/sitename/modules/json"
	"golang.org/x/text/currency"
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
	InitialBalance       *checkout.Money  `json:"initial_balance" db:"-"`
	CurrentBalanceAmount *decimal.Decimal `json:"current_balance_amount"`
	CurrentBalance       *checkout.Money  `json:"current_balance" db:"-"`
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

func (gc *GiftCard) createAppError(fieldName string) *model.AppError {
	id := fmt.Sprintf("model.gift_card.is_valid.%s.app_error", fieldName)
	var details string
	if !strings.EqualFold(fieldName, "id") {
		details = "gift_card_id=" + gc.Id
	}

	return model.NewAppError("GiftCard.IsValid", id, nil, details, http.StatusBadRequest)
}

func (gc *GiftCard) IsValid() *model.AppError {
	if !model.IsValidId(gc.Id) {
		return gc.createAppError("id")
	}
	if !model.IsValidId(gc.UserID) {
		return gc.createAppError("user_id")
	}
	if gc.CreateAt == 0 {
		return gc.createAppError("create_at")
	}
	if gc.LastUsedOn == 0 {
		return gc.createAppError("last_used_on")
	}
	if len(gc.Code) > GIFT_CARD_CODE_MAX_LENGTH {
		return gc.createAppError("code")
	}
	if unit, err := currency.ParseISO(gc.Currency); err != nil || !strings.EqualFold(unit.String(), gc.Currency) {
		return gc.createAppError("currency")
	}

	return nil
}
