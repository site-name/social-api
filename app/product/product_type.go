package product

import (
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model/product_and_discount"
	"github.com/sitename/sitename/store"
)

func (a *AppProduct) ProductTypesByCheckoutToken(checkoutToken string) ([]*product_and_discount.ProductType, *model.AppError) {
	types, err := a.Srv().Store.ProductType().FilterProductTypesByCheckoutID(checkoutToken)
	if err != nil {
		return nil, store.AppErrorFromDatabaseLookupError("ProductTypesByCheckoutToken", "app.product.product_types_by_checkout_missing.app_error", err)
	}

	return types, nil
}
