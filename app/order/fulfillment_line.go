package order

import (
	"net/http"

	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model/order"
	"github.com/sitename/sitename/store"
)

// FulfillmentLinesByOption returns all fulfillment lines by option
func (a *AppOrder) FulfillmentLinesByOption(option *order.FulfillmentLineFilterOption) ([]*order.FulfillmentLine, *model.AppError) {
	fulfillmentLines, err := a.Srv().Store.FulfillmentLine().FilterbyOption(option)
	if err != nil {
		return nil, store.AppErrorFromDatabaseLookupError("FulfillmentLinesByOption", "app.order.error_finding_fulfillment_lines_by_option.app_error", err)
	}

	return fulfillmentLines, nil
}

// BulkUpsertFulfillmentLines performs bulk upsert given fulfillment lines and returns them
func (a *AppOrder) BulkUpsertFulfillmentLines(fulfillmentLines []*order.FulfillmentLine) ([]*order.FulfillmentLine, *model.AppError) {
	fulfillmentLines, err := a.Srv().Store.FulfillmentLine().BulkUpsert(fulfillmentLines)
	if err != nil {
		return nil, model.NewAppError("BulkUpsertFulfillmentLines", "app.order.error_bulk_creating_fulfillment_lines.app_error", nil, err.Error(), http.StatusInternalServerError)
	}

	return fulfillmentLines, nil
}

// DeleteFulfillmentLinesByOption tells store to delete fulfillment lines filtered by given option
func (a *AppOrder) DeleteFulfillmentLinesByOption(option *order.FulfillmentLineFilterOption) *model.AppError {
	err := a.Srv().Store.FulfillmentLine().DeleteFulfillmentLinesByOption(option)
	if err != nil {
		return model.NewAppError("DeleteFulfillmentLinesByOption", "app.order.error_deleting_fulfillment_lines_by_option.app_error", nil, err.Error(), http.StatusInternalServerError)
	}

	return nil
}
