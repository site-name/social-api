package discount

import (
	"net/http"

	"github.com/sitename/sitename/model"
)

// SaleProductVariantsByOptions returns a list of sale-product variant relations filtered using given options
func (s *ServiceDiscount) SaleProductVariantsByOptions(options *model.SaleProductVariantFilterOption) ([]*model.SaleProductVariant, *model.AppError) {
	saleProductVariants, err := s.srv.Store.SaleProductVariant().FilterByOption(options)
	if err != nil {
		return nil, model.NewAppError("SaleProductVariantsByOptions", "app.discount.error_finding_sale_product_variants_by_options.app_error", nil, err.Error(), http.StatusInternalServerError)
	}

	return saleProductVariants, nil
}
