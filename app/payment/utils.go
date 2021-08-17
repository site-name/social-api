package payment

import (
	"context"
	"net/http"
	"strings"

	"github.com/site-name/decimal"
	goprices "github.com/site-name/go-prices"
	"github.com/sitename/sitename/app"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model/account"
	"github.com/sitename/sitename/model/checkout"
	"github.com/sitename/sitename/model/order"
	"github.com/sitename/sitename/model/payment"
	"github.com/sitename/sitename/modules/json"
	"github.com/sitename/sitename/modules/util"
	"github.com/sitename/sitename/web/graphql/gqlmodel"
)

// CreatePaymentInformation Extract order information along with payment details.
//
// Returns information required to process payment and additional
// billing/shipping addresses for optional fraud-prevention mechanisms.
func (a *AppPayment) CreatePaymentInformation(payMent *payment.Payment, paymentToken *string, amount *decimal.Decimal, customerId *string, storeSource bool, additionalData map[string]string) (*payment.PaymentData, *model.AppError) {
	var (
		billingAddress  *account.Address
		shippingAddress *account.Address
		amount_         = payMent.Total

		billingAddressID  string
		shippingAddressID string
		email             string = payMent.BillingEmail
		orderId           string
		customerIpAddress string
		appErr            *model.AppError
	)

	if amount != nil {
		amount_ = amount
	}

	// checks if payMent has checkout
	if payMent.CheckoutID != nil {
		checkout, appErr := a.app.CheckoutApp().CheckoutbyToken(*payMent.CheckoutID)
		if appErr != nil {
			return nil, appErr
		}

		// get checkout user
		if checkout.UserID != nil {
			user, appErr := a.app.AccountApp().UserById(context.Background(), *checkout.UserID)
			if appErr != nil {
				return nil, appErr
			}
			email = user.Email
		} else {
			email = checkout.Email
		}

		if checkout.BillingAddressID != nil && checkout.ShippingAddressID != nil {
			billingAddressID = *checkout.BillingAddressID
			shippingAddressID = *checkout.ShippingAddressID
		}
	} else if payMent.OrderID != nil { // checks if payMent has order
		order, appErr := a.app.OrderApp().OrderById(*payMent.OrderID)
		if appErr != nil {
			return nil, appErr
		}

		email = order.UserEmail
		orderId = order.Id

		if order.BillingAddressID != nil && order.ShippingAddressID != nil {
			billingAddressID = *order.BillingAddressID
			shippingAddressID = *order.ShippingAddressID
		}
	}

	var (
		billingAddressData  *payment.AddressData
		shippingAddressData *payment.AddressData
	)

	if model.IsValidId(billingAddressID) && model.IsValidId(shippingAddressID) {
		billingAddress, appErr = a.app.AccountApp().AddressById(billingAddressID)
		if appErr != nil {
			return nil, appErr
		}

		shippingAddress, appErr = a.app.AccountApp().AddressById(shippingAddressID)
		if appErr != nil {
			return nil, appErr
		}
	}

	if billingAddress != nil {
		billingAddressData = payment.AddressDataFromAddress(billingAddress)
	}
	if shippingAddress != nil {
		shippingAddressData = payment.AddressDataFromAddress(shippingAddress)
	}

	if payMent.CustomerIpAddress != nil {
		customerIpAddress = *payMent.CustomerIpAddress
	}

	return &payment.PaymentData{
		Gateway:           payMent.GateWay,
		Amount:            *amount_,
		Currency:          payMent.Currency,
		Billing:           billingAddressData,
		Shipping:          shippingAddressData,
		PaymentID:         payMent.Id,
		GraphqlPaymentID:  payMent.Id,
		OrderID:           orderId,
		CustomerIpAddress: customerIpAddress,
		CustomerEmail:     email,
		Token:             paymentToken,
		CustomerID:        customerId,
		ReuseSource:       storeSource,
		Data:              additionalData,
	}, nil
}

