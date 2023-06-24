package discount

import (
	"net/http"

	"github.com/sitename/sitename/model"
)

func (s *ServiceDiscount) FilterVats(options *model.VatFilterOptions) ([]*model.Vat, *model.AppError) {
	vats, err := s.srv.Store.Vat().FilterByOptions(options)
	if err != nil {
		return nil, model.NewAppError("FilterVats", "app.shop.vats_by_options.app_error", nil, err.Error(), http.StatusInternalServerError)
	}

	return vats, nil
}
