package sqlstore

import (
	"github.com/sitename/sitename/model/attribute"
	"github.com/sitename/sitename/store"
)

type SqlAssignedPageAttributeStore struct {
	*SqlStore
}

func newSqlAssignedPageAttributeStore(s *SqlStore) store.AssignedPageAttributeStore {
	as := &SqlAssignedPageAttributeStore{
		SqlStore: s,
	}

	for _, db := range s.GetAllConns() {
		table := db.AddTableWithName(attribute.AssignedPageAttribute{}, "AssignedPageAttributes").SetKeys(false, "Id")
		table.ColMap("Id").SetMaxSize(UUID_MAX_LENGTH)
		table.ColMap("PageID").SetMaxSize(UUID_MAX_LENGTH)
		table.ColMap("AssignmentID").SetMaxSize(UUID_MAX_LENGTH)

		table.SetUniqueTogether("PageID", "AssignmentID")
	}
	return as
}

func (as *SqlAssignedPageAttributeStore) createIndexesIfNotExists() {}
