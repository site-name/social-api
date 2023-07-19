package product

import (
	"database/sql"

	"github.com/pkg/errors"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/modules/util"
	"github.com/sitename/sitename/store"
	"gorm.io/gorm"
)

type SqlProductTranslationStore struct {
	store.Store
}

func NewSqlProductTranslationStore(s store.Store) store.ProductTranslationStore {
	return &SqlProductTranslationStore{s}
}

func (s *SqlProductTranslationStore) ModelFields(prefix string) util.AnyArray[string] {
	res := util.AnyArray[string]{
		"Id",
		"LanguageCode",
		"ProductID",
		"Name",
		"Description",
		"SeoTitle",
		"SeoDescription",
	}
	if prefix == "" {
		return res
	}

	return res.Map(func(_ int, s string) string {
		return prefix + s
	})
}

// Upsert inserts or update given translation
func (ps *SqlProductTranslationStore) Upsert(translation *model.ProductTranslation) (*model.ProductTranslation, error) {
	var isSaving bool

	if translation.Id == "" {
		translation.PreSave()
		isSaving = true
	} else {
		translation.PreUpdate()
	}

	if err := translation.IsValid(); err != nil {
		return nil, err
	}

	var (
		err        error
		numUpdated int64
	)
	if isSaving {
		query := "INSERT INTO " + model.ProductTranslationTableName + "(" + ps.ModelFields("").Join(",") + ") VALUES (" + ps.ModelFields(":").Join(",") + ")"
		_, err = ps.GetMaster().NamedExec(query, translation)

	} else {
		query := "UPDATE " + model.ProductTranslationTableName + " SET " + ps.
			ModelFields("").
			Map(func(_ int, s string) string {
				return s + "=:" + s
			}).
			Join(",") + " WHERE Id=:Id"

		var result sql.Result
		result, err = ps.GetMaster().NamedExec(query, translation)
		if err == nil && result != nil {
			numUpdated, _ = result.RowsAffected()
		}
	}

	if err != nil {
		if ps.IsUniqueConstraintError(err, []string{"LanguageCode", "ProductID", "producttranslations_languagecode_productid_key"}) {
			return nil, store.NewErrInvalidInput(model.ProductTranslationTableName, "LanguageCode/ProductID", "duplicate")
		}
		return nil, errors.Wrapf(err, "failed to upsert product translation with id=%s", translation.Id)
	}
	if numUpdated > 0 {
		return nil, errors.Errorf("multiple product translations were updated: %d instead of 1", numUpdated)
	}

	return translation, nil
}

// Get finds and returns a product translation by given id
func (ps *SqlProductTranslationStore) Get(translationID string) (*model.ProductTranslation, error) {
	var res model.ProductTranslation
	err := ps.GetReplica().Get(&res, "SELECT * FROM "+model.ProductTranslationTableName+" WHERE Id = ?", translationID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, store.NewErrNotFound(model.ProductTranslationTableName, translationID)
		}
		return nil, errors.Wrapf(err, "failed to find product translation with id=%s", translationID)
	}

	return &res, nil
}

// FilterByOption finds and returns product translations filtered using given options
func (ps *SqlProductTranslationStore) FilterByOption(option *model.ProductTranslationFilterOption) ([]*model.ProductTranslation, error) {
	query := ps.GetQueryBuilder().
		Select(ps.ModelFields(model.ProductTranslationTableName + ".")...).
		From(model.ProductTranslationTableName)

	// parse options
	if option.Id != nil {
		query = query.Where(option.Id)
	}
	if option.LanguageCode != nil {
		query = query.Where(option.LanguageCode)
	}
	if option.ProductID != nil {
		query = query.Where(option.ProductID)
	}
	if option.Name != nil {
		query = query.Where(option.Name)
	}

	queryString, args, err := query.ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "FilterByOption_ToSql")
	}

	var res []*model.ProductTranslation
	err = ps.GetReplica().Select(&res, queryString, args...)
	if err != nil {
		return nil, errors.Wrap(err, "failed to find product translations with given options")
	}

	return res, nil
}
