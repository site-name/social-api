package api

import (
	"fmt"
	"net/http"

	"github.com/Masterminds/squirrel"
	"github.com/samber/lo"
	"github.com/site-name/decimal"
	"github.com/sitename/sitename/app"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/web"
	"gorm.io/gorm"
)

const (
	ErrCannotBeRefunded                                  = "app.order.order_cannot_refund.app_error"
	ErrCannotDeclareRefundAmountWhenOrerHasGiftCardLines = "app.order.cannot_have_refund_amount_when_order_has_giftcard.app_error"
	ErrRefundAmountGreaterThanPossible                   = "app.order.refund_amount_greater_than_possible_amount.app_error"
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

func cleanOrderPayment(where string, payment *model.Payment) *model.AppError {
	if payment == nil || !payment.CanRefund() {
		return model.NewAppError(where, ErrCannotBeRefunded, nil, "order cannot be refunded", http.StatusNotAcceptable)
	}
	return nil
}

func cleanAmountToRefund(embedCtx *web.Context, where string, order *model.Order, payment *model.Payment, amountToRefund *decimal.Decimal) *model.AppError {
	if amountToRefund != nil {
		orderHasGiftCardLines, appErr := embedCtx.App.Srv().GiftcardService().OrderHasGiftcardLines(order)
		if appErr != nil {
			return appErr
		}

		if orderHasGiftCardLines {
			return model.NewAppError(where, ErrCannotDeclareRefundAmountWhenOrerHasGiftCardLines, nil, "cannot specify refund amount when orer has giftcard lines", http.StatusNotAcceptable)
		}

		if payment.CapturedAmount != nil && amountToRefund.GreaterThan(*payment.CapturedAmount) {
			return model.NewAppError(where, ErrRefundAmountGreaterThanPossible, nil, "the required refund amount greater than possible amount to refund", http.StatusNotAcceptable)
		}

		return nil
	}

	return nil
}

func cleanLines(embedCtx *web.Context, where string, linesData []*OrderRefundLineInput) (model.OrderLineDatas, *model.AppError) {
	orderLineIds := lo.Map(linesData, func(item *OrderRefundLineInput, _ int) string { return item.OrderLineID.String() })

	orderLines, appErr := embedCtx.App.Srv().OrderService().OrderLinesByOption(&model.OrderLineFilterOption{
		Conditions: squirrel.Eq{model.OrderLineTableName + ".Id": orderLineIds},
	})
	if appErr != nil {
		return nil, appErr
	}

	lineDataMap := lo.SliceToMap(linesData, func(item *OrderRefundLineInput) (string, *OrderRefundLineInput) {
		return item.OrderLineID.String(), item
	})

	var res model.OrderLineDatas

	for i := 0; i < min(len(orderLines), len(lineDataMap)); i++ {
		orderLine := orderLines[i]
		lineData := lineDataMap[orderLine.Id]
		quantity := int(lineData.Quantity)

		if orderLine.IsGiftcard {
			return nil, model.NewAppError(where, "app.order.cannot_refund_giftcard_line.app_error", nil, "cannot refund or return giftcard line", http.StatusNotAcceptable)
		}

		if orderLine.Quantity < quantity {
			return nil, model.NewAppError(where, model.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "quantity"}, fmt.Sprintf("provided quantity: %d bigger than order line quantity: %d", lineData.Quantity, orderLine.Quantity), http.StatusBadRequest)
		}

		if unfulfilledQuantity := orderLine.QuantityUnFulfilled(); unfulfilledQuantity < quantity {
			return nil, model.NewAppError(where, model.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "quantity"}, fmt.Sprintf("provided quantity: %d bigger than order line unfulfilled quantity: %d", lineData.Quantity, unfulfilledQuantity), http.StatusBadRequest)
		}

		res = append(res, &model.OrderLineData{
			Line:     *orderLine,
			Quantity: quantity,
		})
	}

	return res, nil
}

func cleanFulfillmentLines(embedCtx *web.Context, where string, fulfillmentLinesData []*OrderRefundFulfillmentLineInput, whitelistedStatuses []model.FulfillmentStatus) ([]*model.FulfillmentLineData, *model.AppError) {
	fulfillmentLineIds := lo.Map(fulfillmentLinesData, func(item *OrderRefundFulfillmentLineInput, _ int) string { return item.FulfillmentLineID.String() })

	fulfillmentLines, appErr := embedCtx.App.
		Srv().
		OrderService().
		FulfillmentLinesByOption(&model.FulfillmentLineFilterOption{
			Conditions: squirrel.Eq{model.FulfillmentLineTableName + ".Id": fulfillmentLineIds},
			Preloads:   []string{"Fulfillment", "OrderLine"},
		})

	if appErr != nil {
		return nil, appErr
	}

	fulfillmentLineDataMap := lo.SliceToMap(fulfillmentLinesData, func(item *OrderRefundFulfillmentLineInput) (string, *OrderRefundFulfillmentLineInput) {
		return item.FulfillmentLineID.String(), item
	})

	res := []*model.FulfillmentLineData{}

	for i := 0; i < min(len(fulfillmentLines), len(fulfillmentLineDataMap)); i++ {
		fulfillmentLine := fulfillmentLines[0]
		lineData := fulfillmentLineDataMap[fulfillmentLine.Id]
		quantity := int(lineData.Quantity)

		if fulfillmentLine.OrderLine.IsGiftcard {
			return nil, model.NewAppError(where, "app.order.cannot_refund_giftcard_line.app_error", nil, "cannot refund or return giftcard line", http.StatusNotAcceptable)
		}

		if fulfillmentLine.Quantity < quantity {
			return nil, model.NewAppError(where, model.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "Quantity"}, fmt.Sprintf("provided quantity: %d greater than quantity from fulfillment line: %d", quantity, fulfillmentLine.Quantity), http.StatusNotAcceptable)
		}

		if !lo.Contains(whitelistedStatuses, fulfillmentLine.Fulfillment.Status) {
			statusString := ""
			for _, status := range whitelistedStatuses {
				statusString += string(status) + ","
			}
			return nil, model.NewAppError(where, "app.order.fulfillment_status_not_acceptable.app_error", nil, statusString[:len(statusString)-1], http.StatusNotAcceptable)
		}

		res = append(res, &model.FulfillmentLineData{
			Line:     *fulfillmentLine,
			Quantity: quantity,
		})
	}

	return res, nil
}
