package discount

import (
	"database/sql"

	"github.com/pkg/errors"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model/product_and_discount"
	"github.com/sitename/sitename/store"
)

type SqlVoucherCollectionStore struct {
	store.Store
}

func NewSqlVoucherCollectionStore(s store.Store) store.VoucherCollectionStore {
	return &SqlVoucherCollectionStore{s}
}

func (s *SqlVoucherCollectionStore) ModelFields(prefix string) model.StringArray {
	res := model.StringArray{
		"Id", "VoucherID", "CollectionID",
	}
	if prefix == "" {
		return res
	}

	return res.Map(func(_ int, s string) string {
		return prefix + s
	})
}

// Upsert saves or updates given voucher collection then returns it with an error
func (vcs *SqlVoucherCollectionStore) Upsert(voucherCollection *product_and_discount.VoucherCollection) (*product_and_discount.VoucherCollection, error) {
	var saving bool
	if voucherCollection.Id == "" {
		voucherCollection.PreSave()
		saving = true
	}

	if err := voucherCollection.IsValid(); err != nil {
		return nil, err
	}

	var (
		err        error
		numUpdated int64
	)
	if saving {
		query := "INSERT INTO " + store.VoucherCollectionTableName + "(" + vcs.ModelFields("").Join(",") + ") VALUES (" + vcs.ModelFields(":").Join(",") + ")"
		_, err = vcs.GetMasterX().NamedExec(query, voucherCollection)

	} else {
		query := "UPDATE " + store.VoucherCollectionTableName + " SET " + vcs.
			ModelFields("").
			Map(func(_ int, s string) string {
				return s + "=:" + s
			}).
			Join(",") + " WHERE Id=:Id"

		var result sql.Result
		result, err = vcs.GetMasterX().NamedExec(query, voucherCollection)
		if err == nil && result != nil {
			numUpdated, _ = result.RowsAffected()
		}
	}

	if err != nil {
		if vcs.IsUniqueConstraintError(err, []string{"VoucherID", "CollectionID", "vouchercollections_voucherid_collectionid_key"}) {
			return nil, store.NewErrInvalidInput(store.VoucherCollectionTableName, "VoucherID/CollectionID", "duplicate")
		}
		return nil, errors.Wrapf(err, "failed to upsert voucher-collection relation with id=%s", voucherCollection.Id)
	}

	if numUpdated > 1 {
		return nil, errors.Errorf("multiple voucher-collection relations updated: %d instead of 1", numUpdated)
	}

	return voucherCollection, nil
}

// Get finds a voucher collection with given id, then returns it with an error
func (vcs *SqlVoucherCollectionStore) Get(voucherCollectionID string) (*product_and_discount.VoucherCollection, error) {
	var res product_and_discount.VoucherCollection
	err := vcs.GetReplicaX().Get(&res, "SELECT * FROM "+store.VoucherCollectionTableName+" WHERE Id = ?", voucherCollectionID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, store.NewErrNotFound(store.VoucherCollectionTableName, voucherCollectionID)
		}
		return nil, errors.Wrapf(err, "failed to find voucher-collection relation with id=%s", voucherCollectionID)
	}

	return &res, nil
}
