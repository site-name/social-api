package sqlstore

import "github.com/sitename/sitename/store"

type SqlDigitalContentStore struct {
	*SqlStore
}

func newSqlDigitalContentStore(s *SqlStore) store.DigitalContentStore {
	dcs := &SqlDigitalContentStore{s}

	return dcs
}

func (ps *SqlDigitalContentStore) createIndexesIfNotExists() {

}
