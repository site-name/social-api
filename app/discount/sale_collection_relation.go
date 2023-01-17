package discount

import (
	"net/http"

	"github.com/sitename/sitename/model"
)

// SaleCollectionsByOptions returns a slice of sale-collection relations filtered using given options
func (s *ServiceDiscount) SaleCollectionsByOptions(options *model.SaleCollectionRelationFilterOption) ([]*model.SaleCollectionRelation, *model.AppError) {
	saleCollections, err := s.srv.Store.SaleCollectionRelation().FilterByOption(options)
	if err != nil {
		return nil, model.NewAppError("SaleCollectionsByOptions", "app.discount.error_finding_sale_collections_by_options.app_error", nil, err.Error(), http.StatusInternalServerError)
	}

	return saleCollections, nil
}
