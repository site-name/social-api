package order

import (
	"context"
	"net/http"
	"strings"

	"github.com/mattermost/squirrel"
	"github.com/samber/lo"
	"github.com/site-name/decimal"
	"github.com/sitename/sitename/app/plugin/interfaces"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model_helper"
	"github.com/sitename/sitename/modules/model_types"
	"github.com/sitename/sitename/modules/slog"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"gorm.io/gorm"
)

// OrderCreated. `fromDraft` is default to false
func (a *ServiceOrder) OrderCreated(tx *gorm.DB, order model.Order, user *model.User, _ any, manager interfaces.PluginManagerInterface, fromDraft bool) (*model_helper.InsufficientStock, *model_helper.AppError) {
	// create order created event
	_, appErr := a.OrderCreatedEvent(order, user, nil, fromDraft)
	if appErr != nil {
		return nil, appErr
	}

	_, appErr = manager.OrderCreated(order)
	if appErr != nil {
		return nil, appErr
	}

	lastPaymentOfOrder, appErr := a.srv.Payment.GetLastOrderPayment(order.Id)
	if appErr != nil {
		if appErr.StatusCode == http.StatusInternalServerError {
			return nil, appErr
		}
	}
	if lastPaymentOfOrder != nil {
		orderIsCaptured, appErr := a.OrderIsCaptured(order.Id)
		if appErr != nil {
			return nil, appErr
		}

		if orderIsCaptured {
			InsufficientStock, appErr := a.OrderCaptured(order, user, nil, lastPaymentOfOrder.Total, *lastPaymentOfOrder, manager)
			if InsufficientStock != nil || appErr != nil {
				return InsufficientStock, appErr
			}

			goto finally
		}

		orderIsPreAuthorized, appErr := a.OrderIsPreAuthorized(order.Id)
		if appErr != nil {
			return nil, appErr
		}
		if orderIsPreAuthorized {
			appErr = a.OrderAuthorized(order, user, nil, lastPaymentOfOrder.Total, *lastPaymentOfOrder, manager)
			if appErr != nil {
				return nil, appErr
			}
		}
	}

finally:
	if *a.srv.Config().ShopSettings.AutomaticallyConfirmAllNewOrders {
		appErr = a.OrderConfirmed(tx, order, user, nil, manager, false)
		if appErr != nil {
			return nil, appErr
		}
	}

	return nil, nil
}

// OrderConfirmed Trigger event, plugin hooks and optionally confirmation email.
func (a *ServiceOrder) OrderConfirmed(tx *gorm.DB, order model.Order, user *model.User, _ any, manager interfaces.PluginManagerInterface, sendConfirmationEmail bool) *model_helper.AppError {
	_, appErr := a.OrderConfirmedEvent(tx, order, user, nil)
	if appErr != nil {
		return appErr
	}

	_, appErr = manager.OrderConfirmed(order)
	if appErr != nil {
		return appErr
	}

	if sendConfirmationEmail {
		a.SendOrderConfirmed(order, user, nil, manager)
	}

	return nil
}

// HandleFullyPaidOrder
//
// user can be nil
func (a *ServiceOrder) HandleFullyPaidOrder(manager interfaces.PluginManagerInterface, order model.Order, user *model.User, _ any) (*model_helper.InsufficientStock, *model_helper.AppError) {
	var userID *string
	if user != nil {
		userID = &user.Id
	}

	_, appErr := a.CommonCreateOrderEvent(nil, &model.OrderEventOption{
		OrderID: order.Id,
		Type:    model.ORDER_EVENT_TYPE_ORDER_FULLY_PAID,
		UserID:  userID,
	})
	if appErr != nil {
		return nil, appErr
	}

	customerEmail, appErr := a.srv.OrderService().CustomerEmail(&order)
	if appErr != nil {
		return nil, appErr
	}

	if model.IsValidEmail(customerEmail) {
		appErr = a.SendPaymentConfirmation(order, manager)
		if appErr != nil {
			return nil, appErr
		}

		orderNeedsAutoFulfillment, appErr := a.OrderNeedsAutomaticFulfillment(order)
		if appErr != nil {
			return nil, appErr
		}
		if orderNeedsAutoFulfillment {
			insufficientStock, appErr := a.AutomaticallyFulfillDigitalLines(order, manager)
			if insufficientStock != nil || appErr != nil {
				return insufficientStock, appErr
			}
		}
	}

	// TODO: implement me
	// panic("not implemented")

	_, appErr = manager.OrderFullyPaid(order)
	if appErr != nil {
		return nil, appErr
	}
	_, appErr = manager.OrderUpdated(order)
	return nil, appErr
}

// CancelOrder Release allocation of unfulfilled order items.
func (a *ServiceOrder) CancelOrder(transaction boil.ContextTransactor, order *model.Order, user *model.User, _ any, manager interfaces.PluginManagerInterface) *model_helper.AppError {
	// determine user id
	var userID *string
	if user != nil {
		userID = &user.Id
	}

	_, appErr := a.CommonCreateOrderEvent(transaction, &model.OrderEventOption{
		OrderID: order.Id,
		UserID:  userID,
		Type:    model.ORDER_EVENT_TYPE_CANCELED,
	})
	if appErr != nil {
		return appErr
	}

	appErr = a.srv.WarehouseService().DeAllocateStockForOrder(order, manager)
	if appErr != nil {
		return appErr
	}

	order.Status = model.ORDER_STATUS_CANCELED
	_, appErr = a.UpsertOrder(transaction, order)
	if appErr != nil {
		return appErr
	}

	appErr = a.SendOrderCancelledConfirmation(order, user, nil, manager)
	if appErr != nil {
		return appErr
	}

	_, appErr = manager.OrderCancelled(*order)
	if appErr != nil {
		return appErr
	}

	_, appErr = manager.OrderUpdated(*order)

	return appErr
}

// OrderRefunded
func (a *ServiceOrder) OrderRefunded(order model.Order, user *model.User, _ any, amount decimal.Decimal, payMent model.Payment, manager interfaces.PluginManagerInterface) *model_helper.AppError {
	var userID *string
	if user != nil && model_helper.IsValidId(user.Id) {
		userID = &user.Id
	}
	_, appErr := a.CommonCreateOrderEvent(nil, &model.OrderEventOption{
		OrderID:    order.Id,
		Type:       model.ORDER_EVENT_TYPE_PAYMENT_REFUNDED,
		UserID:     userID,
		Parameters: getPaymentData(&amount, payMent),
	})
	if appErr != nil {
		return appErr
	}

	_, appErr = manager.OrderUpdated(order)
	if appErr != nil {
		return appErr
	}

	return a.SendOrderRefundedConfirmation(order, user, nil, amount, payMent.Currency, manager)
}

// OrderVoided
func (a *ServiceOrder) OrderVoided(order model.Order, user *model.User, _ any, payMent *model.Payment, manager interfaces.PluginManagerInterface) *model_helper.AppError {
	var userID *string
	if user != nil && model_helper.IsValidId(user.Id) {
		userID = &user.Id
	}

	_, appErr := a.CommonCreateOrderEvent(nil, &model.OrderEventOption{
		OrderID:    order.Id,
		UserID:     userID,
		Type:       model.ORDER_EVENT_TYPE_PAYMENT_VOIDED,
		Parameters: getPaymentData(nil, *payMent),
	})
	if appErr != nil {
		return appErr
	}

	_, appErr = manager.OrderUpdated(order)
	return appErr
}

// OrderReturned
func (a *ServiceOrder) OrderReturned(transaction boil.ContextTransactor, order model.Order, user *model.User, _ any, returnedLines []*model.QuantityOrderLine) *model_helper.AppError {
	var userID *string
	if user != nil && model_helper.IsValidId(user.Id) {
		userID = &user.Id
	}

	_, appErr := a.CommonCreateOrderEvent(transaction, &model.OrderEventOption{
		OrderID: order.Id,
		Type:    model.ORDER_EVENT_TYPE_FULFILLMENT_RETURNED,
		UserID:  userID,
		Parameters: model_types.JSONString{
			"lines": a.LinesPerQuantityToLineObjectList(returnedLines),
		},
	})
	if appErr != nil {
		return appErr
	}

	appErr = a.UpdateOrderStatus(transaction, order)
	if appErr != nil {
		return appErr
	}

	return nil
}

// OrderFulfilled
func (a *ServiceOrder) OrderFulfilled(fulfillments []*model.Fulfillment, user *model.User, _ any, fulfillmentLines []*model.FulfillmentLine, manager interfaces.PluginManagerInterface, notifyCustomer bool) *model_helper.AppError {
	transaction := a.srv.Store.GetMaster().Begin()
	if transaction.Error != nil {
		return model_helper.NewAppError("OrderFulfilled", model_helper.ErrorCreatingTransactionErrorID, nil, transaction.Error.Error(), http.StatusInternalServerError)
	}
	defer a.srv.Store.FinalizeTransaction(transaction)

	var order = fulfillments[0].GetOrder()
	if order == nil {
		order_, appErr := a.OrderById(fulfillments[0].OrderID)
		if appErr != nil {
			return appErr
		}
		order = order_
	}

	appErr := a.UpdateOrderStatus(transaction, *order)
	if appErr != nil {
		return appErr
	}
	_, appErr = a.FulfillmentFulfilledItemsEvent(transaction, order, user, nil, fulfillmentLines)
	if appErr != nil {
		return appErr
	}

	if err := transaction.Commit().Error; err != nil {
		return model_helper.NewAppError("CancelOrder", model_helper.ErrorCommittingTransactionErrorID, nil, err.Error(), http.StatusInternalServerError)
	}

	_, appErr = manager.OrderUpdated(*order)
	if appErr != nil {
		return appErr
	}

	for _, fulfillment := range fulfillments {
		_, appErr = manager.FulfillmentCreated(*fulfillment)
		if appErr != nil {
			return appErr
		}
	}

	if order.Status == model.ORDER_STATUS_FULFILLED {
		_, appErr = manager.OrderFulfilled(*order)
		if appErr != nil {
			return appErr
		}
	}

	if notifyCustomer {
		for _, fulfillment := range fulfillments {
			appErr = a.SendFulfillmentConfirmationToCustomer(order, fulfillment, user, nil, manager)
			if appErr != nil {
				return appErr
			}
		}
	}

	return nil
}

