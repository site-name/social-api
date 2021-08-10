package order

import (
	"database/sql"

	"github.com/pkg/errors"
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

func (oes *SqlOrderEventStore) Save(orderEvent *order.OrderEvent) (*order.OrderEvent, error) {
	orderEvent.PreSave()
	if err := orderEvent.IsValid(); err != nil {
		return nil, err
	}

	if err := oes.GetMaster().Insert(orderEvent); err != nil {
		return nil, errors.Wrapf(err, "failed to save order event with id=%s", orderEvent.Id)
	}

	return orderEvent, nil
}

func (oes *SqlOrderEventStore) Get(orderEventID string) (*order.OrderEvent, error) {
	var res order.OrderEvent
	err := oes.GetReplica().SelectOne(&res, "SELECT * FROM "+store.OrderEventTableName+" WHERE Id = :ID", map[string]interface{}{"ID": orderEventID})
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, store.NewErrNotFound(store.OrderEventTableName, orderEventID)
		}
		return nil, errors.Wrapf(err, "failed to find order event iwth id=%s", orderEventID)
	}

	return &res, nil
}
