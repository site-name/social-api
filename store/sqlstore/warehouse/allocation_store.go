package warehouse

import (
	"database/sql"
	"fmt"

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

func (as *SqlAllocationStore) commonQueryBuilder(options model_helper.AllocationFilterOption) []qm.QueryMod {
	conds := options.Conditions
	if options.OrderLineOrderID != nil {
		conds = append(
			conds,
			qm.InnerJoin(fmt.Sprintf("%s ON %s = %s", model.TableNames.OrderLines, model.AllocationTableColumns.OrderLineID, model.OrderLineTableColumns.ID)),
			options.OrderLineOrderID,
		)
	}
	for _, load := range options.Preloads {
		conds = append(conds, qm.Load(load))
	}

	// TODO: check if joining conditions below is valid

	if options.AnnotateStockAvailableQuantity {
		var annotations = model_helper.AnnotationAggregator{
			model_helper.AllocationAnnotationKeys.AvailableQuantity: fmt.Sprintf("%s - COALESCE( SUM(%s), 0 )", model.StockTableColumns.Quantity, model.AllocationTableColumns.QuantityAllocated),
		}

		conds = append(
			conds,
			qm.Select(model.TableNames.Allocations+".*"), // NOTE: this is needed
			annotations,
			qm.InnerJoin(fmt.Sprintf("%s ON %s = %s", model.TableNames.Stocks, model.AllocationTableColumns.StockID, model.StockTableColumns.ID)),
			qm.LeftOuterJoin(fmt.Sprintf("%s ON %s = %s", model.TableNames.Allocations, model.AllocationTableColumns.StockID, model.StockTableColumns.ID)),
			qm.GroupBy(model.AllocationTableColumns.ID+","+model.StockTableColumns.Quantity),
		)
	}

	return conds
}

func (as *SqlAllocationStore) FilterByOption(option model_helper.AllocationFilterOption) (model.AllocationSlice, error) {
	conds := as.commonQueryBuilder(option)
	return model.Allocations(conds...).All(as.GetReplica())
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
