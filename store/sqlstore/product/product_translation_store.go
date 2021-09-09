package product

import (
	"database/sql"

	"github.com/pkg/errors"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model/product_and_discount"
	"github.com/sitename/sitename/store"
)

type SqlProductTranslationStore struct {
	store.Store
}

func NewSqlProductTranslationStore(s store.Store) store.ProductTranslationStore {
	pts := &SqlProductTranslationStore{s}

	for _, db := range s.GetAllConns() {
		table := db.AddTableWithName(product_and_discount.ProductTranslation{}, store.ProductTranslationTableName).SetKeys(false, "Id")
		table.ColMap("Id").SetMaxSize(store.UUID_MAX_LENGTH)
		table.ColMap("ProductID").SetMaxSize(store.UUID_MAX_LENGTH)
		table.ColMap("LanguageCode").SetMaxSize(model.LANGUAGE_CODE_MAX_LENGTH)
		table.ColMap("Name").SetMaxSize(product_and_discount.PRODUCT_NAME_MAX_LENGTH).SetUnique(true)

		table.SetUniqueTogether("LanguageCode", "ProductID")
		s.CommonSeoMaxLength(table)
	}
	return pts
}

func (ps *SqlProductTranslationStore) CreateIndexesIfNotExists() {
	ps.CreateForeignKeyIfNotExists(store.ProductTranslationTableName, "ProductID", store.ProductTableName, "Id", true)
}

// Upsert inserts or update given translation
func (ps *SqlProductTranslationStore) Upsert(translation *product_and_discount.ProductTranslation) (*product_and_discount.ProductTranslation, error) {
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
		err = ps.GetMaster().Insert()
	} else {
		_, err = ps.Get(translation.Id)
		if err != nil {
			return nil, err
		}

		numUpdated, err = ps.GetMaster().Update(translation)
	}

	if err != nil {
		if ps.IsUniqueConstraintError(err, []string{"LanguageCode", "ProductID", "producttranslations_languagecode_productid_key"}) {
			return nil, store.NewErrInvalidInput(store.ProductTranslationTableName, "LanguageCode/ProductID", "duplicate")
		}
		return nil, errors.Wrapf(err, "failed to upsert product translation with id=%s", translation.Id)
	}
	if numUpdated > 0 {
		return nil, errors.Errorf("multiple product translations were updated: %d instead of 1", numUpdated)
	}

	return translation, nil
}

// Get finds and returns a product translation by given id
func (ps *SqlProductTranslationStore) Get(translationID string) (*product_and_discount.ProductTranslation, error) {
	var res product_and_discount.ProductTranslation
	err := ps.GetReplica().SelectOne(&res, "SELECT * FROM "+store.ProductTranslationTableName+" WHERE Id = :ID", map[string]interface{}{"ID": translationID})
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, store.NewErrNotFound(store.ProductTranslationTableName, translationID)
		}
		return nil, errors.Wrapf(err, "failed to find product translation with id=%s", translationID)
	}

	return &res, nil
}

// FilterByOption finds and returns product translations filtered using given options
func (ps *SqlProductTranslationStore) FilterByOption(option *product_and_discount.ProductTranslationFilterOption) ([]*product_and_discount.ProductTranslation, error) {
	query := ps.GetQueryBuilder().
		Select("*").
		From(store.ProductTranslationTableName).
		OrderBy(store.TableOrderingMap[store.ProductTranslationTableName])

	// parse options
	if option.Id != nil {
		query = query.Where(option.Id.ToSquirrel("Id"))
	}
	if option.LanguageCode != nil {
		query = query.Where(option.LanguageCode.ToSquirrel("LanguageCode"))
	}
	if option.ProductID != nil {
		query = query.Where(option.ProductID.ToSquirrel("ProductID"))
	}
	if option.Name != nil {
		query = query.Where(option.Name.ToSquirrel("Name"))
	}

	queryString, args, err := query.ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "FilterByOption_ToSql")
	}

	var res []*product_and_discount.ProductTranslation
	_, err = ps.GetReplica().Select(&res, queryString, args...)
	if err != nil {
		return nil, errors.Wrap(err, "failed to find product translations with given options")
	}

	return res, nil
}
