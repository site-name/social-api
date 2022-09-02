package order

import (
	"database/sql"

	"github.com/pkg/errors"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model/order"
	"github.com/sitename/sitename/store"
	"github.com/sitename/sitename/store/store_iface"
)

type SqlOrderEventStore struct {
	store.Store
}

func NewSqlOrderEventStore(s store.Store) store.OrderEventStore {
	return &SqlOrderEventStore{s}
}

func (s *SqlOrderEventStore) ModelFields(prefix string) model.AnyArray[string] {
	res := model.AnyArray[string]{
		"Id",
		"CreateAt",
		"Type",
		"OrderID",
		"Parameters",
		"UserID",
	}
	if prefix == "" {
		return res
	}

	return res.Map(func(_ int, s string) string {
		return prefix + s
	})
}

func (oes *SqlOrderEventStore) Save(transaction store_iface.SqlxTxExecutor, orderEvent *order.OrderEvent) (*order.OrderEvent, error) {
	var executor store_iface.SqlxExecutor = oes.GetMasterX()
	if transaction != nil {
		executor = transaction
	}

	orderEvent.PreSave()
	if err := orderEvent.IsValid(); err != nil {
		return nil, err
	}

	query := "INSERT INTO " + store.OrderEventTableName + "(" + oes.ModelFields("").Join(",") + ") VALUES (" + oes.ModelFields(":").Join(",") + ")"
	if _, err := executor.NamedExec(query, orderEvent); err != nil {
		return nil, errors.Wrapf(err, "failed to save order event with id=%s", orderEvent.Id)
	}

	return orderEvent, nil
}

func (oes *SqlOrderEventStore) Get(orderEventID string) (*order.OrderEvent, error) {
	var res order.OrderEvent
	err := oes.GetReplicaX().Get(&res, "SELECT * FROM "+store.OrderEventTableName+" WHERE Id = ?", orderEventID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, store.NewErrNotFound(store.OrderEventTableName, orderEventID)
		}
		return nil, errors.Wrapf(err, "failed to find order event iwth id=%s", orderEventID)
	}

	return &res, nil
}