// OrderAwaitsFulfillmentApproval
func (s *ServiceOrder) OrderAwaitsFulfillmentApproval(fulfillments []*model.Fulfillment, user *model.User, _ any, fulfillmentLines model.FulfillmentLines, manager interfaces.PluginManagerInterface, notifyCustomer bool) *model_helper.AppError {
	transaction := s.srv.Store.GetMaster().Begin()
	if transaction.Error != nil {
		return model_helper.NewAppError("OrderAwaitsFulfillmentApproval", model_helper.ErrorCreatingTransactionErrorID, nil, transaction.Error.Error(), http.StatusInternalServerError)
	}
	defer s.srv.Store.FinalizeTransaction(transaction)

	var order *model.Order
	if order_ := fulfillments[0].GetOrder(); order != nil {
		order = order_
	} else {
		order_, appErr := s.OrderById(fulfillments[0].OrderID)
		if appErr != nil {
			return appErr
		}
		order = order_
	}

	appErr := s.UpdateOrderStatus(transaction, *order)
	if appErr != nil {
		return appErr
	}
	// NOTE: from now on, order is guaranteed non-nil

	_, appErr = s.FulfillmentAwaitsApprovalEvent(transaction, order, user, nil, fulfillmentLines)
	if appErr != nil {
		return appErr
	}

	// commit transaction:
	if err := transaction.Commit().Error; err != nil {
		return model_helper.NewAppError("OrderAwaitsFulfillmentApproval", model_helper.ErrorCommittingTransactionErrorID, nil, err.Error(), http.StatusInternalServerError)
	}

	_, appErr = manager.OrderUpdated(*order)
	return appErr
}

// OrderShippingUpdated
func (a *ServiceOrder) OrderShippingUpdated(order model.Order, manager interfaces.PluginManagerInterface) *model_helper.AppError {
	appErr := a.RecalculateOrder(nil, &order, nil)
	if appErr != nil {
		return appErr
	}
	_, appErr = manager.OrderUpdated(order)
	return appErr
}

// OrderAuthorized
func (a *ServiceOrder) OrderAuthorized(order model.Order, user *model.User, _ any, amount *decimal.Decimal, payMent model.Payment, manager interfaces.PluginManagerInterface) *model_helper.AppError {
	var userID *string
	if user != nil && model_helper.IsValidId(user.Id) {
		userID = &user.Id
	}

	_, appErr := a.CommonCreateOrderEvent(nil, &model.OrderEventOption{
		Type:       model.ORDER_EVENT_TYPE_PAYMENT_AUTHORIZED,
		UserID:     userID,
		OrderID:    order.Id,
		Parameters: getPaymentData(amount, payMent),
	})
	if appErr != nil {
		return appErr
	}
	_, appErr = manager.OrderUpdated(order)
	return appErr
}

// OrderCaptured
func (a *ServiceOrder) OrderCaptured(order model.Order, user *model.User, _ any, amount *decimal.Decimal, payMent model.Payment, manager interfaces.PluginManagerInterface) (*model_helper.InsufficientStock, *model_helper.AppError) {
	var userID *string
	if user != nil && model_helper.IsValidId(user.Id) {
		userID = &user.Id
	}

	_, appErr := a.CommonCreateOrderEvent(nil, &model.OrderEventOption{
		OrderID:    order.Id,
		UserID:     userID,
		Type:       model.ORDER_EVENT_TYPE_PAYMENT_CAPTURED,
		Parameters: getPaymentData(amount, payMent),
	})
	if appErr != nil {
		return nil, appErr
	}

	_, appErr = manager.OrderUpdated(order)
	if appErr != nil {
		return nil, appErr
	}

	if order.IsFullyPaid() {
		insufficientStock, appErr := a.HandleFullyPaidOrder(manager, order, user, nil)
		if insufficientStock != nil || appErr != nil {
			return insufficientStock, appErr
		}
	}

	return nil, nil
}

// FulfillmentTrackingUpdated
func (a *ServiceOrder) FulfillmentTrackingUpdated(fulfillment *model.Fulfillment, user *model.User, _ any, trackingNumber string, manager interfaces.PluginManagerInterface) *model_helper.AppError {
	var order = fulfillment.GetOrder()
	if order == nil {
		order_, appErr := a.OrderById(fulfillment.OrderID)
		if appErr != nil {
			return appErr
		}
		order = order_
	}

	_, appErr := a.FulfillmentTrackingUpdatedEvent(order, user, nil, trackingNumber, fulfillment)
	return appErr
}

// CancelFulfillment Return products to corresponding stocks.
func (a *ServiceOrder) CancelFulfillment(fulfillment model.Fulfillment, user *model.User, _ any, warehouse *model.Warehouse, manager interfaces.PluginManagerInterface) (*model.Fulfillment, *model_helper.AppError) {
	// initialize a transaction
	transaction := a.srv.Store.GetMaster().Begin()
	if transaction.Error != nil {
		return nil, model_helper.NewAppError("CancelFulfillment", model_helper.ErrorCreatingTransactionErrorID, nil, transaction.Error.Error(), http.StatusInternalServerError)
	}
	defer a.srv.Store.FinalizeTransaction(transaction)

	// refetch fulfillment from database, lock for update
	fulfillment_, appErr := a.FulfillmentByOption(&model.FulfillmentFilterOption{
		Conditions:         squirrel.Eq{model.FulfillmentTableName + ".Id": fulfillment.Id},
		SelectForUpdate:    true, // this tells store to add `FOR UPDATE` to select query
		SelectRelatedOrder: true, // this make you free to acces `Order` of returning `fulfillment`
		Transaction:        transaction,
	})
	if appErr != nil {
		if appErr.StatusCode == http.StatusInternalServerError {
			return nil, appErr
		}
		// if error is not found error, this mean given `fulfillment` is not valid
		return nil, model_helper.NewAppError("CancelFulfillment", model_helper.InvalidArgumentAppErrorID, map[string]any{"Fields": "fulfillment"}, appErr.DetailedError, http.StatusBadRequest)
	}

	var userID *string
	if user != nil && model_helper.IsValidId(user.Id) {
		userID = &user.Id
	}
	_, appErr = a.FulfillmentCanceledEvent(transaction, fulfillment_.GetOrder(), user, nil, fulfillment_)
	if appErr != nil {
		return nil, appErr
	}

	if warehouse != nil {
		appErr = a.RestockFulfillmentLines(transaction, fulfillment_, warehouse)
		if appErr != nil {
			return nil, appErr
		}

		// get total quantity for order of given fulfillment:
		orderTotalQuantity, appErr := a.OrderTotalQuantity(fulfillment.OrderID)
		if appErr != nil {
			return nil, appErr
		}
		_, appErr = a.CommonCreateOrderEvent(transaction, &model.OrderEventOption{
			OrderID: fulfillment.OrderID,
			UserID:  userID,
			Type:    model.ORDER_EVENT_TYPE_FULFILLMENT_RESTOCKED_ITEMS,
			Parameters: model_types.JSONString{
				"quantity":  orderTotalQuantity,
				"warehouse": warehouse.Id,
			},
		})
		if appErr != nil {
			return nil, appErr
		}
	}

	fulfillment.Status = model.FULFILLMENT_CANCELED
	_, appErr = a.UpsertFulfillment(transaction, fulfillment_)
	if appErr != nil {
		return nil, appErr
	}

	appErr = a.UpdateOrderStatus(transaction, *fulfillment.GetOrder()) // you can access order here since store attached it to fulfillment above
	if appErr != nil {
		return nil, appErr
	}

	// commit transaction:
	if err := transaction.Commit().Error; err != nil {
		return nil, model_helper.NewAppError("CancelOrder", model_helper.ErrorCommittingTransactionErrorID, nil, err.Error(), http.StatusInternalServerError)
	}

	_, appErr = manager.FulfillmentCanceled(*fulfillment_)
	if appErr != nil {
		return nil, appErr
	}

	_, appErr = manager.OrderUpdated(*fulfillment_.GetOrder())
	if appErr != nil {
		return nil, appErr
	}

	return fulfillment_, nil
}

// CancelWaitingFulfillment cancels fulfillments which is in waiting for approval state.
func (s *ServiceOrder) CancelWaitingFulfillment(fulfillment model.Fulfillment, user *model.User, _ any, manager interfaces.PluginManagerInterface) *model_helper.AppError {
	// initialize a transaction
	transaction := s.srv.Store.GetMaster().Begin()
	if transaction.Error != nil {
		return model_helper.NewAppError("CancelWaitingFulfillment", model_helper.ErrorCreatingTransactionErrorID, nil, transaction.Error.Error(), http.StatusInternalServerError)
	}
	defer s.srv.Store.FinalizeTransaction(transaction)

	fulfillment_, appErr := s.FulfillmentByOption(&model.FulfillmentFilterOption{
		Conditions:         squirrel.Eq{model.FulfillmentTableName + ".Id": fulfillment.Id},
		SelectRelatedOrder: true, // this
		Transaction:        transaction,
	})
	if appErr != nil {
		return appErr
	}

	_, appErr = s.FulfillmentCanceledEvent(transaction, fulfillment_.GetOrder(), user, nil, fulfillment_)
	if appErr != nil {
		return appErr
	}

	fulfillmentLinesOfFulfillment, appErr := s.FulfillmentLinesByOption(&model.FulfillmentLineFilterOption{
		Conditions: squirrel.Eq{model.FulfillmentLineTableName + ".FulfillmentID": fulfillment_.Id},
		Preloads:   []string{"OrderLine"},
	})
	if appErr != nil {
		if appErr.StatusCode == http.StatusInternalServerError {
			return appErr
		}
	}

	var orderLines model.OrderLineSlice
	for _, line := range fulfillmentLinesOfFulfillment {
		orderLine := line.OrderLine
		orderLine.QuantityFulfilled -= line.Quantity
		orderLines = append(orderLines, orderLine)
	}

	_, appErr = s.BulkUpsertOrderLines(transaction, orderLines)
	if appErr != nil {
		return appErr
	}

	appErr = s.BulkDeleteFulfillments(transaction, []*model.Fulfillment{fulfillment_})
	if appErr != nil {
		return appErr
	}

	appErr = s.UpdateOrderStatus(transaction, *fulfillment_.GetOrder())
	if appErr != nil {
		return appErr
	}

	if err := transaction.Commit().Error; err != nil {
		return model_helper.NewAppError("CancelWaitingFulfillment", model_helper.ErrorCommittingTransactionErrorID, nil, err.Error(), http.StatusInternalServerError)
	}

	_, appErr = manager.FulfillmentCanceled(*fulfillment_)
	if appErr != nil {
		return appErr
	}

	_, appErr = manager.OrderUpdated(*fulfillment_.GetOrder())
	return appErr
}

