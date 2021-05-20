package sqlstore

import "github.com/sitename/sitename/store"

type SqlCategoryStore struct {
	*SqlStore
}

func newSqlCategoryStore(s *SqlStore) store.CategoryStore {
	cs := &SqlCategoryStore{s}

	return cs
}

func (ps *SqlCategoryStore) createIndexesIfNotExists() {

}
