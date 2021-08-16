package order

import (
	"net/http"

	"github.com/sitename/sitename/app"
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
			if appErr.StatusCode == http.StatusInternalServerError { // returns immediately if error was caused by system
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
		if appErr, ok := err.(*model.AppError); ok {
			return nil, appErr
		}
		if errNotFound, ok := err.(*store.ErrNotFound); ok { // this happens when update an unexisted instance
			return nil, model.NewAppError("UpsertFulfillment", app.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "Id"}, errNotFound.Error(), http.StatusBadRequest)
		}
		return nil, model.NewAppError("UpsertFulfillment", "app.order.error_saving_fulfillment.app_error", nil, err.Error(), http.StatusInternalServerError)
	}

	return fulfillment, nil
}

// GetOrCreateFulfillment take a filtering option, trys finding a fulfillment with given option.
// If a fulfillment found, returns it. Otherwise, creates a new one then returns it.
func (a *AppOrder) GetOrCreateFulfillment(option *order.FulfillmentFilterOption) (*order.Fulfillment, *model.AppError) {
	fulfillmentByOption, err := a.Srv().Store.Fulfillment().GetByOption(option)
	if err != nil {
		if _, ok := err.(*store.ErrNotFound); ok { // fulfillment not found. Creating a new one
			fulfillmentByOption = new(order.Fulfillment)
			// parse options. if any option is provided, take its Eq property:
			if option.Id != nil {
				fulfillmentByOption.Id = option.Id.Eq
			}
			if option.OrderID != nil {
				fulfillmentByOption.OrderID = option.OrderID.Eq
			}
			if option.Status != nil {
				fulfillmentByOption.Status = option.Status.Eq
			}

			fulfillmentByOption, appErr := a.UpsertFulfillment(fulfillmentByOption)
			if appErr != nil {
				appErr.Where = "GetOrCreateFulfillment"
				return nil, appErr
			}

			return fulfillmentByOption, nil
		}
		return nil, model.NewAppError("GetOrCreateFulfillment", "app.order.error_finding_fulfillment_by_option.app_error", nil, err.Error(), http.StatusInternalServerError)
	}

	return fulfillmentByOption, nil
}
