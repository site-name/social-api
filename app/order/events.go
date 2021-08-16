package order

import (
	"net/http"

	"github.com/site-name/decimal"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model/account"
	"github.com/sitename/sitename/model/order"
	"github.com/sitename/sitename/model/payment"
	"github.com/sitename/sitename/model/product_and_discount"
)

// CommonCreateOrderEvent is common method for creating desired order event instance
func (a *AppOrder) CommonCreateOrderEvent(option *order.OrderEventOption) (*order.OrderEvent, *model.AppError) {
	newOrderEvent := &order.OrderEvent{
		OrderID:    option.OrderID,
		Type:       option.Type,
		Parameters: option.Parameters,
		UserID:     option.UserID,
	}

	orderEvent, err := a.Srv().Store.OrderEvent().Save(newOrderEvent)
	if err != nil {
		if appErr, ok := err.(*model.AppError); ok {
			return nil, appErr
		}
		return nil, model.NewAppError("CommonCreateOrderEvent", "app.order.error_creating_order_event.app_error", nil, err.Error(), http.StatusInternalServerError)
	}

	return orderEvent, nil
}

func linePerQuantityToLineObject(quantity int, line *order.OrderLine) map[string]interface{} {
	return map[string]interface{}{
		"quantity": quantity,
		"line_pk":  line.Id,
		"item":     line.String(),
	}
}

type QuantityOrderLine struct {
	Quantity  int
	OrderLine *order.OrderLine
}

func linesPerQuantityToLineObjectList(quantitiesPerOrderLine []*QuantityOrderLine) []map[string]interface{} {
	res := []map[string]interface{}{}
	for _, item := range quantitiesPerOrderLine {
		res = append(res, linePerQuantityToLineObject(item.Quantity, item.OrderLine))
	}

	return res
}

func prepareDiscountObject(orderDiscount *product_and_discount.OrderDiscount, oldOrderDiscount *product_and_discount.OrderDiscount) *model.StringInterface {
	discountParameters := model.StringInterface{
		"value":        orderDiscount.Value,
		"amount_value": orderDiscount.AmountValue,
		"currency":     orderDiscount.Currency,
		"value_type":   orderDiscount.ValueType,
		"reason":       orderDiscount.Reason,
	}
	if oldOrderDiscount != nil {
		discountParameters["old_value"] = oldOrderDiscount.Value
		discountParameters["old_value_type"] = oldOrderDiscount.ValueType
		discountParameters["old_amount_value"] = oldOrderDiscount.AmountValue
	}

	return &discountParameters
}

func (a *AppOrder) OrderDiscountsAutomaticallyUpdatedEvent(ord *order.Order, changedOrderDiscounts [][2]*product_and_discount.OrderDiscount) *model.AppError {
	for _, tuple := range changedOrderDiscounts {
		_, appErr := a.OrderDiscountAutomaticallyUpdatedEvent(
			ord,
			tuple[1],
			tuple[0],
		)
		if appErr != nil {
			appErr.Where = "OrderDiscountsAutomaticallyUpdatedEvent"
			return appErr
		}
	}

	return nil
}

func (a *AppOrder) OrderDiscountAutomaticallyUpdatedEvent(ord *order.Order, orderDiscount *product_and_discount.OrderDiscount, oldOrderDiscount *product_and_discount.OrderDiscount) (*order.OrderEvent, *model.AppError) {
	return a.OrderDiscountEvent(
		order.ORDER_EVENT_TYPE__ORDER_DISCOUNT_AUTOMATICALLY_UPDATED,
		ord,
		nil,
		orderDiscount,
		oldOrderDiscount,
	)
}

func (a *AppOrder) OrderDiscountEvent(eventType string, ord *order.Order, user *account.User, orderDiscount *product_and_discount.OrderDiscount, oldOrderDiscount *product_and_discount.OrderDiscount) (*order.OrderEvent, *model.AppError) {
	var userID *string
	if user == nil || !model.IsValidId(user.Id) {
		userID = nil
	} else {
		userID = model.NewString(user.Id)
	}

	discountParameters := prepareDiscountObject(orderDiscount, oldOrderDiscount)

	return a.CommonCreateOrderEvent(&order.OrderEventOption{
		OrderID:    ord.Id,
		Type:       eventType,
		UserID:     userID,
		Parameters: discountParameters,
	})
}

func getPaymentData(amount *decimal.Decimal, payMent *payment.Payment) map[string]map[string]interface{} {
	return map[string]map[string]interface{}{
		"parameters": {
			"amount":          amount,
			"payment_id":      payMent.Token,
			"payment_gateway": payMent.GateWay,
		},
	}
}

func (a *AppOrder) OrderLineDiscountEvent(eventType string, ord *order.Order, user *account.User, line *order.OrderLine, lineBeforeUpdate *order.OrderLine) (*order.OrderEvent, *model.AppError) {
	var userID *string
	if user == nil || !model.IsValidId(user.Id) {
		userID = nil
	} else {
		userID = model.NewString(user.Id)
	}

	discountParameters := map[string]interface{}{
		"value":        line.UnitDiscountValue,
		"amount_value": line.UnitDiscountAmount,
		"currency":     line.Currency,
		"value_type":   line.UnitDiscountType,
		"reason":       line.UnitDiscountReason,
	}
	if lineBeforeUpdate != nil {
		discountParameters["old_value"] = lineBeforeUpdate.UnitDiscountValue
		discountParameters["old_value_type"] = lineBeforeUpdate.UnitDiscountType
		discountParameters["old_amount_value"] = lineBeforeUpdate.UnitDiscountAmount
	}

	lineData := linePerQuantityToLineObject(int(line.Quantity), line)
	lineData["discount"] = discountParameters

	return a.CommonCreateOrderEvent(&order.OrderEventOption{
		OrderID: ord.Id,
		Type:    eventType,
		UserID:  userID,
		Parameters: &model.StringInterface{
			"lines": []map[string]interface{}{
				lineData,
			},
		},
	})
}
