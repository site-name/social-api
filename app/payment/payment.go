package payment

import (
	"net/http"

	goprices "github.com/site-name/go-prices"
	"github.com/sitename/sitename/app"
	"github.com/sitename/sitename/app/sub_app_iface"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model/payment"
	"github.com/sitename/sitename/modules/util"
	"github.com/sitename/sitename/store"
)

// AppPayment handle all logics related to payment
type AppPayment struct {
	app app.AppIface
}

func init() {
	app.RegisterPaymentApp(func(a app.AppIface) sub_app_iface.PaymentApp {
		return &AppPayment{a}
	})
}

// PaymentsByOption returns all payments that satisfy given option
func (a *AppPayment) PaymentsByOption(option *payment.PaymentFilterOption) ([]*payment.Payment, *model.AppError) {
	payments, err := a.app.Srv().Store.Payment().FilterByOption(option)
	if err != nil {
		return nil, store.AppErrorFromDatabaseLookupError("PaymentsByOption", "app.payment.error_finding_payments_by_option.app_error", err)
	}

	return payments, nil
}

func (a *AppPayment) GetLastOrderPayment(orderID string) (*payment.Payment, *model.AppError) {
	payments, appError := a.PaymentsByOption(&payment.PaymentFilterOption{
		OrderID: orderID,
	})
	if appError != nil {
		appError.Where = "GetLastOrderPayment"
		return nil, appError
	}

	var latestPayment *payment.Payment
	for _, payment := range payments {
		if latestPayment == nil && payment.CreateAt >= latestPayment.CreateAt {
			latestPayment = payment
		}
	}

	return latestPayment, nil
}

func (a *AppPayment) PaymentIsAuthorized(paymentID string) (bool, *model.AppError) {
	trans, err := a.GetAllPaymentTransactions(paymentID)
	if err != nil {
		return false, err
	}

	for _, tran := range trans {
		if tran.Kind == payment.AUTH && tran.IsSuccess && !tran.ActionRequired {
			return true, nil
		}
	}

	return false, nil
}

func (a *AppPayment) PaymentGetAuthorizedAmount(pm *payment.Payment) (*goprices.Money, *model.AppError) {
	authorizedMoney, err := util.ZeroMoney(pm.Currency)
	if err != nil {
		return nil, model.NewAppError("PaymentGetAuthorizedAmount", "app.payment.create_zero_money.app_error", nil, err.Error(), http.StatusInternalServerError)
	}

	trans, appErr := a.GetAllPaymentTransactions(pm.Id)
	if appErr != nil {
		return nil, appErr
	}

	// There is no authorized amount anymore when capture is succeeded
	// since capture can only be made once, even it is a partial capture
	for _, tran := range trans {
		if tran.Kind == payment.CAPTURE && tran.IsSuccess {
			return authorizedMoney, nil
		}
	}

	// Filter the succeeded auth transactions
	for _, tran := range trans {
		if tran.Kind == payment.AUTH && tran.IsSuccess && !tran.ActionRequired {
			authorizedMoney, err = authorizedMoney.Add(&goprices.Money{
				Amount:   tran.Amount,
				Currency: tran.Currency,
			})
			if err != nil {
				return nil, model.NewAppError("PaymentGetAuthorizedAmount", "app.payment.error_calculation_payment_authorized_amount.app_error", nil, err.Error(), http.StatusInternalServerError)
			}
		}
	}

	return authorizedMoney, nil
}

func (a *AppPayment) PaymentCanVoid(payMent *payment.Payment) (bool, *model.AppError) {
	authorized, err := a.PaymentIsAuthorized(payMent.Id)
	if err != nil {
		return false, err
	}

	return *payMent.IsActive && payMent.IsNotCharged() && authorized, nil
}

func (a *AppPayment) CreateOrUpdatePayment(pm *payment.Payment) (*payment.Payment, *model.AppError) {
	var (
		returnedPayment *payment.Payment
		appErr          *model.AppError
		err             error
	)

	if pm.Id == "" { // id not set mean creating new payment
		returnedPayment, err = a.app.Srv().Store.Payment().Save(pm)
		if err != nil {
			if apErr, ok := err.(*model.AppError); ok {
				return nil, apErr
			}
			appErr = model.NewAppError("CreateOrUpdatePayment", "app.payment.save_payment_error.app_error", nil, "", http.StatusInternalServerError)
		}
	} else { // otherwise update
		returnedPayment, err = a.app.Srv().Store.Payment().Update(pm)
		if err != nil {
			if apErr, ok := err.(*model.AppError); ok {
				return nil, apErr
			}
			appErr = model.NewAppError("CreateOrUpdatePayment", "app.payment.update_payment_error.app_error", nil, "", http.StatusInternalServerError)
		}
	}

	return returnedPayment, appErr
}

func (a *AppPayment) GetPaymentToken(paymentID string) (string, *model.AppError) {
	trans, appErr := a.GetAllPaymentTransactions(paymentID)
	if appErr != nil {
		return "", appErr
	}

	var tran *payment.PaymentTransaction

	// find most recent transaction that has kind = "auth" and was made successfully
	for _, tr := range trans {
		if tr.Kind == payment.AUTH && tr.IsSuccess {
			if tran == nil || tran.CreateAt <= tr.CreateAt {
				tran = tr
			}
		}
	}

	if tran == nil {
		return "", model.NewAppError("GetPaymentToken", "app.payment.no_authorized_payment_transaction.app_error", nil, "", http.StatusNotFound)
	}

	return tran.Token, nil
}

func (a *AppPayment) GetAllPaymentsByCheckout(checkoutToken string) ([]*payment.Payment, *model.AppError) {
	payments, appErr := a.PaymentsByOption(&payment.PaymentFilterOption{
		CheckoutToken: checkoutToken,
	})
	if appErr != nil {
		appErr.Where = "GetAllPaymentsByCheckout"
		return nil, appErr
	}
	return payments, nil
}
