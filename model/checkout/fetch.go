package checkout

import (
	"github.com/sitename/sitename/model/account"
	"github.com/sitename/sitename/model/channel"
	"github.com/sitename/sitename/model/product_and_discount"
	"github.com/sitename/sitename/model/shipping"
)

// CheckoutLineInfo contains information of a checkout line
type CheckoutLineInfo struct {
	Line           CheckoutLine
	Variant        *product_and_discount.ProductVariant
	ChannelListing product_and_discount.ProductVariantChannelListing
	Product        product_and_discount.Product
	ProductType    product_and_discount.ProductType
	Collections    []*product_and_discount.Collection
}

// CheckoutInfo contains information of a checkout
type CheckoutInfo struct {
	Checkout                      Checkout
	User                          *account.User
	Channel                       channel.Channel
	BillingAddress                *account.Address
	ShippingAddress               *account.Address
	ShippingMethod                *shipping.ShippingMethod
	ValidShippingMethods          []*shipping.ShippingMethod
	ShippingMethodChannelListings *shipping.ShippingMethodChannelListing
}

func (c *CheckoutInfo) GetCountry() string {
	addr := c.ShippingAddress
	if addr == nil {
		addr = c.BillingAddress
	}

	if addr == nil || addr.Country == "" {
		return c.Checkout.Country
	}

	return addr.Country
}

func (c *CheckoutInfo) GetCustomerEmail() string {
	if c.User != nil {
		return c.User.Email
	}
	return c.Checkout.Email
}
