package sqlstore

import (
	"github.com/sitename/sitename/model/product_and_discount"
	"github.com/sitename/sitename/store"
)

type SqlVariantMediaStore struct {
	*SqlStore
}

func newSqlVariantMediaStore(s *SqlStore) store.VariantMediaStore {
	vms := &SqlVariantMediaStore{s}

	for _, db := range s.GetAllConns() {
		table := db.AddTableWithName(product_and_discount.VariantMedia{}, "VariantMedias").SetKeys(false, "Id")
		table.ColMap("Id").SetMaxSize(UUID_MAX_LENGTH)
		table.ColMap("VariantID").SetMaxSize(UUID_MAX_LENGTH)
		table.ColMap("MediaID").SetMaxSize(UUID_MAX_LENGTH)
	}
	return vms
}

func (ps *SqlVariantMediaStore) createIndexesIfNotExists() {

}
