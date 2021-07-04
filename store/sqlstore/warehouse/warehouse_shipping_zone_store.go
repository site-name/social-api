package warehouse

import (
	"github.com/sitename/sitename/model/warehouse"
	"github.com/sitename/sitename/store"
	"github.com/sitename/sitename/store/sqlstore/shipping"
)

const (
	WarehouseShippingZoneTableName = "WarehouseShippingZones"
)

type SqlWarehouseShippingZoneStore struct {
	store.Store
}

func NewSqlWarehouseShippingZoneStore(s store.Store) store.WarehouseShippingZoneStore {
	ws := &SqlWarehouseShippingZoneStore{s}
	for _, db := range s.GetAllConns() {
		table := db.AddTableWithName(warehouse.WarehouseShippingZone{}, WarehouseShippingZoneTableName).SetKeys(false, "Id")
		table.ColMap("Id").SetMaxSize(store.UUID_MAX_LENGTH)
		table.ColMap("WarehouseID").SetMaxSize(store.UUID_MAX_LENGTH)
		table.ColMap("ShippingZoneID").SetMaxSize(store.UUID_MAX_LENGTH)

		table.SetUniqueTogether("WarehouseID", "ShippingZoneID")
	}

	return ws
}

func (ws *SqlWarehouseShippingZoneStore) CreateIndexesIfNotExists() {
	ws.CreateForeignKeyIfNotExists(WarehouseShippingZoneTableName, "WarehouseID", WarehouseTableName, "Id", false)
	ws.CreateForeignKeyIfNotExists(WarehouseShippingZoneTableName, "ShippingZoneID", shipping.ShippingZoneTableName, "Id", false)
}
