/*
	NOTE: This package is initialized during server startup (modules/imports does that)
	so the init() function get the chance to register a function to create `ServiceAccount`
*/
package product

import (
	"net/http"
	"sync"

	"github.com/sitename/sitename/app"
	"github.com/sitename/sitename/app/sub_app_iface"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model/product_and_discount"
	"github.com/sitename/sitename/store"
)

type ServiceProduct struct {
	srv *app.Server
	sync.WaitGroup
	sync.Mutex
}

func init() {
	app.RegisterProductService(func(s *app.Server) (sub_app_iface.ProductService, error) {
		return &ServiceProduct{
			srv: s,
		}, nil
	})
}

// ProductById returns 1 product by given id
func (a *ServiceProduct) ProductById(productID string) (*product_and_discount.Product, *model.AppError) {
	return a.ProductByOption(&product_and_discount.ProductFilterOption{
		Id: &model.StringFilter{
			StringOption: &model.StringOption{
				Eq: productID,
			},
		},
	})
}

// ProductsByOption returns a list of products that satisfy given option
func (a *ServiceProduct) ProductsByOption(option *product_and_discount.ProductFilterOption) ([]*product_and_discount.Product, *model.AppError) {
	products, err := a.srv.Store.Product().FilterByOption(option)
	var (
		statusCode int
		errMsg     string
	)
	if err != nil {
		statusCode = http.StatusInternalServerError
		errMsg = err.Error()
	}
	if len(products) == 0 {
		statusCode = http.StatusNotFound
	}
	if statusCode != 0 {
		return nil, model.NewAppError("ProductsByOption", "app.product.error_finding_products_by_option.app_error", nil, errMsg, statusCode)
	}

	return products, nil
}

// ProductByOption returns 1 product that satisfy given option
func (a *ServiceProduct) ProductByOption(option *product_and_discount.ProductFilterOption) (*product_and_discount.Product, *model.AppError) {
	product, err := a.srv.Store.Product().GetByOption(option)
	if err != nil {
		return nil, store.AppErrorFromDatabaseLookupError("ProductByOption", "app.error_finding_product_by_option.app_error", err)
	}

	return product, nil
}

// ProductsByVoucherID finds all products that have relationships with given voucher
func (a *ServiceProduct) ProductsByVoucherID(voucherID string) ([]*product_and_discount.Product, *model.AppError) {
	products, appErr := a.ProductsByOption(&product_and_discount.ProductFilterOption{
		VoucherIDs: []string{voucherID},
	})
	if appErr != nil {
		return nil, appErr
	}

	return products, nil
}

// ProductsRequireShipping checks if at least 1 product require shipping, then return true, false otherwise
func (a *ServiceProduct) ProductsRequireShipping(productIDs []string) (bool, *model.AppError) {
	productTypes, appErr := a.ProductTypesByProductIDs(productIDs)
	if appErr != nil { // this error caused by system
		return false, appErr
	}

	for _, productType := range productTypes {
		if *productType.IsShippingRequired { // use pointer directly since this field has default value
			return true, nil
		}
	}

	return false, nil
}

// ProductGetFirstImage returns first media of given product
func (a *ServiceProduct) ProductGetFirstImage(productID string) (*product_and_discount.ProductMedia, *model.AppError) {
	mediasOfProduct, appErr := a.ProductMediasByOption(&product_and_discount.ProductMediaFilterOption{
		ProductID: &model.StringFilter{
			StringOption: &model.StringOption{
				Eq: productID,
			},
		},
		Type: []string{product_and_discount.IMAGE},
	})
	if appErr != nil {
		return nil, appErr
	}

	return mediasOfProduct[0], nil
}
