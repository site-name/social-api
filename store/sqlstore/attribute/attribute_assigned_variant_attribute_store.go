package attribute

import (
	"github.com/sitename/sitename/model/attribute"
	"github.com/sitename/sitename/store"
)

type SqlAssignedVariantAttributeStore struct {
	store.Store
}

func NewSqlAssignedVariantAttributeStore(s store.Store) store.AssignedVariantAttributeStore {
	as := &SqlAssignedVariantAttributeStore{s}

	for _, db := range s.GetAllConns() {
		table := db.AddTableWithName(attribute.AssignedVariantAttribute{}, "AssignedVariantAttributes").SetKeys(false, "Id")
		table.ColMap("Id").SetMaxSize(store.UUID_MAX_LENGTH)
		table.ColMap("VariantID").SetMaxSize(store.UUID_MAX_LENGTH)
		table.ColMap("AssignmentID").SetMaxSize(store.UUID_MAX_LENGTH)

		table.SetUniqueTogether("VariantID", "AssignmentID")
	}
	return as
}

func (as *SqlAssignedVariantAttributeStore) CreateIndexesIfNotExists() {
	as.CreateForeignKeyIfNotExists("AssignedVariantAttributes", "VariantID", "ProductVariants", "Id", true)
	as.CreateForeignKeyIfNotExists("AssignedVariantAttributes", "AssignmentID", "AttributeVariants", "Id", true)
}
