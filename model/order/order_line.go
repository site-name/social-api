package order

import (
	"fmt"
	"io"
	"strings"
	"unicode/utf8"

	"github.com/shopspring/decimal"
	goprices "github.com/site-name/go-prices"
	"github.com/sitename/sitename/model"
	"golang.org/x/text/currency"
)

// max lengths for some fields of OrderLine
const (
	ORDER_LINE_PRODUCT_NAME_MAX_LENGTH       = 386
	ORDER_LINE_VARIANT_NAME_MAX_LENGTH       = 255
	ORDER_LINE_PRODUCT_SKU_MAX_LENGTH        = 255
	ORDER_LINE_UNIT_DISCOUNT_TYPE_MAX_LENGTH = 10
)

// valid values for order line's unit discount type
const (
	FIXED      = "fixed"
	PERCENTAGE = "percentage"
)

var UnitDiscountTypeStrings = map[string]string{
	FIXED:      "fixed",
	PERCENTAGE: "%",
}

type OrderLine struct {
	Id                    string               `json:"id"`
	OrderID               string               `json:"order_id"`
	VariantID             *string              `json:"variant_id"`
	ProductName           string               `json:"product_name"`
	VariantName           string               `json:"variant_name"`
	TranslatedProductName string               `json:"translated_product_name"`
	TranslatedVariantName string               `json:"translated_variant_name"`
	ProductSku            string               `json:"product_sku"`
	IsShippingRequired    bool                 `json:"is_shipping_required"`
	Quantity              int                  `json:"quantity"`
	QuantityFulfilled     int                  `json:"quantity_fulfilled"`
	Currency              string               `json:"currency"`
	UnitDiscountAmount    *decimal.Decimal     `json:"unit_discount_amount"`
	UnitDiscount          *goprices.Money      `json:"unit_dsicount" db:"-"`
	UnitDiscountType      string               `json:"unit_discount_type"`
	UnitDiscountReason    *string              `json:"unit_discount_reason"`
	UnitPriceNetAmount    *decimal.Decimal     `json:"unit_price_net_amount"`
	UnitDiscountValue     *decimal.Decimal     `json:"unit_discount_value"`
	UnitPriceNet          *goprices.Money      `json:"unit_price_net" db:"-"`
	UnitPriceGrossAmount  *decimal.Decimal     `json:"unit_price_gross_amount"`
	UnitPriceGross        *goprices.Money      `json:"unit_price_gross" db:"-"`
	UnitPrice             *goprices.TaxedMoney `json:"unit_price" db:"-"`
	TotalPriceNetAmount   *decimal.Decimal     `json:"total_price_net_amount"`
	TotalPriceNet         *goprices.Money      `json:"total_price_net" db:"-"`
	TotalPriceGrossAmount *decimal.Decimal     `json:"total_price_gross_amount"`
	TotalPriceGross       *goprices.Money      `json:"total_price_gross" db:"-"`
	TotalPrice            *goprices.TaxedMoney `json:"total_price" db:"-"`
	TaxRate               *decimal.Decimal     `json:"tax_rate"`
}

func (o *OrderLine) IsValid() *model.AppError {
	outer := model.CreateAppErrorForModel(
		"model.order_line.is_valid.%s.app_error",
		"order_line_id=",
		"OrderLine.IsValid",
	)
	if !model.IsValidId(o.Id) {
		return outer("Id", nil)
	}
	if !model.IsValidId(o.OrderID) {
		return outer("order_id", &o.Id)
	}
	if o.VariantID != nil && !model.IsValidId(*o.VariantID) {
		return outer("variant_id", &o.Id)
	}
	if utf8.RuneCountInString(o.ProductName) > ORDER_LINE_PRODUCT_NAME_MAX_LENGTH {
		return outer("product_name", &o.Id)
	}
	if utf8.RuneCountInString(o.VariantName) > ORDER_LINE_VARIANT_NAME_MAX_LENGTH {
		return outer("variant_name", &o.Id)
	}
	if utf8.RuneCountInString(o.TranslatedProductName) > ORDER_LINE_PRODUCT_NAME_MAX_LENGTH {
		return outer("translated_product_name", &o.Id)
	}
	if utf8.RuneCountInString(o.TranslatedVariantName) > ORDER_LINE_VARIANT_NAME_MAX_LENGTH {
		return outer("translated_variant_name", &o.Id)
	}
	if len(o.ProductSku) > ORDER_LINE_PRODUCT_SKU_MAX_LENGTH {
		return outer("product_sku", &o.Id)
	}
	if len(o.UnitDiscountType) > ORDER_LINE_UNIT_DISCOUNT_TYPE_MAX_LENGTH {
		return outer("unit_discount_type", &o.Id)
	}
	if o.Quantity < 1 {
		return outer("quantity", &o.Id)
	}
	if o.QuantityFulfilled < 0 {
		return outer("quantity_fulfilled", &o.Id)
	}
	if unit, err := currency.ParseISO(o.Currency); err != nil || !strings.EqualFold(unit.String(), o.Currency) {
		return outer("currency", &o.Id)
	}

	return nil
}

