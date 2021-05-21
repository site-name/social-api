package sqlstore

import (
	"github.com/sitename/sitename/model/product_and_discount"
	"github.com/sitename/sitename/store"
)

type SqlCollectionProductStore struct {
	*SqlStore
}

func newSqlCollectionProductStore(s *SqlStore) store.CollectionProductStore {
	cps := &SqlCollectionProductStore{s}

	for _, db := range s.GetAllConns() {
		table := db.AddTableWithName(product_and_discount.CollectionProduct{}, "CollectionProducts").SetKeys(false, "Id")
		table.ColMap("Id").SetMaxSize(UUID_MAX_LENGTH)
		table.ColMap("CollectionID").SetMaxSize(UUID_MAX_LENGTH)
		table.ColMap("ProductID").SetMaxSize(UUID_MAX_LENGTH)

		table.SetUniqueTogether("CollectionID", "ProductID")
	}
	return cps
}

func (ps *SqlCollectionProductStore) createIndexesIfNotExists() {

}
