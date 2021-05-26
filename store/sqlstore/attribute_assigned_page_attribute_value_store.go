package sqlstore

import (
	"github.com/sitename/sitename/model/attribute"
	"github.com/sitename/sitename/store"
)

type SqlAssignedPageAttributeValueStore struct {
	*SqlStore
}

func newSqlAssignedPageAttributeValueStore(s *SqlStore) store.AssignedPageAttributeValueStore {
	as := &SqlAssignedPageAttributeValueStore{
		SqlStore: s,
	}

	for _, db := range s.GetAllConns() {
		table := db.AddTableWithName(attribute.AssignedPageAttributeValue{}, "AssignedPageAttributeValues").SetKeys(false, "Id")
		table.ColMap("Id").SetMaxSize(UUID_MAX_LENGTH)
		table.ColMap("ValueID").SetMaxSize(UUID_MAX_LENGTH)
		table.ColMap("AssignmentID").SetMaxSize(UUID_MAX_LENGTH)

		table.SetUniqueTogether("ValueID", "AssignmentID")
	}
	return as
}

func (as *SqlAssignedPageAttributeValueStore) createIndexesIfNotExists() {}
