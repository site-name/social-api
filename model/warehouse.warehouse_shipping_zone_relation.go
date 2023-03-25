package model

import (
	"github.com/Masterminds/squirrel"
)

// WarehouseShippingZone represents relationships between warehouses and shipping zones (m2m)
type WarehouseShippingZone struct {
	Id             string `json:"id"`
	WarehouseID    string `json:"warehouse_id"`
	ShippingZoneID string `json:"shipping_zone_id"`
}

// WarehouseShippingZoneFilterOption is used to build squirrel sql queries
type WarehouseShippingZoneFilterOption struct {
	WarehouseID    squirrel.Sqlizer
	ShippingZoneID squirrel.Sqlizer
}

func (w *WarehouseShippingZone) IsValid() *AppError {
	outer := CreateAppErrorForModel(
		"model.warehouse_shipping_zone.is_valid.%s.app_error",
		"warehouse_shipping_zone_id=",
		"WarehouseShippingZone.IsValid",
	)
	if !IsValidId(w.Id) {
		return outer("id", nil)
	}
	if !IsValidId(w.WarehouseID) {
		return outer("warehouse_id", &w.Id)
	}
	if !IsValidId(w.Id) {
		return outer("shipping_zone_id", &w.Id)
	}

	return nil
}

func (w *WarehouseShippingZone) PreSave() {
	if w.Id == "" {
		w.Id = NewId()
	}
}

func (w *WarehouseShippingZone) ToJSON() string {
	return ModelToJson(w)
}
