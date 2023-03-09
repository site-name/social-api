package account

import (
	"database/sql"

	"github.com/pkg/errors"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/modules/util"
	"github.com/sitename/sitename/store"
)

type SqlCustomerNoteStore struct {
	store.Store
}

var customerNoteModelFields = util.AnyArray[string]{
	"Id",
	"UserID",
	"Date",
	"Content",
	"IsPublic",
	"CustomerID",
}

func NewSqlCustomerNoteStore(s store.Store) store.CustomerNoteStore {
	return &SqlCustomerNoteStore{s}
}

func (cs *SqlCustomerNoteStore) ModelFields(prefix string) util.AnyArray[string] {
	if prefix == "" {
		return customerNoteModelFields
	}

	return customerNoteModelFields.Map(func(_ int, item string) string {
		return prefix + item
	})
}

func (cs *SqlCustomerNoteStore) Save(note *model.CustomerNote) (*model.CustomerNote, error) {
	note.PreSave()
	if err := note.IsValid(); err != nil {
		return nil, err
	}

	query := "INSERT INTO " + store.CustomerNoteTableName + " (" + cs.ModelFields("").Join(",") + ") VALUES (" + cs.ModelFields(":").Join(",") + ")"
	if _, err := cs.GetMasterX().NamedExec(query, note); err != nil {
		return nil, errors.Wrapf(err, "failed to save customer note with id=%s", note.Id)
	}

	return note, nil
}

func (cs *SqlCustomerNoteStore) Get(id string) (*model.CustomerNote, error) {
	var res model.CustomerNote

	if err := cs.GetReplicaX().Get(&res, "SELECT * FROM "+store.CustomerNoteTableName+" WHERE Id = ?", id); err != nil {
		if err == sql.ErrNoRows {
			return nil, store.NewErrNotFound(store.CustomerNoteTableName, id)
		}
		return nil, errors.Wrapf(err, "failed to find customer note with id=%s", id)
	}

	return &res, nil
}
