package payment

import (
	"net/http"

	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model/payment"
)

// TransactionsByOption returns a list of transactions filtered based on given option
func (a *AppPayment) TransactionsByOption(option *payment.PaymentTransactionFilterOpts) ([]*payment.PaymentTransaction, *model.AppError) {
	transactions, err := a.app.Srv().Store.PaymentTransaction().FilterByOption(option)

	var statusCode int
	var appErrMsg string
	if err != nil {
		statusCode = http.StatusInternalServerError
		appErrMsg = err.Error()
	}
	if len(transactions) == 0 {
		statusCode = http.StatusNotFound
	}

	if statusCode == 0 {
		return transactions, nil
	}
	return nil, model.NewAppError("TransactionsByOption", "app.payment.error_finding_transactions_by_option.app_error", nil, appErrMsg, statusCode)
}

func (a *AppPayment) GetAllPaymentTransactions(paymentID string) ([]*payment.PaymentTransaction, *model.AppError) {
	transactions, appErr := a.TransactionsByOption(&payment.PaymentTransactionFilterOpts{
		PaymentID: &model.StringFilter{
			StringOption: &model.StringOption{
				Eq: paymentID,
			},
		},
	})
	if appErr != nil {
		return nil, appErr
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

	var lastTran *payment.PaymentTransaction
	for _, tran := range trans {
		if lastTran == nil || tran.CreateAt >= lastTran.CreateAt {
			lastTran = tran
		}
	}

	return lastTran, nil
}

func (a *AppPayment) SaveTransaction(tran *payment.PaymentTransaction) (*payment.PaymentTransaction, *model.AppError) {
	tran, err := a.app.Srv().Store.PaymentTransaction().Save(tran)
	if err != nil {
		if appErr, ok := err.(*model.AppError); ok {
			return nil, appErr
		}
		return nil, model.NewAppError("SaveTransaction", "app.payment.save_transaction_error.app_error", nil, err.Error(), http.StatusInternalServerError)
	}

	return tran, nil
}

func (a *AppPayment) UpdateTransaction(transaction *payment.PaymentTransaction) (*payment.PaymentTransaction, *model.AppError) {
	tran, err := a.app.Srv().Store.PaymentTransaction().Update(transaction)
	if err != nil {
		if appErr, ok := err.(*model.AppError); ok {
			return nil, appErr
		}
		return nil, model.NewAppError("UpdateTransaction", "app.payment.error_updating.transaction.app_error", nil, err.Error(), http.StatusInternalServerError)
	}

	return tran, nil
}
