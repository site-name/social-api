/*
	Since Go does not support @decorator like in python and generic,
	Then you are about to find this code file not elegant at all.
	But that is fine.
*/
package payment

import (
	"fmt"
	"net/http"

	"github.com/mattermost/gorp"
	"github.com/site-name/decimal"
	"github.com/sitename/sitename/app"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model/payment"
)

const (
	ERROR_MSG                 = "Oops! Something went wrong"
	GENERIC_TRANSACTION_ERROR = "Transaction was unsuccessful."
)

// raisePaymentError must be called right before function returns
func (a *ServicePayment) raisePaymentError(where string, transaction *payment.PaymentTransaction) (*payment.PaymentTransaction, *payment.PaymentError) {
	if !transaction.IsSuccess {
		msg := GENERIC_TRANSACTION_ERROR
		if transaction.Error != nil {
			msg = *transaction.Error
		}
		return nil, payment.NewPaymentError(where, msg, payment.INVALID)
	}

	return transaction, nil
}

// paymentPostProcess must be called right before function returns
func (a *ServicePayment) paymentPostProcess(transaction *payment.PaymentTransaction) *model.AppError {
	payMent, appErr := a.PaymentByID(nil, transaction.PaymentID, false)
	if appErr != nil {
		return appErr
	}
	return a.GatewayPostProcess(transaction, payMent)
}

// requireActivePayment must be called in the beginning of the function body
func (a *ServicePayment) requireActivePayment(where string, payMent *payment.Payment) (*payment.Payment, *payment.PaymentError) {
	if !*payMent.IsActive {
		return nil, payment.NewPaymentError(where, "This payment is no longer active", payment.INVALID)
	}
	return payMent, nil
}

// withLockedPayment Lock payment to protect from asynchronous modification.
func (a *ServicePayment) withLockedPayment(where string, payMent *payment.Payment) (*payment.Payment, *gorp.Transaction, *model.AppError) {
	transaction, err := a.srv.Store.GetMaster().Begin()
	if err != nil {
		return nil, nil, model.NewAppError(where, app.ErrorCreatingTransactionErrorID, nil, err.Error(), http.StatusInternalServerError)
	}

	paymentToOperateOn, appErr := a.srv.PaymentService().PaymentByID(transaction, payMent.Id, true)
	if appErr != nil {
		a.srv.Store.FinalizeTransaction(transaction)
		return nil, nil, appErr
	}

	return paymentToOperateOn, transaction, nil
}

func (a *ServicePayment) ProcessPayment(
	payMent *payment.Payment,
	token string,
	manager interface{},
	channelSlug string,
	customerID *string, // can be empty
	storeSource bool, // default to false
	additionalData map[string]interface{}, // can be nil

) (*payment.PaymentTransaction, *payment.PaymentError, *model.AppError) {

	payMent, paymentErr := a.requireActivePayment("ProcessPayment", payMent)
	if paymentErr != nil {
		return nil, paymentErr, nil
	}

	payMent, transaction, appErr := a.withLockedPayment("ProcessPayment", payMent)
	if appErr != nil {
		return nil, nil, appErr
	}
	defer a.srv.Store.FinalizeTransaction(transaction)

	_, appErr = a.CreatePaymentInformation(payMent, &token, nil, customerID, storeSource, additionalData)
	if appErr != nil {
		return nil, nil, appErr
	}

	panic("not implemented")
}

func (a *ServicePayment) Authorize(
	payMent *payment.Payment,
	token string,
	manager interface{},
	channelSlug string,
	customerID *string,
	storeSource bool,

) (*payment.PaymentTransaction, *payment.PaymentError, *model.AppError) {

	payMent, paymentErr := a.requireActivePayment("ProcessPayment", payMent)
	if paymentErr != nil {
		return nil, paymentErr, nil
	}

	payMent, transaction, appErr := a.withLockedPayment("ProcessPayment", payMent)
	if appErr != nil {
		return nil, nil, appErr
	}
	defer a.srv.Store.FinalizeTransaction(transaction)

	paymentErr = a.CleanAuthorize(payMent)
	if paymentErr != nil {
		return nil, paymentErr, nil
	}

	_, appErr = a.CreatePaymentInformation(payMent, &token, nil, customerID, storeSource, nil)
	if appErr != nil {
		return nil, nil, appErr
	}

	panic("not implt")
}

