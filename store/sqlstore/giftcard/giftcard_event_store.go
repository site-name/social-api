package giftcard

import (
	"database/sql"

	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model_helper"
	"github.com/sitename/sitename/store"
	"github.com/volatiletech/sqlboiler/v4/boil"
)

type SqlGiftcardEventStore struct {
	store.Store
}

func NewSqlGiftcardEventStore(s store.Store) store.GiftcardEventStore {
	return &SqlGiftcardEventStore{s}
}

func (gs *SqlGiftcardEventStore) Upsert(transaction boil.ContextTransactor, events model.GiftcardEventSlice) (model.GiftcardEventSlice, error) {
	if transaction == nil {
		transaction = gs.GetMaster()
	}

	for _, event := range events {
		if event == nil {
			continue
		}

		isSaving := event.ID == ""
		if isSaving {
			model_helper.GiftCardEventPreSave(event)
		}

		if err := model_helper.GiftcardEventIsValid(*event); err != nil {
			return nil, err
		}

		var err error
		if isSaving {
			err = event.Insert(transaction, boil.Infer())
		} else {
			_, err = event.Update(transaction, boil.Blacklist(model.GiftcardEventColumns.Date))
		}

		if err != nil {

		}
	}
	return events, nil
}

func (gs *SqlGiftcardEventStore) Get(id string) (*model.GiftcardEvent, error) {
	event, err := model.FindGiftcardEvent(gs.GetReplica(), id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, store.NewErrNotFound(model.TableNames.GiftcardEvents, id)
		}
		return nil, err
	}

	return event, nil
}

func (gs *SqlGiftcardEventStore) FilterByOptions(options model_helper.GiftCardEventFilterOption) (model.GiftcardEventSlice, error) {
	return model.GiftcardEvents(options.Conditions...).All(gs.GetReplica())
}
