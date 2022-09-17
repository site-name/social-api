package checkout

import (
	"net/http"

	"github.com/sitename/sitename/app/plugin/interfaces"
	"github.com/sitename/sitename/model"
)

// CleanCheckoutShipping
func (a *ServiceCheckout) CleanCheckoutShipping(checkoutInfo model.CheckoutInfo, lines model.CheckoutLineInfos) *model.AppError {
	requireShipping, appErr := a.srv.ProductService().ProductsRequireShipping(lines.Products().IDs())
	if appErr != nil {
		return appErr
	}

	if requireShipping {
		deliveruMethodInfo := checkoutInfo.DeliveryMethodInfo

		if deliveruMethodInfo.GetDeliveryMethod() == nil {
			return model.NewAppError("CleanCheckoutShipping", "app.discount.shipping_method_not_set.app_error", nil, "", http.StatusNotImplemented)
		}

		if !deliveruMethodInfo.IsValidDeliveryMethod() {
			return model.NewAppError("CleanCheckoutShipping", "app.discount.shipping_address_not_set.app_error", nil, "", http.StatusNotImplemented)
		}

		if !deliveruMethodInfo.IsMethodInValidMethods(&checkoutInfo) {
			appErr = a.ClearDeliveryMethod(checkoutInfo)
			if appErr != nil {
				return appErr
			}
			return model.NewAppError("CleanCheckoutShipping", "app.discount.shipping_method_not_valid_for_shipping_address.app_error", nil, "", http.StatusNotImplemented)
		}
	}

	return nil
}

func (a *ServiceCheckout) CleanBillingAddress(checkoutInfo model.CheckoutInfo) *model.AppError {
	if checkoutInfo.BillingAddress == nil {
		return model.NewAppError("CleanBillingAddress", "app.discount.billing_address_not_set.app_error", nil, "", http.StatusNotImplemented)
	}

	return nil
}

func (a *ServiceCheckout) CleanCheckoutPayment(manager interfaces.PluginManagerInterface, checkoutInfo model.CheckoutInfo, lines []*model.CheckoutLineInfo, discounts []*model.DiscountInfo, lastPayment *model.Payment) (*model.PaymentError, *model.AppError) {
	if appErr := a.CleanBillingAddress(checkoutInfo); appErr != nil {
		return nil, appErr
	}

	isFullyPaid, appErr := a.IsFullyPaid(manager, checkoutInfo, lines, discounts)
	if appErr != nil {
		return nil, appErr
	}

	if !isFullyPaid {
		paymentErr, appErr := a.srv.PaymentService().PaymentRefundOrVoid(lastPayment, manager, checkoutInfo.Channel.Slug)
		if paymentErr != nil || appErr != nil {
			return paymentErr, appErr
		}

		return nil, model.NewAppError("CleanCHeckoutPayment", "app.checkout.checkout_not_fully_paid.app_error", nil, "", http.StatusNotAcceptable)
	}

	return nil, nil
}
