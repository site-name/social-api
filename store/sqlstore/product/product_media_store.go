package product

import (
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model/product_and_discount"
	"github.com/sitename/sitename/store"
)

type SqlProductMediaStore struct {
	store.Store
}

func NewSqlProductMediaStore(s store.Store) store.ProductMediaStore {
	pms := &SqlProductMediaStore{s}

	for _, db := range s.GetAllConns() {
		table := db.AddTableWithName(product_and_discount.ProductMedia{}, store.ProductMediaTableName).SetKeys(false, "Id")
		table.ColMap("Id").SetMaxSize(store.UUID_MAX_LENGTH)
		table.ColMap("ProductID").SetMaxSize(store.UUID_MAX_LENGTH)
		table.ColMap("Image").SetMaxSize(model.URL_LINK_MAX_LENGTH)
		table.ColMap("Ppoi").SetMaxSize(product_and_discount.PRODUCT_MEDIA_PPOI_MAX_LENGTH)
		table.ColMap("Type").SetMaxSize(product_and_discount.PRODUCT_MEDIA_TYPE_MAX_LENGTH)
		table.ColMap("ExternalUrl").SetMaxSize(product_and_discount.PRODUCT_MEDIA_EXTERNAL_URL_MAX_LENGTH)
		table.ColMap("Alt").SetMaxSize(product_and_discount.PRODUCT_MEDIA_ALT_MAX_LENGTH)
	}
	return pms
}

func (ps *SqlProductMediaStore) CreateIndexesIfNotExists() {

}
