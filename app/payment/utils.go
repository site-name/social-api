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
)

// CreatePaymentInformation Extract order information along with payment details.
//
// Returns information required to process payment and additional
// billing/shipping addresses for optional fraud-prevention mechanisms.
func (a *ServicePayment) CreatePaymentInformation(payMent *payment.Payment, paymentToken *string, amount *decimal.Decimal, customerId *string, storeSource bool, additionalData map[string]interface{}) (*payment.PaymentData, *model.AppError) {

	var (
		billingAddressID  string
		shippingAddressID string
		billingAddress    *account.Address
		shippingAddress   *account.Address
		email             string
		userID            *string
	)

	if payMent.CheckoutID != nil {
		checkoutOfPayment, appErr := a.srv.CheckoutService().CheckoutByOption(&checkout.CheckoutFilterOption{
			Token: &model.StringFilter{
				StringOption: &model.StringOption{
					Eq: *payMent.CheckoutID,
				},
			},
		})
		if appErr != nil && appErr.StatusCode == http.StatusInternalServerError {
			return nil, appErr // ignore not found error
		}

		if checkoutOfPayment != nil {
			if checkoutOfPayment.BillingAddressID != nil {
				billingAddressID = *checkoutOfPayment.BillingAddressID
			}
			if checkoutOfPayment.ShippingAddressID != nil {
				shippingAddressID = *checkoutOfPayment.ShippingAddressID
			}
			emailOfCheckoutUser, appErr := a.srv.CheckoutService().GetCustomerEmail(checkoutOfPayment)
			if appErr != nil { // this is system caused error
				return nil, appErr
			}
			email = emailOfCheckoutUser
			if checkoutOfPayment.UserID != nil {
				userID = checkoutOfPayment.UserID
			}
		}
	} else if payMent.OrderID != nil {
		orderOfPayment, appErr := a.srv.OrderService().OrderById(*payMent.OrderID)
		if appErr != nil && appErr.StatusCode == http.StatusInternalServerError {
			return nil, appErr
		}

		if orderOfPayment != nil {
			if orderOfPayment.BillingAddressID != nil {
				billingAddressID = *orderOfPayment.BillingAddressID
			}
			if orderOfPayment.ShippingAddressID != nil {
				shippingAddressID = *orderOfPayment.ShippingAddressID
			}
			email = orderOfPayment.UserEmail
			if orderOfPayment.UserID != nil {
				userID = orderOfPayment.UserID
			}
		}
	} else {
		email = payMent.BillingEmail
	}

	if billingAddressID != "" || shippingAddressID != "" {
		addresses, appErr := a.srv.AccountService().AddressesByOption(&account.AddressFilterOption{
			Id: &model.StringFilter{
				StringOption: &model.StringOption{
					In: []string{billingAddressID, shippingAddressID},
				},
			},
		})
		if appErr.StatusCode == http.StatusInternalServerError {
			return nil, appErr
		}

		if len(addresses) > 0 {
			if addresses[0].Id == billingAddressID {
				billingAddress = addresses[0]
			} else {
				shippingAddress = addresses[0]
			}
		}
		if len(addresses) > 1 {
			if addresses[0].Id == billingAddressID {
				billingAddress = addresses[0]
				shippingAddress = addresses[1]
			} else {
				billingAddress = addresses[1]
				shippingAddress = addresses[0]
			}
		}
	}

	var (
		billingAddressData  *payment.AddressData
		shippingAddressData *payment.AddressData
	)
	if billingAddress != nil {
		billingAddressData = payment.AddressDataFromAddress(billingAddress)
	}
	if shippingAddress != nil {
		shippingAddressData = payment.AddressDataFromAddress(shippingAddress)
	}

	var orderID *string
	if payMent.OrderID != nil {
		orderID = payMent.OrderID
	}
	if amount == nil {
		amount = payMent.Total
	}
	if additionalData == nil {
		additionalData = make(map[string]interface{})
	}

	return &payment.PaymentData{
		Gateway:           payMent.GateWay,
		Token:             paymentToken,
		Amount:            *amount,
		Currency:          payMent.Currency,
		Billing:           billingAddressData,
		Shipping:          shippingAddressData,
		OrderID:           orderID,
		PaymentID:         payMent.Token,
		GraphqlPaymentID:  payMent.Token,
		CustomerIpAddress: payMent.CustomerIpAddress,
		CustomerID:        customerId,
		CustomerEmail:     email,
		ReuseSource:       storeSource,
		Data:              additionalData,
		GraphqlCustomerID: userID,
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
func (a *ServicePayment) CreatePayment(
	gateway string,
	total *decimal.Decimal,
	currency string,
	email string,
	customerIpAddress string,
	paymentToken string,
	extraData map[string]string,
	ckout *checkout.Checkout,
	ord *order.Order,
	returnUrl string,
	externalReference string,

) (*payment.Payment, *payment.PaymentError, *model.AppError) {
	// must at least provide either checkout or order, both is best :))
	if ckout == nil && ord == nil {
		return nil, nil, model.NewAppError("CreatePayment", app.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "order/checkout"}, "", http.StatusBadRequest)
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

	billingAddress, appErr := a.srv.AccountService().AddressById(billingAddressID)
	if appErr != nil {
		return nil, nil, appErr // this error can be either system error/not found error
	}

	if billingAddress == nil {
		return nil, payment.NewPaymentError("CreatePayment", "Order does not have a billing address.", payment.BILLING_ADDRESS_NOT_SET), nil
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

	payment, appErr = a.srv.PaymentService().UpsertPayment(payment)
	return payment, nil, appErr
}

func (a *ServicePayment) GetAlreadyProcessedTransaction(paymentID string, gatewayResponse *payment.GatewayResponse) (*payment.PaymentTransaction, *model.AppError) {
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
			if processedTran == nil || tran.CreateAt >= processedTran.CreateAt { // this find the most recent
				processedTran = tran
			}
		}
	}

	return processedTran, nil
}

// CreateTransaction reate a transaction based on transaction kind and gateway response.
func (a *ServicePayment) CreateTransaction(paymentID string, kind string, paymentInformation *payment.PaymentData, actionRequired bool, gatewayResponse *payment.GatewayResponse, errorMsg string, isSuccess bool) (*payment.PaymentTransaction, *model.AppError) {
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

func (a *ServicePayment) GetAlreadyProcessedTransactionOrCreateNewTransaction(paymentID, kind string, paymentInformation *payment.PaymentData, actionRequired bool, gatewayResponse *payment.GatewayResponse, errorMsg string) (*payment.PaymentTransaction, *model.AppError) {
	if gatewayResponse != nil && gatewayResponse.TransactionAlreadyProcessed {
		transaction, appErr := a.GetAlreadyProcessedTransaction(paymentID, gatewayResponse)
		if appErr == nil {
			return transaction, nil
		}
		if appErr.StatusCode == http.StatusInternalServerError { // ignore not found error
			return nil, appErr
		}
	}

	return a.CreateTransaction(paymentID, kind, paymentInformation, actionRequired, gatewayResponse, errorMsg, false)
}

// CleanCapture Check if payment can be captured.
func (a *ServicePayment) CleanCapture(pm *payment.Payment, amount decimal.Decimal) *payment.PaymentError {
	if amount.LessThanOrEqual(decimal.Zero) {
		return payment.NewPaymentError("CleanCapture", "Amount should be a positive number.", payment.INVALID)
	}
	if !pm.CanCapture() {
		return payment.NewPaymentError("CleanCapture", "This payment cannot be captured.", payment.INVALID)
	}
	// amount > payment's total || amount > payment's Total - payment's CapturedAmount
	if amount.GreaterThan(*pm.Total) || amount.GreaterThan((*pm.Total).Sub(*pm.CapturedAmount)) {
		return payment.NewPaymentError("CleanCapture", "Unable to charge more than un-captured amount.", payment.INVALID)
	}

	return nil
}

// CleanAuthorize Check if payment can be authorized
func (a *ServicePayment) CleanAuthorize(payMent *payment.Payment) *payment.PaymentError {
	if !payMent.CanAuthorize() {
		return payment.NewPaymentError("CleanAuthorize", "Charged transactions cannot be authorized again.", payment.INVALID)
	}
	return nil
}

// ValidateGatewayResponse Validate response to be a correct format for Saleor to process.
func (a *ServicePayment) ValidateGatewayResponse(response *payment.GatewayResponse) *payment.GatewayError {
	if response == nil {
		return &payment.GatewayError{
			Where:   "ValidateGatewayResponse",
			Message: "Gateway needs to return a GatewayResponse obj",
		}
	}

	// checks if response's Kind is valid transaction kind:
	if _, ok := payment.TransactionKindString[response.Kind]; !ok {
		validTransactionKinds := []string{}
		for key := range payment.TransactionKindString {
			validTransactionKinds = append(validTransactionKinds, key)
		}

		return &payment.GatewayError{
			Where:   "ValidateGatewayResponse",
			Message: "Gateway response kind must be one of " + strings.Join(validTransactionKinds, ", "),
		}
	}

	// checks if response's RawResponse is json encodable
	_, err := json.JSON.Marshal(response.RawResponse)
	if err != nil {
		return &payment.GatewayError{
			Where:   "ValidateGatewayResponse",
			Message: "Gateway response needs to be json serializable",
		}
	}

	return nil
}

// GatewayPostProcess
func (a *ServicePayment) GatewayPostProcess(paymentTransaction *payment.PaymentTransaction, payMent *payment.Payment) *model.AppError {
	tx, err := a.srv.Store.GetMaster().Begin()
	if err != nil {
		return model.NewAppError("GatewayPostProcess", app.ErrorCreatingTransactionErrorID, nil, err.Error(), http.StatusInternalServerError)
	}
	defer a.srv.Store.FinalizeTransaction(tx)

	changedFields := []string{}
	var appErr *model.AppError

	if !paymentTransaction.IsSuccess || paymentTransaction.AlreadyProcessed {
		if len(changedFields) > 0 {
			if _, appErr = a.UpsertPayment(payMent); appErr != nil {
				return appErr
			}
		}
		return nil
	}

	if paymentTransaction.ActionRequired {
		payMent.ToConfirm = true
		changedFields = append(changedFields, "to_confirm")
		if _, appErr = a.UpsertPayment(payMent); appErr != nil {
			return appErr
		}
	}

	// to_confirm is defined by the paymentTransaction.action_required. Payment doesn't
	// require confirmation when we got action_required == False
	if payMent.ToConfirm {
		payMent.ToConfirm = true
		changedFields = append(changedFields, "to_confirm")
	}

	switch paymentTransaction.Kind {
	case payment.CAPTURE, payment.REFUND_REVERSED:
		payMent.CapturedAmount = model.NewDecimal(payMent.CapturedAmount.Add(*paymentTransaction.Amount))
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
		payMent.CapturedAmount = model.NewDecimal(payMent.CapturedAmount.Sub(*paymentTransaction.Amount))
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
			payMent.CapturedAmount = model.NewDecimal(payMent.CapturedAmount.Sub(*paymentTransaction.Amount))
			payMent.ChargeStatus = payment.PARTIALLY_CHARGED
			if payMent.CapturedAmount.LessThanOrEqual(decimal.Zero) {
				payMent.CapturedAmount = &decimal.Zero
			}
			changedFields = append(changedFields, "charge_status", "captured_amount", "update_at")
		}
	}

	if len(changedFields) > 0 {
		if _, appErr := a.UpsertPayment(payMent); appErr != nil {
			return appErr
		}
	}

	paymentTransaction.AlreadyProcessed = true
	if _, appErr = a.UpdateTransaction(paymentTransaction); appErr != nil {
		return appErr
	}

	if util.StringInSlice("captured_amount", changedFields) && payMent.OrderID != nil {
		if appErr = a.srv.OrderService().UpdateOrderTotalPaid(tx, *payMent.OrderID); appErr != nil {
			return appErr
		}
	}

	if err = tx.Commit(); err != nil {
		return model.NewAppError("GatewayPostProcess", app.ErrorCommittingTransactionErrorID, nil, err.Error(), http.StatusInternalServerError)
	}

	return nil
}

// FetchCustomerId Retrieve users customer_id stored for desired gateway.
// returning string could be "" or long string
func (a *ServicePayment) FetchCustomerId(user *account.User, gateway string) (string, *model.AppError) {
	// validate arguments are valid
	var argumentErrorFields string
	if user == nil {
		argumentErrorFields = "'user'"
	}
	if gateway == "" {
		argumentErrorFields += ", 'gateway'"
	}

	if argumentErrorFields != "" {
		return "", model.NewAppError("FetchCustomerId", app.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": argumentErrorFields}, "", http.StatusBadRequest)
	}

	metaKey := prepareKeyForGatewayCustomerId(gateway)
	return user.ModelMetadata.GetValueFromMeta(metaKey, "", account.PrivateMetadata), nil
}

// StoreCustomerId stores new value into given user's PrivateMetadata
func (a *ServicePayment) StoreCustomerId(userID string, gateway string, customerID string) *model.AppError {
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
		return model.NewAppError("StoreCustomerId", app.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": argumentErrFields}, "", http.StatusBadRequest)
	}

	metaKey := prepareKeyForGatewayCustomerId(gateway)
	user, appErr := a.srv.AccountService().UserById(context.Background(), userID)
	if appErr != nil {
		return appErr
	}
	user.StoreValueInMeta(
		map[string]string{
			metaKey: customerID,
		},
		account.PrivateMetadata,
	)
	_, appErr = a.srv.AccountService().UpdateUser(user, false)
	return appErr
}

// prepareKeyForGatewayCustomerId just trims spaces, upper then concatenates ".customer_id" to given `gatewayName`.
//
//  strings.TrimSpace(strings.ToUpper(gatewayName)) + ".customer_id"
func prepareKeyForGatewayCustomerId(gatewayName string) string {
	return strings.TrimSpace(strings.ToUpper(gatewayName)) + ".customer_id"
}

func (a *ServicePayment) UpdatePayment(pm *payment.Payment, gatewayResponse *payment.GatewayResponse) *model.AppError {

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
		_, appErr = a.UpsertPayment(pm)
	}

	return appErr
}

func (a *ServicePayment) updatePaymentMethodDetails(payMent *payment.Payment, gatewayResponse *payment.GatewayResponse, changedFields []string) {
	if changedFields == nil {
		changedFields = []string{}
	}

	if gatewayResponse.PaymentMethodInfo == nil {
		return
	}

	if brand := gatewayResponse.PaymentMethodInfo.Brand; brand != "" {
		payMent.CcBrand = brand
	}
	if last4 := gatewayResponse.PaymentMethodInfo.Last4; last4 != "" {
		payMent.CcLastDigits = last4
	}
	if exprYear := gatewayResponse.PaymentMethodInfo.ExpYear; exprYear != 0 {
		payMent.CcExpYear = &exprYear
	}
	if exprMonth := gatewayResponse.PaymentMethodInfo.ExpMonth; exprMonth != 0 {
		payMent.CcExpMonth = &exprMonth
	}
	if type_ := gatewayResponse.PaymentMethodInfo.Type; type_ != "" {
		payMent.PaymentMethodType = type_
	}
}

func (a *ServicePayment) GetPaymentToken(payMent *payment.Payment) (string, *payment.PaymentError, *model.AppError) {
	authTransactions, appErr := a.TransactionsByOption(&payment.PaymentTransactionFilterOpts{
		Kind: &model.StringFilter{
			StringOption: &model.StringOption{
				Eq: payment.AUTH,
			},
		},
		IsSuccess: model.NewBool(true),
	})
	if appErr != nil {
		if appErr.StatusCode == http.StatusInternalServerError {
			return "", nil, appErr
		}
		return "", payment.NewPaymentError("GetPaymentToken", "Cannot process unauthorized transaction", payment.INVALID), appErr
	}

	return authTransactions[0].Token, nil, nil
}

// IsCurrencySupported checks if given currency is supported by system
func (a *ServicePayment) IsCurrencySupported(currency string, gatewayID string, manager interface{}) (bool, *model.AppError) {
	panic("not implemented")
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
