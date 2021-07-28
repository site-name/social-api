package product

import (
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model/product_and_discount"
	"github.com/sitename/sitename/store"
)

// CollectionsByVoucherID finds all collections that have relationships with given voucher
func (a *AppProduct) CollectionsByVoucherID(voucherID string) ([]*product_and_discount.Collection, *model.AppError) {
	collections, err := a.Srv().Store.VoucherCollection().CollectionsByVoucherID(voucherID)
	if err != nil {
		return nil, store.AppErrorFromDatabaseLookupError("App.Product.CollectionsByVoucherID", "app.product.error_finding_collections_by_voucherID.app_error", err)
	}

	return collections, nil
}
