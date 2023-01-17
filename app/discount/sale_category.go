package discount

import (
	"net/http"

	"github.com/sitename/sitename/model"
)

// SaleCategoriesByOption returns sale-category relations with an app error
func (s *ServiceDiscount) SaleCategoriesByOption(option *model.SaleCategoryRelationFilterOption) ([]*model.SaleCategoryRelation, *model.AppError) {
	saleCategoryRelations, err := s.srv.Store.SaleCategoryRelation().SaleCategoriesByOption(option)
	if err != nil {
		return nil, model.NewAppError("SaleCategoriesByOption", "app.discount.sale_categories_by_options.app_error", nil, err.Error(), http.StatusInternalServerError)
	}

	return saleCategoryRelations, nil
}
