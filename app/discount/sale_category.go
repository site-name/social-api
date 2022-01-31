package discount

import (
	"net/http"

	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model/product_and_discount"
)

// SaleCategoriesByOption returns sale-category relations with an app error
func (s *ServiceDiscount) SaleCategoriesByOption(option *product_and_discount.SaleCategoryRelationFilterOption) ([]*product_and_discount.SaleCategoryRelation, *model.AppError) {
	saleCategoryRelations, err := s.srv.Store.SaleCategoryRelation().SaleCategoriesByOption(option)
	var (
		statusCode int
		errMessage string
	)
	if err != nil {
		statusCode = http.StatusInternalServerError
		errMessage = err.Error()
	}
	if len(saleCategoryRelations) == 0 {
		statusCode = http.StatusNotFound
	}

	if statusCode != 0 {
		return nil, model.NewAppError("SaleCategoriesByOption", "app.discount.error_finding_sale_category_relations.app_error", nil, errMessage, statusCode)
	}

	return saleCategoryRelations, nil
}
