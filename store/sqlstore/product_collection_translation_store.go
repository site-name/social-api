package sqlstore

import "github.com/sitename/sitename/store"

type SqlCollectionTranslationStore struct {
	*SqlStore
}

func newSqlCollectionTranslationStore(s *SqlStore) store.CollectionTranslationStore {
	cts := &SqlCollectionTranslationStore{s}

	return cts
}

func (ps *SqlCollectionTranslationStore) createIndexesIfNotExists() {

}
