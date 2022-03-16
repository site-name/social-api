package wishlist

import (
	"database/sql"

	"github.com/Masterminds/squirrel"
	"github.com/mattermost/gorp"
	"github.com/pkg/errors"
	"github.com/sitename/sitename/model"
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

// BulkUpsert inserts or updates given wishlist items then returns it
func (ws *SqlWishlistItemStore) BulkUpsert(transaction *gorp.Transaction, wishlistItems wishlist.WishlistItems) (wishlist.WishlistItems, error) {
	var (
		isSaving        bool
		err             error
		numUpdated      int64
		oldWishlistItem *wishlist.WishlistItem
		upsertor        gorp.SqlExecutor = ws.GetMaster()
	)
	if transaction != nil {
		upsertor = transaction
	}

	for _, wishlistItem := range wishlistItems {
		isSaving = false // reset

		if !model.IsValidId(wishlistItem.Id) {
			wishlistItem.PreSave()
			isSaving = true
		}

		if err := wishlistItem.IsValid(); err != nil {
			return nil, err
		}

		if isSaving {
			err = upsertor.Insert(wishlistItem)
		} else {
			oldWishlistItem, err = ws.GetById(transaction, wishlistItem.Id)
			if err != nil {
				return nil, err
			}

			wishlistItem.CreateAt = oldWishlistItem.CreateAt

			numUpdated, err = upsertor.Update(wishlistItem)
		}

		if err != nil {
			if ws.IsUniqueConstraintError(err, []string{"WishlistID", "ProductID", "wishlistitems_wishlistid_productid_key"}) {
				return nil, store.NewErrInvalidInput(store.WishlistItemTableName, "WishlistID/ProductID", "duplicate")
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
func (ws *SqlWishlistItemStore) GetById(transaction *gorp.Transaction, id string) (*wishlist.WishlistItem, error) {
	var selectOneFunc func(holder interface{}, query string, args ...interface{}) error = ws.GetReplica().SelectOne
	if transaction != nil {
		selectOneFunc = transaction.SelectOne
	}

	var res wishlist.WishlistItem
	if err := selectOneFunc(&res, "SELECT * FROM "+store.WishlistItemTableName+" WHERE Id = :ID", map[string]interface{}{"ID": id}); err != nil {
		if err == sql.ErrNoRows {
			return nil, store.NewErrNotFound(store.WishlistItemTableName, id)
		}
		return nil, errors.Wrapf(err, "failed to find wishlist item with id=%s", id)
	} else {
		return &res, nil
	}
}

func (ws *SqlWishlistItemStore) commonQueryBuilder(option *wishlist.WishlistItemFilterOption) (string, []interface{}, error) {
	query := ws.GetQueryBuilder().
		Select("*").
		From(store.WishlistItemTableName)

	// parse option
	if option.Id != nil {
		query = query.Where(option.Id.ToSquirrel("Id"))
	}
	if option.WishlistID != nil {
		query = query.Where(option.WishlistID.ToSquirrel("WishlistID"))
	}
	if option.ProductID != nil {
		query = query.Where(option.ProductID.ToSquirrel("ProductID"))
	}

	return query.ToSql()
}

// FilterByOption finds and returns a slice of wishlist items filtered using given options
func (ws *SqlWishlistItemStore) FilterByOption(option *wishlist.WishlistItemFilterOption) ([]*wishlist.WishlistItem, error) {
	queryString, args, err := ws.commonQueryBuilder(option)
	if err != nil {
		return nil, errors.Wrap(err, "FilterByOption_ToSql")
	}

	var items []*wishlist.WishlistItem
	if _, err := ws.GetReplica().Select(&items, queryString, args...); err != nil {
		return nil, errors.Wrapf(err, "failed to find wishlist items by given options")
	} else {
		return items, nil
	}
}

// GetByOption finds and returns a wishlist item filtered by given option
func (ws *SqlWishlistItemStore) GetByOption(option *wishlist.WishlistItemFilterOption) (*wishlist.WishlistItem, error) {
	queryString, args, err := ws.commonQueryBuilder(option)
	if err != nil {
		return nil, errors.Wrap(err, "GetByOption_ToSql")
	}

	var res wishlist.WishlistItem
	err = ws.GetReplica().SelectOne(&res, queryString, args...)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, store.NewErrNotFound(store.WishlistItemTableName, "option")
		}
		return nil, errors.Wrap(err, "failed to find wishlist item by given option")
	}

	return &res, nil
}

// DeleteItemsByOption finds and deletes wishlist items that satisfy given filtering options
func (ws *SqlWishlistItemStore) DeleteItemsByOption(transaction *gorp.Transaction, option *wishlist.WishlistItemFilterOption) (int64, error) {
	var runner squirrel.BaseRunner = ws.GetMaster()
	if transaction != nil {
		runner = transaction
	}

	query := ws.GetQueryBuilder().Delete(store.WishlistItemTableName)

	// parse options
	if option.Id != nil {
		query = query.Where(option.Id.ToSquirrel("Id"))
	}
	if option.WishlistID != nil {
		query = query.Where(option.WishlistID.ToSquirrel("WishlistID"))
	}
	if option.ProductID != nil {
		query = query.Where(option.ProductID.ToSquirrel("ProductID"))
	}

	result, err := query.RunWith(runner).Exec()
	if err != nil {
		return 0, errors.Wrap(err, "failed to delete wishlist item wiht given option")
	}
	numDeleted, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "failed to count number of wishlist items deleted")
	}

	return numDeleted, nil
}