// CreatePayment Create a payment instance.
//
// This method is responsible for creating payment instances that works for
// both Django views and GraphQL mutations.
//
// NOTE: `customerIpAddress`, `paymentToken`, `returnUrl` and `externalReference` can be empty
//
// `extraData`, `ckout`, `ord` can be nil
func (a *AppPayment) CreatePayment(gateway string, total *decimal.Decimal, currency string, email string, customerIpAddress string, paymentToken string, extraData map[string]string, ckout *checkout.Checkout, ord *order.Order, returnUrl string, externalReference string) (*payment.Payment, *model.AppError) {
	// must at least provide either checkout or order, both is best :))
	if ckout == nil && ord == nil {
		return nil, model.NewAppError("CreatePayment", "app.payment.checkout_order_required.app_error", nil, "", http.StatusBadRequest)
	}

	if extraData == nil {
		extraData = make(map[string]string)
	}

	var (
		billingAddress   *account.Address
		billingAddressID string
	)

	if ckout != nil && ckout.BillingAddressID != nil {
		billingAddressID = *ckout.BillingAddressID
	} else if ord != nil && ord.BillingAddressID != nil {
		billingAddressID = *ord.BillingAddressID
	}

	if billingAddressID == "" || !model.IsValidId(billingAddressID) {
		return nil, model.NewAppError("CreatePayment", "app.payment.order_billing_address_not_set.app_error", nil, "", http.StatusBadRequest)
	}

	billingAddress, appErr := a.app.AccountApp().AddressById(billingAddressID)
	if appErr != nil {
		return nil, appErr
	}

	payment := &payment.Payment{
		BillingEmail:       email,
		BillingFirstName:   billingAddress.FirstName,
		BillingLastName:    billingAddress.LastName,
		BillingCompanyName: billingAddress.CompanyName,
		BillingAddress1:    billingAddress.StreetAddress1,
		BillingAddress2:    billingAddress.StreetAddress2,
		BillingCity:        billingAddress.City,
		BillingPostalCode:  billingAddress.PostalCode,
		BillingCountryCode: billingAddress.Country,
		BillingCountryArea: billingAddress.CountryArea,
		Currency:           currency,
		GateWay:            gateway,
		Total:              total,
		ReturnUrl:          &returnUrl,
		PspReference:       &externalReference,
		IsActive:           model.NewBool(true),
		CustomerIpAddress:  &customerIpAddress,
		ExtraData:          model.ModelToJson(extraData),
		Token:              paymentToken,
	}
	if ckout != nil {
		payment.CheckoutID = &ckout.Token
	}
	if ord != nil {
		payment.OrderID = &ord.Id
	}

	return a.app.PaymentApp().CreateOrUpdatePayment(payment)
}

func (a *AppPayment) GetAlreadyProcessedTransaction(paymentID string, gatewayResponse *payment.GatewayResponse) (*payment.PaymentTransaction, *model.AppError) {
	// get all transactions that belong to given payment
	trans, appErr := a.GetAllPaymentTransactions(paymentID)
	if appErr != nil {
		return nil, appErr
	}

	var processedTran *payment.PaymentTransaction

	// find the most recent transaction that satifies:
	for _, tran := range trans {
		if tran.IsSuccess == gatewayResponse.IsSucess &&
			tran.ActionRequired == gatewayResponse.ActionRequired &&
			tran.Token == gatewayResponse.TransactionID &&
			tran.Kind == gatewayResponse.Kind &&
			tran.Amount != nil && tran.Amount.Equal(gatewayResponse.Amount) &&
			tran.Currency == gatewayResponse.Currency {
			if processedTran == nil || tran.CreateAt > processedTran.CreateAt {
				processedTran = tran
			}
		}
	}

	return processedTran, nil
}

// CreateTransaction reate a transaction based on transaction kind and gateway response.
func (a *AppPayment) CreateTransaction(paymentID string, kind string, paymentInformation *payment.PaymentData, actionRequired bool, gatewayResponse *payment.GatewayResponse, errorMsg string, isSuccess bool) (*payment.PaymentTransaction, *model.AppError) {
	// Default values for token, amount, currency are only used in cases where
	// response from gateway was invalid or an exception occured
	if gatewayResponse == nil {
		var transactionId string
		if paymentInformation.Token != nil {
			transactionId = *paymentInformation.Token
		}
		gatewayResponse = &payment.GatewayResponse{
			Kind:           kind,
			ActionRequired: false,
			IsSucess:       isSuccess,
			TransactionID:  transactionId,
			Amount:         paymentInformation.Amount,
			Currency:       paymentInformation.Currency,
			Error:          errorMsg,
			RawResponse:    map[string]string{},
		}
	}

	tran := &payment.PaymentTransaction{
		PaymentID:          paymentID,
		ActionRequired:     actionRequired,
		Kind:               gatewayResponse.Kind,
		Token:              gatewayResponse.TransactionID,
		IsSuccess:          isSuccess,
		Amount:             &gatewayResponse.Amount,
		Currency:           gatewayResponse.Currency,
		Error:              &gatewayResponse.Error,
		CustomerID:         &gatewayResponse.CustomerID,
		GatewayResponse:    gatewayResponse.RawResponse,
		ActionRequiredData: gatewayResponse.ActionRequiredData,
	}

	return a.SaveTransaction(tran)
}

