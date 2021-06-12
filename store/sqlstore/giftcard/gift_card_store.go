package giftcard

import (
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model/giftcard"
	"github.com/sitename/sitename/store"
)

type SqlGiftCardStore struct {
	store.Store
}

func NewSqlGiftCardStore(sqlStore store.Store) store.GiftCardStore {
	gcs := &SqlGiftCardStore{sqlStore}

	for _, db := range sqlStore.GetAllConns() {
		table := db.AddTableWithName(giftcard.GiftCard{}, "GiftCards").SetKeys(false, "Id")
		table.ColMap("Id").SetMaxSize(store.UUID_MAX_LENGTH)
		table.ColMap("Code").SetMaxSize(giftcard.GIFT_CARD_CODE_MAX_LENGTH).SetUnique(true)
		table.ColMap("Currency").SetMaxSize(model.CURRENCY_CODE_MAX_LENGTH).SetDefaultConstraint(model.NewString(model.DEFAULT_CURRENCY))
	}

	return gcs
}

func (gcs *SqlGiftCardStore) CreateIndexesIfNotExists() {
	gcs.CreateIndexIfNotExists("idx_giftcards_code", "GiftCards", "Code")
}
