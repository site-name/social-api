package warehouse

import (
	"github.com/sitename/sitename/model/warehouse"
	"github.com/sitename/sitename/store"
)

type SqlPreorderAllocationStore struct {
	store.Store
}

func NewSqlPreorderAllocationStore(s store.Store) store.PreorderAllocationStore {
	ps := &SqlPreorderAllocationStore{s}

	for _, db := range s.GetAllConns() {
		table := db.AddTableWithName(warehouse.PreorderAllocation{}, store.PreOrderAllocationTableName).SetKeys(false, "Id")
		table.ColMap("Id").SetMaxSize(store.UUID_MAX_LENGTH)
		table.ColMap("OrderLineID").SetMaxSize(store.UUID_MAX_LENGTH)
		table.ColMap("ProductVariantChannelListingID").SetMaxSize(store.UUID_MAX_LENGTH)

		table.SetUniqueTogether("OrderLineID", "ProductVariantChannelListingID")
	}

	return ps
}

func (ws *SqlPreorderAllocationStore) CreateIndexesIfNotExists() {}

func (ws *SqlPreorderAllocationStore) ModelFields() []string {
	return []string{
		"PreorderAllocations.Id",
		"PreorderAllocations.OrderLineID",
		"PreorderAllocations.Quantity",
		"PreorderAllocations.ProductVariantChannelListingID",
	}
}
func (ws *SqlPreorderAllocationStore) ScanFields(preorderAllocation warehouse.PreorderAllocation) []interface{} {
	return []interface{}{
		&preorderAllocation.Id,
		&preorderAllocation.OrderLineID,
		&preorderAllocation.Quantity,
		&preorderAllocation.ProductVariantChannelListingID,
	}
}