func (a *AppPayment) GetAlreadyProcessedTransactionOrCreateNewTransaction(paymentID, kind string, paymentInformation *payment.PaymentData, actionRequired bool, gatewayResponse *payment.GatewayResponse, errorMsg string) (*payment.PaymentTransaction, *model.AppError) {
	if gatewayResponse != nil && gatewayResponse.TransactionAlreadyProcessed {
		transaction, appErr := a.GetAlreadyProcessedTransaction(paymentID, gatewayResponse)
		if appErr == nil {
			return transaction, nil
		} else if appErr.StatusCode == http.StatusInternalServerError {
			// if error caused by internal server, still have to return it
			return nil, appErr
		}
	}

	return a.CreateTransaction(paymentID, kind, paymentInformation, actionRequired, gatewayResponse, errorMsg, false)
}

func (a *AppPayment) CleanCapture(pm *payment.Payment, amount decimal.Decimal) *model.AppError {
	if amount.LessThanOrEqual(decimal.Zero) {
		return model.NewAppError("CleanCapture", app.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "amount"}, "", http.StatusBadRequest)
	}
	if !pm.CanCapture() {
		return model.NewAppError("CleanCapture", "app.payment.cannot_capture.app_error", nil, "", http.StatusNotAcceptable)
	}
	// amount > payment's total || amount > payment's Total - payment's CapturedAmount
	if amount.GreaterThan(*pm.Total) || amount.GreaterThan((*pm.Total).Sub(*pm.CapturedAmount)) {
		return model.NewAppError("CleanCapture", "app.payment.un-captured_must_greater_than_charge.app_error", nil, "", http.StatusNotAcceptable)
	}

	return nil
}

func (a *AppPayment) CleanAuthorize(payment *payment.Payment) *model.AppError {
	if !payment.CanAuthorize() {
		return model.NewAppError("CleanAuthorize", "app.payment.cannot_authorized_again.app_error", nil, "", http.StatusNotAcceptable)
	}
	return nil
}

func (a *AppPayment) ValidateGatewayResponse(response *payment.GatewayResponse) *model.AppError {
	if response == nil {
		return model.NewAppError("ValidateGatewayResponse", "app.payment.argument_required.app_error", nil, "", http.StatusBadRequest)
	}

	// checks if response's Kind is valid transaction kind:
	if _, ok := payment.TransactionKindString[response.Kind]; !ok {
		validTransactionKinds := make([]string, len(payment.TransactionKindString))
		i := 0
		for key := range payment.TransactionKindString {
			validTransactionKinds[i] = key
			i++
		}

		return model.NewAppError("ValidateGatewayResponse",
			"app.payment.invalid_gateway_response_kind.app_error",
			map[string]interface{}{
				"ValidKinds": strings.Join(validTransactionKinds, ","),
			}, "",
			http.StatusNotAcceptable,
		)
	}

	// checks if response's RawResponse is json encodable
	_, err := json.JSON.Marshal(response.RawResponse)
	if err != nil {
		return model.NewAppError("", "app.payment.gateway_response_not_serializable.app_error", nil, err.Error(), http.StatusNotAcceptable)
	}

	return nil
}

