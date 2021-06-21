package payment

import (
	"errors"
	"net/http"

	"github.com/sitename/sitename/app"
	"github.com/sitename/sitename/app/sub_app_iface"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model/payment"
	"github.com/sitename/sitename/store"
)

type AppPayment struct {
	app.AppIface
}

func init() {
	app.RegisterPaymentApp(func(a app.AppIface) sub_app_iface.PaymentApp {
		return &AppPayment{a}
	})
}

func (a *AppPayment) GetAllPaymentsByOrderId(orderID string) ([]*payment.Payment, *model.AppError) {
	payments, err := a.Srv().Store.Payment().GetPaymentsByOrderID(orderID)
	if err != nil {
		var statusCode int = http.StatusInternalServerError
		var nfErr *store.ErrNotFound
		if errors.As(err, &nfErr) {
			statusCode = http.StatusNotFound
		}
		return nil, model.NewAppError("GetAllPaymentsByOrderId", "app.order.get_child_payments.app_error", nil, err.Error(), statusCode)
	}

	return payments, nil
}

func (a *AppPayment) GetLastOrderPayment(orderID string) (*payment.Payment, *model.AppError) {
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
