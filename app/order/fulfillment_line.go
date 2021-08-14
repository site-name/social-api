package order

import (
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
