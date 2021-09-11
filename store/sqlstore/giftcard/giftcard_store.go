package giftcard

import (
	"database/sql"

	"github.com/Masterminds/squirrel"
	"github.com/mattermost/gorp"
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
		table.ColMap("CreatedByID").SetMaxSize(store.UUID_MAX_LENGTH)
		table.ColMap("UsedByID").SetMaxSize(store.UUID_MAX_LENGTH)
		table.ColMap("ProductID").SetMaxSize(store.UUID_MAX_LENGTH)
		table.ColMap("CreatedByEmail").SetMaxSize(model.USER_EMAIL_MAX_LENGTH)
		table.ColMap("UsedByEmail").SetMaxSize(model.USER_EMAIL_MAX_LENGTH)
		table.ColMap("ExpiryType").SetMaxSize(giftcard.GiftcardExpiryTypeMaxLength)
		table.ColMap("ExpiryPeriodType").SetMaxSize(giftcard.GiftcardExpiryPeriodTypeMaxLength)
		table.ColMap("Tag").SetMaxSize(giftcard.GiftcardTagMaxLength)
		table.ColMap("Code").SetMaxSize(giftcard.GiftcardCodeMaxLength).SetUnique(true)
		table.ColMap("Currency").SetMaxSize(model.CURRENCY_CODE_MAX_LENGTH)
	}

	return gcs
}

func (gcs *SqlGiftCardStore) CreateIndexesIfNotExists() {
	gcs.CommonMetaDataIndex(store.GiftcardTableName)

	gcs.CreateIndexIfNotExists("idx_giftcards_tag", store.GiftcardTableName, "Tag")
	gcs.CreateIndexIfNotExists("idx_giftcards_code", store.GiftcardTableName, "Code")
	gcs.CreateForeignKeyIfNotExists(store.GiftcardTableName, "CreatedByID", store.UserTableName, "Id", false)
	gcs.CreateForeignKeyIfNotExists(store.GiftcardTableName, "UsedByID", store.UserTableName, "Id", false)
	gcs.CreateForeignKeyIfNotExists(store.GiftcardTableName, "ProductID", store.ProductTableName, "Id", false)
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
	var res giftcard.GiftCard
	if err := gcs.GetReplica().SelectOne(&res, "SELECT * FROM "+store.GiftcardTableName+" WHERE Id = :ID", map[string]interface{}{"ID": id}); err != nil {
		if err == sql.ErrNoRows {
			return nil, store.NewErrNotFound(store.GiftcardTableName, id)
		}
		return nil, errors.Wrapf(err, "failed to find giftcard with id=%s", id)
	} else {
		return &res, nil
	}
}

// FilterByOption finds giftcards wth option
func (gs *SqlGiftCardStore) FilterByOption(transaction *gorp.Transaction, option *giftcard.GiftCardFilterOption) ([]*giftcard.GiftCard, error) {
	var selector store.Selector = gs.GetReplica()
	if transaction != nil {
		selector = transaction
	}

	query := gs.
		GetQueryBuilder().
		Select("*").
		From(store.GiftcardTableName).
		OrderBy(store.TableOrderingMap[store.GiftcardTableName])

	// check code
	if option.Distinct {
		query = query.Distinct()
	}
	if option.CreatedByID != nil {
		query = query.Where(option.CreatedByID.ToSquirrel("CreatedByID"))
	}
	if option.Code != nil {
		query = query.Where(option.Code.ToSquirrel("Code"))
	}
	if option.ExpiryDate != nil {
		query = query.Where(option.ExpiryDate.ToSquirrel("ExpiryDate"))
	}
	if option.StartDate != nil {
		query = query.Where(option.StartDate.ToSquirrel("StartDate"))
	}
	if option.CheckoutToken != nil {
		subSelect := gs.GetQueryBuilder().
			Select("GiftcardID").
			From(store.GiftcardCheckoutTableName).
			Where(option.CheckoutToken.ToSquirrel("GiftcardCheckouts.CheckoutID"))

		query = query.Where(squirrel.Expr("Id IN ?", subSelect))
	}
	if option.IsActive != nil {
		query = query.Where(squirrel.Eq{"IsActive": *option.IsActive})
	}
	if option.SelectForUpdate {
		query = query.Suffix("FOR UPDATE")
	}

	queryString, args, err := query.ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "query_toSql")
	}

	var giftcards []*giftcard.GiftCard
	_, err = selector.Select(&giftcards, queryString, args...)
	if err != nil {
		return nil, errors.Wrap(err, "failed to finds giftcards with code")
	}

	return giftcards, nil
}
