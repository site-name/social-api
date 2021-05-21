package sqlstore

import (
	"github.com/sitename/sitename/model/warehouse"
	"github.com/sitename/sitename/store"
)

type SqlStockStore struct {
	*SqlStore
}

func newSqlStockStore(s *SqlStore) store.StockStore {
	ws := &SqlStockStore{s}

	for _, db := range s.GetAllConns() {
		table := db.AddTableWithName(warehouse.Stock{}, "Stocks").SetKeys(false, "Id")
		table.ColMap("Id").SetMaxSize(UUID_MAX_LENGTH)
		table.ColMap("WarehouseID").SetMaxSize(UUID_MAX_LENGTH)
		table.ColMap("ProductVariantID").SetMaxSize(UUID_MAX_LENGTH)

		table.SetUniqueTogether("WarehouseID", "ProductVariantID")
	}
	return ws
}

func (ws *SqlStockStore) createIndexesIfNotExists() {

}
