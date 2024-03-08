/*
NOTE: This package is initialized during server startup (modules/imports does that)
so the init() function get the chance to register a function to create `ServiceAccount`
*/
package payment

import (
	"net/http"

	"github.com/mattermost/squirrel"
	goprices "github.com/site-name/go-prices"
	"github.com/sitename/sitename/app"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model_helper"
	"github.com/sitename/sitename/modules/util"
	"github.com/sitename/sitename/store"
	"github.com/volatiletech/sqlboiler/v4/boil"
)

// ServicePayment handle all logics related to payment
type ServicePayment struct {
	srv *app.Server
}

func init() {
	app.RegisterService(func(s *app.Server) error {
		s.Payment = &ServicePayment{s}
		return nil
	})
}

// PaymentByID returns a payment with given id
func (a *ServicePayment) PaymentByID(transaction boil.ContextTransactor, paymentID string, lockForUpdate bool) (*model.Payment, *model_helper.AppError) {
	_, payments, appErr := a.PaymentsByOption(&model.PaymentFilterOption{
		Conditions:    squirrel.Expr(model.PaymentTableName+".Id = ?", paymentID),
		DbTransaction: transaction,
		LockForUpdate: lockForUpdate,
	})
	if appErr != nil {
		return nil, appErr
	}
	return payments[0], nil
}

// PaymentsByOption returns all payments that satisfy given option
func (a *ServicePayment) PaymentsByOption(option *model.PaymentFilterOption) (int64, []*model.Payment, *model_helper.AppError) {
	totalCount, payments, err := a.srv.Store.Payment().FilterByOption(option)
	if err != nil {
		return 0, nil, model_helper.NewAppError("PaymentsByOption", "app.payment.error_finding_payments_by_option.app_error", nil, err.Error(), http.StatusInternalServerError)
	}

	return totalCount, payments, nil
}

func (a *ServicePayment) GetLastOrderPayment(orderID string) (*model.Payment, *model_helper.AppError) {
	_, payments, appError := a.PaymentsByOption(&model.PaymentFilterOption{
		Conditions: squirrel.Eq{model.PaymentTableName + ".OrderID": orderID},
	})
	if appError != nil {
		return nil, appError
	}
	if len(payments) == 0 {
		return nil, model_helper.NewAppError("GetLastOrderPayment", "app.payment.order_has_no_payment.app_error", nil, "order has no payment yet", http.StatusNotFound)
	}

	var latestPayment *model.Payment
	for _, payment := range payments {
		if latestPayment == nil && payment.CreateAt >= latestPayment.CreateAt {
			latestPayment = payment
		}
	}

	return latestPayment, nil
}

func (a *ServicePayment) PaymentIsAuthorized(paymentID string) (bool, *model_helper.AppError) {
	trans, appErr := a.TransactionsByOption(&model.PaymentTransactionFilterOpts{
		Conditions: squirrel.Eq{model.TransactionTableName + "." + model.TransactionColumnPaymentID: paymentID},
	})
	if appErr != nil {
		return false, appErr
	}

	for _, tran := range trans {
		if tran.Kind == model.TRANSACTION_KIND_AUTH && tran.IsSuccess && !tran.ActionRequired {
			return true, nil
		}
	}

	return false, nil
}

func (a *ServicePayment) PaymentGetAuthorizedAmount(payment *model.Payment) (*goprices.Money, *model_helper.AppError) {
	authorizedMoney, err := util.ZeroMoney(payment.Currency)
	if err != nil {
		return nil, model_helper.NewAppError("PaymentGetAuthorizedAmount", "app.payment.create_zero_money.app_error", nil, err.Error(), http.StatusInternalServerError)
	}

	trans, appErr := a.TransactionsByOption(&model.PaymentTransactionFilterOpts{
		Conditions: squirrel.Eq{model.TransactionTableName + "." + model.TransactionColumnPaymentID: payment.Id},
	})
	if appErr != nil {
		return nil, appErr
	}

	// There is no authorized amount anymore when capture is succeeded
	// since capture can only be made once, even it is a partial capture
	for _, tran := range trans {
		if tran.Kind == model.TRANSACTION_KIND_CAPTURE && tran.IsSuccess {
			return authorizedMoney, nil
		}
	}

	// Filter the succeeded auth transactions
	for _, tran := range trans {
		if tran.Kind == model.TRANSACTION_KIND_AUTH && tran.IsSuccess && !tran.ActionRequired {
			authorizedMoney, err = authorizedMoney.Add(&goprices.Money{
				Amount:   *tran.Amount,
				Currency: tran.Currency,
			})
			if err != nil {
				return nil, model_helper.NewAppError("PaymentGetAuthorizedAmount", "app.payment.error_calculation_payment_authorized_amount.app_error", nil, err.Error(), http.StatusInternalServerError)
			}
		}
	}

	return authorizedMoney, nil
}

// PaymentCanVoid checks if given payment is: Active && not charged and authorized
func (a *ServicePayment) PaymentCanVoid(payMent *model.Payment) (bool, *model_helper.AppError) {
	authorized, err := a.PaymentIsAuthorized(payMent.Id)
	if err != nil {
		return false, err
	}

	return *payMent.IsActive && payMent.NotCharged() && authorized, nil
}

// UpsertPayment updates or insert given payment, depends on the validity of its Id
func (a *ServicePayment) UpsertPayment(transaction boil.ContextTransactor, payMent *model.Payment) (*model.Payment, *model_helper.AppError) {
	var err error

	if !model_helper.IsValidId(payMent.Id) {
		payMent, err = a.srv.Store.Payment().Save(transaction, payMent)
	} else {
		payMent, err = a.srv.Store.Payment().Update(transaction, payMent)
	}
	if err != nil {
		if appErr, ok := err.(*model_helper.AppError); ok {
			return nil, appErr
		}
		var statusCode = http.StatusInternalServerError
		if _, ok := err.(*store.ErrNotFound); ok {
			statusCode = http.StatusNotFound
		}
		return nil, model_helper.NewAppError("UpsertPayment", "app.payment.error_upserting_payment.app_error", nil, err.Error(), statusCode)
	}

	return payMent, nil
}

// GetAllPaymentsByCheckout returns all payments that belong to given checkout
func (a *ServicePayment) GetAllPaymentsByCheckout(checkoutToken string) ([]*model.Payment, *model_helper.AppError) {
	_, payments, appErr := a.PaymentsByOption(&model.PaymentFilterOption{
		Conditions: squirrel.Eq{model.PaymentTableName + ".CheckoutID": checkoutToken},
	})
	if appErr != nil {
		return nil, appErr
	}
	return payments, nil
}

// UpdatePaymentsOfCheckout updates payments of given checkout, with parameters specified in option
func (s *ServicePayment) UpdatePaymentsOfCheckout(transaction boil.ContextTransactor, checkoutToken string, option *model.PaymentPatch) *model_helper.AppError {
	err := s.srv.Store.Payment().UpdatePaymentsOfCheckout(transaction, checkoutToken, option)
	if err != nil {
		return model_helper.NewAppError("UpdatePaymentsOfCheckout", "app.payment.error_updating_payments_of_checkout.app_error", nil, err.Error(), http.StatusInternalServerError)
	}

	return nil
}