func (s *ServiceOrder) ApproveFulfillment(fulfillment *model.Fulfillment, user *model.User, _ any, manager interfaces.PluginManagerInterface, settings model.ShopSettings, notifyCustomer bool, allowStockTobeExceeded bool) (*model.Fulfillment, *model_helper.InsufficientStock, *model_helper.AppError) {
	// initialize a transaction
	transaction := s.srv.Store.GetMaster().Begin()
	if transaction.Error != nil {
		return nil, nil, model_helper.NewAppError("ApproveFulfillment", model_helper.ErrorCreatingTransactionErrorID, nil, transaction.Error.Error(), http.StatusInternalServerError)
	}
	defer s.srv.Store.FinalizeTransaction(transaction)

	fulfillment.Status = model.FULFILLMENT_FULFILLED
	_, appErr := s.UpsertFulfillment(transaction, fulfillment)
	if appErr != nil {
		return nil, nil, appErr
	}

	order, appErr := s.OrderById(fulfillment.OrderID)
	if appErr != nil {
		return nil, nil, appErr
	}

	if notifyCustomer {
		appErr = s.SendFulfillmentConfirmationToCustomer(order.DeepCopy(), fulfillment, user, nil, manager)
		if appErr != nil {
			return nil, nil, appErr
		}
	}

	fulfillmentLines, appErr := s.FulfillmentLinesByOption(&model.FulfillmentLineFilterOption{
		Conditions: squirrel.Eq{model.FulfillmentLineTableName + ".FulfillmentID": fulfillment.Id},
		Preloads:   []string{"Orderline.ProductVariant"},
	})
	if appErr != nil {
		return nil, nil, appErr
	}

	_, appErr = s.FulfillmentFulfilledItemsEvent(transaction, order, user, nil, fulfillmentLines)
	if appErr != nil {
		return nil, nil, appErr
	}

	_, stocks, appErr := s.srv.WarehouseService().StocksByOption(&model.StockFilterOption{
		Conditions: squirrel.Eq{model.StockTableName + ".Id": fulfillmentLines.StockIDs()},
	})
	if appErr != nil {
		if appErr.StatusCode == http.StatusInternalServerError {
			return nil, nil, appErr
		}
	}
	// stocksMap has keys are stock ids
	var stocksMap = map[string]*model.Stock{}
	for _, stock := range stocks {
		stocksMap[stock.Id] = stock
	}

	var linesToFulfill []*model.OrderLineData
	for _, line := range fulfillmentLines {

		// determine warehouse id
		var warehouseID *string
		if stockID := line.StockID; stockID != nil && stocksMap[*stockID] != nil {
			warehouseID = &stocksMap[*stockID].WarehouseID
		}

		linesToFulfill = append(linesToFulfill, &model.OrderLineData{
			Line:        *line.OrderLine,
			Quantity:    line.Quantity,
			Variant:     line.OrderLine.ProductVariant, //
			WarehouseID: warehouseID,
		})
	}

	insufficientStock, appErr := s.decreaseStocks(linesToFulfill, manager, allowStockTobeExceeded)
	if insufficientStock != nil || appErr != nil {
		return nil, insufficientStock, appErr
	}

	// refetch order
	order, appErr = s.OrderById(fulfillment.OrderID)
	if appErr != nil {
		return nil, nil, appErr
	}

	appErr = s.UpdateOrderStatus(transaction, *order)
	if appErr != nil {
		return nil, nil, appErr
	}

	appErr = s.CreateGiftcardsWhenApprovingFulfillment(order, linesToFulfill, user, nil, manager, settings)
	if appErr != nil {
		return nil, nil, appErr
	}

	// commit
	if err := transaction.Commit().Error; err != nil {
		return nil, nil, model_helper.NewAppError("ApproveFulfillment", model_helper.ErrorCreatingTransactionErrorID, nil, err.Error(), http.StatusInternalServerError)
	}

	_, appErr = manager.OrderUpdated(*order)
	if appErr != nil {
		return nil, nil, appErr
	}

	if order.Status == model.ORDER_STATUS_FULFILLED {
		_, appErr = manager.OrderFulfilled(*order)
		if appErr != nil {
			return nil, nil, appErr
		}
	}

	return fulfillment, nil, nil
}

// CreateGiftcardsWhenApprovingFulfillment
func (s *ServiceOrder) CreateGiftcardsWhenApprovingFulfillment(order *model.Order, linesData []*model.OrderLineData, user *model.User, _ any, manager interfaces.PluginManagerInterface, settings model.ShopSettings) *model_helper.AppError {
	var (
		giftcardLines = model.OrderLineSlice{}
		quantities    = map[string]int{}
	)

	for _, lineData := range linesData {
		line := lineData.Line
		if line.IsGiftcard {
			giftcardLines = append(giftcardLines, &line)
			quantities[line.Id] = line.Quantity
		}
	}

	_, appErr := s.srv.GiftcardService().GiftcardsCreate(nil, order, giftcardLines, quantities, settings, user, nil, manager)
	return appErr
}

// Mark order as paid.
//
// Allows to create a payment for an order without actually performing any
// payment by the gateway.
//
// externalReference can be empty
func (a *ServiceOrder) MarkOrderAsPaid(order model.Order, requestUser *model.User, _ any, manager interfaces.PluginManagerInterface, externalReference string) (*model_helper.PaymentError, *model_helper.AppError) {
	transaction := a.srv.Store.GetMaster().Begin()
	if transaction.Error != nil {
		return nil, model_helper.NewAppError("MarkOrderAsPaid", model_helper.ErrorCreatingTransactionErrorID, nil, transaction.Error.Error(), http.StatusInternalServerError)
	}
	defer a.srv.Store.FinalizeTransaction(transaction)

	order.PopulateNonDbFields() // this is required

	payMent, paymentErr, appErr := a.srv.PaymentService().CreatePayment(transaction, model.GATE_WAY_MANUAL, nil, order.Total.Gross.Currency, order.UserEmail, "", "", nil, nil, &order, "", externalReference, "", nil)
	if appErr != nil || paymentErr != nil {
		return paymentErr, appErr
	}

	payMent.ChargeStatus = model.PAYMENT_CHARGE_STATUS_FULLY_CHARGED
	payMent.CapturedAmount = &order.Total.Gross.Amount

	savedPayment, appErr := a.srv.PaymentService().UpsertPayment(transaction, payMent)
	if appErr != nil {
		return nil, appErr
	}

	_, appErr = a.srv.PaymentService().SaveTransaction(transaction, &model.PaymentTransaction{
		PaymentID:       savedPayment.Id,
		ActionRequired:  false,
		Kind:            model.TRANSACTION_KIND_EXTERNAL,
		Token:           externalReference,
		IsSuccess:       true,
		Amount:          &order.Total.Gross.Amount,
		Currency:        order.Total.Gross.Currency,
		GatewayResponse: model_types.JSONString{},
	})
	if appErr != nil {
		return nil, appErr
	}

	_, appErr = a.OrderManuallyMarkedAsPaidEvent(transaction, order, requestUser, nil, externalReference)
	if appErr != nil {
		return nil, appErr
	}

	_, appErr = manager.OrderFullyPaid(order)
	if appErr != nil {
		return nil, appErr
	}
	_, appErr = manager.OrderUpdated(order)
	if appErr != nil {
		return nil, appErr
	}

	appErr = a.UpdateOrderTotalPaid(transaction, &order)
	if appErr != nil {
		return nil, appErr
	}

	if err := transaction.Commit().Error; err != nil {
		return nil, model_helper.NewAppError("CancelOrder", model_helper.ErrorCommittingTransactionErrorID, nil, err.Error(), http.StatusInternalServerError)
	}

	return nil, nil
}

// CleanMarkOrderAsPaid Check if an order can be marked as paid.
func (a *ServiceOrder) CleanMarkOrderAsPaid(order *model.Order) *model_helper.AppError {
	_, paymentsForOrder, appErr := a.srv.PaymentService().PaymentsByOption(&model.PaymentFilterOption{
		Conditions: squirrel.Eq{model.PaymentTableName + ".OrderID": order.Id},
	})
	if appErr != nil {
		return appErr
	}
	if len(paymentsForOrder) > 0 {
		return model_helper.NewAppError("CleanMarkOrderAsPaid", "app.order.mark_order_with_payments_as_paid.app_error", nil, "cannot mark order with payments as paid", http.StatusNotAcceptable)
	}

	return nil
}

func (s *ServiceOrder) increaseOrderLineQuantity(transaction boil.ContextTransactor, orderLinesInfo []*model.OrderLineData) *model_helper.AppError {
	orderLines := model.OrderLineSlice{}

	for _, lineInfo := range orderLinesInfo {
		line := lineInfo.Line
		line.Quantity += lineInfo.Quantity
		orderLines = append(orderLines, &line)
	}

	_, appErr := s.BulkUpsertOrderLines(transaction, orderLines)
	return appErr
}

// FulfillOrderLines Fulfill order line with given quantity
func (a *ServiceOrder) FulfillOrderLines(orderLineInfos []*model.OrderLineData, manager interfaces.PluginManagerInterface, allowStockTobeExceeded bool) (*model_helper.InsufficientStock, *model_helper.AppError) {
	transaction := a.srv.Store.GetMaster().Begin()
	if transaction.Error != nil {
		return nil, model_helper.NewAppError("FulfillOrderLines", model_helper.ErrorCreatingTransactionErrorID, nil, transaction.Error.Error(), http.StatusInternalServerError)
	}
	defer a.srv.Store.FinalizeTransaction(transaction)

	insufficientStock, appErr := a.decreaseStocks(orderLineInfos, manager, allowStockTobeExceeded)
	if insufficientStock != nil || appErr != nil {
		return insufficientStock, appErr
	}

	appErr = a.increaseOrderLineQuantity(transaction, orderLineInfos)
	if appErr != nil {
		return nil, appErr
	}

	if err := transaction.Commit().Error; err != nil {
		return nil, model_helper.NewAppError("AutomaticallyFulfillDigitalLines", model_helper.ErrorCommittingTransactionErrorID, nil, err.Error(), http.StatusInternalServerError)
	}

	return nil, nil
}

