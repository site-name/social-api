package payment

import (
	"net/http"

	"github.com/Masterminds/squirrel"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model_helper"
	"gorm.io/gorm"
)

// TransactionsByOption returns a list of transactions filtered based on given option
func (a *ServicePayment) TransactionsByOption(option *model.PaymentTransactionFilterOpts) ([]*model.PaymentTransaction, *model_helper.AppError) {
	transactions, err := a.srv.Store.PaymentTransaction().FilterByOption(option)
	if err != nil {
		return nil, model_helper.NewAppError("TransactionsByOption", "app.payment.error_finding_transactions_by_option.app_error", nil, err.Error(), http.StatusInternalServerError)
	}
	return transactions, nil
}

func (a *ServicePayment) GetLastPaymentTransaction(paymentID string) (*model.PaymentTransaction, *model_helper.AppError) {
	trans, appErr := a.TransactionsByOption(&model.PaymentTransactionFilterOpts{
		Conditions: squirrel.Eq{model.TransactionTableName + "." + model.TransactionColumnPaymentID: paymentID},
	})

	if appErr != nil {
		return nil, appErr
	}

	if len(trans) == 0 {
		return nil, model_helper.NewAppError("GetLastPaymentTransaction", "app.payment.get_last_transaction_of_payment.app_error", nil, "payment has no transaction", http.StatusNotFound)
	}

	var lastTran *model.PaymentTransaction
	for _, tran := range trans {
		if lastTran == nil || tran.CreateAt >= lastTran.CreateAt {
			lastTran = tran
		}
	}

	return lastTran, nil
}

func (a *ServicePayment) SaveTransaction(transaction *gorm.DB, paymentTransaction *model.PaymentTransaction) (*model.PaymentTransaction, *model_helper.AppError) {
	paymentTransaction, err := a.srv.Store.PaymentTransaction().Save(transaction, paymentTransaction)
	if err != nil {
		if appErr, ok := err.(*model_helper.AppError); ok {
			return nil, appErr
		}
		return nil, model_helper.NewAppError("SaveTransaction", "app.payment.save_transaction_error.app_error", nil, err.Error(), http.StatusInternalServerError)
	}

	return paymentTransaction, nil
}

func (a *ServicePayment) UpdateTransaction(transaction *model.PaymentTransaction) (*model.PaymentTransaction, *model_helper.AppError) {
	paymentTransaction, err := a.srv.Store.PaymentTransaction().Update(transaction)
	if err != nil {
		if appErr, ok := err.(*model_helper.AppError); ok {
			return nil, appErr
		}
		return nil, model_helper.NewAppError("UpdateTransaction", "app.payment.error_updating.transaction.app_error", nil, err.Error(), http.StatusInternalServerError)
	}

	return paymentTransaction, nil
}
