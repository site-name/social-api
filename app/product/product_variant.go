package product

import (
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model/product_and_discount"
	"github.com/sitename/sitename/store"
)

func (a *AppProduct) ProductVariantById(id string) (*product_and_discount.ProductVariant, *model.AppError) {
	variant, err := a.Srv().Store.ProductVariant().Get(id)
	if err != nil {
		return nil, store.AppErrorFromDatabaseLookupError("ProductVariantbyId", "app.product.product_variant_missing.app_error", err)
	}

	return variant, nil
}
