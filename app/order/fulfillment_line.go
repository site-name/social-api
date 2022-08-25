package order

import (
	"net/http"

	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model/order"
	"github.com/sitename/sitename/store/store_iface"
)

// FulfillmentLinesByOption returns all fulfillment lines by option
func (a *ServiceOrder) FulfillmentLinesByOption(option *order.FulfillmentLineFilterOption) (order.FulfillmentLines, *model.AppError) {
	fulfillmentLines, err := a.srv.Store.FulfillmentLine().FilterbyOption(option)
	var (
		statusCode int
		errMessage string
	)
	if err != nil {
		statusCode = http.StatusInternalServerError
		errMessage = err.Error()
	} else if len(fulfillmentLines) == 0 {
		statusCode = http.StatusNotFound
	}

	if statusCode != 0 {
		return nil, model.NewAppError("FulfillmentLinesByOption", "app.order.error_finding_fulfillment_lines_by_options.app_error", nil, errMessage, statusCode)
	}

	return fulfillmentLines, nil
}

// BulkUpsertFulfillmentLines performs bulk upsert given fulfillment lines and returns them
func (a *ServiceOrder) BulkUpsertFulfillmentLines(transaction store_iface.SqlxTxExecutor, fulfillmentLines []*order.FulfillmentLine) ([]*order.FulfillmentLine, *model.AppError) {
	fulfillmentLines, err := a.srv.Store.FulfillmentLine().BulkUpsert(transaction, fulfillmentLines)
	if err != nil {
		return nil, model.NewAppError("BulkUpsertFulfillmentLines", "app.order.error_bulk_creating_fulfillment_lines.app_error", nil, err.Error(), http.StatusInternalServerError)
	}

	return fulfillmentLines, nil
}

// DeleteFulfillmentLinesByOption tells store to delete fulfillment lines filtered by given option
func (a *ServiceOrder) DeleteFulfillmentLinesByOption(transaction store_iface.SqlxTxExecutor, option *order.FulfillmentLineFilterOption) *model.AppError {
	err := a.srv.Store.FulfillmentLine().DeleteFulfillmentLinesByOption(transaction, option)
	if err != nil {
		return model.NewAppError("DeleteFulfillmentLinesByOption", "app.order.error_deleting_fulfillment_lines_by_option.app_error", nil, err.Error(), http.StatusInternalServerError)
	}

	return nil
}
