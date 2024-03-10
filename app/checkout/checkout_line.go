package checkout

import (
	"net/http"

	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model_helper"
	"github.com/sitename/sitename/store"
	"github.com/volatiletech/sqlboiler/v4/boil"
)

func (a *ServiceCheckout) CheckoutLinesByCheckoutToken(checkoutToken string) (model.CheckoutLineSlice, *model_helper.AppError) {
	return a.CheckoutLinesByOption(model_helper.CheckoutLineFilterOptions{
		CommonQueryOptions: model_helper.NewCommonQueryOptions(
			model.CheckoutLineWhere.CheckoutID.EQ(checkoutToken),
		),
	})
}

func (a *ServiceCheckout) DeleteCheckoutLines(transaction boil.ContextTransactor, checkoutLineIDs []string) *model_helper.AppError {
	err := a.srv.Store.CheckoutLine().DeleteLines(transaction, checkoutLineIDs)
	if err != nil {
		return model_helper.NewAppError("DeleteCheckoutLines", "app.checkout.error_deleting_checkoutlines.app_error", nil, err.Error(), http.StatusInternalServerError)
	}

	return nil
}

func (a *ServiceCheckout) UpsertCheckoutLine(checkoutLine model.CheckoutLine) (*model.CheckoutLine, *model_helper.AppError) {
	checkoutLines, err := a.srv.Store.CheckoutLine().Upsert(model.CheckoutLineSlice{&checkoutLine})
	if err != nil {
		if appErr, ok := err.(*model_helper.AppError); ok {
			return nil, appErr
		}
		return nil, model_helper.NewAppError("CreateCheckoutLines", "app.checkout.failed_creating_checkoutline.app_error", nil, err.Error(), http.StatusInternalServerError)
	}

	return checkoutLines[0], nil
}

func (a *ServiceCheckout) UpsertCheckoutLines(checkoutLines model.CheckoutLineSlice) (model.CheckoutLineSlice, *model_helper.AppError) {
	checkoutLines, err := a.srv.Store.CheckoutLine().Upsert(checkoutLines)
	if err != nil {
		if appErr, ok := err.(*model_helper.AppError); ok {
			return nil, appErr
		}
		status := http.StatusInternalServerError
		if _, ok := err.(*store.ErrInvalidInput); ok {
			status = http.StatusBadRequest
		}
		return nil, model_helper.NewAppError("UpsertCheckoutLines", "app.checkout.error_bulk_create_lines.app_error", nil, err.Error(), status)
	}

	return checkoutLines, nil
}

func (s *ServiceCheckout) CheckoutLinesByOption(option model_helper.CheckoutLineFilterOptions) (model.CheckoutLineSlice, *model_helper.AppError) {
	checkoutLines, err := s.srv.Store.CheckoutLine().CheckoutLinesByOption(option)
	if err != nil {
		return nil, model_helper.NewAppError("CheckoutLinesByOption", "app.checkout.error_finding_checkout_lines_by_options.app_error", nil, err.Error(), http.StatusInternalServerError)
	}

	return checkoutLines, nil
}
