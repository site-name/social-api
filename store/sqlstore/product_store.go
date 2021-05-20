package sqlstore

import "github.com/sitename/sitename/store"

type SqlProductStore struct {
	*SqlStore
}

func newSqlProductStore(s *SqlStore) store.ProductStore {
	ps := &SqlProductStore{s}

	return ps
}

func (ps *SqlProductStore) createIndexesIfNotExists() {

}
