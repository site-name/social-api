package warehouse

import (
	"github.com/sitename/sitename/app"
	"github.com/sitename/sitename/app/sub_app_iface"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model/warehouse"
	"github.com/sitename/sitename/store"
)

type AppWarehouse struct {
	app.AppIface
}

func init() {
	app.RegisterWarehouseApp(func(a app.AppIface) sub_app_iface.WarehouseApp {
		return &AppWarehouse{a}
	})
}

// WarehouseByOption returns a list of warehouses based on given option
func (a *AppWarehouse) WarehouseByOption(option *warehouse.WarehouseFilterOption) ([]*warehouse.WareHouse, *model.AppError) {
	warehouses, err := a.Srv().Store.Warehouse().FilterByOprion(option)
	if err != nil {
		return nil, store.AppErrorFromDatabaseLookupError("WarehouseByoption", "app.warehouse.error_finding_warehouses_by_option.app_error", err)
	}

	return warehouses, nil
}

// WarehouseByStockID returns a warehouse that owns the given stock
func (a *AppWarehouse) WarehouseByStockID(stockID string) (*warehouse.WareHouse, *model.AppError) {
	warehouse, err := a.Srv().Store.Warehouse().WarehouseByStockID(stockID)
	if err != nil {
		return nil, store.AppErrorFromDatabaseLookupError("WarehouseByStockID", "app.warehouse.error_finding_warehouse_by_stock_id.app_error", err)
	}

	return warehouse, nil
}
