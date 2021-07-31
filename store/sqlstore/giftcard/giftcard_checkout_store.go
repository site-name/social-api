package giftcard

import (
	"database/sql"
	"fmt"

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

func (gs *SqlGiftCardCheckoutStore) CreateIndexesIfNotExists() {
	gs.CreateForeignKeyIfNotExists(store.GiftcardCheckoutTableName, "GiftcardID", store.GiftcardTableName, "Id", false)
	gs.CreateForeignKeyIfNotExists(store.GiftcardCheckoutTableName, "CheckoutID", store.CheckoutTableName, "Token", false)
}

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

// Delete deletes a giftcard-checkout relation with given id
func (gs *SqlGiftCardCheckoutStore) Delete(giftcardID string, checkoutToken string) error {
	var oldRelation *giftcard.GiftCardCheckout
	err := gs.GetReplica().SelectOne(
		&oldRelation,
		`SELECT * FROM `+store.GiftcardCheckoutTableName+`
		WHERE (
			GiftcardID = :GiftCardID AND CheckoutID = :CheckoutToken
		)`,
		map[string]interface{}{
			"GiftCardID": giftcardID,
			"CheckoutID": checkoutToken,
		},
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return store.NewErrNotFound(store.GiftcardCheckoutTableName, fmt.Sprintf("GiftCardID=%s, CheckoutID=%s", giftcardID, checkoutToken))
		}
		return errors.Errorf("failed to delete giftcard-checkout relation with GiftCardID=%s, CheckoutID=%s", giftcardID, checkoutToken)
	}

	numDeleted, err := gs.GetMaster().Delete(oldRelation)
	if err != nil {
		return errors.Wrapf(err, "failed to delete giftcard-checkout relation with GiftCardID=%s, CheckoutToken=%s", giftcardID, checkoutToken)
	}

	if numDeleted > 1 {
		return errors.Errorf("multiple giftcard-checkout relations were deleted: %d instead of 1", numDeleted)
	}

	return nil
}
