/*
NOTE: This package is initialized during server startup (modules/imports does that)
so the init() function get the chance to register a function to create `ServiceAccount`
*/
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
	"github.com/sitename/sitename/store/store_iface"
)

// ServicePayment handle all logics related to payment
type ServicePayment struct {
	srv *app.Server
}

func init() {
	app.RegisterPaymentService(func(s *app.Server) (sub_app_iface.PaymentService, error) {
		return &ServicePayment{
			srv: s,
		}, nil
	})
}

// PaymentByID returns a payment with given id
func (a *ServicePayment) PaymentByID(transaction store_iface.SqlxTxExecutor, paymentID string, lockForUpdate bool) (*payment.Payment, *model.AppError) {
	payMent, err := a.srv.Store.Payment().Get(transaction, paymentID, lockForUpdate)
	if err != nil {
		return nil, store.AppErrorFromDatabaseLookupError("PaymentByID", "app.payment.error_finding_payment_by_id.app_error", err)
	}
	return payMent, nil
}

// PaymentsByOption returns all payments that satisfy given option
func (a *ServicePayment) PaymentsByOption(option *payment.PaymentFilterOption) ([]*payment.Payment, *model.AppError) {
	payments, err := a.srv.Store.Payment().FilterByOption(option)
	var (
		statusCode int
		errMessage string
	)
	if err != nil {
		statusCode = http.StatusInternalServerError
		errMessage = err.Error()
	} else if len(payments) == 0 {
		statusCode = http.StatusNotFound
	}

	if statusCode != 0 {
		return nil, model.NewAppError("PaymentsByOption", "app.payment.error_finding_payments_by_option.app_error", nil, errMessage, statusCode)
	}

	return payments, nil
}

func (a *ServicePayment) GetLastOrderPayment(orderID string) (*payment.Payment, *model.AppError) {
	payments, appError := a.PaymentsByOption(&payment.PaymentFilterOption{
		OrderID: orderID,
	})
	if appError != nil {
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

func (a *ServicePayment) PaymentIsAuthorized(paymentID string) (bool, *model.AppError) {
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

func (a *ServicePayment) PaymentGetAuthorizedAmount(pm *payment.Payment) (*goprices.Money, *model.AppError) {
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
				Amount:   *tran.Amount,
				Currency: tran.Currency,
			})
			if err != nil {
				return nil, model.NewAppError("PaymentGetAuthorizedAmount", "app.payment.error_calculation_payment_authorized_amount.app_error", nil, err.Error(), http.StatusInternalServerError)
			}
		}
	}

	return authorizedMoney, nil
}

// PaymentCanVoid checks if given payment is: Active && not charged and authorized
func (a *ServicePayment) PaymentCanVoid(payMent *payment.Payment) (bool, *model.AppError) {
	authorized, err := a.PaymentIsAuthorized(payMent.Id)
	if err != nil {
		return false, err
	}

	return *payMent.IsActive && payMent.IsNotCharged() && authorized, nil
}

// UpsertPayment updates or insert given payment, depends on the validity of its Id
func (a *ServicePayment) UpsertPayment(transaction store_iface.SqlxTxExecutor, payMent *payment.Payment) (*payment.Payment, *model.AppError) {
	var err error

	if !model.IsValidId(payMent.Id) {
		payMent, err = a.srv.Store.Payment().Save(transaction, payMent)
	} else {
		payMent, err = a.srv.Store.Payment().Update(transaction, payMent)
	}
	if err != nil {
		if appErr, ok := err.(*model.AppError); ok {
			return nil, appErr
		}
		var statusCode = http.StatusInternalServerError
		if _, ok := err.(*store.ErrNotFound); ok {
			statusCode = http.StatusNotFound
		}
		return nil, model.NewAppError("UpsertPayment", "app.payment.error_upserting_payment.app_error", nil, err.Error(), statusCode)
	}

	return payMent, nil
}

// GetAllPaymentsByCheckout returns all payments that belong to given checkout
func (a *ServicePayment) GetAllPaymentsByCheckout(checkoutToken string) ([]*payment.Payment, *model.AppError) {
	payments, appErr := a.PaymentsByOption(&payment.PaymentFilterOption{
		CheckoutToken: checkoutToken,
	})
	if appErr != nil {
		return nil, appErr
	}
	return payments, nil
}

// UpdatePaymentsOfCheckout updates payments of given checkout, with parameters specified in option
func (s *ServicePayment) UpdatePaymentsOfCheckout(transaction store_iface.SqlxTxExecutor, checkoutToken string, option *payment.PaymentPatch) *model.AppError {
	err := s.srv.Store.Payment().UpdatePaymentsOfCheckout(transaction, checkoutToken, option)
	if err != nil {
		return model.NewAppError("UpdatePaymentsOfCheckout", "app.payment.error_updating_payments_of_checkout.app_error", nil, err.Error(), http.StatusInternalServerError)
	}

	return nil
}
