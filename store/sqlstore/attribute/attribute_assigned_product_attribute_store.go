package attribute

import (
	"github.com/sitename/sitename/model/attribute"
	"github.com/sitename/sitename/store"
)

type SqlAssignedProductAttributeStore struct {
	store.Store
}

func NewSqlAssignedProductAttributeStore(s store.Store) store.AssignedProductAttributeStore {
	as := &SqlAssignedProductAttributeStore{s}

	for _, db := range s.GetAllConns() {
		table := db.AddTableWithName(attribute.AssignedProductAttribute{}, "AssignedProductAttributes").SetKeys(false, "Id")
		table.ColMap("Id").SetMaxSize(store.UUID_MAX_LENGTH)
		table.ColMap("ProductID").SetMaxSize(store.UUID_MAX_LENGTH)
		table.ColMap("AssignmentID").SetMaxSize(store.UUID_MAX_LENGTH)

		table.SetUniqueTogether("ProductID", "AssignmentID")
	}
	return as
}

func (as *SqlAssignedProductAttributeStore) CreateIndexesIfNotExists() {
	as.CreateForeignKeyIfNotExists("AssignedProductAttributes", "ProductID", "Products", "Id", true)
	as.CreateForeignKeyIfNotExists("AssignedProductAttributes", "AssignmentID", "AttributeProducts", "Id", true)
}
