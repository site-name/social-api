package attribute

import (
	"github.com/sitename/sitename/model/attribute"
	"github.com/sitename/sitename/store"
)

type SqlAttributeVariantStore struct {
	store.Store
}

func NewSqlAttributeVariantStore(s store.Store) store.AttributeVariantStore {
	as := &SqlAttributeVariantStore{s}

	for _, db := range s.GetAllConns() {
		table := db.AddTableWithName(attribute.AttributeVariant{}, "AttributeVariants").SetKeys(false, "Id")
		table.ColMap("Id").SetMaxSize(store.UUID_MAX_LENGTH)
		table.ColMap("AttributeID").SetMaxSize(store.UUID_MAX_LENGTH)
		table.ColMap("ProductTypeID").SetMaxSize(store.UUID_MAX_LENGTH)

		table.SetUniqueTogether("AttributeID", "ProductTypeID")
	}
	return as
}

func (as *SqlAttributeVariantStore) CreateIndexesIfNotExists() {
	as.CreateForeignKeyIfNotExists("AttributeVariants", "AttributeID", "Attributes", "Id", true)
	as.CreateForeignKeyIfNotExists("AttributeVariants", "ProductTypeID", "ProductTypes", "Id", true)
}
