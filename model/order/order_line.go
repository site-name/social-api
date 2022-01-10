package order

import (
	"fmt"
	"strings"
	"unicode/utf8"

	"github.com/Masterminds/squirrel"
	"github.com/site-name/decimal"
	goprices "github.com/site-name/go-prices"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model/product_and_discount"
	"golang.org/x/text/currency"
)

// max lengths for some fields of OrderLine
const (
	ORDER_LINE_PRODUCT_NAME_MAX_LENGTH       = 386
	ORDER_LINE_VARIANT_NAME_MAX_LENGTH       = 255
	ORDER_LINE_PRODUCT_SKU_MAX_LENGTH        = 255
	ORDER_LINE_PRODUCT_VARIANT_ID_MAX_LENGTH = 255
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
	CreateAt                          int64                `json:"create_at"` // for database ordering
	OrderID                           string               `json:"order_id"`
	VariantID                         *string              `json:"variant_id"` // FOREIGN KEY ProductVariant
	ProductName                       string               `json:"product_name"`
	VariantName                       string               `json:"variant_name"`
	TranslatedProductName             string               `json:"translated_product_name"`
	TranslatedVariantName             string               `json:"translated_variant_name"`
	ProductSku                        *string              `json:"product_sku"`
	ProductVariantID                  *string              `json:"product_variant_id"` // GraphQL ID used as fallback when product SKU is not available
	IsShippingRequired                bool                 `json:"is_shipping_required"`
	IsGiftcard                        bool                 `json:"is_gift_card"`
	Quantity                          int                  `json:"quantity"`
	QuantityFulfilled                 int                  `json:"quantity_fulfilled"`
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
	UnDiscountedTotalPriceGrossAmount *decimal.Decimal     `json:"undiscounted_total_price_gross_amount"`
	UnDiscountedTotalPriceNetAmount   *decimal.Decimal     `json:"undiscounted_total_price_net_amount"`
	UnDiscountedTotalPrice            *goprices.TaxedMoney `json:"undiscounted_total_price" db:"-"`
	TaxRate                           *decimal.Decimal     `json:"tax_rate"` // decimal places: 4

	ProductVariant *product_and_discount.ProductVariant `json:"-" db:"-"` // for storing value returned by prefetching
	Order          *Order                               `json:"-" db:"-"` // related data, get popularized in some calls to database
	allocations    []*ReplicateWarehouseAllocation
}

func (o *OrderLine) SetAllocations(allocations []*ReplicateWarehouseAllocation) {
	o.allocations = allocations
}

func (o *OrderLine) GetAllocations() []*ReplicateWarehouseAllocation {
	return o.allocations
}

// OrderLinePrefetchRelated
type OrderLinePrefetchRelated struct {
	VariantProduct        bool // This tells store to prefetch related ProductVariant(s) and Product(s) as well
	VariantDigitalContent bool
	VariantStocks         bool
	AllocationsStock      bool
}

// OrderLineFilterOption is used for build sql queries
type OrderLineFilterOption struct {
	Id                 squirrel.Sqlizer
	OrderID            squirrel.Sqlizer
	IsShippingRequired *bool
	IsGiftcard         *bool
	VariantID          squirrel.Sqlizer

	VariantProductID        squirrel.Sqlizer // INNER JOIN ProductVariants INNER JOIN Products WHERE Products.Id ...
	VariantDigitalContentID squirrel.Sqlizer // INNER JOIN ProductVariants INNER JOIN DigitalContents WHERE DigitalContents.Id ...

	PrefetchRelated OrderLinePrefetchRelated
}

func (o *OrderLine) String() string {
	if o.VariantName != "" {
		return fmt.Sprintf("%s (%s)", o.ProductName, o.VariantName)
	}
	return o.ProductName
}

type OrderLines []*OrderLine

// ProductVariantIDs returns only non-nil product variant ids
func (o OrderLines) ProductVariantIDs() []string {
	res := []string{}
	for _, orderLine := range o {
		if orderLine != nil && orderLine.VariantID != nil {
			res = append(res, *orderLine.VariantID)
		}
	}

	return res
}

