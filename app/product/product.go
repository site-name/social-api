package product

import (
	"net/http"

	"github.com/sitename/sitename/app"
	"github.com/sitename/sitename/app/sub_app_iface"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model/product_and_discount"
	"github.com/sitename/sitename/store"
)

type AppProduct struct {
	app.AppIface
}

func init() {
	app.RegisterProductApp(func(a app.AppIface) sub_app_iface.ProductApp {
		return &AppProduct{a}
	})
}

func (a *AppProduct) ProductById(productID string) (*product_and_discount.Product, *model.AppError) {
	product, err := a.Srv().Store.Product().Get(productID)
	if err != nil {
		return nil, store.AppErrorFromDatabaseLookupError("ProductById", "app.product.product_missing.app_error", err)
	}

	return product, nil
}

// ProductsByVoucherID finds all products that have relationships with given voucher
func (a *AppProduct) ProductsByVoucherID(voucherID string) ([]*product_and_discount.Product, *model.AppError) {
	products, err := a.Srv().Store.VoucherProduct().ProductsByVoucherID(voucherID)
	if err != nil {
		return nil, store.AppErrorFromDatabaseLookupError("App.Product.ProductsByVoucherID", "app.error_finding_products_by_voucherID.app_error", err)
	}

	return products, nil
}

// ProductsRequireShipping checks if at least 1 product require shipping, then return true, false otherwise
func (a *AppProduct) ProductsRequireShipping(productIDs []string) (bool, *model.AppError) {
	productTypes, appErr := a.ProductTypesByProductIDs(productIDs)
	if appErr != nil {
		if appErr.StatusCode == http.StatusNotFound {
			return false, nil
		}
		return false, appErr
	}

	for _, productType := range productTypes {
		if *productType.IsShippingRequired { // use pointer directly since this field has default value
			return true, nil
		}
	}

	return false, nil
}
