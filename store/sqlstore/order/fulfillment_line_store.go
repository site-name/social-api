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

const (
	fulfillmentLineTableName = "FulfillmentLines"
)

func NewSqlFulfillmentLineStore(s store.Store) store.FulfillmentLineStore {
	fls := &SqlFulfillmentLineStore{s}

	for _, db := range s.GetAllConns() {
		table := db.AddTableWithName(order.FulfillmentLine{}, fulfillmentLineTableName).SetKeys(false, "Id")
		table.ColMap("Id").SetMaxSize(store.UUID_MAX_LENGTH)
		table.ColMap("OrderLineID").SetMaxSize(store.UUID_MAX_LENGTH)
		table.ColMap("FulfillmentID").SetMaxSize(store.UUID_MAX_LENGTH)
		table.ColMap("StockID").SetMaxSize(store.UUID_MAX_LENGTH)
	}

	return fls
}

func (fls *SqlFulfillmentLineStore) CreateIndexesIfNotExists() {}

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
	var ffml order.FulfillmentLine
	if err := fls.GetReplica().SelectOne(&ffml,
		"SELECT * FROM "+fulfillmentLineTableName+" WHERE Id = :id",
		map[string]interface{}{"id": id},
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, store.NewErrNotFound(fulfillmentLineTableName, id)
		}
		return nil, errors.Wrapf(err, "failed to find fulfillment line with id=%s", id)
	}

	return &ffml, nil
}
