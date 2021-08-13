package order

import (
	"fmt"
	"strings"
	"unicode/utf8"

	"github.com/site-name/decimal"
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
	Id                                string               `json:"id"`
	OrderID                           string               `json:"order_id"`
	VariantID                         *string              `json:"variant_id"`
	ProductName                       string               `json:"product_name"`
	VariantName                       string               `json:"variant_name"`
	TranslatedProductName             string               `json:"translated_product_name"`
	TranslatedVariantName             string               `json:"translated_variant_name"`
	ProductSku                        string               `json:"product_sku"`
	IsShippingRequired                bool                 `json:"is_shipping_required"`
	Quantity                          uint                 `json:"quantity"`
	QuantityFulfilled                 uint                 `json:"quantity_fulfilled"`
	Currency                          string               `json:"currency"`
	UnitDiscountAmount                *decimal.Decimal     `json:"unit_discount_amount"`
	UnitDiscount                      *goprices.Money      `json:"unit_dsicount" db:"-"`
	UnitDiscountType                  string               `json:"unit_discount_type"`
	UnitDiscountReason                *string              `json:"unit_discount_reason"`
	UnitPriceNetAmount                *decimal.Decimal     `json:"unit_price_net_amount"`
	UnitDiscountValue                 *decimal.Decimal     `json:"unit_discount_value"` // store the value of the applied discount. Like 20%, default 0
	UnitPriceNet                      *goprices.Money      `json:"unit_price_net" db:"-"`
	UnitPriceGrossAmount              *decimal.Decimal     `json:"unit_price_gross_amount"`
	UnitPriceGross                    *goprices.Money      `json:"unit_price_gross" db:"-"`
	UnitPrice                         *goprices.TaxedMoney `json:"unit_price" db:"-"`
	TotalPriceNetAmount               *decimal.Decimal     `json:"total_price_net_amount"`
	TotalPriceNet                     *goprices.Money      `json:"total_price_net" db:"-"`
	TotalPriceGrossAmount             *decimal.Decimal     `json:"total_price_gross_amount"`
	TotalPriceGross                   *goprices.Money      `json:"total_price_gross" db:"-"`
	TotalPrice                        *goprices.TaxedMoney `json:"total_price" db:"-"`
	UnDiscountedUnitPriceGrossAmount  *decimal.Decimal     `json:"undiscounted_unit_price_gross_amount"`
	UnDiscountedUnitPriceNetAmount    *decimal.Decimal     `json:"undiscounted_unit_price_net_amount"`
	UnDiscountedUnitPrice             *goprices.TaxedMoney `json:"undiscounted_unit_price" db:"-"`
	UnDsicountedTotalPriceGrossAmount *decimal.Decimal     `json:"undiscounted_total_price_gross_amount"`
	UnDiscountedTotalPriceNetAmount   *decimal.Decimal     `json:"undiscounted_total_price_net_amount"`
	UnDiscountedTotalPrice            *goprices.TaxedMoney `json:"undiscounted_total_price" db:"-"`
	TaxRate                           *decimal.Decimal     `json:"tax_rate"` // decimal places: 4
}

func (o *OrderLine) IsValid() *model.AppError {
	outer := model.CreateAppErrorForModel(
		"model.order_line.is_valid.%s.app_error",
		"order_line_id=",
		"OrderLine.IsValid",
	)
	if !model.IsValidId(o.Id) {
		return outer("id", nil)
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
	if unit, err := currency.ParseISO(o.Currency); err != nil || !strings.EqualFold(unit.String(), o.Currency) {
		return outer("currency", &o.Id)
	}

	return nil
}

func (o *OrderLine) String() string {
	if o.VariantName != "" {
		return fmt.Sprintf("%s (%s)", o.ProductName, o.VariantName)
	}
	return o.ProductName
}

func (o *OrderLine) PopulateNonDbFields() {
	o.UnitDiscount, _ = goprices.NewMoney(o.UnitDiscountAmount, o.Currency)
	o.UnitPriceNet, _ = goprices.NewMoney(o.UnitPriceNetAmount, o.Currency)
	o.UnitPriceGross, _ = goprices.NewMoney(o.UnitPriceGrossAmount, o.Currency)
	o.TotalPriceNet, _ = goprices.NewMoney(o.TotalPriceNetAmount, o.Currency)
	o.TotalPriceGross, _ = goprices.NewMoney(o.TotalPriceGrossAmount, o.Currency)

	o.UnitPrice, _ = goprices.NewTaxedMoney(o.UnitPriceNet, o.UnitPriceGross)
	o.TotalPrice, _ = goprices.NewTaxedMoney(o.TotalPriceNet, o.TotalPriceGross)

	net, _ := goprices.NewMoney(o.UnDiscountedUnitPriceNetAmount, o.Currency)
	gross, _ := goprices.NewMoney(o.UnDiscountedUnitPriceGrossAmount, o.Currency)
	o.UnDiscountedUnitPrice, _ = goprices.NewTaxedMoney(net, gross)

	net, _ = goprices.NewMoney(o.UnDiscountedTotalPriceNetAmount, o.Currency)
	gross, _ = goprices.NewMoney(o.UnDsicountedTotalPriceGrossAmount, o.Currency)
	o.UnDiscountedTotalPrice, _ = goprices.NewTaxedMoney(net, gross)
}

func (o *OrderLine) ToJson() string {
	o.PopulateNonDbFields()

	return model.ModelToJson(o)
}

func (o *OrderLine) PreSave() {
	if o.Id == "" {
		o.Id = model.NewId()
	}
	o.commonPre()
}

func (o *OrderLine) commonPre() {
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

	if o.UnitDiscount != nil {
		o.UnitDiscountAmount = o.UnitDiscount.Amount
	} else {
		o.UnitDiscountAmount = &decimal.Zero
	}

	if o.UnDiscountedUnitPrice != nil {
		o.UnDiscountedUnitPriceNetAmount = o.UnDiscountedUnitPrice.Net.Amount
		o.UnDiscountedUnitPriceGrossAmount = o.UnDiscountedUnitPrice.Gross.Amount
	} else {
		o.UnDiscountedUnitPriceNetAmount = &decimal.Zero
		o.UnDiscountedUnitPriceGrossAmount = &decimal.Zero
	}

	if o.UnDiscountedTotalPrice != nil {
		o.UnDiscountedTotalPriceNetAmount = o.UnDiscountedTotalPrice.Net.Amount
		o.UnDsicountedTotalPriceGrossAmount = o.UnDiscountedTotalPrice.Gross.Amount
	} else {
		o.UnDiscountedTotalPriceNetAmount = &decimal.Zero
		o.UnDsicountedTotalPriceGrossAmount = &decimal.Zero
	}

	if o.TaxRate == nil {
		o.TaxRate = &decimal.Zero
	}
}

func (o *OrderLine) PreUpdate() {
	o.commonPre()
}

func (o *OrderLine) QuantityUnFulfilled() uint {
	return o.Quantity - o.QuantityFulfilled
}
