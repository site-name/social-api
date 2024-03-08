package account

import (
	"net/http"

	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model_helper"
)

func (a *ServiceAccount) CustomerEventsByOptions(options model_helper.CustomerEventFilterOptions) (model.CustomerEventSlice, *model_helper.AppError) {
	events, err := a.srv.Store.CustomerEvent().FilterByOptions(options)
	if err != nil {
		return nil, model_helper.NewAppError("CustomerEventsByOptions", "app.customer_event.filter_by_opions.app_error", nil, err.Error(), http.StatusInternalServerError)
	}

	return events, nil
}
