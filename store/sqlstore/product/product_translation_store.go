package product

import (
	"github.com/sitename/sitename/store"
)

type SqlProductTranslationStore struct {
	store.Store
}

func NewSqlProductTranslationStore(s store.Store) store.ProductTranslationStore {
	return &SqlProductTranslationStore{s}
}

// func (ps *SqlProductTranslationStore) Upsert(translation *model.ProductTranslation) (*model.ProductTranslation, error) {
// 	err := ps.GetMaster().Save(translation).Error
// 	if err != nil {
// 		if ps.IsUniqueConstraintError(err, []string{"LanguageCode", "ProductID", "languagecode_productid_key"}) {
// 			return nil, store.NewErrInvalidInput(model.ProductTranslationTableName, "LanguageCode/ProductID", "duplicate")
// 		}
// 		return nil, errors.Wrap(err, "failed to upsert product translation")
// 	}

// 	return translation, nil
// }

// func (ps *SqlProductTranslationStore) Get(translationID string) (*model.ProductTranslation, error) {
// 	var res model.ProductTranslation
// 	err := ps.GetReplica().First(&res, "Id = ?", translationID).Error
// 	if err != nil {
// 		if errors.Is(err, gorm.ErrRecordNotFound) {
// 			return nil, store.NewErrNotFound(model.ProductTranslationTableName, translationID)
// 		}
// 		return nil, errors.Wrapf(err, "failed to find product translation with id=%s", translationID)
// 	}

// 	return &res, nil
// }

// func (ps *SqlProductTranslationStore) FilterByOption(option *model.ProductTranslationFilterOption) ([]*model.ProductTranslation, error) {
// 	args, err := store.BuildSqlizer(option.Conditions, "ProductTranslation_FilterByOption")
// 	if err != nil {
// 		return nil, err
// 	}
// 	var res []*model.ProductTranslation
// 	err = ps.GetReplica().Find(&res, args...).Error
// 	if err != nil {
// 		return nil, errors.Wrap(err, "failed to find product translations with given options")
// 	}

// 	return res, nil
// }
