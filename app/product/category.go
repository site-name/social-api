package product

import (
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model/product_and_discount"
	"github.com/sitename/sitename/store"
)

// CategoriesByVoucherID finds all categories that have relationship with given voucher
func (a *AppProduct) CategoriesByVoucherID(voucherID string) ([]*product_and_discount.Category, *model.AppError) {
	categories, err := a.Srv().Store.VoucherCategory().ProductCategoriesByVoucherID(voucherID)
	if err != nil {
		return nil, store.AppErrorFromDatabaseLookupError("App.Product.CategoriesByVoucherID", "app.product.error_finding_categories_by_voucherID.app_error", err)
	}

	return categories, nil
}