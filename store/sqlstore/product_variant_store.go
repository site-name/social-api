package sqlstore

import "github.com/sitename/sitename/store"

type SqlProductVariantStore struct {
	*SqlStore
}

func newSqlProductVariantStore(s *SqlStore) store.ProductVariantStore {
	pvs := &SqlProductVariantStore{s}

	return pvs
}

func (ps *SqlProductVariantStore) createIndexesIfNotExists() {

}
