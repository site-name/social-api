package sqlstore

import "github.com/sitename/sitename/store"

type SqlCollectionStore struct {
	*SqlStore
}

func newSqlCollectionStore(s *SqlStore) store.CollectionStore {
	cs := &SqlCollectionStore{s}

	return cs
}

func (ps *SqlCollectionStore) createIndexesIfNotExists() {

}
