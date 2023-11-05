package order

import (
	"github.com/pkg/errors"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/store"
	"gorm.io/gorm"
)

type SqlOrderEventStore struct {
	store.Store
}

func NewSqlOrderEventStore(s store.Store) store.OrderEventStore {
	return &SqlOrderEventStore{s}
}

func (oes *SqlOrderEventStore) Save(transaction *gorm.DB, orderEvent *model.OrderEvent) (*model.OrderEvent, error) {
	if transaction == nil {
		transaction = oes.GetMaster()
	}

	if err := transaction.Create(orderEvent).Error; err != nil {
		return nil, errors.Wrap(err, "failed to save order event")
	}

	return orderEvent, nil
}

func (oes *SqlOrderEventStore) Get(orderEventID string) (*model.OrderEvent, error) {
	var res model.OrderEvent
	err := oes.GetReplica().First(&res, "Id = ?", orderEventID).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, store.NewErrNotFound(model.OrderEventTableName, orderEventID)
		}
		return nil, errors.Wrapf(err, "failed to find order event iwth id=%s", orderEventID)
	}

	return &res, nil
}

func (s *SqlOrderEventStore) FilterByOptions(options *model.OrderEventFilterOptions) ([]*model.OrderEvent, error) {
	args, err := store.BuildSqlizer(options.Conditions, "OrderEvent_FilterByOptions")
	if err != nil {
		return nil, err
	}

	var res []*model.OrderEvent
	err = s.GetReplica().Find(&res, args...).Error
	if err != nil {
		return nil, errors.Wrap(err, "failed to find order events by given options")
	}

	return res, nil
}
