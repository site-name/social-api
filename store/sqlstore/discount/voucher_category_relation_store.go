package discount

import (
	"database/sql"

	"github.com/pkg/errors"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/modules/util"
	"github.com/sitename/sitename/store"
)

type SqlVoucherCategoryStore struct {
	store.Store
}

func NewSqlVoucherCategoryStore(s store.Store) store.VoucherCategoryStore {
	return &SqlVoucherCategoryStore{s}
}

func (s *SqlVoucherCategoryStore) ModelFields(prefix string) util.AnyArray[string] {
	res := util.AnyArray[string]{
		"Id", "VoucherID", "CategoryID", "CreateAt",
	}
	if prefix == "" {
		return res
	}

	return res.Map(func(_ int, s string) string {
		return prefix + s
	})
}

// Upsert saves or updates given voucher category then returns it with an error
func (vcs *SqlVoucherCategoryStore) Upsert(voucherCategory *model.VoucherCategory) (*model.VoucherCategory, error) {
	var saving bool
	if !model.IsValidId(voucherCategory.Id) {
		voucherCategory.Id = ""
		voucherCategory.PreSave()
		saving = true
	}
	if err := voucherCategory.IsValid(); err != nil {
		return nil, err
	}

	var (
		err        error
		numUpdated int64
	)
	if saving {
		query := "INSERT INTO " + store.VoucherCategoryTableName + " (" + vcs.ModelFields("").Join(",") + ") VALUES (" + vcs.ModelFields(":").Join(",") + ")"
		_, err = vcs.GetMasterX().NamedExec(query, voucherCategory)

	} else {
		query := "UPDATE " + store.VoucherCategoryTableName + " SET " + vcs.
			ModelFields("").
			Map(func(_ int, s string) string {
				return s + "=:" + s
			}).
			Join(",") + " WHERE Id=:Id"

		var result sql.Result
		result, err = vcs.GetMasterX().NamedExec(query, voucherCategory)
		if err == nil && result != nil {
			numUpdated, _ = result.RowsAffected()
		}
	}

	if err != nil {
		if vcs.IsUniqueConstraintError(err, []string{"VoucherID", "CategoryID", "vouchercategories_voucherid_categoryid_key"}) {
			return nil, store.NewErrInvalidInput(store.VoucherCategoryTableName, "VoucherID/CategoryID", "duplicate")
		}
		return nil, errors.Wrapf(err, "failed to upsert voucher-category relation with id=%s", voucherCategory.Id)
	} else if numUpdated > 1 {
		return nil, errors.Errorf("multiple voucher-category relations updated: %d instead of 1", numUpdated)
	}

	return voucherCategory, nil
}

// Get finds a voucher category with given id, then returns it with an error
func (vcs *SqlVoucherCategoryStore) Get(voucherCategoryID string) (*model.VoucherCategory, error) {
	var res model.VoucherCategory
	err := vcs.GetReplicaX().Get(&res, "SELECT * FROM "+store.VoucherCategoryTableName+" WHERE Id = ?", voucherCategoryID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, store.NewErrNotFound(store.VoucherCategoryTableName, voucherCategoryID)
		}
		return nil, errors.Wrapf(err, "failed to find voucher-category relation with id=%s", voucherCategoryID)
	}

	return &res, nil
}

func (s *SqlVoucherCategoryStore) FilterByOptions(options *model.VoucherCategoryFilterOption) ([]*model.VoucherCategory, error) {
	query := s.GetQueryBuilder().Select("*").From(store.VoucherCategoryTableName)

	if options.Id != nil {
		query = query.Where(options.Id)
	}
	if options.VoucherID != nil {
		query = query.Where(options.VoucherID)
	}
	if options.CategoryID != nil {
		query = query.Where(options.CategoryID)
	}

	queryStr, args, err := query.ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "FilterByOptions_ToSql")
	}

	var res []*model.VoucherCategory
	err = s.GetReplicaX().Select(&res, queryStr, args...)
	if err != nil {
		return nil, errors.Wrap(err, "failed to find voucher category relations by given options")
	}

	return res, nil
}