func (a *AppPayment) GatewayPostProcess(transaction *payment.PaymentTransaction, payMent *payment.Payment) *model.AppError {
	if transaction == nil || payMent == nil {
		return model.NewAppError("GatewayPostProcess", app.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "transaction/payMent"}, "", http.StatusBadRequest)
	}

	changedFields := []string{}
	var appErr *model.AppError

	if !transaction.IsSuccess || transaction.AlreadyProcessed {
		if len(changedFields) > 0 {
			if _, appErr = a.CreateOrUpdatePayment(payMent); appErr != nil {
				return appErr
			}
		}
		return nil
	}

	if transaction.ActionRequired {
		payMent.ToConfirm = true
		changedFields = append(changedFields, "to_confirm")
		if _, appErr = a.CreateOrUpdatePayment(payMent); appErr != nil {
			return appErr
		}
	}

	// to_confirm is defined by the transaction.action_required. Payment doesn't
	// require confirmation when we got action_required == False
	if payMent.ToConfirm {
		payMent.ToConfirm = true
		changedFields = append(changedFields, "to_confirm")
	}

	switch transaction.Kind {
	case payment.CAPTURE, payment.REFUND_REVERSED:
		payMent.CapturedAmount = model.NewDecimal(payMent.CapturedAmount.Add(*transaction.Amount))
		payMent.IsActive = model.NewBool(true)
		// Set payment charge status to fully charged
		// only if there is no more amount needs to charge
		payMent.ChargeStatus = payment.PARTIALLY_CHARGED
		if payMent.GetChargeAmount().LessThanOrEqual(decimal.Zero) {
			payMent.ChargeStatus = payment.FULLY_CHARGED
		}
		changedFields = append(changedFields, "charge_status", "captured_amount", "update_at")
	case payment.VOID:
		payMent.IsActive = model.NewBool(false)
		changedFields = append(changedFields, "is_active", "update_at")
	case payment.REFUND:
		changedFields = append(changedFields, "captured_amount", "update_at")
		payMent.CapturedAmount = model.NewDecimal(payMent.CapturedAmount.Sub(*transaction.Amount))
		payMent.ChargeStatus = payment.PARTIALLY_REFUNDED
		if payMent.CapturedAmount.LessThanOrEqual(decimal.Zero) {
			payMent.CapturedAmount = &decimal.Zero
			payMent.ChargeStatus = payment.FULLY_REFUNDED
			payMent.IsActive = model.NewBool(false)
		}
	case payment.PENDING:
		payMent.ChargeStatus = payment.PENDING
		changedFields = append(changedFields, "charge_status")
	case payment.CANCEL:
		payMent.ChargeStatus = payment.CANCELLED
		payMent.IsActive = model.NewBool(false)
		changedFields = append(changedFields, "charge_status", "is_active")
	case payment.CAPTURE_FAILED:
		if payMent.ChargeStatus == payment.PARTIALLY_CHARGED || payMent.ChargeStatus == payment.FULLY_CHARGED {
			payMent.CapturedAmount = model.NewDecimal(payMent.CapturedAmount.Sub(*transaction.Amount))
			payMent.ChargeStatus = payment.PARTIALLY_CHARGED
			if payMent.CapturedAmount.LessThanOrEqual(decimal.Zero) {
				payMent.CapturedAmount = &decimal.Zero
			}
			changedFields = append(changedFields, "charge_status", "captured_amount", "update_at")
		}
	}

	if len(changedFields) > 0 {
		if _, appErr := a.CreateOrUpdatePayment(payMent); appErr != nil {
			return appErr
		}
	}

	transaction.AlreadyProcessed = true
	if _, appErr = a.UpdateTransaction(transaction); appErr != nil {
		return appErr
	}

	if util.StringInSlice("captured_amount", changedFields) && payMent.OrderID != nil {
		if appErr = a.app.OrderApp().UpdateOrderTotalPaid(*payMent.OrderID); appErr != nil {
			return appErr
		}
	}

	return nil
}

// FetchCustomerId
// user must be either: *model.User OR *gqlmodel.User
// returning string could be "" or long string
func (a *AppPayment) FetchCustomerId(user interface{}, gateway string) (string, *model.AppError) {
	// validate arguments are valid
	var argumentErrorFields string
	if user == nil {
		argumentErrorFields = "'user'"
	}
	if gateway == "" {
		argumentErrorFields += ", 'gateway'"
	}

	if argumentErrorFields != "" {
		return "", model.NewAppError("FetchCustomerId", app.InvalidArgumentAppErrorID,
			map[string]interface{}{
				"Fields": argumentErrorFields,
			}, "", http.StatusBadRequest,
		)
	}

	metaKey := prepareKeyForGatewayCustomerId(gateway)

	switch v := user.(type) {
	case *account.User:
		return v.ModelMetadata.GetValueFromMeta(metaKey, "", account.PrivateMetadata), nil
	case *gqlmodel.User:
		// create new ModelMetadata for concurrent accessing
		meta := &model.ModelMetadata{
			PrivateMetadata: gqlmodel.MetaDataToStringMap(v.PrivateMetadata),
		}
		return meta.GetValueFromMeta(metaKey, "", model.PrivateMetadata), nil
	default:
		return "", model.NewAppError("FetchCustomerId", app.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "user"}, "user param must be wither *account.User or *gqlmodel.User", http.StatusBadRequest)
	}
}

