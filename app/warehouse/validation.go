package warehouse

import (
	"net/http"

	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model/shipping"
	"github.com/sitename/sitename/model/warehouse"
)

// ValidateWarehouseCount
//	Every ShippingZone can be assigned to only one warehouse.
//
// If not there would be issue with automatically selecting stock for operation.
func (a *AppWarehouse) ValidateWarehouseCount(shippingZones shipping.ShippingZones, instance *warehouse.WareHouse) (bool, *model.AppError) {
	shippingZones, appErr := a.ShippingApp().ShippingZonesByOption(&shipping.ShippingZoneFilterOption{
		Id: &model.StringFilter{
			StringOption: &model.StringOption{
				In: shippingZones.IDs(),
			},
		},
		WarehouseID: &model.StringFilter{
			StringOption: &model.StringOption{
				NULL: model.NewBool(false),
			},
		},
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
