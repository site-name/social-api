package model

import (
	"fmt"
	"strings"

	"github.com/mattermost/squirrel"
	"github.com/samber/lo"
	"github.com/site-name/decimal"
	goprices "github.com/site-name/go-prices"
	"golang.org/x/text/currency"
	"gorm.io/gorm"
)

type OrderLine struct {
	Id                                string            `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid();column:Id"`
	CreateAt                          int64             `json:"create_at" gorm:"type:bigint;column:CreateAt;autoCreateTime:milli"` // for database ordering
	OrderID                           string            `json:"order_id" gorm:"type:uuid;column:OrderID"`                          // NOTE editable
	VariantID                         *string           `json:"variant_id" gorm:"type:uuid;column:VariantID"`                      // FOREIGN KEY ProductVariant
	ProductName                       string            `json:"product_name" gorm:"type:varchar(386);column:ProductName"`
	VariantName                       string            `json:"variant_name" gorm:"type:varchar(255);column:VariantName"`
	TranslatedProductName             string            `json:"translated_product_name" gorm:"type:varchar(386);column:TranslatedProductName"`
	TranslatedVariantName             string            `json:"translated_variant_name" gorm:"type:varchar(255);column:TranslatedVariantName"`
	ProductSku                        *string           `json:"product_sku" gorm:"type:varchar(255);column:ProductSku"`
	ProductVariantID                  *string           `json:"product_variant_id" gorm:"type:varchar(255);column:ProductVariantID"` // GraphQL ID used as fallback when product SKU is not available
	IsShippingRequired                bool              `json:"is_shipping_required" gorm:"column:IsShippingRequired"`
	IsGiftcard                        bool              `json:"is_gift_card" gorm:"column:IsGiftcard"`
	Quantity                          int               `json:"quantity" gorm:"type:integer;check:Quantity >= 1;column:Quantity"`
	QuantityFulfilled                 int               `json:"quantity_fulfilled" gorm:"type:integer;check:QuantityFulfilled >= 0;column:QuantityFulfilled"`
	Currency                          string            `json:"currency" gorm:"type:varchar(3);column:Currency"`
	UnitDiscountAmount                *decimal.Decimal  `json:"unit_discount_amount" gorm:"default:0;column:UnitDiscountAmount;type:decimal(12,3)"` // default 0
	UnitDiscountType                  DiscountValueType `json:"unit_discount_type" gorm:"type:varchar(10);column:UnitDiscountType"`                 // default 'fixed'
	UnitDiscountReason                *string           `json:"unit_discount_reason" gorm:"column:UnitDiscountReason"`
	UnitPriceNetAmount                *decimal.Decimal  `json:"unit_price_net_amount" gorm:"default:0;column:UnitPriceNetAmount;type:decimal(12,3)"`     // default 0
	UnitDiscountValue                 *decimal.Decimal  `json:"unit_discount_value" gorm:"default:0;column:UnitDiscountValue;type:decimal(12,3)"`        // store the value of the applied discount. Like 20%, default 0
	UnitPriceGrossAmount              *decimal.Decimal  `json:"unit_price_gross_amount" gorm:"default:0;column:UnitPriceGrossAmount;type:decimal(12,3)"` // default 0
	TotalPriceNetAmount               *decimal.Decimal  `json:"total_price_net_amount" gorm:"column:TotalPriceNetAmount;type:decimal(12,3)"`
	TotalPriceGrossAmount             *decimal.Decimal  `json:"total_price_gross_amount" gorm:"column:TotalPriceGrossAmount;type:decimal(12,3)"`
	UnDiscountedUnitPriceGrossAmount  *decimal.Decimal  `json:"undiscounted_unit_price_gross_amount" gorm:"column:UnDiscountedUnitPriceGrossAmount;default:0;type:decimal(12,3)"`
	UnDiscountedUnitPriceNetAmount    *decimal.Decimal  `json:"undiscounted_unit_price_net_amount" gorm:"column:UnDiscountedUnitPriceNetAmount;default:0;type:decimal(12,3)"`
	UnDiscountedTotalPriceGrossAmount *decimal.Decimal  `json:"undiscounted_total_price_gross_amount" gorm:"column:UnDiscountedTotalPriceGrossAmount;default:0;type:decimal(12,3)"` // default 0
	UnDiscountedTotalPriceNetAmount   *decimal.Decimal  `json:"undiscounted_total_price_net_amount" gorm:"column:UnDiscountedTotalPriceNetAmount;default:0;type:decimal(12,3)"`     // default 0
	TaxRate                           *decimal.Decimal  `json:"tax_rate" gorm:"column:TaxRate;type:decimal(5,4)"`                                                                   // decimal places: 4, default: 0

	UnitDiscount           *goprices.Money      `json:"unit_dsicount" gorm:"-"`
	UnDiscountedTotalPrice *goprices.TaxedMoney `json:"undiscounted_total_price" gorm:"-"`
	UnDiscountedUnitPrice  *goprices.TaxedMoney `json:"undiscounted_unit_price" gorm:"-"`
	TotalPrice             *goprices.TaxedMoney `json:"total_price" gorm:"-"`
	TotalPriceGross        *goprices.Money      `json:"total_price_gross" gorm:"-"`
	TotalPriceNet          *goprices.Money      `json:"total_price_net" gorm:"-"`
	UnitPrice              *goprices.TaxedMoney `json:"unit_price" gorm:"-"`
	UnitPriceGross         *goprices.Money      `json:"unit_price_gross" gorm:"-"`
	UnitPriceNet           *goprices.Money      `json:"unit_price_net" gorm:"-"`

	ProductVariant *ProductVariant `json:"-" gorm:"constraint:OnDelete:SET NULL"` // for storing value returned by prefetching
	Order          *Order          `json:"-" gorm:"constraint:OnDelete:CASCADE"`  // related data, get popularized in some calls to database
	Allocations    Allocations     `json:"-" gorm:"foreignKey:OrderLineID"`
}

