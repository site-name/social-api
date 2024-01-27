package shop

import (
	"github.com/pkg/errors"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/store"
	"gorm.io/gorm"
)

type SqlShopTranslationStore struct {
	store.Store
}

func NewSqlShopTranslationStore(s store.Store) store.ShopTranslationStore {
	return &SqlShopTranslationStore{s}
}

// Upsert depends on translation's Id then decides to update or insert
func (sts *SqlShopTranslationStore) Upsert(translation model.ShopTranslation) (*model.ShopTranslation, error) {
	err := sts.GetMaster().Save(translation).Error
	if err != nil {
		if sts.IsUniqueConstraintError(err, []string{"LanguageCode", "shoptranslations_languagecode_shopid_key"}) {
			return nil, store.NewErrInvalidInput(model.ShopTranslationTableName, "LanguageCode/ShopID", "duplicate value")
		}
		return nil, errors.Wrapf(err, "failed to upsert shop translation with id=%s", translation.Id)
	}

	return translation, nil
}

// Get finds a shop translation with given id then return it with an error
func (sts *SqlShopTranslationStore) Get(id string) (*model.ShopTranslation, error) {
	var res model.ShopTranslation
	err := sts.GetReplica().First(&res, "Id = ?", id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, store.NewErrNotFound(model.ShopTranslationTableName, id)
		}
		return nil, errors.Wrapf(err, "failed to find shop translation with id=%s", id)
	}

	return &res, nil
}
