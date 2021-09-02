package checkout

import (
	"net/http"

	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model/account"
	"github.com/sitename/sitename/model/checkout"
	"github.com/sitename/sitename/model/product_and_discount"
	"github.com/sitename/sitename/model/shipping"
)

// FetchCheckoutLines Fetch checkout lines as CheckoutLineInfo objects.
// It prefetch some related value also
func (a *ServiceCheckout) FetchCheckoutLines(ckout *checkout.Checkout) ([]*checkout.CheckoutLineInfo, *model.AppError) {
	checkoutLineInfos, err := a.srv.Store.Checkout().FetchCheckoutLinesAndPrefetchRelatedValue(ckout)
	if err != nil {
		return nil, model.NewAppError("FetchCheckoutLines", "app.checkout.error_collecting_checkout_line_infos.app_error", nil, err.Error(), http.StatusInternalServerError)
	}

	return checkoutLineInfos, nil
}

func (a *ServiceCheckout) FetCheckoutInfo(ckout *checkout.Checkout, lines []*checkout.CheckoutLineInfo, discounts []*product_and_discount.DiscountInfo) (*checkout.CheckoutInfo, *model.AppError) {
	panic("not implt")

}

func (a *ServiceCheckout) UpdateCheckoutInfoShippingAddress(checkoutInfo *checkout.CheckoutInfo, address *account.Address, lines []*checkout.CheckoutLineInfo) *model.AppError {
	panic("not implt")

}

func (a *ServiceCheckout) GetValidShippingMethodListForCheckoutInfo(checkoutInfo *checkout.CheckoutInfo, shippingAddress *account.Address, lines []*checkout.CheckoutLineInfo, discounts []*product_and_discount.DiscountInfo) ([]*shipping.ShippingMethod, *model.AppError) {
	panic("not implt")

}

// UpdateCheckoutInfoShippingMethod set CheckoutInfo's ShippingMethod to given shippingMethod
// and set new value for checkoutInfo's ShippingMethodChannelListings
func (a *ServiceCheckout) UpdateCheckoutInfoShippingMethod(checkoutInfo *checkout.CheckoutInfo, shippingMethod *shipping.ShippingMethod) *model.AppError {
	checkoutInfo.ShippingMethod = shippingMethod

	checkoutInfo.ShippingMethodChannelListings = nil
	if shippingMethod != nil {
		listings, appErr := a.srv.ShippingService().ShippingMethodChannelListingsByOption(&shipping.ShippingMethodChannelListingFilterOption{
			ShippingMethodID: &model.StringFilter{
				StringOption: &model.StringOption{
					Eq: shippingMethod.Id,
				},
			},
			ChannelID: &model.StringFilter{
				StringOption: &model.StringOption{
					Eq: checkoutInfo.Channel.Id,
				},
			},
		})

		if appErr != nil {
			return appErr
		}
		checkoutInfo.ShippingMethodChannelListings = listings[0]
	}

	return nil
}
