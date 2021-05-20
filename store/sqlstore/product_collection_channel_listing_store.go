package sqlstore

import "github.com/sitename/sitename/store"

type SqlCollectionChannelListingStore struct {
	*SqlStore
}

func newSqlCollectionChannelListingStore(s *SqlStore) store.CollectionChannelListingStore {
	ccls := &SqlCollectionChannelListingStore{s}

	return ccls
}

func (ps *SqlCollectionChannelListingStore) createIndexesIfNotExists() {

}
