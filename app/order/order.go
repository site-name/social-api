/*
	NOTE: This package is initialized during server startup (modules/imports does that)
	so the init() function get the chance to register a function to create `ServiceAccount`
*/
package order

import (
	"context"
	"net/http"

	"github.com/mattermost/gorp"
	"github.com/site-name/decimal"
	goprices "github.com/site-name/go-prices"
	"github.com/sitename/sitename/app"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model/account"
	"github.com/sitename/sitename/model/order"
	"github.com/sitename/sitename/model/payment"
	"github.com/sitename/sitename/modules/util"
	"github.com/sitename/sitename/store"
)

// UpsertOrder depends on given order's Id property to decide update/save it
func (a *ServiceOrder) UpsertOrder(transaction *gorp.Transaction, ord *order.Order) (*order.Order, *model.AppError) {
	var err error

	if ord.Id == "" {
		ord, err = a.srv.Store.Order().Save(transaction, ord)
	} else {
		ord, err = a.srv.Store.Order().Update(transaction, ord)
	}

	if err != nil {
		if appErr, ok := err.(*model.AppError); ok {
			return nil, appErr
		}
		if _, ok := err.(*store.ErrInvalidInput); ok {
			return nil, model.NewAppError("UpsertOrder", app.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "order"}, err.Error(), http.StatusBadRequest)
		}

		return nil, model.NewAppError("UpsertOrder", "app.order.error_upserting_order.app_error", nil, err.Error(), http.StatusInternalServerError)
	}

	return ord, nil
}

