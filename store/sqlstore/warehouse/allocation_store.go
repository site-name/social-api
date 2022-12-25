package warehouse

import (
	"database/sql"

	"github.com/Masterminds/squirrel"
	"github.com/pkg/errors"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/store"
	"github.com/sitename/sitename/store/store_iface"
)

type SqlAllocationStore struct {
	store.Store
}

func NewSqlAllocationStore(s store.Store) store.AllocationStore {
	return &SqlAllocationStore{s}
}

func (as *SqlAllocationStore) ModelFields(prefix string) model.AnyArray[string] {
	res := model.AnyArray[string]{
		"Id",
		"CreateAt",
		"OrderLineID",
		"StockID",
		"QuantityAllocated",
	}
	if prefix == "" {
		return res
	}

	return res.Map(func(_ int, s string) string {
		return prefix + s
	})
}

func (as *SqlAllocationStore) ScanFields(allocation *model.Allocation) []interface{} {
	return []interface{}{
		&allocation.Id,
		&allocation.CreateAt,
		&allocation.OrderLineID,
		&allocation.StockID,
		&allocation.QuantityAllocated,
	}
}

// BulkUpsert performs update, insert given allocations then returns them afterward
func (as *SqlAllocationStore) BulkUpsert(transaction store_iface.SqlxTxExecutor, allocations []*model.Allocation) ([]*model.Allocation, error) {
	var executor store_iface.SqlxExecutor = as.GetMasterX()
	if transaction != nil {
		executor = transaction
	}

	var (
		saveQuery   = "INSERT INTO " + store.AllocationTableName + "(" + as.ModelFields("").Join(",") + ") VALUES (" + as.ModelFields(":").Join(",") + ")"
		updateQuery = "UPDATE " + store.AllocationTableName + " SET " + as.
				ModelFields("").
				Map(func(_ int, s string) string {
				return s + "=:" + s
			}).
			Join(",") + " WHERE Id=:Id"
	)

	for _, allocation := range allocations {
		isSaving := false

		if allocation.Id == "" {
			allocation.PreSave()
			isSaving = true
		} else {
			allocation.PreUpdate()
		}

		if err := allocation.IsValid(); err != nil {
			return nil, err
		}

		var (
			err        error
			numUpdated int64
		)
		if isSaving {
			_, err = executor.NamedExec(saveQuery, allocation)

		} else {
			var result sql.Result
			result, err = executor.NamedExec(updateQuery, allocation)
			if err == nil && result != nil {
				numUpdated, _ = result.RowsAffected()
			}
		}

		if err != nil {
			if as.IsUniqueConstraintError(err, []string{"OrderLineID", "StockID", "allocations_orderlineid_stockid_key"}) {
				return nil, store.NewErrInvalidInput(store.AllocationTableName, "OrderLineID/StockID", "duplicate")
			}
			return nil, errors.Wrapf(err, "failed to upsert allocation with id=%s", allocation.Id)
		}

		if numUpdated > 1 {
			return nil, errors.Errorf("multiple allocations were updated: %d instead of 1 for allocation with id=%s", numUpdated, allocation.Id)
		}
	}

	return allocations, nil
}

// Get finds an allocation with given id then returns it with an error
func (as *SqlAllocationStore) Get(id string) (*model.Allocation, error) {
	var res model.Allocation
	err := as.GetReplicaX().Get(&res, "SELECT * FROM "+store.AllocationTableName+" WHERE Id = ?", id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, store.NewErrNotFound(store.AllocationTableName, id)
		}
		return nil, errors.Wrapf(err, "failed to find allocation with id=%s", id)
	}

	return &res, nil
}

