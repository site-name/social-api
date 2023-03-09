package warehouse

import (
	"github.com/Masterminds/squirrel"
	"github.com/pkg/errors"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/modules/util"
	"github.com/sitename/sitename/store"
	"github.com/sitename/sitename/store/store_iface"
)

type SqlPreorderAllocationStore struct {
	store.Store
}

func NewSqlPreorderAllocationStore(s store.Store) store.PreorderAllocationStore {
	return &SqlPreorderAllocationStore{s}
}

func (ws *SqlPreorderAllocationStore) ModelFields(prefix string) util.AnyArray[string] {
	res := util.AnyArray[string]{
		"Id",
		"OrderLineID",
		"Quantity",
		"ProductVariantChannelListingID",
	}
	if prefix == "" {
		return res
	}

	return res.Map(func(_ int, s string) string {
		return prefix + s
	})
}
func (ws *SqlPreorderAllocationStore) ScanFields(preorderAllocation *model.PreorderAllocation) []interface{} {
	return []interface{}{
		&preorderAllocation.Id,
		&preorderAllocation.OrderLineID,
		&preorderAllocation.Quantity,
		&preorderAllocation.ProductVariantChannelListingID,
	}
}

// BulkCreate bulk inserts given preorderAllocations and returns them
func (ws *SqlPreorderAllocationStore) BulkCreate(transaction store_iface.SqlxTxExecutor, preorderAllocations []*model.PreorderAllocation) ([]*model.PreorderAllocation, error) {
	var upsertor store_iface.SqlxExecutor = ws.GetMasterX()
	if transaction != nil {
		upsertor = transaction
	}

	query := "INSERT INTO " + store.PreOrderAllocationTableName + "(" + ws.ModelFields("").Join(",") + ") VALUES (" + ws.ModelFields(":").Join(",") + ")"
	for _, allocation := range preorderAllocations {
		allocation.PreSave()

		if err := allocation.IsValid(); err != nil {
			return nil, err
		}

		_, err := upsertor.NamedExec(query, allocation)
		if err != nil {
			if ws.IsUniqueConstraintError(err, []string{"OrderLineID", "ProductVariantChannelListingID", "preorderallocations_orderlineid_productvariantchannellistingid_key"}) {
				return nil, store.NewErrInvalidInput(store.PreOrderAllocationTableName, "OrderLineID/ProductVariantChannelListingID", "duplicate")
			}
			return nil, errors.Wrapf(err, "failed to insert preorder allocation with id=%s", allocation.Id)
		}
	}

	return preorderAllocations, nil
}

// FilterByOption finds and returns a list of preorder allocations filtered using given options
func (ws *SqlPreorderAllocationStore) FilterByOption(options *model.PreorderAllocationFilterOption) ([]*model.PreorderAllocation, error) {
	selectFields := ws.ModelFields(store.PreOrderAllocationTableName + ".")

	if options.SelectRelated_OrderLine {
		selectFields = append(selectFields, ws.OrderLine().ModelFields(store.OrderLineTableName+".")...)
	}
	if options.SelectRelated_OrderLine_Order && options.SelectRelated_OrderLine {
		selectFields = append(selectFields, ws.Order().ModelFields(store.OrderTableName+".")...)
	}

	query := ws.GetQueryBuilder().Select(selectFields...).From(store.PreOrderAllocationTableName)

	andConditions := squirrel.And{}
	// parse options
	if options.Id != nil {
		andConditions = append(andConditions, options.Id)
	}
	if options.OrderLineID != nil {
		andConditions = append(andConditions, options.OrderLineID)
	}
	if options.Quantity != nil {
		andConditions = append(andConditions, options.Quantity)
	}
	if options.ProductVariantChannelListingID != nil {
		andConditions = append(andConditions, options.ProductVariantChannelListingID)
	}

	if options.SelectRelated_OrderLine {
		query = query.InnerJoin(store.OrderLineTableName + " ON PreorderAllocations.OrderLineID = Orderlines.Id")
	}
	if options.SelectRelated_OrderLine_Order && options.SelectRelated_OrderLine {
		query = query.InnerJoin(store.OrderTableName + " ON Orderlines.OrderID = Orders.Id")
	}

	queryString, args, err := query.ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "FilterByOption_ToSql")
	}

	rows, err := ws.GetReplicaX().QueryX(queryString, args...)
	if err != nil {
		return nil, errors.Wrap(err, "failed to find preorder allocations with given options")
	}

	var (
		res                []*model.PreorderAllocation
		preorderAllocation model.PreorderAllocation
		orderLine          model.OrderLine
		orDer              model.Order
		scanFields         = ws.ScanFields(&preorderAllocation)
	)
	if options.SelectRelated_OrderLine {
		scanFields = append(scanFields, ws.OrderLine().ScanFields(&orderLine)...)
	}
	if options.SelectRelated_OrderLine && options.SelectRelated_OrderLine_Order {
		scanFields = append(scanFields, ws.Order().ScanFields(&orDer)...)
	}

	for rows.Next() {
		err = rows.Scan(scanFields...)
		if err != nil {
			return nil, errors.Wrap(err, "failed to scan a row of preorder allocation")
		}

		// join data.
		if options.SelectRelated_OrderLine && options.SelectRelated_OrderLine_Order {
			orderLine.SetOrder(&orDer) // no need deep copy here, line 143 takes care of that
		}
		if options.SelectRelated_OrderLine {
			preorderAllocation.SetOrderLine(&orderLine) // no need deep copy here, line 143 takes care of that
		}

		res = append(res, preorderAllocation.DeepCopy())
	}

	if err = rows.Close(); err != nil {
		return nil, errors.Wrap(err, "failed to close rows of preorder allocations")
	}

	return res, nil
}

// Delete deletes preorder-allocations by given ids
func (ws *SqlPreorderAllocationStore) Delete(transaction store_iface.SqlxTxExecutor, preorderAllocationIDs ...string) error {
	var runner store_iface.SqlxExecutor = ws.GetMasterX()
	if transaction != nil {
		runner = transaction
	}

	query, args, err := ws.GetQueryBuilder().Delete(store.PreOrderAllocationTableName).Where(squirrel.Eq{"Id": preorderAllocationIDs}).ToSql()
	if err != nil {
		return errors.Wrap(err, "Delete_ToSql")
	}
	result, err := runner.Exec(query, args...)
	if err != nil {
		return errors.Wrap(err, "failed to delete preorder-allocations with given ids")
	}
	numDeleted, _ := result.RowsAffected()
	if int(numDeleted) != len(preorderAllocationIDs) {
		return errors.Errorf("%d preorder-allocations were deleted instead of %d", numDeleted, len(preorderAllocationIDs))
	}

	return nil
}
