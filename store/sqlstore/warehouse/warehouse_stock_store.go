package warehouse

import (
	"github.com/sitename/sitename/model/warehouse"
	"github.com/sitename/sitename/store"
)

type SqlStockStore struct {
	store.Store
}

func NewSqlStockStore(s store.Store) store.StockStore {
	ws := &SqlStockStore{s}

	for _, db := range s.GetAllConns() {
		table := db.AddTableWithName(warehouse.Stock{}, "Stocks").SetKeys(false, "Id")
		table.ColMap("Id").SetMaxSize(store.UUID_MAX_LENGTH)
		table.ColMap("WarehouseID").SetMaxSize(store.UUID_MAX_LENGTH)
		table.ColMap("ProductVariantID").SetMaxSize(store.UUID_MAX_LENGTH)

		table.SetUniqueTogether("WarehouseID", "ProductVariantID")
	}
	return ws
}

func (ws *SqlStockStore) CreateIndexesIfNotExists() {

}
