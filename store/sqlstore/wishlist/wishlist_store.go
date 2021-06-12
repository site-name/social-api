package wishlist

import (
	"github.com/sitename/sitename/model/wishlist"
	"github.com/sitename/sitename/store"
)

type SqlWishlistStore struct {
	store.Store
}

func NewSqlWishlistStore(s store.Store) store.WishlistStore {
	ws := &SqlWishlistStore{s}
	for _, db := range s.GetAllConns() {
		table := db.AddTableWithName(wishlist.Wishlist{}, "Wishlists").SetKeys(false, "Id")
		table.ColMap("Id").SetMaxSize(store.UUID_MAX_LENGTH)
		table.ColMap("Token").SetMaxSize(store.UUID_MAX_LENGTH)
		table.ColMap("UserID").SetMaxSize(store.UUID_MAX_LENGTH)
	}
	return ws
}

func (ws *SqlWishlistStore) CreateIndexesIfNotExists() {
	ws.CreateIndexIfNotExists("idx_wishlists", "Wishlists", "CreateAt")
}
