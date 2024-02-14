package warehouse

import (
	"database/sql"
	"fmt"

	"github.com/pkg/errors"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model_helper"
	"github.com/sitename/sitename/store"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

type SqlAllocationStore struct {
	store.Store
}

func NewSqlAllocationStore(s store.Store) store.AllocationStore {
	return &SqlAllocationStore{s}
}

func (as *SqlAllocationStore) BulkUpsert(transaction boil.ContextTransactor, allocations model.AllocationSlice) (model.AllocationSlice, error) {
	if transaction == nil {
		transaction = as.GetMaster()
	}

	for _, allocation := range allocations {
		if allocation == nil {
			continue
		}

		isSaving := allocation.ID == ""
		if isSaving {
			model_helper.AllocationPreSave(allocation)
		}

		if err := model_helper.AllocationIsValid(*allocation); err != nil {
			return nil, err
		}

		var err error
		if isSaving {
			err = allocation.Insert(transaction, boil.Infer())
		} else {
			_, err = allocation.Update(transaction, boil.Blacklist(model.AllocationColumns.CreatedAt))
		}

		if err != nil {
			if as.IsUniqueConstraintError(err, []string{model.AllocationColumns.OrderLineID, model.AllocationColumns.StockID, "allocations_order_line_id_stock_id_key"}) {
				return nil, store.NewErrInvalidInput(model.TableNames.Allocations, "OrderLineID/StockID", "unique")
			}
			return nil, err
		}

		return allocations, nil
	}

	return allocations, nil
}

func (as *SqlAllocationStore) Get(id string) (*model.Allocation, error) {
	allocation, err := model.FindAllocation(as.GetReplica(), id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, store.NewErrNotFound(model.TableNames.Allocations, id)
		}
		return nil, err
	}
	return allocation, nil
}

func (as *SqlAllocationStore) FilterByOption(option model_helper.AllocationFilterOption) (model.AllocationSlice, error) {
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

func (as *SqlAllocationStore) Delete(transaction boil.ContextTransactor, ids []string) error {
	if transaction == nil {
		transaction = as.GetMaster()
	}

	_, err := model.Allocations(model.AllocationWhere.ID.IN(ids)).DeleteAll(transaction)
	return err
}

func (as *SqlAllocationStore) CountAvailableQuantityForStock(stock model.Stock) (int, error) {
	var count int

	err := model.Allocations(
		qm.Select(fmt.Sprintf("COALESCE(SUM(%s), 0) AS QuantityAllocated", model.AllocationTableColumns.QuantityAllocated)),
		model.AllocationWhere.StockID.EQ(stock.ID),
	).QueryRow(as.GetReplica()).Scan(&count)
	if err != nil {
		return 0, err
	}

	return max(0, stock.Quantity-count), nil
}
