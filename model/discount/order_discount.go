package discount

import (
	"github.com/shopspring/decimal"
	"github.com/sitename/sitename/model"
)

const (
	ORDER_DISCOUNT_NAME_MAX_LENGTH = 255
)

type OrderDiscount struct {
	Id             string           `json:"id"`
	OrderID        string           `json:"order_id"`
	Type           string           `json:"type"`
	ValueType      string           `json:"value_type"`
	Value          *decimal.Decimal `json:"value"`
	AmountValue    *decimal.Decimal `json:"amount_value"`
	Amount         *model.Money     `json:"amount,omitempty" db:"-"`
	Currency       string           `json:""currency`
	Name           *string          `json:"name"`
	TranslatedName *string          `json:"translated_name"`
	Reason         *string          `json:"reason"`
}