// StoreCustomerId stores new value into given user's PrivateMetadata
func (a *AppPayment) StoreCustomerId(userID string, gateway string, customerID string) *model.AppError {
	// validate arguments are valid:
	var argumentErrFields string
	if !model.IsValidId(userID) {
		argumentErrFields = "'userID'"
	}
	if trimmedGateway := strings.TrimSpace(gateway); trimmedGateway == "" || len(trimmedGateway) > payment.MAX_LENGTH_PAYMENT_GATEWAY {
		argumentErrFields += ", 'gateway'"
	}
	if trimmedCustomerID := strings.TrimSpace(customerID); trimmedCustomerID == "" || len(trimmedCustomerID) > payment.TRANSACTION_CUSTOMER_ID_MAX_LENGTH {
		argumentErrFields += ", 'customerID'"
	}

	if argumentErrFields != "" {
		return model.NewAppError(
			"StoreCustomerId", app.InvalidArgumentAppErrorID,
			map[string]interface{}{
				"Fields": argumentErrFields,
			}, "", http.StatusBadRequest,
		)
	}

	metaKey := prepareKeyForGatewayCustomerId(gateway)
	user, appErr := a.app.AccountApp().UserById(context.Background(), userID)
	if appErr != nil {
		return appErr
	}
	user.StoreValueInMeta(
		map[string]string{
			metaKey: customerID,
		},
		account.PrivateMetadata,
	)
	_, appErr = a.app.AccountApp().UpdateUser(user, false)
	if appErr != nil {
		return appErr
	}

	return nil
}

// prepareKeyForGatewayCustomerId just trims spaces, upper then concatenates ".customer_id" to given `gatewayName`.
//
//  strings.TrimSpace(strings.ToUpper(gatewayName)) + ".customer_id"
func prepareKeyForGatewayCustomerId(gatewayName string) string {
	return strings.TrimSpace(strings.ToUpper(gatewayName)) + ".customer_id"
}

func (a *AppPayment) UpdatePayment(pm *payment.Payment, gatewayResponse *payment.GatewayResponse) *model.AppError {

	var changed bool
	var appErr *model.AppError

	if gatewayResponse.PspReference != "" {
		pm.PspReference = &gatewayResponse.PspReference
		changed = true
	}

	if gatewayResponse.PaymentMethodInfo != nil {
		if brand := gatewayResponse.PaymentMethodInfo.Brand; brand != "" {
			pm.CcBrand = brand
			changed = true
		}
		if last4 := gatewayResponse.PaymentMethodInfo.Last4; last4 != "" {
			pm.CcLastDigits = last4
			changed = true
		}
		if expYear := gatewayResponse.PaymentMethodInfo.ExpYear; expYear > 0 {
			pm.CcExpYear = &expYear
			changed = true
		}
		if expMonth := gatewayResponse.PaymentMethodInfo.ExpMonth; expMonth > 0 {
			pm.CcExpMonth = &expMonth
			changed = true
		}
		if type_ := gatewayResponse.PaymentMethodInfo.Type; type_ != "" {
			pm.PaymentMethodType = type_
			changed = true
		}
	}

	if changed {
		_, appErr = a.CreateOrUpdatePayment(pm)
	}

	return appErr
}

// Convert minor unit (smallest unit of currency) to decimal value.
//
// (value: 1000, currency: USD) will be converted to 10.00
func PriceFromMinorUnit(value string, currency string) (*decimal.Decimal, error) {
	d, err := decimal.NewFromString(value)
	if err != nil {
		return nil, err
	}

	precision, err := goprices.GetCurrencyPrecision(currency)
	if err != nil {
		return nil, err
	}

	d = d.
		Mul(
			decimal.
				NewFromInt32(10).
				Pow(decimal.NewFromInt32(-int32(precision))),
		).
		Round(int32(precision))

	return &d, nil
}

// Convert decimal value to the smallest unit of currency.
//
// Take the value, discover the precision of currency and multiply value by
// Decimal('10.0'), then change quantization to remove the comma.
// Decimal(10.0) -> str(1000)
func PriceToMinorUnit(value *decimal.Decimal, currency string) (string, error) {
	precision, err := goprices.GetCurrencyPrecision(currency)
	if err != nil {
		return "", err
	}

	return value.
		Mul(
			decimal.
				NewFromFloat(10.0).
				Pow(decimal.NewFromInt32(int32(precision))),
		).
		String(), nil
}

// IsCurrencySupported checks if given currency is supported by system
// TODO: implement me
func IsCurrencySupported() bool {
	panic("not implemented")
}
