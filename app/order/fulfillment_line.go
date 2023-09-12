package order

import (
	"net/http"

	"github.com/sitename/sitename/model"
	"gorm.io/gorm"
)

// FulfillmentLinesByOption returns all fulfillment lines by option
func (a *ServiceOrder) FulfillmentLinesByOption(option *model.FulfillmentLineFilterOption) (model.FulfillmentLines, *model.AppError) {
	fulfillmentLines, err := a.srv.Store.FulfillmentLine().FilterbyOption(option)
	if err != nil {
		return nil, model.NewAppError("FulfillmentLinesByOption", "app.order.error_finding_fulfillment_lines_by_options.app_error", nil, err.Error(), http.StatusInternalServerError)
	}

	return fulfillmentLines, nil
}

// BulkUpsertFulfillmentLines performs bulk upsert given fulfillment lines and returns them
func (a *ServiceOrder) BulkUpsertFulfillmentLines(transaction *gorm.DB, fulfillmentLines []*model.FulfillmentLine) ([]*model.FulfillmentLine, *model.AppError) {
	fulfillmentLines, err := a.srv.Store.FulfillmentLine().BulkUpsert(transaction, fulfillmentLines)
	if err != nil {
		return nil, model.NewAppError("BulkUpsertFulfillmentLines", "app.order.error_bulk_creating_fulfillment_lines.app_error", nil, err.Error(), http.StatusInternalServerError)
	}

	return fulfillmentLines, nil
}

// DeleteFulfillmentLinesByOption tells store to delete fulfillment lines filtered by given option
func (a *ServiceOrder) DeleteFulfillmentLinesByOption(transaction *gorm.DB, option *model.FulfillmentLineFilterOption) *model.AppError {
	err := a.srv.Store.FulfillmentLine().DeleteFulfillmentLinesByOption(transaction, option)
	if err != nil {
		return model.NewAppError("DeleteFulfillmentLinesByOption", "app.order.error_deleting_fulfillment_lines_by_option.app_error", nil, err.Error(), http.StatusInternalServerError)
	}

	return nil
}
