package giftcard

import (
	"database/sql"

	"github.com/pkg/errors"
	"github.com/sitename/sitename/model/giftcard"
	"github.com/sitename/sitename/store"
)

type SqlGiftCardCheckoutStore struct {
	store.Store
}

func NewSqlGiftCardCheckoutStore(s store.Store) store.GiftCardCheckoutStore {
	gs := &SqlGiftCardCheckoutStore{s}
	for _, db := range s.GetAllConns() {
		table := db.AddTableWithName(giftcard.GiftCardCheckout{}, store.GiftcardCheckoutTableName).SetKeys(false, "Id")
		table.ColMap("Id").SetMaxSize(store.UUID_MAX_LENGTH)
		table.ColMap("GiftcardID").SetMaxSize(store.UUID_MAX_LENGTH)
		table.ColMap("CheckoutID").SetMaxSize(store.UUID_MAX_LENGTH)

		table.SetUniqueTogether("GiftcardID", "CheckoutID")
	}

	return gs
}

func (gs *SqlGiftCardCheckoutStore) CreateIndexesIfNotExists() {}

func (gs *SqlGiftCardCheckoutStore) Save(giftcardCheckout *giftcard.GiftCardCheckout) (*giftcard.GiftCardCheckout, error) {
	giftcardCheckout.PreSave()
	if err := giftcardCheckout.IsValid(); err != nil {
		return nil, err
	}

	if err := gs.GetMaster().Insert(giftcardCheckout); err != nil {
		if gs.IsUniqueConstraintError(err, []string{"GiftcardID", "OrderID", "giftcardcheckouts_giftcardid_checkoutid_key"}) {
			return nil, store.NewErrInvalidInput(store.GiftcardCheckoutTableName, "GiftcardID/checkoutID", giftcardCheckout.GiftcardID+"/"+giftcardCheckout.CheckoutID)
		}
		return nil, errors.Wrapf(err, "failed to save giftcard-checkout relation with id=%s", giftcardCheckout.Id)
	}

	return giftcardCheckout, nil
}

func (gs *SqlGiftCardCheckoutStore) Get(id string) (*giftcard.GiftCardCheckout, error) {
	res, err := gs.GetReplica().Get(giftcard.GiftCardCheckout{}, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, store.NewErrNotFound(store.GiftcardCheckoutTableName, id)
		}
		return nil, errors.Wrapf(err, "failed to get order-checkout with id=%s", id)
	}

	return res.(*giftcard.GiftCardCheckout), nil
}
