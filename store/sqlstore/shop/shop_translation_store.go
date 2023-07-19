package shop

import (
	"database/sql"

	"github.com/pkg/errors"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/modules/util"
	"github.com/sitename/sitename/store"
	"gorm.io/gorm"
)

type SqlShopTranslationStore struct {
	store.Store
}

func NewSqlShopTranslationStore(s store.Store) store.ShopTranslationStore {
	return &SqlShopTranslationStore{s}
}

func (s *SqlShopTranslationStore) ModelFields(prefix string) util.AnyArray[string] {
	res := util.AnyArray[string]{
		"Id",
		"LanguageCode",
		"Name",
		"Description",
		"CreateAt",
		"UpdateAt",
	}
	if prefix == "" {
		return res
	}

	return res.Map(func(_ int, s string) string {
		return prefix + s
	})
}

// Upsert depends on translation's Id then decides to update or insert
func (sts *SqlShopTranslationStore) Upsert(translation *model.ShopTranslation) (*model.ShopTranslation, error) {
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
		err        error
		numUpdated int64
	)
	if saving {
		query := "INSERT INTO " + model.ShopTranslationTableName + "(" + sts.ModelFields("").Join(",") + ") VALUES (" + sts.ModelFields(":").Join(",") + ")"
		_, err = sts.GetMaster().NamedExec(query, translation)

	} else {
		query := "UPDATE " + model.ShopTranslationTableName + " SET " + sts.
			ModelFields("").
			Map(func(_ int, s string) string {
				return s + "=:" + s
			}).
			Join(",") + " WHERE Id=:Id"

		var result sql.Result
		result, err = sts.GetMaster().NamedExec(query, translation)
		if err == nil && result != nil {
			numUpdated, _ = result.RowsAffected()
		}
	}

	if err != nil {
		if sts.IsUniqueConstraintError(err, []string{"LanguageCode", "shoptranslations_languagecode_shopid_key"}) {
			return nil, store.NewErrInvalidInput(model.ShopTranslationTableName, "LanguageCode/ShopID", "duplicate value")
		}
		return nil, errors.Wrapf(err, "failed to upsert shop translation with id=%s", translation.Id)
	}

	if numUpdated > 1 {
		return nil, errors.Errorf("multiple shop translations were updated: %d instead of 1", numUpdated)
	}

	return translation, nil
}

// Get finds a shop translation with given id then return it with an error
func (sts *SqlShopTranslationStore) Get(id string) (*model.ShopTranslation, error) {
	var res model.ShopTranslation
	err := sts.GetReplica().Get(&res, "SELECT * FROM "+model.ShopTranslationTableName+" WHERE Id = ?", id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, store.NewErrNotFound(model.ShopTranslationTableName, id)
		}
		return nil, errors.Wrapf(err, "failed to find shop translation with id=%s", id)
	}

	return &res, nil
}
