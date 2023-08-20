/*
NOTE: This package is initialized during server startup (modules/imports does that)
so the init() function get the chance to register a function to create `ServiceAccount`
*/
package order

import (
	"context"
	"net/http"
	"strings"

	"github.com/Masterminds/squirrel"
	"github.com/samber/lo"
	"github.com/site-name/decimal"
	goprices "github.com/site-name/go-prices"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/modules/util"
	"github.com/sitename/sitename/store"
	"gorm.io/gorm"
)

// UpsertOrder depends on given order's Id property to decide update/save it
func (a *ServiceOrder) UpsertOrder(transaction *gorm.DB, order *model.Order) (*model.Order, *model.AppError) {
	orders, err := a.srv.Store.Order().BulkUpsert(transaction, []*model.Order{order})

	if err != nil {
		if appErr, ok := err.(*model.AppError); ok {
			return nil, appErr
		}
		if _, ok := err.(*store.ErrInvalidInput); ok {
			return nil, model.NewAppError("UpsertOrder", model.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "order"}, err.Error(), http.StatusBadRequest)
		}

		return nil, model.NewAppError("UpsertOrder", "app.order.error_upserting_order.app_error", nil, err.Error(), http.StatusInternalServerError)
	}

	return orders[0], nil
}

// BulkUpsertOrders performs bulk upsert given orders
func (a *ServiceOrder) BulkUpsertOrders(transaction *gorm.DB, orders []*model.Order) ([]*model.Order, *model.AppError) {
	orders, err := a.srv.Store.Order().BulkUpsert(transaction, orders)
	if err != nil {
		if appErr, ok := err.(*model.AppError); ok { // error caused by IsValid()
			return nil, appErr
		}
		statusCode := http.StatusInternalServerError
		if _, ok := err.(*store.ErrNotFound); ok { // error caused by SelectOne()
			statusCode = http.StatusNotFound
		}

		return nil, model.NewAppError("BulkUpsertOrders", "app.order.error_bulk_upsert_orders.app_error", nil, err.Error(), statusCode)
	}

	return orders, nil
}

// FilterOrdersByOptions is common method for filtering orders by given option
func (a *ServiceOrder) FilterOrdersByOptions(option *model.OrderFilterOption) ([]*model.Order, *model.AppError) {
	orders, err := a.srv.Store.Order().FilterByOption(option)
	if err != nil {
		return nil, model.NewAppError("FilterOrdersbyOption", "app.order.error_finding_orders_by_option.app_error", nil, err.Error(), http.StatusInternalServerError)
	}

	return orders, nil
}

// OrderById retuns an order with given id
func (a *ServiceOrder) OrderById(id string) (*model.Order, *model.AppError) {
	order, err := a.srv.Store.Order().Get(id)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if _, ok := err.(*store.ErrNotFound); ok {
			statusCode = http.StatusNotFound
		}
		return nil, model.NewAppError("OrderById", "app.order.order_missing.app_error", nil, err.Error(), statusCode)
	}

	return order, nil
}

// OrderShippingIsRequired returns a boolean value indicating that given order requires shipping or not
func (a *ServiceOrder) OrderShippingIsRequired(orderID string) (bool, *model.AppError) {
	lines, appErr := a.OrderLinesByOption(&model.OrderLineFilterOption{
		Conditions: squirrel.Eq{model.OrderLineTableName + ".OrderID": orderID},
	})
	if appErr != nil {
		return false, appErr
	}

	return lo.SomeBy(lines, func(o *model.OrderLine) bool { return o != nil && o.IsShippingRequired }), nil
}

// OrderTotalQuantity return total quantity of given order
func (a *ServiceOrder) OrderTotalQuantity(orderID string) (int, *model.AppError) {
	lines, appErr := a.OrderLinesByOption(&model.OrderLineFilterOption{
		Conditions: squirrel.Eq{model.OrderLineTableName + ".OrderID": orderID},
	})
	if appErr != nil {
		return 0, appErr
	}

	var total int = 0
	for _, line := range lines {
		total += line.Quantity
	}

	return total, nil
}

// UpdateOrderTotalPaid update given order's total paid amount
func (a *ServiceOrder) UpdateOrderTotalPaid(transaction *gorm.DB, orDer *model.Order) *model.AppError {
	payments, appErr := a.srv.PaymentService().PaymentsByOption(&model.PaymentFilterOption{
		Conditions: squirrel.Eq{model.PaymentTableName + ".OrderID": orDer.Id},
	})
	if appErr != nil {
		return appErr
	}

	total := decimal.Zero
	for _, payment := range payments {
		if payment.CapturedAmount != nil {
			total = total.Add(*payment.CapturedAmount)
		}
	}

	orDer.TotalPaidAmount = &total

	_, appErr = a.UpsertOrder(transaction, orDer)
	if appErr != nil {
		return appErr
	}

	return nil
}

