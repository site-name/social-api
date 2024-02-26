package account

import (
	"net/http"

	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model_helper"
	"github.com/sitename/sitename/modules/model_types"
	"github.com/sitename/sitename/store"
)

func (a *ServiceAccount) CommonCustomerCreateEvent(
	tx store.ContextRunner,
	userID *string,
	orderID *string,
	eventType model.CustomerEventType,
	params model_types.JSONString,
) (*model.CustomerEvent, *model_helper.AppError) {
	event := model.CustomerEvent{
		Type:       eventType,
		Parameters: params,
		OrderID:    model_types.NullString{String: orderID},
		UserID:     model_types.NullString{String: userID},
	}

	savedEvent, err := a.srv.Store.CustomerEvent().Upsert(tx, event)
	if err != nil {
		if appErr, ok := err.(*model_helper.AppError); ok {
			return nil, appErr
		}
		return nil, model_helper.NewAppError("CommonCustomerCreateEvent", "app.account.customer_event_save_error.app_error", nil, err.Error(), http.StatusInternalServerError)
	}

	return savedEvent, nil
}

// CustomerPlacedOrderEvent creates an customer event, if given user is not valid, it returns immediately.
func (s *ServiceAccount) CustomerPlacedOrderEvent(tx store.ContextRunner, user *model.User, order model.Order) (*model.CustomerEvent, *model_helper.AppError) {
	if user == nil || !model_helper.IsValidId(user.ID) {
		return nil, nil
	}

	return s.CommonCustomerCreateEvent(tx, &user.ID, &order.ID, model.CustomerEventTypePLACED_ORDER, nil)
}
