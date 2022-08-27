package warehouse

import (
	"github.com/pkg/errors"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model/warehouse"
	"github.com/sitename/sitename/store"
)

type SqlWarehouseShippingZoneStore struct {
	store.Store
}

func NewSqlWarehouseShippingZoneStore(s store.Store) store.WarehouseShippingZoneStore {
	return &SqlWarehouseShippingZoneStore{s}
}

func (ws *SqlWarehouseShippingZoneStore) ModelFields(prefix string) model.StringArray {
	res := model.StringArray{
		"Id",
		"WarehouseID",
		"ShippingZoneID",
	}
	if prefix == "" {
		return res
	}

	return res.Map(func(_ int, s string) string {
		return prefix + s
	})
}

// Save inserts given warehouse-shipping zone relation into database
func (ws *SqlWarehouseShippingZoneStore) Save(warehouseShippingZone *warehouse.WarehouseShippingZone) (*warehouse.WarehouseShippingZone, error) {
	warehouseShippingZone.PreSave()
	if err := warehouseShippingZone.IsValid(); err != nil {
		return nil, err
	}

	query := "INSERT INTO " + store.WarehouseShippingZoneTableName + "(" + ws.ModelFields("").Join(",") + ") VALUES (" + ws.ModelFields(":").Join(",") + ")"
	_, err := ws.GetMasterX().NamedExec(query, warehouseShippingZone)
	if err != nil {
		if ws.IsUniqueConstraintError(err, []string{"WarehouseID", "ShippingZoneID", "warehouseshippingzones_warehouseid_shippingzoneid_key"}) {
			return nil, store.NewErrInvalidInput(store.WarehouseShippingZoneTableName, "WarehouseID/ShippingZoneID", "duplicate")
		}
		return nil, errors.Wrapf(err, "failed to save warehouse-shipping zone relation with id=%s", warehouseShippingZone.Id)
	}

	return warehouseShippingZone, nil
}
