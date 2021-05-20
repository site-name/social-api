package sqlstore

import "github.com/sitename/sitename/store"

type SqlProductChannelListingStore struct {
	*SqlStore
}

func newSqlProductChannelListingStore(s *SqlStore) store.ProductChannelListingStore {
	pcls := &SqlProductChannelListingStore{s}

	return pcls
}

func (ps *SqlProductChannelListingStore) createIndexesIfNotExists() {

}
