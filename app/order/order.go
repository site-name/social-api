package order

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"github.com/site-name/decimal"
	goprices "github.com/site-name/go-prices"
	"github.com/sitename/sitename/app"
	"github.com/sitename/sitename/app/sub_app_iface"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model/order"
	"github.com/sitename/sitename/model/payment"
	"github.com/sitename/sitename/modules/util"
	"github.com/sitename/sitename/store"
)

type AppOrder struct {
	app.AppIface
}

func init() {
	app.RegisterOrderApp(func(a app.AppIface) sub_app_iface.OrderApp {
		return &AppOrder{a}
	})
}

func (a *AppOrder) OrderById(id string) (*order.Order, *model.AppError) {
	order, err := a.Srv().Store.Order().Get(id)
	if err != nil {
		return nil, store.AppErrorFromDatabaseLookupError("OrderById", "app.order.order_missing.app_error", err)
	}

	return order, nil
}

func (a *AppOrder) OrderShippingIsRequired(orderID string) (bool, *model.AppError) {
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
func (a *AppOrder) OrderTotalQuantity(orderID string) (uint, *model.AppError) {
	lines, err := a.GetAllOrderLinesByOrderId(orderID)
	if err != nil {
		return 0, err
	}

	var total uint = 0
	for _, line := range lines {
		total += line.Quantity
	}

	return total, nil
}

// UpdateOrderTotalPaid update given order's total paid amount
func (a *AppOrder) UpdateOrderTotalPaid(orderID string) *model.AppError {
	payments, appErr := a.PaymentApp().GetAllPaymentsByOrderId(orderID)
	if appErr != nil {
		return appErr
	}

	total := decimal.Zero
	for _, payment := range payments {
		if payment.CapturedAmount != nil {
			total = total.Add(*payment.CapturedAmount)
		}
	}

	if err := a.Srv().Store.Order().UpdateTotalPaid(orderID, &total); err != nil {
		return model.NewAppError("UpdateOrderTotalPaid", "app.order.update_total_paid.app_error", nil, err.Error(), http.StatusInternalServerError)
	}

	return nil
}

// OrderIsPreAuthorized checks if order is pre-authorized
func (a *AppOrder) OrderIsPreAuthorized(orderID string) (bool, *model.AppError) {
	filterOptions := &payment.PaymentFilterOpts{
		OrderID:  orderID,
		IsActive: true,
		PaymentTransactionFilterOpts: payment.PaymentTransactionFilterOpts{
			Kind:           payment.AUTH,
			ActionRequired: false,
			IsSuccess:      true,
		},
	}
	exist, err := a.Srv().Store.Payment().PaymentExistWithOptions(filterOptions)
	if err != nil {
		// this err means system error, not sql not found
		return false, model.NewAppError("OrderIsPreAuthorized", "app.order.order_is_pre_authorized.app_error", nil, err.Error(), http.StatusInternalServerError)
	}

	return exist, nil
}

// OrderIsCaptured checks if given order is captured
func (a *AppOrder) OrderIsCaptured(orderID string) (bool, *model.AppError) {
	filterOptions := &payment.PaymentFilterOpts{
		OrderID:  orderID,
		IsActive: true,
		PaymentTransactionFilterOpts: payment.PaymentTransactionFilterOpts{
			Kind:           payment.CAPTURE,
			ActionRequired: false,
			IsSuccess:      true,
		},
	}
	exist, err := a.Srv().Store.Payment().PaymentExistWithOptions(filterOptions)
	if err != nil {
		// this err means system error, not sql not found error
		return false, model.NewAppError("OrderIsCaptured", "app.order.order_is_captured.app_error", nil, err.Error(), http.StatusInternalServerError)
	}

	return exist, nil
}

// OrderSubTotal returns sum of TotalPrice of all order lines that belong to given order
func (a *AppOrder) OrderSubTotal(orderID string, orderCurrency string) (*goprices.TaxedMoney, *model.AppError) {
	orderLines, appErr := a.GetAllOrderLinesByOrderId(orderID)
	if appErr != nil {
		return nil, appErr
	}

	// check if order and its child order lines have  same currencies
	if !strings.EqualFold(orderCurrency, orderLines[0].Currency) {
		return nil, model.NewAppError("OrderSubTotal", "app.order.currency_integrity.app_error", nil, "orders and order lines must have same currencies", http.StatusInternalServerError)
	}

	subTotal, _ := util.ZeroTaxedMoney(orderLines[0].Currency)

	var err error
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
func (a *AppOrder) OrderCanCancel(ord *order.Order) (bool, *model.AppError) {
	fulfillments, err := a.Srv().Store.Fulfillment().
		FilterByoption(&order.FulfillmentFilterOption{
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
func (a *AppOrder) OrderCanCapture(ord *order.Order, payment *payment.Payment) (bool, *model.AppError) {
	var err *model.AppError

	if payment == nil {
		payment, err = a.PaymentApp().GetLastOrderPayment(ord.Id)
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
func (a *AppOrder) OrderCanVoid(ord *order.Order, payment *payment.Payment) (bool, *model.AppError) {
	var err *model.AppError
	if payment == nil {
		payment, err = a.PaymentApp().GetLastOrderPayment(ord.Id)
		if err != nil {
			return false, err
		}
	}

	if payment == nil {
		return false, nil
	}

	return a.PaymentApp().PaymentCanVoid(payment)
}

// OrderCanRefund checks if order can refund
func (a *AppOrder) OrderCanRefund(ord *order.Order, payments []*payment.Payment) (bool, *model.AppError) {
	var appErr *model.AppError
	if len(payments) == 0 {
		payments, appErr = a.PaymentApp().GetAllPaymentsByOrderId(ord.Id)
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
func (a *AppOrder) CanMarkOrderAsPaid(ord *order.Order, payments []*payment.Payment) (bool, *model.AppError) {
	var appErr *model.AppError
	if len(payments) == 0 {
		payments, appErr = a.PaymentApp().GetAllPaymentsByOrderId(ord.Id)
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
func (a *AppOrder) OrderTotalAuthorized(ord *order.Order) (*goprices.Money, *model.AppError) {
	lastPayment, appErr := a.PaymentApp().GetLastOrderPayment(ord.Id)
	if appErr != nil {
		return nil, appErr
	}
	if lastPayment != nil && *lastPayment.IsActive {
		return a.PaymentApp().PaymentGetAuthorizedAmount(lastPayment)
	}

	zeroMoney, _ := util.ZeroMoney(ord.Currency)
	return zeroMoney, nil
}

// GetOrderCountryCode is helper function, returns contry code of given order
func (a *AppOrder) GetOrderCountryCode(ord *order.Order) (string, *model.AppError) {
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

	address, err := a.Srv().Store.Address().Get(*addressID)
	if err != nil {
		var errNf *store.ErrNotFound
		var statusCode int = http.StatusInternalServerError
		if errors.As(err, &errNf) {
			statusCode = http.StatusNotFound
		}
		return "", model.NewAppError("GetOrderCountryCode", "app.order.get_address.app_error", nil, err.Error(), statusCode)
	}

	return address.Country, nil
}

// CustomerEmail try finding order's owner's email. If order has no user or error occured during the finding process, returns order's UserEmail property instead
func (a *AppOrder) CustomerEmail(ord *order.Order) (string, *model.AppError) {
	if ord.UserID != nil {
		user, appErr := a.AccountApp().UserById(context.Background(), *ord.UserID)
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
