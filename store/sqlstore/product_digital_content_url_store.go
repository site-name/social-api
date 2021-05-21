package sqlstore

import (
	"github.com/sitename/sitename/model/product_and_discount"
	"github.com/sitename/sitename/store"
)

type SqlDigitalContentUrlStore struct {
	*SqlStore
}

func newSqlDigitalContentUrlStore(s *SqlStore) store.DigitalContentUrlStore {
	dcs := &SqlDigitalContentUrlStore{s}

	for _, db := range s.GetAllConns() {
		table := db.AddTableWithName(product_and_discount.DigitalContentUrl{}, "DigitalContentUrls").SetKeys(false, "Id")
		table.ColMap("Id").SetMaxSize(UUID_MAX_LENGTH)
		table.ColMap("Token").SetMaxSize(UUID_MAX_LENGTH).SetUnique(true)
		table.ColMap("ContentID").SetMaxSize(UUID_MAX_LENGTH)
		table.ColMap("LineID").SetMaxSize(UUID_MAX_LENGTH)
	}
	return dcs
}

func (ps *SqlDigitalContentUrlStore) createIndexesIfNotExists() {

}
