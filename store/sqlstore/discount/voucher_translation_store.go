package discount

import (
	"database/sql"

	"github.com/Masterminds/squirrel"
	"github.com/pkg/errors"
	"github.com/sitename/sitename/model/product_and_discount"
	"github.com/sitename/sitename/store"
)

type SqlVoucherTranslationStore struct {
	store.Store
}

func NewSqlVoucherTranslationStore(sqlStore store.Store) store.VoucherTranslationStore {
	vts := &SqlVoucherTranslationStore{sqlStore}

	for _, db := range sqlStore.GetAllConns() {
		table := db.AddTableWithName(product_and_discount.VoucherTranslation{}, store.VoucherTranslationTableName).SetKeys(false, "Id")
		table.ColMap("Id").SetMaxSize(store.UUID_MAX_LENGTH)
		table.ColMap("VoucherID").SetMaxSize(store.UUID_MAX_LENGTH)
		table.ColMap("Name").SetMaxSize(product_and_discount.VOUCHER_NAME_MAX_LENGTH)
		table.ColMap("LanguageCode").SetMaxSize(10)

		table.SetUniqueTogether("LanguageCode", "VoucherID")
	}

	return vts
}

func (vts *SqlVoucherTranslationStore) CreateIndexesIfNotExists() {
	vts.CreateForeignKeyIfNotExists(store.VoucherTranslationTableName, "VoucherID", store.VoucherTableName, "Id", true)
}

// Save inserts given translation into database and returns it
func (vts *SqlVoucherTranslationStore) Save(translation *product_and_discount.VoucherTranslation) (*product_and_discount.VoucherTranslation, error) {
	translation.PreSave()
	if err := translation.IsValid(); err != nil {
		return nil, err
	}

	err := vts.GetMaster().Insert(translation)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to save voucher translation with id=%s", translation.Id)
	}

	return translation, nil
}

// Get finds and returns a voucher translation with given id
func (vts *SqlVoucherTranslationStore) Get(id string) (*product_and_discount.VoucherTranslation, error) {
	var res product_and_discount.VoucherTranslation
	err := vts.GetReplica().SelectOne(&res, "SELECT * FROM "+store.VoucherTranslationTableName+" WHERE Id = :ID", map[string]interface{}{"ID": id})
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
	_, err = vts.GetReplica().Select(&res, queryString, args...)
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
	err = vts.GetReplica().SelectOne(&res, queryString, args...)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, store.NewErrNotFound(store.VoucherTranslationTableName, "options")
		}
		return nil, errors.Wrap(err, "failed to find a voucher translation by given option")
	}

	return &res, nil
}
