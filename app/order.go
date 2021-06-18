package app

import (
	"errors"
	"net/http"
	"strings"

	"github.com/shopspring/decimal"
	goprices "github.com/site-name/go-prices"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model/order"
	"github.com/sitename/sitename/model/payment"
	"github.com/sitename/sitename/modules/util"
	"github.com/sitename/sitename/store"
)

// GetAllOrderLinesByOrderId returns a slice of order lines that belong to given order
func (a *App) GetAllOrderLinesByOrderId(orderID string) ([]*order.OrderLine, *model.AppError) {
	lines, err := a.srv.Store.OrderLine().GetAllByOrderID(orderID)
	if err != nil {
		var statusCode int
		switch err.(type) {
		case *store.ErrNotFound:
			statusCode = http.StatusNotFound
		default:
			statusCode = http.StatusInternalServerError
		}
		return nil, model.NewAppError("GetAllOrderLinesByOrderId", "app.order.get_child_order_lines.app_error", nil, err.Error(), statusCode)
	}

	return lines, nil
}

// GetAllPaymentsByOrderId returns all payments that belong to order with given orderID
func (a *App) GetAllPaymentsByOrderId(orderID string) ([]*payment.Payment, *model.AppError) {
	payments, err := a.srv.Store.Payment().GetPaymentsByOrderID(orderID)
	if err != nil {
		var statusCode int
		switch err.(type) {
		case *store.ErrNotFound:
			statusCode = http.StatusNotFound
		default:
			statusCode = http.StatusInternalServerError
		}
		return nil, model.NewAppError("GetAllPaymentsByOrderId", "app.order.get_child_payments.app_error", nil, err.Error(), statusCode)
	}

	return payments, nil
}

// OrderShippingIsRequired checks if an order requires ship or not by:
//
// 1) Find all child order lines that belong to given order
//
// 2) iterates over resulting slice to check if at least one order line requires shipping
func (a *App) OrderShippingIsRequired(orderID string) (bool, *model.AppError) {
	lines, err := a.GetAllOrderLinesByOrderId(orderID)
	if err != nil {
		return false, err
	}

	for _, line := range lines {
		if line.IsShippingRequired {
			return true, nil
		}
	}

	return false, nil
}

// OrderTotalQuantity return total quantity of given order
func (a *App) OrderTotalQuantity(orderID string) (int, *model.AppError) {
	lines, err := a.GetAllOrderLinesByOrderId(orderID)
	if err != nil {
		return 0, err
	}

	var total int = 0
	for _, line := range lines {
		total += line.Quantity
	}

	return total, nil
}

// GetLastOrderPayment get most recent payment made for given order
func (a *App) GetLastOrderPayment(orderID string) (*payment.Payment, *model.AppError) {
	payments, err := a.GetAllPaymentsByOrderId(orderID)
	if err != nil {
		return nil, err
	}

	if len(payments) == 0 {
		return nil, nil
	}

	if len(payments) == 1 {
		return payments[0], nil
	}

	latestPayment := payments[0]
	for _, payment := range payments[1:] {
		if payment != nil && payment.CreateAt >= latestPayment.CreateAt {
			latestPayment = payment
		}
	}

	return latestPayment, nil
}

// UpdateOrderTotalPaid update given order's total paid amount
func (a *App) UpdateOrderTotalPaid(orderID string) *model.AppError {
	payments, appErr := a.GetAllPaymentsByOrderId(orderID)
	if appErr != nil {
		return appErr
	}

	total := decimal.Zero
	for _, payment := range payments {
		if payment.CapturedAmount != nil {
			total = total.Add(*payment.CapturedAmount)
		}
	}

	if err := a.srv.Store.Order().UpdateTotalPaid(orderID, &total); err != nil {
		return model.NewAppError("UpdateOrderTotalPaid", "app.order.update_total_paid.app_error", nil, err.Error(), http.StatusInternalServerError)
	}

	return nil
}

// OrderIsPreAuthorized checks if order is pre-authorized
func (a *App) OrderIsPreAuthorized(orderID string) (bool, *model.AppError) {
	filterOptions := &payment.PaymentFilterOpts{
		OrderID:  orderID,
		IsActive: true,
		PaymentTransactionFilterOpts: payment.PaymentTransactionFilterOpts{
			Kind:           payment.AUTH,
			ActionRequired: false,
			IsSuccess:      true,
		},
	}
	exist, err := a.srv.Store.Payment().PaymentExistWithOptions(filterOptions)
	if err != nil {
		// this err means system error, not sql not found
		return false, model.NewAppError("OrderIsPreAuthorized", "app.order.order_is_pre_authorized.app_error", nil, err.Error(), http.StatusInternalServerError)
	}

	return exist, nil
}

// OrderIsCaptured checks if given order is captured
func (a *App) OrderIsCaptured(orderID string) (bool, *model.AppError) {
	filterOptions := &payment.PaymentFilterOpts{
		OrderID:  orderID,
		IsActive: true,
		PaymentTransactionFilterOpts: payment.PaymentTransactionFilterOpts{
			Kind:           payment.CAPTURE,
			ActionRequired: false,
			IsSuccess:      true,
		},
	}
	exist, err := a.srv.Store.Payment().PaymentExistWithOptions(filterOptions)
	if err != nil {
		// this err means system error, not sql not found error
		return false, model.NewAppError("OrderIsCaptured", "app.order.order_is_captured.app_error", nil, err.Error(), http.StatusInternalServerError)
	}

	return exist, nil
}

