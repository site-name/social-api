package payment

import (
	"errors"
	"net/http"
	"strings"

	goprices "github.com/site-name/go-prices"
	"github.com/sitename/sitename/app"
	"github.com/sitename/sitename/app/sub_app_iface"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model/payment"
	"github.com/sitename/sitename/modules/util"
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

func (a *AppPayment) PaymentIsAuthorized(paymentID string) (bool, *model.AppError) {
	trans, err := a.GetAllPaymentTransactions(paymentID)
	if err != nil {
		return false, err
	}

	for _, tran := range trans {
		if tran.Kind == payment.AUTH && tran.IsSuccess && !tran.ActionRequired {
			return true, nil
		}
	}

	return false, nil
}

func (a *AppPayment) PaymentGetAuthorizedAmount(pm *payment.Payment) (*goprices.Money, *model.AppError) {
	authorizedMoney, err := util.ZeroMoney(pm.Currency)
	if err != nil {
		return nil, model.NewAppError("PaymentGetAuthorizedAmount", "app.payment.create_zero_money.app_error", nil, err.Error(), http.StatusInternalServerError)
	}
	trans, appErr := a.GetAllPaymentTransactions(pm.Id)
	if appErr != nil {
		return nil, appErr
	}

	// check if payment's Currency is same as transactions's currencies
	for _, tran := range trans {
		if !strings.EqualFold(tran.Currency, pm.Currency) {
			return nil, model.NewAppError("PaymentGetAuthorizedAmount", "app.payment.payment_transactions_currency_integrity.app_error", nil, "payment and its transactions must have same money currencies", http.StatusInternalServerError)
		}
	}

	// There is no authorized amount anymore when capture is succeeded
	// since capture can only be made once, even it is a partial capture
	for _, tran := range trans {
		if tran.Kind == payment.CAPTURE && tran.IsSuccess {
			return authorizedMoney, nil
		}
	}

	// Filter the succeeded auth transactions
	for _, tran := range trans {
		if tran.Kind == payment.AUTH && tran.IsSuccess && !tran.ActionRequired {
			// resulting error can be ignored here:
			authorizedMoney, _ = authorizedMoney.Add(&goprices.Money{
				Amount:   tran.Amount,
				Currency: tran.Currency,
			})
		}
	}

	return authorizedMoney, nil
}

func (a *AppPayment) PaymentCanVoid(pm *payment.Payment) (bool, *model.AppError) {
	authorized, err := a.PaymentIsAuthorized(pm.Id)
	if err != nil {
		return false, err
	}

	return pm.IsActive && pm.IsNotCharged() && authorized, nil
}

func (a *AppPayment) SavePayment(pm *payment.Payment) (*payment.Payment, *model.AppError) {
	newPm, err := a.Srv().Store.Payment().Save(pm)
	if err != nil {
		if appErr, ok := err.(*model.AppError); ok {
			return nil, appErr
		}
		return nil, model.NewAppError("SavePayment", "app.payment.save_payment_error.app_error", nil, "", http.StatusInternalServerError)
	}

	return newPm, nil
}

func (a *AppPayment) GetPaymentToken(paymentID string) (string, *model.AppError) {
	trans, appErr := a.GetAllPaymentTransactions(paymentID)
	if appErr != nil {
		return "", appErr
	}

	var tran *payment.PaymentTransaction
	for _, tr := range trans {
		if tr.Kind == payment.AUTH && tr.IsSuccess {
			if tran == nil || tran.CreateAt > tr.CreateAt {
				tran = tr
			}
		}
	}

	if tran == nil {
		return "", model.NewAppError("GetPaymentToken", "app.payment.no_authorized_payment_transaction.app_error", nil, "", http.StatusNotFound)
	}

	return tran.Token, nil
}