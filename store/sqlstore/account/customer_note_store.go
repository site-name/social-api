package account

import (
	"database/sql"

	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model_helper"
	"github.com/sitename/sitename/store"
	"github.com/volatiletech/sqlboiler/v4/boil"
)

type SqlCustomerNoteStore struct {
	store.Store
}

func NewSqlCustomerNoteStore(s store.Store) store.CustomerNoteStore {
	return &SqlCustomerNoteStore{s}
}

func (cs *SqlCustomerNoteStore) Upsert(note model.CustomerNote) (*model.CustomerNote, error) {
	if err := model_helper.CustomerNoteIsValid(note); err != nil {
		return nil, err
	}
	isSaving := note.ID == ""

	var err error
	if isSaving {
		err = note.Insert(cs.GetMaster(), boil.Infer())
	} else {
		_, err = note.Update(cs.GetMaster(), boil.Infer())
	}
	if err != nil {
		return nil, err
	}

	return &note, nil
}

func (cs *SqlCustomerNoteStore) Get(id string) (*model.CustomerNote, error) {
	note, err := model.FindCustomerNote(cs.GetReplica(), id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, store.NewErrNotFound(model.TableNames.CustomerNotes, id)
		}
		return nil, err
	}

	return note, nil
}
