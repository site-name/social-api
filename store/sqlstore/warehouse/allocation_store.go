package warehouse

import (
	"database/sql"
	"fmt"

	"github.com/pkg/errors"
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

// Save takes an allocation and inserts it into database
func (as *SqlAllocationStore) Save(allocation *warehouse.Allocation) (*warehouse.Allocation, error) {
	allocation.PreSave()
	if err := allocation.IsValid(); err != nil {
		return nil, err
	}

	if err := as.GetMaster().Insert(allocation); err != nil {
		if as.IsUniqueConstraintError(err, []string{"OrderLineID", "StockID", "allocations_orderlineid_stockid_key"}) {
			return nil, store.NewErrInvalidInput(store.AllocationTableName, "OrderLineID/StockID", fmt.Sprintf("%s/%s", allocation.OrderLineID, allocation.StockID))
		}
		return nil, errors.Wrapf(err, "failed to save allocation with id=%s", allocation.Id)
	}

	return allocation, nil
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
func (as *SqlAllocationStore) FilterByOption(option *warehouse.AllocationFilterOption) ([]*warehouse.Allocation, error) {
	query := as.GetQueryBuilder().
		Select("*").
		From(store.AllocationTableName).
		OrderBy(store.TableOrderingMap[store.AllocationTableName])

	// parse option
	if option.Id != nil {
		query.Where(option.Id.ToSquirrel("Id"))
	}
	if option.OrderLineID != nil {
		query.Where(option.OrderLineID.ToSquirrel("OrderLineID"))
	}
	if option.StockID != nil {
		query.Where(option.StockID.ToSquirrel("StockID"))
	}

	queryString, args, err := query.ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "FilterbyOption_ToSql")
	}

	var res []*warehouse.Allocation
	_, err = as.GetReplica().Select(&res, queryString, args...)
	if err != nil {
		return nil, errors.Wrap(err, "failed to find allocations with given option")
	}

	return res, nil
}
