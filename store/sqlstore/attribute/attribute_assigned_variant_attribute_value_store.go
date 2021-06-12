package attribute

import (
	"github.com/sitename/sitename/model/attribute"
	"github.com/sitename/sitename/store"
)

type SqlAssignedVariantAttributeValueStore struct {
	store.Store
}

func NewSqlAssignedVariantAttributeValueStore(s store.Store) store.AssignedVariantAttributeValueStore {
	as := &SqlAssignedVariantAttributeValueStore{s}

	for _, db := range s.GetAllConns() {
		table := db.AddTableWithName(attribute.AssignedVariantAttributeValue{}, "AssignedVariantAttributeValues").SetKeys(false, "Id")
		table.ColMap("Id").SetMaxSize(store.UUID_MAX_LENGTH)
		table.ColMap("ValueID").SetMaxSize(store.UUID_MAX_LENGTH)
		table.ColMap("AssignmentID").SetMaxSize(store.UUID_MAX_LENGTH)

		table.SetUniqueTogether("ValueID", "AssignmentID")
	}
	return as
}

func (as *SqlAssignedVariantAttributeValueStore) CreateIndexesIfNotExists() {
	as.CreateForeignKeyIfNotExists("AssignedVariantAttributeValues", "ValueID", "AttributeValues", "Id", true)
	as.CreateForeignKeyIfNotExists("AssignedVariantAttributeValues", "AssignmentID", "AssignedVariantAttributes", "Id", true)
}
