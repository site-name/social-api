package attribute

import (
	"github.com/sitename/sitename/model/attribute"
	"github.com/sitename/sitename/store"
)

type SqlAssignedPageAttributeValueStore struct {
	store.Store
}

func NewSqlAssignedPageAttributeValueStore(s store.Store) store.AssignedPageAttributeValueStore {
	as := &SqlAssignedPageAttributeValueStore{s}

	for _, db := range s.GetAllConns() {
		table := db.AddTableWithName(attribute.AssignedPageAttributeValue{}, "AssignedPageAttributeValues").SetKeys(false, "Id")
		table.ColMap("Id").SetMaxSize(store.UUID_MAX_LENGTH)
		table.ColMap("ValueID").SetMaxSize(store.UUID_MAX_LENGTH)
		table.ColMap("AssignmentID").SetMaxSize(store.UUID_MAX_LENGTH)

		table.SetUniqueTogether("ValueID", "AssignmentID")
	}
	return as
}

func (as *SqlAssignedPageAttributeValueStore) CreateIndexesIfNotExists() {
	as.CreateForeignKeyIfNotExists("AssignedPageAttributeValues", "ValueID", "AttributeValues", "Id", true)
	as.CreateForeignKeyIfNotExists("AssignedPageAttributeValues", "AssignmentID", "AssignedPageAttributes", "Id", true)
}