func (o *OrderLine) ToJson() string {
	if o.UnitDiscount == nil {
		o.UnitDiscount = &goprices.Money{
			Amount:   o.UnitDiscountAmount,
			Currency: o.Currency,
		}
	}
	if o.UnitPriceNet == nil {
		o.UnitDiscount = &goprices.Money{
			Amount:   o.UnitPriceNetAmount,
			Currency: o.Currency,
		}
	}
	if o.UnitPriceGross == nil {
		o.UnitDiscount = &goprices.Money{
			Amount:   o.UnitPriceGrossAmount,
			Currency: o.Currency,
		}
	}
	if o.TotalPriceNet == nil {
		o.UnitDiscount = &goprices.Money{
			Amount:   o.TotalPriceNetAmount,
			Currency: o.Currency,
		}
	}
	if o.TotalPriceGross == nil {
		o.UnitDiscount = &goprices.Money{
			Amount:   o.TotalPriceGrossAmount,
			Currency: o.Currency,
		}
	}
	if o.UnitPrice == nil {
		o.UnitPrice = &goprices.TaxedMoney{
			Net:      o.UnitPriceNet,
			Gross:    o.UnitPriceGross,
			Currency: o.Currency,
		}
	}
	if o.TotalPrice == nil {
		o.UnitPrice = &goprices.TaxedMoney{
			Net:      o.TotalPriceNet,
			Gross:    o.TotalPriceGross,
			Currency: o.Currency,
		}
	}

	return model.ModelToJson(o)
}

func OrderLineFromJson(data io.Reader) *OrderLine {
	var o OrderLine
	model.ModelFromJson(&o, data)
	return &o
}

func (o *OrderLine) PreSave() {
	if o.Id == "" {
		o.Id = model.NewId()
	}
	o.commonPrePostActions()
}

func (o *OrderLine) commonPrePostActions() {
	o.ProductName = model.SanitizeUnicode(o.ProductName)
	o.VariantName = model.SanitizeUnicode(o.VariantName)
	o.TranslatedProductName = model.SanitizeUnicode(o.TranslatedProductName)
	o.TranslatedVariantName = model.SanitizeUnicode(o.TranslatedVariantName)

	if o.UnitDiscountReason != nil {
		o.UnitDiscountReason = model.NewString(model.SanitizeUnicode(*o.UnitDiscountReason))
	}
	if o.UnitDiscountType == "" {
		o.UnitDiscountType = FIXED
	}
	if o.UnitDiscountValue == nil {
		o.UnitDiscountValue = &decimal.Zero
	}
	if o.UnitDiscountAmount == nil {
		o.UnitDiscountAmount = &decimal.Zero
	}
	if o.TaxRate == nil {
		o.TaxRate = &decimal.Zero
	}
}

func (o *OrderLine) PreUpdate() {
	o.commonPrePostActions()
}

func (o *OrderLine) String() string {
	return fmt.Sprintf("%s (%s)", o.ProductName, o.VariantName)
}

func (o *OrderLine) UndiscountedUnitPrice() (*goprices.TaxedMoney, error) {
	return o.UnitPrice.Add(o.UnitDiscount)
}

func (o *OrderLine) QuantityUnFulfilled() int {
	return o.Quantity - o.QuantityFulfilled
}