/*
  // Sample pure SQL query when the option AnnotateStockAvailableQuantity is set to true

SELECT
  "warehouse_allocation"."id",
  "warehouse_allocation"."order_line_id",
  "warehouse_allocation"."stock_id",
  "warehouse_allocation"."quantity_allocated",
  (
    "warehouse_stock"."quantity" - COALESCE(SUM(T3."quantity_allocated"), 0)
  ) AS "stock_available_quantity"
FROM
  "warehouse_allocation"
  INNER JOIN "warehouse_stock" ON (
    "warehouse_allocation"."stock_id" = "warehouse_stock"."id"
  )
  LEFT OUTER JOIN "warehouse_allocation" T3 ON ("warehouse_stock"."id" = T3."stock_id")
WHERE
  "warehouse_allocation"."id" > 1
GROUP BY
  "warehouse_allocation"."id",
  "warehouse_stock"."quantity";
*/
// FilterByOption finds and returns a list of allocation based on given option
func (as *SqlAllocationStore) FilterByOption(transaction store_iface.SqlxTxExecutor, option *model.AllocationFilterOption) ([]*model.Allocation, error) {
	// define fields to select:
	selectFields := as.ModelFields(store.AllocationTableName + ".")
	if option.SelectedRelatedStock {
		selectFields = append(selectFields, as.Stock().ModelFields(store.StockTableName+".")...)
	}
	if option.SelectRelatedOrderLine {
		selectFields = append(selectFields, as.OrderLine().ModelFields(store.OrderLineTableName+".")...)
	}

	query := as.GetQueryBuilder().
		Select(selectFields...).
		From(store.AllocationTableName).
		OrderBy(store.TableOrderingMap[store.AllocationTableName])

	var (
		joined_OrderLines_tableName bool
		joined_Stock_table          bool
	)

	// parse option
	if option.AnnotateStockAvailableQuantity {
		query.
			Column(`Stocks.Quantity - COALESCE( SUM( T3.QuantityAllocated ), 0 ) AS StockAvailableQuantity`). // NOTE: `T3` alias of `Allocations`
			InnerJoin("Stocks ON (Stocks.Id = Allocations.StockID)").
			LeftJoin(store.AllocationTableName+" AS T3 ON (T3.StockID = Stocks.Id)").
			GroupBy("Allocations.Id", "Stocks.Quantity")

		joined_Stock_table = true // indicate for later check
	}

	if option.SelectRelatedOrderLine {
		query = query.InnerJoin(store.OrderLineTableName + " ON Orderlines.Id = Allocations.OrderLineID")

		joined_OrderLines_tableName = true // indicate for later check
	}
	if option.SelectedRelatedStock && !joined_Stock_table {
		query = query.InnerJoin(store.StockTableName + " ON Stocks.Id = Allocations.StockID")
	}
	if option.Id != nil {
		query = query.Where(option.Id)
	}
	if option.OrderLineID != nil {
		query = query.Where(option.OrderLineID)
	}
	if option.StockID != nil {
		query = query.Where(option.StockID)
	}
	if option.QuantityAllocated != nil {
		query = query.Where(option.QuantityAllocated)
	}
	if option.OrderLineOrderID != nil {
		if !joined_OrderLines_tableName { // check if have joined Orderlines table
			query = query.InnerJoin(store.OrderLineTableName + " ON Orderlines.Id = Allocations.OrderLineID")
		}
		query = query.Where(option.OrderLineOrderID)
	}
	if option.LockForUpdate {
		suf := "FOR UPDATE"
		if option.ForUpdateOf != "" {
			suf += " OF " + option.ForUpdateOf
		}
		query = query.Suffix(suf)
	}

	queryString, args, err := query.ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "FilterbyOption_ToSql")
	}

	var (
		returnAllocations      []*model.Allocation
		allocation             model.Allocation
		orderLine              model.OrderLine
		stock                  model.Stock
		stockAvailableQuantity int
		scanFields             = as.ScanFields(&allocation)

		queryer store_iface.SqlxExecutor = as.GetReplicaX()
	)

	// check if transaction is non-nil to promote it to be actual queryer:
	if transaction != nil {
		queryer = transaction
	}

	// check if we need to modify scan list:
	if option.SelectRelatedOrderLine {
		scanFields = append(scanFields, as.OrderLine().ScanFields(&orderLine)...)
	}
	if option.SelectedRelatedStock {
		scanFields = append(scanFields, as.Stock().ScanFields(&stock)...)
	}
	if option.AnnotateStockAvailableQuantity {
		scanFields = append(scanFields, &stockAvailableQuantity)
	}

	rows, err := queryer.QueryX(queryString, args...)
	if err != nil {
		return nil, errors.Wrap(err, "failed to find allocations with given option")
	}

	for rows.Next() {
		err = rows.Scan(scanFields...)
		if err != nil {
			return nil, errors.Wrap(err, "failed to scan a row")
		}

		if option.SelectRelatedOrderLine {
			allocation.OrderLine = orderLine.DeepCopy()
		}
		if option.SelectedRelatedStock {
			allocation.Stock = stock.DeepCopy()
		}
		if option.AnnotateStockAvailableQuantity {
			allocation.StockAvailableQuantity = stockAvailableQuantity
		}
		returnAllocations = append(returnAllocations, allocation.DeepCopy())
	}

	if err = rows.Close(); err != nil {
		return nil, errors.Wrap(err, "failed to close rows")
	}

	return returnAllocations, nil
}

// BulkDelete perform bulk deletes given allocations.
func (as *SqlAllocationStore) BulkDelete(transaction store_iface.SqlxTxExecutor, allocationIDs []string) error {
	var (
		executor store_iface.SqlxExecutor = as.GetMasterX()
	)
	if transaction != nil {
		executor = transaction
	}

	query, args, err := as.GetQueryBuilder().Delete(store.AllocationTableName).Where(squirrel.Eq{"Id": allocationIDs}).ToSql()
	if err != nil {
		return errors.Wrap(err, "BulkDelete_ToSql")
	}
	result, err := executor.Exec(query, args...)
	if err != nil {
		return errors.Wrap(err, "failed to delete allocations")
	}
	numDeleted, err := result.RowsAffected()
	if err != nil {
		return errors.Wrap(err, "failed to count number of allocations were deleted")
	}
	if numDeleted != int64(len(allocationIDs)) {
		return errors.Errorf("%d allocations were deleted instead of %d", numDeleted, len(allocationIDs))
	}

	return nil
}

// CountAvailableQuantityForStock counts and returns available quantity of given stock
func (as *SqlAllocationStore) CountAvailableQuantityForStock(stock *model.Stock) (int, error) {
	var count int
	err := as.GetReplicaX().Get(
		&count,
		`SELECT COALESCE(
			SUM (
				Allocations.QuantityAllocated
			), 0
		)
		FROM 
			Allocations 
		WHERE StockID = ?`,
		stock.Id,
	)
	if err != nil {
		return 0, errors.Wrapf(err, "failed to count allocated quantity of stock with id=%s", stock.Id)
	}

	if sub := stock.Quantity - count; sub > 0 {
		return sub, nil
	}
	return 0, nil
}
