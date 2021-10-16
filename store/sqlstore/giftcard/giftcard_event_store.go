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
		table.ColMap("Parameters").SetDataType("jsonb")
	}
	return gcs
}

func (gcs *SqlGiftcardEventStore) CreateIndexesIfNotExists() {
	gcs.CreateIndexIfNotExists("idx_giftcardevents_date", store.GiftcardEventTableName, "Date")
	gcs.CreateForeignKeyIfNotExists(store.GiftcardTableName, "UserID", store.UserTableName, "Id", false)
	gcs.CreateForeignKeyIfNotExists(store.GiftcardTableName, "GiftcardID", store.GiftcardTableName, "Id", true)
}

func (gs *SqlGiftcardEventStore) TableName(withField string) string {
	if withField == "" {
		return "GiftcardEvents"
	}
	return "GiftcardEvents." + withField
}

func (gs *SqlGiftcardEventStore) Ordering() string {
	return "GiftcardEvents.Date ASC"
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

// FilterByOptions finds and returns a list of giftcard events with given options
func (gs *SqlGiftcardEventStore) FilterByOptions(options *giftcard.GiftCardEventFilterOption) ([]*giftcard.GiftCardEvent, error) {
	query := gs.GetQueryBuilder().
		Select("*").
		From(gs.TableName("")).
		OrderBy(gs.Ordering())

	// parse options
	if options.Id != nil {
		query = query.Where(options.Id)
	}
	if options.Type != nil {
		query = query.Where(options.Type)
	}
	if options.Parameters != nil {
		query = query.Where(options.Parameters)
	}

	queryString, args, err := query.ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "FilterByOptions_ToSql")
	}

	var res []*giftcard.GiftCardEvent
	_, err = gs.GetReplica().Select(&res, queryString, args...)
	if err != nil {
		return nil, errors.Wrap(err, "failed to find giftcard events with given options")
	}

	return res, nil
}
