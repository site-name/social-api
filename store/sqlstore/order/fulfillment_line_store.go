package order

import (
	"database/sql"

	"github.com/pkg/errors"
	"github.com/sitename/sitename/model/order"
	"github.com/sitename/sitename/store"
)

type SqlFulfillmentLineStore struct {
	store.Store
}

func NewSqlFulfillmentLineStore(s store.Store) store.FulfillmentLineStore {
	fls := &SqlFulfillmentLineStore{s}

	for _, db := range s.GetAllConns() {
		table := db.AddTableWithName(order.FulfillmentLine{}, store.FulfillmentLineTableName).SetKeys(false, "Id")
		table.ColMap("Id").SetMaxSize(store.UUID_MAX_LENGTH)
		table.ColMap("OrderLineID").SetMaxSize(store.UUID_MAX_LENGTH)
		table.ColMap("FulfillmentID").SetMaxSize(store.UUID_MAX_LENGTH)
		table.ColMap("StockID").SetMaxSize(store.UUID_MAX_LENGTH)
	}

	return fls
}

func (fls *SqlFulfillmentLineStore) CreateIndexesIfNotExists() {
	fls.CreateForeignKeyIfNotExists(store.FulfillmentLineTableName, "OrderLineID", store.OrderLineTableName, "Id", true)
	fls.CreateForeignKeyIfNotExists(store.FulfillmentLineTableName, "FulfillmentID", store.FulfillmentTableName, "Id", true)
	fls.CreateForeignKeyIfNotExists(store.FulfillmentLineTableName, "StockID", store.StockTableName, "Id", false)
}

func (fls *SqlFulfillmentLineStore) Save(ffml *order.FulfillmentLine) (*order.FulfillmentLine, error) {
	ffml.PreSave()
	if err := ffml.IsValid(); err != nil {
		return nil, err
	}

	if err := fls.GetMaster().Insert(ffml); err != nil {
		return nil, errors.Wrapf(err, "failed to save fulfillment line with id=%s", ffml.Id)
	}
	return ffml, nil
}

func (fls *SqlFulfillmentLineStore) Get(id string) (*order.FulfillmentLine, error) {
	var res order.FulfillmentLine
	if err := fls.GetReplica().SelectOne(&res, "SELECT * FROM "+store.FulfillmentLineTableName+" WHERE Id = :ID", map[string]interface{}{"ID": id}); err != nil {
		if err == sql.ErrNoRows {
			return nil, store.NewErrNotFound(store.FulfillmentLineTableName, id)
		}
		return nil, errors.Wrapf(err, "failed to find fulfillment line with id=%s", id)
	} else {
		return &res, nil
	}
}
