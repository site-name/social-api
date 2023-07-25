package model

import (
	"strings"

	"github.com/Masterminds/squirrel"
	"github.com/samber/lo"
	"github.com/site-name/decimal"
	goprices "github.com/site-name/go-prices"
	"golang.org/x/text/currency"
	"gorm.io/gorm"
)

// max lengths for order discount
const (
	ORDER_DISCOUNT_NAME_MAX_LENGTH       = 255
	ORDER_DISCOUNT_TYPE_MAX_LENGTH       = 10
	ORDER_DISCOUNT_VALUE_TYPE_MAX_LENGTH = 10
)

type OrderDiscountType string

func (o OrderDiscountType) IsValid() bool {
	return OrderDiscountTypeStrings[o] != ""
}

// order discount type's values
const (
	VOUCHER OrderDiscountType = "voucher"
	MANUAL  OrderDiscountType = "manual"
)

var OrderDiscountTypeStrings = map[OrderDiscountType]string{
	VOUCHER: "Voucher",
	MANUAL:  "Manual",
}

type OrderDiscount struct {
	Id             string            `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid();column:Id"`
	OrderID        *string           `json:"order_id" gorm:"type:uuid;column:OrderID"`
	Type           OrderDiscountType `json:"type" gorm:"type:varchar(10);column:Type"`
	ValueType      DiscountType      `json:"value_type" gorm:"type:varchar(10);column:ValueType"`
	Value          *decimal.Decimal  `json:"value" gorm:"default:0;column:Value"`              // default 0
	AmountValue    *decimal.Decimal  `json:"amount_value" gorm:"default:0;column:AmountValue"` // default 0
	Currency       string            `json:"currency" gorm:"type:varchar(3);column:Currency"`
	Name           *string           `json:"name" gorm:"type:varchar(255);column:Name"`
	TranslatedName *string           `json:"translated_name" gorm:"type:varchar(255);column:TranslatedName"`
	Reason         *string           `json:"reason" gorm:"column:Reason"`

	Amount *goprices.Money `json:"amount,omitempty" gorm:"-"`
}

func (c *OrderDiscount) BeforeCreate(_ *gorm.DB) error { c.commonPre(); return c.IsValid() }
func (c *OrderDiscount) BeforeUpdate(_ *gorm.DB) error { c.commonPre(); return c.IsValid() }
func (c *OrderDiscount) TableName() string             { return OrderDiscountTableName }

type OrderDiscountFilterOption struct {
	Conditions squirrel.Sqlizer
}

type OrderDiscounts []*OrderDiscount

func (o OrderDiscounts) IDs() []string {
	return lo.Map(o, func(item *OrderDiscount, _ int) string { return item.Id })
}

func (o *OrderDiscount) DeepCopy() *OrderDiscount {
	res := *o

	res.OrderID = CopyPointer(o.OrderID)
	res.Name = CopyPointer(o.Name)
	res.TranslatedName = CopyPointer(o.TranslatedName)
	res.Reason = CopyPointer(o.Reason)
	res.Value = CopyPointer(o.Value)
	res.AmountValue = CopyPointer(o.AmountValue)

	return &res
}

func (o *OrderDiscount) IsValid() *AppError {
	outer := CreateAppErrorForModel(
		"model.order_dicount.is_valid.%s.app_error",
		"order_discount_id=",
		"OrderDiscount.IsValid",
	)

	if o.OrderID != nil && !IsValidId(*o.OrderID) {
		return outer("order_id", &o.Id)
	}
	if !o.Type.IsValid() {
		return outer("type", &o.Id)
	}
	if !o.ValueType.IsValid() {
		return outer("value_type", &o.Id)
	}
	if unit, err := currency.ParseISO(o.Currency); err != nil || !strings.EqualFold(unit.String(), o.Currency) {
		return outer("currency", &o.Id)
	}
	if err := ValidateDecimal("OrderDiscount.IsValid.Value", o.Value, DECIMAL_TOTAL_DIGITS_ALLOWED, DECIMAL_MAX_DECIMAL_PLACES_ALLOWED); err != nil {
		return err
	}
	if err := ValidateDecimal("OrderDiscount.IsValid.AmountValue", o.AmountValue, DECIMAL_TOTAL_DIGITS_ALLOWED, DECIMAL_MAX_DECIMAL_PLACES_ALLOWED); err != nil {
		return err
	}

	return nil
}

func (o *OrderDiscount) PopulateNonDbFields() {
	o.Amount = &goprices.Money{
		Amount:   *o.AmountValue,
		Currency: o.Currency,
	}
}

func (o *OrderDiscount) commonPre() {
	if !o.Type.IsValid() {
		o.Type = MANUAL
	}
	if o.ValueType == "" {
		o.ValueType = FIXED
	}
	if o.Value == nil {
		o.Value = &decimal.Zero
	}
	if o.AmountValue == nil {
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
	if o.Amount != nil {
		o.AmountValue = &o.Amount.Amount
	}
}
