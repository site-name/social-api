package shop

import (
	"database/sql"

	"github.com/Masterminds/squirrel"
	"github.com/pkg/errors"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/modules/util"
	"github.com/sitename/sitename/store"
	"gorm.io/gorm"
)

type sqlVatStore struct {
	store.Store
}

func NewSqlVatStore(s store.Store) store.VatStore {
	return &sqlVatStore{s}
}

func (s *sqlVatStore) ModelFields(prefix string) util.AnyArray[string] {
	res := util.AnyArray[string]{
		"Id", "CountryCode", "Data",
	}
	if prefix == "" {
		return res
	}
	return res.Map(func(_ int, item string) string {
		return prefix + "item"
	})
}

func (s *sqlVatStore) Upsert(transaction *gorm.DB, vats []*model.Vat) ([]*model.Vat, error) {
	runner := s.GetMaster()
	if transaction != nil {
		runner = transaction
	}

	saveQuery := "INSERT INTO " + model.VatTableName + "(Id, CountryCode, Data) VALUES (" + s.ModelFields(":").Join(",") + ")"
	updateQuery := "UPDATE " + model.VatTableName + " SET " + s.ModelFields("").
		Map(func(_ int, item string) string {
			return item + ":=" + item
		}).
		Join(",") + " WHERE Id=:Id"

	for _, vat := range vats {
		isSaving := false

		if !model.IsValidId(vat.Id) {
			vat.Id = ""
			isSaving = true
			vat.PreSave()
		} else {
			vat.PreUpdate()
		}

		if err := vat.IsValid(); err != nil {
			return nil, err
		}

		var err error
		var result sql.Result

		if isSaving {
			result, err = runner.Exec(saveQuery, vat)
		} else {
			result, err = runner.Exec(updateQuery, vat)
		}

		if err != nil {
			return nil, errors.Wrap(err, "failed to upsert a vat")
		}
		numUpserted, _ := result.RowsAffected()
		if numUpserted != 1 {
			return nil, errors.Errorf("$d vat object(s) upserted instead of 1", numUpserted)
		}
	}

	return vats, nil
}

func (s *sqlVatStore) FilterByOptions(options *model.VatFilterOptions) ([]*model.Vat, error) {
	query := s.
		GetQueryBuilder().
		Select(s.ModelFields(model.VatTableName + ".")...).
		From(model.VatTableName)

	if options == nil {
		options = new(model.VatFilterOptions)
	}

	for _, opt := range []squirrel.Sqlizer{options.Id, options.CountryCode} {
		if opt != nil {
			query = query.Where(opt)
		}
	}

	queryString, args, err := query.ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "FilterbyOptions_ToSql")
	}

	var res []*model.Vat
	err = s.GetReplica().Select(&res, queryString, args...)
	if err != nil {
		return nil, errors.Wrap(err, "failed to find vat objects by options")
	}

	return res, nil
}
