package product

import (
	"database/sql"

	"github.com/pkg/errors"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model/product_and_discount"
	"github.com/sitename/sitename/store"
)

type SqlProductVariantTranslationStore struct {
	store.Store
}

func NewSqlProductVariantTranslationStore(s store.Store) store.ProductVariantTranslationStore {
	pvts := &SqlProductVariantTranslationStore{s}

	for _, db := range s.GetAllConns() {
		table := db.AddTableWithName(product_and_discount.ProductVariantTranslation{}, store.ProductVariantTranslationTableName).SetKeys(false, "Id")
		table.ColMap("Id").SetMaxSize(store.UUID_MAX_LENGTH)
		table.ColMap("LanguageCode").SetMaxSize(model.LANGUAGE_CODE_MAX_LENGTH)
		table.ColMap("ProductVariantID").SetMaxSize(store.UUID_MAX_LENGTH)
		table.ColMap("Name").SetMaxSize(product_and_discount.PRODUCT_VARIANT_NAME_MAX_LENGTH)

		table.SetUniqueTogether("LanguageCode", "ProductVariantID")
	}
	return pvts
}

func (ps *SqlProductVariantTranslationStore) CreateIndexesIfNotExists() {
	ps.CreateForeignKeyIfNotExists(store.ProductVariantTranslationTableName, "ProductVariantID", store.ProductVariantTableName, "Id", true)
}

// Upsert inserts or updates given translation then returns it
func (ps *SqlProductVariantTranslationStore) Upsert(translation *product_and_discount.ProductVariantTranslation) (*product_and_discount.ProductVariantTranslation, error) {
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
		err = ps.GetMaster().Insert(translation)
	} else {
		_, err = ps.Get(translation.Id)
		if err != nil {
			return nil, err
		}

		numUpdated, err = ps.GetMaster().Update(translation)
	}

	if err != nil {
		if ps.IsUniqueConstraintError(err, []string{"LanguageCode", "ProductVariantID", "productvarianttranslations_languagecode_productvariantid_key", "idx_productvarianttranslations_languagecode_productvariantid_unique"}) {
			return nil, store.NewErrInvalidInput(store.ProductVariantTranslationTableName, "LanguageCode/ProductVariantID", "duplicate")
		}
		return nil, errors.Wrapf(err, "failed to upsert product variant translation with id=%s", translation.Id)
	}

	if numUpdated > 1 {
		return nil, errors.Errorf("multiple product variant translations were updated: %d instead of 1", numUpdated)
	}

	return translation, nil
}

// Get finds and returns 1 product variant translation with given id
func (ps *SqlProductVariantTranslationStore) Get(translationID string) (*product_and_discount.ProductVariantTranslation, error) {
	var res product_and_discount.ProductVariantTranslation
	err := ps.GetReplica().SelectOne(&res, "SELECT * FROM "+store.ProductVariantTranslationTableName+" WHERE Id = :ID", map[string]interface{}{"ID": translationID})
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, store.NewErrNotFound(store.ProductVariantTranslationTableName, translationID)
		}
		return nil, errors.Wrapf(err, "failed to find product variant translation with id=%s", translationID)
	}

	return &res, nil
}

// FilterByOption finds and returns product variant translations filtered using given options
func (ps *SqlProductVariantTranslationStore) FilterByOption(option *product_and_discount.ProductVariantTranslationFilterOption) ([]*product_and_discount.ProductVariantTranslation, error) {
	query := ps.GetQueryBuilder().
		Select("*").
		From(store.ProductVariantTranslationTableName)

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

	var res []*product_and_discount.ProductVariantTranslation
	_, err = ps.GetReplica().Select(&res, queryString, args...)
	if err != nil {
		return nil, errors.Wrap(err, "failed to find product variant translations with given options")
	}

	return res, nil
}
