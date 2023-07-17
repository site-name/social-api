package wishlist

import (
	"database/sql"

	"github.com/pkg/errors"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/modules/util"
	"github.com/sitename/sitename/store"
	"github.com/sitename/sitename/store/store_iface"
	"gorm.io/gorm"
)

type SqlWishlistItemStore struct {
	store.Store
}

func NewSqlWishlistItemStore(s store.Store) store.WishlistItemStore {
	return &SqlWishlistItemStore{s}
}

func (s *SqlWishlistItemStore) ModelFields(prefix string) util.AnyArray[string] {
	res := util.AnyArray[string]{
		"Id",
		"WishlistID",
		"ProductID",
		"CreateAt",
	}
	if prefix == "" {
		return res
	}

	return res.Map(func(_ int, s string) string {
		return prefix + s
	})
}

// BulkUpsert inserts or updates given wishlist items then returns it
func (ws *SqlWishlistItemStore) BulkUpsert(transaction store_iface.SqlxExecutor, wishlistItems model.WishlistItems) (model.WishlistItems, error) {
	var (
		upsertor store_iface.SqlxExecutor = ws.GetMasterX()
	)
	if transaction != nil {
		upsertor = transaction
	}

	var (
		saveQuery   = "INSERT INTO " + model.WishlistItemTableName + "(" + ws.ModelFields("").Join(",") + ") VALUES (" + ws.ModelFields(":").Join(",") + ")"
		updateQuery = "UPDATE " + model.WishlistItemTableName + " SET " + ws.
				ModelFields("").
				Map(func(_ int, s string) string {
				return s + "=:" + s
			}).
			Join(",") + " WHERE Id=:Id"
	)

	for _, wishlistItem := range wishlistItems {
		var (
			err        error
			numUpdated int64
			isSaving   bool // false
		)

		if !model.IsValidId(wishlistItem.Id) {
			wishlistItem.PreSave()
			isSaving = true
		}

		if err := wishlistItem.IsValid(); err != nil {
			return nil, err
		}

		if isSaving {
			_, err = upsertor.NamedExec(saveQuery, wishlistItem)
		} else {
			var result sql.Result
			result, err = upsertor.NamedExec(updateQuery, wishlistItem)
			if err == nil && result != nil {
				numUpdated, _ = result.RowsAffected()
			}
		}

		if err != nil {
			if ws.IsUniqueConstraintError(err, []string{"WishlistID", "ProductID", "wishlistitems_wishlistid_productid_key"}) {
				return nil, store.NewErrInvalidInput(model.WishlistItemTableName, "WishlistID/ProductID", "duplicate")
			}
			return nil, errors.Wrapf(err, "failed to upsert wishlist item with id=%s", wishlistItem.Id)
		}
		if numUpdated > 1 {
			return nil, errors.Errorf("multiple wishlist items (id=%s) were updated for: %d instead of 1", wishlistItem.Id, numUpdated)
		}
	}

	return wishlistItems, nil
}

// GetById finds and returns a wishlist item by given id
func (ws *SqlWishlistItemStore) GetById(transaction store_iface.SqlxExecutor, id string) (*model.WishlistItem, error) {
	var executor store_iface.SqlxExecutor = ws.GetReplicaX()
	if transaction != nil {
		executor = transaction
	}

	var res model.WishlistItem
	if err := executor.Get(&res, "SELECT * FROM "+model.WishlistItemTableName+" WHERE Id = ?", id); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, store.NewErrNotFound(model.WishlistItemTableName, id)
		}
		return nil, errors.Wrapf(err, "failed to find wishlist item with id=%s", id)
	} else {
		return &res, nil
	}
}

func (ws *SqlWishlistItemStore) commonQueryBuilder(option *model.WishlistItemFilterOption) (string, []interface{}, error) {
	query := ws.GetQueryBuilder().
		Select("*").
		From(model.WishlistItemTableName)

	// parse option
	if option.Id != nil {
		query = query.Where(option.Id)
	}
	if option.WishlistID != nil {
		query = query.Where(option.WishlistID)
	}
	if option.ProductID != nil {
		query = query.Where(option.ProductID)
	}

	return query.ToSql()
}

// FilterByOption finds and returns a slice of wishlist items filtered using given options
func (ws *SqlWishlistItemStore) FilterByOption(option *model.WishlistItemFilterOption) ([]*model.WishlistItem, error) {
	queryString, args, err := ws.commonQueryBuilder(option)
	if err != nil {
		return nil, errors.Wrap(err, "FilterByOption_ToSql")
	}

	var items []*model.WishlistItem
	if err := ws.GetReplicaX().Select(&items, queryString, args...); err != nil {
		return nil, errors.Wrapf(err, "failed to find wishlist items by given options")
	} else {
		return items, nil
	}
}

// GetByOption finds and returns a wishlist item filtered by given option
func (ws *SqlWishlistItemStore) GetByOption(option *model.WishlistItemFilterOption) (*model.WishlistItem, error) {
	queryString, args, err := ws.commonQueryBuilder(option)
	if err != nil {
		return nil, errors.Wrap(err, "GetByOption_ToSql")
	}

	var res model.WishlistItem
	err = ws.GetReplicaX().Get(&res, queryString, args...)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, store.NewErrNotFound(model.WishlistItemTableName, "option")
		}
		return nil, errors.Wrap(err, "failed to find wishlist item by given option")
	}

	return &res, nil
}

// DeleteItemsByOption finds and deletes wishlist items that satisfy given filtering options
func (ws *SqlWishlistItemStore) DeleteItemsByOption(transaction store_iface.SqlxExecutor, option *model.WishlistItemFilterOption) (int64, error) {
	var runner store_iface.SqlxExecutor = ws.GetMasterX()
	if transaction != nil {
		runner = transaction
	}

	query := ws.GetQueryBuilder().Delete(model.WishlistItemTableName)

	// parse options
	if option.Id != nil {
		query = query.Where(option.Id)
	}
	if option.WishlistID != nil {
		query = query.Where(option.WishlistID)
	}
	if option.ProductID != nil {
		query = query.Where(option.ProductID)
	}

	queryString, args, err := query.ToSql()
	if err != nil {
		return 0, errors.Wrap(err, "DeleteItemsByOption_ToSql")
	}

	result, err := runner.Exec(queryString, args...)
	if err != nil {
		return 0, errors.Wrap(err, "failed to delete wishlist item wiht given option")
	}
	numDeleted, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "failed to count number of wishlist items deleted")
	}

	return numDeleted, nil
}
