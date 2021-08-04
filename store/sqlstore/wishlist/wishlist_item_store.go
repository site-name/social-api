package wishlist

import (
	"database/sql"

	"github.com/pkg/errors"
	"github.com/sitename/sitename/model/wishlist"
	"github.com/sitename/sitename/store"
)

type SqlWishlistItemStore struct {
	store.Store
}

func NewSqlWishlistItemStore(s store.Store) store.WishlistItemStore {
	ws := &SqlWishlistItemStore{s}
	for _, db := range s.GetAllConns() {
		table := db.AddTableWithName(wishlist.WishlistItem{}, store.WishlistItemTableName).SetKeys(false, "Id")
		table.ColMap("Id").SetMaxSize(store.UUID_MAX_LENGTH)
		table.ColMap("WishlistID").SetMaxSize(store.UUID_MAX_LENGTH)
		table.ColMap("ProductID").SetMaxSize(store.UUID_MAX_LENGTH)

		table.SetUniqueTogether("WishlistID", "ProductID")
	}
	return ws
}

func (ws *SqlWishlistItemStore) CreateIndexesIfNotExists() {
	ws.CreateIndexIfNotExists("idx_wishlist_items", store.WishlistItemTableName, "CreateAt")
	ws.CreateForeignKeyIfNotExists(store.WishlistItemTableName, "WishlistID", store.WishlistTableName, "Id", true)
	ws.CreateForeignKeyIfNotExists(store.WishlistItemTableName, "ProductID", store.ProductVariantTableName, "Id", true)
}

func (ws *SqlWishlistItemStore) Save(item *wishlist.WishlistItem) (*wishlist.WishlistItem, error) {
	item.PreSave()
	if err := item.IsValid(); err != nil {
		return nil, err
	}

	if err := ws.GetMaster().Insert(item); err != nil {
		return nil, err
	}

	return item, nil
}

func (ws *SqlWishlistItemStore) GetById(id string) (*wishlist.WishlistItem, error) {
	var res wishlist.WishlistItem
	if err := ws.GetReplica().SelectOne(&res, "SELECT * FROM "+store.WishlistItemTableName+" WHERE Id = :ID", map[string]interface{}{"ID": id}); err != nil {
		if err == sql.ErrNoRows {
			return nil, store.NewErrNotFound(store.WishlistItemTableName, id)
		}
		return nil, errors.Wrapf(err, "failed to find wishlist item with id=%s", id)
	} else {
		return &res, nil
	}
}

func (ws *SqlWishlistItemStore) WishlistItemsByWishlistId(id string) ([]*wishlist.WishlistItem, error) {
	var items []*wishlist.WishlistItem
	if _, err := ws.GetReplica().Select(&items, "SELECT * FROM "+store.WishlistItemTableName+" WHERE WishlistID = :WishlistID", map[string]interface{}{"WishlistID": id}); err != nil {
		return nil, errors.Wrapf(err, "failed to find wishlist items belong to wishlistId=%s", id)
	} else {
		return items, nil
	}
}
