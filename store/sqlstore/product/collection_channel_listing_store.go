package product

import (
	"github.com/sitename/sitename/model/product_and_discount"
	"github.com/sitename/sitename/store"
)

type SqlCollectionChannelListingStore struct {
	store.Store
}

func NewSqlCollectionChannelListingStore(s store.Store) store.CollectionChannelListingStore {
	ccls := &SqlCollectionChannelListingStore{s}

	for _, db := range s.GetAllConns() {
		table := db.AddTableWithName(product_and_discount.CollectionChannelListing{}, "CollectionChannelListings").SetKeys(false, "Id")
		table.ColMap("Id").SetMaxSize(store.UUID_MAX_LENGTH)
		table.ColMap("CollectionID").SetMaxSize(store.UUID_MAX_LENGTH)
		table.ColMap("ChannelID").SetMaxSize(store.UUID_MAX_LENGTH)

		table.SetUniqueTogether("ChannelID", "CollectionID")
	}
	return ccls
}

func (ps *SqlCollectionChannelListingStore) CreateIndexesIfNotExists() {

}
