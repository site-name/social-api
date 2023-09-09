package product

import (
	"net/http"

	"github.com/sitename/sitename/model"
)

// ProductMediasByOption returns a list of product medias that satisfy given option
func (a *ServiceProduct) ProductMediasByOption(option *model.ProductMediaFilterOption) ([]*model.ProductMedia, *model.AppError) {
	productMedias, err := a.srv.Store.ProductMedia().FilterByOption(option)
	if err != nil {
		return nil, model.NewAppError("ProductMediasByOption", "app.product.error_finding_product_medias_by_option.app_error", nil, err.Error(), http.StatusInternalServerError)
	}

	return productMedias, nil
}
