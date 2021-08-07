package gqlmodel

import (
	"strings"
	"time"

	"github.com/site-name/decimal"
	goprices "github.com/site-name/go-prices"
	"github.com/sitename/sitename/model/order"
	"github.com/sitename/sitename/modules/measurement"
	"github.com/sitename/sitename/modules/util"
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
	UnitPrice             *TaxedMoney            `json:"unitPrice"`
	UndiscountedUnitPrice *TaxedMoney            `json:"undiscountedUnitPrice"`
	UnitDiscount          *Money                 `json:"unitDiscount"`
	UnitDiscountValue     *decimal.Decimal       `json:"unitDiscountValue"`
	TotalPrice            *TaxedMoney            `json:"totalPrice"`
	TranslatedProductName string                 `json:"translatedProductName"`
	TranslatedVariantName string                 `json:"translatedVariantName"`
	UnitDiscountType      *DiscountValueTypeEnum `json:"unitDiscountType"`
	Thumbnail             func() *Image          `json:"thumbnail"`         // *Image
	DigitalContentURLID   *string                `json:"digitalContentUrl"` // *DigitalContentURL
	VariantID             *string                `json:"variant"`           // *ProductVariant
	AllocationIDs         []string               `json:"allocations"`       // []*Allocation
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
		Quantity:              int(o.Quantity),
		QuantityFulfilled:     int(o.QuantityFulfilled),
		UnitDiscountReason:    o.UnitDiscountReason,
		TaxRate:               taxRate,
		UnitPrice:             NormalTaxedMoneyToGraphqlTaxedMoney(o.UnitPrice),
		UndiscountedUnitPrice: NormalTaxedMoneyToGraphqlTaxedMoney(o.UnDiscountedUnitPrice),
		UnitDiscount:          NormalMoneyToGraphqlMoney(o.UnitDiscount),
		UnitDiscountValue:     o.UnitDiscountValue,
		VariantID:             o.VariantID,
		TranslatedProductName: o.TranslatedProductName,
		TranslatedVariantName: o.TranslatedVariantName,
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

// Order represents data of a typical merchandise order
// When initialize new Order, you can ignore fields that have type of functions
type Order struct {
	ID                       string                         `json:"id"`
	Created                  time.Time                      `json:"created"`
	Status                   OrderStatus                    `json:"status"`
	TrackingClientID         string                         `json:"trackingClientId"`
	ShippingMethodName       *string                        `json:"shippingMethodName"`
	ShippingPrice            *TaxedMoney                    `json:"shippingPrice"`
	ShippingTaxRate          float64                        `json:"shippingTaxRate"`
	Token                    string                         `json:"token"`
	UserID                   *string                        `json:"user"`            // *User
	BillingAddressID         *string                        `json:"billingAddress"`  // *Address
	ShippingAddressID        *string                        `json:"shippingAddress"` // *Address
	ShippingMethodID         *string                        `json:"shippingMethod"`  // *ShippingMethod
	ChannelID                string                         `json:"channel"`         // Channel
	VoucherID                *string                        `json:"voucher"`         // Voucher
	DisplayGrossPrices       bool                           `json:"displayGrossPrices"`
	CustomerNote             string                         `json:"customerNote"`
	Weight                   *Weight                        `json:"weight"`
	RedirectURL              *string                        `json:"redirectUrl"`
	PrivateMetadata          []*MetadataItem                `json:"privateMetadata"`
	Metadata                 []*MetadataItem                `json:"metadata"`
	Total                    *TaxedMoney                    `json:"total"`
	UndiscountedTotal        *TaxedMoney                    `json:"undiscountedTotal"`
	Number                   *string                        `json:"number"`
	Original                 *string                        `json:"original"`
	Origin                   OrderOriginEnum                `json:"origin"`
	IsPaid                   bool                           `json:"isPaid"`
	TotalBalance             *Money                         `json:"totalBalance"`
	UserEmail                *string                        `json:"userEmail"`
	LanguageCodeEnum         LanguageCodeEnum               `json:"languageCodeEnum"`
	PaymentStatus            func() PaymentChargeStatusEnum `json:"paymentStatus"`            // PaymentChargeStatusEnum
	PaymentStatusDisplay     func() string                  `json:"paymentStatusDisplay"`     // string
	Payments                 func() []*Payment              `json:"payments"`                 // []*Payment
	GiftCards                func() []*GiftCard             `json:"giftCards"`                // []*GiftCard
	Fulfillments             func() []*Fulfillment          `json:"fulfillments"`             // []*Fulfillment
	Lines                    func() []*OrderLine            `json:"lines"`                    // []*OrderLine
	Actions                  func() []*OrderAction          `json:"actions"`                  // []*OrderAction
	AvailableShippingMethods func() []*ShippingMethod       `json:"availableShippingMethods"` // []*ShippingMethod
	Invoices                 func() []*Invoice              `json:"invoices"`                 // []*Invoice
	Subtotal                 func() *TaxedMoney             `json:"subtotal"`                 // *TaxedMoney
	StatusDisplay            func() *string                 `json:"statusDisplay"`            // *string
	CanFinalize              func() bool                    `json:"canFinalize"`              // bool
	TotalAuthorized          func() *Money                  `json:"totalAuthorized"`          // *Money
	TotalCaptured            func() *Money                  `json:"totalCaptured"`            // *Money
	Events                   func() []*OrderEvent           `json:"events"`                   // []*OrderEvent
	IsShippingRequired       func() bool                    `json:"isShippingRequired"`       // bool
	Discounts                func() []*OrderDiscount        `json:"discounts"`                // []*OrderDiscount
}

func (Order) IsNode()               {}
func (Order) IsObjectWithMetadata() {}

// DatabaseOrderToGraphqlOrder converts 1 database order to 1 graphql order model
func DatabaseOrderToGraphqlOrder(o *order.Order) *Order {

	shippingTaxRate, _ := o.ShippingTaxRate.Float64()
	totalBalance, _ := o.TotalBalance()

	return &Order{
		ID:                 o.Id,
		Created:            util.TimeFromMillis(o.CreateAt),
		Status:             OrderStatus(strings.ToUpper(o.Status)),
		TrackingClientID:   o.TrackingClientID,
		ShippingMethodName: o.ShippingMethodName,
		ShippingPrice:      NormalTaxedMoneyToGraphqlTaxedMoney(o.ShippingPrice),
		ShippingTaxRate:    shippingTaxRate,
		Token:              o.Token,
		UserID:             o.UserID,
		BillingAddressID:   o.BillingAddressID,
		ShippingAddressID:  o.ShippingAddressID,
		ShippingMethodID:   o.ShippingMethodID,
		ChannelID:          o.ChannelID,
		VoucherID:          o.VoucherID,
		DisplayGrossPrices: *o.DisplayGrossPrices,
		CustomerNote:       o.CustomerNote,
		Weight:             NormalWeightToGraphqlWeight(o.Weight),
		RedirectURL:        o.RedirectUrl,
		PrivateMetadata:    MapToGraphqlMetaDataItems(o.PrivateMetadata),
		Metadata:           MapToGraphqlMetaDataItems(o.Metadata),
		Total:              NormalTaxedMoneyToGraphqlTaxedMoney(o.Total),
		UndiscountedTotal:  NormalTaxedMoneyToGraphqlTaxedMoney(o.UnDiscountedTotal),
		Number:             &o.Id,
		Original:           o.OriginalID,
		Origin:             OrderOriginEnum(strings.ToUpper(o.Origin)),
		IsPaid:             o.IsFullyPaid(),
		TotalBalance:       NormalMoneyToGraphqlMoney(totalBalance),
		UserEmail:          &o.UserEmail,
		LanguageCodeEnum:   LanguageCodeEnum(strings.ToUpper(o.LanguageCode)),
	}
}

// NormalWeightToGraphqlWeight converts weight to graphql weight
func NormalWeightToGraphqlWeight(w *measurement.Weight) *Weight {
	return &Weight{
		Value: float64(*w.Amount),
		Unit:  WeightUnitsEnum(strings.ToUpper(string(w.Unit))),
	}
}
