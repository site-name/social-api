package sqlstore

import (
	"github.com/sitename/sitename/model/wishlist"
	"github.com/sitename/sitename/store"
)

type SqlWishlistStore struct {
	*SqlStore
}

func newSqlWishlistStore(s *SqlStore) store.WishlistStore {
	ws := &SqlWishlistStore{s}
	for _, db := range s.GetAllConns() {
		table := db.AddTableWithName(wishlist.Wishlist{}, "Wishlists").SetKeys(false, "Id")
		table.ColMap("Id").SetMaxSize(UUID_MAX_LENGTH)
		table.ColMap("Token").SetMaxSize(UUID_MAX_LENGTH)
		table.ColMap("UserID").SetMaxSize(UUID_MAX_LENGTH)
	}
	return ws
}

func (ws *SqlWishlistStore) createIndexesIfNotExists() {
	ws.CreateIndexIfNotExists("idx_wishlists", "Wishlists", "CreateAt")
}
