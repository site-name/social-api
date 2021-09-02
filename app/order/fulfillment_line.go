package order

import (
	"net/http"

	"github.com/mattermost/gorp"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model/order"
	"github.com/sitename/sitename/store"
)

// FulfillmentLinesByOption returns all fulfillment lines by option
func (a *ServiceOrder) FulfillmentLinesByOption(option *order.FulfillmentLineFilterOption) ([]*order.FulfillmentLine, *model.AppError) {
	fulfillmentLines, err := a.srv.Store.FulfillmentLine().FilterbyOption(option)
	if err != nil {
		return nil, store.AppErrorFromDatabaseLookupError("FulfillmentLinesByOption", "app.order.error_finding_fulfillment_lines_by_option.app_error", err)
	}

	return fulfillmentLines, nil
}

// BulkUpsertFulfillmentLines performs bulk upsert given fulfillment lines and returns them
func (a *ServiceOrder) BulkUpsertFulfillmentLines(transaction *gorp.Transaction, fulfillmentLines []*order.FulfillmentLine) ([]*order.FulfillmentLine, *model.AppError) {
	fulfillmentLines, err := a.srv.Store.FulfillmentLine().BulkUpsert(transaction, fulfillmentLines)
	if err != nil {
		return nil, model.NewAppError("BulkUpsertFulfillmentLines", "app.order.error_bulk_creating_fulfillment_lines.app_error", nil, err.Error(), http.StatusInternalServerError)
	}

	return fulfillmentLines, nil
}

// DeleteFulfillmentLinesByOption tells store to delete fulfillment lines filtered by given option
func (a *ServiceOrder) DeleteFulfillmentLinesByOption(transaction *gorp.Transaction, option *order.FulfillmentLineFilterOption) *model.AppError {
	err := a.srv.Store.FulfillmentLine().DeleteFulfillmentLinesByOption(transaction, option)
	if err != nil {
		return model.NewAppError("DeleteFulfillmentLinesByOption", "app.order.error_deleting_fulfillment_lines_by_option.app_error", nil, err.Error(), http.StatusInternalServerError)
	}

	return nil
}
