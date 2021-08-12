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

// CollectionsByProductID finds and returns all collections related to given product
func (a *AppProduct) CollectionsByProductID(productID string) ([]*product_and_discount.Collection, *model.AppError) {
	collections, err := a.Srv().Store.Collection().CollectionsByProductID(productID)
	if err != nil {
		return nil, store.AppErrorFromDatabaseLookupError("CollectionsByProductID", "app.product.error_finding_collections_by_productID", err)
	}

	var (
		res     []*product_and_discount.Collection
		meetMap = map[string]bool{}
	)
	for _, collection := range collections {
		if _, met := meetMap[collection.Id]; !met {
			res = append(res, collection)
			meetMap[collection.Id] = true
		}
	}
	return res, nil
}
