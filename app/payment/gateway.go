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
	"github.com/sitename/sitename/model/payment"
)

const (
	ERROR_MSG                 = "Oops! Something went wrong"
	GENERIC_TRANSACTION_ERROR = "Transaction was unsuccessful."
)

// raisePaymentError must be called right before function returns
func (a *ServicePayment) raisePaymentError(where string, transaction payment.PaymentTransaction) *payment.PaymentError {
	if !transaction.IsSuccess {
		msg := GENERIC_TRANSACTION_ERROR
		if transaction.Error != nil {
			msg = *transaction.Error
		}
		return payment.NewPaymentError(where, msg, payment.INVALID)
	}

	return nil
}

// paymentPostProcess must be called right before function returns
func (a *ServicePayment) paymentPostProcess(transaction payment.PaymentTransaction) *model.AppError {
	payMent, appErr := a.PaymentByID(nil, transaction.PaymentID, false)
	if appErr != nil {
		return appErr
	}
	return a.GatewayPostProcess(transaction, payMent)
}

// requireActivePayment must be called in the beginning of the function body
func (a *ServicePayment) requireActivePayment(where string, payMent payment.Payment) *payment.PaymentError {
	if payMent.IsActive == nil || !*payMent.IsActive {
		return payment.NewPaymentError(where, "This payment is no longer active", payment.INVALID)
	}
	return nil
}

