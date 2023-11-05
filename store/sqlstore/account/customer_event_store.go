package account

import (
	"github.com/Masterminds/squirrel"
	"github.com/pkg/errors"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/store"
	"gorm.io/gorm"
)

type SqlCustomerEventStore struct {
	store.Store
}

func NewSqlCustomerEventStore(s store.Store) store.CustomerEventStore {
	return &SqlCustomerEventStore{s}
}

func (cs *SqlCustomerEventStore) Save(tx *gorm.DB, event *model.CustomerEvent) (*model.CustomerEvent, error) {
	if tx == nil {
		tx = cs.GetMaster()
	}
	err := tx.Save(event).Error
	if err != nil {
		return nil, errors.Wrap(err, "failed to save customer event")
	}

	return event, nil
}

func (cs *SqlCustomerEventStore) Get(id string) (*model.CustomerEvent, error) {
	var res model.CustomerEvent
	err := cs.GetReplica().First(&res, "Id = ?", id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, store.NewErrNotFound(model.CustomerEventTableName, id)
		}
		return nil, errors.Wrap(err, "failed to find customer event with id = "+id)
	}

	return &res, nil
}

func (cs *SqlCustomerEventStore) Count() (int64, error) {
	var count int64
	return count, cs.GetReplica().Table(model.CustomerEventTableName).Count(&count).Error
}

func (cs *SqlCustomerEventStore) FilterByOptions(options squirrel.Sqlizer) ([]*model.CustomerEvent, error) {
	args, err := store.BuildSqlizer(options, "FilterByOptions")
	if err != nil {
		return nil, err
	}

	var res []*model.CustomerEvent
	err = cs.GetReplica().Find(&res, args...).Error
	if err != nil {
		return nil, errors.Wrap(err, "failed to find customer events by given options")
	}
	return res, nil
}
