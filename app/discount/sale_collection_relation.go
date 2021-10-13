package discount

import (
	"net/http"

	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model/product_and_discount"
)

// SaleCollectionsByOptions returns a slice of sale-collection relations filtered using given options
func (s *ServiceDiscount) SaleCollectionsByOptions(options *product_and_discount.SaleCollectionRelationFilterOption) ([]*product_and_discount.SaleCollectionRelation, *model.AppError) {
	saleCollections, err := s.srv.Store.SaleCollectionRelation().FilterByOption(options)
	var (
		statusCode int
		errMessage string
	)
	if err != nil {
		statusCode = http.StatusInternalServerError
		errMessage = err.Error()
	} else if len(saleCollections) == 0 {
		statusCode = http.StatusNotFound
	}

	if statusCode != 0 {
		return nil, model.NewAppError("SaleCollectionsByOptions", "app.discount.error_finding_sale_collections_by_options.app_error", nil, errMessage, statusCode)
	}

	return saleCollections, nil
}
