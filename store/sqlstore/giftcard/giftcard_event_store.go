package giftcard

import (
	"database/sql"

	"github.com/mattermost/gorp"
	"github.com/pkg/errors"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model/giftcard"
	"github.com/sitename/sitename/store"
)

type SqlGiftcardEventStore struct {
	store.Store
}

func NewSqlGiftcardEventStore(s store.Store) store.GiftcardEventStore {
	gcs := &SqlGiftcardEventStore{s}

	for _, db := range s.GetAllConns() {
		table := db.AddTableWithName(giftcard.GiftCardEvent{}, store.GiftcardEventTableName).SetKeys(false, "Id")
		table.ColMap("Id").SetMaxSize(store.UUID_MAX_LENGTH)
		table.ColMap("UserID").SetMaxSize(store.UUID_MAX_LENGTH)
		table.ColMap("GiftcardID").SetMaxSize(store.UUID_MAX_LENGTH)
		table.ColMap("Type").SetMaxSize(giftcard.GiftCardEventTypeMaxLength)
	}
	return gcs
}

func (gcs *SqlGiftcardEventStore) CreateIndexesIfNotExists() {
	gcs.CreateIndexIfNotExists("idx_giftcardevents_date", store.GiftcardEventTableName, "Date")
	gcs.CreateForeignKeyIfNotExists(store.GiftcardTableName, "UserID", store.UserTableName, "Id", false)
	gcs.CreateForeignKeyIfNotExists(store.GiftcardTableName, "GiftcardID", store.GiftcardTableName, "Id", true)
}

// BulkUpsert upserts and returns given giftcard events
func (gs *SqlGiftcardEventStore) BulkUpsert(transaction *gorp.Transaction, events ...*giftcard.GiftCardEvent) ([]*giftcard.GiftCardEvent, error) {
	var isSaving bool
	var upsertSelector store.SelectUpsertor = gs.GetMaster()
	if transaction != nil {
		upsertSelector = transaction
	}

	for _, event := range events {
		isSaving = false
		if !model.IsValidId(event.Id) {
			event.PreSave()
			isSaving = true
		}

		if err := event.IsValid(); err != nil {
			return nil, err
		}

		var (
			err        error
			numUpdated int64
			oldEvent   giftcard.GiftCardEvent
		)
		if isSaving {
			err = upsertSelector.Insert(event)
		} else {
			err = upsertSelector.SelectOne(&oldEvent, "SELECT * FROM "+store.GiftcardEventTableName+" WHERE Id = :ID", map[string]interface{}{"ID": event.Id})
			if err != nil {
				if err == sql.ErrNoRows {
					return nil, store.NewErrNotFound(store.GiftcardEventTableName, event.Id)
				}
				return nil, errors.Wrapf(err, "failed to find giftcard event with id=%s", event.Id)
			}

			event.Date = oldEvent.Date
			numUpdated, err = upsertSelector.Update(event)
		}

		if err != nil {
			return nil, errors.Wrapf(err, "failed to upsert a giftcard event with id=%s", event.Id)
		}

		if numUpdated != 1 {
			return nil, errors.Errorf("%d giftcard event(s) were updated instead of 1", numUpdated)
		}
	}

	return events, nil
}

func (gs *SqlGiftcardEventStore) Save(event *giftcard.GiftCardEvent) (*giftcard.GiftCardEvent, error) {
	event.PreSave()
	if err := event.IsValid(); err != nil {
		return nil, err
	}

	err := gs.GetMaster().Insert(event)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to save giftcard event with id=%s", event.Id)
	}

	return event, nil
}

func (gs *SqlGiftcardEventStore) Get(eventId string) (*giftcard.GiftCardEvent, error) {
	var res giftcard.GiftCardEvent
	err := gs.GetReplica().SelectOne(&res, "SELECT * FROM "+store.GiftcardEventTableName+" WHERE Id = :ID", map[string]interface{}{"ID": eventId})
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, store.NewErrNotFound(store.GiftcardEventTableName, eventId)
		}
		return nil, errors.Wrapf(err, "failed to find giftcard event with id=%s", eventId)
	}

	return &res, nil
}
