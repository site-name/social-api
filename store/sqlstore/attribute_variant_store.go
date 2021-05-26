package sqlstore

import (
	"github.com/sitename/sitename/model/attribute"
	"github.com/sitename/sitename/store"
)

type SqlAttributeVariantStore struct {
	*SqlStore
}

func newSqlAttributeVariantStore(s *SqlStore) store.AttributeVariantStore {
	as := &SqlAttributeVariantStore{
		SqlStore: s,
	}

	for _, db := range s.GetAllConns() {
		table := db.AddTableWithName(attribute.AttributeVariant{}, "Attributes").SetKeys(false, "Id")
		table.ColMap("Id").SetMaxSize(UUID_MAX_LENGTH)
		table.ColMap("AttributeID").SetMaxSize(UUID_MAX_LENGTH)
		table.ColMap("ProductTypeID").SetMaxSize(UUID_MAX_LENGTH)

		table.SetUniqueTogether("AttributeID", "ProductTypeID")
	}
	return as
}

func (as *SqlAttributeVariantStore) createIndexesIfNotExists() {}
