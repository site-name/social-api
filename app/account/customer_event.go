package account

import (
	"net/http"

	"github.com/sitename/sitename/model"
)

func (a *ServiceAccount) CustomerEventsByOptions(option *model.CustomerEventFilterOptions) ([]*model.CustomerEvent, *model.AppError) {
	events, err := a.srv.Store.CustomerEvent().FilterByOptions(option)
	if err != nil {
		return nil, model.NewAppError("CustomerEventsByOptions", "app.customer_event.filter_by_opions.app_error", nil, err.Error(), http.StatusInternalServerError)
	}

	return events, nil
}
