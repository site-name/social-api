package payment

import (
	"context"
	"net/http"

	"github.com/shopspring/decimal"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model/account"
	"github.com/sitename/sitename/model/checkout"
	"github.com/sitename/sitename/model/order"
	"github.com/sitename/sitename/model/payment"
)

func (a *AppPayment) CreatePaymentInformation(pm *payment.Payment, paymentToken *string, amount *decimal.Decimal, customerId *string, storeSource bool, additionalData map[string]string) (*payment.PaymentData, *model.AppError) {
	var (
		billingAddress  *account.Address
		shippingAddress *account.Address
		amount_         *decimal.Decimal = pm.Total

		billingAddressID  string
		shippingAddressID string
		email             string = pm.BillingEmail
		orderId           string
		customerIpAddress string
		appErr            *model.AppError
	)

	if amount != nil {
		amount_ = amount
	}

	// checks if pm has checkout
	if pm.CheckoutID != nil && model.IsValidId(*pm.CheckoutID) {
		checkout, appErr := a.CheckoutApp().CheckoutbyToken(*pm.CheckoutID)
		if appErr != nil {
			return nil, appErr
		}

		// get checkout user
		if checkout.UserID != nil {
			user, appErr := a.AccountApp().UserById(context.Background(), *checkout.UserID)
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
	} else if pm.OrderID != nil && model.IsValidId(*pm.OrderID) { // checks if pm has order
		order, appErr := a.OrderApp().OrderById(*pm.OrderID)
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
		billingAddress, appErr = a.AccountApp().AddressById(billingAddressID)
		if appErr != nil {
			return nil, appErr
		}

		shippingAddress, appErr = a.AccountApp().AddressById(shippingAddressID)
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

	if pm.CustomerIpAddress != nil {
		customerIpAddress = *pm.CustomerIpAddress
	}

	return &payment.PaymentData{
		Gateway:           pm.GateWay,
		Amount:            *amount_,
		Currency:          pm.Currency,
		Billing:           billingAddressData,
		Shipping:          shippingAddressData,
		PaymentID:         pm.Id,
		GraphqlPaymentID:  pm.Id,
		OrderID:           orderId,
		CustomerIpAddress: customerIpAddress,
		CustomerEmail:     email,
		Token:             paymentToken,
		CustomerID:        customerId,
		ReuseSource:       storeSource,
		Data:              additionalData,
	}, nil
}

func (a *AppPayment) GetAlreadyProcessedTransaction(paymentID string, gatewayResponse *payment.GatewayResponse) (*payment.PaymentTransaction, *model.AppError) {
	// get all transactions that belong to given payment

	trans, appErr := a.PaymentApp().GetAllPaymentTransactions(paymentID)
	if appErr != nil {
		return nil, appErr
	}

	var processedTran *payment.PaymentTransaction

	// find the most recent transaction
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

func (a *AppPayment) CreatePayment(gateway, currency, email, customerIpAddress, paymentToken, returnUrl, externalReference string, total decimal.Decimal, extraData map[string]string, checkOut *checkout.Checkout, orDer *order.Order) (*payment.Payment, *model.AppError) {
	// must at least provider either checkout or order, both is best :))
	if checkOut == nil && orDer == nil {
		return nil, model.NewAppError("CreatePayment", "app.payment.checkout_order_required.app_error", nil, "", http.StatusBadRequest)
	}

	if extraData == nil {
		extraData = make(map[string]string)
	}

	var (
		billingAddress   *account.Address
		billingAddressID string
	)

	if checkOut != nil && checkOut.BillingAddressID != nil {
		billingAddressID = *checkOut.BillingAddressID
	} else if orDer != nil && orDer.BillingAddressID != nil {
		billingAddressID = *orDer.BillingAddressID
	}

	if billingAddressID == "" || !model.IsValidId(billingAddressID) {
		return nil, model.NewAppError("CreatePayment", "app.payment.order_billing_address_not_set.app_error", nil, "", http.StatusBadRequest)
	}

	billingAddress, appErr := a.AccountApp().AddressById(billingAddressID)
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
		Total:              &total,
		ReturnUrl:          &returnUrl,
		PspReference:       &externalReference,
		IsActive:           true,
		CustomerIpAddress:  &customerIpAddress,
		ExtraData:          model.ModelToJson(extraData),
		Token:              paymentToken,
	}
	if checkOut != nil {
		payment.CheckoutID = &checkOut.Token
	}
	if orDer != nil {
		payment.OrderID = &orDer.Id
	}

	return a.PaymentApp().SavePayment(payment)
}

func (a *AppPayment) CreatePaymentTransaction(paymentID string, kind string, paymentInformation *payment.PaymentData, actionRequired bool, gatewayResponse *payment.GatewayResponse, errorMsg string, isSuccess bool) (*payment.PaymentTransaction, *model.AppError) {
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

	return a.PaymentApp().SaveTransaction(tran)
}

func (a *AppPayment) GetAlreadyProcessedTransactionOrCreateNewTransaction(paymentID, kind string, paymentInformation *payment.PaymentData, actionRequired bool, gatewayResponse *payment.GatewayResponse, errorMsg string) (*payment.PaymentTransaction, *model.AppError) {
	if gatewayResponse != nil && gatewayResponse.TransactionAlreadyProcessed {
		return a.GetAlreadyProcessedTransaction(paymentID, gatewayResponse)
	}

	return a.CreatePaymentTransaction(paymentID, kind, paymentInformation, actionRequired, gatewayResponse, errorMsg, false)
}

func (a *AppPayment) CleanCapture(pm *payment.Payment, amount decimal.Decimal) *model.AppError {
	// where := "CleanCapture"
	// if amount.LessThanOrEqual(decimal.Zero) {
	// 	return model.NewAppError(where, "", )
	// }
	// TODO: fixme
	panic("not implemented")
}
