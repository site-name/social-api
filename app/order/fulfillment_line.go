package order

import (
	"net/http"

	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model_helper"
	"github.com/volatiletech/sqlboiler/boil"
)

// FulfillmentLinesByOption returns all fulfillment lines by option
func (a *ServiceOrder) FulfillmentLinesByOption(option *model.FulfillmentLineFilterOption) (model.FulfillmentLines, *model_helper.AppError) {
	fulfillmentLines, err := a.srv.Store.FulfillmentLine().FilterbyOption(option)
	if err != nil {
		return nil, model_helper.NewAppError("FulfillmentLinesByOption", "app.order.error_finding_fulfillment_lines_by_options.app_error", nil, err.Error(), http.StatusInternalServerError)
	}

	return fulfillmentLines, nil
}

// BulkUpsertFulfillmentLines performs bulk upsert given fulfillment lines and returns them
func (a *ServiceOrder) BulkUpsertFulfillmentLines(transaction boil.ContextTransactor, fulfillmentLines []*model.FulfillmentLine) ([]*model.FulfillmentLine, *model_helper.AppError) {
	fulfillmentLines, err := a.srv.Store.FulfillmentLine().BulkUpsert(transaction, fulfillmentLines)
	if err != nil {
		return nil, model_helper.NewAppError("BulkUpsertFulfillmentLines", "app.order.error_bulk_creating_fulfillment_lines.app_error", nil, err.Error(), http.StatusInternalServerError)
	}

	return fulfillmentLines, nil
}

// DeleteFulfillmentLinesByOption tells store to delete fulfillment lines filtered by given option
func (a *ServiceOrder) DeleteFulfillmentLinesByOption(transaction boil.ContextTransactor, option *model.FulfillmentLineFilterOption) *model_helper.AppError {
	err := a.srv.Store.FulfillmentLine().DeleteFulfillmentLinesByOption(transaction, option)
	if err != nil {
		return model_helper.NewAppError("DeleteFulfillmentLinesByOption", "app.order.error_deleting_fulfillment_lines_by_option.app_error", nil, err.Error(), http.StatusInternalServerError)
	}

	return nil
}
