package shipping

import (
	"github.com/sitename/sitename/model/shipping"
	"github.com/sitename/sitename/store"
)

type SqlShippingMethodExcludedProductStore struct {
	store.Store
}

func NewSqlShippingMethodExcludedProductStore(s store.Store) store.ShippingMethodExcludedProductStore {
	ss := &SqlShippingMethodExcludedProductStore{s}
	for _, db := range s.GetAllConns() {
		table := db.AddTableWithName(shipping.ShippingMethodExcludedProduct{}, store.ShippingMethodExcludedProductTableName).SetKeys(false, "Id")
		table.ColMap("Id").SetMaxSize(store.UUID_MAX_LENGTH)
		table.ColMap("ShippingMethodID").SetMaxSize(store.UUID_MAX_LENGTH)
		table.ColMap("ProductID").SetMaxSize(store.UUID_MAX_LENGTH)

		table.SetUniqueTogether("ShippingMethodID", "ProductID")
	}

	return ss
}

func (ss *SqlShippingMethodExcludedProductStore) CreateIndexesIfNotExists() {
	ss.CreateForeignKeyIfNotExists(store.ShippingMethodExcludedProductTableName, "ShippingMethodID", store.ShippingMethodTableName, "Id", false)
	ss.CreateForeignKeyIfNotExists(store.ShippingMethodExcludedProductTableName, "ProductID", store.ProductTableName, "Id", false)
}
