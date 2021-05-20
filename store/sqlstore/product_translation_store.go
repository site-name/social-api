package sqlstore

import "github.com/sitename/sitename/store"

type SqlProductTranslationStore struct {
	*SqlStore
}

func newSqlProductTranslationStore(s *SqlStore) store.ProductTranslationStore {
	pts := &SqlProductTranslationStore{s}

	return pts
}

func (ps *SqlProductTranslationStore) createIndexesIfNotExists() {

}
