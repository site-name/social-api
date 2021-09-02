package product

import (
	"net/http"

	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model/product_and_discount"
)

// ProductMediasByOption returns a list of product medias that satisfy given option
func (a *ServiceProduct) ProductMediasByOption(option *product_and_discount.ProductMediaFilterOption) ([]*product_and_discount.ProductMedia, *model.AppError) {
	productMedias, err := a.srv.Store.ProductMedia().FilterByOption(option)
	var (
		errMsg     string
		statusCode int
	)
	if err != nil {
		errMsg = err.Error()
		statusCode = http.StatusInternalServerError
	}
	if len(productMedias) == 0 {
		statusCode = http.StatusNotFound
	}

	if statusCode != 0 {
		return nil, model.NewAppError("ProductMediasByOption", "app.product.error_finding_product_medias_by_option.app_error", nil, errMsg, statusCode)
	}

	return productMedias, nil
}
