package shop

import (
	"database/sql"

	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/store"
)

type SqlShopTranslationStore struct {
	store.Store
}

func NewSqlShopTranslationStore(s store.Store) store.ShopTranslationStore {
	return &SqlShopTranslationStore{s}
}

// Upsert depends on translation's Id then decides to update or insert
func (sts *SqlShopTranslationStore) Upsert(translation model.ShopTranslation) (*model.ShopTranslation, error) {
	panic("not implemented")
}

// Get finds a shop translation with given id then return it with an error
func (sts *SqlShopTranslationStore) Get(id string) (*model.ShopTranslation, error) {
	record, err := model.FindShopTranslation(sts.GetReplica(), id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, store.NewErrNotFound(model.TableNames.ShopTranslations, id)
		}
		return nil, err
	}

	return record, nil
}