// OrderIsPreAuthorized checks if order is pre-authorized
func (a *ServiceOrder) OrderIsPreAuthorized(orderID string) (bool, *model.AppError) {
	payments, appErr := a.srv.PaymentService().PaymentsByOption(&model.PaymentFilterOption{
		TransactionsKind:           squirrel.Eq{model.TransactionTableName + ".Kind": model.AUTH},
		TransactionsActionRequired: squirrel.Eq{model.TransactionTableName + ".ActionRequired": false},
		TransactionsIsSuccess:      squirrel.Eq{model.TransactionTableName + ".IsSuccess": true},

		Conditions: squirrel.Eq{
			model.PaymentTableName + ".OrderID":  orderID,
			model.PaymentTableName + ".IsActive": true,
		},
	})
	if appErr != nil {
		return false, appErr
	}

	return len(payments) > 0, nil
}

// OrderIsCaptured checks if given order is captured
func (a *ServiceOrder) OrderIsCaptured(orderID string) (bool, *model.AppError) {
	payments, appErr := a.srv.PaymentService().PaymentsByOption(&model.PaymentFilterOption{
		TransactionsKind:           squirrel.Eq{model.TransactionTableName + ".Kind": model.CAPTURE},
		TransactionsActionRequired: squirrel.Eq{model.TransactionTableName + ".ActionRequired": false},
		TransactionsIsSuccess:      squirrel.Eq{model.TransactionTableName + ".IsSuccess": true},

		Conditions: squirrel.Eq{
			model.PaymentTableName + ".OrderID":  orderID,
			model.PaymentTableName + ".IsActive": true,
		},
	})
	if appErr != nil {
		return false, appErr
	}

	return len(payments) > 0, nil
}

// OrderSubTotal returns sum of TotalPrice of all order lines that belong to given order
func (a *ServiceOrder) OrderSubTotal(ord *model.Order) (*goprices.TaxedMoney, *model.AppError) {
	lines, appErr := a.OrderLinesByOption(&model.OrderLineFilterOption{
		Conditions: squirrel.Eq{model.OrderLineTableName + ".OrderID": ord.Id},
	})
	if appErr != nil {
		return nil, appErr
	}

	return a.srv.PaymentService().GetSubTotal(lines, ord.Currency)
}

// OrderCanCalcel checks if given order can be canceled
func (a *ServiceOrder) OrderCanCancel(ord *model.Order) (bool, *model.AppError) {
	fulfillments, err := a.FulfillmentsByOption(&model.FulfillmentFilterOption{
		Conditions: squirrel.And{
			squirrel.Eq{model.FulfillmentTableName + ".OrderID": ord.Id},
			squirrel.NotEq{model.FulfillmentTableName + ".Status": []model.FulfillmentStatus{
				model.FULFILLMENT_CANCELED,
				model.FULFILLMENT_REFUNDED,
				model.FULFILLMENT_RETURNED,
				model.FULFILLMENT_REFUNDED_AND_RETURNED,
				model.FULFILLMENT_REPLACED,
			}},
		},
	})

	if err != nil {
		// this means system error
		return false, model.NewAppError("OrderCanCancel", "app.order.fulfillments_by_option.app_error", nil, err.Error(), http.StatusInternalServerError)
	}

	return len(fulfillments) == 0 && ord.Status != model.ORDER_STATUS_CANCELED && ord.Status != model.ORDER_STATUS_DRAFT, nil
}

// OrderCanCapture
func (a *ServiceOrder) OrderCanCapture(ord *model.Order, payment *model.Payment) (bool, *model.AppError) {
	if payment == nil {
		var appErr *model.AppError
		payment, appErr = a.srv.PaymentService().GetLastOrderPayment(ord.Id)
		if appErr != nil {
			return false, appErr
		}
	}

	if payment == nil {
		return false, nil
	}

	return payment.CanCapture() &&
		ord.Status != model.ORDER_STATUS_DRAFT &&
		ord.Status != model.ORDER_STATUS_CANCELED, nil
}

