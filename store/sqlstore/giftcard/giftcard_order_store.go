package giftcard

import (
	"database/sql"

	"github.com/pkg/errors"
	"github.com/sitename/sitename/model/giftcard"
	"github.com/sitename/sitename/store"
)

type SqlGiftCardOrderStore struct {
	store.Store
}

func NewSqlGiftCardOrderStore(s store.Store) store.GiftCardOrderStore {
	gs := &SqlGiftCardOrderStore{s}
	for _, db := range s.GetAllConns() {
		table := db.AddTableWithName(giftcard.OrderGiftCard{}, store.OrderGiftCardTableName).SetKeys(false, "Id")
		table.ColMap("Id").SetMaxSize(store.UUID_MAX_LENGTH)
		table.ColMap("GiftCardID").SetMaxSize(store.UUID_MAX_LENGTH)
		table.ColMap("OrderID").SetMaxSize(store.UUID_MAX_LENGTH)

		table.SetUniqueTogether("GiftCardID", "OrderID")
	}

	return gs
}

func (gs *SqlGiftCardOrderStore) CreateIndexesIfNotExists() {
	gs.CreateForeignKeyIfNotExists(store.OrderGiftCardTableName, "GiftCardID", store.GiftcardTableName, "id", false)
	gs.CreateForeignKeyIfNotExists(store.OrderGiftCardTableName, "OrderID", store.OrderTableName, "Id", false)
}

func (gs *SqlGiftCardOrderStore) Save(giftCardOrder *giftcard.OrderGiftCard) (*giftcard.OrderGiftCard, error) {
	giftCardOrder.PreSave()
	if err := giftCardOrder.IsValid(); err != nil {
		return nil, err
	}

	if err := gs.GetMaster().Insert(giftCardOrder); err != nil {
		if gs.IsUniqueConstraintError(err, []string{"GiftCardID", "OrderID", "ordergiftcards_giftcardid_orderid_key"}) {
			return nil, store.NewErrInvalidInput(store.OrderGiftCardTableName, "GiftCardID/OrderID", giftCardOrder.GiftCardID+"/"+giftCardOrder.OrderID)
		}
		return nil, errors.Wrapf(err, "failed to save giftcard-order relation with id=%s", giftCardOrder.Id)
	}

	return giftCardOrder, nil
}

func (gs *SqlGiftCardOrderStore) Get(id string) (*giftcard.OrderGiftCard, error) {
	res, err := gs.GetReplica().Get(giftcard.OrderGiftCard{}, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, store.NewErrNotFound(store.OrderGiftCardTableName, id)
		}
		return nil, errors.Wrapf(err, "failed to get order-giftcard with id=%s", id)
	}

	return res.(*giftcard.OrderGiftCard), nil
}
