/*
Since Go does not support @decorator like in python and generic,
Then you are about to find this code file not elegant at all.
But that is fine.
*/
package payment

import (
	"fmt"
	"net/http"

	"github.com/site-name/decimal"
	"github.com/sitename/sitename/app/plugin/interfaces"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model_helper"
	"github.com/volatiletech/sqlboiler/v4/boil"
)

const (
	ERROR_MSG                 = "Oops! Something went wrong"
	GENERIC_TRANSACTION_ERROR = "Transaction was unsuccessful."
)

// PaymentMethod is type for some methods of PluginManager.
// They are:
//
// 1) AuthorizePayment
//
// 2) CapturePayment
//
// 3) ConfirmPayment
//
// 4) ProcessPayment
//
// 5) RefundPayment
//
// 6) VoidPayment
type PaymentMethod func(gateway string, paymentInformation model_helper.PaymentData, channelID string) (*model_helper.GatewayResponse, error)

func (a *ServicePayment) raisePaymentError(where string, transaction model.PaymentTransaction) *model_helper.PaymentError {
	if !transaction.IsSuccess {
		msg := GENERIC_TRANSACTION_ERROR
		if !transaction.Error.IsNil() {
			msg = *transaction.Error.String
		}
		return model_helper.NewPaymentError(where, msg, model_helper.INVALID)
	}

	return nil
}

func (a *ServicePayment) paymentPostProcess(transaction model.PaymentTransaction) *model_helper.AppError {
	payment, appErr := a.PaymentByID(nil, transaction.PaymentID, false)
	if appErr != nil {
		return appErr
	}
	return a.GatewayPostProcess(transaction, *payment)
}

func (a *ServicePayment) requireActivePayment(where string, payment model.Payment) *model_helper.PaymentError {
	if !payment.IsActive {
		return model_helper.NewPaymentError(where, "This payment is no longer active", model_helper.INVALID)
	}
	return nil
}

func (a *ServicePayment) withLockedPayment(dbTransaction boil.ContextTransactor, where string, payment model.Payment) (*model.Payment, *model_helper.AppError) {
	paymentToOperateOn, appErr := a.PaymentByID(dbTransaction, payment.ID, true)
	if appErr != nil {
		return nil, appErr
	}

	return paymentToOperateOn, nil
}

// @requireActivePayment
//
// @withLockedPayment
//
// @raisePaymentError
//
// @paymentPostProcess
func (a *ServicePayment) ProcessPayment(
	dbTransaction boil.ContextTransactor,
	payment model.Payment,
	token string,
	manager interfaces.PluginManagerInterface,
	channelID string, // originally is channelSlug in saleor
	customerID *string,
	storeSource bool,
	additionalData map[string]any,
) (*model.PaymentTransaction, *model_helper.PaymentError, *model_helper.AppError) {
	paymentErr := a.requireActivePayment("ProcessPayment", payment)
	if paymentErr != nil {
		return nil, paymentErr, nil
	}

	lockedPayment, appErr := a.withLockedPayment(dbTransaction, "ProcessPayment", payment)
	if appErr != nil {
		return nil, nil, appErr
	}

	paymentData, appErr := a.CreatePaymentInformation(*lockedPayment, &token, nil, customerID, storeSource, additionalData)
	if appErr != nil {
		return nil, nil, appErr
	}

	response, errMsg := a.fetchGatewayResponse(manager.ProcessPayment, lockedPayment.Gateway, *paymentData, channelID)
	actionRequired := response != nil && response.ActionRequired

	if response != nil {
		appErr = a.UpdatePayment(*lockedPayment, response)
		if appErr != nil {
			return nil, nil, appErr
		}
	}

	paymentTransaction, appErr := a.GetAlreadyProcessedTransactionOrCreateNewTransaction(lockedPayment.ID, model.TransactionKindCapture, paymentData, actionRequired, response, errMsg)
	if appErr != nil {
		return nil, nil, appErr
	}

	paymentErr = a.raisePaymentError("ProcessPayment", *paymentTransaction)
	if paymentErr != nil {
		return nil, paymentErr, nil
	}

	appErr = a.paymentPostProcess(*paymentTransaction)
	if appErr != nil {
		return nil, nil, appErr
	}

	return paymentTransaction, nil, nil
}

