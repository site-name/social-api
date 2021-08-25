package warehouse

import (
	"database/sql"

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
	if option.SelectRelatedOrderLine {
		query = query.InnerJoin(store.OrderLineTableName + " ON Orderlines.Id = Allocations.OrderLineID")
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
	if option.LockForUpdate {
		query = query.Suffix("FOR UPDATE")
	}
	if option.ForUpdateOf != "" {
		query = query.Suffix("OF " + option.ForUpdateOf)
	}

	queryString, args, err := query.ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "FilterbyOption_ToSql")
	}

	var rows *sql.Rows
	if transaction == nil {
		rows, err = as.GetReplica().Query(queryString, args...)
	} else {
		rows, err = transaction.Query(queryString, args...)
	}
	if err != nil {
		return nil, errors.Wrap(err, "failed to find allocations with given option")
	}

	var (
		returnAllocations []*warehouse.Allocation
		allocation        warehouse.Allocation
		orderLine         order.OrderLine
		stock             warehouse.Stock
	)
	var scanFields []interface{} = []interface{}{
		&allocation.Id,
		&allocation.CreateAt,
		&allocation.OrderLineID,
		&allocation.StockID,
		&allocation.QuantityAllocated,
	}
	if option.SelectRelatedOrderLine {
		scanFields = append(
			scanFields,

			&orderLine.Id,
			&orderLine.CreateAt,
			&orderLine.OrderID,
			&orderLine.VariantID,
			&orderLine.ProductName,
			&orderLine.VariantName,
			&orderLine.TranslatedProductName,
			&orderLine.TranslatedVariantName,
			&orderLine.ProductSku,
			&orderLine.IsShippingRequired,
			&orderLine.Quantity,
			&orderLine.QuantityFulfilled,
			&orderLine.Currency,
			&orderLine.UnitDiscountAmount,
			&orderLine.UnitDiscountType,
			&orderLine.UnitDiscountReason,
			&orderLine.UnitPriceNetAmount,
			&orderLine.UnitDiscountValue,
			&orderLine.UnitPriceGrossAmount,
			&orderLine.TotalPriceNetAmount,
			&orderLine.TotalPriceGrossAmount,
			&orderLine.UnDiscountedUnitPriceGrossAmount,
			&orderLine.UnDiscountedUnitPriceNetAmount,
			&orderLine.UnDiscountedTotalPriceGrossAmount,
			&orderLine.UnDiscountedTotalPriceNetAmount,
			&orderLine.TaxRate,
		)
	}
	if option.SelectedRelatedStock {
		scanFields = append(
			scanFields,

			&stock.Id,
			&stock.CreateAt,
			&stock.WarehouseID,
			&stock.ProductVariantID,
			&stock.Quantity,
		)
	}

	for rows.Next() {
		err = rows.Scan(scanFields...)
		if err != nil {
			return nil, errors.Wrap(err, "failed to scan a row")
		}

		if option.SelectRelatedOrderLine {
			allocation.OrderLine = &orderLine
		}
		if option.SelectedRelatedStock {
			allocation.Stock = &stock
		}
		returnAllocations = append(returnAllocations, &allocation)
	}

	if err = rows.Close(); err != nil {
		return nil, errors.Wrap(err, "failed to close rows")
	}

	return returnAllocations, nil
}
