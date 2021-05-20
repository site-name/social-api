package sqlstore

import (
	"github.com/sitename/sitename/model/order"
	"github.com/sitename/sitename/store"
)

type SqlFulfillmentLineStore struct {
	*SqlStore
}

func newSqlFulfillmentLineStore(s *SqlStore) store.FulfillmentLineStore {
	fls := &SqlFulfillmentLineStore{s}

	for _, db := range s.GetAllConns() {
		table := db.AddTableWithName(order.FulfillmentLine{}, "FulfillmentLines").SetKeys(false, "Id")
		table.ColMap("Id").SetMaxSize(UUID_MAX_LENGTH)
		table.ColMap("OrderLineID").SetMaxSize(UUID_MAX_LENGTH)
		table.ColMap("FulfillmentID").SetMaxSize(UUID_MAX_LENGTH)
		table.ColMap("StockID").SetMaxSize(UUID_MAX_LENGTH)
	}

	return fls
}

func (fls *SqlFulfillmentLineStore) createIndexesIfNotExists() {

}