// @requireActivePayment
//
// @withLockedPayment
//
// @raisePaymentError
//
// @paymentPostProcess
func (a *ServicePayment) Authorize(
	dbTransaction boil.ContextTransactor,
	payment model.Payment,
	token string,
	manager interfaces.PluginManagerInterface,
	channelID string,
	customerID *string,
	storeSource bool,

) (*model.PaymentTransaction, *model_helper.PaymentError, *model_helper.AppError) {
	paymentErr := a.requireActivePayment("Authorize", payment)
	if paymentErr != nil {
		return nil, paymentErr, nil
	}

	lockedPayment, appErr := a.withLockedPayment(dbTransaction, "Authorize", payment)
	if appErr != nil {
		return nil, nil, appErr
	}

	paymentErr = a.CleanAuthorize(*lockedPayment)
	if paymentErr != nil {
		return nil, paymentErr, nil
	}

	paymentData, appErr := a.CreatePaymentInformation(*lockedPayment, &token, nil, customerID, storeSource, nil)
	if appErr != nil {
		return nil, nil, appErr
	}

	response, errMsg := a.fetchGatewayResponse(manager.AuthorizePayment, lockedPayment.Gateway, *paymentData, channelID)
	if response != nil {
		appErr = a.UpdatePayment(*lockedPayment, response)
		if appErr != nil {
			return nil, nil, appErr
		}
	}

	paymentTransaction, appErr := a.GetAlreadyProcessedTransactionOrCreateNewTransaction(lockedPayment.ID, model.TransactionKindCapture, paymentData, false, response, errMsg)
	if appErr != nil {
		return nil, nil, appErr
	}

	paymentErr = a.raisePaymentError("Authorize", *paymentTransaction)
	if paymentErr != nil {
		return nil, paymentErr, nil
	}

	appErr = a.paymentPostProcess(*paymentTransaction)
	if appErr != nil {
		return nil, nil, appErr
	}

	return paymentTransaction, nil, nil
}

// @requireActivePayment
//
// @withLockedPayment
//
// @raisePaymentError
//
// @paymentPostProcess
func (a *ServicePayment) Capture(
	dbTransaction boil.ContextTransactor,
	payment model.Payment,
	manager interfaces.PluginManagerInterface,
	channelID string,
	amount *decimal.Decimal, // can be nil
	customerID *string, // can be nil
	storeSource bool, // default false

) (*model.PaymentTransaction, *model_helper.PaymentError, *model_helper.AppError) {

	paymentErr := a.requireActivePayment("Capture", payment)
	if paymentErr != nil {
		return nil, paymentErr, nil
	}

	lockedPayment, appErr := a.withLockedPayment(dbTransaction, "Capture", payment)
	if appErr != nil {
		return nil, nil, appErr
	}

	if amount == nil {
		amount = model_helper.GetPointerOfValue(model_helper.PaymentGetChargeAmount(*lockedPayment))
	}

	paymentErr = a.CleanCapture(*lockedPayment, *amount)
	if paymentErr != nil {
		return nil, paymentErr, nil
	}

	paymentData, appErr := a.CreatePaymentInformation(*lockedPayment, &lockedPayment.Token, amount, customerID, storeSource, nil)
	if appErr != nil {
		return nil, nil, appErr
	}

	response, errMsg := a.fetchGatewayResponse(manager.CapturePayment, lockedPayment.Gateway, *paymentData, channelID)
	if response != nil {
		appErr = a.UpdatePayment(*lockedPayment, response)
		if appErr != nil {
			return nil, nil, appErr
		}
	}

	paymentTransaction, appErr := a.GetAlreadyProcessedTransactionOrCreateNewTransaction(lockedPayment.ID, model.TransactionKindCapture, paymentData, false, response, errMsg)
	if appErr != nil {
		return nil, nil, appErr
	}

	paymentErr = a.raisePaymentError("Capture", *paymentTransaction)
	if paymentErr != nil {
		return nil, paymentErr, nil
	}

	appErr = a.paymentPostProcess(*paymentTransaction)
	if appErr != nil {
		return nil, nil, appErr
	}

	return paymentTransaction, nil, nil
}

