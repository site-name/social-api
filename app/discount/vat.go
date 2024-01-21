package discount

import (
	"net/http"

	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model_helper"
)

func (s *ServiceDiscount) FilterVats(options *model.VatFilterOptions) ([]*model.Vat, *model_helper.AppError) {
	vats, err := s.srv.Store.Vat().FilterByOptions(options)
	if err != nil {
		return nil, model_helper.NewAppError("FilterVats", "app.shop.vats_by_options.app_error", nil, err.Error(), http.StatusInternalServerError)
	}

	return vats, nil
}
