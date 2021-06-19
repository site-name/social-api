package giftcard

import (
	"database/sql"

	"github.com/pkg/errors"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model/giftcard"
	"github.com/sitename/sitename/store"
)

type SqlGiftCardStore struct {
	store.Store
}

const (
	giftcardTableName = "GiftCards"
)

func NewSqlGiftCardStore(sqlStore store.Store) store.GiftCardStore {
	gcs := &SqlGiftCardStore{sqlStore}

	for _, db := range sqlStore.GetAllConns() {
		table := db.AddTableWithName(giftcard.GiftCard{}, giftcardTableName).SetKeys(false, "Id")
		table.ColMap("Id").SetMaxSize(store.UUID_MAX_LENGTH)
		table.ColMap("UserID").SetMaxSize(store.UUID_MAX_LENGTH)
		table.ColMap("Code").SetMaxSize(giftcard.GIFT_CARD_CODE_MAX_LENGTH).SetUnique(true)
		table.ColMap("Currency").SetMaxSize(model.CURRENCY_CODE_MAX_LENGTH).SetDefaultConstraint(model.NewString(model.DEFAULT_CURRENCY))
	}

	return gcs
}

func (gcs *SqlGiftCardStore) CreateIndexesIfNotExists() {
	gcs.CreateIndexIfNotExists("idx_giftcards_code", giftcardTableName, "Code")
}

func (gcs *SqlGiftCardStore) Save(giftCard *giftcard.GiftCard) (*giftcard.GiftCard, error) {
	giftCard.PreSave()
	if err := giftCard.IsValid(); err != nil {
		return nil, err
	}
	if err := gcs.GetMaster().Insert(giftCard); err != nil {
		if gcs.IsUniqueConstraintError(err, []string{"Code", "giftcards_code_key", "idx_giftcards_code_unique"}) {
			return nil, store.NewErrInvalidInput(giftcardTableName, "Code", giftCard.Code)
		}
		return nil, errors.Wrapf(err, "failed to save giftcard with id=%s", giftCard.Id)
	}

	return giftCard, nil
}

func (gcs *SqlGiftCardStore) GetById(id string) (*giftcard.GiftCard, error) {
	if res, err := gcs.GetReplica().Get(giftcard.GiftCard{}, id); err != nil {
		if err == sql.ErrNoRows {
			return nil, store.NewErrNotFound(giftcardTableName, id)
		}
		return nil, errors.Wrapf(err, "failed to find giftcard with id=%s", id)
	} else {
		return res.(*giftcard.GiftCard), nil
	}
}

func (gcs *SqlGiftCardStore) GetAllByUserId(userID string) ([]*giftcard.GiftCard, error) {
	var giftcards []*giftcard.GiftCard
	if _, err := gcs.GetReplica().Select(&giftcards, "SELECT * FROM "+giftcardTableName+" WHERE UserID = :userID",
		map[string]interface{}{"userID": userID}); err != nil {
		if err != nil {
			if err == sql.ErrNoRows {
				return nil, store.NewErrNotFound(giftcardTableName, "userID="+userID)
			}
			return nil, errors.Wrapf(err, "failed to find giftcards with userID=%s", userID)
		}
	}
	return giftcards, nil
}