// OrderCanVoid
func (a *ServiceOrder) OrderCanVoid(ord *model.Order, payment *model.Payment) (bool, *model.AppError) {
	var appErr *model.AppError
	if payment == nil {
		payment, appErr = a.srv.PaymentService().GetLastOrderPayment(ord.Id)
	}

	if appErr != nil {
		return false, appErr
	}

	if payment == nil {
		return false, nil
	}

	return a.srv.PaymentService().PaymentCanVoid(payment)
}

// OrderCanRefund checks if order can refund
func (a *ServiceOrder) OrderCanRefund(ord *model.Order, payment *model.Payment) (bool, *model.AppError) {
	var appErr *model.AppError
	if payment == nil {
		payment, appErr = a.srv.PaymentService().GetLastOrderPayment(ord.Id)
	}

	if appErr != nil {
		if appErr.StatusCode == http.StatusNotFound {
			// this means order has no payments yet
			return false, nil
		}
		return false, appErr
	}

	if payment == nil {
		return false, nil
	}

	return payment.CanRefund(), nil
}

// CanMarkOrderAsPaid checks if given order can be marked as paid.
func (a *ServiceOrder) CanMarkOrderAsPaid(ord *model.Order, payments []*model.Payment) (bool, *model.AppError) {
	var appErr *model.AppError
	if len(payments) == 0 {
		payments, appErr = a.srv.PaymentService().PaymentsByOption(&model.PaymentFilterOption{
			Conditions: squirrel.Eq{model.PaymentTableName + ".OrderID": ord.Id},
		})
	}

	if appErr != nil {
		return false, appErr
	}

	return len(payments) == 0, nil
}

// OrderTotalAuthorized returns order's total authorized amount
func (a *ServiceOrder) OrderTotalAuthorized(ord *model.Order) (*goprices.Money, *model.AppError) {
	lastPayment, appErr := a.srv.PaymentService().GetLastOrderPayment(ord.Id)
	if appErr != nil {
		return nil, appErr
	}
	if lastPayment != nil && *lastPayment.IsActive {
		return a.srv.PaymentService().PaymentGetAuthorizedAmount(lastPayment)
	}

	zeroMoney, _ := util.ZeroMoney(ord.Currency)
	return zeroMoney, nil
}

// CustomerEmail try finding order's owner's email. If order has no user or error occured during the finding process, returns order's UserEmail property instead
func (a *ServiceOrder) CustomerEmail(ord *model.Order) (string, *model.AppError) {
	if ord.UserID != nil {
		user, appErr := a.srv.AccountService().UserById(context.Background(), *ord.UserID)
		if appErr != nil {
			if appErr.StatusCode == http.StatusInternalServerError {
				return "", appErr
			}
			if ord.UserEmail != "" {
				return ord.UserEmail, nil
			}
		}
		return user.Email, nil
	}

	return ord.UserEmail, nil
}

// AnAddressOfOrder returns shipping address of given order if presents
func (a *ServiceOrder) AnAddressOfOrder(orderID string, whichAddressID model.WhichOrderAddressID) (*model.Address, *model.AppError) {
	addresses, appErr := a.srv.AccountService().AddressesByOption(&model.AddressFilterOption{
		OrderID: &model.AddressFilterOrderOption{
			Id: squirrel.Eq{model.OrderTableName + ".Id": orderID},
			On: whichAddressID,
		},
	})
	if appErr != nil {
		return nil, appErr
	}

	return addresses[0], nil
}

func (a *ServiceOrder) PrepareInsufficientStockOrderValidationAppError(where string, err *model.InsufficientStock) *model.AppError {
	var warehouseIDs = make([]string, 0, len(err.Items))
	var orderLineIDs = make([]string, 0, len(err.Items))
	for _, item := range err.Items {
		if item.WarehouseID != nil {
			warehouseIDs = append(warehouseIDs, *item.WarehouseID)
		}
		if item.OrderLine != nil {
			orderLineIDs = append(orderLineIDs, item.OrderLine.Id)
		}
	}

	return model.NewAppError(where, "app.order.insufficient_stock.app_error", map[string]interface{}{"orderLines": strings.Join(orderLineIDs, ", "), "warehouses": strings.Join(warehouseIDs, ", ")}, "insufficient product stock", http.StatusNotAcceptable)
}

func (s *ServiceOrder) DeleteOrders(transaction *gorm.DB, ids []string) (int64, *model.AppError) {
	count, err := s.srv.Store.Order().Delete(transaction, ids)
	if err != nil {
		return 0, model.NewAppError("DeleteOrders", "app.order.delete_orders.app_error", nil, err.Error(), http.StatusInternalServerError)
	}

	return count, nil
}
