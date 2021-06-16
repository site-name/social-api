package gqlmodel

import (
	"strings"

	"time"

	"github.com/shopspring/decimal"
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

func NormalMoneyToGraphqlMoney(m *goprices.Money) *Money {
	float64Amount, _ := m.Amount.Float64()

	return &Money{
		Currency: m.Currency,
		Amount:   float64Amount,
	}
}

func NormalTaxedMoneyToGraphqlTaxedMoney(t *goprices.TaxedMoney) *TaxedMoney {
	taxMoney, _ := t.Tax()

	return &TaxedMoney{
		Currency: t.Currency,
		Gross:    NormalMoneyToGraphqlMoney(t.Gross),
		Net:      NormalMoneyToGraphqlMoney(t.Net),
		Tax:      NormalMoneyToGraphqlMoney(taxMoney),
	}
}

type Order struct {
	ID                         string                  `json:"id"`
	Created                    time.Time               `json:"created"`
	Status                     OrderStatus             `json:"status"`
	User                       *User                   `json:"user"`
	TrackingClientID           string                  `json:"trackingClientId"`
	BillingAddressID           *string                 `json:"billingAddress"`     // *Address
	ShippingAddressID          *string                 `json:"shippingAddress"`    // *Address
	ShippingMethodID           *string                 `json:"shippingMethod"`     // *ShippingMethod
	ShippingMethodName         *string                 `json:"shippingMethodName"` //
	ChannelID                  *string                 `json:"channel"`            // Channel
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
