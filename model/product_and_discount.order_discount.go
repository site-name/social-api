package model

import (
	"strings"
	"unicode/utf8"

	"github.com/Masterminds/squirrel"
	"github.com/site-name/decimal"
	goprices "github.com/site-name/go-prices"
	"golang.org/x/text/currency"
)

// max lengths for order discount
const (
	ORDER_DISCOUNT_NAME_MAX_LENGTH       = 255
	ORDER_DISCOUNT_TYPE_MAX_LENGTH       = 10
	ORDER_DISCOUNT_VALUE_TYPE_MAX_LENGTH = 10
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
	Value          *decimal.Decimal `json:"value"`        // default 0
	AmountValue    *decimal.Decimal `json:"amount_value"` // default 0
	Amount         *goprices.Money  `json:"amount,omitempty" db:"-"`
	Currency       string           `json:"currency"`
	Name           *string          `json:"name"`
	TranslatedName *string          `json:"translated_name"`
	Reason         *string          `json:"reason"`
}

// OrderDiscountFilterOption is used to build sql queries
type OrderDiscountFilterOption struct {
	Id      squirrel.Sqlizer
	OrderID squirrel.Sqlizer
	Type    squirrel.Sqlizer
}

type OrderDiscounts []*OrderDiscount

func (o OrderDiscounts) IDs() []string {
	res := make([]string, len(o))
	for i := range o {
		res[i] = o[i].Id
	}

	return res
}

func (o *OrderDiscount) DeepCopy() *OrderDiscount {
	newOrderDiscount := *o

	return &newOrderDiscount
}

func (o *OrderDiscount) IsValid() *AppError {
	outer := CreateAppErrorForModel(
		"order_dicount.is_valid.%s.app_error",
		"order_discount_id=",
		"OrderDiscount.IsValid",
	)

	if !IsValidId(o.Id) {
		return outer("id", nil)
	}
	if o.OrderID != nil && !IsValidId(*o.OrderID) {
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

func (o *OrderDiscount) PopulateNonDbFields() {
	o.Amount = &goprices.Money{
		Amount:   *o.AmountValue,
		Currency: o.Currency,
	}
}

func (o *OrderDiscount) PreSave() {
	if o.Id == "" {
		o.Id = NewId()
	}
	o.commonPre()
}

func (o *OrderDiscount) commonPre() {
	if o.Type == "" {
		o.Type = MANUAL
	}
	if o.ValueType == "" {
		o.ValueType = FIXED
	}
	if o.Value == nil {
		o.Value = &decimal.Zero
	}

	if o.Amount != nil {
		o.AmountValue = &o.Amount.Amount
	} else {
		o.AmountValue = &decimal.Zero
	}

	if o.Name != nil {
		*o.Name = SanitizeUnicode(*o.Name)
	}
	if o.TranslatedName != nil {
		*o.TranslatedName = SanitizeUnicode(*o.TranslatedName)
	}
	if o.Reason != nil {
		*o.Reason = SanitizeUnicode(*o.Reason)
	}
	if o.Currency != "" {
		o.Currency = strings.ToUpper(o.Currency)
	} else {
		o.Currency = DEFAULT_CURRENCY
	}
}

func (o *OrderDiscount) PreUpdate() {
	o.commonPre()
}

func (o *OrderDiscount) ToJSON() string {
	o.PopulateNonDbFields()
	return ModelToJson(o)
}