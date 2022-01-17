package discount

import (
	"net/http"

	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model/product_and_discount"
)

// SaleProductsByOptions returns a slice of sale-product relations filtered using given options
func (s *ServiceDiscount) SaleProductsByOptions(options *product_and_discount.SaleProductRelationFilterOption) ([]*product_and_discount.SaleProductRelation, *model.AppError) {
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
