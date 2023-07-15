package giftcard

import (
	"database/sql"

	"github.com/pkg/errors"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/modules/util"
	"github.com/sitename/sitename/store"
	"github.com/sitename/sitename/store/store_iface"
)

type SqlGiftcardEventStore struct {
	store.Store
}

func NewSqlGiftcardEventStore(s store.Store) store.GiftcardEventStore {
	return &SqlGiftcardEventStore{s}
}

func (s *SqlGiftcardEventStore) ModelFields(prefix string) util.AnyArray[string] {
	res := util.AnyArray[string]{
		"Id",
		"Date",
		"Type",
		"Parameters",
		"UserID",
		"GiftcardID",
	}
	if prefix == "" {
		return res
	}

	return res.Map(func(_ int, s string) string {
		return prefix + s
	})
}

// BulkUpsert upserts and returns given giftcard events
func (gs *SqlGiftcardEventStore) BulkUpsert(transaction store_iface.SqlxExecutor, events ...*model.GiftCardEvent) ([]*model.GiftCardEvent, error) {
	var executor store_iface.SqlxExecutor = gs.GetMasterX()
	if transaction != nil {
		executor = transaction
	}

	for _, event := range events {
		isSaving := false

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
		)
		if isSaving {
			query := "INSERT INTO " + model.GiftcardEventTableName + "(" + gs.ModelFields("").Join(",") + ") VALUES (" + gs.ModelFields(":").Join(",") + ")"
			_, err = executor.NamedExec(query, event)

		} else {
			// check if an event exist:
			var oldEvent model.GiftCardEvent
			err = executor.Get(&oldEvent, "SELECT * FROM "+model.GiftcardEventTableName+" WHERE Id = ?", event.Id)
			if err != nil {
				if err == sql.ErrNoRows {
					return nil, store.NewErrNotFound(model.GiftcardEventTableName, event.Id)
				}
				return nil, errors.Wrapf(err, "failed to find giftcard event with id=%s", event.Id)
			}

			event.Date = oldEvent.Date

			query := "UPDATE " + model.GiftcardEventTableName + " SET " + gs.
				ModelFields("").
				Map(func(_ int, s string) string {
					return s + "=:" + s
				}).
				Join(",") + " WHERE Id = :Id"

			var result sql.Result
			result, err = executor.NamedExec(query, event)
			if err == nil && result != nil {
				numUpdated, _ = result.RowsAffected()
			}
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

func (gs *SqlGiftcardEventStore) Save(event *model.GiftCardEvent) (*model.GiftCardEvent, error) {
	event.PreSave()
	if err := event.IsValid(); err != nil {
		return nil, err
	}

	query := "INSERT INTO " + model.GiftcardEventTableName + "(" + gs.ModelFields("").Join(",") + ") VALUES (" + gs.ModelFields(":").Join(",") + ")"

	_, err := gs.GetMasterX().NamedExec(query, event)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to save giftcard event with id=%s", event.Id)
	}

	return event, nil
}

func (gs *SqlGiftcardEventStore) Get(eventId string) (*model.GiftCardEvent, error) {
	var res model.GiftCardEvent
	err := gs.GetReplicaX().Get(&res, "SELECT * FROM "+model.GiftcardEventTableName+" WHERE Id = ?", eventId)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, store.NewErrNotFound(model.GiftcardEventTableName, eventId)
		}
		return nil, errors.Wrapf(err, "failed to find giftcard event with id=%s", eventId)
	}

	return &res, nil
}

// FilterByOptions finds and returns a list of giftcard events with given options
func (gs *SqlGiftcardEventStore) FilterByOptions(options *model.GiftCardEventFilterOption) ([]*model.GiftCardEvent, error) {
	query := gs.GetQueryBuilder().
		Select("*").
		From(model.GiftcardEventTableName)

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
	if options.GiftcardID != nil {
		query = query.Where(options.GiftcardID)
	}

	queryString, args, err := query.ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "FilterByOptions_ToSql")
	}

	var res []*model.GiftCardEvent
	err = gs.GetReplicaX().Select(&res, queryString, args...)
	if err != nil {
		return nil, errors.Wrap(err, "failed to find giftcard events with given options")
	}

	return res, nil
}
