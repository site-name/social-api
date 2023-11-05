package wishlist

import (
	"github.com/pkg/errors"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/store"
	"gorm.io/gorm"
)

type SqlWishlistItemStore struct {
	store.Store
}

func NewSqlWishlistItemStore(s store.Store) store.WishlistItemStore {
	return &SqlWishlistItemStore{s}
}

// BulkUpsert inserts or updates given wishlist items then returns it
func (ws *SqlWishlistItemStore) BulkUpsert(transaction *gorm.DB, wishlistItems model.WishlistItems) (model.WishlistItems, error) {
	if transaction == nil {
		transaction = ws.GetMaster()
	}

	for _, wishlistItem := range wishlistItems {
		err := transaction.Save(wishlistItem).Error
		if err != nil {
			if ws.IsUniqueConstraintError(err, []string{"WishlistID", "ProductID", "wishlistitems_wishlistid_productid_key"}) {
				return nil, store.NewErrInvalidInput(model.WishlistItemTableName, "WishlistID/ProductID", "duplicate")
			}
			return nil, errors.Wrapf(err, "failed to upsert wishlist item with id=%s", wishlistItem.Id)
		}
	}

	return wishlistItems, nil
}

// GetById finds and returns a wishlist item by given id
func (ws *SqlWishlistItemStore) GetById(id string) (*model.WishlistItem, error) {
	var res model.WishlistItem
	if err := ws.GetReplica().First(&res, "Id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, store.NewErrNotFound(model.WishlistItemTableName, id)
		}
		return nil, errors.Wrapf(err, "failed to find wishlist item with id=%s", id)
	}
	return &res, nil
}

// FilterByOption finds and returns a slice of wishlist items filtered using given options
func (ws *SqlWishlistItemStore) FilterByOption(option *model.WishlistItemFilterOption) ([]*model.WishlistItem, error) {
	var items []*model.WishlistItem
	args, err := store.BuildSqlizer(option.Conditions, "WishlistItem_FilterByOption")
	if err != nil {
		return nil, err
	}
	if err := ws.GetReplica().Find(&items, args...).Error; err != nil {
		return nil, errors.Wrapf(err, "failed to find wishlist items by given options")
	} else {
		return items, nil
	}
}

// GetByOption finds and returns a wishlist item filtered by given option
func (ws *SqlWishlistItemStore) GetByOption(option *model.WishlistItemFilterOption) (*model.WishlistItem, error) {
	var res model.WishlistItem
	args, err := store.BuildSqlizer(option.Conditions, "WishlstItem_GetByOption")
	if err != nil {
		return nil, err
	}
	err = ws.GetReplica().First(&res, args...).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, store.NewErrNotFound(model.WishlistItemTableName, "option")
		}
		return nil, errors.Wrap(err, "failed to find wishlist item by given option")
	}

	return &res, nil
}

// DeleteItemsByOption finds and deletes wishlist items that satisfy given filtering options
func (ws *SqlWishlistItemStore) DeleteItemsByOption(transaction *gorm.DB, option *model.WishlistItemFilterOption) (int64, error) {
	if transaction == nil {
		transaction = ws.GetMaster()
	}

	query := ws.GetQueryBuilder().Delete(model.WishlistItemTableName).Where(option.Conditions)
	queryString, args, err := query.ToSql()
	if err != nil {
		return 0, errors.Wrap(err, "DeleteItemsByOption_ToSql")
	}

	result := transaction.Raw(queryString, args...)
	if result.Error != nil {
		return 0, errors.Wrap(result.Error, "failed to delete wishlist item wiht given option")
	}

	return result.RowsAffected, nil
}