// withLockedPayment Lock payment to protect from asynchronous modification.
func (a *ServicePayment) withLockedPayment(where string, payMent payment.Payment) (*payment.Payment, *model.AppError) {
	paymentToOperateOn, appErr := a.srv.PaymentService().PaymentByID(nil, payMent.Id, true)
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
	payMent payment.Payment,
	token string,
	manager interfaces.PluginManagerInterface,
	channelID string, // originally is channelSlug in saleor
	customerID *string,
	storeSource bool,
	additionalData map[string]interface{},
) (*payment.PaymentTransaction, *payment.PaymentError, *model.AppError) {

	paymentErr := a.requireActivePayment("ProcessPayment", payMent)
	if paymentErr != nil {
		return nil, paymentErr, nil
	}

	lockedPayment, appErr := a.withLockedPayment("ProcessPayment", payMent)
	if appErr != nil {
		return nil, nil, appErr
	}

	paymentData, appErr := a.CreatePaymentInformation(lockedPayment, &token, nil, customerID, storeSource, additionalData)
	if appErr != nil {
		return nil, nil, appErr
	}

	response, errMsg := a.fetchGatewayResponse(manager.ProcessPayment, lockedPayment.GateWay, *paymentData, channelID)
	actionRequired := response != nil && response.ActionRequired

	if response != nil {
		appErr = a.UpdatePayment(*lockedPayment, response)
		if appErr != nil {
			return nil, nil, appErr
		}
	}

	paymentTransaction, appErr := a.GetAlreadyProcessedTransactionOrCreateNewTransaction(lockedPayment.Id, payment.CAPTURE, paymentData, actionRequired, response, errMsg)
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
	payMent payment.Payment,
	token string,
	manager interfaces.PluginManagerInterface,
	channelID string,
	customerID *string,
	storeSource bool,

) (*payment.PaymentTransaction, *payment.PaymentError, *model.AppError) {

	paymentErr := a.requireActivePayment("Authorize", payMent)
	if paymentErr != nil {
		return nil, paymentErr, nil
	}

	lockedPayment, appErr := a.withLockedPayment("Authorize", payMent)
	if appErr != nil {
		return nil, nil, appErr
	}

	paymentErr = a.CleanAuthorize(lockedPayment)
	if paymentErr != nil {
		return nil, paymentErr, nil
	}

	paymentData, appErr := a.CreatePaymentInformation(lockedPayment, &token, nil, customerID, storeSource, nil)
	if appErr != nil {
		return nil, nil, appErr
	}

	response, errMsg := a.fetchGatewayResponse(manager.AuthorizePayment, lockedPayment.GateWay, *paymentData, channelID)
	if response != nil {
		appErr = a.UpdatePayment(*lockedPayment, response)
		if appErr != nil {
			return nil, nil, appErr
		}
	}

	paymentTransaction, appErr := a.GetAlreadyProcessedTransactionOrCreateNewTransaction(lockedPayment.Id, payment.CAPTURE, paymentData, false, response, errMsg)
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
	payMent payment.Payment,
	manager interfaces.PluginManagerInterface,
	channelID string,
	amount *decimal.Decimal, // can be nil
	customerID *string, // can be nil
	storeSource bool, // default false

) (*payment.PaymentTransaction, *payment.PaymentError, *model.AppError) {

	paymentErr := a.requireActivePayment("Capture", payMent)
	if paymentErr != nil {
		return nil, paymentErr, nil
	}

	lockedPayment, appErr := a.withLockedPayment("Capture", payMent)
	if appErr != nil {
		return nil, nil, appErr
	}

	if amount == nil {
		amount = lockedPayment.GetChargeAmount()
	}

	paymentErr = a.CleanCapture(lockedPayment, *amount)
	if paymentErr != nil {
		return nil, paymentErr, nil
	}

	paymentData, appErr := a.CreatePaymentInformation(lockedPayment, &lockedPayment.Token, amount, customerID, storeSource, nil)
	if appErr != nil {
		return nil, nil, appErr
	}

	response, errMsg := a.fetchGatewayResponse(manager.CapturePayment, lockedPayment.GateWay, *paymentData, channelID)
	if response != nil {
		appErr = a.UpdatePayment(*lockedPayment, response)
		if appErr != nil {
			return nil, nil, appErr
		}
	}

	paymentTransaction, appErr := a.GetAlreadyProcessedTransactionOrCreateNewTransaction(lockedPayment.Id, payment.CAPTURE, paymentData, false, response, errMsg)
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
	payMent payment.Payment,
	manager interfaces.PluginManagerInterface,
	channelID string,
	amount *decimal.Decimal, // can be nil

) (*payment.PaymentTransaction, *payment.PaymentError, *model.AppError) {

	paymentErr := a.requireActivePayment("Refund", payMent)
	if paymentErr != nil {
		return nil, paymentErr, nil
	}

	lockedPayment, appErr := a.withLockedPayment("Refund", payMent)
	if appErr != nil {
		return nil, nil, appErr
	}

	if amount == nil {
		amount = lockedPayment.CapturedAmount
	}

	paymentErr = a.validateRefundAmount(lockedPayment, amount)
	if paymentErr != nil {
		return nil, paymentErr, nil
	}

	if !lockedPayment.CanRefund() {
		return nil, payment.NewPaymentError("Refund", "This payment cannot be refunded", payment.INVALID), nil
	}

	kind := payment.CAPTURE
	if lockedPayment.IsManual() {
		kind = payment.EXTERNAL
	}

	token, paymentErr, appErr := a.getPastTransactionToken(lockedPayment, kind)
	if paymentErr != nil || appErr != nil {
		return nil, paymentErr, appErr
	}

	paymentData, appErr := a.CreatePaymentInformation(lockedPayment, &token, amount, nil, false, nil)
	if appErr != nil {
		return nil, nil, appErr
	}

	var paymentTransaction *payment.PaymentTransaction

	if lockedPayment.IsManual() {
		// for manual payment we just need to mark payment as a refunded
		paymentTransaction, appErr = a.CreateTransaction(lockedPayment.Id, payment.REFUND, paymentData, false, nil, "", true)
		if appErr != nil {
			return nil, nil, appErr
		}
	} else {
		response, errMsg := a.fetchGatewayResponse(manager.RefundPayment, lockedPayment.GateWay, *paymentData, channelID)
		paymentTransaction, appErr = a.GetAlreadyProcessedTransactionOrCreateNewTransaction(
			lockedPayment.Id,
			payment.REFUND,
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
func (a *ServicePayment) Void(payMent payment.Payment, manager interfaces.PluginManagerInterface, channelID string) (*payment.PaymentTransaction, *payment.PaymentError, *model.AppError) {
	paymentErr := a.requireActivePayment("Refund", payMent)
	if paymentErr != nil {
		return nil, paymentErr, nil
	}

	lockedPayment, appErr := a.withLockedPayment("Refund", payMent)
	if appErr != nil {
		return nil, nil, appErr
	}

	token, paymentErr, appErr := a.getPastTransactionToken(lockedPayment, payment.AUTH)
	if paymentErr != nil || appErr != nil {
		return nil, paymentErr, appErr
	}

	paymentData, appErr := a.CreatePaymentInformation(lockedPayment, &token, nil, nil, false, nil)
	if appErr != nil {
		return nil, nil, appErr
	}

	response, errMsg := a.fetchGatewayResponse(manager.VoidPayment, lockedPayment.GateWay, *paymentData, channelID)
	paymentTransaction, appErr := a.GetAlreadyProcessedTransactionOrCreateNewTransaction(
		lockedPayment.Id,
		payment.VOID,
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
	payMent payment.Payment,
	manager interfaces.PluginManagerInterface,
	channelID string,
	additionalData map[string]interface{}, // can be none

) (*payment.PaymentTransaction, *payment.PaymentError, *model.AppError) {
	panic("not implt")
}

func (a *ServicePayment) ListPaymentSources(
	gateway string,
	customerID string,
	manager interfaces.PluginManagerInterface,
	channelID string,

) ([]*payment.CustomerSource, *model.AppError) {

	source, err := manager.ListPaymentSources(gateway, customerID, channelID)
	if err != nil {
		return nil, model.NewAppError("ListPaymentSources", "app.payment.error_listing_payment_sources.app_error", nil, err.Error(), http.StatusInternalServerError)
	}
	return source, nil
}

func (a *ServicePayment) ListGateways(
	manager interfaces.PluginManagerInterface,
	channelID string,
) []*payment.PaymentGateway {

	return manager.ListPaymentGateways("", nil, channelID, true)
}

func (a *ServicePayment) fetchGatewayResponse(paymentFunc interfaces.PaymentMethod, gateway string, paymentData payment.PaymentData, channelID string) (res *payment.GatewayResponse, errMsg string) {
	res, _ = paymentFunc(gateway, paymentData, channelID)
	gatewayErr := a.ValidateGatewayResponse(res)
	if gatewayErr != nil {
		a.srv.Log.Warn("Gateway response validation failed!")
		errMsg = "Ops! Something went wrong."
	}

	return res, errMsg
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
func (a *ServicePayment) PaymentRefundOrVoid(payMent *payment.Payment, manager interfaces.PluginManagerInterface, channelSlug string) (*payment.PaymentError, *model.AppError) {
	if payMent == nil {
		return nil, nil
	}

	paymentCanVoid, appErr := a.srv.PaymentService().PaymentCanVoid(payMent)
	if appErr != nil {
		return nil, appErr
	}

	var paymentErr *payment.PaymentError

	if payMent.CanRefund() {
		_, paymentErr, appErr = a.Refund(*payMent, manager, channelSlug, nil)
	} else if paymentCanVoid {
		_, paymentErr, appErr = a.Void(*payMent, manager, channelSlug)
	}

	return paymentErr, appErr
}
