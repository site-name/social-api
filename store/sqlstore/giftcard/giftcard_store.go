package giftcard

import (
	"database/sql"

	"github.com/Masterminds/squirrel"
	"github.com/pkg/errors"
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
		table := db.AddTableWithName(giftcard.GiftCard{}, store.GiftcardTableName).SetKeys(false, "Id")
		table.ColMap("Id").SetMaxSize(store.UUID_MAX_LENGTH)
		table.ColMap("UserID").SetMaxSize(store.UUID_MAX_LENGTH)
		table.ColMap("Code").SetMaxSize(giftcard.GIFT_CARD_CODE_MAX_LENGTH).SetUnique(true)
		table.ColMap("Currency").SetMaxSize(model.CURRENCY_CODE_MAX_LENGTH)
	}

	return gcs
}

func (gcs *SqlGiftCardStore) CreateIndexesIfNotExists() {
	gcs.CreateIndexIfNotExists("idx_giftcards_code", store.GiftcardTableName, "Code")
	gcs.CreateForeignKeyIfNotExists(store.GiftcardTableName, "UserID", store.UserTableName, "Id", false)
}

// Upsert depends on given giftcard's Id property then perform according operation
func (gcs *SqlGiftCardStore) Upsert(giftCard *giftcard.GiftCard) (*giftcard.GiftCard, error) {
	var saving bool
	if giftCard.Id == "" {
		giftCard.PreSave()
		saving = true
	} else {
		giftCard.PreUpdate()
	}

	if err := giftCard.IsValid(); err != nil {
		return nil, err
	}

	var (
		err         error
		oldGiftcard *giftcard.GiftCard
		numUpdated  int64
	)
	if saving {
		err = gcs.GetMaster().Insert(giftCard)
	} else {
		oldGiftcard, err = gcs.GetById(giftCard.Id)
		if err != nil {
			return nil, err
		}

		giftCard.CreateAt = oldGiftcard.CreateAt
		giftCard.Code = oldGiftcard.Code

		numUpdated, err = gcs.GetMaster().Update(giftCard)
	}

	if err != nil {
		if gcs.IsUniqueConstraintError(err, []string{"Code", "giftcards_code_key", "idx_giftcards_code_unique"}) {
			return nil, store.NewErrInvalidInput(store.GiftcardTableName, "Code", giftCard.Code)
		}
		return nil, errors.Wrapf(err, "failed to upsert giftcard with id=%s", giftCard.Id)
	}

	if numUpdated > 1 {
		return nil, errors.Errorf("multiple giftcards were updated: %d instead of 1", numUpdated)
	}

	return giftCard, nil
}

func (gcs *SqlGiftCardStore) GetById(id string) (*giftcard.GiftCard, error) {
	if res, err := gcs.GetReplica().Get(giftcard.GiftCard{}, id); err != nil {
		if err == sql.ErrNoRows {
			return nil, store.NewErrNotFound(store.GiftcardTableName, id)
		}
		return nil, errors.Wrapf(err, "failed to find giftcard with id=%s", id)
	} else {
		return res.(*giftcard.GiftCard), nil
	}
}

func (gcs *SqlGiftCardStore) GetAllByUserId(userID string) ([]*giftcard.GiftCard, error) {
	var giftcards []*giftcard.GiftCard
	if _, err := gcs.GetReplica().Select(&giftcards, "SELECT * FROM "+store.GiftcardTableName+" WHERE UserID = :userID",
		map[string]interface{}{"userID": userID}); err != nil {
		if err != nil {
			if err == sql.ErrNoRows {
				return nil, store.NewErrNotFound(store.GiftcardTableName, "userID="+userID)
			}
			return nil, errors.Wrapf(err, "failed to find giftcards with userID=%s", userID)
		}
	}
	return giftcards, nil
}

func (gs *SqlGiftCardStore) GetAllByCheckout(checkoutID string) ([]*giftcard.GiftCard, error) {
	query := `SELECT * FROM ` + store.GiftcardTableName + ` AS Gc
		WHERE Gc.Id IN (
			SELECT GcCk.GiftcardID FROM ` + store.GiftcardCheckoutTableName + ` AS GcCk
		)
		WHERE GcCk.CheckoutID = :CheckoutID`

	var giftcards []*giftcard.GiftCard
	_, err := gs.GetReplica().Select(&giftcards, query, map[string]interface{}{"CheckoutID": checkoutID})
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, store.NewErrNotFound(store.GiftcardTableName, "checkoutID="+checkoutID)
		}
		return nil, errors.Wrapf(err, "failed to find giftcards belong to checkout with id=%s", checkoutID)
	}

	return giftcards, nil
}

// FilterByOption finds giftcards wth option
func (gs *SqlGiftCardStore) FilterByOption(option *giftcard.GiftCardFilterOption) ([]*giftcard.GiftCard, error) {

	query := gs.
		GetQueryBuilder().
		Select(store.GiftcardTableName).
		OrderBy("CreateAt ASC")

	// check code
	if option.Code != nil {
		query = query.Where(option.Code.ToSquirrel("Code"))
	}

	// check end date
	if option.EndDate != nil {
		query = query.Where(option.EndDate.ToSquirrel("EndDate"))
	}

	// check start date
	if option.StartDate != nil {
		query = query.Where(option.StartDate.ToSquirrel("StartDate"))
	}

	// check is active
	if option.IsActive != nil {
		query = query.Where(squirrel.Eq{"IsActive": *option.IsActive})
	}

	queryString, args, err := query.ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "query_toSql")
	}

	var giftcards []*giftcard.GiftCard
	_, err = gs.GetReplica().Select(&giftcards, queryString, args...)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, store.NewErrNotFound(store.GiftcardTableName, "code")
		}
		return nil, errors.Wrap(err, "failed to finds giftcards with code")
	}

	return giftcards, nil
}
