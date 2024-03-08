package order

import (
	"net/http"

	"github.com/samber/lo"
	"github.com/site-name/decimal"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model_helper"
	"github.com/sitename/sitename/modules/model_types"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"gorm.io/gorm"
)

// CommonCreateOrderEvent is common method for creating desired order event instance
func (a *ServiceOrder) CommonCreateOrderEvent(transaction boil.ContextTransactor, option *model.OrderEventOption) (*model.OrderEvent, *model_helper.AppError) {
	newOrderEvent := &model.OrderEvent{
		OrderID:    option.OrderID,
		Type:       option.Type,
		Parameters: option.Parameters,
		UserID:     option.UserID,
	}

	orderEvent, err := a.srv.Store.OrderEvent().Save(transaction, newOrderEvent)
	if err != nil {
		if appErr, ok := err.(*model_helper.AppError); ok {
			return nil, appErr
		}
		return nil, model_helper.NewAppError("CommonCreateOrderEvent", "app.order.error_creating_order_event.app_error", nil, err.Error(), http.StatusInternalServerError)
	}

	return orderEvent, nil
}

func (s *ServiceOrder) LinePerQuantityToLineObject(quantity int, line *model.OrderLine) model_types.JSONString {
	return model_types.JSONString{
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

func (s *ServiceOrder) LinesPerQuantityToLineObjectList(quantitiesPerOrderLine []*model.QuantityOrderLine) []model_types.JSONString {
	return lo.Map(quantitiesPerOrderLine, func(item *model.QuantityOrderLine, _ int) model_types.JSONString {
		return s.LinePerQuantityToLineObject(item.Quantity, item.OrderLine)
	})
}

func (s *ServiceOrder) PrepareDiscountObject(orderDiscount *model.OrderDiscount, oldOrderDiscount *model.OrderDiscount) model_types.JSONString {
	discountParameters := model_types.JSONString{
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

func (a *ServiceOrder) OrderDiscountsAutomaticallyUpdatedEvent(transaction boil.ContextTransactor, ord *model.Order, changedOrderDiscounts [][2]*model.OrderDiscount) *model_helper.AppError {
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

func (a *ServiceOrder) OrderDiscountAutomaticallyUpdatedEvent(transaction boil.ContextTransactor, ord *model.Order, orderDiscount *model.OrderDiscount, oldOrderDiscount *model.OrderDiscount) (*model.OrderEvent, *model_helper.AppError) {
	return a.OrderDiscountEvent(
		transaction,
		model.ORDER_EVENT_TYPE_ORDER_DISCOUNT_AUTOMATICALLY_UPDATED,
		ord,
		nil,
		orderDiscount,
		oldOrderDiscount,
	)
}

func (a *ServiceOrder) OrderDiscountEvent(transaction boil.ContextTransactor, eventType model.OrderEventType, ord *model.Order, user *model.User, orderDiscount *model.OrderDiscount, oldOrderDiscount *model.OrderDiscount) (*model.OrderEvent, *model_helper.AppError) {
	var userID *string
	if user == nil || !model_helper.IsValidId(user.Id) {
		userID = nil
	} else {
		userID = model_helper.GetPointerOfValue(user.Id)
	}

	discountParameters := a.PrepareDiscountObject(orderDiscount, oldOrderDiscount)

	return a.CommonCreateOrderEvent(transaction, &model.OrderEventOption{
		OrderID:    ord.Id,
		Type:       eventType,
		UserID:     userID,
		Parameters: discountParameters,
	})
}

func getPaymentData(amount *decimal.Decimal, payMent model.Payment) map[string]any {
	return map[string]any{
		"amount":          amount,
		"payment_id":      payMent.Token,
		"payment_gateway": payMent.GateWay,
	}
}

func (a *ServiceOrder) OrderLineDiscountEvent(eventType model.OrderEventType, ord *model.Order, user *model.User, line *model.OrderLine, lineBeforeUpdate *model.OrderLine) (*model.OrderEvent, *model_helper.AppError) {
	var userID *string
	if user != nil || model_helper.IsValidId(user.Id) {
		userID = &user.Id
	}
	discountParameters := map[string]any{
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
		Parameters: model_types.JSONString{
			"lines": []map[string]any{
				lineData,
			},
		},
	})
}

func (s *ServiceOrder) FulfillmentCanceledEvent(transaction boil.ContextTransactor, orDer *model.Order, user *model.User, _ any, fulfillment *model.Fulfillment) (*model.OrderEvent, *model_helper.AppError) {
	var userID *string
	if user != nil {
		userID = &user.Id
	}

	params := model_types.JSONString{}
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

func (s *ServiceOrder) FulfillmentFulfilledItemsEvent(transaction boil.ContextTransactor, orDer *model.Order, user *model.User, _ any, fulfillmentLines model.FulfillmentLines) (*model.OrderEvent, *model_helper.AppError) {
	var userID *string
	if user != nil {
		userID = &user.Id
	}

	return s.CommonCreateOrderEvent(transaction, &model.OrderEventOption{
		UserID:  userID,
		OrderID: orDer.Id,
		Type:    model.ORDER_EVENT_TYPE_FULFILLMENT_FULFILLED_ITEMS,
		Parameters: model_types.JSONString{
			"fulfilled_items": fulfillmentLines.IDs(),
		},
	})
}

func (s *ServiceOrder) OrderCreatedEvent(orDer model.Order, user *model.User, _ any, fromDraft bool) (*model.OrderEvent, *model_helper.AppError) {
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
	if user != nil && model_helper.IsValidId(user.Id) {
		userID = &user.Id
	}

	return s.CommonCreateOrderEvent(nil, &model.OrderEventOption{
		OrderID: orDer.Id,
		UserID:  userID,
		Type:    eventType,
	})
}

func (s *ServiceOrder) OrderConfirmedEvent(tx *gorm.DB, orDer model.Order, user *model.User, _ any) (*model.OrderEvent, *model_helper.AppError) {
	var userID *string
	if user != nil && model_helper.IsValidId(user.Id) {
		userID = &user.Id
	}

	return s.CommonCreateOrderEvent(tx, &model.OrderEventOption{
		UserID:  userID,
		OrderID: orDer.Id,
		Type:    model.ORDER_EVENT_TYPE_CONFIRMED,
	})
}

func (s *ServiceOrder) FulfillmentAwaitsApprovalEvent(transaction boil.ContextTransactor, orDer *model.Order, user *model.User, _ any, fulfillmentLines model.FulfillmentLines) (*model.OrderEvent, *model_helper.AppError) {
	var userID *string
	if user != nil && model_helper.IsValidId(user.Id) {
		userID = &user.Id
	}

	return s.CommonCreateOrderEvent(transaction, &model.OrderEventOption{
		OrderID: orDer.Id,
		UserID:  userID,
		Type:    model.ORDER_EVENT_TYPE_FULFILLMENT_AWAITS_APPROVAL,
		Parameters: model_types.JSONString{
			"awaiting_fulfillments": fulfillmentLines.IDs(),
		},
	})
}

func (s *ServiceOrder) FulfillmentTrackingUpdatedEvent(orDer *model.Order, user *model.User, _ any, trackingNumber string, fulfillment *model.Fulfillment) (*model.OrderEvent, *model_helper.AppError) {
	var userID *string
	if user != nil && model_helper.IsValidId(user.Id) {
		userID = &user.Id
	}

	return s.CommonCreateOrderEvent(nil, &model.OrderEventOption{
		OrderID: orDer.Id,
		UserID:  userID,
		Type:    model.ORDER_EVENT_TYPE_TRACKING_UPDATED,
		Parameters: model_types.JSONString{
			"tracking_number": trackingNumber,
			"fulfillment":     fulfillment.ComposedId(),
		},
	})
}

func (s *ServiceOrder) OrderManuallyMarkedAsPaidEvent(transaction boil.ContextTransactor, orDer model.Order, user *model.User, _ any, transactionReference string) (*model.OrderEvent, *model_helper.AppError) {
	var (
		userID     *string
		parameters = model_types.JSONString{}
	)
	if user != nil && model_helper.IsValidId(user.Id) {
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

func (s *ServiceOrder) DraftOrderCreatedFromReplaceEvent(transaction boil.ContextTransactor, draftOrder model.Order, originalOrder model.Order, user *model.User, _ any, lines []*model.QuantityOrderLine) (*model.OrderEvent, *model_helper.AppError) {
	var userID *string
	if user != nil && model_helper.IsValidId(user.Id) {
		userID = &user.Id
	}

	return s.CommonCreateOrderEvent(transaction, &model.OrderEventOption{
		OrderID: draftOrder.Id,
		Type:    model.ORDER_EVENT_TYPE_DRAFT_CREATED_FROM_REPLACE,
		UserID:  userID,
		Parameters: model_types.JSONString{
			"related_order_pk": originalOrder.Id,
			"lines":            s.LinesPerQuantityToLineObjectList(lines),
		},
	})
}

func (s *ServiceOrder) FulfillmentReplacedEvent(transaction boil.ContextTransactor, orDer model.Order, user *model.User, _ any, replacedLines []*model.QuantityOrderLine) (*model.OrderEvent, *model_helper.AppError) {
	var userID *string

	if user != nil && model_helper.IsValidId(user.Id) {
		userID = &user.Id
	}

	return s.CommonCreateOrderEvent(transaction, &model.OrderEventOption{
		OrderID: orDer.Id,
		UserID:  userID,
		Type:    model.ORDER_EVENT_TYPE_FULFILLMENT_REPLACED,
		Parameters: model_types.JSONString{
			"lines": s.LinesPerQuantityToLineObjectList(replacedLines),
		},
	})
}

func (s *ServiceOrder) OrderReplacementCreated(transaction boil.ContextTransactor, originalOrder model.Order, replaceOrder *model.Order, user *model.User, _ any) (*model.OrderEvent, *model_helper.AppError) {
	var userID *string

	if user != nil && model_helper.IsValidId(user.Id) {
		userID = &user.Id
	}

	return s.CommonCreateOrderEvent(transaction, &model.OrderEventOption{
		OrderID: originalOrder.Id,
		UserID:  userID,
		Type:    model.ORDER_EVENT_TYPE_ORDER_REPLACEMENT_CREATED,
		Parameters: model_types.JSONString{
			"related_order_pk": replaceOrder.Id,
		},
	})
}

func (s *ServiceOrder) FilterOrderEventsByOptions(options *model.OrderEventFilterOptions) ([]*model.OrderEvent, *model_helper.AppError) {
	events, err := s.srv.Store.OrderEvent().FilterByOptions(options)
	if err != nil {
		return nil, model_helper.NewAppError("FilterOrderEventsByOptions", "app.order.order_events_by_options.app_error", nil, err.Error(), http.StatusInternalServerError)
	}

	return events, nil
}

func (s *ServiceOrder) OrderNoteAddedEvent(tx *gorm.DB, order *model.Order, user *model.User, message string) (*model.OrderEvent, *model_helper.AppError) {
	if user != nil && model_helper.IsValidId(user.Id) {
		if order.UserID != nil && *order.UserID == user.Id {
			_, appErr := s.srv.AccountService().CommonCustomerCreateEvent(tx, &user.Id, &order.Id, model.CUSTOMER_EVENT_TYPE_NOTE_ADDED_TO_ORDER, model_types.JSONString{"message": message})
			if appErr != nil {
				return nil, appErr
			}
		}
	}

	return s.CommonCreateOrderEvent(tx, &model.OrderEventOption{
		OrderID: order.Id,
		UserID:  &user.Id,
		Type:    model.ORDER_EVENT_TYPE_NOTE_ADDED,
		Parameters: model_types.JSONString{
			"message": message,
		},
	})
}
