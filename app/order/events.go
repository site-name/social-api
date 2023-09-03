package order

import (
	"net/http"

	"github.com/samber/lo"
	"github.com/site-name/decimal"
	"github.com/sitename/sitename/model"
	"gorm.io/gorm"
)

// CommonCreateOrderEvent is common method for creating desired order event instance
func (a *ServiceOrder) CommonCreateOrderEvent(transaction *gorm.DB, option *model.OrderEventOption) (*model.OrderEvent, *model.AppError) {
	newOrderEvent := &model.OrderEvent{
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

func (s *ServiceOrder) LinePerQuantityToLineObject(quantity int, line *model.OrderLine) model.StringInterface {
	return model.StringInterface{
		"quantity": quantity,
		"line_pk":  line.Id,
		"item":     line.String(),
	}
}

func orderLinesToQuantityOrderLine(orderLines []*model.OrderLine) []*model.QuantityOrderLine {
	var res []*model.QuantityOrderLine
	for _, line := range orderLines {
		res = append(res, &model.QuantityOrderLine{
			Quantity:  line.Quantity,
			OrderLine: line,
		})
	}

	return res
}

func (s *ServiceOrder) LinesPerQuantityToLineObjectList(quantitiesPerOrderLine []*model.QuantityOrderLine) []model.StringInterface {
	return lo.Map(quantitiesPerOrderLine, func(item *model.QuantityOrderLine, _ int) model.StringInterface {
		return s.LinePerQuantityToLineObject(item.Quantity, item.OrderLine)
	})
}

func (s *ServiceOrder) PrepareDiscountObject(orderDiscount *model.OrderDiscount, oldOrderDiscount *model.OrderDiscount) model.StringInterface {
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

func (a *ServiceOrder) OrderDiscountsAutomaticallyUpdatedEvent(transaction *gorm.DB, ord *model.Order, changedOrderDiscounts [][2]*model.OrderDiscount) *model.AppError {
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

func (a *ServiceOrder) OrderDiscountAutomaticallyUpdatedEvent(transaction *gorm.DB, ord *model.Order, orderDiscount *model.OrderDiscount, oldOrderDiscount *model.OrderDiscount) (*model.OrderEvent, *model.AppError) {
	return a.OrderDiscountEvent(
		transaction,
		model.ORDER_EVENT_TYPE_ORDER_DISCOUNT_AUTOMATICALLY_UPDATED,
		ord,
		nil,
		orderDiscount,
		oldOrderDiscount,
	)
}

func (a *ServiceOrder) OrderDiscountEvent(transaction *gorm.DB, eventType model.OrderEventType, ord *model.Order, user *model.User, orderDiscount *model.OrderDiscount, oldOrderDiscount *model.OrderDiscount) (*model.OrderEvent, *model.AppError) {
	var userID *string
	if user == nil || !model.IsValidId(user.Id) {
		userID = nil
	} else {
		userID = model.NewPrimitive(user.Id)
	}

	discountParameters := a.PrepareDiscountObject(orderDiscount, oldOrderDiscount)

	return a.CommonCreateOrderEvent(transaction, &model.OrderEventOption{
		OrderID:    ord.Id,
		Type:       eventType,
		UserID:     userID,
		Parameters: discountParameters,
	})
}

func getPaymentData(amount *decimal.Decimal, payMent model.Payment) map[string]interface{} {
	return map[string]interface{}{
		"amount":          amount,
		"payment_id":      payMent.Token,
		"payment_gateway": payMent.GateWay,
	}
}

func (a *ServiceOrder) OrderLineDiscountEvent(eventType model.OrderEventType, ord *model.Order, user *model.User, line *model.OrderLine, lineBeforeUpdate *model.OrderLine) (*model.OrderEvent, *model.AppError) {
	var userID *string
	if user != nil || model.IsValidId(user.Id) {
		userID = &user.Id
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

	lineData := a.LinePerQuantityToLineObject(int(line.Quantity), line)
	lineData["discount"] = discountParameters

	return a.CommonCreateOrderEvent(nil, &model.OrderEventOption{
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

func (s *ServiceOrder) FulfillmentCanceledEvent(transaction *gorm.DB, orDer *model.Order, user *model.User, _ interface{}, fulfillment *model.Fulfillment) (*model.OrderEvent, *model.AppError) {
	var userID *string
	if user != nil {
		userID = &user.Id
	}

	params := model.StringInterface{}
	if fulfillment != nil {
		params["composed_id"] = fulfillment.ComposedId()
	}

	return s.CommonCreateOrderEvent(transaction, &model.OrderEventOption{
		OrderID:    orDer.Id,
		UserID:     userID,
		Type:       model.ORDER_EVENT_TYPE_FULFILLMENT_CANCELED,
		Parameters: params,
	})
}

func (s *ServiceOrder) FulfillmentFulfilledItemsEvent(transaction *gorm.DB, orDer *model.Order, user *model.User, _ interface{}, fulfillmentLines model.FulfillmentLines) (*model.OrderEvent, *model.AppError) {
	var userID *string
	if user != nil {
		userID = &user.Id
	}

	return s.CommonCreateOrderEvent(transaction, &model.OrderEventOption{
		UserID:  userID,
		OrderID: orDer.Id,
		Type:    model.ORDER_EVENT_TYPE_FULFILLMENT_FULFILLED_ITEMS,
		Parameters: model.StringInterface{
			"fulfilled_items": fulfillmentLines.IDs(),
		},
	})
}

func (s *ServiceOrder) OrderCreatedEvent(orDer model.Order, user *model.User, _ interface{}, fromDraft bool) (*model.OrderEvent, *model.AppError) {
	var (
		eventType = model.ORDER_EVENT_TYPE_PLACED_FROM_DRAFT
		userID    *string
	)
	if !fromDraft {
		eventType = model.ORDER_EVENT_TYPE_PLACED
		_, appErr := s.srv.AccountService().CustomerPlacedOrderEvent(user, orDer)
		if appErr != nil {
			return nil, appErr
		}
	}
	if user != nil && model.IsValidId(user.Id) {
		userID = &user.Id
	}

	return s.CommonCreateOrderEvent(nil, &model.OrderEventOption{
		OrderID: orDer.Id,
		UserID:  userID,
		Type:    eventType,
	})
}

func (s *ServiceOrder) OrderConfirmedEvent(orDer model.Order, user *model.User, _ interface{}) (*model.OrderEvent, *model.AppError) {
	var userID *string
	if user != nil && model.IsValidId(user.Id) {
		userID = &user.Id
	}

	return s.CommonCreateOrderEvent(nil, &model.OrderEventOption{
		UserID:  userID,
		OrderID: orDer.Id,
		Type:    model.ORDER_EVENT_TYPE_CONFIRMED,
	})
}

func (s *ServiceOrder) FulfillmentAwaitsApprovalEvent(transaction *gorm.DB, orDer *model.Order, user *model.User, _ interface{}, fulfillmentLines model.FulfillmentLines) (*model.OrderEvent, *model.AppError) {
	var userID *string
	if user != nil && model.IsValidId(user.Id) {
		userID = &user.Id
	}

	return s.CommonCreateOrderEvent(transaction, &model.OrderEventOption{
		OrderID: orDer.Id,
		UserID:  userID,
		Type:    model.ORDER_EVENT_TYPE_FULFILLMENT_AWAITS_APPROVAL,
		Parameters: model.StringInterface{
			"awaiting_fulfillments": fulfillmentLines.IDs(),
		},
	})
}

func (s *ServiceOrder) FulfillmentTrackingUpdatedEvent(orDer *model.Order, user *model.User, _ interface{}, trackingNumber string, fulfillment *model.Fulfillment) (*model.OrderEvent, *model.AppError) {
	var userID *string
	if user != nil && model.IsValidId(user.Id) {
		userID = &user.Id
	}

	return s.CommonCreateOrderEvent(nil, &model.OrderEventOption{
		OrderID: orDer.Id,
		UserID:  userID,
		Type:    model.ORDER_EVENT_TYPE_TRACKING_UPDATED,
		Parameters: model.StringInterface{
			"tracking_number": trackingNumber,
			"fulfillment":     fulfillment.ComposedId(),
		},
	})
}

func (s *ServiceOrder) OrderManuallyMarkedAsPaidEvent(transaction *gorm.DB, orDer model.Order, user *model.User, _ interface{}, transactionReference string) (*model.OrderEvent, *model.AppError) {
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

	return s.CommonCreateOrderEvent(transaction, &model.OrderEventOption{
		OrderID:    orDer.Id,
		UserID:     userID,
		Type:       model.ORDER_EVENT_TYPE_ORDER_MARKED_AS_PAID,
		Parameters: parameters,
	})
}

func (s *ServiceOrder) DraftOrderCreatedFromReplaceEvent(transaction *gorm.DB, draftOrder model.Order, originalOrder model.Order, user *model.User, _ interface{}, lines []*model.QuantityOrderLine) (*model.OrderEvent, *model.AppError) {
	var userID *string
	if user != nil && model.IsValidId(user.Id) {
		userID = &user.Id
	}

	return s.CommonCreateOrderEvent(transaction, &model.OrderEventOption{
		OrderID: draftOrder.Id,
		Type:    model.ORDER_EVENT_TYPE_DRAFT_CREATED_FROM_REPLACE,
		UserID:  userID,
		Parameters: model.StringInterface{
			"related_order_pk": originalOrder.Id,
			"lines":            s.LinesPerQuantityToLineObjectList(lines),
		},
	})
}

func (s *ServiceOrder) FulfillmentReplacedEvent(transaction *gorm.DB, orDer model.Order, user *model.User, _ interface{}, replacedLines []*model.QuantityOrderLine) (*model.OrderEvent, *model.AppError) {
	var userID *string

	if user != nil && model.IsValidId(user.Id) {
		userID = &user.Id
	}

	return s.CommonCreateOrderEvent(transaction, &model.OrderEventOption{
		OrderID: orDer.Id,
		UserID:  userID,
		Type:    model.ORDER_EVENT_TYPE_FULFILLMENT_REPLACED,
		Parameters: model.StringInterface{
			"lines": s.LinesPerQuantityToLineObjectList(replacedLines),
		},
	})
}

func (s *ServiceOrder) OrderReplacementCreated(transaction *gorm.DB, originalOrder model.Order, replaceOrder *model.Order, user *model.User, _ interface{}) (*model.OrderEvent, *model.AppError) {
	var userID *string

	if user != nil && model.IsValidId(user.Id) {
		userID = &user.Id
	}

	return s.CommonCreateOrderEvent(transaction, &model.OrderEventOption{
		OrderID: originalOrder.Id,
		UserID:  userID,
		Type:    model.ORDER_EVENT_TYPE_ORDER_REPLACEMENT_CREATED,
		Parameters: model.StringInterface{
			"related_order_pk": replaceOrder.Id,
		},
	})
}

func (s *ServiceOrder) FilterOrderEventsByOptions(options *model.OrderEventFilterOptions) ([]*model.OrderEvent, *model.AppError) {
	events, err := s.srv.Store.OrderEvent().FilterByOptions(options)
	if err != nil {
		return nil, model.NewAppError("FilterOrderEventsByOptions", "app.order.order_events_by_options.app_error", nil, err.Error(), http.StatusInternalServerError)
	}

	return events, nil
}