// AutomaticallyFulfillDigitalLines
// Fulfill all digital lines which have enabled automatic fulfillment setting. Send confirmation email afterward.
func (a *ServiceOrder) AutomaticallyFulfillDigitalLines(order model.Order, manager interfaces.PluginManagerInterface) (*model_helper.InsufficientStock, *model_helper.AppError) {
	transaction := a.srv.Store.GetMaster().Begin()
	if transaction.Error != nil {
		return nil, model_helper.NewAppError("AutomaticallyFulfillDigitalLines", model_helper.ErrorCreatingTransactionErrorID, nil, transaction.Error.Error(), http.StatusInternalServerError)
	}
	defer a.srv.Store.FinalizeTransaction(transaction)

	// find order lines of given order that are:
	// 1) NOT require shipping
	// 2) has ProductVariant attached AND that productVariant has a digitalContent accompanies
	digitalOrderLinesOfOrder, appErr := a.OrderLinesByOption(&model.OrderLineFilterOption{
		Conditions: squirrel.Expr("?.OrderID = ? AND ?.IsShippingRequired = false", model.OrderLineTableName, order.Id, model.OrderLineTableName),
		// VariantDigitalContentID: squirrel.NotEq{model.DigitalContentTableName + ".Id": nil},
		VariantDigitalContentID: squirrel.Expr(model.DigitalContentTableName + ".Id IS NOT NULL"),
		Preload:                 []string{"ProductVariant.DigitalContent"}, // TODO: check if this works
	})
	if appErr != nil {
		return nil, appErr
	}

	if len(digitalOrderLinesOfOrder) == 0 {
		return nil, nil
	}

	fulfillment, appErr := a.FulfillmentByOption(&model.FulfillmentFilterOption{
		Conditions:  squirrel.Eq{model.FulfillmentTableName + ".OrderID": order.Id},
		Transaction: transaction,
	})
	if appErr != nil {
		if appErr.StatusCode != http.StatusNotFound {
			return nil, appErr
		}

		fulfillment, appErr = a.UpsertFulfillment(transaction, &model.Fulfillment{
			OrderID: order.Id,
		})
		if appErr != nil {
			return nil, appErr
		}
	}

	var (
		fulfillmentLines []*model.FulfillmentLine
		orderLineDatas   []*model.OrderLineData
	)

	for _, orderLine := range digitalOrderLinesOfOrder {
		orderLineNeedsAutomaticFulfillment, appErr := a.OrderLineNeedsAutomaticFulfillment(orderLine)
		if appErr != nil {
			return nil, appErr // must return if error occured
		}
		if !orderLineNeedsAutomaticFulfillment {
			continue
		}

		if orderLine.ProductVariant != nil { // ProductVariant is available to use, prefetch option is enabled above
			_, appErr = a.srv.ProductService().UpsertDigitalContentURL(&model.DigitalContentUrl{
				LineID: &orderLine.Id,
			})
			if appErr != nil {
				return nil, appErr
			}
		}

		fulfillmentLines = append(fulfillmentLines, &model.FulfillmentLine{
			FulfillmentID: fulfillment.Id,
			OrderLineID:   orderLine.Id,
			Quantity:      orderLine.Quantity,
		})

		allocationsOfOrderLine, appErr := a.srv.WarehouseService().AllocationsByOption(&model.AllocationFilterOption{
			Conditions: squirrel.Eq{model.AllocationTableName + ".OrderLineID": orderLine.Id},
		})
		if appErr != nil {
			return nil, appErr
		}

		stock, appErr := a.srv.WarehouseService().GetStockById(allocationsOfOrderLine[0].StockID)
		if appErr != nil {
			return nil, appErr
		}

		orderLineDatas = append(orderLineDatas, &model.OrderLineData{
			Line:        *orderLine,
			Quantity:    orderLine.Quantity,
			Variant:     orderLine.ProductVariant,
			WarehouseID: &stock.WarehouseID,
		})
	}

	_, appErr = a.BulkUpsertFulfillmentLines(transaction, fulfillmentLines)
	if appErr != nil {
		return nil, appErr
	}
	insufficientStock, appErr := a.FulfillOrderLines(orderLineDatas, manager, false)

	if insufficientStock != nil || appErr != nil {
		return insufficientStock, appErr
	}

	var user *model.User // can be nil
	if order.UserID != nil {
		user, appErr = a.srv.AccountService().UserById(context.Background(), *order.UserID)
		if appErr != nil {
			return nil, appErr
		}
	}
	appErr = a.SendFulfillmentConfirmationToCustomer(&order, fulfillment, user, nil, manager)
	if appErr != nil {
		return nil, appErr
	}

	appErr = a.UpdateOrderStatus(transaction, order)
	if appErr != nil {
		return nil, appErr
	}

	if err := transaction.Commit().Error; err != nil {
		return nil, model_helper.NewAppError("AutomaticallyFulfillDigitalLines", model_helper.ErrorCommittingTransactionErrorID, nil, err.Error(), http.StatusInternalServerError)
	}

	return nil, nil
}

// Modify stocks and allocations. Return list of unsaved FulfillmentLines.
//
//	Args:
//	    fulfillment (Fulfillment): Fulfillment to create lines
//	    warehouse_pk (str): Warehouse to fulfill order.
//	    lines_data (List[Dict]): List with information from which system
//	        create FulfillmentLines. Example:
//	            [
//	                {
//	                    "order_line": (OrderLine),
//	                    "quantity": (int),
//	                },
//	                ...
//	            ]
//	    channel_slug (str): Channel for which fulfillment lines should be created.
//
//	Return:
//	    List[FulfillmentLine]: Unsaved fulfillmet lines created for this fulfillment
//	        based on information form `lines`
//
//	Raise:
//	    InsufficientStock: If system hasn't containt enough item in stock for any line.
func (a *ServiceOrder) createFulfillmentLines(fulfillment *model.Fulfillment, warehouseID string, lineDatas model.QuantityOrderLines, channelID string, manager interfaces.PluginManagerInterface, decreaseStock bool, allowStockTobeExceeded bool) ([]*model.FulfillmentLine, *model_helper.InsufficientStock, *model_helper.AppError) {

	var variantIDs = lineDatas.OrderLines().ProductVariantIDs()

	stocks, appErr := a.srv.WarehouseService().FilterStocksForChannel(&model.StockFilterForChannelOption{
		ChannelID: channelID,
		Conditions: squirrel.Eq{
			model.StockTableName + ".WarehouseID":      warehouseID,
			model.StockTableName + ".ProductVariantID": variantIDs,
		},
		SelectRelatedProductVariant: true,
	})
	if appErr != nil {
		if appErr.StatusCode == http.StatusInternalServerError {
			return nil, nil, appErr
		}
	}

	productVariants, appErr := a.srv.ProductService().ProductVariantsByOption(&model.ProductVariantFilterOption{
		Conditions: squirrel.Eq{model.ProductVariantTableName + ".Id": variantIDs},
		Preloads:   []string{"DigitalContent"},
	})
	if appErr != nil {
		if appErr.StatusCode == http.StatusInternalServerError {
			return nil, nil, appErr
		}
	}

	// productVariantsMap has keys are product variant ids
	productVariantsMap := map[string]*model.ProductVariant{}
	for _, variant := range productVariants {
		productVariantsMap[variant.Id] = variant
	}

	// variantToStock map has keys are product variant ids
	variantToStock := map[string][]*model.Stock{}
	for _, stock := range stocks {
		variantToStock[stock.ProductVariantID] = append(variantToStock[stock.ProductVariantID], stock)
	}

	var (
		insufficientStocks              []*model.InsufficientStockData
		fulfillmentLines                []*model.FulfillmentLine
		linesInfo                       []*model.OrderLineData
		quantity                        int
		orderLine                       *model.OrderLine
		productVariantOfOrderLine       model.ProductVariant
		productVariantOfOrderLineIsReal bool
	)

	for _, line := range lineDatas {

		productVariantOfOrderLineIsReal = false
		quantity = line.Quantity
		orderLine = line.OrderLine
		productVariantOfOrderLine = model.ProductVariant{}

		if orderLine.VariantID != nil && productVariantsMap[*orderLine.VariantID] != nil {
			productVariantOfOrderLine = *(productVariantsMap[*orderLine.VariantID])
			productVariantOfOrderLineIsReal = true
		}

		if quantity > 0 {
			if orderLine.VariantID == nil || variantToStock[*orderLine.VariantID] == nil {
				insufficientStocks = append(insufficientStocks, &model.InsufficientStockData{
					Variant:     productVariantOfOrderLine,
					OrderLine:   orderLine,
					WarehouseID: &warehouseID,
				})
				continue
			}

			stock := variantToStock[*orderLine.VariantID][0]
			linesInfo = append(linesInfo, &model.OrderLineData{
				Line:        *orderLine,
				Quantity:    line.Quantity,
				Variant:     &productVariantOfOrderLine,
				WarehouseID: &warehouseID,
			})

			orderLineIsDigital, appErr := a.srv.OrderService().OrderLineIsDigital(orderLine)
			if appErr != nil {
				return nil, nil, appErr
			}
			if orderLineIsDigital && productVariantOfOrderLineIsReal {

				_, appErr = a.srv.ProductService().UpsertDigitalContentURL(&model.DigitalContentUrl{
					ContentID: productVariantOfOrderLine.DigitalContent.Id, // check out 2nd goroutine above to see why is it possible to access DigitalContent.
					LineID:    &orderLine.Id,
				})
				if appErr != nil {
					return nil, nil, appErr
				}
			}

			fulfillmentLines = append(fulfillmentLines, &model.FulfillmentLine{
				OrderLineID:   orderLine.Id,
				FulfillmentID: fulfillment.Id,
				Quantity:      line.Quantity,
				StockID:       &stock.Id,
			})
		}
	}

	if len(insufficientStocks) > 0 {
		return nil,
			&model_helper.InsufficientStock{
				Items: insufficientStocks,
			},
			nil
	}

	if len(linesInfo) > 0 {
		if decreaseStock {
			insufficientStock, appErr := a.decreaseStocks(linesInfo, manager, allowStockTobeExceeded)
			if insufficientStock != nil || appErr != nil {
				return nil, insufficientStock, appErr
			}
		}

		appErr := a.increaseOrderLineQuantity(nil, linesInfo)
		if appErr != nil {
			return nil, nil, appErr
		}
	}

	return fulfillmentLines, nil, nil
}

// Fulfill order.
//
//	Function create fulfillments with lines.
//	Next updates Order based on created fulfillments.
//
//	Args:
//	    requester (User): Requester who trigger this action.
//	    order (Order): Order to fulfill
//	    fulfillment_lines_for_warehouses (Dict): Dict with information from which
//	        system create fulfillments. Example:
//	            {
//	                (Warehouse.pk): [
//	                    {
//	                        "order_line": (OrderLine),
//	                        "quantity": (int),
//	                    },
//	                    ...
//	                ]
//	            }
//	    manager (PluginsManager): Base manager for handling plugins logic.
//	    notify_customer (bool): If `True` system send email about
//	        fulfillments to customer.
//
//	Return:
//	    List[Fulfillment]: Fulfillmet with lines created for this order
//	        based on information form `fulfillment_lines_for_warehouses`
//
//
//	Raise:
//	    InsufficientStock: If system hasn't containt enough item in stock for any line.
func (a *ServiceOrder) CreateFulfillments(user *model.User, _ any, order *model.Order, fulfillmentLinesForWarehouses map[string][]*model.QuantityOrderLine, manager interfaces.PluginManagerInterface, notifyCustomer bool, approved bool, allowStockTobeExceeded bool) ([]*model.Fulfillment, *model_helper.InsufficientStock, *model_helper.AppError) {
	transaction := a.srv.Store.GetMaster().Begin()
	if transaction.Error != nil {
		return nil, nil, model_helper.NewAppError("CreateFulfillments", model_helper.ErrorCreatingTransactionErrorID, nil, transaction.Error.Error(), http.StatusInternalServerError)
	}
	defer a.srv.Store.FinalizeTransaction(transaction)

	var (
		fulfillments     []*model.Fulfillment
		fulfillmentLines []*model.FulfillmentLine
	)

	channel, appErr := a.srv.ChannelService().ChannelByOption(&model.ChannelFilterOption{
		Conditions: squirrel.Eq{model.ChannelTableName + ".Id": order.ChannelID},
	})
	if appErr != nil {
		return nil, nil, appErr
	}

	status := model.FULFILLMENT_FULFILLED
	if !approved {
		status = model.FULFILLMENT_WAITING_FOR_APPROVAL
	}
	for warehouseID, quantityOrderLine := range fulfillmentLinesForWarehouses {
		fulfillment, appErr := a.UpsertFulfillment(transaction, &model.Fulfillment{
			Status:  status,
			OrderID: order.Id,
		})
		if appErr != nil {
			return nil, nil, appErr
		}

		fulfillments = append(fulfillments, fulfillment)

		filmentLines, insufficientStockErr, appErr := a.createFulfillmentLines(
			fulfillment,
			warehouseID,
			quantityOrderLine,
			channel.Slug,
			manager,
			approved,
			allowStockTobeExceeded,
		)
		if insufficientStockErr != nil || appErr != nil {
			return nil, insufficientStockErr, appErr
		}

		fulfillmentLines = append(fulfillmentLines, filmentLines...)
	}

	fulfillmentLines, appErr = a.BulkUpsertFulfillmentLines(transaction, fulfillmentLines)
	if appErr != nil {
		return nil, nil, appErr
	}

	// commit transaction
	if err := transaction.Commit().Error; err != nil {
		return nil, nil, model_helper.NewAppError("CreateFulfillments", model_helper.ErrorCommittingTransactionErrorID, nil, err.Error(), http.StatusInternalServerError)
	}

	if approved {
		appErr = a.OrderFulfilled(fulfillments, user, nil, fulfillmentLines, manager, notifyCustomer)
	} else {
		appErr = a.OrderAwaitsFulfillmentApproval(fulfillments, user, nil, fulfillmentLines, manager, notifyCustomer)
	}

	if appErr != nil {
		return nil, nil, appErr
	}

	return fulfillments, nil, nil
}

