package account

import (
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model/account"
	"github.com/sitename/sitename/store"
)

// CustomerEventsByUser returns customer events belong to given user
func (a *ServiceAccount) CustomerEventsByUser(userID string) ([]*account.CustomerEvent, *model.AppError) {
	events, err := a.srv.Store.CustomerEvent().GetEventsByUserID(userID)
	if err != nil {
		return nil, store.AppErrorFromDatabaseLookupError("CustomerEventsByUser", "app.account.customer_event_missing", err)
	}

	return events, nil
}
