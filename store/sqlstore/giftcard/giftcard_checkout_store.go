package giftcard

import (
	"database/sql"

	"github.com/pkg/errors"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model/giftcard"
	"github.com/sitename/sitename/store"
)

type SqlGiftCardCheckoutStore struct {
	store.Store
}

func NewSqlGiftCardCheckoutStore(s store.Store) store.GiftCardCheckoutStore {
	return &SqlGiftCardCheckoutStore{s}
}

func (s *SqlGiftCardCheckoutStore) ModelFields(prefix string) model.StringArray {
	res := model.StringArray{"Id", "GiftcardID", "CheckoutID"}
	if prefix == "" {
		return res
	}

	return res.Map(func(_ int, s string) string {
		return prefix + s
	})
}

func (gs *SqlGiftCardCheckoutStore) Save(giftcardCheckout *giftcard.GiftCardCheckout) (*giftcard.GiftCardCheckout, error) {
	giftcardCheckout.PreSave()
	if err := giftcardCheckout.IsValid(); err != nil {
		return nil, err
	}

	query := "INSERT INTO " + store.GiftcardCheckoutTableName + "(" + gs.ModelFields("").Join(",") + ") VALUES (" + gs.ModelFields(":").Join(",") + ")"
	if _, err := gs.GetMasterX().NamedExec(query, giftcardCheckout); err != nil {
		if gs.IsUniqueConstraintError(err, []string{"GiftcardID", "CheckoutID", "giftcardcheckouts_giftcardid_checkoutid_key"}) {
			return nil, store.NewErrInvalidInput(store.GiftcardCheckoutTableName, "GiftcardID/checkoutID", giftcardCheckout.GiftcardID+"/"+giftcardCheckout.CheckoutID)
		}
		return nil, errors.Wrapf(err, "failed to save giftcard-checkout relation with id=%s", giftcardCheckout.Id)
	}

	return giftcardCheckout, nil
}

func (gs *SqlGiftCardCheckoutStore) Get(id string) (*giftcard.GiftCardCheckout, error) {
	var res giftcard.GiftCardCheckout
	err := gs.GetReplicaX().Get(&res, "SELECT * FROM "+store.GiftcardCheckoutTableName+" WHERE Id = ?", id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, store.NewErrNotFound(store.GiftcardCheckoutTableName, id)
		}
		return nil, errors.Wrapf(err, "failed to get order-checkout with id=%s", id)
	}

	return &res, nil
}

// Delete deletes a giftcard-checkout relation with given id
func (gs *SqlGiftCardCheckoutStore) Delete(giftcardID string, checkoutToken string) error {
	_, err := gs.GetMasterX().Exec("DELETE FROM "+store.GiftcardCheckoutTableName+" WHERE GiftcardID = ? AND CheckoutID = ?", giftcardID, checkoutToken)
	if err != nil {
		return errors.Wrapf(err, "failed to delete giftcard-checkout relation with GiftCardID=%s, CheckoutToken=%s", giftcardID, checkoutToken)
	}

	return nil
}
