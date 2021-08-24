package warehouse

import (
	"github.com/pkg/errors"
	"github.com/sitename/sitename/model/warehouse"
	"github.com/sitename/sitename/store"
)

type SqlWarehouseShippingZoneStore struct {
	store.Store
}

func NewSqlWarehouseShippingZoneStore(s store.Store) store.WarehouseShippingZoneStore {
	ws := &SqlWarehouseShippingZoneStore{s}
	for _, db := range s.GetAllConns() {
		table := db.AddTableWithName(warehouse.WarehouseShippingZone{}, store.WarehouseShippingZoneTableName).SetKeys(false, "Id")
		table.ColMap("Id").SetMaxSize(store.UUID_MAX_LENGTH)
		table.ColMap("WarehouseID").SetMaxSize(store.UUID_MAX_LENGTH)
		table.ColMap("ShippingZoneID").SetMaxSize(store.UUID_MAX_LENGTH)

		table.SetUniqueTogether("WarehouseID", "ShippingZoneID")
	}

	return ws
}

func (ws *SqlWarehouseShippingZoneStore) CreateIndexesIfNotExists() {
	ws.CreateForeignKeyIfNotExists(store.WarehouseShippingZoneTableName, "WarehouseID", store.WarehouseTableName, "Id", false)
	ws.CreateForeignKeyIfNotExists(store.WarehouseShippingZoneTableName, "ShippingZoneID", store.ShippingZoneTableName, "Id", false)
}

// Save inserts given warehouse-shipping zone relation into database
func (ws *SqlWarehouseShippingZoneStore) Save(warehouseShippingZone *warehouse.WarehouseShippingZone) (*warehouse.WarehouseShippingZone, error) {
	warehouseShippingZone.PreSave()
	if err := warehouseShippingZone.IsValid(); err != nil {
		return nil, err
	}

	err := ws.GetMaster().Insert(warehouseShippingZone)
	if err != nil {
		if ws.IsUniqueConstraintError(err, []string{"WarehouseID", "ShippingZoneID", "warehouseshippingzones_warehouseid_shippingzoneid_key"}) {
			return nil, store.NewErrInvalidInput(store.WarehouseShippingZoneTableName, "WarehouseID/ShippingZoneID", "duplicate")
		}
		return nil, errors.Wrapf(err, "failed to save warehouse-shipping zone relation with id=%s", warehouseShippingZone.Id)
	}

	return warehouseShippingZone, nil
}
