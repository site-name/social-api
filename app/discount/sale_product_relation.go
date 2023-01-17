package discount

import (
	"net/http"

	"github.com/sitename/sitename/model"
)

// SaleProductsByOptions returns a slice of sale-product relations filtered using given options
func (s *ServiceDiscount) SaleProductsByOptions(options *model.SaleProductRelationFilterOption) ([]*model.SaleProductRelation, *model.AppError) {
	saleProducts, err := s.srv.Store.SaleProductRelation().SaleProductsByOption(options)
	if err != nil {
		return nil, model.NewAppError("SaleProductsByOptions", "app.discount.error_finding_sale_product_relation_with_options.app_error", nil, err.Error(), http.StatusInternalServerError)
	}

	return saleProducts, nil
}