// getFulfillmentLineIfExists
//
// NOTE: stockID can be empty
func (a *ServiceOrder) getFulfillmentLineIfExists(fulfillmentLines []*model.FulfillmentLine, orderLineID string, stockID *string) *model.FulfillmentLine {
	for _, line := range fulfillmentLines {
		if line.OrderLineID == orderLineID &&
			(line.StockID != nil && stockID != nil && *line.StockID == *stockID) {
			return line
		}
	}

	return nil
}

type AResult struct {
	MovedFulfillmentLine *model.FulfillmentLine
	FulfillmentLineExist bool
}

// getFulfillmentLine Get fulfillment line if extists or create new fulfillment line object.
//
// NOTE: stockID can be empty
func (a *ServiceOrder) getFulfillmentLine(targetFulfillment *model.Fulfillment, linesInTargetFulfillment []*model.FulfillmentLine, orderLineID string, stockID *string) *AResult {
	// Check if line for order_line_id and stock_id does not exist in DB.
	movedFulfillmentLine := a.getFulfillmentLineIfExists(linesInTargetFulfillment, orderLineID, stockID)

	fulfillmentLineExisted := true

	if movedFulfillmentLine == nil {
		// Create new not saved FulfillmentLine object and assign it to target fulfillment
		fulfillmentLineExisted = false
		movedFulfillmentLine = &model.FulfillmentLine{
			FulfillmentID: targetFulfillment.Id,
			OrderLineID:   orderLineID,
			StockID:       stockID,
			Quantity:      0,
		}
	}

	return &AResult{
		MovedFulfillmentLine: movedFulfillmentLine,
		FulfillmentLineExist: fulfillmentLineExisted,
	}
}

// moveOrderLinesToTargetFulfillment Move order lines with given quantity to the target fulfillment
func (a *ServiceOrder) moveOrderLinesToTargetFulfillment(orderLinesToMove []*model.OrderLineData, targetFulfillment *model.Fulfillment, manager interfaces.PluginManagerInterface) (fulfillmentLineToCreate []*model.FulfillmentLine, appErr *model_helper.AppError) {
	transaction := a.srv.Store.GetMaster().Begin()
	if transaction.Error != nil {
		return nil, model_helper.NewAppError("moveOrderLinesToTargetFulfillment", model_helper.ErrorCreatingTransactionErrorID, nil, transaction.Error.Error(), http.StatusInternalServerError)
	}
	defer a.srv.Store.FinalizeTransaction(transaction)

	var (
		orderLinesToUpdate        model.OrderLineSlice
		orderLineDatasToDeAlocate []*model.OrderLineData
	)

	for _, lineData := range orderLinesToMove {
		// calculate the quantity fulfilled/unfulfilled to move
		unFulfilledToMove := min(lineData.Line.QuantityUnFulfilled(), lineData.Quantity)
		lineData.Line.QuantityFulfilled += unFulfilledToMove

		// update current lines with new value of quantity
		orderLinesToUpdate = append(orderLinesToUpdate, &lineData.Line)
		fulfillmentLineToCreate = append(fulfillmentLineToCreate, &model.FulfillmentLine{
			FulfillmentID: targetFulfillment.Id,
			OrderLineID:   lineData.Line.Id,
			StockID:       nil,
			Quantity:      unFulfilledToMove,
		})

		allocationsOfOrderLine, appErr := a.srv.WarehouseService().AllocationsByOption(&model.AllocationFilterOption{
			Conditions: squirrel.Eq{model.AllocationTableName + ".OrderLineID": lineData.Line.Id},
		})
		if appErr != nil {
			if appErr.StatusCode == http.StatusInternalServerError {
				return nil, appErr
			}
		}

		if len(allocationsOfOrderLine) > 0 {
			orderLineDatasToDeAlocate = append(orderLineDatasToDeAlocate, &model.OrderLineData{
				Line:     lineData.Line,
				Quantity: unFulfilledToMove,
			})
		}
	}

	if len(orderLineDatasToDeAlocate) > 0 {
		allocationErr, appErr := a.srv.WarehouseService().DeallocateStock(orderLineDatasToDeAlocate, manager)
		if appErr != nil {
			return nil, appErr
		}

		if allocationErr != nil {
			slog.Warn("Unable to deallocate stock for order lines", slog.String("lines", strings.Join(allocationErr.OrderLines.IDs(), ", ")))
		}
	}

	fulfillmentLineToCreate, appErr = a.BulkUpsertFulfillmentLines(transaction, fulfillmentLineToCreate)
	if appErr != nil {
		return nil, appErr
	}

	_, appErr = a.BulkUpsertOrderLines(transaction, orderLinesToUpdate)
	if appErr != nil {
		return nil, appErr
	}

	if err := transaction.Commit().Error; err != nil {
		return nil, model_helper.NewAppError("moveOrderLinesToTargetFulfillment", model_helper.ErrorCommittingTransactionErrorID, nil, err.Error(), http.StatusInternalServerError)
	}

	return fulfillmentLineToCreate, nil
}

// moveFulfillmentLinesToTargetFulfillment Move fulfillment lines with given quantity to the target fulfillment
func (a *ServiceOrder) moveFulfillmentLinesToTargetFulfillment(fulfillmentLinesToMove []*model.FulfillmentLineData, linesInTargetFulfillment []*model.FulfillmentLine, targetFulfillment *model.Fulfillment) *model_helper.AppError {
	transaction := a.srv.Store.GetMaster().Begin()
	if transaction.Error != nil {
		return model_helper.NewAppError("moveFulfillmentLinesToTargetFulfillment", model_helper.ErrorCreatingTransactionErrorID, nil, transaction.Error.Error(), http.StatusInternalServerError)
	}
	defer a.srv.Store.FinalizeTransaction(transaction)

	var (
		fulfillmentLinesToCreate      []*model.FulfillmentLine
		fulfillmentLinesToUpdate      []*model.FulfillmentLine
		emptyFulfillmentLinesToDelete model.FulfillmentLines
	)

	for _, fulfillmentLineData := range fulfillmentLinesToMove {
		fulfillmentLine := fulfillmentLineData.Line
		quantityToMove := fulfillmentLineData.Quantity

		res := a.getFulfillmentLine(targetFulfillment, linesInTargetFulfillment, fulfillmentLine.OrderLineID, fulfillmentLine.StockID)

		// calculate the quantity fulfilled/unfulfilled/to move
		fulfilledToMove := min(fulfillmentLine.Quantity, quantityToMove)
		quantityToMove -= fulfilledToMove
		res.MovedFulfillmentLine.Quantity += fulfilledToMove
		fulfillmentLine.Quantity -= fulfilledToMove

		if fulfillmentLine.Quantity == 0 {
			// the fulfillment line without any items will be deleted
			emptyFulfillmentLinesToDelete = append(emptyFulfillmentLinesToDelete, &fulfillmentLine)
		} else {
			// update with new quantity value
			fulfillmentLinesToUpdate = append(fulfillmentLinesToUpdate, &fulfillmentLine)
		}

		if res.MovedFulfillmentLine.Quantity > 0 && !res.FulfillmentLineExist {
			// If this is new type of (order_line, stock) then we create new fulfillment line
			fulfillmentLinesToCreate = append(fulfillmentLinesToCreate, res.MovedFulfillmentLine)
		} else if res.FulfillmentLineExist {
			// if target fulfillment already have the same line, we  just update the quantity
			fulfillmentLinesToUpdate = append(fulfillmentLinesToUpdate, res.MovedFulfillmentLine)
		}
	}

	_, appErr := a.BulkUpsertFulfillmentLines(transaction, fulfillmentLinesToUpdate)
	if appErr != nil {
		return appErr
	}

	_, appErr = a.BulkUpsertFulfillmentLines(transaction, fulfillmentLinesToCreate)
	if appErr != nil {
		return appErr
	}

	appErr = a.DeleteFulfillmentLinesByOption(transaction, &model.FulfillmentLineFilterOption{
		Conditions: squirrel.Eq{model.FulfillmentLineTableName + ".Id": emptyFulfillmentLinesToDelete.IDs()},
	})
	if appErr != nil {
		return appErr
	}

	if err := transaction.Commit().Error; err != nil {
		return model_helper.NewAppError("moveFulfillmentLinesToTargetFulfillment", model_helper.ErrorCommittingTransactionErrorID, nil, err.Error(), http.StatusInternalServerError)
	}

	return nil
}

func getShippingRefundAmount(refundShippingCosts bool, refundAmount *decimal.Decimal, shippingPrice *decimal.Decimal) *decimal.Decimal {
	// We set shipping refund amount only when refund amount is calculated
	var shippingRefundAmount *decimal.Decimal
	if refundShippingCosts && refundAmount == nil {
		shippingRefundAmount = shippingPrice
	}
	return shippingRefundAmount
}

