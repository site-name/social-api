package payment

import (
	"net/http"

	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model_helper"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

func (a *ServicePayment) TransactionsByOption(option model_helper.PaymentTransactionFilterOpts) ([]*model.PaymentTransaction, *model_helper.AppError) {
	transactions, err := a.srv.Store.PaymentTransaction().FilterByOption(option)
	if err != nil {
		return nil, model_helper.NewAppError("TransactionsByOption", "app.payment.error_finding_transactions_by_option.app_error", nil, err.Error(), http.StatusInternalServerError)
	}
	return transactions, nil
}

func (a *ServicePayment) GetLastPaymentTransaction(paymentID string) (*model.PaymentTransaction, *model_helper.AppError) {
	trans, appErr := a.TransactionsByOption(model_helper.PaymentTransactionFilterOpts{
		CommonQueryOptions: model_helper.NewCommonQueryOptions(
			model.PaymentTransactionWhere.PaymentID.EQ(paymentID),
			qm.OrderBy(model.PaymentTransactionColumns.CreatedAt+" "+model_helper.DESC.String()),
			qm.Limit(1),
		),
	})

	if appErr != nil {
		return nil, appErr
	}
	if len(trans) == 0 {
		return nil, model_helper.NewAppError("GetLastPaymentTransaction", "app.payment.get_last_transaction_of_payment.app_error", nil, "payment has no transaction", http.StatusNotFound)
	}

	return trans[0], nil
}

func (a *ServicePayment) UpsertTransaction(transaction boil.ContextTransactor, paymentTransaction model.PaymentTransaction) (*model.PaymentTransaction, *model_helper.AppError) {
	upsertPaymentTransaction, err := a.srv.Store.PaymentTransaction().Upsert(transaction, paymentTransaction)
	if err != nil {
		if appErr, ok := err.(*model_helper.AppError); ok {
			return nil, appErr
		}
		return nil, model_helper.NewAppError("UpsertTransaction", "app.payment.save_transaction_error.app_error", nil, err.Error(), http.StatusInternalServerError)
	}

	return upsertPaymentTransaction, nil
}
