package sqlstore

import "github.com/sitename/sitename/store"

type SqlDigitalContentUrlStore struct {
	*SqlStore
}

func newSqlDigitalContentUrlStore(s *SqlStore) store.DigitalContentUrlStore {
	dcs := &SqlDigitalContentUrlStore{s}

	return dcs
}

func (ps *SqlDigitalContentUrlStore) createIndexesIfNotExists() {

}
