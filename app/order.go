package app

import (
	"net/http"

	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model/order"
	"github.com/sitename/sitename/model/payment"
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

// OrderShippingIsRequired checks if an order requires ship or not
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

	latestPayment := payments[0]
	for _, payment := range payments[1:] {
		if payment.CreateAt > latestPayment.CreateAt {
			latestPayment = payment
		}
	}

	return latestPayment, nil
}
