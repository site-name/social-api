package giftcard

import (
	"github.com/pkg/errors"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/store"
	"gorm.io/gorm"
)

type SqlGiftcardEventStore struct {
	store.Store
}

func NewSqlGiftcardEventStore(s store.Store) store.GiftcardEventStore {
	return &SqlGiftcardEventStore{s}
}

// BulkUpsert upserts and returns given giftcard events
func (gs *SqlGiftcardEventStore) BulkUpsert(transaction *gorm.DB, events ...*model.GiftCardEvent) ([]*model.GiftCardEvent, error) {
	if transaction == nil {
		transaction = gs.GetMaster()
	}

	for _, event := range events {
		err := transaction.Save(event).Error
		if err != nil {
			return nil, errors.Wrap(err, "failed to upsert giftcard event with")
		}
	}
	return events, nil
}

func (gs *SqlGiftcardEventStore) Save(event *model.GiftCardEvent) (*model.GiftCardEvent, error) {
	err := gs.GetMaster().Create(event).Error
	if err != nil {
		return nil, errors.Wrapf(err, "failed to save giftcard event with id=%s", event.Id)
	}

	return event, nil
}

func (gs *SqlGiftcardEventStore) Get(eventId string) (*model.GiftCardEvent, error) {
	var res model.GiftCardEvent
	err := gs.GetReplica().First(&res, "Id = ?", eventId).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, store.NewErrNotFound(model.GiftcardEventTableName, eventId)
		}
		return nil, errors.Wrapf(err, "failed to find giftcard event with id=%s", eventId)
	}

	return &res, nil
}

// FilterByOptions finds and returns a list of giftcard events with given options
func (gs *SqlGiftcardEventStore) FilterByOptions(options *model.GiftCardEventFilterOption) ([]*model.GiftCardEvent, error) {
	var res []*model.GiftCardEvent
	err := gs.GetReplica().Find(&res, store.BuildSqlizer(options.Conditions)...).Error
	if err != nil {
		return nil, errors.Wrap(err, "failed to find giftcard events with given options")
	}

	return res, nil
}
