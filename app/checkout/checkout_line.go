package checkout

import (
	"net/http"

	"github.com/mattermost/squirrel"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model_helper"
	"github.com/sitename/sitename/store"
	"gorm.io/gorm"
)

func (a *ServiceCheckout) CheckoutLinesByCheckoutToken(checkoutToken string) ([]*model.CheckoutLine, *model_helper.AppError) {
	return a.CheckoutLinesByOption(&model.CheckoutLineFilterOption{
		Conditions: squirrel.Eq{model.CheckoutLineTableName + ".CheckoutID": checkoutToken},
	})
}

func (a *ServiceCheckout) DeleteCheckoutLines(transaction *gorm.DB, checkoutLineIDs []string) *model_helper.AppError {
	err := a.srv.Store.CheckoutLine().DeleteLines(transaction, checkoutLineIDs)
	if err != nil {
		return model_helper.NewAppError("DeleteCheckoutLines", "app.checkout.error_deleting_checkoutlines.app_error", nil, err.Error(), http.StatusInternalServerError)
	}

	return nil
}

func (a *ServiceCheckout) UpsertCheckoutLine(checkoutLine *model.CheckoutLine) (*model.CheckoutLine, *model_helper.AppError) {
	checkoutLine, err := a.srv.Store.CheckoutLine().Upsert(checkoutLine)
	if err != nil {
		if appErr, ok := err.(*model_helper.AppError); ok {
			return nil, appErr
		}
		return nil, model_helper.NewAppError("CreateCheckoutLines", "app.checkout.failed_creating_checkoutline.app_error", nil, err.Error(), http.StatusInternalServerError)
	}

	return checkoutLine, nil
}

func (a *ServiceCheckout) BulkCreateCheckoutLines(checkoutLines []*model.CheckoutLine) ([]*model.CheckoutLine, *model_helper.AppError) {
	checkoutLines, err := a.srv.Store.CheckoutLine().BulkCreate(checkoutLines)
	if err != nil {
		if appErr, ok := err.(*model_helper.AppError); ok {
			return nil, appErr
		}
		status := http.StatusInternalServerError
		if _, ok := err.(*store.ErrInvalidInput); ok {
			status = http.StatusBadRequest
		}
		return nil, model_helper.NewAppError("BulkCreateCheckoutLines", "app.checkout.error_bulk_create_lines.app_error", nil, err.Error(), status)
	}

	return checkoutLines, nil
}

func (a *ServiceCheckout) BulkUpdateCheckoutLines(checkoutLines []*model.CheckoutLine) *model_helper.AppError {
	err := a.srv.Store.CheckoutLine().BulkUpdate(checkoutLines)
	if err != nil {
		if appErr, ok := err.(*model_helper.AppError); ok {
			return appErr
		}
		status := http.StatusInternalServerError
		if _, ok := err.(*store.ErrInvalidInput); ok {
			status = http.StatusBadRequest
		}
		return model_helper.NewAppError("BulkUpdateCheckoutLines", "app.checkout.error_bulk_update_lines.app_error", nil, err.Error(), status)
	}

	return nil
}

// CheckoutLinesByOption returns a list of checkout lines filtered using given option
func (s *ServiceCheckout) CheckoutLinesByOption(option *model.CheckoutLineFilterOption) ([]*model.CheckoutLine, *model_helper.AppError) {
	checkoutLines, err := s.srv.Store.CheckoutLine().CheckoutLinesByOption(option)
	if err != nil {
		return nil, model_helper.NewAppError("CheckoutLinesByOption", "app.checkout.error_finding_checkout_lines_by_options.app_error", nil, err.Error(), http.StatusInternalServerError)
	}

	return checkoutLines, nil
}
