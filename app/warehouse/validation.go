package warehouse

import (
	"net/http"

	"github.com/mattermost/squirrel"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model_helper"
)

// ValidateWarehouseCount
//
//	Every ShippingZone can be assigned to only one warehouse.
//
// If not there would be issue with automatically selecting stock for operation.
func (a *ServiceWarehouse) ValidateWarehouseCount(shippingZones model.ShippingZoneSlice, instance model.Warehouse) (bool, *model_helper.AppError) {
	shippingZones, appErr := a.srv.Shipping.ShippingZonesByOption(model_helper.ShippingZoneFilterOption{
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
	if instance.ID == "" {
		return false, nil
	}

	return warehouseIdMap[instance.ID], nil
}
