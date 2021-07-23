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
		table := db.AddTableWithName(attribute.AssignedVariantAttributeValue{}, store.AssignedVariantAttributeValueTableName).SetKeys(false, "Id")
		table.ColMap("Id").SetMaxSize(store.UUID_MAX_LENGTH)
		table.ColMap("ValueID").SetMaxSize(store.UUID_MAX_LENGTH)
		table.ColMap("AssignmentID").SetMaxSize(store.UUID_MAX_LENGTH)

		table.SetUniqueTogether("ValueID", "AssignmentID")
	}
	return as
}

func (as *SqlAssignedVariantAttributeValueStore) CreateIndexesIfNotExists() {
	as.CreateForeignKeyIfNotExists(store.AssignedVariantAttributeValueTableName, "ValueID", store.AttributeValueTableName, "Id", true)
	as.CreateForeignKeyIfNotExists(store.AssignedVariantAttributeValueTableName, "AssignmentID", store.AssignedVariantAttributeTableName, "Id", true)
}
