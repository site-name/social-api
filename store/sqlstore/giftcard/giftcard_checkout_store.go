package giftcard

import (
	"github.com/pkg/errors"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/modules/util"
	"github.com/sitename/sitename/store"
	"gorm.io/gorm"
)

type SqlGiftCardCheckoutStore struct {
	store.Store
}

func NewSqlGiftCardCheckoutStore(s store.Store) store.GiftCardCheckoutStore {
	return &SqlGiftCardCheckoutStore{s}
}

func (s *SqlGiftCardCheckoutStore) ModelFields(prefix string) util.AnyArray[string] {
	res := util.AnyArray[string]{"Id", "GiftcardID", "CheckoutID"}
	if prefix == "" {
		return res
	}

	return res.Map(func(_ int, s string) string {
		return prefix + s
	})
}

func (gs *SqlGiftCardCheckoutStore) Save(giftcardCheckout *model.GiftCardCheckout) (*model.GiftCardCheckout, error) {
	giftcardCheckout.PreSave()
	if err := giftcardCheckout.IsValid(); err != nil {
		return nil, err
	}

	query := "INSERT INTO " + model.GiftcardCheckoutTableName + "(" + gs.ModelFields("").Join(",") + ") VALUES (" + gs.ModelFields(":").Join(",") + ")"
	if _, err := gs.GetMasterX().NamedExec(query, giftcardCheckout); err != nil {
		if gs.IsUniqueConstraintError(err, []string{"GiftcardID", "CheckoutID", "giftcardcheckouts_giftcardid_checkoutid_key"}) {
			return nil, store.NewErrInvalidInput(model.GiftcardCheckoutTableName, "GiftcardID/checkoutID", giftcardCheckout.GiftcardID+"/"+giftcardCheckout.CheckoutID)
		}
		return nil, errors.Wrapf(err, "failed to save giftcard-checkout relation with id=%s", giftcardCheckout.Id)
	}

	return giftcardCheckout, nil
}

func (gs *SqlGiftCardCheckoutStore) Get(id string) (*model.GiftCardCheckout, error) {
	var res model.GiftCardCheckout
	err := gs.GetReplicaX().Get(&res, "SELECT * FROM "+model.GiftcardCheckoutTableName+" WHERE Id = ?", id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, store.NewErrNotFound(model.GiftcardCheckoutTableName, id)
		}
		return nil, errors.Wrapf(err, "failed to get order-checkout with id=%s", id)
	}

	return &res, nil
}

// Delete deletes a giftcard-checkout relation with given id
func (gs *SqlGiftCardCheckoutStore) Delete(giftcardID string, checkoutToken string) error {
	_, err := gs.GetMasterX().Exec("DELETE FROM "+model.GiftcardCheckoutTableName+" WHERE GiftcardID = ? AND CheckoutID = ?", giftcardID, checkoutToken)
	if err != nil {
		return errors.Wrapf(err, "failed to delete giftcard-checkout relation with GiftCardID=%s, CheckoutToken=%s", giftcardID, checkoutToken)
	}

	return nil
}
