package discount

import (
	"github.com/Masterminds/squirrel"
	"github.com/pkg/errors"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/modules/util"
	"github.com/sitename/sitename/store"
	"gorm.io/gorm"
)

type SqlVoucherTranslationStore struct {
	store.Store
}

func NewSqlVoucherTranslationStore(sqlStore store.Store) store.VoucherTranslationStore {
	return &SqlVoucherTranslationStore{sqlStore}
}

func (s *SqlVoucherTranslationStore) ModelFields(prefix string) util.AnyArray[string] {
	res := util.AnyArray[string]{
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
func (vts *SqlVoucherTranslationStore) Save(translation *model.VoucherTranslation) (*model.VoucherTranslation, error) {
	translation.PreSave()
	if err := translation.IsValid(); err != nil {
		return nil, err
	}

	query := "INSERT INTO " + model.VoucherTranslationTableName + "(" + vts.ModelFields("").Join(",") + ") VALUES (" + vts.ModelFields(":").Join(",") + ")"
	_, err := vts.GetMasterX().NamedExec(query, translation)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to save voucher translation with id=%s", translation.Id)
	}

	return translation, nil
}

// Get finds and returns a voucher translation with given id
func (vts *SqlVoucherTranslationStore) Get(id string) (*model.VoucherTranslation, error) {
	var res model.VoucherTranslation
	err := vts.GetReplicaX().Get(&res, "SELECT * FROM "+model.VoucherTranslationTableName+" WHERE Id = ?", id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, store.NewErrNotFound(model.VoucherTranslationTableName, id)
		}
		return nil, errors.Wrapf(err, "failed to find voucher translation with id=%s", id)
	}

	return &res, nil
}

func (vts *SqlVoucherTranslationStore) commonQueryBuilder(option *model.VoucherTranslationFilterOption) squirrel.SelectBuilder {
	query := vts.GetQueryBuilder().Select("*").From(model.VoucherTranslationTableName)

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
func (vts *SqlVoucherTranslationStore) FilterByOption(option *model.VoucherTranslationFilterOption) ([]*model.VoucherTranslation, error) {
	query := vts.commonQueryBuilder(option)

	queryString, args, err := query.ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "FilterByOption_ToSql")
	}

	var res []*model.VoucherTranslation
	err = vts.GetReplicaX().Select(&res, queryString, args...)
	if err != nil {
		return nil, errors.Wrap(err, "failed to find voucher translations with given options")
	}

	return res, nil
}

// GetByOption finds and returns 1 voucher translation by given options
func (vts *SqlVoucherTranslationStore) GetByOption(option *model.VoucherTranslationFilterOption) (*model.VoucherTranslation, error) {
	query := vts.commonQueryBuilder(option)

	queryString, args, err := query.ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "GetByOption_ToSql")
	}

	var res model.VoucherTranslation
	err = vts.GetReplicaX().Get(&res, queryString, args...)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, store.NewErrNotFound(model.VoucherTranslationTableName, "options")
		}
		return nil, errors.Wrap(err, "failed to find a voucher translation by given option")
	}

	return &res, nil
}
