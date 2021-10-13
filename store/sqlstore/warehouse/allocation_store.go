package warehouse

import (
	"database/sql"

	"github.com/Masterminds/squirrel"
	"github.com/mattermost/gorp"
	"github.com/pkg/errors"
	"github.com/sitename/sitename/model/order"
	"github.com/sitename/sitename/model/warehouse"
	"github.com/sitename/sitename/store"
)

type SqlAllocationStore struct {
	store.Store
}

func NewSqlAllocationStore(s store.Store) store.AllocationStore {
	ws := &SqlAllocationStore{s}

	for _, db := range s.GetAllConns() {
		table := db.AddTableWithName(warehouse.Allocation{}, store.AllocationTableName).SetKeys(false, "Id")
		table.ColMap("Id").SetMaxSize(store.UUID_MAX_LENGTH)
		table.ColMap("OrderLineID").SetMaxSize(store.UUID_MAX_LENGTH)
		table.ColMap("StockID").SetMaxSize(store.UUID_MAX_LENGTH)

		table.SetUniqueTogether("OrderLineID", "StockID")
	}
	return ws
}

func (ws *SqlAllocationStore) CreateIndexesIfNotExists() {
	ws.CreateForeignKeyIfNotExists(store.AllocationTableName, "OrderLineID", store.StockTableName, "Id", true)
	ws.CreateForeignKeyIfNotExists(store.AllocationTableName, "StockID", store.OrderLineTableName, "Id", true)
}

func (ws *SqlAllocationStore) ModelFields() []string {
	return []string{
		"Allocations.Id",
		"Allocations.CreateAt",
		"Allocations.OrderLineID",
		"Allocations.StockID",
		"Allocations.QuantityAllocated",
	}
}

func (ws *SqlAllocationStore) ScanFields(allocation warehouse.Allocation) []interface{} {
	return []interface{}{
		&allocation.Id,
		&allocation.CreateAt,
		&allocation.OrderLineID,
		&allocation.StockID,
		&allocation.QuantityAllocated,
	}
}

