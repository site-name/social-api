package discount

import (
	"net/http"

	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model/product_and_discount"
)

// SaleProductVariantsByOptions returns a list of sale-product variant relations filtered using given options
func (s *ServiceDiscount) SaleProductVariantsByOptions(options *product_and_discount.SaleProductVariantFilterOption) ([]*product_and_discount.SaleProductVariant, *model.AppError) {
	saleProductVariants, err := s.srv.Store.SaleProductVariant().FilterByOption(options)
	var (
		statusCode int
		errMessage string
	)
	if err != nil {
		statusCode = http.StatusInternalServerError
		errMessage = err.Error()
	} else if len(saleProductVariants) == 0 {
		statusCode = http.StatusNotFound
	}

	if statusCode != 0 {
		return nil, model.NewAppError("SaleProductVariantsByOptions", "app.discount.error_finding_sale_product_variants_by_options.app_error", nil, errMessage, statusCode)
	}

	return saleProductVariants, nil
}