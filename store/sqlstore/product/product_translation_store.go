package product

import (
	"github.com/pkg/errors"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/store"
	"gorm.io/gorm"
)

type SqlProductTranslationStore struct {
	store.Store
}

func NewSqlProductTranslationStore(s store.Store) store.ProductTranslationStore {
	return &SqlProductTranslationStore{s}
}

// Upsert inserts or update given translation
func (ps *SqlProductTranslationStore) Upsert(translation *model.ProductTranslation) (*model.ProductTranslation, error) {
	err := ps.GetMaster().Save(translation).Error
	if err != nil {
		if ps.IsUniqueConstraintError(err, []string{"LanguageCode", "ProductID", "languagecode_productid_key"}) {
			return nil, store.NewErrInvalidInput(model.ProductTranslationTableName, "LanguageCode/ProductID", "duplicate")
		}
		return nil, errors.Wrap(err, "failed to upsert product translation")
	}

	return translation, nil
}

// Get finds and returns a product translation by given id
func (ps *SqlProductTranslationStore) Get(translationID string) (*model.ProductTranslation, error) {
	var res model.ProductTranslation
	err := ps.GetReplica().First(&res, "Id = ?", translationID).Error
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
	var res []*model.ProductTranslation
	err := ps.GetReplica().Find(&res, store.BuildSqlizer(option.Conditions)...).Error
	if err != nil {
		return nil, errors.Wrap(err, "failed to find product translations with given options")
	}

	return res, nil
}