// OrderSubTotal returns sum of TotalPrice of all order lines that belong to given order
func (a *App) OrderSubTotal(orderID string, orderCurrency string) (*goprices.TaxedMoney, *model.AppError) {
	orderLines, appErr := a.GetAllOrderLinesByOrderId(orderID)
	if appErr != nil {
		return nil, appErr
	}

	// check if order and its child order lines have  same currencies
	if !strings.EqualFold(orderCurrency, orderLines[0].Currency) {
		return nil, model.NewAppError("OrderSubTotal", "app.order.currency_integrity.app_error", nil, "orders and order lines must have same currencies", http.StatusInternalServerError)
	}

	subTotal, err := util.ZeroTaxedMoney(orderLines[0].Currency)
	if err != nil {
		return nil, model.NewAppError("OrderSubTotal", "app.order.get_sub_total.app_error", nil, err.Error(), http.StatusInternalServerError)
	}

	for _, line := range orderLines {
		if line.TotalPrice != nil {
			subTotal, err = subTotal.Add(line.TotalPrice)
			if err != nil {
				return nil, model.NewAppError("OrderSubTotal", "app.order.get_sub_total.app_error", nil, err.Error(), http.StatusInternalServerError)
			}
		}
	}

	return subTotal, nil
}

// OrderCanCalcel checks if given order can be canceled
func (a *App) OrderCanCancel(ord *order.Order) (bool, *model.AppError) {
	exist, err := a.srv.Store.Fulfillment().FilterByExcludeStatuses(ord.Id, []string{
		order.FULFILLMENT_CANCELED,
		order.FULFILLMENT_REFUNDED,
		order.FULFILLMENT_RETURNED,
		order.FULFILLMENT_REFUNDED_AND_RETURNED,
		order.FULFILLMENT_REPLACED,
	})
	if err != nil {
		// this means system error
		return false, model.NewAppError("OrderCanCancel", "app.order.fulfillment_exist.app_error", nil, err.Error(), http.StatusInternalServerError)
	}

	return !exist && ord.Status != order.CANCELED && ord.Status != order.DRAFT, nil
}

// OrderCanCapture
func (a *App) OrderCanCapture(ord *order.Order, payment *payment.Payment) (bool, *model.AppError) {
	var err *model.AppError

	if payment == nil {
		payment, err = a.GetLastOrderPayment(ord.Id)
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
func (a *App) OrderCanVoid(ord *order.Order, payment *payment.Payment) (bool, *model.AppError) {
	var err *model.AppError
	if payment == nil {
		payment, err = a.GetLastOrderPayment(ord.Id)
		if err != nil {
			return false, err
		}
	}

	if payment == nil {
		return false, nil
	}

	return a.PaymentCanVoid(payment)
}

// OrderCanRefund checks if order can refund
func (a *App) OrderCanRefund(ord *order.Order, payments []*payment.Payment) (bool, *model.AppError) {
	var appErr *model.AppError
	if len(payments) == 0 {
		payments, appErr = a.GetAllPaymentsByOrderId(ord.Id)
	}

	if appErr != nil {
		if appErr.StatusCode == http.StatusNotFound {
			// this means order has no payments yet
			return true, nil
		} else {
			// other errors mean system error
			return false, appErr
		}
	}

	return len(payments) == 0, nil
}

// CanMarkOrderAsPaid checks if given order can be marked as paid.
func (a *App) CanMarkOrderAsPaid(ord *order.Order, payments []*payment.Payment) (bool, *model.AppError) {
	var appErr *model.AppError
	if len(payments) == 0 {
		payments, appErr = a.GetAllPaymentsByOrderId(ord.Id)
	}

	if appErr != nil {
		if appErr.StatusCode == http.StatusNotFound {
			return true, nil
		} else {
			return false, appErr
		}
	}

	return len(payments) == 0, nil
}

// OrderTotalAuthorized returns order's total authorized amount
func (a *App) OrderTotalAuthorized(ord *order.Order) (*goprices.Money, *model.AppError) {
	lastPayment, appErr := a.GetLastOrderPayment(ord.Id)
	if appErr != nil {
		return nil, appErr
	}
	if lastPayment != nil && lastPayment.IsActive {
		return a.PaymentGetAuthorizedAmount(lastPayment)
	}

	zeroMoney, err := util.ZeroMoney(ord.Currency)
	if err != nil {
		return nil, model.NewAppError("OrderTotalAuthorized", "app.order.create_zero_money.app_error", nil, err.Error(), http.StatusInternalServerError)
	}
	return zeroMoney, nil
}

// GetOrderCountryCode is helper function, returns contry code of given order
func (a *App) GetOrderCountryCode(ord *order.Order) (string, *model.AppError) {
	addressID := ord.BillingAddressID
	requireShipping, appErr := a.OrderShippingIsRequired(ord.Id)
	if appErr != nil {
		return "", appErr
	}

	if requireShipping {
		addressID = ord.ShippingAddressID
	}
	if addressID == nil {
		return *a.Config().LocalizationSettings.DefaultCountryCode, nil
	}

	address, err := a.srv.Store.Address().Get(*addressID)
	if err != nil {
		var errNf *store.ErrNotFound
		var statusCode int
		if errors.As(err, &errNf) {
			statusCode = http.StatusNotFound
		} else {
			statusCode = http.StatusInternalServerError
		}
		return "", model.NewAppError("GetOrderCountryCode", "app.order.get_address.app_error", nil, errNf.Error(), statusCode)
	}

	return address.Country, nil
}
