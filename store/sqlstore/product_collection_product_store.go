package sqlstore

import "github.com/sitename/sitename/store"

type SqlCollectionProductStore struct {
	*SqlStore
}

func newSqlCollectionProductStore(s *SqlStore) store.CollectionProductStore {
	cps := &SqlCollectionProductStore{s}

	return cps
}

func (ps *SqlCollectionProductStore) createIndexesIfNotExists() {

}
