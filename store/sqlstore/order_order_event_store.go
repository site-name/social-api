package sqlstore

import (
	"github.com/sitename/sitename/model/order"
	"github.com/sitename/sitename/store"
)

type SqlOrderEventStore struct {
	*SqlStore
}

func newSqlOrderEventStore(s *SqlStore) store.OrderEventStore {
	oes := &SqlOrderEventStore{s}

	for _, db := range s.GetAllConns() {
		table := db.AddTableWithName(order.OrderEvent{}, "OrderEvents").SetKeys(false, "Id")
		table.ColMap("Id").SetMaxSize(UUID_MAX_LENGTH)
		table.ColMap("OrderID").SetMaxSize(UUID_MAX_LENGTH)
		table.ColMap("UserID").SetMaxSize(UUID_MAX_LENGTH)
		table.ColMap("Type").SetMaxSize(order.ORDER_EVENT_TYPE_MAX_LENGTH)
	}
	return oes
}

func (oes *SqlOrderEventStore) createIndexesIfNotExists() {

}