// @requireActivePayment
//
// @withLockedPayment
//
// @raisePaymentError
//
// @paymentPostProcess
func (a *ServicePayment) Refund(
	dbTransaction boil.ContextTransactor,
	payment model.Payment,
	manager interfaces.PluginManagerInterface,
	channelID string,
	amount *decimal.Decimal, // can be nil
) (*model.PaymentTransaction, *model_helper.PaymentError, *model_helper.AppError) {
	paymentErr := a.requireActivePayment("Refund", payment)
	if paymentErr != nil {
		return nil, paymentErr, nil
	}

	lockedPayment, appErr := a.withLockedPayment(dbTransaction, "Refund", payment)
	if appErr != nil {
		return nil, nil, appErr
	}

	if amount == nil {
		amount = &lockedPayment.CapturedAmount
	}

	paymentErr = a.validateRefundAmount(*lockedPayment, *amount)
	if paymentErr != nil {
		return nil, paymentErr, nil
	}

	if !model_helper.PaymentCanRefund(*lockedPayment) {
		return nil, model_helper.NewPaymentError("Refund", "This payment cannot be refunded", model_helper.INVALID), nil
	}

	kind := model.TransactionKindCapture
	if model_helper.PaymentIsManual(*lockedPayment) {
		kind = model.TransactionKindExternal
	}

	token, paymentErr, appErr := a.getPastTransactionToken(*lockedPayment, kind)
	if paymentErr != nil || appErr != nil {
		return nil, paymentErr, appErr
	}

	paymentData, appErr := a.CreatePaymentInformation(*lockedPayment, &token, amount, nil, false, nil)
	if appErr != nil {
		return nil, nil, appErr
	}

	var paymentTransaction *model.PaymentTransaction

	if model_helper.PaymentIsManual(*lockedPayment) {
		// for manual payment we just need to mark payment as a refunded
		paymentTransaction, appErr = a.CreateTransaction(lockedPayment.ID, model.TransactionKindRefund, paymentData, false, nil, "", true)
		if appErr != nil {
			return nil, nil, appErr
		}
	} else {
		response, errMsg := a.fetchGatewayResponse(manager.RefundPayment, lockedPayment.Gateway, *paymentData, channelID)
		paymentTransaction, appErr = a.GetAlreadyProcessedTransactionOrCreateNewTransaction(
			lockedPayment.ID,
			model.TransactionKindRefund,
			paymentData,
			false,
			response,
			errMsg,
		)
	}

	paymentErr = a.raisePaymentError("Refund", *paymentTransaction)
	if paymentErr != nil {
		return nil, paymentErr, nil
	}

	appErr = a.paymentPostProcess(*paymentTransaction)
	if appErr != nil {
		return nil, nil, appErr
	}

	return paymentTransaction, nil, nil
}

// @requireActivePayment
//
// @withLockedPayment
//
// @raisePaymentError
//
// @paymentPostProcess
func (a *ServicePayment) Void(dbTransaction boil.ContextTransactor, payment model.Payment, manager interfaces.PluginManagerInterface, channelID string) (*model.PaymentTransaction, *model_helper.PaymentError, *model_helper.AppError) {
	paymentErr := a.requireActivePayment("Refund", payment)
	if paymentErr != nil {
		return nil, paymentErr, nil
	}

	lockedPayment, appErr := a.withLockedPayment(dbTransaction, "Refund", payment)
	if appErr != nil {
		return nil, nil, appErr
	}

	token, paymentErr, appErr := a.getPastTransactionToken(*lockedPayment, model.TransactionKindAuth)
	if paymentErr != nil || appErr != nil {
		return nil, paymentErr, appErr
	}

	paymentData, appErr := a.CreatePaymentInformation(*lockedPayment, &token, nil, nil, false, nil)
	if appErr != nil {
		return nil, nil, appErr
	}

	response, errMsg := a.fetchGatewayResponse(manager.VoidPayment, lockedPayment.Gateway, *paymentData, channelID)
	paymentTransaction, appErr := a.GetAlreadyProcessedTransactionOrCreateNewTransaction(
		lockedPayment.ID,
		model.TransactionKindVoid,
		paymentData,
		false,
		response,
		errMsg,
	)

	paymentErr = a.raisePaymentError("Refund", *paymentTransaction)
	if paymentErr != nil {
		return nil, paymentErr, nil
	}

	appErr = a.paymentPostProcess(*paymentTransaction)
	if appErr != nil {
		return nil, nil, appErr
	}

	return paymentTransaction, nil, nil
}

// @requireActivePayment
//
// @withLockedPayment
//
// @raisePaymentError
//
// @paymentPostProcess
// Confirm confirms payment
func (a *ServicePayment) Confirm(
	dbTransaction boil.ContextTransactor,
	payment model.Payment,
	manager interfaces.PluginManagerInterface,
	channelID string,
	additionalData map[string]any, // can be none
) (*model.PaymentTransaction, *model_helper.PaymentError, *model_helper.AppError) {
	paymentErr := a.requireActivePayment("Confirm", payment)
	if paymentErr != nil {
		return nil, paymentErr, nil
	}

	lockedPayment, appErr := a.withLockedPayment(dbTransaction, "Confirm", payment)
	if appErr != nil {
		return nil, nil, appErr
	}

	transactionsOfPayment, appErr := a.TransactionsByOption(model_helper.PaymentTransactionFilterOpts{
		CommonQueryOptions: model_helper.NewCommonQueryOptions(
			model.PaymentTransactionWhere.Kind.EQ(model.TransactionKindActionToConfirm),
			model.PaymentTransactionWhere.IsSuccess.EQ(true),
		),
	})
	if appErr != nil {
		return nil, nil, appErr
	}

	var (
		lastTransaction = transactionsOfPayment[len(transactionsOfPayment)-1]
		token           string
	)
	if lastTransaction.Token != "" {
		token = lastTransaction.Token
	}

	paymentData, appErr := a.CreatePaymentInformation(*lockedPayment, &token, nil, nil, false, additionalData)
	if appErr != nil {
		return nil, nil, appErr
	}

	response, errMsg := a.fetchGatewayResponse(manager.ConfirmPayment, lockedPayment.Gateway, *paymentData, channelID)
	actionRequired := response != nil && response.ActionRequired

	if response != nil {
		appErr = a.UpdatePayment(*lockedPayment, response)
		if appErr != nil {
			return nil, nil, appErr
		}
	}

	paymentTransaction, appErr := a.GetAlreadyProcessedTransactionOrCreateNewTransaction(lockedPayment.ID, model.TransactionKindConfirm, paymentData, actionRequired, response, errMsg)
	if appErr != nil {
		return nil, nil, appErr
	}

	paymentErr = a.raisePaymentError("Refund", *paymentTransaction)
	if paymentErr != nil {
		return nil, paymentErr, nil
	}

	appErr = a.paymentPostProcess(*paymentTransaction)
	if appErr != nil {
		return nil, nil, appErr
	}

	return paymentTransaction, nil, nil
}

