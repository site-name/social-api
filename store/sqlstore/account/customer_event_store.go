package account

import (
	"github.com/Masterminds/squirrel"
	"github.com/pkg/errors"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/modules/util"
	"github.com/sitename/sitename/store"
)

type SqlCustomerEventStore struct {
	store.Store
}

var customerModelFields = util.AnyArray[string]{
	"Id",
	"Date",
	"Type",
	"OrderID",
	"UserID",
	"Parameters",
}

func NewSqlCustomerEventStore(s store.Store) store.CustomerEventStore {
	return &SqlCustomerEventStore{s}
}

func (cs *SqlCustomerEventStore) ModelFields(prefix string) util.AnyArray[string] {
	if prefix == "" {
		return customerModelFields
	}

	return customerModelFields.Map(func(_ int, item string) string {
		return prefix + item
	})
}

func (cs *SqlCustomerEventStore) Save(event *model.CustomerEvent) (*model.CustomerEvent, error) {
	err := cs.GetMaster().Create(event).Error
	if err != nil {
		return nil, errors.Wrap(err, "failed to save customer event")
	}

	return event, nil
}

func (cs *SqlCustomerEventStore) Get(id string) (*model.CustomerEvent, error) {
	var res model.CustomerEvent

	err := cs.GetReplica().First(&res, "Id = ?", id).Error
	if err != nil {
		return nil, err
	}

	return &res, nil
}

func (cs *SqlCustomerEventStore) Count() (int64, error) {
	var count int64
	err := cs.GetReplica().Table(model.CustomerEventTableName).Count(&count).Error
	return count, err
}

func (cs *SqlCustomerEventStore) FilterByOptions(options squirrel.Sqlizer) ([]*model.CustomerEvent, error) {
	var res []*model.CustomerEvent
	err := cs.GetReplica().Find(&res, store.BuildSqlizer(options)...).Error
	if err != nil {
		return nil, errors.Wrap(err, "failed to find customer events by given options")
	}
	return res, nil
}
