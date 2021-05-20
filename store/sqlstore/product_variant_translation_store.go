package sqlstore

import "github.com/sitename/sitename/store"

type SqlProductVariantTranslationStore struct {
	*SqlStore
}

func newSqlProductVariantTranslationStore(s *SqlStore) store.ProductVariantTranslationStore {
	pvts := &SqlProductVariantTranslationStore{s}

	return pvts
}

func (ps *SqlProductVariantTranslationStore) createIndexesIfNotExists() {

}
