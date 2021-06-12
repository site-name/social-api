package order

import (
	"github.com/sitename/sitename/model/order"
	"github.com/sitename/sitename/store"
)

type SqlFulfillmentLineStore struct {
	store.Store
}

func NewSqlFulfillmentLineStore(s store.Store) store.FulfillmentLineStore {
	fls := &SqlFulfillmentLineStore{s}

	for _, db := range s.GetAllConns() {
		table := db.AddTableWithName(order.FulfillmentLine{}, "FulfillmentLines").SetKeys(false, "Id")
		table.ColMap("Id").SetMaxSize(store.UUID_MAX_LENGTH)
		table.ColMap("OrderLineID").SetMaxSize(store.UUID_MAX_LENGTH)
		table.ColMap("FulfillmentID").SetMaxSize(store.UUID_MAX_LENGTH)
		table.ColMap("StockID").SetMaxSize(store.UUID_MAX_LENGTH)
	}

	return fls
}

func (fls *SqlFulfillmentLineStore) CreateIndexesIfNotExists() {

}
