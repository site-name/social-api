package product

import (
	"github.com/sitename/sitename/store"
)

type SqlCollectionChannelListingStore struct {
	store.Store
}

func NewSqlCollectionChannelListingStore(s store.Store) store.CollectionChannelListingStore {
	return &SqlCollectionChannelListingStore{s}
}
