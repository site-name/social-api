package model

import (
	"net/http"
	"strings"

	"github.com/site-name/decimal"
	"gorm.io/gorm"
)

type OpenExchangeRate struct {
	Id         string           `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid();column:Id"`
	ToCurrency string           `json:"to_currency" gorm:"type:varchar(5);column:ToCurrency"` // db_index
	Rate       *decimal.Decimal `json:"rate" gorm:"column:Rate;default:0"`                    // default 0
}

func (c *OpenExchangeRate) BeforeCreate(_ *gorm.DB) error { c.commonPre(); return c.IsValid() }
func (c *OpenExchangeRate) BeforeUpdate(_ *gorm.DB) error { c.commonPre(); return c.IsValid() }
func (c *OpenExchangeRate) TableName() string             { return OpenExchangeRateTableName }

func (o *OpenExchangeRate) IsValid() *AppError {
	if o.Rate != nil && o.Rate.LessThan(decimal.Zero) {
		return NewAppError("OpenExchangeRate.IsValid", "model.open_exchange.is_valid.rate.app_error", nil, "rate cannot be less than zero", http.StatusBadRequest)
	}

	return nil
}

func (o *OpenExchangeRate) commonPre() {
	o.ToCurrency = strings.ToUpper(o.ToCurrency)
	if o.Rate == nil {
		o.Rate = GetPointerOfValue(decimal.Zero)
	}
}
