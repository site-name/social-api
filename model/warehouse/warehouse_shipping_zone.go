package warehouse

import "github.com/sitename/sitename/model"

// WarehouseShippingZone represents relationships between warehouses and shipping zones (m2m)
type WarehouseShippingZone struct {
	Id             string `json:"id"`
	WarehouseID    string `json:"warehouse_id"`
	ShippingZoneID string `json:"shipping_zone_id"`
}

// WarehouseShippingZoneFilterOption is used to build squirrel sql queries
type WarehouseShippingZoneFilterOption struct {
	WarehouseID    *model.StringFilter
	ShippingZoneID *model.StringFilter
}

func (w *WarehouseShippingZone) IsValid() *model.AppError {
	outer := model.CreateAppErrorForModel(
		"model.warehouse_shipping_zone.is_valid.%s.app_error",
		"warehouse_shipping_zone_id=",
		"WarehouseShippingZone.IsValid",
	)
	if !model.IsValidId(w.Id) {
		return outer("id", nil)
	}
	if !model.IsValidId(w.WarehouseID) {
		return outer("warehouse_id", &w.Id)
	}
	if !model.IsValidId(w.Id) {
		return outer("shipping_zone_id", &w.Id)
	}

	return nil
}

func (w *WarehouseShippingZone) PreSave() {
	if w.Id == "" {
		w.Id = model.NewId()
	}
}

func (w *WarehouseShippingZone) ToJson() string {
	return model.ModelToJson(w)
}
