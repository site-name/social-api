package order

import (
	"github.com/sitename/sitename/model/order"
	"github.com/sitename/sitename/store"
)

type SqlOrderEventStore struct {
	store.Store
}

func NewSqlOrderEventStore(s store.Store) store.OrderEventStore {
	oes := &SqlOrderEventStore{s}

	for _, db := range s.GetAllConns() {
		table := db.AddTableWithName(order.OrderEvent{}, store.OrderEventTableName).SetKeys(false, "Id")
		table.ColMap("Id").SetMaxSize(store.UUID_MAX_LENGTH)
		table.ColMap("OrderID").SetMaxSize(store.UUID_MAX_LENGTH)
		table.ColMap("UserID").SetMaxSize(store.UUID_MAX_LENGTH)
		table.ColMap("Type").SetMaxSize(order.ORDER_EVENT_TYPE_MAX_LENGTH)
	}
	return oes
}

func (oes *SqlOrderEventStore) CreateIndexesIfNotExists() {
	oes.CreateForeignKeyIfNotExists(store.OrderEventTableName, "OrderID", store.OrderTableName, "Id", true)
	oes.CreateForeignKeyIfNotExists(store.OrderEventTableName, "UserID", store.UserTableName, "Id", false)
}
