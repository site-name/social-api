package model

import (
	"strings"

	"github.com/site-name/decimal"
)

type OpenExchangeRate struct {
	Id         string           `json:"id"`
	ToCurrency string           `json:"to_currency"` // db_index
	Rate       *decimal.Decimal `json:"rate"`        // default 0
}

func (o *OpenExchangeRate) IsValid() *AppError {
	outer := CreateAppErrorForModel(
		"model.open_exchange.is_valid.%s.app_error",
		"open_exchange_id=",
		"OpenExchangeRate.IsValid",
	)

	if !IsValidId(o.Id) {
		return outer("id", nil)
	}
	if len(o.ToCurrency) > CURRENCY_CODE_MAX_LENGTH {
		return outer("to_currency", &o.Id)
	}
	if o.Rate.LessThan(decimal.Zero) || o.Rate == nil {
		return outer("rate", &o.Id)
	}

	return nil
}

func (o *OpenExchangeRate) PreSave() {
	if o.Id == "" {
		o.Id = NewId()
	}
	o.commonPre()
}

func (o *OpenExchangeRate) commonPre() {
	if o.ToCurrency != "" {
		o.ToCurrency = strings.ToUpper(o.ToCurrency)
	}
	if o.Rate == nil {
		o.Rate = &decimal.Zero
	}
}

func (o *OpenExchangeRate) PreUpdate() {
	o.commonPre()
}
