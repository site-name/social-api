package api

import (
	"net/http"

	"github.com/samber/lo"
	"github.com/sitename/sitename/app"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/web"
	"gorm.io/gorm"
)

func cleanOrderUpdateShipping(where string, ap app.AppIface, order *model.Order, method *model.ShippingMethod) *model.AppError {
	if order.ShippingAddressID == nil {
		return model.NewAppError(where, "app.order.shipping_address_not_set.app_error", nil, "cannot choose a shipping method for an order without shipping address", http.StatusNotAcceptable)
	}

	validMethods, appErr := ap.Srv().OrderService().GetValidShippingMethodsForOrder(order)
	if appErr != nil {
		return appErr
	}

	if len(validMethods) == 0 || !lo.SomeBy(validMethods, func(item *model.ShippingMethod) bool {
		return item != nil && method != nil && item.Id == method.Id
	}) {
		return model.NewAppError(where, "app.order.shipping_method_not_usable_for_order.app_error", nil, "shipping method cannot be used with this order.", http.StatusNotAcceptable)
	}

	return nil
}

func cleanOrderCancel(where string, app app.AppIface, order *model.Order) *model.AppError {
	if order != nil {
		orderCanCancel, appErr := app.Srv().OrderService().OrderCanCancel(order)
		if appErr != nil {
			return appErr
		}

		if !orderCanCancel {
			return model.NewAppError(where, "app.order.order_can_cancel.app_error", nil, "this order cannot be canceled", http.StatusNotAcceptable)
		}
	}

	return nil
}

// cleanPayment simply checks if payment is nil, return non-nil error.
// return nil otherwise
func cleanPayment(where string, orderPayment *model.Payment) *model.AppError {
	if orderPayment == nil {
		return model.NewAppError(where, "app.order.order_has_no_payment.app_error", nil, "there is no payment for order", http.StatusNotFound)
	}

	return nil
}

func cleanOrderCapture(where string, payment *model.Payment) *model.AppError {
	appErr := cleanPayment(where, payment)
	if appErr != nil {
		return appErr
	}

	if payment.IsActive != nil && !*payment.IsActive {
		return model.NewAppError(where, "app.payment.payment_cannot_capture.app_error", nil, "only pre-authorized payments can be captured", http.StatusNotAcceptable)
	}

	return nil
}

func cleanVoidPayment(where string, payment *model.Payment) *model.AppError {
	appErr := cleanPayment(where, payment)
	if appErr != nil {
		return appErr
	}

	if payment.IsActive != nil && !*payment.IsActive {
		return model.NewAppError(where, "app.payment.payment_cannot_void.app_error", nil, "only pre-authorized payments can be voided", http.StatusNotAcceptable)
	}

	return nil
}

func cleanRefundPayment(where string, payment *model.Payment) *model.AppError {
	appErr := cleanPayment(where, payment)
	if appErr != nil {
		return appErr
	}

	if !payment.CanRefund() {
		return model.NewAppError(where, "app.payment.payment_cannot_refund.app_error", nil, "payment cannot be refunded", http.StatusNotAcceptable)
	}

	return nil
}

// func tryPaymentAction(app app.AppIface, order *model.Order, user *model.User, payment *model.Payment)

func cleanOrderRefund(where string, app app.AppIface, order *model.Order) *model.AppError {
	orderHasGiftcardLines, appErr := app.Srv().GiftcardService().OrderHasGiftcardLines(order)
	if appErr != nil {
		return appErr
	}

	if orderHasGiftcardLines {
		return model.NewAppError(where, "app.order.order_with_giftcard_refund.app_error", nil, "cannot refund order with giftcard lines", http.StatusNotAcceptable)
	}

	return nil
}

func logAndReturnPaymentFailedAppError(where string, ctx *web.Context, tx *gorm.DB, paymentErr *model.PaymentError, order *model.Order, payment *model.Payment) *model.AppError {
	// create payment failed event
	params := model.StringInterface{
		"message": paymentErr.Error(),
	}
	if payment != nil {
		params["gateway"] = payment.GateWay
		params["payment_id"] = payment.Token
	}

	_, appErr := ctx.App.Srv().OrderService().CommonCreateOrderEvent(tx, &model.OrderEventOption{
		OrderID:    order.Id,
		Type:       model.ORDER_EVENT_TYPE_PAYMENT_FAILED,
		UserID:     &ctx.AppContext.Session().UserId,
		Parameters: params,
	})
	if appErr != nil {
		return appErr
	}

	// raise payment failed error
	return model.NewAppError(where, "app.order.payment_failed.app_error", nil, paymentErr.Error(), http.StatusInternalServerError)
}
