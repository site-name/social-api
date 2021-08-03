package account

import (
	"github.com/pkg/errors"
	"github.com/sitename/sitename/model/account"
	"github.com/sitename/sitename/store"
)

type SqlCustomerNoteStore struct {
	store.Store
}

func NewSqlCustomerNoteStore(s store.Store) store.CustomerNoteStore {
	cs := &SqlCustomerNoteStore{s}
	for _, db := range s.GetAllConns() {
		table := db.AddTableWithName(account.CustomerNote{}, store.CustomerNoteTableName)
		table.ColMap("Id").SetMaxSize(store.UUID_MAX_LENGTH)
		table.ColMap("UserID").SetMaxSize(store.UUID_MAX_LENGTH)
		table.ColMap("CustomerID").SetMaxSize(store.UUID_MAX_LENGTH)
	}

	return cs
}

func (cs *SqlCustomerNoteStore) CreateIndexesIfNotExists() {
	cs.CreateIndexIfNotExists("idx_customer_notes_date", store.CustomerNoteTableName, "Date")
	cs.CreateForeignKeyIfNotExists(store.CustomerNoteTableName, "UserID", store.UserTableName, "Id", false)
	cs.CreateForeignKeyIfNotExists(store.CustomerNoteTableName, "CustomerID", store.UserTableName, "Id", true)
}

func (cs *SqlCustomerNoteStore) Save(note *account.CustomerNote) (*account.CustomerNote, error) {
	note.PreSave()
	if err := note.IsValid(); err != nil {
		return nil, err
	}

	if err := cs.GetMaster().Insert(note); err != nil {
		return nil, errors.Wrapf(err, "failed to save customer note with id=%s", note.Id)
	}

	return note, nil
}

func (cs *SqlCustomerNoteStore) Get(id string) (*account.CustomerNote, error) {
	if res, err := cs.GetReplica().Get(account.CustomerNote{}, id); err != nil {
		return nil, errors.Wrapf(err, "failed to find customer note with id=%s", id)
	} else {
		return res.(*account.CustomerNote), nil
	}
}
