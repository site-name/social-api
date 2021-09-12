package checkout

import (
	"net/http"

	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model/checkout"
	"github.com/sitename/sitename/model/payment"
	"github.com/sitename/sitename/model/product_and_discount"
)

func (a *ServiceCheckout) CleanCheckoutShipping(checkoutInfo *checkout.CheckoutInfo, lines checkout.CheckoutLineInfos) *model.AppError {
	productIDs := lines.Products().IDs()

	requireShipping, appErr := a.srv.ProductService().ProductsRequireShipping(productIDs)
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

		if !deliveruMethodInfo.IsMethodInValidMethods(checkoutInfo) {
			appErr = a.ClearDeliveryMethod(checkoutInfo)
			if appErr != nil {
				return appErr
			}
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

func (a *ServiceCheckout) CleanCheckoutPayment(manager interface{}, checkoutInfo *checkout.CheckoutInfo, lines []*checkout.CheckoutLineInfo, discounts []*product_and_discount.DiscountInfo, lastPayment *payment.Payment) (*payment.PaymentError, *model.AppError) {
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
