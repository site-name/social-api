package discount

import (
	"net/http"

	"github.com/Masterminds/squirrel"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/store"
)

// SaleCollectionsByOptions returns a slice of sale-collection relations filtered using given options
func (s *ServiceDiscount) SaleCollectionsByOptions(options squirrel.Sqlizer) ([]*model.SaleCollection, *model.AppError) {
	var res []*model.SaleCollection
	err := s.srv.Store.GetReplica().Table("sale_collections").Find(&res, store.BuildSqlizer(options)...).Error
	if err != nil {
		return nil, model.NewAppError("SaleCollectionsByOptions", "app.discount.sale_collections_by_options.app_error", nil, err.Error(), http.StatusInternalServerError)
	}

	return res, nil
}