func (a *ServicePayment) ListPaymentSources(
	gateway string,
	customerID string,
	manager interfaces.PluginManagerInterface,
	channelID string,
) ([]*model_helper.CustomerSource, *model_helper.AppError) {
	source, err := manager.ListPaymentSources(gateway, customerID, channelID)
	if err != nil {
		return nil, model_helper.NewAppError("ListPaymentSources", "app.payment.error_listing_payment_sources.app_error", nil, err.Error(), http.StatusInternalServerError)
	}
	return source, nil
}

func (a *ServicePayment) ListGateways(
	manager interfaces.PluginManagerInterface,
	channelID string,
) []*model_helper.PaymentGateway {
	return manager.ListPaymentGateways("", nil, channelID, true)
}

func (a *ServicePayment) fetchGatewayResponse(paymentFunc PaymentMethod, gateway string, paymentData model_helper.PaymentData, channelID string) (res *model_helper.GatewayResponse, errMsg string) {
	res, _ = paymentFunc(gateway, paymentData, channelID)
	gatewayErr := a.ValidateGatewayResponse(res)
	if gatewayErr != nil {
		a.srv.Log.Warn("Gateway response validation failed!")
		errMsg = "Ops! Something went wrong."
	}

	return res, errMsg
}

func (a *ServicePayment) getPastTransactionToken(payment model.Payment, kind model.TransactionKind) (string, *model_helper.PaymentError, *model_helper.AppError) {
	transactions, appErr := a.TransactionsByOption(model_helper.PaymentTransactionFilterOpts{
		CommonQueryOptions: model_helper.NewCommonQueryOptions(
			model.PaymentTransactionWhere.PaymentID.EQ(payment.ID),
			model.PaymentTransactionWhere.Kind.EQ(kind),
			model.PaymentTransactionWhere.IsSuccess.EQ(true),
		),
	})
	if appErr != nil && appErr.StatusCode == http.StatusInternalServerError {
		return "", nil, appErr
	}

	if length := len(transactions); length == 0 {
		return "", model_helper.NewPaymentError("getPastTransactionToken", fmt.Sprintf("Cannot find successful %s transaction", kind), model_helper.NOT_FOUND), nil
	} else {
		return transactions[length-1].Token, nil, nil
	}
}

func (a *ServicePayment) validateRefundAmount(payment model.Payment, amount decimal.Decimal) *model_helper.PaymentError {
	if amount.LessThan(decimal.Zero) {
		return model_helper.NewPaymentError("validateRefundAmount", "Amount should be positive number", model_helper.INVALID)
	}
	if amount.GreaterThan(payment.CapturedAmount) {
		return model_helper.NewPaymentError("validateRefundAmount", "Cannot refund more than captures.", model_helper.INVALID)
	}

	return nil
}

func (a *ServicePayment) PaymentRefundOrVoid(dbTransaction boil.ContextTransactor, payment model.Payment, manager interfaces.PluginManagerInterface, channelSlug string) (*model_helper.PaymentError, *model_helper.AppError) {
	paymentCanVoid, appErr := a.PaymentCanVoid(payment)
	if appErr != nil {
		return nil, appErr
	}

	var paymentErr *model_helper.PaymentError

	if model_helper.PaymentCanRefund(payment) {
		_, paymentErr, appErr = a.Refund(dbTransaction, payment, manager, channelSlug, nil)
	} else if paymentCanVoid {
		_, paymentErr, appErr = a.Void(dbTransaction, payment, manager, channelSlug)
	}

	return paymentErr, appErr
}
