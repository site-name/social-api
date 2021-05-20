package product_and_discount

import (
	"io"
	"strings"
	"unicode/utf8"

	"github.com/shopspring/decimal"
	"github.com/sitename/sitename/model"
	"golang.org/x/text/currency"
)

// max lengths for order discount
const (
	ORDER_DISCOUNT_NAME_MAX_LENGTH = 255
)

// order discount type's values
const (
	VOUCHER = "voucher"
	MANUAL  = "manual"
)

var OrderDiscountTypeStrings = map[string]string{
	VOUCHER: "Voucher",
	MANUAL:  "Manual",
}

var OrderDiscountValueTypeStrings = map[string]string{
	FIXED:      "Fixed",
	PERCENTAGE: "%",
}

type OrderDiscount struct {
	Id             string           `json:"id"`
	OrderID        *string          `json:"order_id"`
	Type           string           `json:"type"`
	ValueType      string           `json:"value_type"`
	Value          *decimal.Decimal `json:"value"`
	AmountValue    *decimal.Decimal `json:"amount_value"`
	Amount         *model.Money     `json:"amount,omitempty" db:"-"`
	Currency       string           `json:"currency"`
	Name           *string          `json:"name"`
	TranslatedName *string          `json:"translated_name"`
	Reason         *string          `json:"reason"`
}

func (o *OrderDiscount) IsValid() *model.AppError {
	outer := model.CreateAppErrorForModel(
		"model.order_dicount.is_valid.%s.app_error",
		"order_discount_id=",
		"OrderDiscount.IsValid",
	)

	if !model.IsValidId(o.Id) {
		return outer("id", nil)
	}
	if o.OrderID != nil && !model.IsValidId(*o.OrderID) {
		return outer("order_id", &o.Id)
	}
	if OrderDiscountTypeStrings[strings.ToLower(o.Type)] == "" {
		return outer("type", &o.Id)
	}
	if OrderDiscountValueTypeStrings[strings.ToLower(o.ValueType)] == "" {
		return outer("value_type", &o.Id)
	}
	if o.Name != nil && utf8.RuneCountInString(*o.Name) > ORDER_DISCOUNT_NAME_MAX_LENGTH {
		return outer("name", &o.Id)
	}
	if o.TranslatedName != nil && utf8.RuneCountInString(*o.TranslatedName) > ORDER_DISCOUNT_NAME_MAX_LENGTH {
		return outer("translated_name", &o.Id)
	}
	if unit, err := currency.ParseISO(o.Currency); err != nil || !strings.EqualFold(unit.String(), o.Currency) {
		return outer("currency", &o.Id)
	}

	return nil
}

func (o *OrderDiscount) PreSave() {
	if o.Id == "" {
		o.Id = model.NewId()
	}
	if o.Type == "" {
		o.Type = MANUAL
	}
	if o.ValueType == "" {
		o.ValueType = FIXED
	}
}

func (o *OrderDiscount) ToJson() string {
	o.Amount = &model.Money{
		Amount:   o.AmountValue,
		Currency: o.Currency,
	}
	return model.ModelToJson(o)
}

func OrderDiscountFromJson(data io.Reader) *OrderDiscount {
	var o OrderDiscount
	model.ModelFromJson(&o, data)
	return &o
}