// Proceed with all steps required for refunding products.
//
// Calculate refunds for products based on the order's lines and fulfillment
// lines.  The logic takes the list of order lines, fulfillment lines, and their
// quantities which is used to create the refund fulfillment. The stock for
// unfulfilled lines will be deallocated.
//
// NOTE: `refundShippingCosts` default to false
func (a *ServiceOrder) CreateRefundFulfillment(
	requester *model.User,
	_ any,
	order model.Order,
	payMent model.Payment,
	orderLinesToRefund []*model.OrderLineData,
	fulfillmentLinesToRefund []*model.FulfillmentLineData,
	manager interfaces.PluginManagerInterface,
	amount *decimal.Decimal,
	refundShippingCosts bool,

) (*model.Fulfillment, *model_helper.PaymentError, *model_helper.AppError) {
	shippingRefundAmount := getShippingRefundAmount(refundShippingCosts, amount, order.ShippingPriceGrossAmount)

	// begin transaction
	transaction := a.srv.Store.GetMaster().Begin()
	if transaction.Error != nil {
		return nil, nil, model_helper.NewAppError("CreateRefundFulfillment", model_helper.ErrorCreatingTransactionErrorID, nil, transaction.Error.Error(), http.StatusInternalServerError)
	}
	defer a.srv.Store.FinalizeTransaction(transaction)

	totalRefundAmount, paymentErr, appErr := a.processRefund(requester, nil, order, payMent, orderLinesToRefund, fulfillmentLinesToRefund, amount, refundShippingCosts, manager)
	if paymentErr != nil || appErr != nil {
		return nil, paymentErr, appErr
	}

	refundedFulfillment, appErr := a.UpsertFulfillment(transaction, &model.Fulfillment{
		Status:               model.FULFILLMENT_REFUNDED,
		OrderID:              order.Id,
		TotalRefundAmount:    totalRefundAmount,
		ShippingRefundAmount: shippingRefundAmount,
	})
	if appErr != nil {
		return nil, nil, appErr
	}

	createdFulfillmentLines, appErr := a.moveOrderLinesToTargetFulfillment(orderLinesToRefund, refundedFulfillment, manager)
	if appErr != nil {
		return nil, nil, appErr
	}

	appErr = a.moveFulfillmentLinesToTargetFulfillment(fulfillmentLinesToRefund, createdFulfillmentLines, refundedFulfillment)
	if appErr != nil {
		return nil, nil, appErr
	}

	// delete fulfillments without lines after lines are removed
	fulfillments, appErr := a.FulfillmentsByOption(&model.FulfillmentFilterOption{
		Conditions: squirrel.Eq{
			model.FulfillmentTableName + ".OrderID": order.Id,
			model.FulfillmentTableName + ".Status": []model.FulfillmentStatus{
				model.FULFILLMENT_FULFILLED,
				model.FULFILLMENT_WAITING_FOR_APPROVAL,
			},
		},
		FulfillmentLineID: squirrel.Eq{model.FulfillmentLineTableName + ".Id": nil},
		Transaction:       transaction,
	})
	if appErr != nil {
		return nil, nil, appErr
	}

	if len(fulfillments) != 0 {
		appErr = a.BulkDeleteFulfillments(transaction, fulfillments)
		if appErr != nil {
			return nil, nil, appErr
		}
	}

	// commit transaction
	if err := transaction.Commit().Error; err != nil {
		return nil, nil, model_helper.NewAppError("CreateRefundFulfillment", model_helper.ErrorCommittingTransactionErrorID, nil, err.Error(), http.StatusInternalServerError)
	}

	_, appErr = manager.OrderUpdated(order)
	if appErr != nil {
		return nil, nil, appErr
	}

	return refundedFulfillment, nil, nil
}

// populateReplaceOrderFields create new order based on the state of given originalOrder
//
// If original order has shippingAddress/billingAddress, the new order copy these address(es) and change their IDs
func (a *ServiceOrder) populateReplaceOrderFields(transaction boil.ContextTransactor, originalOrder model.Order) (replaceOrder *model.Order, appErr *model_helper.AppError) {
	replaceOrder = &model.Order{
		Status:             model.ORDER_STATUS_DRAFT,
		UserID:             originalOrder.UserID,
		LanguageCode:       originalOrder.LanguageCode,
		UserEmail:          originalOrder.UserEmail,
		Currency:           originalOrder.Currency,
		ChannelID:          originalOrder.ChannelID,
		DisplayGrossPrices: originalOrder.DisplayGrossPrices,
		RedirectUrl:        originalOrder.RedirectUrl,
		OriginalID:         &originalOrder.Id,
		Origin:             model.ORDER_ORIGIN_REISSUE,
		ModelMetadata: model.ModelMetadata{
			Metadata:        originalOrder.Metadata,
			PrivateMetadata: originalOrder.PrivateMetadata,
		},
	}

	originalOrderAddressIDs := []string{}
	if originalOrder.BillingAddressID != nil {
		originalOrderAddressIDs = append(originalOrderAddressIDs, *originalOrder.BillingAddressID)
	}
	if originalOrder.ShippingAddressID != nil {
		originalOrderAddressIDs = append(originalOrderAddressIDs, *originalOrder.ShippingAddressID)
	}

	if len(originalOrderAddressIDs) > 0 {
		addressesOfOriginalOrder, appErr := a.srv.AccountService().AddressesByOption(&model.AddressFilterOption{
			Conditions: squirrel.Eq{model.AddressTableName + ".Id": originalOrderAddressIDs},
		})
		if appErr != nil {
			return nil, appErr
		}

		for _, address := range addressesOfOriginalOrder {
			originalOrderAddressID := address.Id
			address.Id = ""
			newAddress, appErr := a.srv.AccountService().UpsertAddress(transaction, address)
			if appErr != nil {
				return nil, appErr
			}

			if originalOrder.BillingAddressID != nil && originalOrderAddressID == *originalOrder.BillingAddressID {
				replaceOrder.BillingAddressID = &newAddress.Id
			} else if originalOrder.ShippingAddressID != nil && originalOrderAddressID == *originalOrder.ShippingAddressID {
				replaceOrder.ShippingAddressID = &newAddress.Id
			}
		}
	}

	return a.UpsertOrder(transaction, replaceOrder)
}

// CreateReplaceOrder Create draft order with lines to replace
func (a *ServiceOrder) CreateReplaceOrder(user *model.User, _ any, originalOrder model.Order, orderLinesToReplace []*model.OrderLineData, fulfillmentLinesToReplace []*model.FulfillmentLineData) (*model.Order, *model_helper.AppError) {
	transaction := a.srv.Store.GetMaster().Begin()
	if transaction.Error != nil {
		return nil, model_helper.NewAppError("CreateReplaceOrder", model_helper.ErrorCreatingTransactionErrorID, nil, transaction.Error.Error(), http.StatusInternalServerError)
	}
	defer a.srv.Store.FinalizeTransaction(transaction)

	replaceOrder, appErr := a.populateReplaceOrderFields(transaction, originalOrder)
	if appErr != nil {
		return nil, appErr
	}

	orderLinesToCreateMap := map[string]*model.OrderLine{}

	// iterate over lines without fulfillment to get the items for replace.
	// deepcopy to not lose the reference for lines assigned to original order
	for _, orderLineData := range model.OrderLineDatas(orderLinesToReplace).DeepCopy() {
		orderLine := orderLineData.Line
		orderLineID := orderLine.Id

		orderLine.Id = ""
		orderLine.OrderID = replaceOrder.Id
		orderLine.Quantity = orderLineData.Quantity
		orderLine.QuantityFulfilled = 0
		// we set order_line_id as a key to use it for iterating over fulfillment items
		orderLinesToCreateMap[orderLineID] = &orderLine
	}

	orderLineWithFulfillmentIDs := []string{}
	for _, lineData := range fulfillmentLinesToReplace {
		orderLineWithFulfillmentIDs = append(orderLineWithFulfillmentIDs, lineData.Line.OrderLineID)
	}

	orderLinesWithFulfillment, appErr := a.OrderLinesByOption(&model.OrderLineFilterOption{
		Conditions: squirrel.Eq{model.OrderLineTableName + ".Id": orderLineWithFulfillmentIDs},
	})
	if appErr != nil {
		return nil, appErr
	}

	orderLinesWithFulfillmentMap := map[string]*model.OrderLine{}
	for _, id := range orderLineWithFulfillmentIDs {
		for _, orderLine := range orderLinesWithFulfillment {
			if id == orderLine.Id {
				orderLinesWithFulfillmentMap[id] = orderLine
			}
		}
	}

	for _, fulfillmentLineData := range fulfillmentLinesToReplace {
		fulfillmentLine := fulfillmentLineData.Line
		orderLineID := fulfillmentLine.OrderLineID

		// if order_line_id exists in order_line_to_create, it means that we already have
		// prepared new order_line for this fulfillment. In that case we need to increase
		// quantity amount of new order_line by fulfillment_line.quantity
		if item, exist := orderLinesToCreateMap[orderLineID]; exist && item != nil {
			orderLinesToCreateMap[orderLineID].Quantity += fulfillmentLineData.Quantity
			continue
		}

		orderLine := orderLinesWithFulfillmentMap[orderLineID]
		orderLineID = orderLine.Id
		orderLine.Id = ""
		orderLine.OrderID = replaceOrder.Id
		orderLine.Quantity = fulfillmentLineData.Quantity
		orderLine.QuantityFulfilled = 0
		orderLinesToCreateMap[orderLineID] = orderLine
	}

	orderLinesToCreate := model.OrderLineSlice{}
	for _, orderLine := range orderLinesToCreateMap {
		orderLinesToCreate = append(orderLinesToCreate, orderLine)
	}

	_, appErr = a.BulkUpsertOrderLines(transaction, orderLinesToCreate)
	if appErr != nil {
		return nil, appErr
	}

	appErr = a.RecalculateOrder(transaction, replaceOrder, nil)
	if appErr != nil {
		return nil, appErr
	}

	lines := []*model.QuantityOrderLine{}
	for _, line := range orderLinesToCreate {
		lines = append(lines, &model.QuantityOrderLine{Quantity: line.Quantity, OrderLine: line})
	}
	_, appErr = a.DraftOrderCreatedFromReplaceEvent(transaction, *replaceOrder, originalOrder, user, nil, lines)
	if appErr != nil {
		return nil, appErr
	}

	// commit transaction
	err := transaction.Commit().Error
	if err != nil {
		return nil, model_helper.NewAppError("CreateReplaceOrder", model_helper.ErrorCommittingTransactionErrorID, nil, err.Error(), http.StatusInternalServerError)
	}

	return replaceOrder, nil
}

