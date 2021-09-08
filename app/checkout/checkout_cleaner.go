package checkout

import (
	"net/http"

	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model/checkout"
)

func (a *ServiceCheckout) CleanCheckoutShipping(checkoutInfo *checkout.CheckoutInfo, lines []*checkout.CheckoutLineInfo) *model.AppError {
	productIDs := []string{}
	for _, line := range lines {
		productIDs = append(productIDs, line.Product.Id)
	}

	requireShipping, appErr := a.srv.ProductService().ProductsRequireShipping(productIDs)
	if appErr != nil {
		return appErr
	}

	if requireShipping {
		deliveruMethodInfo := checkoutInfo.DeliveryMethodInfo

		if deliveruMethodInfo.DeliveryMethod == nil {
			return model.NewAppError("CleanCheckoutShipping", "app.discount.shipping_method_not_set.app_error", nil, "", http.StatusNotImplemented)
		}

		if !deliveruMethodInfo.IsValidDeliveryMethod() {
			return model.NewAppError("CleanCheckoutShipping", "app.discount.shipping_address_not_set.app_error", nil, "", http.StatusNotImplemented)
		}

		if !deliveruMethodInfo.IsMethodInValidMethods(checkoutInfo) {
			a.ClearDeliveryMethod(checkoutInfo)
			return model.NewAppError("CleanCheckoutShipping", "app.discount.shipping_method_not_valid_for_shipping_address.app_error", nil, "", http.StatusNotImplemented)
		}
	}

	return nil
}

func (a *ServiceCheckout) CleanBillingAddress(checkoutInfo *checkout.CheckoutInfo) *model.AppError {
	if checkoutInfo.BillingAddress == nil {
		return model.NewAppError("CleanBillingAddress", "app.discount.billing_address_not_set.app_error", nil, "", http.StatusNotImplemented)
	}

	return nil
}

func (a *ServiceCheckout) CleanCheckoutPayment() {
	panic("not implemented")
}
