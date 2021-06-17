package app

import (
	"errors"
	"net/http"
	"strings"

	goprices "github.com/site-name/go-prices"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model/payment"
	"github.com/sitename/sitename/modules/util"
	"github.com/sitename/sitename/store"
)

// GetAllPaymentTransactions returns all transactions belong to given payment
func (a *App) GetAllPaymentTransactions(paymentID string) ([]*payment.PaymentTransaction, *model.AppError) {
	transactions, err := a.srv.Store.PaymentTransaction().GetAllByPaymentID(paymentID)
	if err != nil {
		var nfErr *store.ErrNotFound
		var statusCode int
		if errors.As(err, &nfErr) {
			statusCode = http.StatusNotFound
		} else {
			statusCode = http.StatusInternalServerError
		}

		return nil, model.NewAppError("GetAllPaymentTransactions", "app.payment.get_associated_transactions.app_error", nil, err.Error(), statusCode)
	}

	return transactions, nil
}

// GetLastPaymentTransaction return most recent transaction made for given payment
func (a *App) GetLastPaymentTransaction(paymentID string) (*payment.PaymentTransaction, *model.AppError) {
	trans, appErr := a.GetAllPaymentTransactions(paymentID)
	if appErr != nil {
		return nil, appErr
	}

	if len(trans) == 0 {
		return nil, nil
	}

	if len(trans) == 1 {
		return trans[0], nil
	}

	lastTran := trans[0]
	for _, tran := range trans {
		if tran != nil && tran.CreateAt >= lastTran.CreateAt {
			lastTran = tran
		}
	}

	return lastTran, nil
}

// PaymentIsAuthorized checks if given payment is authorized
func (a *App) PaymentIsAuthorized(paymentID string) (bool, *model.AppError) {
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

// PaymentGetAuthorizedAmount
func (a *App) PaymentGetAuthorizedAmount(pm *payment.Payment) (*goprices.Money, *model.AppError) {
	zeroMoney, err := util.ZeroMoney(pm.Currency)
	if err != nil {
		return nil, model.NewAppError("PaymentGetAuthorizedAmount", "app.payment.create_zero_money.app_error", nil, err.Error(), http.StatusInternalServerError)
	}
	trans, appErr := a.GetAllPaymentTransactions(pm.Id)
	if appErr != nil {
		return nil, appErr
	}

	if !strings.EqualFold(pm.Currency, trans[0].Currency) {
		return nil, model.NewAppError("PaymentGetAuthorizedAmount", "app.payment.payment_transactions_currency_integrity.app_error", nil, "payment and its transactions must have same money currencies", http.StatusInternalServerError)
	}

	for _, tran := range trans {
		if tran.Kind == payment.CAPTURE && tran.IsSuccess {
			return zeroMoney, nil
		}
	}

	for _, tran := range trans {
		if tran.Kind == payment.AUTH && tran.IsSuccess && !tran.ActionRequired {
			zeroMoney, err = zeroMoney.Add(&goprices.Money{Amount: tran.Amount, Currency: tran.Currency})
			if err != nil {
				return nil, model.NewAppError("PaymentGetAuthorizedAmount", "app.payment.add_money.app_error", nil, err.Error(), http.StatusInternalServerError)
			}
		}
	}

	return zeroMoney, nil
}

// PaymentCanVoid
func (a *App) PaymentCanVoid(pm *payment.Payment) (bool, *model.AppError) {
	authorized, err := a.PaymentIsAuthorized(pm.Id)
	if err != nil {
		return false, err
	}

	return pm.IsActive && pm.IsNotCharged() && authorized, nil
}
