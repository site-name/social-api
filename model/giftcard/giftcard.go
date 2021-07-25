package giftcard

import (
	"io"
	"strings"
	"time"

	"github.com/shopspring/decimal"
	goprices "github.com/site-name/go-prices"
	"github.com/sitename/sitename/model"
	"golang.org/x/text/currency"
)

const (
	GIFT_CARD_CODE_MAX_LENGTH = 16
)

type GiftCard struct {
	Id                   string           `json:"id"`
	UserID               *string          `json:"user_id"`
	Code                 string           `json:"code"` // unique, db_index
	CreateAt             int64            `json:"created_at"`
	StartDate            *time.Time       `json:"start_date"`
	EndDate              *time.Time       `json:"end_date"`
	LastUsedOn           int64            `json:"last_used_on"`
	IsActive             *bool            `json:"is_active"`
	Currency             string           `json:"currency"`
	InitialBalanceAmount *decimal.Decimal `json:"initial_balance_amount"`
	InitialBalance       *goprices.Money  `json:"initial_balance,omitempty" db:"-"`
	CurrentBalanceAmount *decimal.Decimal `json:"current_balance_amount"`
	CurrentBalance       *goprices.Money  `json:"current_balance,omitempty" db:"-"`
}

func (gc *GiftCard) DisplayCode() string {
	return "****" + gc.Code[len(gc.Code)-4:]
}

func (gc *GiftCard) ToJson() string {
	gc.PopulateNonDbFields()
	return model.ModelToJson(gc)
}

// PopulateNonDbFields populates money fields for giftcard
func (gc *GiftCard) PopulateNonDbFields() {
	money := gc.InitialBalanceAmount
	if money == nil {
		money = &decimal.Zero
	}
	gc.InitialBalance = &goprices.Money{
		Amount:   money,
		Currency: gc.Currency,
	}

	money = gc.CurrentBalanceAmount
	if money == nil {
		money = &decimal.Zero
	}
	gc.CurrentBalance = &goprices.Money{
		Amount:   money,
		Currency: gc.Currency,
	}
}

func GiftCardFromJson(data io.Reader) *GiftCard {
	var gc GiftCard
	model.ModelFromJson(&gc, data)
	return &gc
}

func (gc *GiftCard) IsValid() *model.AppError {
	outer := model.CreateAppErrorForModel(
		"model.gift_card.is_valid.%s.app_error",
		"gift_card_id=",
		"GiftCard.IsValid",
	)

	if !model.IsValidId(gc.Id) {
		return outer("id", nil)
	}
	if gc.UserID != nil && !model.IsValidId(*gc.UserID) {
		return outer("user_id", &gc.Id)
	}
	if gc.CreateAt == 0 {
		return outer("create_at", &gc.Id)
	}
	if gc.LastUsedOn == 0 {
		return outer("last_used_on", &gc.Id)
	}
	if len(gc.Code) > GIFT_CARD_CODE_MAX_LENGTH {
		return outer("code", &gc.Id)
	}
	if unit, err := currency.ParseISO(gc.Currency); err != nil || !strings.EqualFold(unit.String(), gc.Currency) {
		return outer("currency", &gc.Id)
	}

	return nil
}

func (gc *GiftCard) PreSave() {
	if gc.Id == "" {
		gc.Id = model.NewId()
	}
	gc.CreateAt = model.GetMillis()
	if gc.IsActive == nil {
		gc.IsActive = model.NewBool(true)
	}
	if gc.StartDate == nil {
		today := time.Now()
		gc.StartDate = &today
	}
	if gc.Currency == "" {
		gc.Currency = model.DEFAULT_CURRENCY
	} else {
		gc.Currency = strings.ToUpper(strings.TrimSpace(gc.Currency))
	}
}
