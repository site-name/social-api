package product

import (
	"net/http"

	"github.com/sitename/sitename/model"
)

// CollectionProductRelationsByOptions finds and returns a list of product-collection relations based on given filter options
func (s *ServiceProduct) CollectionProductRelationsByOptions(options *model.CollectionProductFilterOptions) ([]*model.CollectionProduct, *model.AppError) {
	relations, err := s.srv.Store.CollectionProduct().FilterByOptions(options)
	if err != nil {
		return nil, model.NewAppError("CollectionProductRelationsByOptions", "app.product.error_finding_product-collection_relations.app_error", nil, err.Error(), http.StatusInternalServerError)
	}

	return relations, nil
}
