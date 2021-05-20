package sqlstore

import "github.com/sitename/sitename/store"

type SqlProductTypeStore struct {
	*SqlStore
}

func newSqlProductTypeStore(s *SqlStore) store.ProductTypeStore {
	pts := &SqlProductTypeStore{s}

	return pts
}

func (ps *SqlProductTypeStore) createIndexesIfNotExists() {

}
