package checkout

import (
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model/account"
	"github.com/sitename/sitename/model/checkout"
	"github.com/sitename/sitename/model/product_and_discount"
	"github.com/sitename/sitename/model/shipping"
)

func (a *AppCheckout) FetchCheckoutLines(checkout *checkout.Checkout) ([]*checkout.CheckoutLineInfo, *model.AppError) {
	panic("not implt")
}

func (a *AppCheckout) FetCheckoutInfo(ckout *checkout.Checkout, lines []*checkout.CheckoutLineInfo, discounts []*product_and_discount.DiscountInfo) (*checkout.CheckoutInfo, *model.AppError) {
	panic("not implt")

}

func (a *AppCheckout) UpdateCheckoutInfoShippingAddress(checkoutInfo *checkout.CheckoutInfo, address *account.Address, lines []*checkout.CheckoutLineInfo) *model.AppError {
	panic("not implt")

}

func (a *AppCheckout) GetValidShippingMethodListForCheckoutInfo(checkoutInfo *checkout.CheckoutInfo, shippingAddress *account.Address, lines []*checkout.CheckoutLineInfo, discounts []*product_and_discount.DiscountInfo) ([]*shipping.ShippingMethod, *model.AppError) {
	panic("not implt")

}

func (a *AppCheckout) UpdateCheckoutInfoShippingMethod(checkoutInfo *checkout.CheckoutInfo, shippingMethod *shipping.ShippingMethod) *model.AppError {
	checkoutInfo.ShippingMethod = shippingMethod
	checkoutInfo.ShippingMethodChannelListings = nil

	if shippingMethod != nil {
		panic("not implt")

	}
	panic("not implt")

}
