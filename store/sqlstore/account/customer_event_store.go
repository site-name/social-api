package account

import (
	"database/sql"

	"github.com/pkg/errors"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model_helper"
	"github.com/sitename/sitename/store"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

type SqlCustomerEventStore struct {
	store.Store
}

func NewSqlCustomerEventStore(s store.Store) store.CustomerEventStore {
	return &SqlCustomerEventStore{s}
}

func (cs *SqlCustomerEventStore) Upsert(tx boil.ContextTransactor, event model.CustomerEvent) (*model.CustomerEvent, error) {
	if err := model_helper.CustomerEventIsValid(event); err != nil {
		return nil, err
	}
	if tx == nil {
		tx = cs.GetMaster()
	}
	isSaving := event.ID == ""

	var err error
	if isSaving {
		err = event.Insert(tx, boil.Infer())
	} else {
		_, err = event.Update(tx, boil.Infer())
	}
	if err != nil {
		return nil, err
	}

	return &event, nil
}

func (cs *SqlCustomerEventStore) Get(id string) (*model.CustomerEvent, error) {
	event, err := model.FindCustomerEvent(cs.GetReplica(), id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, store.NewErrNotFound(model.TableNames.CustomerEvents, id)
		}
		return nil, errors.Wrap(err, "failed to find customer event with id = "+id)
	}

	return event, nil
}

func (cs *SqlCustomerEventStore) Count() (int64, error) {
	return model.CustomerEvents().Count(cs.GetReplica())
}

func (cs *SqlCustomerEventStore) FilterByOptions(queryMods ...qm.QueryMod) (model.CustomerEventSlice, error) {
	return model.CustomerEvents(queryMods...).All(cs.GetReplica())
}
