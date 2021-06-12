package attribute

import (
	"github.com/sitename/sitename/model/attribute"
	"github.com/sitename/sitename/store"
)

type SqlAttributeProductStore struct {
	store.Store
}

func NewSqlAttributeProductStore(s store.Store) store.AttributeProductStore {
	as := &SqlAttributeProductStore{s}

	for _, db := range s.GetAllConns() {
		table := db.AddTableWithName(attribute.AttributeProduct{}, "AttributeProducts").SetKeys(false, "Id")
		table.ColMap("Id").SetMaxSize(store.UUID_MAX_LENGTH)
		table.ColMap("AttributeID").SetMaxSize(store.UUID_MAX_LENGTH)
		table.ColMap("ProductTypeID").SetMaxSize(store.UUID_MAX_LENGTH)

		table.SetUniqueTogether("ProductTypeID", "AttributeID")

	}
	return as
}

func (as *SqlAttributeProductStore) CreateIndexesIfNotExists() {
	as.CreateForeignKeyIfNotExists("AttributeProducts", "AttributeID", "Attributes", "Id", true)
	as.CreateForeignKeyIfNotExists("AttributeProducts", "ProductTypeID", "ProductTypes", "Id", true)
}
