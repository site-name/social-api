package order

import (
	"net/http"

	"github.com/mattermost/squirrel"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model_helper"
	"github.com/sitename/sitename/store"
	"gorm.io/gorm"
)

// FulfillmentsByOption returns a list of fulfillments be given options
func (a *ServiceOrder) FulfillmentsByOption(option *model.FulfillmentFilterOption) (model.Fulfillments, *model_helper.AppError) {
	fulfillments, err := a.srv.Store.Fulfillment().FilterByOption(option)
	if err != nil {
		return nil, model_helper.NewAppError("FulfillmentsByOption", "app.model.error_finding_fulfillments_by_option.app_error", nil, err.Error(), http.StatusInternalServerError)
	}

	return fulfillments, nil
}

// UpsertFulfillment performs some actions then save given fulfillment
func (a *ServiceOrder) UpsertFulfillment(transaction *gorm.DB, fulfillment *model.Fulfillment) (*model.Fulfillment, *model_helper.AppError) {
	// Assign an auto incremented value as a fulfillment order.
	if fulfillment.Id == "" {
		fulfillmentsByOrder, appErr := a.FulfillmentsByOption(&model.FulfillmentFilterOption{
			Conditions: squirrel.Eq{model.FulfillmentTableName + ".OrderID": fulfillment.OrderID},
		})
		if appErr != nil {
			if appErr.StatusCode == http.StatusInternalServerError { // returns immediately if error was caused by system
				return nil, appErr
			}
			// this means the order has no fulfillment, this is the first one and not saved yet
			fulfillment.FulfillmentOrder = 1
		} else {
			var max int
			for _, fulfillment := range fulfillmentsByOrder {
				if num := fulfillment.FulfillmentOrder; num > max {
					max = num
				}
			}

			fulfillment.FulfillmentOrder = max + 1
		}
	}

	upsertedFulfillment, err := a.srv.Store.Fulfillment().Upsert(transaction, fulfillment)
	if err != nil {
		if appErr, ok := err.(*model_helper.AppError); ok {
			return nil, appErr
		}
		if errNotFound, ok := err.(*store.ErrNotFound); ok { // this happens when update an unexisted instance
			return nil, model_helper.NewAppError("UpsertFulfillment", model_helper.InvalidArgumentAppErrorID, map[string]any{"Fields": "Id"}, errNotFound.Error(), http.StatusBadRequest)
		}
		return nil, model_helper.NewAppError("UpsertFulfillment", "app.order.error_saving_fulfillment.app_error", nil, err.Error(), http.StatusInternalServerError)
	}

	return upsertedFulfillment, nil
}

// FulfillmentByOption returns 1 fulfillment filtered using given options
func (a *ServiceOrder) FulfillmentByOption(option *model.FulfillmentFilterOption) (*model.Fulfillment, *model_helper.AppError) {
	fulfillment, err := a.srv.Store.Fulfillment().GetByOption(option)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if _, ok := err.(*store.ErrNotFound); ok {
			statusCode = http.StatusNotFound
		}
		return nil, model_helper.NewAppError("FulfillmentByOption", "app.order.error_finding_fulfillment_by_option.app_error", nil, err.Error(), statusCode)
	}

	return fulfillment, nil
}

// BulkDeleteFulfillments tells store to delete fulfillments that satisfy given option
func (a *ServiceOrder) BulkDeleteFulfillments(transaction *gorm.DB, fulfillments model.Fulfillments) *model_helper.AppError {
	err := a.srv.Store.Fulfillment().BulkDeleteFulfillments(transaction, fulfillments)
	if err != nil {
		return model_helper.NewAppError("BulkDeleteFulfillments", "app.order.error_deleting_fulfillments.app_error", nil, err.Error(), http.StatusInternalServerError)
	}

	return nil
}
