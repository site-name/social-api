package discount

import (
	"database/sql"

	"github.com/pkg/errors"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/modules/util"
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

func (s *SqlVoucherProductStore) ModelFields(prefix string) util.AnyArray[string] {
	res := util.AnyArray[string]{
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
func (vps *SqlVoucherProductStore) Upsert(voucherProduct *model.VoucherProduct) (*model.VoucherProduct, error) {
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
func (vps *SqlVoucherProductStore) Get(voucherProductID string) (*model.VoucherProduct, error) {
	var res model.VoucherProduct
	err := vps.GetReplicaX().Get(&res, "SELECT * FROM "+store.VoucherProductTableName+" WHERE Id = ?", voucherProductID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, store.NewErrNotFound(store.VoucherProductTableName, voucherProductID)
		}
		return nil, errors.Wrapf(err, "failed to find voucher-product relation with id=%s", voucherProductID)
	}

	return &res, nil
}

func (s *SqlVoucherProductStore) FilterByOptions(options *model.VoucherProductFilterOptions) ([]*model.VoucherProduct, error) {
	query := s.GetQueryBuilder().Select("*").From(store.VoucherProductTableName)

	if options.ProductID != nil {
		query = query.Where(options.ProductID)
	}
	if options.VoucherID != nil {
		query = query.Where(options.VoucherID)
	}

	queryString, args, err := query.ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "FilterByOptions_ToSql")
	}

	var res []*model.VoucherProduct
	err = s.GetReplicaX().Select(&res, queryString, args...)
	if err != nil {
		return nil, errors.Wrap(err, "failed to find voucher product relations by options")
	}

	return res, nil
}
