package shop

import (
	"database/sql"

	"github.com/pkg/errors"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model/shop"
	"github.com/sitename/sitename/store"
)

type SqlShopTranslationStore struct {
	store.Store
}

func NewSqlShopTranslationStore(s store.Store) store.ShopTranslationStore {
	sts := &SqlShopTranslationStore{s}

	for _, db := range s.GetAllConns() {
		table := db.AddTableWithName(shop.ShopTranslation{}, store.ShopTranslationTableName).SetKeys(false, "Id")
		table.ColMap("Id").SetMaxSize(store.UUID_MAX_LENGTH)
		table.ColMap("ShopID").SetMaxSize(store.UUID_MAX_LENGTH)
		table.ColMap("LanguageCode").SetMaxSize(model.LANGUAGE_CODE_MAX_LENGTH)
		table.ColMap("Name").SetMaxSize(shop.SHOP_TRANSLATION_NAME_MAX_LENGTH)
		table.ColMap("Description").SetMaxSize(shop.SHOP_TRANSLATION_NAME_MAX_LENGTH)

		table.SetUniqueTogether("LanguageCode", "ShopID")
	}
	return sts
}

func (sts *SqlShopTranslationStore) CreateIndexesIfNotExists() {
	sts.CreateForeignKeyIfNotExists(store.ShopTranslationTableName, "ShopID", store.ShopTableName, "Id", true)
}

// Upsert depends on translation's Id then decides to update or insert
func (sts *SqlShopTranslationStore) Upsert(translation *shop.ShopTranslation) (*shop.ShopTranslation, error) {
	var saving bool
	if translation.Id == "" {
		translation.PreSave()
		saving = true
	} else {
		translation.PreUpdate()
	}

	if err := translation.IsValid(); err != nil {
		return nil, err
	}

	var (
		err           error
		numUpdated    int64
		oldTraslation *shop.ShopTranslation
	)
	if saving {
		err = sts.GetMaster().Insert(translation)
	} else {
		oldTraslation, err = sts.Get(translation.Id)
		if err != nil {
			return nil, err
		}
		translation.CreateAt = oldTraslation.CreateAt

		numUpdated, err = sts.GetMaster().Update(translation)
	}

	if err != nil {
		if sts.IsUniqueConstraintError(err, []string{"LanguageCode", "ShopID", "shoptranslations_languagecode_shopid_key"}) {
			return nil, store.NewErrInvalidInput(store.ShopTranslationTableName, "LanguageCode/ShopID", "duplicate value")
		}
		return nil, errors.Wrapf(err, "failed to upsert shop translation with id=%s", translation.Id)
	}

	if numUpdated > 1 {
		return nil, errors.Errorf("multiple shop translations were updated: %d instead of 1", numUpdated)
	}

	return translation, nil
}

// Get finds a shop translation with given id then return it with an error
func (sts *SqlShopTranslationStore) Get(id string) (*shop.ShopTranslation, error) {
	var res shop.ShopTranslation
	err := sts.GetReplica().SelectOne(&res, "SELECT * FROM "+store.ShopTranslationTableName+" WHERE Id = :ID", map[string]interface{}{"ID": id})
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, store.NewErrNotFound(store.ShopTranslationTableName, id)
		}
		return nil, errors.Wrapf(err, "failed to find shop translation with id=%s", id)
	}

	return &res, nil
}
