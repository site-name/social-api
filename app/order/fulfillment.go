package order

import (
	"net/http"

	"github.com/Masterminds/squirrel"
	"github.com/sitename/sitename/app"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model/order"
	"github.com/sitename/sitename/store"
	"github.com/sitename/sitename/store/store_iface"
)

// FulfillmentsByOption returns a list of fulfillments be given options
func (a *ServiceOrder) FulfillmentsByOption(transaction store_iface.SqlxTxExecutor, option *order.FulfillmentFilterOption) (order.Fulfillments, *model.AppError) {
	fulfillments, err := a.srv.Store.Fulfillment().FilterByOption(transaction, option)
	if err != nil {
		return nil, store.AppErrorFromDatabaseLookupError("FulfillmentsByOption", "app.order.error_finding_fulfillments_by_option.app_error", err)
	}

	return fulfillments, nil
}

// UpsertFulfillment performs some actions then save given fulfillment
func (a *ServiceOrder) UpsertFulfillment(transaction store_iface.SqlxTxExecutor, fulfillment *order.Fulfillment) (*order.Fulfillment, *model.AppError) {
	// Assign an auto incremented value as a fulfillment order.
	if fulfillment.Id == "" {
		fulfillmentsByOrder, appErr := a.FulfillmentsByOption(nil, &order.FulfillmentFilterOption{
			OrderID: squirrel.Eq{store.FulfillmentTableName + ".OrderID": fulfillment.OrderID},
		})
		if appErr != nil {
			if appErr.StatusCode == http.StatusInternalServerError { // returns immediately if error was caused by system
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

	fulfillment, err := a.srv.Store.Fulfillment().Upsert(transaction, fulfillment)
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

// FulfillmentByOption returns 1 fulfillment filtered using given options
func (a *ServiceOrder) FulfillmentByOption(transaction store_iface.SqlxTxExecutor, option *order.FulfillmentFilterOption) (*order.Fulfillment, *model.AppError) {
	fulfillment, err := a.srv.Store.Fulfillment().GetByOption(transaction, option)
	if err != nil {
		return nil, store.AppErrorFromDatabaseLookupError("FulfillmentByOption", "app.order.error_finding_fulfillment_by_option.app_error", err)
	}

	return fulfillment, nil
}

// GetOrCreateFulfillment take a filtering option, trys finding a fulfillment with given option.
// If a fulfillment found, returns it. Otherwise, creates a new one then returns it.
func (a *ServiceOrder) GetOrCreateFulfillment(transaction store_iface.SqlxTxExecutor, option *order.FulfillmentFilterOption) (*order.Fulfillment, *model.AppError) {
	fulfillmentByOption, appErr := a.FulfillmentByOption(transaction, option)
	if appErr != nil {
		if appErr.StatusCode == http.StatusInternalServerError {
			return nil, appErr
		}

		// not found, create a new one:
		var parser = func(expr squirrel.Sqlizer) interface{} {
			if expr == nil {
				return nil
			}

			if equal, ok := expr.(squirrel.Eq); ok && len(equal) > 0 {
				// get first value only
				index := 0
				for _, value := range equal {
					if index == 0 {
						return value
					}
					return nil
				}
			}

			return nil
		}

		fulfillmentByOption = new(order.Fulfillment)
		// parse options. if any option is provided, take its Eq property:
		if option.Id != nil {
			value := parser(option.Id)
			if value != nil {
				fulfillmentByOption.Id = value.(string)
			}
		}
		if option.OrderID != nil {
			value := parser(option.OrderID)
			if value != nil {
				fulfillmentByOption.OrderID = value.(string)
			}
		}
		if option.Status != nil {
			value := parser(option.Status)
			if value != nil {
				fulfillmentByOption.Status = order.FulfillmentStatus(value.(string))
			}
		}

		fulfillmentByOption, appErr = a.UpsertFulfillment(transaction, fulfillmentByOption)
		if appErr != nil {
			return nil, appErr
		}

		return fulfillmentByOption, nil
	}

	return fulfillmentByOption, nil
}

// BulkDeleteFulfillments tells store to delete fulfillments that satisfy given option
func (a *ServiceOrder) BulkDeleteFulfillments(transaction store_iface.SqlxTxExecutor, fulfillments order.Fulfillments) *model.AppError {
	err := a.srv.Store.Fulfillment().BulkDeleteFulfillments(transaction, fulfillments)
	if err != nil {
		return model.NewAppError("BulkDeleteFulfillments", "app.order.error_deleting_fulfillments.app_error", nil, err.Error(), http.StatusInternalServerError)
	}

	return nil
}
