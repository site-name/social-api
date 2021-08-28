package wishlist

import (
	"database/sql"

	"github.com/pkg/errors"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model/wishlist"
	"github.com/sitename/sitename/store"
)

type SqlWishlistItemProductVariantStore struct {
	store.Store
}

func NewSqlWishlistItemProductVariantStore(s store.Store) store.WishlistItemProductVariantStore {
	ws := &SqlWishlistItemProductVariantStore{s}

	for _, db := range s.GetAllConns() {
		table := db.AddTableWithName(wishlist.WishlistItemProductVariant{}, store.WishlistProductVariantTableName).SetKeys(false, "Id")
		table.ColMap("Id").SetMaxSize(store.UUID_MAX_LENGTH)
		table.ColMap("WishlistItemID").SetMaxSize(store.UUID_MAX_LENGTH)
		table.ColMap("ProductVariantID").SetMaxSize(store.UUID_MAX_LENGTH)

		table.SetUniqueTogether("WishlistItemID", "ProductVariantID")
	}
	return ws
}

func (w *SqlWishlistItemProductVariantStore) CreateIndexesIfNotExists() {
	w.CreateForeignKeyIfNotExists(store.WishlistProductVariantTableName, "WishlistItemID", store.WishlistItemTableName, "Id", true)
	w.CreateForeignKeyIfNotExists(store.WishlistProductVariantTableName, "ProductVariantID", store.ProductVariantTableName, "Id", true)
}

// Save inserts given wishlist item-product variant relation into database and returns it
func (w *SqlWishlistItemProductVariantStore) Save(item *wishlist.WishlistItemProductVariant) (*wishlist.WishlistItemProductVariant, error) {
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

// GetById finds and returns a product variant-wishlist item relation and returns it
func (w *SqlWishlistItemProductVariantStore) GetById(id string) (*wishlist.WishlistItemProductVariant, error) {
	var res wishlist.WishlistItemProductVariant
	if err := w.GetReplica().SelectOne(&res, "SELECT * FROM "+store.WishlistProductVariantTableName+" WHERE Id = :ID", map[string]interface{}{"ID": id}); err != nil {
		if err == sql.ErrNoRows {
			return nil, store.NewErrNotFound(store.WishlistProductVariantTableName, id)
		}
		return nil, errors.Wrapf(err, "failed to find item with Id=%s", id)
	} else {
		return &res, nil
	}
}

// DeleteRelation deletes a product variant-wishlist item relation and counts numeber of relations left in database
func (w *SqlWishlistItemProductVariantStore) DeleteRelation(relation *wishlist.WishlistItemProductVariant) (int64, error) {
	transaction, err := w.GetMaster().Begin()
	if err != nil {
		return 0, errors.Wrap(err, "transaction_begin")
	}
	defer w.FinalizeTransaction(transaction)

	query := w.GetQueryBuilder().
		Delete(store.WishlistProductVariantTableName)
	if model.IsValidId(relation.Id) {
		query = query.Where("Id = ?", relation.Id)
	}
	if model.IsValidId(relation.WishlistItemID) {
		query = query.Where("WishlistItemID = ?", relation.WishlistItemID)
	}
	if model.IsValidId(relation.ProductVariantID) {
		query = query.Where("ProductVariantID = ?", relation.ProductVariantID)
	}

	queryString, args, err := query.ToSql()
	if err != nil {
		return 0, errors.Wrap(err, "DeleteRelation_ToSql")
	}

	result, err := transaction.Exec(queryString, args...)
	if err != nil {
		return 0, errors.Wrap(err, "failed to delete a wishlist item-product variant relation")
	}
	numDeleted, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "failed to count number of relations deleted")
	}
	if numDeleted > 1 {
		return 0, errors.Errorf("multiple wishlist item-product variant relations were deleted: %d instead of 1", numDeleted)
	}

	numOfRelationsLeft, err := transaction.SelectInt("SELECT COUNT(Id) FROM " + store.WishlistProductVariantTableName)
	if err != nil {
		return 0, errors.Wrap(err, "failed to count number of wishlist item-product variant left")
	}

	if err = transaction.Commit(); err != nil {
		return 0, errors.Wrap(err, "transaction_commit")
	}

	return numOfRelationsLeft, nil
}
