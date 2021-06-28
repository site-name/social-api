package payment

import (
	"context"
	"net/http"

	"github.com/shopspring/decimal"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model/account"
	"github.com/sitename/sitename/model/checkout"
	"github.com/sitename/sitename/model/order"
	modelPayment "github.com/sitename/sitename/model/payment"
	"github.com/sitename/sitename/store"
)

func (a *AppPayment) CreatePaymentInformation(
	payment *modelPayment.Payment, paymentToken *string,
	amount *decimal.Decimal, customerId *string,
	storeSource bool, additionalData map[string]string,
) (*modelPayment.PaymentData, *model.AppError) {

	var (
		billingAddress  *account.Address
		shippingAddress *account.Address
		amount_         *decimal.Decimal

		billingAddressID  string
		shippingAddressID string
		email             string = payment.BillingEmail
		err               error
		orderId           string
		customerIpAddress string
	)

	// checks if payment has checkout
	if payment.CheckoutID != nil && model.IsValidId(*payment.CheckoutID) {
		checkout, err := a.Srv().Store.Checkout().Get(*payment.CheckoutID)
		if err != nil {
			return nil, store.AppErrorFromDatabaseLookupError("CreatePaymentInformation", "app.payment.payment_missing.app_error", err)
		}

		// get checkout user
		if checkout.UserID != nil {
			user, err := a.Srv().Store.User().Get(context.Background(), *checkout.UserID)
			if err != nil {
				return nil, store.AppErrorFromDatabaseLookupError("CreatePaymentInformation", "app.account.user_missing.app_error", err)
			}
			email = user.Email
		} else {
			email = checkout.Email
		}

		if checkout.BillingAddressID != nil && checkout.ShippingAddressID != nil {
			billingAddressID = *checkout.BillingAddressID
			shippingAddressID = *checkout.ShippingAddressID
		}
	} else if payment.OrderID != nil && model.IsValidId(*payment.OrderID) { // checks if payment has order
		order, err := a.Srv().Store.Order().Get(*payment.OrderID)
		if err != nil {
			return nil, store.AppErrorFromDatabaseLookupError("CreatePaymentInformation", "app.order.order_missing.app_error", err)
		}

		email = order.UserEmail
		orderId = order.Id

		if order.BillingAddressID != nil && order.ShippingAddressID != nil {
			billingAddressID = *order.BillingAddressID
			shippingAddressID = *order.ShippingAddressID
		}
	}

	var (
		billingAddressData  *modelPayment.AddressData
		shippingAddressData *modelPayment.AddressData
	)

	if model.IsValidId(billingAddressID) && model.IsValidId(shippingAddressID) {
		billingAddress, err = a.Srv().Store.Address().Get(billingAddressID)
		if err != nil {
			return nil, store.AppErrorFromDatabaseLookupError("", "app.account.address_missing.app_error", err)
		}

		shippingAddress, err = a.Srv().Store.Address().Get(shippingAddressID)
		if err != nil {
			return nil, store.AppErrorFromDatabaseLookupError("", "app.account.address_missing.app_error", err)
		}
	}

	if billingAddress != nil {
		billingAddressData = modelPayment.AddressDataFromAddress(billingAddress)
	}
	if shippingAddress != nil {
		shippingAddressData = modelPayment.AddressDataFromAddress(shippingAddress)
	}

	if amount != nil {
		amount_ = amount
	} else {
		amount_ = payment.Total
	}
	if payment.CustomerIpAddress != nil {
		customerIpAddress = *payment.CustomerIpAddress
	}

	return &modelPayment.PaymentData{
		Gateway:           payment.GateWay,
		Amount:            amount_,
		Currency:          payment.Currency,
		Billing:           billingAddressData,
		Shipping:          shippingAddressData,
		PaymentID:         payment.Id,
		GraphqlPaymentID:  payment.Id,
		OrderID:           orderId,
		CustomerIpAddress: customerIpAddress,
		CustomerEmail:     email,
		Token:             paymentToken,
		CustomerID:        customerId,
		ReuseSource:       storeSource,
		Data:              additionalData,
	}, nil
}

func (a *AppPayment) GetAlreadyProcessedTransaction(payment *modelPayment.Payment, gatewayResponse *modelPayment.GatewayResponse) (*modelPayment.PaymentTransaction, *model.AppError) {
	// get all transactions that belong to given payment
	trans, err := a.Srv().Store.PaymentTransaction().GetAllByPaymentID(payment.Id)
	if err != nil {
		return nil, store.AppErrorFromDatabaseLookupError("", "app.payment.payment_transaction_missing.app_error", err)
	}

	var processedTran *modelPayment.PaymentTransaction

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

func (a *AppPayment) CreatePayment(gateway, currency, email, customerIpAddress, paymentToken, returnUrl, externalReference string, total decimal.Decimal, extraData map[string]string, checkOut *checkout.Checkout, orDer *order.Order) (*modelPayment.Payment, *model.AppError) {
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

	billingAddress, err := a.Srv().Store.Address().Get(billingAddressID)
	if err != nil {
		return nil, store.AppErrorFromDatabaseLookupError("CreatePayment", "app.account.address_missing.app_error", err)
	}

	payment := &modelPayment.Payment{
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
		payment.CheckoutID = &checkOut.Id
	}
	if orDer != nil {
		payment.OrderID = &orDer.Id
	}

	payment, err = a.Srv().Store.Payment().Save(payment)
	if err != nil {
		if appErr, ok := err.(*model.AppError); ok {
			return nil, appErr
		}
		return nil, model.NewAppError("CreatePayment", "app.payment.error_saving.app_error", nil, err.Error(), http.StatusInternalServerError)
	}

	return payment, nil
}
