package product

import (
	"net/http"

	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model_helper"
	"github.com/sitename/sitename/store"
	"gorm.io/gorm"
)

// CollectionProductRelationsByOptions finds and returns a list of product-collection relations based on given filter options
func (s *ServiceProduct) CollectionProductRelationsByOptions(options *model.CollectionProductFilterOptions) ([]*model.CollectionProduct, *model_helper.AppError) {
	relations, err := s.srv.Store.CollectionProduct().FilterByOptions(options)
	if err != nil {
		return nil, model_helper.NewAppError("CollectionProductRelationsByOptions", "app.product.error_finding_product-collection_relations.app_error", nil, err.Error(), http.StatusInternalServerError)
	}

	return relations, nil
}

func (s *ServiceProduct) CreateCollectionProductRelations(transaction *gorm.DB, relations []*model.CollectionProduct) ([]*model.CollectionProduct, *model_helper.AppError) {
	relations, err := s.srv.Store.CollectionProduct().BulkSave(transaction, relations)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if _, ok := err.(*store.ErrInvalidInput); ok {
			statusCode = http.StatusBadRequest
		}

		return nil, model_helper.NewAppError("CreateCollectionProductRelations", "app.product.error_saving_collection_product.app_error", nil, err.Error(), statusCode)
	}

	return relations, nil
}
