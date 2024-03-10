package checkout

import (
	"net/http"

	"github.com/samber/lo"
	"github.com/sitename/sitename/app/plugin/interfaces"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model_helper"
	"github.com/volatiletech/sqlboiler/v4/boil"
)

func (a *ServiceCheckout) CleanCheckoutShipping(checkoutInfo model_helper.CheckoutInfo, lines model_helper.CheckoutLineInfos) *model_helper.AppError {
	productIDs := lo.Map(lines.Products(), func(item *model.Product, _ int) string { return item.ID })
	requireShipping, appErr := a.srv.Product.ProductsRequireShipping(productIDs)
	if appErr != nil {
		return appErr
	}

	if requireShipping {
		deliveruMethodInfo := checkoutInfo.DeliveryMethodInfo

		if deliveruMethodInfo.GetDeliveryMethod() == nil {
			return model_helper.NewAppError("CleanCheckoutShipping", "app.discount.shipping_method_not_set.app_error", nil, "", http.StatusNotImplemented)
		}

		if !deliveruMethodInfo.IsValidDeliveryMethod() {
			return model_helper.NewAppError("CleanCheckoutShipping", "app.discount.shipping_address_not_set.app_error", nil, "", http.StatusNotImplemented)
		}

		if !deliveruMethodInfo.IsMethodInValidMethods(checkoutInfo) {
			appErr = a.ClearDeliveryMethod(checkoutInfo)
			if appErr != nil {
				return appErr
			}
			return model_helper.NewAppError("CleanCheckoutShipping", "app.discount.shipping_method_not_valid_for_shipping_address.app_error", nil, "", http.StatusNotImplemented)
		}
	}

	return nil
}

func (a *ServiceCheckout) CleanBillingAddress(checkoutInfo model_helper.CheckoutInfo) *model_helper.AppError {
	if checkoutInfo.BillingAddress == nil {
		return model_helper.NewAppError("CleanBillingAddress", "app.discount.billing_address_not_set.app_error", nil, "", http.StatusNotImplemented)
	}

	return nil
}

func (a *ServiceCheckout) CleanCheckoutPayment(tx boil.ContextTransactor, manager interfaces.PluginManagerInterface, checkoutInfo model_helper.CheckoutInfo, lines model_helper.CheckoutLineInfos, discounts []*model_helper.DiscountInfo, lastPayment *model.Payment) (*model_helper.PaymentError, *model_helper.AppError) {
	if appErr := a.CleanBillingAddress(checkoutInfo); appErr != nil {
		return nil, appErr
	}

	isFullyPaid, appErr := a.IsFullyPaid(manager, checkoutInfo, lines, discounts)
	if appErr != nil {
		return nil, appErr
	}

	if !isFullyPaid {
		paymentErr, appErr := a.srv.Payment.PaymentRefundOrVoid(tx, lastPayment, manager, checkoutInfo.Channel.Slug)
		if paymentErr != nil || appErr != nil {
			return paymentErr, appErr
		}

		return nil, model_helper.NewAppError("CleanCHeckoutPayment", "app.checkout.checkout_not_fully_paid.app_error", nil, "", http.StatusNotAcceptable)
	}

	return nil, nil
}