// BulkUpsert performs update, insert given allocations then returns them afterward
func (as *SqlAllocationStore) BulkUpsert(transaction *gorp.Transaction, allocations []*warehouse.Allocation) ([]*warehouse.Allocation, error) {

	var isSaving bool
	for _, allocation := range allocations {
		isSaving = false

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
			err           error
			numUpdated    int64
			oldAllocation warehouse.Allocation
		)
		if isSaving {
			err = transaction.Insert(allocation)
		} else {
			err = transaction.SelectOne(&oldAllocation, "SELECT * FROM "+store.AllocationTableName+" WHERE Id = :ID", map[string]interface{}{"ID": allocation.Id})
			if err != nil {
				if err == sql.ErrNoRows {
					return nil, store.NewErrNotFound(store.AllocationTableName, allocation.Id)
				}
				return nil, errors.Wrapf(err, "failed to find allocation with id=%s", allocation.Id)
			}

			allocation.CreateAt = oldAllocation.CreateAt

			// update
			numUpdated, err = transaction.Update(allocation)
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
func (as *SqlAllocationStore) Get(id string) (*warehouse.Allocation, error) {
	var res warehouse.Allocation
	err := as.GetReplica().SelectOne(&res, "SELECT * FROM "+store.AllocationTableName+" WHERE Id = :ID", map[string]interface{}{"ID": id})
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
func (as *SqlAllocationStore) FilterByOption(transaction *gorp.Transaction, option *warehouse.AllocationFilterOption) ([]*warehouse.Allocation, error) {
	// define fields to select:
	selectFields := as.ModelFields()
	if option.SelectedRelatedStock {
		selectFields = append(selectFields, as.Stock().ModelFields()...)
	}
	if option.SelectRelatedOrderLine {
		selectFields = append(selectFields, as.OrderLine().ModelFields()...)
	}

	query := as.GetQueryBuilder().
		Select(selectFields...).
		From(store.AllocationTableName).
		OrderBy(store.TableOrderingMap[store.AllocationTableName])

	// parse option
	if option.AnnotateStockAvailableQuantity {
		query.
			Column(`Stocks.Quantity - COALESCE( SUM( T3.QuantityAllocated ), 0 ) AS StockAvailableQuantity`). // NOTE: `T3` alias of `Allocations`
			InnerJoin(store.StockTableName+" ON (Stocks.Id = Allocations.StockID)").
			LeftJoin(store.AllocationTableName+" AS T3 ON (T3.StockID = Stocks.Id)").
			GroupBy("Allocations.Id", "Stocks.Quantity")
	}

	var joined_OrderLines_tableName bool

	if option.SelectRelatedOrderLine {
		query = query.InnerJoin(store.OrderLineTableName + " ON Orderlines.Id = Allocations.OrderLineID")
		joined_OrderLines_tableName = true // indicate for later check
	}
	if option.SelectedRelatedStock {
		query = query.InnerJoin(store.StockTableName + " ON Stocks.Id = Allocations.StockID")
	}
	if option.Id != nil {
		query = query.Where(option.Id.ToSquirrel("Id"))
	}
	if option.OrderLineID != nil {
		query = query.Where(option.OrderLineID.ToSquirrel("OrderLineID"))
	}
	if option.StockID != nil {
		query = query.Where(option.StockID.ToSquirrel("StockID"))
	}
	if option.QuantityAllocated != nil {
		query = query.Where(option.QuantityAllocated.ToSquirrel("QuantityAllocated"))
	}
	if option.OrderLineOrderID != nil {
		if !joined_OrderLines_tableName { // check if have joined Orderlines table
			query = query.InnerJoin(store.OrderLineTableName + " ON Orderlines.Id = Allocations.OrderLineID")
		}
		query = query.Where(option.OrderLineOrderID.ToSquirrel("Orderlines.OrderID"))
	}
	if option.LockForUpdate {
		query = query.Suffix("FOR UPDATE")
	}
	if option.ForUpdateOf != "" && option.LockForUpdate {
		query = query.Suffix("OF " + option.ForUpdateOf)
	}

	queryString, args, err := query.ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "FilterbyOption_ToSql")
	}

	var (
		returnAllocations      []*warehouse.Allocation
		allocation             warehouse.Allocation
		orderLine              order.OrderLine
		stock                  warehouse.Stock
		stockAvailableQuantity int
		queryer                squirrel.Queryer = as.GetReplica()
		scanFields             []interface{}    = as.ScanFields(allocation)
	)

	// check if transaction is non-nil to promote it to be actual queryer:
	if transaction != nil {
		queryer = transaction
	}

	// check if we need to modify scan list:
	if option.SelectRelatedOrderLine {
		scanFields = append(scanFields, as.OrderLine().ScanFields(orderLine)...)
	}
	if option.SelectedRelatedStock {
		scanFields = append(scanFields, as.Stock().ScanFields(stock)...)
	}
	if option.AnnotateStockAvailableQuantity {
		scanFields = append(scanFields, &stockAvailableQuantity)
	}

	rows, err := queryer.Query(queryString, args...)
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
func (as *SqlAllocationStore) BulkDelete(transaction *gorp.Transaction, allocationIDs []string) error {
	// decide which exec function to use:
	var (
		execFunc func(query string, args ...interface{}) (sql.Result, error) = as.GetMaster().Exec
	)
	if transaction != nil {
		execFunc = transaction.Exec
	}

	result, err := execFunc("DELETE FROM "+store.AllocationTableName+" WHERE Id IN $1", allocationIDs)
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
func (as *SqlAllocationStore) CountAvailableQuantityForStock(stock *warehouse.Stock) (int, error) {
	allocatedQuantity, err := as.GetReplica().SelectInt(
		`SELECT COALESCE(
			SUM (
				Allocations.QuantityAllocated
			), 0
		)
		FROM 
			Allocations 
		WHERE StockID = :StockID`,
		map[string]interface{}{"StockID": stock.Id},
	)
	if err != nil {
		return 0, errors.Wrapf(err, "failed to count allocated quantity of stock with id=%s", stock.Id)
	}

	if sub := stock.Quantity - int(allocatedQuantity); sub > 0 {
		return sub, nil
	}
	return 0, nil
}
