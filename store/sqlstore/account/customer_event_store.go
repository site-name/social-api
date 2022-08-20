package account

import (
	"database/sql"

	"github.com/pkg/errors"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model/account"
	"github.com/sitename/sitename/store"
)

type SqlCustomerEventStore struct {
	store.Store
}

var customerModelFields = model.StringArray{
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

func (cs *SqlCustomerEventStore) ModelFields(prefix string) model.StringArray {
	if prefix == "" {
		return customerModelFields
	}

	return customerModelFields.Map(func(_ int, item string) string {
		return prefix + item
	})
}

func (cs *SqlCustomerEventStore) Save(event *account.CustomerEvent) (*account.CustomerEvent, error) {
	event.PreSave()
	if err := event.IsValid(); err != nil {
		return nil, err
	}

	query := "INSERT INTO " + store.CustomerEventTableName + " (" + cs.ModelFields("").Join(",") + ") VALUES (" + cs.ModelFields(":").Join(",") + ")"
	if _, err := cs.GetMasterX().NamedExec(query, event); err != nil {
		return nil, errors.Wrapf(err, "failed to save CustomerEvent with Id=%s", event.Id)
	}

	return event, nil
}

func (cs *SqlCustomerEventStore) Get(id string) (*account.CustomerEvent, error) {
	var res account.CustomerEvent
	err := cs.GetMasterX().Get(&res, "SELECT * FROM "+store.CustomerEventTableName+" WHERE Id = ?", id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, store.NewErrNotFound(store.CustomerEventTableName, id)
		}
		return nil, errors.Wrapf(err, "failed to find CustomerEvent with Id=%s", id)
	}

	return &res, nil
}

func (cs *SqlCustomerEventStore) Count() (int64, error) {
	var count int64
	err := cs.GetReplicaX().Select(&count, "SELECT COUNT(Id) FROM "+store.CustomerEventTableName)
	if err != nil {
		return 0, errors.Wrap(err, "failed to count number of "+store.CustomerEventTableName)
	}

	return count, nil
}

func (cs *SqlCustomerEventStore) GetEventsByUserID(userID string) ([]*account.CustomerEvent, error) {
	var events []*account.CustomerEvent
	err := cs.GetReplicaX().Select(&events, "SELECT * FROM "+store.CustomerEventTableName+" WHERE UserID = ?", userID)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to find customer events with userId=%s", userID)
	}

	return events, nil
}