func (a *ServiceOrder) moveLinesToReturnFulfillment(
	orderLineDatas []*model.OrderLineData,
	fulfillmentLineDatas []*model.FulfillmentLineData,
	fulfillmentStatus model.FulfillmentStatus,
	order model.Order,
	totalRefundAmount *decimal.Decimal,
	shippingRefundAmount *decimal.Decimal,
	manager interfaces.PluginManagerInterface,

) (*model.Fulfillment, *model_helper.AppError) {

	targetFulfillment, appErr := a.UpsertFulfillment(nil, &model.Fulfillment{
		Status:               fulfillmentStatus,
		OrderID:              order.Id,
		TotalRefundAmount:    totalRefundAmount,
		ShippingRefundAmount: shippingRefundAmount,
	})
	if appErr != nil {
		return nil, appErr
	}

	LinesInTargetFulfillment, appErr := a.moveOrderLinesToTargetFulfillment(orderLineDatas, targetFulfillment, manager)
	if appErr != nil {
		return nil, appErr
	}

	fulfillmentLinesAlreadyRefunded, appErr := a.FulfillmentLinesByOption(&model.FulfillmentLineFilterOption{
		FulfillmentOrderID: squirrel.Eq{model.FulfillmentTableName + ".OrderID": order.Id},
		FulfillmentStatus:  squirrel.Eq{model.FulfillmentTableName + ".Status": model.FULFILLMENT_REFUNDED},
	})
	if appErr != nil {
		return nil, appErr
	}
	fulfillmentLinesAlreadyRefundedMap := lo.SliceToMap(fulfillmentLinesAlreadyRefunded, func(f *model.FulfillmentLine) (string, *model.FulfillmentLine) {
		return f.Id, f
	})

	var (
		refundedFulfillmentLinesToReturn []*model.FulfillmentLineData
		fulfillmentLinesToReturn         []*model.FulfillmentLineData
	)

	for _, lineData := range fulfillmentLineDatas {
		if item, exist := fulfillmentLinesAlreadyRefundedMap[lineData.Line.Id]; exist && item != nil {
			refundedFulfillmentLinesToReturn = append(refundedFulfillmentLinesToReturn, lineData)
			continue
		}

		fulfillmentLinesToReturn = append(fulfillmentLinesToReturn, lineData)
	}

	appErr = a.moveFulfillmentLinesToTargetFulfillment(fulfillmentLinesToReturn, LinesInTargetFulfillment, targetFulfillment)
	if appErr != nil {
		return nil, appErr
	}

	if len(refundedFulfillmentLinesToReturn) > 0 {
		var refundAndReturnFulfillment *model.Fulfillment
		if fulfillmentStatus == model.FULFILLMENT_REFUNDED_AND_RETURNED {
			refundAndReturnFulfillment = targetFulfillment
		} else {
			refundAndReturnFulfillment, appErr = a.UpsertFulfillment(nil, &model.Fulfillment{
				Status:  model.FULFILLMENT_REFUNDED_AND_RETURNED,
				OrderID: order.Id,
			})
			if appErr != nil {
				return nil, appErr
			}
		}

		appErr = a.moveFulfillmentLinesToTargetFulfillment(refundedFulfillmentLinesToReturn, []*model.FulfillmentLine{}, refundAndReturnFulfillment)
		if appErr != nil {
			return nil, appErr
		}
	}

	return targetFulfillment, nil
}

func (a *ServiceOrder) moveLinesToReplaceFulfillment(
	orderLinesToReplace []*model.OrderLineData,
	fulfillmentLinesToReplace []*model.FulfillmentLineData,
	order model.Order,
	manager interfaces.PluginManagerInterface,

) (*model.Fulfillment, *model_helper.AppError) {

	targetFulfillment, appErr := a.UpsertFulfillment(nil, &model.Fulfillment{
		Status:  model.FULFILLMENT_REPLACED,
		OrderID: order.Id,
	})
	if appErr != nil {
		return nil, appErr
	}

	linesInTargetFulfillment, appErr := a.moveOrderLinesToTargetFulfillment(orderLinesToReplace, targetFulfillment, manager)
	if appErr != nil {
		return nil, appErr
	}

	appErr = a.moveFulfillmentLinesToTargetFulfillment(
		fulfillmentLinesToReplace,
		linesInTargetFulfillment,
		targetFulfillment,
	)

	return targetFulfillment, appErr
}

func (a *ServiceOrder) CreateReturnFulfillment(
	requester *model.User, // can be nil
	order model.Order,
	orderLineDatas []*model.OrderLineData,
	fulfillmentLineDatas []*model.FulfillmentLineData,
	totalRefundAmount *decimal.Decimal, // can be nil
	shippingRefundAmount *decimal.Decimal, // can be nil
	manager interfaces.PluginManagerInterface,

) (*model.Fulfillment, *model_helper.AppError) {

	// begin transaction
	transaction := a.srv.Store.GetMaster().Begin()
	if transaction.Error != nil {
		return nil, model_helper.NewAppError("CreateReturnFulfillment", model_helper.ErrorCreatingTransactionErrorID, nil, transaction.Error.Error(), http.StatusInternalServerError)
	}
	defer a.srv.Store.FinalizeTransaction(transaction)

	status := model.FULFILLMENT_RETURNED
	if totalRefundAmount != nil {
		status = model.FULFILLMENT_REFUNDED_AND_RETURNED
	}

	returnFulfillment, appErr := a.moveLinesToReturnFulfillment(
		orderLineDatas,
		fulfillmentLineDatas,
		status,
		order,
		totalRefundAmount,
		shippingRefundAmount,
		manager,
	)
	if appErr != nil {
		return nil, appErr
	}

	returnedLines := map[string]*model.QuantityOrderLine{}

	orderLineIDs := []string{}
	for _, lineData := range fulfillmentLineDatas {
		orderLineIDs = append(orderLineIDs, lineData.Line.OrderLineID)
	}
	orderLinesByIDs, appErr := a.OrderLinesByOption(&model.OrderLineFilterOption{
		Conditions: squirrel.Eq{model.OrderLineTableName + ".Id": orderLineIDs},
	})
	if appErr != nil {
		return nil, appErr
	}

	orderLinesByIDsMap := lo.SliceToMap(orderLinesByIDs, func(o *model.OrderLine) (string, *model.OrderLine) {
		return o.Id, o
	})

	for _, orderLineData := range orderLineDatas {
		returnedLines[orderLineData.Line.Id] = &model.QuantityOrderLine{
			Quantity:  orderLineData.Quantity,
			OrderLine: &orderLineData.Line,
		}
	}

	for _, fulfillmentLineData := range fulfillmentLineDatas {
		if orderLine := orderLinesByIDsMap[fulfillmentLineData.Line.OrderLineID]; orderLine != nil {
			returnedLine := returnedLines[orderLine.Id]

			if returnedLine != nil {
				returnedLines[orderLine.Id] = &model.QuantityOrderLine{
					Quantity:  returnedLine.Quantity + fulfillmentLineData.Quantity,
					OrderLine: returnedLine.OrderLine,
				}
			} else {
				returnedLines[orderLine.Id] = &model.QuantityOrderLine{
					Quantity:  fulfillmentLineData.Quantity,
					OrderLine: orderLine,
				}
			}
		}
	}

	sliceOfQuantityOrderLine := []*model.QuantityOrderLine{}
	for _, value := range returnedLines {
		sliceOfQuantityOrderLine = append(sliceOfQuantityOrderLine, value)
	}

	// commit
	if err := transaction.Commit().Error; err != nil {
		return nil, model_helper.NewAppError("CreateReturnFulfillment", model_helper.ErrorCommittingTransactionErrorID, nil, err.Error(), http.StatusInternalServerError)
	}

	// NOTE: this is called after transaction commit
	appErr = a.OrderReturned(transaction, order, requester, nil, sliceOfQuantityOrderLine)
	if appErr != nil {
		return nil, appErr
	}
	return returnFulfillment, nil
}

// ProcessReplace Create replace fulfillment and new draft order.
//
// Move all requested lines to fulfillment with status replaced. Based on original
// order create the draft order with all user details, and requested lines.
func (a *ServiceOrder) ProcessReplace(
	requester *model.User,
	order model.Order,
	orderLineDatas []*model.OrderLineData,
	fulfillmentLineDatas []*model.FulfillmentLineData,
	manager interfaces.PluginManagerInterface,

) (*model.Fulfillment, *model.Order, *model_helper.AppError) {
	transaction := a.srv.Store.GetMaster().Begin()
	if transaction.Error != nil {
		return nil, nil, model_helper.NewAppError("ProcessReplace", model_helper.ErrorCreatingTransactionErrorID, nil, transaction.Error.Error(), http.StatusInternalServerError)
	}
	defer a.srv.Store.FinalizeTransaction(transaction)

	replaceFulfillment, appErr := a.moveLinesToReplaceFulfillment(orderLineDatas, fulfillmentLineDatas, order, manager)
	if appErr != nil {
		return nil, nil, appErr
	}

	newOrder, appErr := a.CreateReplaceOrder(requester, nil, order, orderLineDatas, fulfillmentLineDatas)
	if appErr != nil {
		return nil, nil, appErr
	}

	orderLinesOfOrder, appErr := a.OrderLinesByOption(&model.OrderLineFilterOption{
		Conditions: squirrel.Eq{model.OrderLineTableName + ".OrderID": newOrder.Id},
	})
	if appErr != nil {
		return nil, nil, appErr
	}

	replacedLines := []*model.QuantityOrderLine{}
	for _, orderLine := range orderLinesOfOrder {
		replacedLines = append(replacedLines, &model.QuantityOrderLine{
			Quantity:  orderLine.Quantity,
			OrderLine: orderLine,
		})
	}

	_, appErr = a.FulfillmentReplacedEvent(transaction, order, requester, nil, replacedLines)
	if appErr != nil {
		return nil, nil, appErr
	}

	_, appErr = a.OrderReplacementCreated(transaction, order, newOrder, requester, nil)
	if appErr != nil {
		return nil, nil, appErr
	}

	if err := transaction.Commit().Error; err != nil {
		return nil, nil, model_helper.NewAppError("ProcessReplace", model_helper.ErrorCommittingTransactionErrorID, nil, err.Error(), http.StatusInternalServerError)
	}

	return replaceFulfillment, newOrder, nil
}

