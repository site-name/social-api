package sqlstore

import (
	"github.com/sitename/sitename/model/product_and_discount"
	"github.com/sitename/sitename/store"
)

type SqlCollectionChannelListingStore struct {
	*SqlStore
}

func newSqlCollectionChannelListingStore(s *SqlStore) store.CollectionChannelListingStore {
	ccls := &SqlCollectionChannelListingStore{s}

	for _, db := range s.GetAllConns() {
		table := db.AddTableWithName(product_and_discount.CollectionChannelListing{}, "CollectionChannelListings").SetKeys(false, "Id")
		table.ColMap("Id").SetMaxSize(UUID_MAX_LENGTH)
		table.ColMap("CollectionID").SetMaxSize(UUID_MAX_LENGTH)
		table.ColMap("ChannelID").SetMaxSize(UUID_MAX_LENGTH)

		table.SetUniqueTogether("ChannelID", "CollectionID")
	}
	return ccls
}

func (ps *SqlCollectionChannelListingStore) createIndexesIfNotExists() {

}
