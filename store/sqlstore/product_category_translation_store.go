package sqlstore

import "github.com/sitename/sitename/store"

type SqlCategoryTranslationStore struct {
	*SqlStore
}

func newSqlCategoryTranslationStore(s *SqlStore) store.CategoryTranslationStore {
	cts := &SqlCategoryTranslationStore{s}

	return cts
}

func (ps *SqlCategoryTranslationStore) createIndexesIfNotExists() {

}
