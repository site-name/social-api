package warehouse

import (
	"github.com/Masterminds/squirrel"
	"github.com/pkg/errors"
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

// FilterByOption finds and returns a list of preorder allocations filtered using given options
func (ws *SqlPreorderAllocationStore) FilterByOption(options *warehouse.PreorderAllocationFilterOption) ([]*warehouse.PreorderAllocation, error) {
	query := ws.GetQueryBuilder().Select(ws.ModelFields()...).From(store.PreOrderAllocationTableName)

	and := squirrel.And{}
	// parse options
	if options.Id != nil {
		and = append(and, options.Id.ToSquirrel("PreorderAllocations.Id"))
	}
	if options.OrderLineID != nil {
		and = append(and, options.OrderLineID.ToSquirrel("PreorderAllocations.OrderLineID"))
	}
	if options.Quantity != nil {
		and = append(and, options.Quantity.ToSquirrel("PreorderAllocations.Quantity"))
	}
	if options.ProductVariantChannelListingID != nil {
		and = append(and, options.ProductVariantChannelListingID.ToSquirrel("PreorderAllocations.ProductVariantChannelListingID"))
	}

	queryString, args, err := query.Where(and).ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "FilterByOption_ToSql")
	}

	var res []*warehouse.PreorderAllocation
	_, err = ws.GetReplica().Select(&res, queryString, args...)
	if err != nil {
		return nil, errors.Wrap(err, "failed to find preorder allocations with given options")
	}

	return res, nil
}
