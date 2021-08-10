package order

import (
	"net/http"

	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model/order"
	"github.com/sitename/sitename/store"
)

// CreateFulfillment performs some actions then save given fulfillment
func (a *AppOrder) CreateFulfillment(fulfillment *order.Fulfillment) (*order.Fulfillment, *model.AppError) {
	// Assign an auto incremented value as a fulfillment order.
	if fulfillment.Id == "" {
		fulfillmentsByOrder, appErr := a.FulfillmentsByOrderID(fulfillment.OrderID)
		if appErr != nil {
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

	fulfillment, err := a.Srv().Store.Fulfillment().Save(fulfillment)
	if err != nil {
		return nil, model.NewAppError("CreateFulfillment", "app.order.error_saving_fulfillment.app_error", nil, err.Error(), http.StatusInternalServerError)
	}

	return fulfillment, nil
}

// FulfillmentsByOrderID returns all fulfillments belong to given order
func (a *AppOrder) FulfillmentsByOrderID(orderID string) ([]*order.Fulfillment, *model.AppError) {
	fulfillmentsByOrder, err := a.Srv().SqlStore.
		Fulfillment().
		FilterByoption(&order.FulfillmentFilterOption{
			OrderID: &model.StringFilter{
				StringOption: &model.StringOption{
					Eq: orderID,
				},
			},
		})
	if err != nil {
		return nil, store.AppErrorFromDatabaseLookupError("FulfillmentsByOrderID", "app.order.fulfillments_by_option.app_error", err)
	}

	return fulfillmentsByOrder, nil
}

func (a *AppOrder) FulfillmentLinesByFulfillmentID(fulfillmentID string) ([]*order.FulfillmentLine, *model.AppError) {
	fulfillmentLines, err := a.Srv().Store.FulfillmentLine().FilterbyOption(&order.FulfillmentLineFilterOption{
		FulfillmentID: &model.StringFilter{
			StringOption: &model.StringOption{
				Eq: fulfillmentID,
			},
		},
	})

	if err != nil {
		return nil, store.AppErrorFromDatabaseLookupError("FulfillmentLinesByFulfillmentID", "app.order.error_finding_fulfillment_lines_by_fulfillment.app_error", err)
	}

	return fulfillmentLines, nil
}
