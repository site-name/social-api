package product

import (
	"github.com/sitename/sitename/model/product_and_discount"
	"github.com/sitename/sitename/store"
)

type SqlDigitalContentUrlStore struct {
	store.Store
}

func NewSqlDigitalContentUrlStore(s store.Store) store.DigitalContentUrlStore {
	dcs := &SqlDigitalContentUrlStore{s}

	for _, db := range s.GetAllConns() {
		table := db.AddTableWithName(product_and_discount.DigitalContentUrl{}, "DigitalContentUrls").SetKeys(false, "Id")
		table.ColMap("Id").SetMaxSize(store.UUID_MAX_LENGTH)
		table.ColMap("Token").SetMaxSize(store.UUID_MAX_LENGTH).SetUnique(true)
		table.ColMap("ContentID").SetMaxSize(store.UUID_MAX_LENGTH)
		table.ColMap("LineID").SetMaxSize(store.UUID_MAX_LENGTH)
	}
	return dcs
}

func (ps *SqlDigitalContentUrlStore) CreateIndexesIfNotExists() {

}
