package warehouse

import (
	"github.com/Masterminds/squirrel"
	"github.com/mattermost/gorp"
	"github.com/pkg/errors"
	"github.com/sitename/sitename/model/order"
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
	selectFields := ws.ModelFields()

	if options.SelectRelated_OrderLine {
		selectFields = append(selectFields, ws.OrderLine().ModelFields()...)
	}
	if options.SelectRelated_OrderLine_Order && options.SelectRelated_OrderLine {
		selectFields = append(selectFields, ws.Order().ModelFields()...)
	}

	query := ws.GetQueryBuilder().Select(selectFields...).From(store.PreOrderAllocationTableName)

	andConditions := squirrel.And{}
	// parse options
	if options.Id != nil {
		andConditions = append(andConditions, options.Id.ToSquirrel("PreorderAllocations.Id"))
	}
	if options.OrderLineID != nil {
		andConditions = append(andConditions, options.OrderLineID.ToSquirrel("PreorderAllocations.OrderLineID"))
	}
	if options.Quantity != nil {
		andConditions = append(andConditions, options.Quantity.ToSquirrel("PreorderAllocations.Quantity"))
	}
	if options.ProductVariantChannelListingID != nil {
		andConditions = append(andConditions, options.ProductVariantChannelListingID.ToSquirrel("PreorderAllocations.ProductVariantChannelListingID"))
	}

	if options.SelectRelated_OrderLine {
		query = query.InnerJoin(store.OrderLineTableName + " ON PreorderAllocations.OrderLineID = Orderlines.Id")
	}
	if options.SelectRelated_OrderLine_Order && options.SelectRelated_OrderLine {
		query = query.InnerJoin(store.OrderTableName + " ON Orderlines.OrderID = Orders.Id")
	}

	rows, err := query.Where(andConditions).RunWith(ws.GetReplica()).Query()
	if err != nil {
		return nil, errors.Wrap(err, "failed to find preorder allocations with given options")
	}

	var (
		res                []*warehouse.PreorderAllocation
		preorderAllocation warehouse.PreorderAllocation
		orderLine          order.OrderLine
		orDer              order.Order
	)
	scanFields := ws.ScanFields(preorderAllocation)
	if options.SelectRelated_OrderLine {
		scanFields = append(scanFields, ws.OrderLine().ScanFields(orderLine)...)
	}
	if options.SelectRelated_OrderLine_Order && options.SelectRelated_OrderLine {
		scanFields = append(scanFields, ws.Order().ScanFields(orDer)...)
	}

	for rows.Next() {
		err = rows.Scan(scanFields...)
		if err != nil {
			return nil, errors.Wrap(err, "failed to scan a row of preorder allocation")
		}

		// join data.
		if options.SelectRelated_OrderLine_Order && options.SelectRelated_OrderLine {
			orderLine.Order = orDer.DeepCopy()
		}
		if options.SelectRelated_OrderLine {
			preorderAllocation.OrderLine = orderLine.DeepCopy()
		}

		res = append(res, preorderAllocation.DeepCopy())
	}

	if err = rows.Close(); err != nil {
		return nil, errors.Wrap(err, "failed to close rows of preorder allocations")
	}

	return res, nil
}

// Delete deletes preorder-allocations by given ids
func (ws *SqlPreorderAllocationStore) Delete(transaction *gorp.Transaction, preorderAllocationIDs ...string) error {
	var runner squirrel.BaseRunner = ws.GetMaster()
	if transaction != nil {
		runner = transaction
	}

	result, err := runner.Exec("DELETE FROM "+store.PreOrderAllocationTableName+" WHERE Id IN $1", preorderAllocationIDs)
	if err != nil {
		return errors.Wrap(err, "failed to delete preorder-allocations with given ids")
	}
	numDeleted, _ := result.RowsAffected()
	if int(numDeleted) != len(preorderAllocationIDs) {
		return errors.Errorf("%d preorder-allocations were deleted instead of %d", numDeleted, len(preorderAllocationIDs))
	}

	return nil
}
