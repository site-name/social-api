package discount

import (
	"database/sql"

	"github.com/Masterminds/squirrel"
	"github.com/pkg/errors"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model/product_and_discount"
	"github.com/sitename/sitename/store"
)

type SqlVoucherTranslationStore struct {
	store.Store
}

func NewSqlVoucherTranslationStore(sqlStore store.Store) store.VoucherTranslationStore {
	return &SqlVoucherTranslationStore{sqlStore}
}

func (s *SqlVoucherTranslationStore) ModelFields(prefix string) model.StringArray {
	res := model.StringArray{
		"Id",
		"LanguageCode",
		"Name",
		"VoucherID",
		"CreateAt",
	}
	if prefix == "" {
		return res
	}

	return res.Map(func(_ int, s string) string {
		return prefix + s
	})
}

// Save inserts given translation into database and returns it
func (vts *SqlVoucherTranslationStore) Save(translation *product_and_discount.VoucherTranslation) (*product_and_discount.VoucherTranslation, error) {
	translation.PreSave()
	if err := translation.IsValid(); err != nil {
		return nil, err
	}

	query := "INSERT INTO " + store.VoucherTranslationTableName + "(" + vts.ModelFields("").Join(",") + ") VALUES (" + vts.ModelFields(":").Join(",") + ")"
	_, err := vts.GetMasterX().NamedExec(query, translation)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to save voucher translation with id=%s", translation.Id)
	}

	return translation, nil
}

// Get finds and returns a voucher translation with given id
func (vts *SqlVoucherTranslationStore) Get(id string) (*product_and_discount.VoucherTranslation, error) {
	var res product_and_discount.VoucherTranslation
	err := vts.GetReplicaX().Get(&res, "SELECT * FROM "+store.VoucherTranslationTableName+" WHERE Id = ?", id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, store.NewErrNotFound(store.VoucherTranslationTableName, id)
		}
		return nil, errors.Wrapf(err, "failed to find voucher translation with id=%s", id)
	}

	return &res, nil
}

func (vts *SqlVoucherTranslationStore) commonQueryBuilder(option *product_and_discount.VoucherTranslationFilterOption) squirrel.SelectBuilder {
	query := vts.GetQueryBuilder().Select("*").From(store.VoucherTranslationTableName)

	// parse option
	if option.Id != nil {
		query = query.Where(option.Id)
	}
	if option.LanguageCode != nil {
		query = query.Where(option.LanguageCode)
	}
	if option.VoucherID != nil {
		query = query.Where(option.VoucherID)
	}

	return query
}

// FilterByOption returns a list of voucher translations filtered using given options
func (vts *SqlVoucherTranslationStore) FilterByOption(option *product_and_discount.VoucherTranslationFilterOption) ([]*product_and_discount.VoucherTranslation, error) {
	query := vts.commonQueryBuilder(option)

	queryString, args, err := query.ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "FilterByOption_ToSql")
	}

	var res []*product_and_discount.VoucherTranslation
	err = vts.GetReplicaX().Select(&res, queryString, args...)
	if err != nil {
		return nil, errors.Wrap(err, "failed to find voucher translations with given options")
	}

	return res, nil
}

// GetByOption finds and returns 1 voucher translation by given options
func (vts *SqlVoucherTranslationStore) GetByOption(option *product_and_discount.VoucherTranslationFilterOption) (*product_and_discount.VoucherTranslation, error) {
	query := vts.commonQueryBuilder(option)

	queryString, args, err := query.ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "GetByOption_ToSql")
	}

	var res product_and_discount.VoucherTranslation
	err = vts.GetReplicaX().Get(&res, queryString, args...)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, store.NewErrNotFound(store.VoucherTranslationTableName, "options")
		}
		return nil, errors.Wrap(err, "failed to find a voucher translation by given option")
	}

	return &res, nil
}
