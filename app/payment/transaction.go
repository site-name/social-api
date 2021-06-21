package payment

import (
	"errors"
	"net/http"

	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model/payment"
	"github.com/sitename/sitename/store"
)

func (a *AppPayment) GetAllPaymentTransactions(paymentID string) ([]*payment.PaymentTransaction, *model.AppError) {
	transactions, err := a.Srv().Store.PaymentTransaction().GetAllByPaymentID(paymentID)
	if err != nil {
		var nfErr *store.ErrNotFound
		var statusCode int = http.StatusInternalServerError
		if errors.As(err, &nfErr) {
			statusCode = http.StatusNotFound
		}

		return nil, model.NewAppError("GetAllPaymentTransactions", "app.payment.get_associated_transactions.app_error", nil, err.Error(), statusCode)
	}

	return transactions, nil
}

func (a *AppPayment) GetLastPaymentTransaction(paymentID string) (*payment.PaymentTransaction, *model.AppError) {
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
