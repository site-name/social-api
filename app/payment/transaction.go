package payment

import (
	"net/http"

	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model/payment"
	"github.com/sitename/sitename/store"
)

func (a *AppPayment) GetAllPaymentTransactions(paymentID string) ([]*payment.PaymentTransaction, *model.AppError) {
	transactions, err := a.Srv().Store.PaymentTransaction().GetAllByPaymentID(paymentID)
	if err != nil {
		return nil, store.AppErrorFromDatabaseLookupError("GetAllPaymentTransactions", "app.payment.payment_transactions_not_found.app_error", err)
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

func (a *AppPayment) SaveTransaction(tran *payment.PaymentTransaction) (*payment.PaymentTransaction, *model.AppError) {
	tran, err := a.Srv().Store.PaymentTransaction().Save(tran)
	if err != nil {
		if appErr, ok := err.(*model.AppError); ok {
			return nil, appErr
		}
		return nil, model.NewAppError("SaveTransaction", "app.payment.save_transaction_error.app_error", nil, err.Error(), http.StatusInternalServerError)
	}

	return tran, nil
}
