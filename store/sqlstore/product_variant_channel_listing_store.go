package sqlstore

import "github.com/sitename/sitename/store"

type SqlProductVariantChannelListingStore struct {
	*SqlStore
}

func newSqlProductVariantChannelListingStore(s *SqlStore) store.ProductVariantChannelListingStore {
	pvcls := &SqlProductVariantChannelListingStore{s}

	return pvcls
}

func (ps *SqlProductVariantChannelListingStore) createIndexesIfNotExists() {

}
