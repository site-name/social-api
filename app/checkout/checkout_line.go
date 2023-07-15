package checkout

import (
	"net/http"

	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/store"
	"github.com/sitename/sitename/store/store_iface"
)

func (a *ServiceCheckout) CheckoutLinesByCheckoutToken(checkoutToken string) ([]*model.CheckoutLine, *model.AppError) {
	lines, err := a.srv.Store.CheckoutLine().CheckoutLinesByCheckoutID(checkoutToken)
	if err != nil {
		return nil, model.NewAppError("CheckoutLinesByCheckoutToken", "app.checkout.checkout_lines_by_checkout.app_error", nil, err.Error(), http.StatusInternalServerError)
	}

	return lines, nil
}

func (a *ServiceCheckout) DeleteCheckoutLines(transaction store_iface.SqlxExecutor, checkoutLineIDs []string) *model.AppError {
	err := a.srv.Store.CheckoutLine().DeleteLines(transaction, checkoutLineIDs)
	if err != nil {
		return model.NewAppError("DeleteCheckoutLines", "app.checkout.error_deleting_checkoutlines.app_error", nil, err.Error(), http.StatusInternalServerError)
	}

	return nil
}

func (a *ServiceCheckout) UpsertCheckoutLine(checkoutLine *model.CheckoutLine) (*model.CheckoutLine, *model.AppError) {
	checkoutLine, err := a.srv.Store.CheckoutLine().Upsert(checkoutLine)
	if err != nil {
		if appErr, ok := err.(*model.AppError); ok {
			return nil, appErr
		}
		return nil, model.NewAppError("CreateCheckoutLines", "app.checkout.failed_creating_checkoutline.app_error", nil, err.Error(), http.StatusInternalServerError)
	}

	return checkoutLine, nil
}

func (a *ServiceCheckout) BulkCreateCheckoutLines(checkoutLines []*model.CheckoutLine) ([]*model.CheckoutLine, *model.AppError) {
	checkoutLines, err := a.srv.Store.CheckoutLine().BulkCreate(checkoutLines)
	if err != nil {
		if appErr, ok := err.(*model.AppError); ok {
			return nil, appErr
		}
		status := http.StatusInternalServerError
		if _, ok := err.(*store.ErrInvalidInput); ok {
			status = http.StatusBadRequest
		}
		return nil, model.NewAppError("BulkCreateCheckoutLines", "app.checkout.error_bulk_create_lines.app_error", nil, err.Error(), status)
	}

	return checkoutLines, nil
}

func (a *ServiceCheckout) BulkUpdateCheckoutLines(checkoutLines []*model.CheckoutLine) *model.AppError {
	err := a.srv.Store.CheckoutLine().BulkUpdate(checkoutLines)
	if err != nil {
		if appErr, ok := err.(*model.AppError); ok {
			return appErr
		}
		status := http.StatusInternalServerError
		if _, ok := err.(*store.ErrInvalidInput); ok {
			status = http.StatusBadRequest
		}
		return model.NewAppError("BulkUpdateCheckoutLines", "app.checkout.error_bulk_update_lines.app_error", nil, err.Error(), status)
	}

	return nil
}

// CheckoutLinesByOption returns a list of checkout lines filtered using given option
func (s *ServiceCheckout) CheckoutLinesByOption(option *model.CheckoutLineFilterOption) ([]*model.CheckoutLine, *model.AppError) {
	checkoutLines, err := s.srv.Store.CheckoutLine().CheckoutLinesByOption(option)
	if err != nil {
		return nil, model.NewAppError("CheckoutLinesByOption", "app.checkout.error_finding_checkout_lines_by_options.app_error", nil, err.Error(), http.StatusInternalServerError)
	}

	return checkoutLines, nil
}
