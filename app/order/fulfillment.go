package order

import (
	"net/http"

	"github.com/Masterminds/squirrel"
	"github.com/sitename/sitename/app"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/store"
	"github.com/sitename/sitename/store/store_iface"
)

// FulfillmentsByOption returns a list of fulfillments be given options
func (a *ServiceOrder) FulfillmentsByOption(transaction store_iface.SqlxTxExecutor, option *model.FulfillmentFilterOption) (model.Fulfillments, *model.AppError) {
	fulfillments, err := a.srv.Store.Fulfillment().FilterByOption(transaction, option)
	if err != nil {
		return nil, model.NewAppError("FulfillmentsByOption", "app.model.error_finding_fulfillments_by_option.app_error", nil, err.Error(), http.StatusInternalServerError)
	}

	return fulfillments, nil
}

// UpsertFulfillment performs some actions then save given fulfillment
func (a *ServiceOrder) UpsertFulfillment(transaction store_iface.SqlxTxExecutor, fulfillment *model.Fulfillment) (*model.Fulfillment, *model.AppError) {
	// Assign an auto incremented value as a fulfillment order.
	if fulfillment.Id == "" {
		fulfillmentsByOrder, appErr := a.FulfillmentsByOption(nil, &model.FulfillmentFilterOption{
			OrderID: squirrel.Eq{store.FulfillmentTableName + ".OrderID": fulfillment.OrderID},
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
		if appErr, ok := err.(*model.AppError); ok {
			return nil, appErr
		}
		if errNotFound, ok := err.(*store.ErrNotFound); ok { // this happens when update an unexisted instance
			return nil, model.NewAppError("UpsertFulfillment", app.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "Id"}, errNotFound.Error(), http.StatusBadRequest)
		}
		return nil, model.NewAppError("UpsertFulfillment", "app.order.error_saving_fulfillment.app_error", nil, err.Error(), http.StatusInternalServerError)
	}

	return upsertedFulfillment, nil
}

// FulfillmentByOption returns 1 fulfillment filtered using given options
func (a *ServiceOrder) FulfillmentByOption(transaction store_iface.SqlxTxExecutor, option *model.FulfillmentFilterOption) (*model.Fulfillment, *model.AppError) {
	fulfillment, err := a.srv.Store.Fulfillment().GetByOption(transaction, option)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if _, ok := err.(*store.ErrNotFound); ok {
			statusCode = http.StatusNotFound
		}
		return nil, model.NewAppError("FulfillmentByOption", "app.order.error_finding_fulfillment_by_option.app_error", nil, err.Error(), statusCode)
	}

	return fulfillment, nil
}

// GetOrCreateFulfillment take a filtering option, trys finding a fulfillment with given option.
// If a fulfillment found, returns it. Otherwise, creates a new one then returns it.
func (a *ServiceOrder) GetOrCreateFulfillment(transaction store_iface.SqlxTxExecutor, option *model.FulfillmentFilterOption) (*model.Fulfillment, *model.AppError) {
	fulfillmentByOption, appErr := a.FulfillmentByOption(transaction, option)
	if appErr != nil {
		if appErr.StatusCode == http.StatusInternalServerError {
			return nil, appErr
		}

		upsertFulfillment := &model.Fulfillment{}

		// parse options. if any option is provided, take its Eq property:
		if option.Id != nil {
			eq, isEqual := option.Id.(squirrel.Eq)
			if isEqual && eq != nil {
				eqExpr := eq[store.FulfillmentTableName+".Id"]
				if eqExpr != nil {
					if strID, ok := eqExpr.(string); ok {
						upsertFulfillment.Id = strID
					}
				}
			}
		}

		if option.OrderID != nil {
			if eq, isEqual := option.OrderID.(squirrel.Eq); isEqual && eq != nil {
				if eqExpr := eq[store.FulfillmentTableName+".OrderID"]; eqExpr != nil {
					if strEq, ok := eqExpr.(string); ok {
						upsertFulfillment.OrderID = strEq
					}
				}
			}
		}

		if option.Status != nil {
			if eq, isEqual := option.Status.(squirrel.Eq); isEqual && eq != nil {
				if eqExpr := eq[store.FulfillmentTableName+".Status"]; eqExpr != nil {
					if strEq, ok := eqExpr.(string); ok {
						upsertFulfillment.Status = model.FulfillmentStatus(strEq)
					}
				}
			}
		}

		fulfillmentByOption, appErr = a.UpsertFulfillment(transaction, upsertFulfillment)
		if appErr != nil {
			return nil, appErr
		}

		return fulfillmentByOption, nil
	}

	return fulfillmentByOption, nil
}

// BulkDeleteFulfillments tells store to delete fulfillments that satisfy given option
func (a *ServiceOrder) BulkDeleteFulfillments(transaction store_iface.SqlxTxExecutor, fulfillments model.Fulfillments) *model.AppError {
	err := a.srv.Store.Fulfillment().BulkDeleteFulfillments(transaction, fulfillments)
	if err != nil {
		return model.NewAppError("BulkDeleteFulfillments", "app.order.error_deleting_fulfillments.app_error", nil, err.Error(), http.StatusInternalServerError)
	}

	return nil
}
