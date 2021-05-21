package sqlstore

import (
	"github.com/sitename/sitename/model/wishlist"
	"github.com/sitename/sitename/store"
)

type SqlWishlistItemStore struct {
	*SqlStore
}

func newSqlWishlistItemStore(s *SqlStore) store.WishlistItemStore {
	ws := &SqlWishlistItemStore{s}
	for _, db := range s.GetAllConns() {
		table := db.AddTableWithName(wishlist.WishlistItem{}, "WishlistItems").SetKeys(false, "Id")
		table.ColMap("Id").SetMaxSize(UUID_MAX_LENGTH)
		table.ColMap("WishlistID").SetMaxSize(UUID_MAX_LENGTH)
		table.ColMap("ProductID").SetMaxSize(UUID_MAX_LENGTH)

		table.SetUniqueTogether("WishlistID", "ProductID")
	}
	return ws
}

func (ws *SqlWishlistItemStore) createIndexesIfNotExists() {
	ws.CreateIndexIfNotExists("idx_wishlist_items", "WishlistItems", "CreateAt")
}
