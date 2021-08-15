package order

import (
	"net/http"

	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model/order"
	"github.com/sitename/sitename/store"
)

// FulfillmentsByOption returns a list of fulfillments be given options
func (a *AppOrder) FulfillmentsByOption(option *order.FulfillmentFilterOption) ([]*order.Fulfillment, *model.AppError) {
	fulfillments, err := a.Srv().Store.Fulfillment().FilterByoption(option)
	if err != nil {
		return nil, store.AppErrorFromDatabaseLookupError("FulfillmentsByOption", "app.order.error_finding_fulfillments_by_option.app_error", err)
	}

	return fulfillments, nil
}

// UpsertFulfillment performs some actions then save given fulfillment
func (a *AppOrder) UpsertFulfillment(fulfillment *order.Fulfillment) (*order.Fulfillment, *model.AppError) {
	// Assign an auto incremented value as a fulfillment order.
	if fulfillment.Id == "" {
		fulfillmentsByOrder, appErr := a.FulfillmentsByOption(&order.FulfillmentFilterOption{
			OrderID: &model.StringFilter{
				StringOption: &model.StringOption{
					Eq: fulfillment.OrderID,
				},
			},
		})
		if appErr != nil {
			if appErr.StatusCode == http.StatusInternalServerError {
				appErr.Where = "UpsertFulfillment"
				return nil, appErr
			}
			// this means the order has no fulfillment, this is the first one and not saved yet
			fulfillment.FulfillmentOrder = 1
		} else {
			var max uint
			for _, fulfillment := range fulfillmentsByOrder {
				if num := fulfillment.FulfillmentOrder; num > max {
					max = num
				}
			}

			fulfillment.FulfillmentOrder = max + 1
		}
	}

	fulfillment, err := a.Srv().Store.Fulfillment().Upsert(fulfillment)
	if err != nil {
		return nil, model.NewAppError("CreateFulfillment", "app.order.error_saving_fulfillment.app_error", nil, err.Error(), http.StatusInternalServerError)
	}

	return fulfillment, nil
}
