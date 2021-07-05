package wishlist

import (
	"database/sql"

	"github.com/pkg/errors"
	"github.com/sitename/sitename/model/wishlist"
	"github.com/sitename/sitename/store"
)

type SqlWishlistProductVariantStore struct {
	store.Store
}

func NewSqlWishlistProductVariantStore(s store.Store) store.WishlistProductVariantStore {
	ws := &SqlWishlistProductVariantStore{s}

	for _, db := range s.GetAllConns() {
		table := db.AddTableWithName(wishlist.WishlistProductVariant{}, store.WishlistProductVariantTableName).SetKeys(false, "Id")
		table.ColMap("Id").SetMaxSize(store.UUID_MAX_LENGTH)
		table.ColMap("WishlistItemID").SetMaxSize(store.UUID_MAX_LENGTH)
		table.ColMap("ProductVariantID").SetMaxSize(store.UUID_MAX_LENGTH)

		table.SetUniqueTogether("WishlistItemID", "ProductVariantID")
	}
	return ws
}

func (w *SqlWishlistProductVariantStore) CreateIndexesIfNotExists() {
	w.CreateForeignKeyIfNotExists(store.WishlistProductVariantTableName, "WishlistItemID", store.WishlistItemTableName, "Id", true)
	w.CreateForeignKeyIfNotExists(store.WishlistProductVariantTableName, "ProductVariantID", store.ProductVariantTableName, "Id", true)
}

func (w *SqlWishlistProductVariantStore) Save(item *wishlist.WishlistProductVariant) (*wishlist.WishlistProductVariant, error) {
	item.PreSave()
	if err := item.IsValid(); err != nil {
		return nil, err
	}

	if err := w.GetMaster().Insert(item); err != nil {
		if w.IsUniqueConstraintError(err, []string{"WishlistItemID", "ProductVariantID", "wishlistitemproductvariants_wishlistitemid_productvariantid_key"}) {
			return nil, store.NewErrInvalidInput(store.WishlistProductVariantTableName, "WishlistItemID/ProductVariantID", item.WishlistItemID+"/"+item.ProductVariantID)
		}
		return nil, errors.Wrapf(err, "failed to save wishlist product variant with id=%s", item.Id)
	} else {
		return item, nil
	}
}

func (w *SqlWishlistProductVariantStore) GetById(id string) (*wishlist.WishlistProductVariant, error) {
	if res, err := w.GetReplica().Get(wishlist.WishlistProductVariant{}, id); err != nil {
		if err == sql.ErrNoRows {
			return nil, store.NewErrNotFound(store.WishlistProductVariantTableName, id)
		}
		return nil, errors.Wrapf(err, "failed to find item with Id=%s", id)
	} else {
		return res.(*wishlist.WishlistProductVariant), nil
	}
}
