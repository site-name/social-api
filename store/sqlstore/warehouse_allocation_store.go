package sqlstore

import (
	"github.com/sitename/sitename/model/warehouse"
	"github.com/sitename/sitename/store"
)

type SqlAllocationStore struct {
	*SqlStore
}

func newSqlAllocationStore(s *SqlStore) store.AllocationStore {
	ws := &SqlAllocationStore{s}

	for _, db := range s.GetAllConns() {
		table := db.AddTableWithName(warehouse.Allocation{}, "Allocations").SetKeys(false, "Id")
		table.ColMap("Id").SetMaxSize(UUID_MAX_LENGTH)
		table.ColMap("OrderLineID").SetMaxSize(UUID_MAX_LENGTH).SetNotNull(true)
		table.ColMap("StockID").SetMaxSize(UUID_MAX_LENGTH).SetNotNull(true)

		table.SetUniqueTogether("OrderLineID", "StockID")
	}
	return ws
}

func (ws *SqlAllocationStore) createIndexesIfNotExists() {

}
