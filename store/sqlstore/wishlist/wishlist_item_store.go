package wishlist

import (
	"database/sql"

	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model_helper"
	"github.com/sitename/sitename/store"
	"github.com/volatiletech/sqlboiler/v4/boil"
)

type SqlWishlistItemStore struct {
	store.Store
}

func NewSqlWishlistItemStore(s store.Store) store.WishlistItemStore {
	return &SqlWishlistItemStore{s}
}

// BulkUpsert inserts or updates given wishlist items then returns it
func (ws *SqlWishlistItemStore) BulkUpsert(transaction boil.ContextTransactor, wishlistItems model.WishlistItemSlice) (model.WishlistItemSlice, error) {
	if transaction == nil {
		transaction = ws.GetMaster()
	}

	for _, wishlistItem := range wishlistItems {
		if wishlistItem == nil {
			continue
		}

		isSaving := false
		if wishlistItem.ID == "" {
			isSaving = true
			model_helper.WishlistItemPreSave(wishlistItem)
		}

		if err := model_helper.WishlistItemIsValid(*wishlistItem); err != nil {
			return nil, err
		}

		var err error
		if isSaving {
			err = wishlistItem.Insert(transaction, boil.Infer())
		} else {
			_, err = wishlistItem.Update(transaction, boil.Blacklist(model.WishlistItemColumns.CreatedAt))
		}

		if err != nil {
			if ws.IsUniqueConstraintError(err, []string{model.WishlistItemColumns.WishlistID, model.WishlistItemColumns.ProductID, "wishlist_items_wishlist_id_product_id_key"}) {
				return nil, store.NewErrInvalidInput(model.TableNames.WishlistItems, model.WishlistItemColumns.WishlistID+"/"+model.WishlistItemColumns.ProductID, "unique")
			}
			return nil, err
		}
	}

	return wishlistItems, nil
}

// GetById finds and returns a wishlist item by given id
func (ws *SqlWishlistItemStore) GetById(id string) (*model.WishlistItem, error) {
	record, err := model.FindWishlistItem(ws.GetReplica(), id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, store.NewErrNotFound(model.TableNames.WishlistItems, id)
		}
		return nil, err
	}
	return record, nil
}

// FilterByOption finds and returns a slice of wishlist items filtered using given options
func (ws *SqlWishlistItemStore) FilterByOption(option model_helper.WishlistItemFilterOption) (model.WishlistItemSlice, error) {
	return model.WishlistItems(option.Conditions...).All(ws.GetReplica())
}

// GetByOption finds and returns a wishlist item filtered by given option
func (ws *SqlWishlistItemStore) GetByOption(option model_helper.WishlistItemFilterOption) (*model.WishlistItem, error) {
	item, err := model.WishlistItems(option.Conditions...).One(ws.GetReplica())
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, store.NewErrNotFound(model.TableNames.WishlistItems, "options")
		}
		return nil, err
	}
	return item, nil
}

// DeleteItemsByOption finds and deletes wishlist items that satisfy given filtering options
func (ws *SqlWishlistItemStore) Delete(transaction boil.ContextTransactor, ids []string) (int64, error) {
	if transaction == nil {
		transaction = ws.GetMaster()
	}

	return model.WishlistItems(model.WishlistItemWhere.ID.IN(ids)).DeleteAll(transaction)
}
