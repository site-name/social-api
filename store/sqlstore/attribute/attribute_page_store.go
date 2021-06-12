package attribute

import (
	"github.com/sitename/sitename/model/attribute"
	"github.com/sitename/sitename/store"
)

type SqlAttributePageStore struct {
	store.Store
}

func NewSqlAttributePageStore(s store.Store) store.AttributePageStore {
	as := &SqlAttributePageStore{s}

	for _, db := range s.GetAllConns() {
		table := db.AddTableWithName(attribute.AttributePage{}, "AttributePages").SetKeys(false, "Id")
		table.ColMap("Id").SetMaxSize(store.UUID_MAX_LENGTH)
		table.ColMap("AttributeID").SetMaxSize(store.UUID_MAX_LENGTH)
		table.ColMap("PageTypeID").SetMaxSize(store.UUID_MAX_LENGTH)

		table.SetUniqueTogether("AttributeID", "PageTypeID")
	}
	return as
}

func (as *SqlAttributePageStore) CreateIndexesIfNotExists() {
	as.CreateForeignKeyIfNotExists("AttributePages", "AttributeID", "Attributes", "Id", true)
	as.CreateForeignKeyIfNotExists("AttributePages", "PageTypeID", "PageTypes", "Id", true)
}
