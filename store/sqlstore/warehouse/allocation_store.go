package warehouse

import (
	"github.com/pkg/errors"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/store"
	"gorm.io/gorm"
)

type SqlAllocationStore struct {
	store.Store
}

func NewSqlAllocationStore(s store.Store) store.AllocationStore {
	return &SqlAllocationStore{s}
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
func (as *SqlAllocationStore) BulkUpsert(transaction *gorm.DB, allocations []*model.Allocation) ([]*model.Allocation, error) {
	if transaction == nil {
		transaction = as.GetMaster()
	}

	for _, allocation := range allocations {
		err := transaction.Save(allocation).Error

		if err != nil {
			if as.IsUniqueConstraintError(err, []string{"OrderLineID", "StockID", "allocations_orderlineid_stockid_key"}) {
				return nil, store.NewErrInvalidInput(model.AllocationTableName, "OrderLineID/StockID", "duplicate")
			}
			return nil, errors.Wrapf(err, "failed to upsert allocation with id=%s", allocation.Id)
		}
	}

	return allocations, nil
}

// Get finds an allocation with given id then returns it with an error
func (as *SqlAllocationStore) Get(id string) (*model.Allocation, error) {
	var res model.Allocation
	err := as.GetReplica().First(&res, "Id = ?", id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, store.NewErrNotFound(model.AllocationTableName, id)
		}
		return nil, errors.Wrapf(err, "failed to find allocation with id=%s", id)
	}

	return &res, nil
}

// FilterByOption finds and returns a list of allocation based on given option
func (as *SqlAllocationStore) FilterByOption(option *model.AllocationFilterOption) ([]*model.Allocation, error) {
	// define fields to select:
	selectFields := []string{model.AllocationTableName + ".*"}
	if option.SelectedRelatedStock {
		selectFields = append(selectFields, model.StockTableName+".*")
	}
	if option.SelectRelatedOrderLine {
		selectFields = append(selectFields, model.OrderLineTableName+".*")
	}

	query := as.GetQueryBuilder().
		Select(selectFields...).
		From(model.AllocationTableName).
		Where(option.Conditions)

	if option.AnnotateStockAvailableQuantity || option.SelectedRelatedStock {
		query = query.InnerJoin(model.StockTableName + " ON Stocks.Id = Allocations.StockID")

		if option.AnnotateStockAvailableQuantity {
			query = query.
				Column(`Stocks.Quantity - COALESCE( SUM( Allocations.QuantityAllocated ), 0 ) AS StockAvailableQuantity`).
				LeftJoin(model.AllocationTableName+" ON Allocations.StockID = Stocks.Id").
				GroupBy("Allocations.Id", "Stocks.Quantity")
		}
	}

	// parse options
	if option.SelectRelatedOrderLine ||
		option.OrderLineOrderID != nil {
		query = query.InnerJoin(model.OrderLineTableName + " ON Orderlines.Id = Allocations.OrderLineID")
		if option.OrderLineOrderID != nil {
			query = query.Where(option.OrderLineOrderID)
		}
	}

	if option.LockForUpdate && option.Transaction != nil {
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

	runner := as.GetReplica()
	if option.Transaction != nil {
		runner = option.Transaction
	}

	rows, err := runner.Raw(queryString, args...).Rows()
	if err != nil {
		return nil, errors.Wrap(err, "failed to find allocations with given option")
	}
	defer rows.Close()

	var returnAllocations model.Allocations

	for rows.Next() {
		var (
			allocation model.Allocation
			orderLine  model.OrderLine
			stock      model.Stock
			scanFields = as.ScanFields(&allocation)
		)

		// NOTE: scan order must be identical to select order (like above)
		if option.SelectedRelatedStock {
			scanFields = append(scanFields, as.Stock().ScanFields(&stock)...)
		}
		if option.SelectRelatedOrderLine {
			scanFields = append(scanFields, as.OrderLine().ScanFields(&orderLine)...)
		}
		if option.AnnotateStockAvailableQuantity {
			scanFields = append(scanFields, &allocation.StockAvailableQuantity)
		}

		err = rows.Scan(scanFields...)
		if err != nil {
			return nil, errors.Wrap(err, "failed to scan a row of allocation")
		}

		if option.SelectRelatedOrderLine {
			allocation.OrderLine = &orderLine
		}
		if option.SelectedRelatedStock {
			allocation.Stock = &stock
		}
		returnAllocations = append(returnAllocations, &allocation)
	}

	return returnAllocations, nil
}

// BulkDelete perform bulk deletes given allocations.
func (as *SqlAllocationStore) BulkDelete(transaction *gorm.DB, allocationIDs []string) error {
	if transaction == nil {
		transaction = as.GetMaster()
	}

	err := transaction.Raw("DELETE FROM "+model.AllocationTableName+" WHERE Id IN ?", allocationIDs).Error
	if err != nil {
		return errors.Wrap(err, "failed to delete allocations")
	}

	return nil
}

// CountAvailableQuantityForStock counts and returns available quantity of given stock
func (as *SqlAllocationStore) CountAvailableQuantityForStock(stock *model.Stock) (int, error) {
	var count int
	err := as.GetReplica().Raw(
		`SELECT COALESCE(
			SUM (
				Allocations.QuantityAllocated
			), 0
		)
		FROM 
			Allocations 
		WHERE StockID = ?`,
		stock.Id,
	).
		Scan(&count).
		Error
	if err != nil {
		return 0, errors.Wrapf(err, "failed to count allocated quantity of stock with id=%s", stock.Id)
	}

	if sub := stock.Quantity - count; sub > 0 {
		return sub, nil
	}
	return 0, nil
}
