package gqlmodel

import (
	"strings"

	"time"

	"github.com/shopspring/decimal"
	goprices "github.com/site-name/go-prices"
	"github.com/sitename/sitename/model/order"
)

// OrderLine represents data format will be returned to end user
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
	Thumbnail             func(size *int) *Image `json:"thumbnail"`
	UnitPrice             *TaxedMoney            `json:"unitPrice"`
	UndiscountedUnitPrice *TaxedMoney            `json:"undiscountedUnitPrice"`
	UnitDiscount          *Money                 `json:"unitDiscount"`
	UnitDiscountValue     *decimal.Decimal       `json:"unitDiscountValue"`
	TotalPrice            *TaxedMoney            `json:"totalPrice"`
	VariantID             *string                `json:"variant"`               // *ProductVariant
	TranslatedProductName string                 `json:"translatedProductName"` //
	TranslatedVariantName string                 `json:"translatedVariantName"` //
	AllocationIDs         []string               `json:"allocations"`           // []*Allocation
	UnitDiscountType      *DiscountValueTypeEnum `json:"unitDiscountType"`
}

func (OrderLine) IsNode() {}

// DatabaseOrderLineToGraphqlOrderLine converts database order line to graphql order line
func DatabaseOrderLineToGraphqlOrderLine(o *order.OrderLine) *OrderLine {

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
		UnitPrice:             NormalTaxedMoneyToGraphqlTaxedMoney(o.UnitPrice),
		UndiscountedUnitPrice: NormalTaxedMoneyToGraphqlTaxedMoney(o.UnDiscountedUnitPrice),
		UnitDiscount:          NormalMoneyToGraphqlMoney(o.UnitDiscount),
		UnitDiscountValue:     o.UnitDiscountValue,
		VariantID:             o.VariantID,
		TranslatedProductName: o.TranslatedProductName,
		TranslatedVariantName: o.TranslatedVariantName,
		AllocationIDs:         []string{},
		UnitDiscountType:      &unitDiscountType,
		TotalPrice:            NormalTaxedMoneyToGraphqlTaxedMoney(o.TotalPrice),
	}
}

// NormalMoneyToGraphqlMoney converts money with amount is Decimal into float-amount money
func NormalMoneyToGraphqlMoney(m *goprices.Money) *Money {
	float64Amount, _ := m.Amount.Float64()

	return &Money{
		Currency: m.Currency,
		Amount:   float64Amount,
	}
}

// NormalTaxedMoneyToGraphqlTaxedMoney
func NormalTaxedMoneyToGraphqlTaxedMoney(t *goprices.TaxedMoney) *TaxedMoney {
	taxMoney, _ := t.Tax()

	return &TaxedMoney{
		Currency: t.Currency,
		Gross:    NormalMoneyToGraphqlMoney(t.Gross),
		Net:      NormalMoneyToGraphqlMoney(t.Net),
		Tax:      NormalMoneyToGraphqlMoney(taxMoney),
	}
}

