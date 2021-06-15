package gqlmodel

import (
	"strings"

	goprices "github.com/site-name/go-prices"
	"github.com/sitename/sitename/model/order"
)

type OrderLine struct {
	ID                    string                 `json:"id"`
	ProductName           string                 `json:"productName"`
	VariantName           string                 `json:"variantName"`
	ProductSku            string                 `json:"productSku"`
	IsShippingRequired    bool                   `json:"isShippingRequired"`
	Quantity              int                    `json:"quantity"`
	QuantityFulfilled     int                    `json:"quantityFulfilled"`
	UnitDiscountReason    *string                `json:"unitDiscountReason"`
	TaxRate               float64                `json:"taxRate"`
	DigitalContentURLID   *string                `json:"digitalContentUrl"` // *DigitalContentURL
	Thumbnail             *Image                 `json:"thumbnail"`
	UnitPrice             *TaxedMoney            `json:"unitPrice"`
	UndiscountedUnitPrice *TaxedMoney            `json:"undiscountedUnitPrice"`
	UnitDiscount          *Money                 `json:"unitDiscount"`
	UnitDiscountValue     string                 `json:"unitDiscountValue"`
	TotalPrice            *TaxedMoney            `json:"totalPrice"`
	VariantID             *string                `json:"variant"` // *ProductVariant
	TranslatedProductName string                 `json:"translatedProductName"`
	TranslatedVariantName string                 `json:"translatedVariantName"`
	AllocationIDs         []string               `json:"allocations"` // []*Allocation
	UnitDiscountType      *DiscountValueTypeEnum `json:"unitDiscountType"`
}

func (OrderLine) IsNode() {}

func FromDatabaseOrderLine(o *order.OrderLine) *OrderLine {

	unitDiscountType := DiscountValueTypeEnum(strings.ToUpper(o.UnitDiscountType))

	taxRate, _ := o.TaxRate.Float64()

	return &OrderLine{
		ID:                    o.Id,
		ProductName:           o.ProductName,
		VariantName:           o.VariantName,
		ProductSku:            o.ProductSku,
		IsShippingRequired:    o.IsShippingRequired,
		Quantity:              o.Quantity,
		QuantityFulfilled:     o.QuantityFulfilled,
		UnitDiscountReason:    o.UnitDiscountReason,
		TaxRate:               taxRate,
		DigitalContentURLID:   nil,
		Thumbnail:             nil,
		UnitPrice:             nil,
		UndiscountedUnitPrice: nil,
		UnitDiscount:          NormalMoneyToGraphqlMoney(o.UnitDiscount),
		UnitDiscountValue:     o.UnitDiscountValue,
		TotalPrice:            o.TotalPrice,
		VariantID:             o.VariantID,
		TranslatedProductName: o.TranslatedProductName,
		TranslatedVariantName: o.TranslatedVariantName,
		AllocationIDs:         []string{},
		UnitDiscountType:      &unitDiscountType,
	}
}

func NormalMoneyToGraphqlMoney(m *goprices.Money) *Money {
	float64Amount, _ := m.Amount.Float64()

	return &Money{
		Currency: m.Currency,
		Amount:   float64Amount,
	}
}
