package payment

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/Masterminds/squirrel"
	"github.com/site-name/decimal"
	goprices "github.com/site-name/go-prices"
	"github.com/sitename/sitename/app/plugin/interfaces"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/modules/util"
	"gorm.io/gorm"
)

// CreatePaymentInformation Extract order information along with payment details.
//
// Returns information required to process payment and additional
// billing/shipping addresses for optional fraud-prevention mechanisms.
func (a *ServicePayment) CreatePaymentInformation(payMent *model.Payment, paymentToken *string, amount *decimal.Decimal, customerId *string, storeSource bool, additionalData map[string]interface{}) (*model.PaymentData, *model.AppError) {

	var (
		billingAddressID  string
		shippingAddressID string
		billingAddress    *model.Address
		shippingAddress   *model.Address
		email             string
		userID            *string
	)

	if payMent.CheckoutID != nil {
		checkoutOfPayment, appErr := a.srv.CheckoutService().CheckoutByOption(&model.CheckoutFilterOption{
			Conditions: squirrel.Eq{model.CheckoutTableName + ".Token": *payMent.CheckoutID},
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
		addresses, appErr := a.srv.AccountService().AddressesByOption(&model.AddressFilterOption{
			Conditions: squirrel.Eq{model.AddressTableName + ".Id": []string{billingAddressID, shippingAddressID}},
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
		billingAddressData  *model.AddressData
		shippingAddressData *model.AddressData
	)
	if billingAddress != nil {
		billingAddressData = model.AddressDataFromAddress(billingAddress)
	}
	if shippingAddress != nil {
		shippingAddressData = model.AddressDataFromAddress(shippingAddress)
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

	return &model.PaymentData{
		Gateway:            payMent.GateWay,
		Token:              paymentToken,
		Amount:             *amount,
		Currency:           payMent.Currency,
		Billing:            billingAddressData,
		Shipping:           shippingAddressData,
		OrderID:            orderID,
		PaymentID:          payMent.Token,
		GraphqlPaymentID:   payMent.Token,
		CustomerIpAddress:  payMent.CustomerIpAddress,
		CustomerID:         customerId,
		CustomerEmail:      email,
		ReuseSource:        storeSource,
		Data:               additionalData,
		GraphqlCustomerID:  userID,
		StorePaymentMethod: payMent.StorePaymentMethod,
		PaymentMetadata:    model.StringMap(payMent.Metadata),
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
	transaction *gorm.DB,
	gateway string,
	total *decimal.Decimal,
	currency string,
	email string,
	customerIpAddress string,
	paymentToken string,
	extraData map[string]string,
	checkOut *model.Checkout,
	orDer *model.Order,
	returnUrl string,
	externalReference string,
	storePaymentMethod model.StorePaymentMethod,
	metadata model.StringMap, // can be nil

) (*model.Payment, *model.PaymentError, *model.AppError) {
	// must at least provide either checkout or order, both is best :))
	if checkOut == nil && orDer == nil {
		return nil, nil, model.NewAppError("CreatePayment", model.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "order/checkout"}, "please provide both order and checkout", http.StatusBadRequest)
	}

	if extraData == nil {
		extraData = make(map[string]string)
	}
	if metadata == nil {
		metadata = make(model.StringMap)
	}

	var billingAddressID string
	if checkOut != nil && checkOut.BillingAddressID != nil {
		billingAddressID = *checkOut.BillingAddressID
	} else if orDer != nil && orDer.BillingAddressID != nil {
		billingAddressID = *orDer.BillingAddressID
	}

	var billingAddress *model.Address

	if billingAddressID != "" {
		var appErr *model.AppError
		billingAddress, appErr = a.srv.AccountService().AddressById(billingAddressID)
		if appErr != nil {
			return nil, nil, appErr // this error can be either system error/not found error
		}
	}

	if billingAddress == nil {
		return nil, model.NewPaymentError("CreatePayment", "Order does not have a billing address.", model.BILLING_ADDRESS_NOT_SET), nil
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
		GateWay:            gateway,
		Total:              total,
		ReturnUrl:          &returnUrl,
		PspReference:       &externalReference,
		IsActive:           model.GetPointerOfValue(true),
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

func (a *ServicePayment) GetAlreadyProcessedTransaction(paymentID string, gatewayResponse *model.GatewayResponse) (*model.PaymentTransaction, *model.AppError) {
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
func (a *ServicePayment) CreateTransaction(paymentID string, kind model.TransactionKind, paymentInformation *model.PaymentData, actionRequired bool, gatewayResponse *model.GatewayResponse, errorMsg string, isSuccess bool) (*model.PaymentTransaction, *model.AppError) {
	// Default values for token, amount, currency are only used in cases where
	// response from gateway was invalid or an exception occured
	if gatewayResponse == nil {
		var transactionId string
		if paymentInformation.Token != nil {
			transactionId = *paymentInformation.Token
		}
		gatewayResponse = &model.GatewayResponse{
			Kind:           kind,
			ActionRequired: false,
			IsSucess:       isSuccess,
			TransactionID:  transactionId,
			Amount:         paymentInformation.Amount,
			Currency:       paymentInformation.Currency,
			Error:          errorMsg,
			RawResponse:    make(model.StringInterface),
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

func (a *ServicePayment) GetAlreadyProcessedTransactionOrCreateNewTransaction(paymentID string, kind model.TransactionKind, paymentInformation *model.PaymentData, actionRequired bool, gatewayResponse *model.GatewayResponse, errorMsg string) (*model.PaymentTransaction, *model.AppError) {
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
func (a *ServicePayment) CleanAuthorize(payMent *model.Payment) *model.PaymentError {
	if !payMent.CanAuthorize() {
		return model.NewPaymentError("CleanAuthorize", "Charged transactions cannot be authorized again.", model.INVALID)
	}
	return nil
}

// ValidateGatewayResponse Validate response to be a correct format for Saleor to process.
func (a *ServicePayment) ValidateGatewayResponse(response *model.GatewayResponse) *model.GatewayError {
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
func (a *ServicePayment) GatewayPostProcess(paymentTransaction model.PaymentTransaction, payMent *model.Payment) *model.AppError {
	// create transaction
	transaction := a.srv.Store.GetMaster().Begin()
	if transaction.Error != nil {
		return model.NewAppError("GatewayPostProcess", model.ErrorCreatingTransactionErrorID, nil, transaction.Error.Error(), http.StatusInternalServerError)
	}
	defer a.srv.Store.FinalizeTransaction(transaction)

	var (
		changedFields util.AnyArray[string]
		appErr        *model.AppError
	)

	if !paymentTransaction.IsSuccess || paymentTransaction.AlreadyProcessed {
		if len(changedFields) > 0 {
			if _, appErr = a.UpsertPayment(nil, payMent); appErr != nil {
				return appErr
			}
		}
		return nil
	}

	if paymentTransaction.ActionRequired {
		payMent.ToConfirm = true
		// changedFields = append(changedFields, "to_confirm")
		if _, appErr = a.UpsertPayment(transaction, payMent); appErr != nil {
			return appErr
		}

		return nil
	}

	// to_confirm is defined by the paymentTransaction.action_required. Payment doesn't
	// require confirmation when we got action_required == False
	if payMent.ToConfirm {
		payMent.ToConfirm = false
		changedFields = append(changedFields, "to_confirm")
	}

	switch paymentTransaction.Kind {
	case model.TRANSACTION_KIND_CAPTURE, model.TRANSACTION_KIND_REFUND_REVERSED:
		payMent.CapturedAmount = model.GetPointerOfValue(payMent.CapturedAmount.Add(*paymentTransaction.Amount))
		payMent.IsActive = model.GetPointerOfValue(true)
		// Set payment charge status to fully charged
		// only if there is no more amount needs to charge
		payMent.ChargeStatus = model.PAYMENT_CHARGE_STATUS_PARTIALLY_CHARGED
		if payMent.GetChargeAmount().LessThanOrEqual(decimal.Zero) {
			payMent.ChargeStatus = model.PAYMENT_CHARGE_STATUS_FULLY_CHARGED
		}
		changedFields = append(changedFields, "charge_status", "captured_amount", "update_at")

	case model.TRANSACTION_KIND_VOID:
		payMent.IsActive = model.GetPointerOfValue(false)
		changedFields = append(changedFields, "is_active", "update_at")

	case model.TRANSACTION_KIND_REFUND:
		changedFields = append(changedFields, "captured_amount", "update_at")
		payMent.CapturedAmount = model.GetPointerOfValue(payMent.CapturedAmount.Sub(*paymentTransaction.Amount))
		payMent.ChargeStatus = model.PAYMENT_CHARGE_STATUS_PARTIALLY_REFUNDED
		if payMent.CapturedAmount.LessThanOrEqual(decimal.Zero) {
			payMent.CapturedAmount = model.GetPointerOfValue(decimal.Zero)
			payMent.ChargeStatus = model.PAYMENT_CHARGE_STATUS_FULLY_REFUNDED
			payMent.IsActive = model.GetPointerOfValue(false)
		}

	case model.TRANSACTION_KIND_PENDING:
		payMent.ChargeStatus = model.PAYMENT_CHARGE_STATUS_PENDING
		changedFields = append(changedFields, "charge_status")

	case model.TRANSACTION_KIND_CANCEL:
		payMent.ChargeStatus = model.PAYMENT_CHARGE_STATUS_CANCELLED
		payMent.IsActive = model.GetPointerOfValue(false)
		changedFields = append(changedFields, "charge_status", "is_active")

	case model.TRANSACTION_KIND_CAPTURE_FAILED:
		if payMent.ChargeStatus == model.PAYMENT_CHARGE_STATUS_PARTIALLY_CHARGED || payMent.ChargeStatus == model.PAYMENT_CHARGE_STATUS_FULLY_CHARGED {
			payMent.CapturedAmount = model.GetPointerOfValue(payMent.CapturedAmount.Sub(*paymentTransaction.Amount))
			payMent.ChargeStatus = model.PAYMENT_CHARGE_STATUS_PARTIALLY_CHARGED
			if payMent.CapturedAmount.LessThanOrEqual(decimal.Zero) {
				payMent.CapturedAmount = model.GetPointerOfValue(decimal.Zero)
			}
			changedFields = append(changedFields, "charge_status", "captured_amount", "update_at")
		}
	}

	if len(changedFields) > 0 {
		if _, appErr := a.UpsertPayment(transaction, payMent); appErr != nil {
			return appErr
		}
	}

	paymentTransaction.AlreadyProcessed = true
	if _, appErr = a.UpdateTransaction(&paymentTransaction); appErr != nil {
		return appErr
	}

	if changedFields.Contains("captured_amount") && payMent.OrderID != nil {
		orDer, appErr := a.srv.OrderService().OrderById(*payMent.OrderID)
		if appErr != nil {
			return appErr
		}
		if appErr = a.srv.OrderService().UpdateOrderTotalPaid(transaction, orDer); appErr != nil {
			return appErr
		}
	}

	// commit transaction
	if err := transaction.Commit().Error; err != nil {
		return model.NewAppError("GatewayPostProcess", model.ErrorCommittingTransactionErrorID, nil, err.Error(), http.StatusInternalServerError)
	}

	return nil
}

// FetchCustomerId Retrieve users customer_id stored for desired gateway.
// returning string could be "" or long string
func (a *ServicePayment) FetchCustomerId(user *model.User, gateway string) (string, *model.AppError) {
	metaKey := prepareKeyForGatewayCustomerId(gateway)
	return user.PrivateMetadata.Get(metaKey, ""), nil
}

// StoreCustomerId stores new value into given user's PrivateMetadata
func (a *ServicePayment) StoreCustomerId(userID string, gateway string, customerID string) *model.AppError {
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
func (a *ServicePayment) UpdatePayment(payMent model.Payment, gatewayResponse *model.GatewayResponse) *model.AppError {
	var firstChange, secondChange bool

	if gatewayResponse.PspReference != "" {
		payMent.PspReference = &gatewayResponse.PspReference
		firstChange = true
	}

	if gatewayResponse.PaymentMethodInfo != nil {
		secondChange = a.UpdatePaymentMethodDetails(payMent, gatewayResponse.PaymentMethodInfo)
	}

	if firstChange || secondChange {
		_, appErr := a.UpsertPayment(nil, &payMent)
		if appErr != nil {
			return appErr
		}
	}

	return nil
}

func (a *ServicePayment) UpdatePaymentMethodDetails(payMent model.Payment, paymentMethodInfo *model.PaymentMethodInfo) (changed bool) {
	changed = true

	if paymentMethodInfo == nil {
		changed = false
		return
	}

	if brand := paymentMethodInfo.Brand; brand != nil {
		payMent.CcBrand = *brand
	}
	if last4 := paymentMethodInfo.Last4; last4 != nil {
		payMent.CcLastDigits = *last4
	}
	if exprYear := paymentMethodInfo.ExpYear; exprYear != nil {
		payMent.CcExpYear = exprYear
	}
	if exprMonth := paymentMethodInfo.ExpMonth; exprMonth != nil {
		payMent.CcExpMonth = exprMonth
	}
	if paymentType := paymentMethodInfo.Type; paymentType != nil {
		payMent.PaymentMethodType = *paymentType
	}

	return
}

func (a *ServicePayment) GetPaymentToken(payMent *model.Payment) (string, *model.PaymentError, *model.AppError) {
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
// func (s *ServicePayment) PaymentOwnedByUser(paymentID string, userID string) (bool, *model.AppError) {
// 	s.srv.Store.Payment().FilterByOption(&model.PaymentFilterOption{})
// }