// Order represents data format will be returned to end user
type Order struct {
	ID                         string                  `json:"id"`
	Created                    time.Time               `json:"created"`
	Status                     OrderStatus             `json:"status"`
	UserID                     *string                 `json:"user"`               // *User
	TrackingClientID           string                  `json:"trackingClientId"`   //
	BillingAddressID           *string                 `json:"billingAddress"`     // *Address
	ShippingAddressID          *string                 `json:"shippingAddress"`    // *Address
	ShippingMethodID           *string                 `json:"shippingMethod"`     // *ShippingMethod
	ShippingMethodName         *string                 `json:"shippingMethodName"` //
	ChannelID                  string                  `json:"channel"`            // Channel
	ShippingPrice              *TaxedMoney             `json:"shippingPrice"`      //
	ShippingTaxRate            float64                 `json:"shippingTaxRate"`    //
	Token                      string                  `json:"token"`              //
	VoucherID                  *string                 `json:"voucher"`            // Voucher
	GiftCardIDs                []string                `json:"giftCards"`          // []*GiftCard
	DisplayGrossPrices         bool                    `json:"displayGrossPrices"`
	CustomerNote               string                  `json:"customerNote"`
	Weight                     *Weight                 `json:"weight"`
	RedirectURL                *string                 `json:"redirectUrl"`
	PrivateMetadata            []*MetadataItem         `json:"privateMetadata"`
	Metadata                   []*MetadataItem         `json:"metadata"`
	FulfillmentIDs             []string                `json:"fulfillments"`             // []*Fulfillment
	LineIDs                    []string                `json:"lines"`                    // []*OrderLine
	Actions                    []*OrderAction          `json:"actions"`                  //
	AvailableShippingMethodIDs []string                `json:"availableShippingMethods"` // []*ShippingMethod
	InvoiceIDs                 []string                `json:"invoices"`                 // []*Invoice
	Number                     *string                 `json:"number"`
	Original                   *string                 `json:"original"`
	Origin                     OrderOriginEnum         `json:"origin"`
	IsPaid                     bool                    `json:"isPaid"`
	PaymentStatus              PaymentChargeStatusEnum `json:"paymentStatus"`
	PaymentStatusDisplay       string                  `json:"paymentStatusDisplay"`
	PaymentIDs                 []string                `json:"payments"` // []*Payment
	Total                      *TaxedMoney             `json:"total"`
	UndiscountedTotal          *TaxedMoney             `json:"undiscountedTotal"`
	Subtotal                   *TaxedMoney             `json:"subtotal"`
	StatusDisplay              *string                 `json:"statusDisplay"`
	CanFinalize                bool                    `json:"canFinalize"`
	TotalAuthorized            *Money                  `json:"totalAuthorized"`
	TotalCaptured              *Money                  `json:"totalCaptured"`
	EventIDs                   []string                `json:"events"` // []*OrderEvent
	TotalBalance               *Money                  `json:"totalBalance"`
	UserEmail                  *string                 `json:"userEmail"`
	IsShippingRequired         bool                    `json:"isShippingRequired"`
	LanguageCodeEnum           LanguageCodeEnum        `json:"languageCodeEnum"`
	DiscountIDs                []string                `json:"discounts"` // []*OrderDiscount
}

func (Order) IsNode()               {}
func (Order) IsObjectWithMetadata() {}

// DatabaseOrderToGraphqlOrder converts 1 database order to 1 graphql order model
func DatabaseOrderToGraphqlOrder(o *order.Order) *Order {

	shippingTaxRate, _ := o.ShippingTaxRate.Float64()

	canFinalize := true
	if o.Status == order.DRAFT {

	}

	return &Order{
		// ID                         : o.Id,
		// Created                    : o.CreateAt,
		// Status                     : OrderStatus(strings.ToUpper(o.Status)),
		// UserID                       : o.UserID,
		// TrackingClientID           : o.TrackingClientID,
		// BillingAddressID           : o.BillingAddressID,
		// ShippingAddressID          : o.ShippingAddressID,
		// ShippingMethodID           : o.ShippingMethodID,
		// ShippingMethodName         : o.ShippingMethodName,
		// ChannelID                  : o.ChannelID,
		// ShippingPrice              : o.ShippingPrice,
		// ShippingTaxRate            : shippingTaxRate,
		// Token                      : o.Token,
		// VoucherID                  : o.VoucherID,
		// // GiftCardIDs                : o.GiftCards,
		// DisplayGrossPrices         : *o.DisplayGrossPrices,
		// CustomerNote               : o.CustomerNote,
		// Weight                     : o.Weight,
		// RedirectURL                : o.RedirectUrl,
		// PrivateMetadata            : MapToGraphqlMetaDataItems(o.PrivateMetadata),
		// Metadata                   : MapToGraphqlMetaDataItems(o.Metadata),
		// FulfillmentIDs             : nil,
		// LineIDs                    : []string{},
		// Actions                    : nil,
		// AvailableShippingMethodIDs : []string{},
		// InvoiceIDs                 : nil,
		// Number                     : &o.Id,
		// Original                   : o.OriginalID,
		// Origin                     : OrderOriginEnum(strings.ToUpper(o.Origin)),
		// IsPaid                     : o.IsFullyPaid(),
		// PaymentStatus              : PaymentChargeStatusEnumCancelled,
		// // PaymentStatusDisplay       : ,
		// PaymentIDs                 : []string{},
		// Total                      : o.Total,
		// UndiscountedTotal          : NormalTaxedMoneyToGraphqlTaxedMoney(o.UnDiscountedTotal),
		// // Subtotal                   : NormalTaxedMoneyToGraphqlTaxedMoney(o.),
		// // StatusDisplay              : o.,
		// CanFinalize                : ,
		// TotalAuthorized            : ,
		// TotalCaptured              : ,
		// EventIDs                   : ,
		// TotalBalance               : ,
		// UserEmail                  : ,
		// IsShippingRequired         : ,
		// LanguageCodeEnum           : ,
		// DiscountIDs                : ,
	}
}

func SystemWeightToGraphqlWeight() *Weight {

}
