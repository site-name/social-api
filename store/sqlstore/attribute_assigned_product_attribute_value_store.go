package sqlstore

import (
	"github.com/sitename/sitename/model/attribute"
	"github.com/sitename/sitename/store"
)

type SqlAssignedProductAttributeValueStore struct {
	*SqlStore
}

func newSqlAssignedProductAttributeValueStore(s *SqlStore) store.AssignedProductAttributeValueStore {
	as := &SqlAssignedProductAttributeValueStore{
		SqlStore: s,
	}

	for _, db := range s.GetAllConns() {
		table := db.AddTableWithName(attribute.AssignedProductAttributeValue{}, "AssignedProductAttributeValues").SetKeys(false, "Id")
		table.ColMap("Id").SetMaxSize(UUID_MAX_LENGTH)
		table.ColMap("ValueID").SetMaxSize(UUID_MAX_LENGTH)
		table.ColMap("AssignmentID").SetMaxSize(UUID_MAX_LENGTH)

		table.SetUniqueTogether("ValueID", "AssignmentID")
	}
	return as
}

func (as *SqlAssignedProductAttributeValueStore) createIndexesIfNotExists() {
	as.CreateForeignKeyIfNotExists("AssignedProductAttributeValues", "ValueID", "AttributeValues", "Id", true)
	as.CreateForeignKeyIfNotExists("AssignedProductAttributeValues", "AssignmentID", "AssignedProductAttributes", "Id", true)
}
