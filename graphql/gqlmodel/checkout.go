package gqlmodel

import (
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/sitename/sitename/model/checkout"
	"github.com/sitename/sitename/modules/util"
)

// OLD implementation

// type Checkout struct {
// 	Created                   time.Time         `json:"created"`
// 	LastChange                time.Time         `json:"lastChange"`
// 	User                      *User             `json:"user"`
// 	Channel                   *Channel          `json:"channel"`
// 	BillingAddress            *Address          `json:"billingAddress"`
// 	ShippingAddress           *Address          `json:"shippingAddress"`
// 	Note                      string            `json:"note"`
// 	Discount                  *Money            `json:"discount"`
// 	DiscountName              *string           `json:"discountName"`
// 	TranslatedDiscountName    *string           `json:"translatedDiscountName"`
// 	VoucherCode               *string           `json:"voucherCode"`
// 	GiftCards                 []*GiftCard       `json:"giftCards"`
// 	ID                        string            `json:"id"`
// 	PrivateMetadata           []*MetadataItem   `json:"privateMetadata"`
// 	Metadata                  []*MetadataItem   `json:"metadata"`
// 	AvailableShippingMethods  []*ShippingMethod `json:"availableShippingMethods"`
// 	AvailableCollectionPoints []*Warehouse      `json:"availableCollectionPoints"`
// 	AvailablePaymentGateways  []*PaymentGateway `json:"availablePaymentGateways"`
// 	Email                     string            `json:"email"`
// 	IsShippingRequired        bool              `json:"isShippingRequired"`
// 	Quantity                  int               `json:"quantity"`
// 	Lines                     []*CheckoutLine   `json:"lines"`
// 	ShippingPrice             *TaxedMoney       `json:"shippingPrice"`
// 	DeliveryMethod            DeliveryMethod    `json:"deliveryMethod"`
// 	SubtotalPrice             *TaxedMoney       `json:"subtotalPrice"`
// 	Token                     uuid.UUID         `json:"token"`
// 	TotalPrice                *TaxedMoney       `json:"totalPrice"`
// 	LanguageCode              LanguageCodeEnum  `json:"languageCode"`
// }

type Checkout struct {
	ID                          string            `json:"id"`
	Created                     time.Time         `json:"created"`
	LastChange                  time.Time         `json:"lastChange"`
	UserID                      *string           `json:"user"`
	ChannelID                   *string           `json:"channel"`
	BillingAddressID            *string           `json:"billingAddress"`
	ShippingAddressID           *string           `json:"shippingAddress"`
	Note                        string            `json:"note"`
	DiscountName                *string           `json:"discountName"`
	TranslatedDiscountName      *string           `json:"translatedDiscountName"`
	VoucherCode                 *string           `json:"voucherCode"`
	GiftCardIDs                 []string          `json:"giftCards"` //
	PrivateMetadata             []*MetadataItem   `json:"privateMetadata"`
	Metadata                    []*MetadataItem   `json:"metadata"`
	AvailableShippingMethodIDs  []string          `json:"availableShippingMethods"`  //
	AvailableCollectionPointIDs []string          `json:"availableCollectionPoints"` //
	AvailablePaymentGateways    []*PaymentGateway `json:"availablePaymentGateways"`
	Email                       string            `json:"email"`
	IsShippingRequired          bool              `json:"isShippingRequired"`
	Quantity                    int               `json:"quantity"`
	LineIDs                     []string          `json:"lines"`
	ShippingPrice               *TaxedMoney       `json:"shippingPrice"`
	DeliveryMethod              DeliveryMethod    `json:"deliveryMethod"`
	SubtotalPrice               *TaxedMoney       `json:"subtotalPrice"`
	Token                       uuid.UUID         `json:"token"`
	TotalPrice                  *TaxedMoney       `json:"totalPrice"`
	LanguageCode                LanguageCodeEnum  `json:"languageCode"`
	Discount                    *Money            `json:"discount"`
}

func (Checkout) IsNode()               {}
func (Checkout) IsObjectWithMetadata() {}

// SystemCheckoutToGraphqlCheckout converts given system checkout to graphql checkout
func SystemCheckoutToGraphqlCheckout(c *checkout.Checkout) *Checkout {
	if c == nil {
		return nil
	}

	res := &Checkout{
		ID:                     c.Token,
		Created:                util.TimeFromMillis(c.CreateAt),
		LastChange:             util.TimeFromMillis(c.UpdateAt),
		UserID:                 c.UserID,
		ChannelID:              &c.ChannelID,
		BillingAddressID:       c.BillingAddressID,
		ShippingAddressID:      c.ShippingAddressID,
		Note:                   c.Note,
		DiscountName:           c.DiscountName,
		TranslatedDiscountName: c.TranslatedDiscountName,
		VoucherCode:            c.VoucherCode,
		PrivateMetadata:        MapToGraphqlMetaDataItems(c.PrivateMetadata),
		Metadata:               MapToGraphqlMetaDataItems(c.Metadata),
		Email:                  c.Email,
		LanguageCode:           LanguageCodeEnum(strings.ToUpper(c.LanguageCode)),
		// Quantity:               int(c.Quantity),
		// IsShippingRequired: ,
		// AvailablePaymentGateways: ,
	}

	if c.DiscountAmount != nil {
		fl, _ := c.DiscountAmount.Float64()
		res.Discount = &Money{
			Amount:   fl,
			Currency: c.Currency,
		}
	}

	return res
}
