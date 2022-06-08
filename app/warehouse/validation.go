package warehouse

import (
	"net/http"

	"github.com/Masterminds/squirrel"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model/shipping"
	"github.com/sitename/sitename/model/warehouse"
	"github.com/sitename/sitename/store"
)

// ValidateWarehouseCount
//	Every ShippingZone can be assigned to only one warehouse.
//
// If not there would be issue with automatically selecting stock for operation.
func (a *ServiceWarehouse) ValidateWarehouseCount(shippingZones shipping.ShippingZones, instance *warehouse.WareHouse) (bool, *model.AppError) {
	shippingZones, appErr := a.srv.ShippingService().ShippingZonesByOption(&shipping.ShippingZoneFilterOption{
		Id:                       squirrel.Eq{store.ShippingZoneTableName + ".Id": shippingZones.IDs()},
		WarehouseID:              squirrel.NotEq{store.WarehouseShippingZoneTableName + ".WarehouseID": nil},
		SelectRelatedThroughData: true, // this tells store to populate `RelativeWarehouseIDs` of returning shipping zones
	})
	if appErr != nil {
		if appErr.StatusCode == http.StatusInternalServerError {
			return false, appErr
		}
		shippingZones = []*shipping.ShippingZone{}
	}

	warehouseIDs := shippingZones.RelativeWarehouseIDsFlat(false)

	return len(warehouseIDs) == 0 ||
		(model.IsValidId(instance.Id) && len(warehouseIDs) == 1 && warehouseIDs[0] == instance.Id), nil
}
