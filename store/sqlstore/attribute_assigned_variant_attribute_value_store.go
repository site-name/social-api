package sqlstore

import (
	"github.com/sitename/sitename/model/attribute"
	"github.com/sitename/sitename/store"
)

type SqlAssignedVariantAttributeValueStore struct {
	*SqlStore
}

func newSqlAssignedVariantAttributeValueStore(s *SqlStore) store.AssignedVariantAttributeValueStore {
	as := &SqlAssignedVariantAttributeValueStore{
		SqlStore: s,
	}

	for _, db := range s.GetAllConns() {
		table := db.AddTableWithName(attribute.AssignedVariantAttributeValue{}, "AssignedVariantAttributeValues").SetKeys(false, "Id")
		table.ColMap("Id").SetMaxSize(UUID_MAX_LENGTH)
		table.ColMap("ValueID").SetMaxSize(UUID_MAX_LENGTH)
		table.ColMap("AssignmentID").SetMaxSize(UUID_MAX_LENGTH)

		table.SetUniqueTogether("ValueID", "AssignmentID")
	}
	return as
}

func (as *SqlAssignedVariantAttributeValueStore) createIndexesIfNotExists() {
	as.CreateForeignKeyIfNotExists("AssignedVariantAttributeValues", "ValueID", "AttributeValues", "Id", true)
	as.CreateForeignKeyIfNotExists("AssignedVariantAttributeValues", "AssignmentID", "AssignedVariantAttributes", "Id", true)
}
