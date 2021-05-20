package sqlstore

import "github.com/sitename/sitename/store"

type SqlVariantMediaStore struct {
	*SqlStore
}

func newSqlVariantMediaStore(s *SqlStore) store.VariantMediaStore {
	vms := &SqlVariantMediaStore{s}

	return vms
}

func (ps *SqlVariantMediaStore) createIndexesIfNotExists() {

}
