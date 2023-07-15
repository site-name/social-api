package discount

import (
	"net/http"

	"github.com/Masterminds/squirrel"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/store"
)

// SaleProductVariantsByOptions returns a list of sale-product variant relations filtered using given options
func (s *ServiceDiscount) SaleProductVariantsByOptions(options squirrel.Sqlizer) ([]*model.SaleProductVariant, *model.AppError) {
	var res []*model.SaleProductVariant
	err := s.srv.Store.GetReplica().Table("sale_productvariants").Find(&res, store.BuildSqlizer(options)...).Error
	if err != nil {
		return nil, model.NewAppError("SaleProductVariantsByOptions", "app.discount.error_finding_sale_product_variants_by_options.app_error", nil, err.Error(), http.StatusInternalServerError)
	}

	return res, nil
}
