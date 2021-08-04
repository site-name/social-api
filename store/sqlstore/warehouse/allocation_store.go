package warehouse

import (
	"database/sql"
	"fmt"
	"strings"

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

func (as *SqlAllocationStore) AllocationsByWhich(parentID string, toWhich warehouse.AllocationsBy) ([]*warehouse.Allocation, error) {
	var id string
	if toWhich == warehouse.ByOrderLine {
		id = "OrderLineID"
	} else if toWhich == warehouse.ByStock {
		id = "StockID"
	} else {
		return nil, store.NewErrInvalidInput(store.AllocationTableName, "to which", toWhich)
	}

	var allocations []*warehouse.Allocation
	_, err := as.GetReplica().Select(
		&allocations,
		"SELECT * FROM "+store.AllocationTableName+" WHERE "+id+" = :ParentID",
		map[string]interface{}{"ParentID": parentID},
	)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to find allocations with %s = %s", id, parentID)
	}

	return allocations, nil
}

func (as *SqlAllocationStore) AllocationsByParentIDs(parentIDs []string, which warehouse.AllocationsBy) ([]*warehouse.Allocation, error) {
	var id string
	if which == warehouse.ByOrderLine {
		id = "OrderLineID"
	} else if which == warehouse.ByStock {
		id = "StockID"
	} else {
		return nil, store.NewErrInvalidInput(store.AllocationTableName, "to which", which)
	}

	var allocations []*warehouse.Allocation
	_, err := as.GetReplica().
		Select(
			&allocations,
			"SELECT * FROM "+store.AllocationTableName+" WHERE "+id+" IN :ParentIDs",
			map[string]interface{}{
				"ParentIDs": parentIDs,
			},
		)

	if err != nil {
		return nil, errors.Wrapf(err, "failed to find allocations with %s = (%s)", id, strings.Join(parentIDs, ", "))
	}

	return allocations, nil
}
