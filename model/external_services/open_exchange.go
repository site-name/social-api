package external_services

import (
	"strings"

	"github.com/site-name/decimal"
	goprices "github.com/site-name/go-prices"
	"github.com/sitename/sitename/model"
)

type OpenExchangeRate struct {
	Id         string           `json:"id"`
	ToCurrency string           `json:"to_currency"` // db_index
	Rate       *decimal.Decimal `json:"rate"`
	UpdateAt   int64            `json:"update_at"`
}

func (o *OpenExchangeRate) IsValid() *model.AppError {
	outer := model.CreateAppErrorForModel(
		"model.open_exchange.is_valid.%s.app_error",
		"open_exchange_id=",
		"OpenExchangeRate.IsValid",
	)

	if !model.IsValidId(o.Id) {
		return outer("id", nil)
	}
	if len(o.ToCurrency) > model.CURRENCY_CODE_MAX_LENGTH ||
		goprices.CurrenciesMap[o.ToCurrency] == "" {
		return outer("to_currency", &o.Id)
	}
	if o.Rate.LessThan(decimal.Zero) {
		return outer("rate", &o.Id)
	}

	return nil
}

func (o *OpenExchangeRate) PreSave() {
	if o.Id == "" {
		o.Id = model.NewId()
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
	o.UpdateAt = model.GetMillis()
	o.commonPre()
}
