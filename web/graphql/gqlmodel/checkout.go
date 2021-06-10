package gqlmodel

import (
	"time"

	"github.com/google/uuid"
)

type Checkout struct {
	Created                    time.Time        `json:"created"`
	LastChange                 time.Time        `json:"lastChange"`
	UserID                     *string          `json:"user"`            // User
	ChannelID                  *string          `json:"channel"`         // Channel
	BillingAddressID           *string          `json:"billingAddress"`  // Address
	ShippingAddressID          *string          `json:"shippingAddress"` // Address
	Note                       string           `json:"note"`
	Discount                   *Money           `json:"discount"`
	DiscountName               *string          `json:"discountName"`
	TranslatedDiscountName     *string          `json:"translatedDiscountName"`
	VoucherCode                *string          `json:"voucherCode"`
	GiftCardIDs                []*string        `json:"giftCards"` // GiftCard
	ID                         string           `json:"id"`
	PrivateMetadata            []*MetadataItem  `json:"privateMetadata"`
	Metadata                   []*MetadataItem  `json:"metadata"`
	AvailableShippingMethodIDs []*string        `json:"availableShippingMethods"` // ShippingMethod
	AvailablePaymentGatewayIDs []string         `json:"availablePaymentGateways"` // PaymentGateway
	Email                      string           `json:"email"`
	IsShippingRequired         bool             `json:"isShippingRequired"`
	Quantity                   int              `json:"quantity"`
	LineIDs                    []*string        `json:"lines"` // CheckoutLine
	ShippingPrice              *TaxedMoney      `json:"shippingPrice"`
	ShippingMethod             *ShippingMethod  `json:"shippingMethod"`
	SubtotalPrice              *TaxedMoney      `json:"subtotalPrice"`
	Token                      uuid.UUID        `json:"token"`
	TotalPrice                 *TaxedMoney      `json:"totalPrice"`
	LanguageCode               LanguageCodeEnum `json:"languageCode"`
}

func (Checkout) IsNode()               {}
func (Checkout) IsObjectWithMetadata() {}
