package discount

import (
	"net/http"

	"github.com/Masterminds/squirrel"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/store"
)

// SaleProductsByOptions returns a slice of sale-product relations filtered using given options
func (s *ServiceDiscount) SaleProductsByOptions(options squirrel.Sqlizer) ([]*model.SaleProduct, *model.AppError) {
	var res []*model.SaleProduct
	err := s.srv.Store.GetReplica().Table("sale_collections").Find(&res, store.BuildSqlizer(options)...).Error
	if err != nil {
		return nil, model.NewAppError("SaleProductsByOptions", "app.discount.sale_product_relations.app_error", nil, err.Error(), http.StatusInternalServerError)
	}

	return res, nil
}
