package order

import (
	"net/http"

	"github.com/mattermost/gorp"
	"github.com/site-name/decimal"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model/account"
	"github.com/sitename/sitename/model/order"
	"github.com/sitename/sitename/model/payment"
	"github.com/sitename/sitename/model/product_and_discount"
)

// CommonCreateOrderEvent is common method for creating desired order event instance
func (a *ServiceOrder) CommonCreateOrderEvent(transaction *gorp.Transaction, option *order.OrderEventOption) (*order.OrderEvent, *model.AppError) {
	newOrderEvent := &order.OrderEvent{
		OrderID:    option.OrderID,
		Type:       option.Type,
		Parameters: option.Parameters,
		UserID:     option.UserID,
	}

	orderEvent, err := a.srv.Store.OrderEvent().Save(transaction, newOrderEvent)
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

func orderLinesToQuantityOrderLine(orderLines []*order.OrderLine) []*order.QuantityOrderLine {
	var res []*order.QuantityOrderLine
	for _, line := range orderLines {
		res = append(res, &order.QuantityOrderLine{
			Quantity:  line.Quantity,
			OrderLine: line,
		})
	}

	return res
}

func linesPerQuantityToLineObjectList(quantitiesPerOrderLine []*order.QuantityOrderLine) []map[string]interface{} {
	res := []map[string]interface{}{}
	for _, item := range quantitiesPerOrderLine {
		res = append(res, linePerQuantityToLineObject(item.Quantity, item.OrderLine))
	}

	return res
}

func prepareDiscountObject(orderDiscount *product_and_discount.OrderDiscount, oldOrderDiscount *product_and_discount.OrderDiscount) model.StringInterface {
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

	return discountParameters
}

func (a *ServiceOrder) OrderDiscountsAutomaticallyUpdatedEvent(transaction *gorp.Transaction, ord *order.Order, changedOrderDiscounts [][2]*product_and_discount.OrderDiscount) *model.AppError {
	for _, tuple := range changedOrderDiscounts {
		_, appErr := a.OrderDiscountAutomaticallyUpdatedEvent(
			transaction,
			ord,
			tuple[1],
			tuple[0],
		)
		if appErr != nil {
			return appErr
		}
	}

	return nil
}

func (a *ServiceOrder) OrderDiscountAutomaticallyUpdatedEvent(transaction *gorp.Transaction, ord *order.Order, orderDiscount *product_and_discount.OrderDiscount, oldOrderDiscount *product_and_discount.OrderDiscount) (*order.OrderEvent, *model.AppError) {
	return a.OrderDiscountEvent(
		transaction,
		order.ORDER_DISCOUNT_AUTOMATICALLY_UPDATED,
		ord,
		nil,
		orderDiscount,
		oldOrderDiscount,
	)
}

func (a *ServiceOrder) OrderDiscountEvent(transaction *gorp.Transaction, eventType order.OrderEvents, ord *order.Order, user *account.User, orderDiscount *product_and_discount.OrderDiscount, oldOrderDiscount *product_and_discount.OrderDiscount) (*order.OrderEvent, *model.AppError) {
	var userID *string
	if user == nil || !model.IsValidId(user.Id) {
		userID = nil
	} else {
		userID = model.NewString(user.Id)
	}

	discountParameters := prepareDiscountObject(orderDiscount, oldOrderDiscount)

	return a.CommonCreateOrderEvent(transaction, &order.OrderEventOption{
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

func (a *ServiceOrder) OrderLineDiscountEvent(eventType order.OrderEvents, ord *order.Order, user *account.User, line *order.OrderLine, lineBeforeUpdate *order.OrderLine) (*order.OrderEvent, *model.AppError) {
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

	return a.CommonCreateOrderEvent(nil, &order.OrderEventOption{
		OrderID: ord.Id,
		Type:    eventType,
		UserID:  userID,
		Parameters: model.StringInterface{
			"lines": []map[string]interface{}{
				lineData,
			},
		},
	})
}

func (s *ServiceOrder) FulfillmentCanceledEvent(transaction *gorp.Transaction, orDer *order.Order, user *account.User, _ interface{}, fulfillment *order.Fulfillment) (*order.OrderEvent, *model.AppError) {
	var userID *string
	if user != nil {
		userID = &user.Id
	}

	params := model.StringInterface{}
	if fulfillment != nil {
		params["composed_id"] = fulfillment.ComposedId()
	}

	return s.CommonCreateOrderEvent(transaction, &order.OrderEventOption{
		OrderID:    orDer.Id,
		UserID:     userID,
		Type:       order.FULFILLMENT_CANCELED_,
		Parameters: params,
	})
}

func (s *ServiceOrder) FulfillmentFulfilledItemsEvent(transaction *gorp.Transaction, orDer *order.Order, user *account.User, _ interface{}, fulfillmentLines order.FulfillmentLines) (*order.OrderEvent, *model.AppError) {
	var userID *string
	if user != nil {
		userID = &user.Id
	}

	return s.CommonCreateOrderEvent(transaction, &order.OrderEventOption{
		UserID:  userID,
		OrderID: orDer.Id,
		Type:    order.FULFILLMENT_FULFILLED_ITEMS,
		Parameters: model.StringInterface{
			"fulfilled_items": fulfillmentLines.IDs(),
		},
	})
}

func (s *ServiceOrder) OrderCreatedEvent(orDer order.Order, user *account.User, _ interface{}, fromDraft bool) (*order.OrderEvent, *model.AppError) {
	var (
		eventType = order.PLACED_FROM_DRAFT
		userID    *string
	)
	if !fromDraft {
		eventType = order.PLACED
		_, appErr := s.srv.AccountService().CustomerPlacedOrderEvent(user, orDer)
		if appErr != nil {
			return nil, appErr
		}
	}
	if user != nil && model.IsValidId(user.Id) {
		userID = &user.Id
	}

	return s.CommonCreateOrderEvent(nil, &order.OrderEventOption{
		OrderID: orDer.Id,
		UserID:  userID,
		Type:    eventType,
	})
}

func (s *ServiceOrder) OrderConfirmedEvent(orDer order.Order, user *account.User, _ interface{}) (*order.OrderEvent, *model.AppError) {
	var userID *string
	if user != nil && model.IsValidId(user.Id) {
		userID = &user.Id
	}

	return s.CommonCreateOrderEvent(nil, &order.OrderEventOption{
		UserID:  userID,
		OrderID: orDer.Id,
		Type:    order.CONFIRMED,
	})
}

func (s *ServiceOrder) FulfillmentAwaitsApprovalEvent(transaction *gorp.Transaction, orDer *order.Order, user *account.User, _ interface{}, fulfillmentLines order.FulfillmentLines) (*order.OrderEvent, *model.AppError) {
	var userID *string
	if user != nil && model.IsValidId(user.Id) {
		userID = &user.Id
	}

	return s.CommonCreateOrderEvent(transaction, &order.OrderEventOption{
		OrderID: orDer.Id,
		UserID:  userID,
		Type:    order.FULFILLMENT_AWAITS_APPROVAL,
		Parameters: model.StringInterface{
			"awaiting_fulfillments": fulfillmentLines.IDs(),
		},
	})
}

func (s *ServiceOrder) FulfillmentTrackingUpdatedEvent(orDer *order.Order, user *account.User, _ interface{}, trackingNumber string, fulfillment *order.Fulfillment) (*order.OrderEvent, *model.AppError) {
	var userID *string
	if user != nil && model.IsValidId(user.Id) {
		userID = &user.Id
	}

	return s.CommonCreateOrderEvent(nil, &order.OrderEventOption{
		OrderID: orDer.Id,
		UserID:  userID,
		Type:    order.TRACKING_UPDATED,
		Parameters: model.StringInterface{
			"tracking_number": trackingNumber,
			"fulfillment":     fulfillment.ComposedId(),
		},
	})
}

func (s *ServiceOrder) OrderManuallyMarkedAsPaidEvent(transaction *gorp.Transaction, orDer *order.Order, user *account.User, _ interface{}, transactionReference string) (*order.OrderEvent, *model.AppError) {
	var (
		userID     *string
		parameters = model.StringInterface{}
	)
	if user != nil && model.IsValidId(user.Id) {
		userID = &user.Id
	}
	if transactionReference != "" {
		parameters["transaction_reference"] = transactionReference
	}

	return s.CommonCreateOrderEvent(transaction, &order.OrderEventOption{
		OrderID:    orDer.Id,
		UserID:     userID,
		Type:       order.ORDER_MARKED_AS_PAID,
		Parameters: parameters,
	})
}

func (s *ServiceOrder) DraftOrderCreatedFromReplaceEvent(transaction *gorp.Transaction, draftOrder *order.Order, originalOrder *order.Order, user *account.User, _ interface{}, lines []*order.QuantityOrderLine) (*order.OrderEvent, *model.AppError) {
	var userID *string
	if user != nil && model.IsValidId(user.Id) {
		userID = &user.Id
	}

	return s.CommonCreateOrderEvent(transaction, &order.OrderEventOption{
		OrderID: draftOrder.Id,
		Type:    order.DRAFT_CREATED_FROM_REPLACE,
		UserID:  userID,
		Parameters: model.StringInterface{
			"related_order_pk": originalOrder.Id,
			"lines":            linesPerQuantityToLineObjectList(lines),
		},
	})
}

func (s *ServiceOrder) FulfillmentReplacedEvent(transaction *gorp.Transaction, orDer *order.Order, user *account.User, _ interface{}, replacedLines []*order.QuantityOrderLine) (*order.OrderEvent, *model.AppError) {
	var userID *string

	if user != nil && model.IsValidId(user.Id) {
		userID = &user.Id
	}

	return s.CommonCreateOrderEvent(transaction, &order.OrderEventOption{
		OrderID: orDer.Id,
		UserID:  userID,
		Type:    order.FULFILLMENT_REPLACED_,
		Parameters: model.StringInterface{
			"lines": linesPerQuantityToLineObjectList(replacedLines),
		},
	})
}

func (s *ServiceOrder) OrderReplacementCreated(transaction *gorp.Transaction, originalOrder *order.Order, replaceOrder *order.Order, user *account.User, _ interface{}) (*order.OrderEvent, *model.AppError) {
	var userID *string

	if user != nil && model.IsValidId(user.Id) {
		userID = &user.Id
	}

	return s.CommonCreateOrderEvent(transaction, &order.OrderEventOption{
		OrderID: originalOrder.Id,
		UserID:  userID,
		Type:    order.ORDER_REPLACEMENT_CREATED,
		Parameters: model.StringInterface{
			"related_order_pk": replaceOrder.Id,
		},
	})
}
