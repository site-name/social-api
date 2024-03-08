package order

import (
	"github.com/sitename/sitename/store"
)

type SqlOrderEventStore struct {
	store.Store
}

func NewSqlOrderEventStore(s store.Store) store.OrderEventStore {
	return &SqlOrderEventStore{s}
}

// func (oes *SqlOrderEventStore) Save(transaction boil.ContextTransactor, orderEvent *model.OrderEvent) (*model.OrderEvent, error) {
// 	if transaction == nil {
// 		transaction = oes.GetMaster()
// 	}

// 	if err := transaction.Create(orderEvent).Error; err != nil {
// 		return nil, errors.Wrap(err, "failed to save order event")
// 	}

// 	return orderEvent, nil
// }

// func (oes *SqlOrderEventStore) Get(orderEventID string) (*model.OrderEvent, error) {
// 	model.FindOrderEvent()
// }

// func (s *SqlOrderEventStore) FilterByOptions(options *model.OrderEventFilterOptions) ([]*model.OrderEvent, error) {
// 	args, err := store.BuildSqlizer(options.Conditions, "OrderEvent_FilterByOptions")
// 	if err != nil {
// 		return nil, err
// 	}

// 	var res []*model.OrderEvent
// 	err = s.GetReplica().Find(&res, args...).Error
// 	if err != nil {
// 		return nil, errors.Wrap(err, "failed to find order events by given options")
// 	}

// 	return res, nil
// }
