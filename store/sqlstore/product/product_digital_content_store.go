package product

import (
	"github.com/sitename/sitename/model/product_and_discount"
	"github.com/sitename/sitename/store"
)

type SqlDigitalContentStore struct {
	store.Store
}

func NewSqlDigitalContentStore(s store.Store) store.DigitalContentStore {
	dcs := &SqlDigitalContentStore{s}

	for _, db := range s.GetAllConns() {
		table := db.AddTableWithName(product_and_discount.DigitalContent{}, "DigitalContents").SetKeys(false, "Id")
		table.ColMap("Id").SetMaxSize(store.UUID_MAX_LENGTH)
		table.ColMap("ProductVariantID").SetMaxSize(store.UUID_MAX_LENGTH)
		table.ColMap("ContentType").SetMaxSize(product_and_discount.DIGITAL_CONTENT_CONTENT_TYPE_MAX_LENGTH)
	}
	return dcs
}

func (ps *SqlDigitalContentStore) CreateIndexesIfNotExists() {

}
