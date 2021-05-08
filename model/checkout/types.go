package checkout

import (
	"github.com/sitename/sitename/model/account"
	"github.com/sitename/sitename/model/channel"
	"github.com/sitename/sitename/model/product_and_discount"
	"github.com/sitename/sitename/model/shipping"
)

type CheckoutLineInfo struct {
	Line           *CheckoutLine
	Variant        *product_and_discount.ProductVariant
	ChannelListing *product_and_discount.ProductVariantChannelListing
	Product        *product_and_discount.Product
	ProductType    *product_and_discount.ProductType
	Collections    []*product_and_discount.Collection
}

type CheckoutInfo struct {
	Checkout                      *Checkout // required
	User                          *account.User
	Channel                       *channel.Channel
	BillingAddress                *account.Address
	ShippingAddress               *account.Address
	ShippingMethod                *shipping.ShippingMethod
	ValidShippingMethods          []*shipping.ShippingMethod
	ShippingMethodChannelListings *shipping.ShippingMethodChannelListing
}

func (c *CheckoutInfo) GetCountry() string {
	var add *account.Address
	if c.ShippingAddress != nil {
		add = c.ShippingAddress
	} else if c.BillingAddress != nil {
		add = c.BillingAddress
	}

	if add == nil || add.Country == "" {
		return c.Checkout.Country
	}

	return add.Country
}
