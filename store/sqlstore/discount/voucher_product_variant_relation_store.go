package discount

import (
	"github.com/sitename/sitename/model/product_and_discount"
	"github.com/sitename/sitename/store"
)

type SqlVoucherProductVariantStore struct {
	store.Store
}

func NewSqlVoucherProductVariantStore(s store.Store) store.VoucherProductVariantStore {
	v := &SqlVoucherProductVariantStore{s}
	for _, db := range s.GetAllConns() {
		table := db.AddTableWithName(product_and_discount.VoucherProductVariant{}, store.VoucherProductVariantTableName).SetKeys(false, "Id")
		table.ColMap("Id").SetMaxSize(store.UUID_MAX_LENGTH)
		table.ColMap("VoucherID").SetMaxSize(store.UUID_MAX_LENGTH)
		table.ColMap("ProductVariantID").SetMaxSize(store.UUID_MAX_LENGTH)

		table.SetUniqueTogether("VoucherID", "ProductVariantID")
	}
	return v
}

func (vs *SqlVoucherProductVariantStore) CreateIndexesIfNotExists() {}
