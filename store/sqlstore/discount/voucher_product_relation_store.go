package discount

import (
	"database/sql"

	"github.com/pkg/errors"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model/product_and_discount"
	"github.com/sitename/sitename/store"
)

type SqlVoucherProductStore struct {
	store.Store
}

var VoucherProductDuplicateList = []string{
	"VoucherID", "ProductID", "voucherproducts_voucherid_productid_key",
}

func NewSqlVoucherProductStore(s store.Store) store.VoucherProductStore {
	return &SqlVoucherProductStore{s}
}

func (s *SqlVoucherProductStore) ModelFields(prefix string) model.AnyArray[string] {
	res := model.AnyArray[string]{
		"Id", "VoucherID", "ProductID",
	}
	if prefix == "" {
		return res
	}

	return res.Map(func(_ int, s string) string {
		return prefix + s
	})
}

// Upsert saves or updates given voucher product then returns it with an error
func (vps *SqlVoucherProductStore) Upsert(voucherProduct *product_and_discount.VoucherProduct) (*product_and_discount.VoucherProduct, error) {
	var saving bool
	if voucherProduct.Id == "" {
		voucherProduct.PreSave()
		saving = true
	}
	if err := voucherProduct.IsValid(); err != nil {
		return nil, err
	}

	var (
		err        error
		numUpdated int64
	)
	if saving {
		query := "INSERT INTO " + store.VoucherProductTableName + "(" + vps.ModelFields("").Join(",") + ") VALUES (" + vps.ModelFields(":").Join(",") + ")"
		_, err = vps.GetMasterX().NamedExec(query, voucherProduct)

	} else {
		query := "UPDATE " + store.VoucherProductTableName + " SET " + vps.
			ModelFields("").
			Map(func(_ int, s string) string {
				return s + "=:" + s
			}).
			Join(",") + " WHERE Id=:Id"

		var result sql.Result
		result, err = vps.GetMasterX().NamedExec(query, voucherProduct)
		if err == nil && result != nil {
			numUpdated, _ = result.RowsAffected()
		}
	}

	if err != nil {
		if vps.IsUniqueConstraintError(err, VoucherProductDuplicateList) {
			return nil, store.NewErrInvalidInput(store.VoucherProductTableName, "VoucherID/ProductID", "duplicate")
		}
		return nil, errors.Wrapf(err, "failed to upsert voucher-product relation with id=%s", voucherProduct.Id)
	}

	if numUpdated > 1 {
		return nil, errors.Errorf("multiple voucher-product relations updated: %d instead of 1", numUpdated)
	}

	return voucherProduct, nil
}

// Get finds a voucher product with given id, then returns it with an error
func (vps *SqlVoucherProductStore) Get(voucherProductID string) (*product_and_discount.VoucherProduct, error) {
	var res product_and_discount.VoucherProduct
	err := vps.GetReplicaX().Get(&res, "SELECT * FROM "+store.VoucherProductTableName+" WHERE Id = ?", voucherProductID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, store.NewErrNotFound(store.VoucherProductTableName, voucherProductID)
		}
		return nil, errors.Wrapf(err, "failed to find voucher-product relation with id=%s", voucherProductID)
	}

	return &res, nil
}