// BulkUpsertOrders performs bulk upsert given orders
func (a *ServiceOrder) BulkUpsertOrders(orders []*order.Order) ([]*order.Order, *model.AppError) {
	orders, err := a.srv.Store.Order().BulkUpsert(orders)
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
func (a *ServiceOrder) FilterOrdersByOptions(option *order.OrderFilterOption) ([]*order.Order, *model.AppError) {
	orders, err := a.srv.Store.Order().FilterByOption(option)
	if err != nil {
		return nil, store.AppErrorFromDatabaseLookupError("FilterOrdersbyOption", "app.order.error_finding_orders_by_option.app_error", err)
	}

	return orders, nil
}

// OrderById retuns an order with given id
func (a *ServiceOrder) OrderById(id string) (*order.Order, *model.AppError) {
	order, err := a.srv.Store.Order().Get(id)
	if err != nil {
		return nil, store.AppErrorFromDatabaseLookupError("OrderById", "app.order.order_missing.app_error", err)
	}

	return order, nil
}

// OrderShippingIsRequired returns a boolean value indicating that given order requires shipping or not
func (a *ServiceOrder) OrderShippingIsRequired(orderID string) (bool, *model.AppError) {
	lines, appErr := a.OrderLinesByOption(&order.OrderLineFilterOption{
		OrderID: &model.StringFilter{
			StringOption: &model.StringOption{
				Eq: orderID,
			},
		},
	})
	if appErr != nil {
		appErr.Where = "OrderShippingIsRequired"
		return false, appErr
	}

	for _, line := range lines {
		if line.IsShippingRequired {
			return true, nil
		}
	}

	return false, nil
}

// OrderTotalQuantity return total quantity of given order
func (a *ServiceOrder) OrderTotalQuantity(orderID string) (int, *model.AppError) {
	lines, appErr := a.OrderLinesByOption(&order.OrderLineFilterOption{
		OrderID: &model.StringFilter{
			StringOption: &model.StringOption{
				Eq: orderID,
			},
		},
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
func (a *ServiceOrder) UpdateOrderTotalPaid(transaction *gorp.Transaction, orderID string) *model.AppError {
	order, appErr := a.OrderById(orderID)
	if appErr != nil {
		return appErr
	}
	payments, appErr := a.srv.PaymentService().PaymentsByOption(&payment.PaymentFilterOption{
		OrderID: orderID,
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

	order.TotalPaidAmount = &total

	_, appErr = a.UpsertOrder(transaction, order)
	if appErr != nil {
		return appErr
	}

	return nil
}

// OrderIsPreAuthorized checks if order is pre-authorized
func (a *ServiceOrder) OrderIsPreAuthorized(orderID string) (bool, *model.AppError) {
	payments, appErr := a.srv.PaymentService().PaymentsByOption(&payment.PaymentFilterOption{
		OrderID:  orderID,
		IsActive: model.NewBool(true),
		TransactionsKind: &model.StringFilter{
			StringOption: &model.StringOption{
				Eq: payment.AUTH,
			},
		},
		TransactionsActionRequired: model.NewBool(false),
		TransactionsIsSuccess:      model.NewBool(true),
	})
	if appErr != nil {
		return false, appErr
	}

	return len(payments) > 0, nil
}

// OrderIsCaptured checks if given order is captured
func (a *ServiceOrder) OrderIsCaptured(orderID string) (bool, *model.AppError) {
	payments, appErr := a.srv.PaymentService().PaymentsByOption(&payment.PaymentFilterOption{
		OrderID:  orderID,
		IsActive: model.NewBool(true),
		TransactionsKind: &model.StringFilter{
			StringOption: &model.StringOption{
				Eq: payment.CAPTURE,
			},
		},
		TransactionsActionRequired: model.NewBool(false),
		TransactionsIsSuccess:      model.NewBool(true),
	})
	if appErr != nil {
		return false, appErr
	}

	return len(payments) > 0, nil
}

// OrderSubTotal returns sum of TotalPrice of all order lines that belong to given order
func (a *ServiceOrder) OrderSubTotal(ord *order.Order) (*goprices.TaxedMoney, *model.AppError) {
	lines, appErr := a.OrderLinesByOption(&order.OrderLineFilterOption{
		OrderID: &model.StringFilter{
			StringOption: &model.StringOption{
				Eq: ord.Id,
			},
		},
	})
	if appErr != nil {
		appErr.Where = "OrderSubTotal"
		return nil, appErr
	}

	return a.srv.PaymentService().GetSubTotal(lines, ord.Currency)
}

// OrderCanCalcel checks if given order can be canceled
func (a *ServiceOrder) OrderCanCancel(ord *order.Order) (bool, *model.AppError) {
	fulfillments, err := a.FulfillmentsByOption(nil, &order.FulfillmentFilterOption{
		OrderID: &model.StringFilter{
			StringOption: &model.StringOption{
				Eq: ord.Id,
			},
		},
		Status: &model.StringFilter{
			StringOption: &model.StringOption{
				NotIn: []string{
					order.FULFILLMENT_CANCELED,
					order.FULFILLMENT_REFUNDED,
					order.FULFILLMENT_RETURNED,
					order.FULFILLMENT_REFUNDED_AND_RETURNED,
					order.FULFILLMENT_REPLACED,
				},
			},
		},
	})

	if err != nil {
		// this means system error
		return false, model.NewAppError("OrderCanCancel", "app.order.fulfillments_by_option.app_error", nil, err.Error(), http.StatusInternalServerError)
	}

	return len(fulfillments) == 0 && ord.Status != order.CANCELED && ord.Status != order.DRAFT, nil
}

// OrderCanCapture
func (a *ServiceOrder) OrderCanCapture(ord *order.Order, payment *payment.Payment) (bool, *model.AppError) {
	var err *model.AppError

	if payment == nil {
		payment, err = a.srv.PaymentService().GetLastOrderPayment(ord.Id)
		if err != nil {
			return false, err
		}
	}

	if payment == nil {
		return false, nil
	}

	return payment.CanCapture() &&
		ord.Status != order.DRAFT &&
		ord.Status != order.CANCELED, nil
}

// OrderCanVoid
func (a *ServiceOrder) OrderCanVoid(ord *order.Order, payment *payment.Payment) (bool, *model.AppError) {
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
func (a *ServiceOrder) OrderCanRefund(ord *order.Order, payment *payment.Payment) (bool, *model.AppError) {
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
func (a *ServiceOrder) CanMarkOrderAsPaid(ord *order.Order, payments []*payment.Payment) (bool, *model.AppError) {
	var appErr *model.AppError
	if len(payments) == 0 {
		payments, appErr = a.srv.PaymentService().PaymentsByOption(&payment.PaymentFilterOption{
			OrderID: ord.Id,
		})
	}

	if appErr != nil {
		return false, appErr
	}

	return len(payments) == 0, nil
}

// OrderTotalAuthorized returns order's total authorized amount
func (a *ServiceOrder) OrderTotalAuthorized(ord *order.Order) (*goprices.Money, *model.AppError) {
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

// GetOrderCountryCode is helper function, returns contry code of given order
func (a *ServiceOrder) GetOrderCountryCode(ord *order.Order) (string, *model.AppError) {
	addressID := ord.BillingAddressID
	requireShipping, appErr := a.OrderShippingIsRequired(ord.Id)
	if appErr != nil {
		return "", appErr
	}

	if requireShipping {
		addressID = ord.ShippingAddressID
	}
	if addressID == nil {
		return *a.srv.Config().LocalizationSettings.DefaultCountryCode, nil
	}

	address, err := a.srv.Store.Address().Get(*addressID)
	if err != nil {

		var statusCode int = http.StatusInternalServerError
		if _, ok := err.(*store.ErrNotFound); ok {
			statusCode = http.StatusNotFound
		}
		return "", model.NewAppError("GetOrderCountryCode", "app.order.get_address.app_error", nil, err.Error(), statusCode)
	}

	return address.Country, nil
}

// CustomerEmail try finding order's owner's email. If order has no user or error occured during the finding process, returns order's UserEmail property instead
func (a *ServiceOrder) CustomerEmail(ord *order.Order) (string, *model.AppError) {
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
func (a *ServiceOrder) AnAddressOfOrder(orderID string, whichAddressID order.WhichOrderAddressID) (*account.Address, *model.AppError) {
	addresses, appErr := a.srv.AccountService().AddressesByOption(&account.AddressFilterOption{
		OrderID: &account.AddressFilterOrderOption{
			Id: &model.StringFilter{
				StringOption: &model.StringOption{
					Eq: orderID,
				},
			},
			On: string(whichAddressID),
		},
	})
	if appErr != nil {
		return nil, appErr
	}

	return addresses[0], nil
}
