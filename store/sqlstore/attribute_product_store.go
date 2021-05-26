package sqlstore

import (
	"github.com/sitename/sitename/model/attribute"
	"github.com/sitename/sitename/store"
)

type SqlAttributeProductStore struct {
	*SqlStore
}

func newSqlAttributeProductStore(s *SqlStore) store.AttributeProductStore {
	as := &SqlAttributeProductStore{
		SqlStore: s,
	}

	for _, db := range s.GetAllConns() {
		table := db.AddTableWithName(attribute.AttributeProduct{}, "AttributeProducts").SetKeys(false, "Id")
		table.ColMap("Id").SetMaxSize(UUID_MAX_LENGTH)
		table.ColMap("AttributeID").SetMaxSize(UUID_MAX_LENGTH)
		table.ColMap("ProductTypeID").SetMaxSize(UUID_MAX_LENGTH)

		table.SetUniqueTogether("ProductTypeID", "AttributeID")

	}
	return as
}

func (as *SqlAttributeProductStore) createIndexesIfNotExists() {
	as.CreateForeignKeyIfNotExists("AttributeProducts", "AttributeID", "Attributes", "Id", true)
	as.CreateForeignKeyIfNotExists("AttributeProducts", "ProductTypeID", "ProductTypes", "Id", true)
}
