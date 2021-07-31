package checkout

import (
	"net/http"

	"github.com/sitename/sitename/app"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model/checkout"
	"github.com/sitename/sitename/store"
)

func (a *AppCheckout) CheckoutLinesByCheckoutToken(checkoutToken string) ([]*checkout.CheckoutLine, *model.AppError) {
	lines, err := a.app.Srv().Store.CheckoutLine().CheckoutLinesByCheckoutID(checkoutToken)
	if err != nil {
		return nil, store.AppErrorFromDatabaseLookupError("CheckoutLinesByCheckoutID", "app.checkout.checkout_lines_by_checkout.app_error", err)
	}

	return lines, nil
}

func (a *AppCheckout) DeleteCheckoutLines(checkoutLineIDs []string) *model.AppError {
	err := a.app.Srv().Store.CheckoutLine().DeleteLines(checkoutLineIDs)
	if err != nil {
		var errID string
		var errArgs map[string]interface{}
		statusCode := http.StatusInternalServerError

		if invlErr, ok := err.(*store.ErrInvalidInput); ok {
			errID = app.InvalidArgumentAppErrorID
			statusCode = http.StatusBadRequest
			errArgs = map[string]interface{}{"Fields": invlErr.Field}
		} else {
			errID = "app.checkout.error_deleting_checkoutlines.app_error"
		}
		return model.NewAppError("DeleteCheckoutLines", errID, errArgs, err.Error(), statusCode)
	}

	return nil
}

func (a *AppCheckout) UpsertCheckoutLine(checkoutLine *checkout.CheckoutLine) (*checkout.CheckoutLine, *model.AppError) {
	checkoutLine, err := a.app.Srv().SqlStore.CheckoutLine().Upsert(checkoutLine)
	if err != nil {
		if appErr, ok := err.(*model.AppError); ok {
			return nil, appErr
		}
		return nil, model.NewAppError("CreateCheckoutLines", "app.checkout.failed_creating_checkoutline.app_error", nil, err.Error(), http.StatusInternalServerError)
	}

	return checkoutLine, nil
}

func (a *AppCheckout) BulkCreateCheckoutLines(checkoutLines []*checkout.CheckoutLine) ([]*checkout.CheckoutLine, *model.AppError) {
	checkoutLines, err := a.app.Srv().Store.CheckoutLine().BulkCreate(checkoutLines)
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

func (a *AppCheckout) BulkUpdateCheckoutLines(checkoutLines []*checkout.CheckoutLine) *model.AppError {
	err := a.app.Srv().Store.CheckoutLine().BulkUpdate(checkoutLines)
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
