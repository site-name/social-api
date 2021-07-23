package attribute

import (
	"github.com/sitename/sitename/model/attribute"
	"github.com/sitename/sitename/store"
)

type SqlAssignedPageAttributeStore struct {
	store.Store
}

func NewSqlAssignedPageAttributeStore(s store.Store) store.AssignedPageAttributeStore {
	as := &SqlAssignedPageAttributeStore{
		Store: s,
	}

	for _, db := range s.GetAllConns() {
		table := db.AddTableWithName(attribute.AssignedPageAttribute{}, store.AssignedPageAttributeTableName).SetKeys(false, "Id")
		table.ColMap("Id").SetMaxSize(store.UUID_MAX_LENGTH)
		table.ColMap("PageID").SetMaxSize(store.UUID_MAX_LENGTH)
		table.ColMap("AssignmentID").SetMaxSize(store.UUID_MAX_LENGTH)

		table.SetUniqueTogether("PageID", "AssignmentID")
	}
	return as
}

func (as *SqlAssignedPageAttributeStore) CreateIndexesIfNotExists() {
	as.CreateForeignKeyIfNotExists(store.AssignedPageAttributeTableName, "AssignmentID", store.AttributePageTableName, "Id", true)
	as.CreateForeignKeyIfNotExists(store.AssignedPageAttributeTableName, "PageID", "Pages", "Id", true)
}
