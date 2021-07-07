package account

import (
	"net/http"

	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model/account"
)

func (a *AppAccount) CommonCustomerCreateEvent(userID *string, orderID *string, eventType string, params model.StringInterface) (*account.CustomerEvent, *model.AppError) {
	event := &account.CustomerEvent{
		Type:       eventType,
		Parameters: params,
		OrderID:    orderID,
		UserID:     userID,
	}

	event, err := a.Srv().Store.CustomerEvent().Save(event)
	if err != nil {
		if appErr, ok := err.(*model.AppError); ok {
			return nil, appErr
		}
		return nil, model.NewAppError("CommonCustomerCreateEvent", "app.account.customer_event_save_error.app_error", nil, err.Error(), http.StatusInternalServerError)
	}

	return event, nil
}