// OrderLine table column names
const (
	OrderLineColumnId                                = "Id"
	OrderLineColumnCreateAt                          = "CreateAt"
	OrderLineColumnOrderID                           = "OrderID"
	OrderLineColumnVariantID                         = "VariantID"
	OrderLineColumnProductName                       = "ProductName"
	OrderLineColumnVariantName                       = "VariantName"
	OrderLineColumnTranslatedProductName             = "TranslatedProductName"
	OrderLineColumnTranslatedVariantName             = "TranslatedVariantName"
	OrderLineColumnProductSku                        = "ProductSku"
	OrderLineColumnProductVariantID                  = "ProductVariantID"
	OrderLineColumnIsShippingRequired                = "IsShippingRequired"
	OrderLineColumnIsGiftcard                        = "IsGiftcard"
	OrderLineColumnQuantity                          = "Quantity"
	OrderLineColumnQuantityFulfilled                 = "QuantityFulfilled"
	OrderLineColumnCurrency                          = "Currency"
	OrderLineColumnUnitDiscountAmount                = "UnitDiscountAmount"
	OrderLineColumnUnitDiscountType                  = "UnitDiscountType"
	OrderLineColumnUnitDiscountReason                = "UnitDiscountReason"
	OrderLineColumnUnitPriceNetAmount                = "UnitPriceNetAmount"
	OrderLineColumnUnitDiscountValue                 = "UnitDiscountValue"
	OrderLineColumnUnitPriceGrossAmount              = "UnitPriceGrossAmount"
	OrderLineColumnTotalPriceNetAmount               = "TotalPriceNetAmount"
	OrderLineColumnTotalPriceGrossAmount             = "TotalPriceGrossAmount"
	OrderLineColumnUnDiscountedUnitPriceGrossAmount  = "UnDiscountedUnitPriceGrossAmount"
	OrderLineColumnUnDiscountedUnitPriceNetAmount    = "UnDiscountedUnitPriceNetAmount"
	OrderLineColumnUnDiscountedTotalPriceGrossAmount = "UnDiscountedTotalPriceGrossAmount"
	OrderLineColumnUnDiscountedTotalPriceNetAmount   = "UnDiscountedTotalPriceNetAmount"
	OrderLineColumnTaxRate                           = "TaxRate"
)

