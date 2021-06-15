package account

import (
	"database/sql"

	"github.com/pkg/errors"
	"github.com/sitename/sitename/model/account"
	"github.com/sitename/sitename/store"
)

type SqlCustomerEventStore struct {
	store.Store
}

const (
	customerEventTableName = "CustomerEvents"
)

func NewSqlCustomerEventStore(s store.Store) store.CustomerEventStore {
	cs := &SqlCustomerEventStore{s}

	for _, db := range s.GetAllConns() {
		table := db.AddTableWithName(account.CustomerEvent{}, customerEventTableName).SetKeys(false, "Id")
		table.ColMap("Id").SetMaxSize(store.UUID_MAX_LENGTH)
		table.ColMap("OrderID").SetMaxSize(store.UUID_MAX_LENGTH)
		table.ColMap("UserID").SetMaxSize(store.UUID_MAX_LENGTH)
		table.ColMap("Type").SetMaxSize(account.CUSTOMER_EVENT_TYPE_MAX_LENGTH)
	}
	return cs
}

func (cs *SqlCustomerEventStore) CreateIndexesIfNotExists() {}

func (cs *SqlCustomerEventStore) Save(event *account.CustomerEvent) (*account.CustomerEvent, error) {
	event.PreSave()
	if err := event.IsValid(); err != nil {
		return nil, err
	}
	if err := cs.GetMaster().Insert(event); err != nil {
		return nil, errors.Wrapf(err, "failed to save CustomerEvent with Id=%s", event.Id)
	}

	return event, nil
}

func (cs *SqlCustomerEventStore) Get(id string) (*account.CustomerEvent, error) {
	var event account.CustomerEvent
	err := cs.GetReplica().SelectOne(&event, "SELECT * FROM "+customerEventTableName+" WHERE Id = :ID", map[string]interface{}{"ID": id})
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, store.NewErrNotFound(customerEventTableName, id)
		}
		return nil, errors.Wrapf(err, "failed to find CustomerEvent with Id=%s", id)
	}

	return &event, nil
}

func (cs *SqlCustomerEventStore) Count() (int64, error) {
	count, err := cs.GetReplica().SelectInt("SELECT COUNT(Id) FROM " + customerEventTableName)
	if err != nil {
		return 0, errors.Wrap(err, "failed to count number of "+customerEventTableName)
	}

	return count, nil
}
