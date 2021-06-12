package attribute

import (
	"github.com/sitename/sitename/model/attribute"
	"github.com/sitename/sitename/store"
)

type SqlAssignedProductAttributeValueStore struct {
	store.Store
}

func NewSqlAssignedProductAttributeValueStore(s store.Store) store.AssignedProductAttributeValueStore {
	as := &SqlAssignedProductAttributeValueStore{s}

	for _, db := range s.GetAllConns() {
		table := db.AddTableWithName(attribute.AssignedProductAttributeValue{}, "AssignedProductAttributeValues").SetKeys(false, "Id")
		table.ColMap("Id").SetMaxSize(store.UUID_MAX_LENGTH)
		table.ColMap("ValueID").SetMaxSize(store.UUID_MAX_LENGTH)
		table.ColMap("AssignmentID").SetMaxSize(store.UUID_MAX_LENGTH)

		table.SetUniqueTogether("ValueID", "AssignmentID")
	}
	return as
}

func (as *SqlAssignedProductAttributeValueStore) CreateIndexesIfNotExists() {
	as.CreateForeignKeyIfNotExists("AssignedProductAttributeValues", "ValueID", "AttributeValues", "Id", true)
	as.CreateForeignKeyIfNotExists("AssignedProductAttributeValues", "AssignmentID", "AssignedProductAttributes", "Id", true)
}
