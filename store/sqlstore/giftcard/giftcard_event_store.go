package giftcard

import (
	"database/sql"

	"github.com/pkg/errors"
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
