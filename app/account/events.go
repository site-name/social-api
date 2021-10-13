package account

import (
	"net/http"

	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model/account"
	"github.com/sitename/sitename/model/order"
)

func (a *ServiceAccount) CommonCustomerCreateEvent(userID *string, orderID *string, eventType string, params model.StringInterface) (*account.CustomerEvent, *model.AppError) {
	event := &account.CustomerEvent{
		Type:       eventType,
		Parameters: params,
		OrderID:    orderID,
		UserID:     userID,
	}

	event, err := a.srv.Store.CustomerEvent().Save(event)
	if err != nil {
		if appErr, ok := err.(*model.AppError); ok {
			return nil, appErr
		}
		return nil, model.NewAppError("CommonCustomerCreateEvent", "app.account.customer_event_save_error.app_error", nil, err.Error(), http.StatusInternalServerError)
	}

	return event, nil
}

// CustomerPlacedOrderEvent creates an customer event, if given user is not valid, it returns immediately.
func (s *ServiceAccount) CustomerPlacedOrderEvent(user *account.User, orDer *order.Order) (*account.CustomerEvent, *model.AppError) {
	if user == nil || !model.IsValidId(user.Id) {
		return nil, nil
	}

	return s.CommonCustomerCreateEvent(&user.Id, &orDer.Id, account.PLACED_ORDER, nil)
}
