package checkout

import (
	"net/http"

	"github.com/sitename/sitename/app"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model/checkout"
	"github.com/sitename/sitename/store"
)

func (a *ServiceCheckout) CheckoutLinesByCheckoutToken(checkoutToken string) ([]*checkout.CheckoutLine, *model.AppError) {
	lines, err := a.srv.Store.CheckoutLine().CheckoutLinesByCheckoutID(checkoutToken)
	if err != nil {
		return nil, store.AppErrorFromDatabaseLookupError("CheckoutLinesByCheckoutID", "app.checkout.checkout_lines_by_checkout.app_error", err)
	}

	return lines, nil
}

func (a *ServiceCheckout) DeleteCheckoutLines(checkoutLineIDs []string) *model.AppError {
	// validate id list
	for _, id := range checkoutLineIDs {
		if !model.IsValidId(id) {
			return model.NewAppError("DeleteCheckoutLines", app.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "checkoutLineIDs"}, "", http.StatusBadRequest)
		}
	}

	err := a.srv.Store.CheckoutLine().DeleteLines(checkoutLineIDs)
	if err != nil {
		return model.NewAppError("DeleteCheckoutLines", "app.checkout.error_deleting_checkoutlines.app_error", nil, err.Error(), http.StatusInternalServerError)
	}

	return nil
}

func (a *ServiceCheckout) UpsertCheckoutLine(checkoutLine *checkout.CheckoutLine) (*checkout.CheckoutLine, *model.AppError) {
	checkoutLine, err := a.srv.SqlStore.CheckoutLine().Upsert(checkoutLine)
	if err != nil {
		if appErr, ok := err.(*model.AppError); ok {
			return nil, appErr
		}
		return nil, model.NewAppError("CreateCheckoutLines", "app.checkout.failed_creating_checkoutline.app_error", nil, err.Error(), http.StatusInternalServerError)
	}

	return checkoutLine, nil
}

func (a *ServiceCheckout) BulkCreateCheckoutLines(checkoutLines []*checkout.CheckoutLine) ([]*checkout.CheckoutLine, *model.AppError) {
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

func (a *ServiceCheckout) BulkUpdateCheckoutLines(checkoutLines []*checkout.CheckoutLine) *model.AppError {
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
