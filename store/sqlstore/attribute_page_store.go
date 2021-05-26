package sqlstore

import (
	"github.com/sitename/sitename/model/attribute"
	"github.com/sitename/sitename/store"
)

type SqlAttributePageStore struct {
	*SqlStore
}

func newSqlAttributePageStore(s *SqlStore) store.AttributePageStore {
	as := &SqlAttributePageStore{
		SqlStore: s,
	}

	for _, db := range s.GetAllConns() {
		table := db.AddTableWithName(attribute.AttributePage{}, "AttributePages").SetKeys(false, "Id")
		table.ColMap("Id").SetMaxSize(UUID_MAX_LENGTH)
		table.ColMap("AttributeID").SetMaxSize(UUID_MAX_LENGTH)
		table.ColMap("PageTypeID").SetMaxSize(UUID_MAX_LENGTH)

		table.SetUniqueTogether("AttributeID", "PageTypeID")
	}
	return as
}

func (as *SqlAttributePageStore) createIndexesIfNotExists() {}
