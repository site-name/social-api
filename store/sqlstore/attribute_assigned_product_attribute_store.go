package sqlstore

import (
	"github.com/sitename/sitename/model/attribute"
	"github.com/sitename/sitename/store"
)

type SqlAssignedProductAttributeStore struct {
	*SqlStore
}

func newSqlAssignedProductAttributeStore(s *SqlStore) store.AssignedProductAttributeStore {
	as := &SqlAssignedProductAttributeStore{
		SqlStore: s,
	}

	for _, db := range s.GetAllConns() {
		table := db.AddTableWithName(attribute.AssignedProductAttribute{}, "AssignedProductAttributes").SetKeys(false, "Id")
		table.ColMap("Id").SetMaxSize(UUID_MAX_LENGTH)
		table.ColMap("ProductID").SetMaxSize(UUID_MAX_LENGTH)
		table.ColMap("AssignmentID").SetMaxSize(UUID_MAX_LENGTH)

		table.SetUniqueTogether("ProductID", "AssignmentID")
	}
	return as
}

func (as *SqlAssignedProductAttributeStore) createIndexesIfNotExists() {
	as.CreateForeignKeyIfNotExists("AssignedProductAttributes", "ProductID", "Products", "Id", true)
	as.CreateForeignKeyIfNotExists("AssignedProductAttributes", "AssignmentID", "AttributeProducts", "Id", true)
}
