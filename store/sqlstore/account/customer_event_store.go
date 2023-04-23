package account

import (
	"database/sql"

	"github.com/pkg/errors"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/modules/util"
	"github.com/sitename/sitename/store"
)

type SqlCustomerEventStore struct {
	store.Store
}

var customerModelFields = util.AnyArray[string]{
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

func (cs *SqlCustomerEventStore) ModelFields(prefix string) util.AnyArray[string] {
	if prefix == "" {
		return customerModelFields
	}

	return customerModelFields.Map(func(_ int, item string) string {
		return prefix + item
	})
}

func (cs *SqlCustomerEventStore) Save(event *model.CustomerEvent) (*model.CustomerEvent, error) {
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

func (cs *SqlCustomerEventStore) Get(id string) (*model.CustomerEvent, error) {
	var res model.CustomerEvent
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

func (cs *SqlCustomerEventStore) FilterByOptions(options *model.CustomerEventFilterOptions) ([]*model.CustomerEvent, error) {
	if options == nil {
		options = new(model.CustomerEventFilterOptions)
	}

	query := cs.GetQueryBuilder().Select(store.CustomerEventTableName + ".").From(store.CustomerEventTableName)

	if options.Id != nil {
		query = query.Where(options.Id)
	}
	if options.UserID != nil {
		query = query.Where(options.UserID)
	}

	str, args, err := query.ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "FilterByOptions_ToSql")
	}

	var res []*model.CustomerEvent
	err = cs.GetReplicaX().Select(&res, str, args...)
	if err != nil {
		return nil, errors.Wrap(err, "failed to find customer events by given options")
	}

	return res, nil
}