func (c *OrderLine) BeforeCreate(_ *gorm.DB) error { c.commonPre(); return c.IsValid() }
func (c *OrderLine) BeforeUpdate(_ *gorm.DB) error { c.commonPre(); return c.IsValid() }
func (c *OrderLine) TableName() string             { return OrderLineTableName }

// OrderLineFilterOption is used for build sql queries
type OrderLineFilterOption struct {
	Conditions squirrel.Sqlizer

	RelatedOrderConditions squirrel.Sqlizer // INNER JOIN Orders ON ... WHERE Orders....

	VariantProductID squirrel.Sqlizer // INNER JOIN ProductVariants ON ... WHERE ProductVariants.ProductID ...

	// INNER JOIN ProductVariants ON OrderLines.VariantID = ProductVariants.Id
	// INNER JOIN DigitalContents ON ProductVariants.Id = DigitalContents.ProductVariantID WHERE DigitalContents.Id ...
	VariantDigitalContentID squirrel.Sqlizer

	// Thanks to Gorm's Preload feature, we can select related values easily
	//
	// E.g
	//  "ProductVariant" // will fetch related product variant(s)
	//  "ProductVariant.Product" // will fetch related variants, product
	//  "Order"
	Preload []string
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

func (o OrderLines) OrderIDs() []string {
	return lo.Map(o, func(l *OrderLine, _ int) string { return l.OrderID })
}

func (o OrderLines) IDs() []string {
	return lo.Map(o, func(l *OrderLine, _ int) string { return l.Id })
}

func (o OrderLines) FilterNils() OrderLines {
	return lo.Filter(o, func(l *OrderLine, _ int) bool {
		return l != nil
	})
}

func (o *OrderLine) IsValid() *AppError {
	outer := CreateAppErrorForModel(
		"model.order_line.is_valid.%s.app_error",
		"order_line_id=",
		"OrderLine.IsValid",
	)

	if !IsValidId(o.OrderID) {
		return outer("order_id", &o.Id)
	}
	if o.VariantID != nil && !IsValidId(*o.VariantID) {
		return outer("variant_id", &o.Id)
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

func (o *OrderLine) commonPre() {
	o.ProductName = SanitizeUnicode(o.ProductName)
	o.VariantName = SanitizeUnicode(o.VariantName)
	o.TranslatedProductName = SanitizeUnicode(o.TranslatedProductName)
	o.TranslatedVariantName = SanitizeUnicode(o.TranslatedVariantName)

	if o.UnitDiscountReason != nil {
		o.UnitDiscountReason = GetPointerOfValue(SanitizeUnicode(*o.UnitDiscountReason))
	}
	if !o.UnitDiscountType.IsValid() {
		o.UnitDiscountType = DISCOUNT_VALUE_TYPE_FIXED
	}
	if o.UnitDiscountValue == nil {
		o.UnitDiscountValue = GetPointerOfValue(decimal.Zero)
	}

	if o.UnitDiscount != nil {
		o.UnitDiscountAmount = &o.UnitDiscount.Amount
	} else {
		o.UnitDiscountAmount = GetPointerOfValue(decimal.Zero)
	}

	if o.UnitPrice != nil {
		o.UnitPriceNetAmount = &o.UnitPrice.Net.Amount
		o.UnitPriceGrossAmount = &o.UnitPrice.Gross.Amount
	} else {
		o.UnitPriceNetAmount = GetPointerOfValue(decimal.Zero)
		o.UnitPriceGrossAmount = GetPointerOfValue(decimal.Zero)
	}

	if o.TotalPrice != nil {
		o.TotalPriceNetAmount = &o.TotalPrice.Net.Amount
		o.TotalPriceGrossAmount = &o.TotalPrice.Gross.Amount
	}

	if o.UnDiscountedUnitPrice != nil {
		o.UnDiscountedUnitPriceNetAmount = &o.UnDiscountedUnitPrice.Net.Amount
		o.UnDiscountedUnitPriceGrossAmount = &o.UnDiscountedUnitPrice.Gross.Amount
	} else {
		o.UnDiscountedUnitPriceNetAmount = GetPointerOfValue(decimal.Zero)
		o.UnDiscountedUnitPriceGrossAmount = GetPointerOfValue(decimal.Zero)
	}

	if o.UnDiscountedTotalPrice != nil {
		o.UnDiscountedTotalPriceNetAmount = &o.UnDiscountedTotalPrice.Net.Amount
		o.UnDiscountedTotalPriceGrossAmount = &o.UnDiscountedTotalPrice.Gross.Amount
	} else {
		o.UnDiscountedTotalPriceNetAmount = GetPointerOfValue(decimal.Zero)
		o.UnDiscountedTotalPriceGrossAmount = GetPointerOfValue(decimal.Zero)
	}

	if o.TaxRate == nil {
		o.TaxRate = GetPointerOfValue(decimal.Zero)
	}
}

// QuantityUnFulfilled return current order's Quantity subtract QuantityFulfilled
func (o *OrderLine) QuantityUnFulfilled() int {
	return o.Quantity - o.QuantityFulfilled
}

func (o *OrderLine) DeepCopy() *OrderLine {
	orderLine := *o

	if o.VariantID != nil {
		orderLine.VariantID = GetPointerOfValue(*o.VariantID)
	}
	if o.ProductSku != nil {
		orderLine.ProductSku = GetPointerOfValue(*o.ProductSku)
	}
	if o.ProductVariantID != nil {
		orderLine.ProductVariantID = GetPointerOfValue(*o.ProductVariantID)
	}
	if o.UnitDiscountReason != nil {
		orderLine.UnitDiscountReason = GetPointerOfValue(*o.UnitDiscountReason)
	}

	if o.UnitDiscountAmount != nil {
		orderLine.UnitDiscountAmount = GetPointerOfValue(*o.UnitDiscountAmount)
	}
	if o.UnitPriceNetAmount != nil {
		orderLine.UnitPriceNetAmount = GetPointerOfValue(*o.UnitPriceNetAmount)
	}
	if o.UnitDiscountValue != nil {
		orderLine.UnitDiscountValue = GetPointerOfValue(*o.UnitDiscountValue)
	}
	if o.UnitPriceGrossAmount != nil {
		orderLine.UnitPriceGrossAmount = GetPointerOfValue(*o.UnitPriceGrossAmount)
	}
	if o.TotalPriceNetAmount != nil {
		orderLine.TotalPriceNetAmount = GetPointerOfValue(*o.TotalPriceNetAmount)
	}
	if o.TotalPriceGrossAmount != nil {
		orderLine.TotalPriceGrossAmount = GetPointerOfValue(*o.TotalPriceGrossAmount)
	}
	if o.UnDiscountedUnitPriceNetAmount != nil {
		orderLine.UnDiscountedUnitPriceNetAmount = GetPointerOfValue(*o.UnDiscountedUnitPriceNetAmount)
	}
	if o.UnDiscountedUnitPriceGrossAmount != nil {
		orderLine.UnDiscountedUnitPriceGrossAmount = GetPointerOfValue(*o.UnDiscountedUnitPriceGrossAmount)
	}
	if o.UnDiscountedTotalPriceGrossAmount != nil {
		orderLine.UnDiscountedTotalPriceGrossAmount = GetPointerOfValue(*o.UnDiscountedTotalPriceGrossAmount)
	}
	if o.UnDiscountedTotalPriceNetAmount != nil {
		orderLine.UnDiscountedTotalPriceNetAmount = GetPointerOfValue(*o.UnDiscountedTotalPriceNetAmount)
	}
	if o.TaxRate != nil {
		orderLine.TaxRate = GetPointerOfValue(*o.TaxRate)
	}
	if o.ProductVariant != nil {
		orderLine.ProductVariant = o.ProductVariant.DeepCopy()
	}
	if o.Order != nil {
		orderLine.Order = o.Order.DeepCopy()
	}
	orderLine.Allocations = o.Allocations.DeepCopy()

	return &orderLine
}
