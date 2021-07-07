package account

import (
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
	cs.CreateForeignKeyIfNotExists(store.CustomerNoteTableName, "UserID", "Users", "Id", false)
	cs.CreateForeignKeyIfNotExists(store.CustomerNoteTableName, "CustomerID", "Users", "Id", false)
}
