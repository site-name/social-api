/*
NOTE: This package is initialized during server startup (modules/imports does that)
so the init() function get the chance to register a function to create `ServiceAccount`
*/
package payment

import (
	"net/http"

	goprices "github.com/site-name/go-prices"
	"github.com/sitename/sitename/app"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model_helper"
	"github.com/sitename/sitename/modules/model_types"
	"github.com/sitename/sitename/modules/util"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

type ServicePayment struct {
	srv *app.Server
}

func init() {
	app.RegisterService(func(s *app.Server) error {
		s.Payment = &ServicePayment{s}
		return nil
	})
}

func (a *ServicePayment) PaymentByID(transaction boil.ContextTransactor, paymentID string, lockForUpdate bool) (*model.Payment, *model_helper.AppError) {
	payments, appErr := a.PaymentsByOption(model_helper.PaymentFilterOptions{
		CommonQueryOptions: model_helper.NewCommonQueryOptions(
			model.PaymentWhere.ID.EQ(paymentID),
		),
	})
	if appErr != nil {
		return nil, appErr
	}
	return payments[0], nil
}

func (a *ServicePayment) PaymentsByOption(option model_helper.PaymentFilterOptions) (model.PaymentSlice, *model_helper.AppError) {
	payments, err := a.srv.Store.Payment().FilterByOption(option)
	if err != nil {
		return nil, model_helper.NewAppError("PaymentsByOption", "app.payment.error_finding_payments_by_option.app_error", nil, err.Error(), http.StatusInternalServerError)
	}

	return payments, nil
}

func (a *ServicePayment) GetLastOrderPayment(orderID string) (*model.Payment, *model_helper.AppError) {
	payments, appError := a.PaymentsByOption(model_helper.PaymentFilterOptions{
		CommonQueryOptions: model_helper.NewCommonQueryOptions(
			model.PaymentWhere.OrderID.EQ(model_types.NewNullString(orderID)),
			qm.Limit(1),
			qm.OrderBy(model.PaymentColumns.CreatedAt+" "+model_helper.DESC.String()),
		),
	})
	if appError != nil {
		return nil, appError
	}
	if len(payments) == 0 {
		return nil, model_helper.NewAppError("GetLastOrderPayment", "app.payment.order_has_no_payment.app_error", nil, "order has no payment yet", http.StatusNotFound)
	}

	return payments[0], nil
}

func (a *ServicePayment) PaymentIsAuthorized(paymentID string) (bool, *model_helper.AppError) {
	trans, appErr := a.TransactionsByOption(model_helper.PaymentTransactionFilterOpts{
		CommonQueryOptions: model_helper.NewCommonQueryOptions(
			model.PaymentTransactionWhere.PaymentID.EQ(paymentID),
		),
	})
	if appErr != nil {
		return false, appErr
	}

	for _, tran := range trans {
		if tran.Kind == model.TransactionKindAuth && tran.IsSuccess && !tran.ActionRequired {
			return true, nil
		}
	}

	return false, nil
}

func (a *ServicePayment) PaymentGetAuthorizedAmount(payment model.Payment) (*goprices.Money, *model_helper.AppError) {
	authorizedMoney, err := util.ZeroMoney(payment.Currency.String())
	if err != nil {
		return nil, model_helper.NewAppError("PaymentGetAuthorizedAmount", "app.payment.create_zero_money.app_error", nil, err.Error(), http.StatusInternalServerError)
	}

	transactions, appErr := a.TransactionsByOption(model_helper.PaymentTransactionFilterOpts{
		CommonQueryOptions: model_helper.NewCommonQueryOptions(
			model.PaymentTransactionWhere.PaymentID.EQ(payment.ID),
		),
	})
	if appErr != nil {
		return nil, appErr
	}

	// There is no authorized amount anymore when capture is succeeded
	// since capture can only be made once, even it is a partial capture
	for _, tran := range transactions {
		if tran.Kind == model.TransactionKindCapture && tran.IsSuccess {
			return authorizedMoney, nil
		}
	}

	// Filter the succeeded auth transactions
	for _, transaction := range transactions {
		if transaction.Kind == model.TransactionKindAuth && transaction.IsSuccess && !transaction.ActionRequired {
			authorizedMoney, err = authorizedMoney.Add(goprices.Money{
				Amount:   transaction.Amount,
				Currency: transaction.Currency.String(),
			})
			if err != nil {
				return nil, model_helper.NewAppError("PaymentGetAuthorizedAmount", "app.payment.error_calculation_payment_authorized_amount.app_error", nil, err.Error(), http.StatusInternalServerError)
			}
		}
	}

	return authorizedMoney, nil
}

func (a *ServicePayment) PaymentCanVoid(payment model.Payment) (bool, *model_helper.AppError) {
	authorized, err := a.PaymentIsAuthorized(payment.ID)
	if err != nil {
		return false, err
	}

	return payment.IsActive && model_helper.PaymentIsNotCharged(payment) && authorized, nil
}

func (a *ServicePayment) UpsertPayment(transaction boil.ContextTransactor, payment model.Payment) (*model.Payment, *model_helper.AppError) {
	savedPayment, err := a.srv.Store.Payment().Upsert(transaction, payment)
	if err != nil {
		return nil, model_helper.NewAppError("UpsertPayment", "app.payment.error_upserting_payment.app_error", nil, err.Error(), http.StatusInternalServerError)
	}

	return savedPayment, nil
}

func (a *ServicePayment) GetAllPaymentsByCheckout(checkoutToken string) (model.PaymentSlice, *model_helper.AppError) {
	payments, appErr := a.PaymentsByOption(model_helper.PaymentFilterOptions{
		CommonQueryOptions: model_helper.NewCommonQueryOptions(
			model.PaymentWhere.CheckoutID.EQ(model_types.NewNullString(checkoutToken)),
		),
	})
	if appErr != nil {
		return nil, appErr
	}
	return payments, nil
}

func (s *ServicePayment) UpdatePaymentsOfCheckout(transaction boil.ContextTransactor, checkoutToken string, option model_helper.PaymentPatch) *model_helper.AppError {
	err := s.srv.Store.Payment().UpdatePaymentsOfCheckout(transaction, checkoutToken, option)
	if err != nil {
		return model_helper.NewAppError("UpdatePaymentsOfCheckout", "app.payment.error_updating_payments_of_checkout.app_error", nil, err.Error(), http.StatusInternalServerError)
	}

	return nil
}
