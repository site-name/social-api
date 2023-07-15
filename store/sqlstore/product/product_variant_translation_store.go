package product

import (
	"database/sql"

	"github.com/pkg/errors"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/modules/util"
	"github.com/sitename/sitename/store"
)

type SqlProductVariantTranslationStore struct {
	store.Store
}

func NewSqlProductVariantTranslationStore(s store.Store) store.ProductVariantTranslationStore {
	return &SqlProductVariantTranslationStore{s}
}

func (s *SqlProductVariantTranslationStore) ModelFields(prefix string) util.AnyArray[string] {
	res := util.AnyArray[string]{
		"Id",
		"LanguageCode",
		"ProductVariantID",
		"Name",
	}
	if prefix == "" {
		return res
	}

	return res.Map(func(_ int, s string) string {
		return prefix + s
	})
}

// Upsert inserts or updates given translation then returns it
func (ps *SqlProductVariantTranslationStore) Upsert(translation *model.ProductVariantTranslation) (*model.ProductVariantTranslation, error) {
	var isSaving bool

	if !model.IsValidId(translation.Id) {
		isSaving = true
		translation.PreSave()
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
		query := "INSERT INTO " + model.ProductVariantTranslationTableName + "(" + ps.ModelFields("").Join(",") + ") VALUES (" + ps.ModelFields(":").Join(",") + ")"
		_, err = ps.GetMasterX().NamedExec(query, translation)

	} else {
		query := "UPDATE " + model.ProductVariantTranslationTableName + " SET " + ps.
			ModelFields("").
			Map(func(_ int, s string) string {
				return s + "=:" + s
			}).
			Join(",") + " WHERE Id=:Id"

		var result sql.Result
		result, err = ps.GetMasterX().NamedExec(query, translation)
		if err == nil && result != nil {
			numUpdated, _ = result.RowsAffected()
		}
	}

	if err != nil {
		if ps.IsUniqueConstraintError(err, []string{"LanguageCode", "ProductVariantID", "productvarianttranslations_languagecode_productvariantid_key", "idx_productvarianttranslations_languagecode_productvariantid_unique"}) {
			return nil, store.NewErrInvalidInput(model.ProductVariantTranslationTableName, "LanguageCode/ProductVariantID", "duplicate")
		}
		return nil, errors.Wrapf(err, "failed to upsert product variant translation with id=%s", translation.Id)
	}

	if numUpdated > 1 {
		return nil, errors.Errorf("multiple product variant translations were updated: %d instead of 1", numUpdated)
	}

	return translation, nil
}

// Get finds and returns 1 product variant translation with given id
func (ps *SqlProductVariantTranslationStore) Get(translationID string) (*model.ProductVariantTranslation, error) {
	var res model.ProductVariantTranslation
	err := ps.GetReplicaX().Get(&res, "SELECT * FROM "+model.ProductVariantTranslationTableName+" WHERE Id = ?", translationID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, store.NewErrNotFound(model.ProductVariantTranslationTableName, translationID)
		}
		return nil, errors.Wrapf(err, "failed to find product variant translation with id=%s", translationID)
	}

	return &res, nil
}

// FilterByOption finds and returns product variant translations filtered using given options
func (ps *SqlProductVariantTranslationStore) FilterByOption(option *model.ProductVariantTranslationFilterOption) ([]*model.ProductVariantTranslation, error) {
	query := ps.GetQueryBuilder().
		Select("*").
		From(model.ProductVariantTranslationTableName)

	// parse options
	if option.Id != nil {
		query = query.Where(option.Id)
	}
	if option.LanguageCode != nil {
		query = query.Where(option.LanguageCode)
	}
	if option.ProductVariantID != nil {
		query = query.Where(option.ProductVariantID)
	}
	if option.Name != nil {
		query = query.Where(option.Name)
	}

	queryString, args, err := query.ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "FilterByOption_ToSql")
	}

	var res []*model.ProductVariantTranslation
	err = ps.GetReplicaX().Select(&res, queryString, args...)
	if err != nil {
		return nil, errors.Wrap(err, "failed to find product variant translations with given options")
	}

	return res, nil
}
