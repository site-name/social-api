package product

import (
	"github.com/pkg/errors"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/store"
	"gorm.io/gorm"
)

type SqlProductVariantTranslationStore struct {
	store.Store
}

func NewSqlProductVariantTranslationStore(s store.Store) store.ProductVariantTranslationStore {
	return &SqlProductVariantTranslationStore{s}
}

// Upsert inserts or updates given translation then returns it
func (ps *SqlProductVariantTranslationStore) Upsert(translation *model.ProductVariantTranslation) (*model.ProductVariantTranslation, error) {
	err := ps.GetMaster().Save(translation).Error
	if err != nil {
		if ps.IsUniqueConstraintError(err, []string{"LanguageCode", "ProductVariantID", "productvarianttranslations_languagecode_productvariantid_key", "idx_productvarianttranslations_languagecode_productvariantid_unique"}) {
			return nil, store.NewErrInvalidInput(model.ProductVariantTranslationTableName, "LanguageCode/ProductVariantID", "duplicate")
		}
		return nil, errors.Wrapf(err, "failed to upsert product variant translation with id=%s", translation.Id)
	}

	return translation, nil
}

// Get finds and returns 1 product variant translation with given id
func (ps *SqlProductVariantTranslationStore) Get(translationID string) (*model.ProductVariantTranslation, error) {
	var res model.ProductVariantTranslation
	err := ps.GetReplica().First(&res, "Id = ?", translationID).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, store.NewErrNotFound(model.ProductVariantTranslationTableName, translationID)
		}
		return nil, errors.Wrapf(err, "failed to find product variant translation with id=%s", translationID)
	}

	return &res, nil
}

// FilterByOption finds and returns product variant translations filtered using given options
func (ps *SqlProductVariantTranslationStore) FilterByOption(option *model.ProductVariantTranslationFilterOption) ([]*model.ProductVariantTranslation, error) {
	args, err := store.BuildSqlizer(option.Conditions, "ProductVariantTranslation_FilterByOption_ToSql")
	if err != nil {
		return nil, err
	}
	var res []*model.ProductVariantTranslation
	err = ps.GetReplica().Find(&res, args...).Error
	if err != nil {
		return nil, errors.Wrap(err, "failed to find product variant translations with given options")
	}

	return res, nil
}
