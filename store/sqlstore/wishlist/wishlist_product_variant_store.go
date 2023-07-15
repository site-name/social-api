package wishlist

import (
	"database/sql"

	"github.com/pkg/errors"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/modules/util"
	"github.com/sitename/sitename/store"
	"github.com/sitename/sitename/store/store_iface"
)

type SqlWishlistItemProductVariantStore struct {
	store.Store
}

func NewSqlWishlistItemProductVariantStore(s store.Store) store.WishlistItemProductVariantStore {
	return &SqlWishlistItemProductVariantStore{s}
}

func (s *SqlWishlistItemProductVariantStore) ModelFields(prefix string) util.AnyArray[string] {
	res := util.AnyArray[string]{
		"Id",
		"WishlistItemID",
		"ProductVariantID",
	}
	if prefix == "" {
		return res
	}

	return res.Map(func(_ int, s string) string {
		return prefix + s
	})
}

// Save inserts given wishlist item-product variant relation into database and returns it
func (w *SqlWishlistItemProductVariantStore) Save(item *model.WishlistItemProductVariant) (*model.WishlistItemProductVariant, error) {
	item.PreSave()
	if err := item.IsValid(); err != nil {
		return nil, err
	}

	query := "INSERT INTO " + model.WishlistItemProductVariantTableName + "(" + w.ModelFields("").Join(",") + ") VALUES (" + w.ModelFields(":").Join(",") + ")"
	if _, err := w.GetMasterX().NamedExec(query, item); err != nil {
		if w.IsUniqueConstraintError(err, []string{"WishlistItemID", "ProductVariantID", "wishlistitemproductvariants_wishlistitemid_productvariantid_key"}) {
			return nil, store.NewErrInvalidInput(model.WishlistItemProductVariantTableName, "WishlistItemID/ProductVariantID", item.WishlistItemID+"/"+item.ProductVariantID)
		}
		return nil, errors.Wrapf(err, "failed to save wishlist product variant with id=%s", item.Id)
	}

	return item, nil
}

func (w *SqlWishlistItemProductVariantStore) BulkUpsert(transaction store_iface.SqlxExecutor, relations []*model.WishlistItemProductVariant) ([]*model.WishlistItemProductVariant, error) {
	var (
		executor    store_iface.SqlxExecutor = w.GetMasterX()
		saveQuery                            = "INSERT INTO " + model.WishlistItemProductVariantTableName + "(" + w.ModelFields("").Join(",") + ") VALUES (" + w.ModelFields(":").Join(",") + ")"
		updateQuery                          = "UPDATE " + model.WishlistItemProductVariantTableName + " SET " + w.
				ModelFields("").
				Map(func(_ int, s string) string {
				return s + "=:" + s
			}).
			Join(",") + " WHERE Id=:Id"
	)
	if transaction != nil {
		executor = transaction
	}

	for _, relation := range relations {
		var (
			isSaving   bool // false
			err        error
			numUpdated int64
		)

		if !model.IsValidId(relation.Id) {
			relation.PreSave()
			isSaving = true
		}

		if err := relation.IsValid(); err != nil {
			return nil, err
		}

		if isSaving {
			_, err = executor.NamedExec(saveQuery, relation)

		} else {
			var result sql.Result
			result, err = executor.NamedExec(updateQuery, relation)
			if err == nil && result != nil {
				numUpdated, _ = result.RowsAffected()
			}
		}

		if err != nil {
			if w.IsUniqueConstraintError(err, []string{"WishlistItemID", "ProductVariantID", "wishlistitemproductvariants_wishlistitemid_productvariantid_key"}) {
				return nil, store.NewErrInvalidInput(model.WishlistItemProductVariantTableName, "WishlistItemID/ProductVariantID", "duplicate")
			}
			return nil, errors.Wrapf(err, "failed to upsert relation with id=%s", relation.Id)
		}
		if numUpdated > 1 {
			return nil, errors.Errorf("multiple wishlist item-product variant relations were updated: %d instead of 1", numUpdated)
		}
	}

	return relations, nil
}

// GetById finds and returns a product variant-wishlist item relation and returns it
func (w *SqlWishlistItemProductVariantStore) GetById(transaction store_iface.SqlxExecutor, id string) (*model.WishlistItemProductVariant, error) {
	var selector store_iface.SqlxExecutor = w.GetReplicaX()
	if transaction != nil {
		selector = transaction
	}

	var res model.WishlistItemProductVariant
	if err := selector.Get(&res, "SELECT * FROM "+model.WishlistItemProductVariantTableName+" WHERE Id = ?", id); err != nil {
		if err == sql.ErrNoRows {
			return nil, store.NewErrNotFound(model.WishlistItemProductVariantTableName, id)
		}
		return nil, errors.Wrapf(err, "failed to find item with Id=%s", id)
	}
	return &res, nil
}

// DeleteRelation deletes a product variant-wishlist item relation and counts numeber of relations left in database
func (w *SqlWishlistItemProductVariantStore) DeleteRelation(relation *model.WishlistItemProductVariant) (int64, error) {
	query := w.GetQueryBuilder().
		Delete(model.WishlistItemProductVariantTableName)
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

	result, err := w.GetMasterX().Exec(queryString, args...)
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

	var numOfRelationsLeft int64
	err = w.GetMasterX().Get(&numOfRelationsLeft, "SELECT COUNT(Id) FROM "+model.WishlistItemProductVariantTableName)
	if err != nil {
		return 0, errors.Wrap(err, "failed to count number of wishlist item-product variant left")
	}

	return numOfRelationsLeft, nil
}
