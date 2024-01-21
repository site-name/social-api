package account

import (
	"net/http"

	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model_helper"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

func (a *ServiceAccount) CustomerEventsByOptions(conds ...qm.QueryMod) (model.CustomerEventSlice, *model_helper.AppError) {
	events, err := a.srv.Store.CustomerEvent().FilterByOptions(conds...)
	if err != nil {
		return nil, model_helper.NewAppError("CustomerEventsByOptions", "app.customer_event.filter_by_opions.app_error", nil, err.Error(), http.StatusInternalServerError)
	}

	return events, nil
}
