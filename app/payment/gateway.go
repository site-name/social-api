package payment

import (
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model/payment"
)

const (
	ERROR_MSG                 = "Oops! Something went wrong"
	GENERIC_TRANSACTION_ERROR = "Transaction was unsuccessful."
)

func (a *AppPayment) raisePaymentError(where string, transaction *payment.PaymentTransaction) (*payment.PaymentTransaction, *payment.PaymentError) {
	if !transaction.IsSuccess {
		msg := GENERIC_TRANSACTION_ERROR
		if transaction.Error != nil {
			msg = *transaction.Error
		}
		return nil, payment.NewPaymentError(where, msg, payment.INVALID)
	}

	return transaction, nil
}

func (a *AppPayment) paymentPostProcess(transaction *payment.PaymentTransaction) *model.AppError {
	payMent, appErr := a.PaymentByID(transaction.PaymentID, false)
	if appErr != nil {
		return appErr
	}
	return a.GatewayPostProcess(transaction, payMent)
}

func (a *AppPayment) requireActivePayment(where string, payMent *payment.Payment) (*payment.Payment, *payment.PaymentError) {
	if !*payMent.IsActive {
		return nil, payment.NewPaymentError(where, "This payment is no longer active", payment.INVALID)
	}
	return payMent, nil
}

// withLockedPayment Lock payment to protect from asynchronous modification.
func (a *AppPayment) withLockedPayment(where string, payMent *payment.Payment) (*payment.Payment, *model.AppError) {
	return a.PaymentByID(payMent.Id, true)
}

func (a *AppPayment) ProcessPayment(
	payMent *payment.Payment,
	token string,
	manager interface{},
	channelSlug string,
	customerID string, // can be empty
	storeSource bool, // default to false
	additionalData map[string]string, // can be nil

) (*payment.PaymentTransaction, *payment.PaymentError, *model.AppError) {

	payMent, paymentErr := a.requireActivePayment("ProcessPayment", payMent)
	if paymentErr != nil {
		return nil, paymentErr, nil
	}

	payMent, appErr := a.withLockedPayment("ProcessPayment", payMent)
	if appErr != nil {
		return nil, nil, appErr
	}

	_, appErr = a.CreatePaymentInformation(payMent, &token, nil, &customerID, storeSource, additionalData)
	if appErr != nil {
		return nil, nil, appErr
	}

	panic("not implemented")
}
