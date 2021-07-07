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

func NewSqlCustomerEventStore(s store.Store) store.CustomerEventStore {
	cs := &SqlCustomerEventStore{s}

	for _, db := range s.GetAllConns() {
		table := db.AddTableWithName(account.CustomerEvent{}, store.CustomerEventTableName).SetKeys(false, "Id")
		table.ColMap("Id").SetMaxSize(store.UUID_MAX_LENGTH)
		table.ColMap("OrderID").SetMaxSize(store.UUID_MAX_LENGTH)
		table.ColMap("UserID").SetMaxSize(store.UUID_MAX_LENGTH)
		table.ColMap("Type").SetMaxSize(account.CUSTOMER_EVENT_TYPE_MAX_LENGTH)
	}
	return cs
}

func (cs *SqlCustomerEventStore) CreateIndexesIfNotExists() {
	cs.CreateForeignKeyIfNotExists(store.CustomerEventTableName, "OrderID", store.OrderTableName, "Id", false)
	cs.CreateForeignKeyIfNotExists(store.CustomerEventTableName, "UserID", store.UserTableName, "Id", false)
}

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
	res, err := cs.GetReplica().Get(account.CustomerEvent{}, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, store.NewErrNotFound(store.CustomerEventTableName, id)
		}
		return nil, errors.Wrapf(err, "failed to find CustomerEvent with Id=%s", id)
	}

	return res.(*account.CustomerEvent), nil
}

func (cs *SqlCustomerEventStore) Count() (int64, error) {
	count, err := cs.GetReplica().SelectInt("SELECT COUNT(Id) FROM " + store.CustomerEventTableName)
	if err != nil {
		return 0, errors.Wrap(err, "failed to count number of "+store.CustomerEventTableName)
	}

	return count, nil
}

func (cs *SqlCustomerEventStore) GetEventsByUserID(userID string) ([]*account.CustomerEvent, error) {
	var events []*account.CustomerEvent
	_, err := cs.GetReplica().Select(&events, "SELECT * FROM "+store.CustomerEventTableName+" WHERE UserID = :userID", map[string]interface{}{"userID": userID})
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, store.NewErrNotFound(store.CustomerEventTableName, "userId="+userID)
		}
		return nil, errors.Wrapf(err, "failed to find customer events with userId=%s", userID)
	}

	return events, nil
}
