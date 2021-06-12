package wishlist

import (
	"github.com/sitename/sitename/model/wishlist"
	"github.com/sitename/sitename/store"
)

type SqlWishlistItemStore struct {
	store.Store
}

func NewSqlWishlistItemStore(s store.Store) store.WishlistItemStore {
	ws := &SqlWishlistItemStore{s}
	for _, db := range s.GetAllConns() {
		table := db.AddTableWithName(wishlist.WishlistItem{}, "WishlistItems").SetKeys(false, "Id")
		table.ColMap("Id").SetMaxSize(store.UUID_MAX_LENGTH)
		table.ColMap("WishlistID").SetMaxSize(store.UUID_MAX_LENGTH)
		table.ColMap("ProductID").SetMaxSize(store.UUID_MAX_LENGTH)

		table.SetUniqueTogether("WishlistID", "ProductID")
	}
	return ws
}

func (ws *SqlWishlistItemStore) CreateIndexesIfNotExists() {
	ws.CreateIndexIfNotExists("idx_wishlist_items", "WishlistItems", "CreateAt")
}
