package product

import (
	"github.com/sitename/sitename/model/product_and_discount"
	"github.com/sitename/sitename/store"
)

type SqlVariantMediaStore struct {
	store.Store
}

func NewSqlVariantMediaStore(s store.Store) store.VariantMediaStore {
	vms := &SqlVariantMediaStore{s}

	for _, db := range s.GetAllConns() {
		table := db.AddTableWithName(product_and_discount.VariantMedia{}, store.ProductVariantMediaTableName).SetKeys(false, "Id")
		table.ColMap("Id").SetMaxSize(store.UUID_MAX_LENGTH)
		table.ColMap("VariantID").SetMaxSize(store.UUID_MAX_LENGTH)
		table.ColMap("MediaID").SetMaxSize(store.UUID_MAX_LENGTH)

		table.SetUniqueTogether("VariantID", "MediaID")
	}
	return vms
}

func (ps *SqlVariantMediaStore) CreateIndexesIfNotExists() {
	ps.CreateForeignKeyIfNotExists(store.ProductVariantMediaTableName, "VariantID", store.ProductVariantTableName, "Id", true)
	ps.CreateForeignKeyIfNotExists(store.ProductVariantMediaTableName, "MediaID", store.ProductMediaTableName, "Id", true)
}
