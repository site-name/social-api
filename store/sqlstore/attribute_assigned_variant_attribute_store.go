package sqlstore

import (
	"github.com/sitename/sitename/model/attribute"
	"github.com/sitename/sitename/store"
)

type SqlAssignedVariantAttributeStore struct {
	*SqlStore
}

func newSqlAssignedVariantAttributeStore(s *SqlStore) store.AssignedVariantAttributeStore {
	as := &SqlAssignedVariantAttributeStore{
		SqlStore: s,
	}

	for _, db := range s.GetAllConns() {
		table := db.AddTableWithName(attribute.AssignedVariantAttribute{}, "AssignedVariantAttributes").SetKeys(false, "Id")
		table.ColMap("Id").SetMaxSize(UUID_MAX_LENGTH)
		table.ColMap("VariantID").SetMaxSize(UUID_MAX_LENGTH)
		table.ColMap("AssignmentID").SetMaxSize(UUID_MAX_LENGTH)

		table.SetUniqueTogether("VariantID", "AssignmentID")
	}
	return as
}

func (as *SqlAssignedVariantAttributeStore) createIndexesIfNotExists() {}
