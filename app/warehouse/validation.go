package warehouse

import (
	"net/http"

	"github.com/Masterminds/squirrel"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/store"
)

// ValidateWarehouseCount
//
//	Every ShippingZone can be assigned to only one warehouse.
//
// If not there would be issue with automatically selecting stock for operation.
func (a *ServiceWarehouse) ValidateWarehouseCount(shippingZones model.ShippingZones, instance *model.WareHouse) (bool, *model.AppError) {
	shippingZones, appErr := a.srv.ShippingService().ShippingZonesByOption(&model.ShippingZoneFilterOption{
		Id:                        squirrel.Eq{store.ShippingZoneTableName + ".Id": shippingZones.IDs()},
		WarehouseID:               squirrel.NotEq{store.WarehouseShippingZoneTableName + ".WarehouseID": nil},
		SelectRelatedWarehouseIDs: true,
	})
	if appErr != nil {
		if appErr.StatusCode == http.StatusInternalServerError {
			return false, appErr
		}
		shippingZones = []*model.ShippingZone{}
	}

	warehouseIdMap := map[string]bool{}
	for _, zone := range shippingZones {
		for _, warehouseId := range zone.RelativeWarehouseIDs {
			warehouseIdMap[warehouseId] = true
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