// Process the request for replacing or returning the products.
//
// Process the refund when the refund is set to True. The amount of refund will be
// calculated for all lines with statuses different from refunded.  The lines which
// are set to replace will not be included in the refund amount.
//
// If the amount is provided, the refund will be used for this amount.
//
// If refund_shipping_costs is True, the calculated refund amount will include
// shipping costs.
//
// All lines with replace set to True will be used to create a new draft order, with
// the same order details as the original order.  These lines will be moved to
// fulfillment with status replaced. The events with relation to new order will be
// created.
//
// All lines with replace set to False will be moved to fulfillment with status
// returned/refunded_and_returned - depends on refund flag and current line status.
// If the fulfillment line has refunded status it will be moved to
// returned_and_refunded
//
// NOTE: `payMent`, `amount` , `user` are optional.
//
// `refund` and `refundShippingCosts` default to false.
func (a *ServiceOrder) CreateFulfillmentsForReturnedProducts(
	user *model.User,
	_ any,
	order model.Order,
	payMent *model.Payment,
	orderLineDatas []*model.OrderLineData,
	fulfillmentLineDatas []*model.FulfillmentLineData,
	manager interfaces.PluginManagerInterface,
	refund bool,
	amount *decimal.Decimal,
	refundShippingCosts bool,

) (*model.Fulfillment, *model.Fulfillment, *model.Order, *model_helper.PaymentError, *model_helper.AppError) {

	var (
		returnOrderLines        []*model.OrderLineData
		returnFulfillmentLines  []*model.FulfillmentLineData
		replaceOrderLines       []*model.OrderLineData
		replaceFulfillmentLines []*model.FulfillmentLineData
	)
	for _, lineData := range orderLineDatas {
		if !lineData.Replace {
			returnOrderLines = append(returnOrderLines, lineData)
			continue
		}
		replaceOrderLines = append(replaceOrderLines, lineData)
	}
	for _, lineData := range fulfillmentLineDatas {
		if !lineData.Replace {
			returnFulfillmentLines = append(returnFulfillmentLines, lineData)
			continue
		}
		replaceFulfillmentLines = append(replaceFulfillmentLines, lineData)
	}

	shippingRefundAmount := getShippingRefundAmount(refundShippingCosts, amount, order.ShippingPriceGrossAmount)

	// create transaction
	transaction := a.srv.Store.GetMaster().Begin()
	if transaction.Error != nil {
		return nil, nil, nil, nil, model_helper.NewAppError("CreateFulfillmentsForReturnedProducts", model_helper.ErrorCreatingTransactionErrorID, nil, transaction.Error.Error(), http.StatusInternalServerError)
	}
	defer a.srv.Store.FinalizeTransaction(transaction)

	var (
		totalRefundAmount *decimal.Decimal
		appErr            *model_helper.AppError
		paymentErr        *model_helper.PaymentError
	)
	if refund && payMent != nil {
		totalRefundAmount, paymentErr, appErr = a.processRefund(
			user,
			nil,
			order,
			*payMent,
			returnOrderLines,
			returnFulfillmentLines,
			amount,
			refundShippingCosts,
			manager,
		)
		if paymentErr != nil || appErr != nil {
			return nil, nil, nil, paymentErr, appErr
		}
	}

	var (
		replaceFulfillment *model.Fulfillment
		newOrder           *model.Order
	)
	if len(replaceFulfillmentLines) > 0 || len(replaceOrderLines) > 0 {
		replaceFulfillment, newOrder, appErr = a.ProcessReplace(
			user,
			order,
			replaceOrderLines,
			replaceFulfillmentLines,
			manager,
		)
		if appErr != nil {
			return nil, nil, nil, nil, appErr
		}
	}

	returnFulfillment, appErr := a.CreateReturnFulfillment(
		user,
		order,
		returnOrderLines,
		returnFulfillmentLines,
		totalRefundAmount,
		shippingRefundAmount,
		manager,
	)
	if appErr != nil {
		return nil, nil, nil, nil, appErr
	}

	fulfillmentsToDelete, appErr := a.FulfillmentsByOption(&model.FulfillmentFilterOption{
		Conditions: squirrel.Eq{
			model.FulfillmentTableName + ".OrderID": order.Id,
			model.FulfillmentTableName + ".Status": []model.FulfillmentStatus{
				model.FULFILLMENT_FULFILLED,
				model.FULFILLMENT_WAITING_FOR_APPROVAL,
			},
		},
		FulfillmentLineID: squirrel.Eq{model.FulfillmentLineTableName + ".Id": nil},
		Transaction:       transaction,
	})
	if appErr != nil && appErr.StatusCode == http.StatusInternalServerError { // ignore not found err
		return nil, nil, nil, nil, appErr
	}

	if len(fulfillmentsToDelete) > 0 {
		appErr = a.BulkDeleteFulfillments(transaction, fulfillmentsToDelete)
		if appErr != nil {
			return nil, nil, nil, nil, appErr
		}
	}

	// commit transaction
	if err := transaction.Commit().Error; err != nil {
		return nil, nil, nil, nil, model_helper.NewAppError("CreateFulfillmentsForReturnedProducts", model_helper.ErrorCommittingTransactionErrorID, nil, err.Error(), http.StatusInternalServerError)
	}

	_, appErr = manager.OrderUpdated(order)
	if appErr != nil {
		return nil, nil, nil, nil, appErr
	}

	return returnFulfillment, replaceFulfillment, newOrder, nil, nil
}

func (a *ServiceOrder) calculateRefundAmount(
	returnOrderLineDatas []*model.OrderLineData,
	returnFulfillmentLineDatas []*model.FulfillmentLineData,
	linesToRefund map[string]*model.QuantityOrderLine,

) (*decimal.Decimal, *model_helper.AppError) {

	refundAmount := decimal.NewFromInt(0)
	for _, lineData := range returnOrderLineDatas {
		if unitPriceGrossAmount := lineData.Line.UnitPriceGrossAmount; unitPriceGrossAmount != nil {
			refundAmount = refundAmount.Add(
				unitPriceGrossAmount.Mul(decimal.NewFromInt32(int32(lineData.Quantity))),
			)
		}
		linesToRefund[lineData.Line.Id] = &model.QuantityOrderLine{
			Quantity:  lineData.Quantity,
			OrderLine: &lineData.Line,
		}
	}

	if len(returnFulfillmentLineDatas) == 0 {
		return &refundAmount, nil
	}

	orderLineIDs := []string{}
	fulfillmentIDs := []string{}

	for _, lineData := range returnFulfillmentLineDatas {
		orderLineIDs = append(orderLineIDs, lineData.Line.OrderLineID)
		fulfillmentIDs = append(fulfillmentIDs, lineData.Line.FulfillmentID)
	}

	orderLines, appErr := a.OrderLinesByOption(&model.OrderLineFilterOption{
		Conditions: squirrel.Eq{model.OrderLineTableName + ".Id": orderLineIDs},
	})
	if appErr != nil {
		return nil, appErr
	}

	fulfillments, appErr := a.FulfillmentsByOption(&model.FulfillmentFilterOption{
		Conditions: squirrel.Eq{model.FulfillmentTableName + ".Id": fulfillmentIDs},
	})
	if appErr != nil {
		return nil, appErr
	}

	orderLinesMap := lo.SliceToMap(orderLines, func(o *model.OrderLine) (string, *model.OrderLine) { return o.Id, o })

	fulfillmentsMap := lo.SliceToMap(fulfillments, func(f *model.Fulfillment) (string, *model.Fulfillment) { return f.Id, f })

	for _, lineData := range returnFulfillmentLineDatas {
		// skip lines which were already refunded
		fulfillment := fulfillmentsMap[lineData.Line.FulfillmentID]
		if fulfillment != nil && fulfillment.Status == model.FULFILLMENT_REFUNDED {
			continue
		}

		if orderLine := orderLinesMap[lineData.Line.OrderLineID]; orderLine != nil {
			if unitPriceGrossAmount := orderLine.UnitPriceGrossAmount; unitPriceGrossAmount != nil {
				refundAmount = refundAmount.Add(
					unitPriceGrossAmount.Mul(decimal.NewFromInt32(int32(lineData.Quantity))),
				)
			}

			dataFromAllRefundedLines := linesToRefund[orderLine.Id]
			if dataFromAllRefundedLines != nil {
				linesToRefund[orderLine.Id] = &model.QuantityOrderLine{
					Quantity:  dataFromAllRefundedLines.Quantity + lineData.Quantity,
					OrderLine: dataFromAllRefundedLines.OrderLine,
				}
			} else {
				linesToRefund[orderLine.Id] = &model.QuantityOrderLine{
					Quantity:  lineData.Quantity,
					OrderLine: orderLine,
				}
			}
		}
	}

	return &refundAmount, nil
}

// `requester` and `amount` can be nil
func (a *ServiceOrder) processRefund(
	user *model.User,
	_ any,
	order model.Order,
	payMent model.Payment,
	orderLinesToRefund []*model.OrderLineData,
	fulfillmentLinesToRefund []*model.FulfillmentLineData,
	amount *decimal.Decimal,
	refundShippingCosts bool,
	manager interfaces.PluginManagerInterface,
) (*decimal.Decimal, *model_helper.PaymentError, *model_helper.AppError) {
	// transaction begin
	transaction := a.srv.Store.GetMaster().Begin()
	if transaction.Error != nil {
		return nil, nil, model_helper.NewAppError("processRefund", model_helper.ErrorCreatingTransactionErrorID, nil, transaction.Error.Error(), http.StatusInternalServerError)
	}
	defer a.srv.Store.FinalizeTransaction(transaction)

	linesToRefund := map[string]*model.QuantityOrderLine{}

	if amount == nil {
		var appErr *model_helper.AppError
		amount, appErr = a.calculateRefundAmount(orderLinesToRefund, fulfillmentLinesToRefund, linesToRefund)
		if appErr != nil {
			return nil, nil, appErr
		}
		// we take into consideration the shipping costs only when amount is not provided.
		if refundShippingCosts && order.ShippingPriceGrossAmount != nil {
			addResult := amount.Add(*order.ShippingPriceGrossAmount)
			amount = &addResult
		}
	}

	var (
		createPaymentRefundedEvent   bool
		sendOrderRefunddConfirmation bool
	)

	if amount != nil && !amount.Equal(decimal.Zero) {
		if payMent.CapturedAmount.LessThan(*amount) {
			amount = payMent.CapturedAmount
		}

		_, paymentErr, appErr := a.srv.PaymentService().Refund(transaction, payMent, manager, order.ChannelID, amount)
		if paymentErr != nil || appErr != nil {
			return nil, paymentErr, appErr
		}

		createPaymentRefundedEvent = true
		sendOrderRefunddConfirmation = true
	}

	// commit transaction
	if err := transaction.Commit().Error; err != nil {
		return nil, nil, model_helper.NewAppError("processRefund", model_helper.ErrorCommittingTransactionErrorID, nil, err.Error(), http.StatusInternalServerError)
	}

	var userID *string
	if user != nil && model_helper.IsValidId(user.Id) {
		userID = &user.Id
	}

	if createPaymentRefundedEvent {
		_, appErr := a.CommonCreateOrderEvent(nil, &model.OrderEventOption{
			OrderID:    order.Id,
			UserID:     userID,
			Type:       model.ORDER_EVENT_TYPE_PAYMENT_REFUNDED,
			Parameters: getPaymentData(amount, payMent),
		})
		if appErr != nil {
			return nil, nil, appErr
		}
	}

	if sendOrderRefunddConfirmation {
		appErr := a.SendOrderRefundedConfirmation(order, user, nil, *amount, payMent.Currency, manager)
		if appErr != nil {
			return nil, nil, appErr
		}
	}

	_, appErr := a.CommonCreateOrderEvent(nil, &model.OrderEventOption{
		OrderID: order.Id,
		Type:    model.ORDER_EVENT_TYPE_FULFILLMENT_REFUNDED,
		UserID:  userID,
		Parameters: model_types.JSONString{
			"lines":                   a.LinesPerQuantityToLineObjectList(lo.Values(linesToRefund)),
			"amount":                  amount,
			"shipping_costs_included": refundShippingCosts,
		},
	})
	if appErr != nil {
		return nil, nil, appErr
	}

	return amount, nil, nil
}

func (s *ServiceOrder) decreaseStocks(orderLinesInfo []*model.OrderLineData, manager interfaces.PluginManagerInterface, allowStockToBeExceeded bool) (*model_helper.InsufficientStock, *model_helper.AppError) {
	linesToDecreaseStock := s.srv.WarehouseService().GetOrderLinesWithTrackInventory(orderLinesInfo)
	if len(linesToDecreaseStock) > 0 {
		insufficientStock, appErr := s.srv.WarehouseService().DecreaseStock(linesToDecreaseStock, manager, true, allowStockToBeExceeded)
		if insufficientStock != nil || appErr != nil {
			return insufficientStock, appErr
		}
	}

	return nil, nil
}
