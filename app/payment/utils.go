package payment

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/mattermost/squirrel"
	"github.com/site-name/decimal"
	goprices "github.com/site-name/go-prices"
	"github.com/sitename/sitename/app/plugin/interfaces"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model_helper"
	"github.com/sitename/sitename/modules/model_types"
	"github.com/sitename/sitename/modules/util"
	"github.com/volatiletech/sqlboiler/v4/boil"
)

// CreatePaymentInformation Extract order information along with payment details.
//
// Returns information required to process payment and additional
// billing/shipping addresses for optional fraud-prevention mechanisms.
func (a *ServicePayment) CreatePaymentInformation(payment model.Payment, paymentToken *string, amount *decimal.Decimal, customerId *string, storeSource bool, additionalData map[string]any) (*model_helper.PaymentData, *model_helper.AppError) {
	var (
		billingAddressID  string
		shippingAddressID string
		billingAddress    *model.Address
		shippingAddress   *model.Address
		email             string
		userID            *string
	)

	if !payment.CheckoutID.IsNil() {
		checkoutOfPayment, appErr := a.srv.Checkout.CheckoutByOption(model_helper.CheckoutFilterOptions{
			CommonQueryOptions: model_helper.NewCommonQueryOptions(
				model.CheckoutWhere.Token.EQ(*payment.CheckoutID.String),
			),
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
			emailOfCheckoutUser, appErr := a.srv.Checkout.GetCustomerEmail(checkoutOfPayment)
			if appErr != nil { // this is system caused error
				return nil, appErr
			}
			email = emailOfCheckoutUser
			if checkoutOfPayment.UserID != nil {
				userID = checkoutOfPayment.UserID
			}
		}
	} else if !payment.OrderID.IsNil() {
		orderOfPayment, appErr := a.srv.Order.OrderById(*payment.OrderID.String)
		if appErr != nil && appErr.StatusCode == http.StatusInternalServerError {
			return nil, appErr
		}

		if orderOfPayment != nil {
			if !orderOfPayment.BillingAddressID.IsNil() {
				billingAddressID = *orderOfPayment.BillingAddressID.String
			}
			if !orderOfPayment.ShippingAddressID.IsNil() {
				shippingAddressID = *orderOfPayment.ShippingAddressID.String
			}
			email = orderOfPayment.UserEmail
			if !orderOfPayment.UserID.IsNil() {
				userID = orderOfPayment.UserID.String
			}
		}
	} else {
		email = payment.BillingEmail
	}

	if billingAddressID != "" || shippingAddressID != "" {
		addresses, appErr := a.srv.Account.AddressesByOption(model_helper.AddressFilterOptions{
			CommonQueryOptions: model_helper.NewCommonQueryOptions(
				model.AddressWhere.ID.IN([]string{billingAddressID, shippingAddressID}),
			),
		})
		if appErr.StatusCode == http.StatusInternalServerError {
			return nil, appErr
		}

		if len(addresses) > 0 {
			if addresses[0].ID == billingAddressID {
				billingAddress = addresses[0]
			} else {
				shippingAddress = addresses[0]
			}
		}
		if len(addresses) > 1 {
			if addresses[0].ID == billingAddressID {
				billingAddress = addresses[0]
				shippingAddress = addresses[1]
			} else {
				billingAddress = addresses[1]
				shippingAddress = addresses[0]
			}
		}
	}

	var (
		billingAddressData  *model_helper.AddressData
		shippingAddressData *model_helper.AddressData
	)
	if billingAddress != nil {
		billingAddressData = model_helper.AddressDataFromAddress(billingAddress)
	}
	if shippingAddress != nil {
		shippingAddressData = model_helper.AddressDataFromAddress(shippingAddress)
	}

	var orderID *string
	if !payment.OrderID.IsNil() {
		orderID = payment.OrderID.String
	}
	if amount == nil {
		amount = &payment.Total
	}
	if additionalData == nil {
		additionalData = make(map[string]any)
	}

	return &model_helper.PaymentData{
		Gateway:            payment.Gateway,
		Token:              paymentToken,
		Amount:             *amount,
		Currency:           payment.Currency.String(),
		Billing:            billingAddressData,
		Shipping:           shippingAddressData,
		OrderID:            orderID,
		PaymentID:          payment.Token,
		GraphqlPaymentID:   payment.Token,
		CustomerIpAddress:  payment.CustomerIPAddress.String,
		CustomerID:         customerId,
		CustomerEmail:      email,
		ReuseSource:        storeSource,
		Data:               additionalData,
		GraphqlCustomerID:  userID,
		StorePaymentMethod: payment.StorePaymentMethod,
		PaymentMetadata:    payment.Metadata,
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
//
// `storePaymentMethod` default to model.StorePaymentMethod.NONE
func (a *ServicePayment) CreatePayment(
	transaction boil.ContextTransactor,
	gateway string,
	total *decimal.Decimal,
	currency model.Currency,
	email string,
	customerIpAddress string,
	paymentToken string,
	extraData map[string]string,
	checkOut *model.Checkout,
	orDer *model.Order,
	returnUrl string,
	externalReference string,
	storePaymentMethod model.StorePaymentMethod,
	metadata model_helper.StringMap, // can be nil

) (*model.Payment, *model_helper.PaymentError, *model_helper.AppError) {
	// must at least provide either checkout or order, both is best :))
	if checkOut == nil && orDer == nil {
		return nil, nil, model_helper.NewAppError("CreatePayment", model_helper.InvalidArgumentAppErrorID, map[string]any{"Fields": "order/checkout"}, "please provide both order and checkout", http.StatusBadRequest)
	}

	if extraData == nil {
		extraData = make(map[string]string)
	}
	if metadata == nil {
		metadata = make(model_helper.StringMap)
	}

	var billingAddressID string
	if checkOut != nil && checkOut.BillingAddressID != nil {
		billingAddressID = *checkOut.BillingAddressID
	} else if orDer != nil && orDer.BillingAddressID != nil {
		billingAddressID = *orDer.BillingAddressID
	}

	var billingAddress *model.Address

	if billingAddressID != "" {
		var appErr *model_helper.AppError
		billingAddress, appErr = a.srv.Account.AddressById(billingAddressID)
		if appErr != nil {
			return nil, nil, appErr // this error can be either system error/not found error
		}
	}

	if billingAddress == nil {
		return nil, model_helper.NewPaymentError("CreatePayment", "Order does not have a billing address.", model_helper.BILLING_ADDRESS_NOT_SET), nil
	}

	payment := &model.Payment{
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
		Gateway:            gateway,
		Total:              total,
		ReturnUrl:          &returnUrl,
		PspReference:       &externalReference,
		IsActive:           model_helper.GetPointerOfValue(true),
		CustomerIpAddress:  &customerIpAddress,
		ExtraData:          model.ModelToJson(extraData),
		Token:              paymentToken,
		StorePaymentMethod: storePaymentMethod,
		ModelMetadata: model.ModelMetadata{
			Metadata: model.StringMAP(metadata),
		},
	}
	if checkOut != nil {
		payment.CheckoutID = &checkOut.Token
	}
	if orDer != nil {
		payment.OrderID = &orDer.Id
	}

	payment, appErr := a.srv.PaymentService().UpsertPayment(transaction, payment)
	return payment, nil, appErr
}

func (a *ServicePayment) GetAlreadyProcessedTransaction(paymentID string, gatewayResponse *model_helper.GatewayResponse) (*model.PaymentTransaction, *model_helper.AppError) {
	// get all transactions that belong to given payment
	trans, appErr := a.TransactionsByOption(&model.PaymentTransactionFilterOpts{
		Conditions: squirrel.Eq{model.TransactionTableName + "." + model.TransactionColumnPaymentID: paymentID},
	})
	if appErr != nil {
		return nil, appErr
	}

	var processedTran *model.PaymentTransaction

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
func (a *ServicePayment) CreateTransaction(paymentID string, kind model.TransactionKind, paymentInformation *model_helper.PaymentData, actionRequired bool, gatewayResponse *model_helper.GatewayResponse, errorMsg string, isSuccess bool) (*model.PaymentTransaction, *model_helper.AppError) {
	// Default values for token, amount, currency are only used in cases where
	// response from gateway was invalid or an exception occured
	if gatewayResponse == nil {
		var transactionId string
		if paymentInformation.Token != nil {
			transactionId = *paymentInformation.Token
		}
		gatewayResponse = &model_helper.GatewayResponse{
			Kind:           kind,
			ActionRequired: false,
			IsSucess:       isSuccess,
			TransactionID:  transactionId,
			Amount:         paymentInformation.Amount,
			Currency:       paymentInformation.Currency,
			Error:          errorMsg,
			RawResponse:    make(model_types.JSONString),
		}
	}

	tran := &model.PaymentTransaction{
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

	return a.SaveTransaction(nil, tran)
}

func (a *ServicePayment) GetAlreadyProcessedTransactionOrCreateNewTransaction(paymentID string, kind model.TransactionKind, paymentInformation *model_helper.PaymentData, actionRequired bool, gatewayResponse *model_helper.GatewayResponse, errorMsg string) (*model.PaymentTransaction, *model_helper.AppError) {
	if gatewayResponse != nil && gatewayResponse.TransactionAlreadyProcessed {
		transaction, appErr := a.GetAlreadyProcessedTransaction(paymentID, gatewayResponse)
		if appErr != nil {
			if appErr.StatusCode == http.StatusInternalServerError {
				return nil, appErr
			}
			// ignore not found error
		}
		if transaction != nil {
			return transaction, nil
		}
	}

	return a.CreateTransaction(paymentID, kind, paymentInformation, actionRequired, gatewayResponse, errorMsg, false)
}

// CleanCapture Check if payment can be captured.
func (a *ServicePayment) CleanCapture(pm *model.Payment, amount decimal.Decimal) *model.PaymentError {
	if amount.LessThanOrEqual(decimal.Zero) {
		return model.NewPaymentError("CleanCapture", "Amount should be a positive number.", model.INVALID)
	}
	if !pm.CanCapture() {
		return model.NewPaymentError("CleanCapture", "This payment cannot be captured.", model.INVALID)
	}
	// amount > payment's total || amount > payment's Total - payment's CapturedAmount
	if amount.GreaterThan(*pm.Total) || amount.GreaterThan((*pm.Total).Sub(*pm.CapturedAmount)) {
		return model.NewPaymentError("CleanCapture", "Unable to charge more than un-captured amount.", model.INVALID)
	}

	return nil
}

// CleanAuthorize Check if payment can be authorized
func (a *ServicePayment) CleanAuthorize(payment *model.Payment) *model.PaymentError {
	if !payment.CanAuthorize() {
		return model.NewPaymentError("CleanAuthorize", "Charged transactions cannot be authorized again.", model.INVALID)
	}
	return nil
}

// ValidateGatewayResponse Validate response to be a correct format for Saleor to process.
func (a *ServicePayment) ValidateGatewayResponse(response *model_helper.GatewayResponse) *model.GatewayError {
	if response == nil {
		return &model.GatewayError{
			Where:   "ValidateGatewayResponse",
			Message: "Gateway needs to return a GatewayResponse obj",
		}
	}

	// checks if response's Kind is valid transaction kind:
	if _, ok := model.TransactionKindString[response.Kind]; !ok {
		validTransactionKinds := []string{}
		for key := range model.TransactionKindString {
			validTransactionKinds = append(validTransactionKinds, key.String())
		}

		return &model.GatewayError{
			Where:   "ValidateGatewayResponse",
			Message: "Gateway response kind must be one of " + strings.Join(validTransactionKinds, ", "),
		}
	}

	// checks if response's RawResponse is json encodable
	_, err := json.Marshal(response.RawResponse)
	if err != nil {
		return &model.GatewayError{
			Where:   "ValidateGatewayResponse",
			Message: "Gateway response needs to be json serializable",
		}
	}

	return nil
}

// GatewayPostProcess
func (a *ServicePayment) GatewayPostProcess(paymentTransaction model.PaymentTransaction, payment *model.Payment) *model_helper.AppError {
	// create transaction
	transaction := a.srv.Store.GetMaster().Begin()
	if transaction.Error != nil {
		return model_helper.NewAppError("GatewayPostProcess", model_helper.ErrorCreatingTransactionErrorID, nil, transaction.Error.Error(), http.StatusInternalServerError)
	}
	defer a.srv.Store.FinalizeTransaction(transaction)

	var (
		changedFields util.AnyArray[string]
		appErr        *model_helper.AppError
	)

	if !paymentTransaction.IsSuccess || paymentTransaction.AlreadyProcessed {
		if len(changedFields) > 0 {
			if _, appErr = a.UpsertPayment(nil, payment); appErr != nil {
				return appErr
			}
		}
		return nil
	}

	if paymentTransaction.ActionRequired {
		payment.ToConfirm = true
		// changedFields = append(changedFields, "to_confirm")
		if _, appErr = a.UpsertPayment(transaction, payment); appErr != nil {
			return appErr
		}

		return nil
	}

	// to_confirm is defined by the paymentTransaction.action_required. Payment doesn't
	// require confirmation when we got action_required == False
	if payment.ToConfirm {
		payment.ToConfirm = false
		changedFields = append(changedFields, "to_confirm")
	}

	switch paymentTransaction.Kind {
	case model.TRANSACTION_KIND_CAPTURE, model.TRANSACTION_KIND_REFUND_REVERSED:
		payment.CapturedAmount = model_helper.GetPointerOfValue(payment.CapturedAmount.Add(*paymentTransaction.Amount))
		payment.IsActive = model_helper.GetPointerOfValue(true)
		// Set payment charge status to fully charged
		// only if there is no more amount needs to charge
		payment.ChargeStatus = model.PAYMENT_CHARGE_STATUS_PARTIALLY_CHARGED
		if payment.GetChargeAmount().LessThanOrEqual(decimal.Zero) {
			payment.ChargeStatus = model.PAYMENT_CHARGE_STATUS_FULLY_CHARGED
		}
		changedFields = append(changedFields, "charge_status", "captured_amount", "update_at")

	case model.TRANSACTION_KIND_VOID:
		payment.IsActive = model_helper.GetPointerOfValue(false)
		changedFields = append(changedFields, "is_active", "update_at")

	case model.TRANSACTION_KIND_REFUND:
		changedFields = append(changedFields, "captured_amount", "update_at")
		payment.CapturedAmount = model_helper.GetPointerOfValue(payment.CapturedAmount.Sub(*paymentTransaction.Amount))
		payment.ChargeStatus = model.PAYMENT_CHARGE_STATUS_PARTIALLY_REFUNDED
		if payment.CapturedAmount.LessThanOrEqual(decimal.Zero) {
			payment.CapturedAmount = model_helper.GetPointerOfValue(decimal.Zero)
			payment.ChargeStatus = model.PAYMENT_CHARGE_STATUS_FULLY_REFUNDED
			payment.IsActive = model_helper.GetPointerOfValue(false)
		}

	case model.TRANSACTION_KIND_PENDING:
		payment.ChargeStatus = model.PAYMENT_CHARGE_STATUS_PENDING
		changedFields = append(changedFields, "charge_status")

	case model.TRANSACTION_KIND_CANCEL:
		payment.ChargeStatus = model.PAYMENT_CHARGE_STATUS_CANCELLED
		payment.IsActive = model_helper.GetPointerOfValue(false)
		changedFields = append(changedFields, "charge_status", "is_active")

	case model.TRANSACTION_KIND_CAPTURE_FAILED:
		if payment.ChargeStatus == model.PAYMENT_CHARGE_STATUS_PARTIALLY_CHARGED || payment.ChargeStatus == model.PAYMENT_CHARGE_STATUS_FULLY_CHARGED {
			payment.CapturedAmount = model_helper.GetPointerOfValue(payment.CapturedAmount.Sub(*paymentTransaction.Amount))
			payment.ChargeStatus = model.PAYMENT_CHARGE_STATUS_PARTIALLY_CHARGED
			if payment.CapturedAmount.LessThanOrEqual(decimal.Zero) {
				payment.CapturedAmount = model_helper.GetPointerOfValue(decimal.Zero)
			}
			changedFields = append(changedFields, "charge_status", "captured_amount", "update_at")
		}
	}

	if len(changedFields) > 0 {
		if _, appErr := a.UpsertPayment(transaction, payment); appErr != nil {
			return appErr
		}
	}

	paymentTransaction.AlreadyProcessed = true
	if _, appErr = a.UpdateTransaction(&paymentTransaction); appErr != nil {
		return appErr
	}

	if changedFields.Contains("captured_amount") && payment.OrderID != nil {
		orDer, appErr := a.srv.OrderService().OrderById(*payment.OrderID)
		if appErr != nil {
			return appErr
		}
		if appErr = a.srv.OrderService().UpdateOrderTotalPaid(transaction, orDer); appErr != nil {
			return appErr
		}
	}

	// commit transaction
	if err := transaction.Commit().Error; err != nil {
		return model_helper.NewAppError("GatewayPostProcess", model_helper.ErrorCommittingTransactionErrorID, nil, err.Error(), http.StatusInternalServerError)
	}

	return nil
}

// FetchCustomerId Retrieve users customer_id stored for desired gateway.
// returning string could be "" or long string
func (a *ServicePayment) FetchCustomerId(user *model.User, gateway string) (string, *model_helper.AppError) {
	metaKey := prepareKeyForGatewayCustomerId(gateway)
	return user.PrivateMetadata.Get(metaKey, ""), nil
}

// StoreCustomerId stores new value into given user's PrivateMetadata
func (a *ServicePayment) StoreCustomerId(userID string, gateway string, customerID string) *model_helper.AppError {
	metaKey := prepareKeyForGatewayCustomerId(gateway)
	user, appErr := a.srv.AccountService().UserById(context.Background(), userID)
	if appErr != nil {
		return appErr
	}
	user.PrivateMetadata.Set(metaKey, customerID)
	_, appErr = a.srv.AccountService().UpdateUser(user, false)
	return appErr
}

// prepareKeyForGatewayCustomerId just trims spaces, upper then concatenates ".customer_id" to given `gatewayName`.
//
//	strings.TrimSpace(strings.ToUpper(gatewayName)) + ".customer_id"
func prepareKeyForGatewayCustomerId(gatewayName string) string {
	return strings.TrimSpace(strings.ToUpper(gatewayName)) + ".customer_id"
}

// UpdatePayment
func (a *ServicePayment) UpdatePayment(payment model.Payment, gatewayResponse *model_helper.GatewayResponse) *model_helper.AppError {
	var firstChange, secondChange bool

	if gatewayResponse.PspReference != "" {
		payment.PspReference = &gatewayResponse.PspReference
		firstChange = true
	}

	if gatewayResponse.PaymentMethodInfo != nil {
		secondChange = a.UpdatePaymentMethodDetails(payment, gatewayResponse.PaymentMethodInfo)
	}

	if firstChange || secondChange {
		_, appErr := a.UpsertPayment(nil, &payment)
		if appErr != nil {
			return appErr
		}
	}

	return nil
}

func (a *ServicePayment) UpdatePaymentMethodDetails(payment model.Payment, paymentMethodInfo *model.PaymentMethodInfo) (changed bool) {
	changed = true

	if paymentMethodInfo == nil {
		changed = false
		return
	}

	if brand := paymentMethodInfo.Brand; brand != nil {
		payment.CcBrand = *brand
	}
	if last4 := paymentMethodInfo.Last4; last4 != nil {
		payment.CcLastDigits = *last4
	}
	if exprYear := paymentMethodInfo.ExpYear; exprYear != nil {
		payment.CcExpYear = exprYear
	}
	if exprMonth := paymentMethodInfo.ExpMonth; exprMonth != nil {
		payment.CcExpMonth = exprMonth
	}
	if paymentType := paymentMethodInfo.Type; paymentType != nil {
		payment.PaymentMethodType = *paymentType
	}

	return
}

func (a *ServicePayment) GetPaymentToken(payment *model.Payment) (string, *model.PaymentError, *model_helper.AppError) {
	authTransactions, appErr := a.TransactionsByOption(&model.PaymentTransactionFilterOpts{
		Conditions: squirrel.Eq{
			model.TransactionTableName + ".Kind":      model.TRANSACTION_KIND_AUTH,
			model.TransactionTableName + ".IsSuccess": true,
		},
	})
	if appErr != nil {
		if appErr.StatusCode == http.StatusInternalServerError {
			return "", nil, appErr
		}
		return "", model.NewPaymentError("GetPaymentToken", "Cannot process unauthorized transaction", model.INVALID), appErr
	}

	return authTransactions[0].Token, nil, nil
}

// IsCurrencySupported Return true if the given gateway supports given currency.
func (a *ServicePayment) IsCurrencySupported(currency string, gatewayID string, manager interfaces.PluginManagerInterface) bool {
	for _, gateway := range manager.ListPaymentGateways(currency, nil, "", true) {
		if gateway.Id == gatewayID {
			return true
		}
	}

	return false
}

// Convert minor unit (smallest unit of currency) to decimal value.
//
// (value: 1000, currency: USD) will be converted to 10.00
func PriceFromMinorUnit(value string, currency string) (*decimal.Decimal, error) {
	deci, err := decimal.NewFromString(value)
	if err != nil {
		return nil, err
	}

	precision, err := goprices.GetCurrencyPrecision(currency)
	if err != nil {
		return nil, err
	}

	numberPlaces := decimal.NewFromInt(10).Pow(decimal.NewFromInt32(int32(-precision)))
	mul := deci.Mul(numberPlaces)
	return &mul, nil
}

// Convert decimal value to the smallest unit of currency.
//
// Take the value, discover the precision of currency and multiply value by
// Decimal('10.0'), then change quantization to remove the comma.
// Decimal(10.0) -> str(1000)
// func PriceToMinorUnit(value decimal.Decimal, currency string) (string, error) {
// 	precision, err := goprices.GetCurrencyPrecision(currency)
// 	if err != nil {
// 		return "", err
// 	}
// 	value = value.RoundUp(int32(precision))
// }

// PaymentOwnedByUser checks if given user is authenticated and owns given payment
//
// NOTE: if the `user` is unauthenticated, don't call me, just returns false
// func (s *ServicePayment) PaymentOwnedByUser(paymentID string, userID string) (bool, *model_helper.AppError) {
// 	s.srv.Store.Payment().FilterByOption(&model.PaymentFilterOption{})
// }
