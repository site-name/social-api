package warehouse

import (
	"net/http"

	"github.com/Masterminds/squirrel"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model_helper"
)

// ValidateWarehouseCount
//
//	Every ShippingZone can be assigned to only one warehouse.
//
// If not there would be issue with automatically selecting stock for operation.
func (a *ServiceWarehouse) ValidateWarehouseCount(shippingZones model.ShippingZones, instance *model.WareHouse) (bool, *model_helper.AppError) {
	shippingZones, appErr := a.srv.ShippingService().ShippingZonesByOption(&model.ShippingZoneFilterOption{
		SelectRelatedWarehouses: true,
		Conditions:              squirrel.Eq{model.ShippingZoneTableName + ".Id": shippingZones.IDs()},
		WarehouseID:             squirrel.NotEq{model.WarehouseShippingZoneTableName + ".WarehouseID": nil},
	})
	if appErr != nil {
		if appErr.StatusCode == http.StatusInternalServerError {
			return false, appErr
		}
		shippingZones = []*model.ShippingZone{}
	}

	warehouseIdMap := map[string]bool{}
	for _, zone := range shippingZones {
		for _, warehouse := range zone.Warehouses {
			warehouseIdMap[warehouse.Id] = true
		}
	}

	if len(warehouseIdMap) == 0 {
		return true, nil
	}
	if len(warehouseIdMap) > 1 {
		return false, nil
	}
	if instance.Id == "" {
		return false, nil
	}

	return warehouseIdMap[instance.Id], nil
}
