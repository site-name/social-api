package sqlstore

import "github.com/sitename/sitename/store"

type SqlProductMediaStore struct {
	*SqlStore
}

func newSqlProductMediaStore(s *SqlStore) store.ProductMediaStore {
	pms := &SqlProductMediaStore{s}

	return pms
}

func (ps *SqlProductMediaStore) createIndexesIfNotExists() {

}