func (a *ServicePayment) Capture(
	payMent *payment.Payment,
	manager interface{},
	channelSlug string,
	amount *decimal.Decimal, // can be nil
	customerID *string, // can be nil
	storeSource bool, // default false

) (*payment.PaymentTransaction, *payment.PaymentError, *model.AppError) {

	payMent, paymentErr := a.requireActivePayment("ProcessPayment", payMent)
	if paymentErr != nil {
		return nil, paymentErr, nil
	}

	payMent, transaction, appErr := a.withLockedPayment("ProcessPayment", payMent)
	if appErr != nil {
		return nil, nil, appErr
	}
	defer a.srv.Store.FinalizeTransaction(transaction)

	panic("not implemented")
}

func (a *ServicePayment) Refund(
	payMent *payment.Payment,
	manager interface{},
	channelSlug string,
	amount *decimal.Decimal, // can be nil

) (*payment.PaymentTransaction, *payment.PaymentError, *model.AppError) {
	panic("not implt")
}

func (a *ServicePayment) Void(payMent *payment.Payment, manager interface{}, channelSlug string) (*payment.PaymentTransaction, *payment.PaymentError, *model.AppError) {
	panic("not implt")
}

// Confirm confirms payment
func (a *ServicePayment) Confirm(
	payMent *payment.Payment,
	manager interface{},
	channelSlug string,
	additionalData map[string]interface{}, // can be none

) (*payment.PaymentTransaction, *payment.PaymentError, *model.AppError) {
	panic("not implt")
}

func (a *ServicePayment) ListPaymentSources(
	gateway string,
	customerID string,
	manager interface{},
	channelSlug string,

) ([]*payment.CustomerSource, *model.AppError) {
	panic("not implemented")
}

func (a *ServicePayment) ListGateways(
	manager interface{},
	channelSlug string,
) ([]*payment.PaymentGateway, *model.AppError) {
	panic("not implemented")
}

func (a *ServicePayment) fetchGatewayResponse(fn func()) {
	panic("not implemented")
}

func (a *ServicePayment) getPastTransactionToken(payMent *payment.Payment, kind string) (string, *payment.PaymentError, *model.AppError) {
	transactions, appErr := a.TransactionsByOption(&payment.PaymentTransactionFilterOpts{
		PaymentID: &model.StringFilter{
			StringOption: &model.StringOption{
				Eq: payMent.Id,
			},
		},
		Kind: &model.StringFilter{
			StringOption: &model.StringOption{
				Eq: kind,
			},
		},
		IsSuccess: model.NewBool(true),
	})
	if appErr != nil && appErr.StatusCode == http.StatusInternalServerError {
		return "", nil, appErr
	}

	if length := len(transactions); length == 0 {
		return "", payment.NewPaymentError("getPastTransactionToken", fmt.Sprintf("Cannot find successful %s transaction", kind), payment.NOT_FOUND), nil
	} else {
		return transactions[length-1].Token, nil, nil
	}
}

func (a *ServicePayment) validateRefundAmount(payMent *payment.Payment, amount *decimal.Decimal) *payment.PaymentError {
	if amount.LessThan(decimal.Zero) {
		return payment.NewPaymentError("validateRefundAmount", "Amount should be positive number", payment.INVALID)
	}
	if payMent.CapturedAmount == nil {
		payMent.CapturedAmount = &decimal.Zero
	}
	if amount.GreaterThan(*payMent.CapturedAmount) {
		return payment.NewPaymentError("validateRefundAmount", "Cannot refund more than captures.", payment.INVALID)
	}

	return nil
}

// PaymentRefundOrVoid
func (a *ServicePayment) PaymentRefundOrVoid(payMent *payment.Payment, manager interface{}, channelSlug string) (*payment.PaymentError, *model.AppError) {
	if payMent == nil {
		return nil, nil
	}

	paymentCanVoid, appErr := a.srv.PaymentService().PaymentCanVoid(payMent)
	if appErr != nil {
		return nil, appErr
	}

	var paymentErr *payment.PaymentError

	if payMent.CanRefund() {
		_, paymentErr, appErr = a.Refund(payMent, manager, channelSlug, nil)
	} else if paymentCanVoid {
		_, paymentErr, appErr = a.Void(payMent, manager, channelSlug)
	}

	return paymentErr, appErr
}