func (o OrderLines) IDs() []string {
	res := []string{}
	for _, orderLine := range o {
		if orderLine != nil {
			res = append(res, orderLine.Id)
		}
	}

	return res
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
	if o.CreateAt == 0 {
		return outer("create_at", &o.Id)
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
	if o.ProductSku != nil && len(*o.ProductSku) > ORDER_LINE_PRODUCT_SKU_MAX_LENGTH {
		return outer("product_sku", &o.Id)
	}
	if o.ProductVariantID != nil && len(*o.ProductVariantID) > ORDER_LINE_PRODUCT_VARIANT_ID_MAX_LENGTH {
		return outer("product_variant_id", &o.Id)
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

func (o *OrderLine) PopulateNonDbFields() {
	if o.UnitDiscountAmount != nil {
		o.UnitDiscount = &goprices.Money{
			Amount:   *o.UnitDiscountAmount,
			Currency: o.Currency,
		}
	}
	if o.UnitPriceNetAmount != nil && o.UnitPriceGrossAmount != nil {
		o.UnitPriceNet = &goprices.Money{
			Amount:   *o.UnitPriceNetAmount,
			Currency: o.Currency,
		}
		o.UnitPriceGross = &goprices.Money{
			Amount:   *o.UnitPriceGrossAmount,
			Currency: o.Currency,
		}
		o.UnitPrice, _ = goprices.NewTaxedMoney(o.UnitPriceNet, o.UnitPriceGross)
	}

	if o.TotalPriceNetAmount != nil && o.TotalPriceGrossAmount != nil {
		o.TotalPriceNet = &goprices.Money{
			Amount:   *o.TotalPriceNetAmount,
			Currency: o.Currency,
		}
		o.TotalPriceGross = &goprices.Money{
			Amount:   *o.TotalPriceGrossAmount,
			Currency: o.Currency,
		}
		o.TotalPrice, _ = goprices.NewTaxedMoney(o.TotalPriceNet, o.TotalPriceGross)

	}

	if o.UnDiscountedUnitPriceNetAmount != nil && o.UnDiscountedUnitPriceGrossAmount != nil {
		net := &goprices.Money{
			Amount:   *o.UnDiscountedUnitPriceNetAmount,
			Currency: o.Currency,
		}
		gross := &goprices.Money{
			Amount:   *o.UnDiscountedUnitPriceGrossAmount,
			Currency: o.Currency,
		}
		o.UnDiscountedUnitPrice, _ = goprices.NewTaxedMoney(net, gross)
	}

	if o.UnDiscountedTotalPriceNetAmount != nil && o.UnDiscountedTotalPriceGrossAmount != nil {
		net := &goprices.Money{
			Amount:   *o.UnDiscountedTotalPriceNetAmount,
			Currency: o.Currency,
		}
		gross := &goprices.Money{
			Amount:   *o.UnDiscountedTotalPriceGrossAmount,
			Currency: o.Currency,
		}
		o.UnDiscountedTotalPrice, _ = goprices.NewTaxedMoney(net, gross)
	}
}

func (o *OrderLine) PreSave() {
	if o.Id == "" {
		o.Id = model.NewId()
	}
	o.CreateAt = model.GetMillis()

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
		o.UnitDiscountAmount = &o.UnitDiscount.Amount
	} else {
		o.UnitDiscountAmount = &decimal.Zero
	}

	if o.UnitPrice != nil {
		o.UnitPriceNet = o.UnitPrice.Net
		o.UnitPriceGross = o.UnitPrice.Gross
		o.Currency = o.UnitPrice.Currency
	}

	if o.UnDiscountedUnitPrice != nil {
		o.UnDiscountedUnitPriceNetAmount = &o.UnDiscountedUnitPrice.Net.Amount
		o.UnDiscountedUnitPriceGrossAmount = &o.UnDiscountedUnitPrice.Gross.Amount
	} else {
		o.UnDiscountedUnitPriceNetAmount = &decimal.Zero
		o.UnDiscountedUnitPriceGrossAmount = &decimal.Zero
	}

	if o.UnDiscountedTotalPrice != nil {
		o.UnDiscountedTotalPriceNetAmount = &o.UnDiscountedTotalPrice.Net.Amount
		o.UnDiscountedTotalPriceGrossAmount = &o.UnDiscountedTotalPrice.Gross.Amount
	} else {
		o.UnDiscountedTotalPriceNetAmount = &decimal.Zero
		o.UnDiscountedTotalPriceGrossAmount = &decimal.Zero
	}

	if o.TaxRate == nil {
		o.TaxRate = &decimal.Zero
	}
}

func (o *OrderLine) PreUpdate() {
	o.commonPre()
}

// QuantityUnFulfilled return current order's Quantity subtract QuantityFulfilled
func (o *OrderLine) QuantityUnFulfilled() int {
	return o.Quantity - o.QuantityFulfilled
}

func (o *OrderLine) DeepCopy() *OrderLine {
	orderLine := *o

	return &orderLine
}
