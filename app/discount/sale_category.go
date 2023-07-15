package discount

import (
	"net/http"

	"github.com/Masterminds/squirrel"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/store"
)

// SaleCategoriesByOption returns sale-category relations with an app error
func (s *ServiceDiscount) SaleCategoriesByOption(option squirrel.Sqlizer) ([]*model.SaleCategory, *model.AppError) {

	var res []*model.SaleCategory
	err := s.srv.Store.GetReplica().Table("sale_categories").Find(&res, store.BuildSqlizer(option)...).Error
	if err != nil {
		return nil, model.NewAppError("SaleCategoriesByOption", "app.discount.sale_categories_by_options.app_error", nil, err.Error(), http.StatusInternalServerError)
	}

	return res, nil
}
