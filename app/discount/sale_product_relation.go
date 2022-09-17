package discount

import (
	"net/http"

	"github.com/sitename/sitename/model"
)

// SaleProductsByOptions returns a slice of sale-product relations filtered using given options
func (s *ServiceDiscount) SaleProductsByOptions(options *model.SaleProductRelationFilterOption) ([]*model.SaleProductRelation, *model.AppError) {
	saleProducts, err := s.srv.Store.SaleProductRelation().SaleProductsByOption(options)
	var (
		statusCode int
		errMessage string
	)
	if err != nil {
		statusCode = http.StatusInternalServerError
		errMessage = err.Error()
	} else if len(saleProducts) == 0 {
		statusCode = http.StatusNotFound
	}

	if statusCode != 0 {
		return nil, model.NewAppError("SaleProductsByOptions", "app.discount.error_finding_sale_product_relation_with_options.app_error", nil, errMessage, statusCode)
	}

	return saleProducts, nil
}
