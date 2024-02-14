package model_helper

import (
	"net/http"

	"github.com/gosimple/slug"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/modules/model_types"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

func AllocationPreSave(allocation *model.Allocation) {
	if allocation.ID == "" {
		allocation.ID = NewId()
	}
	allocation.CreatedAt = GetMillis()
}

func AllocationIsValid(allocation model.Allocation) *AppError {
	if !IsValidId(allocation.ID) {
		return NewAppError("AllocationIsValid", "model.allocation.is_valid.id.app_error", nil, "please provide valid id", http.StatusBadRequest)
	}
	if !IsValidId(allocation.StockID) {
		return NewAppError("AllocationIsValid", "model.allocation.is_valid.stock_id.app_error", nil, "please provide valid stock id", http.StatusBadRequest)
	}
	if !IsValidId(allocation.OrderLineID) {
		return NewAppError("AllocationIsValid", "model.allocation.is_valid.order_line_id.app_error", nil, "please provide valid order line id", http.StatusBadRequest)
	}
	if allocation.QuantityAllocated < 0 {
		return NewAppError("AllocationIsValid", "model.allocation.is_valid.quantity_allocated.app_error", nil, "please provide valid quantity allocated", http.StatusBadRequest)
	}
	return nil
}

type AllocationFilterOption struct {
	CommonQueryOptions
	OrderLineOrderID qm.QueryMod
}

func WarehousePreSave(w *model.Warehouse) {
	if w.ID == "" {
		w.ID = NewId()
	}
	w.CreatedAt = GetMillis()
	WarehouseCommonPre(w)
}

func WarehousePreUpdate(w *model.Warehouse) {
	WarehouseCommonPre(w)
}

func WarehouseCommonPre(w *model.Warehouse) {
	w.Name = SanitizeUnicode(w.Name)
	w.Slug = slug.Make(w.Name)
	if w.ClickAndCollectOption.IsValid() != nil {
		w.ClickAndCollectOption = model.WarehouseClickAndCollectOptionDisabled
	}
	if w.IsPrivate.IsNil() {
		w.IsPrivate = model_types.NewNullBool(true)
	}
}

func WarehouseIsValid(w model.Warehouse) *AppError {
	if !IsValidId(w.ID) {
		return NewAppError("WarehouseIsValid", "model.warehouse.is_valid.id.app_error", nil, "please provide valid id", http.StatusBadRequest)
	}
	if !w.AddressID.IsNil() && !IsValidId(*w.AddressID.String) {
		return NewAppError("WarehouseIsValid", "model.warehouse.is_valid.address_id.app_error", nil, "please provide valid address id", http.StatusBadRequest)
	}
	if !IsValidEmail(w.Email) {
		return NewAppError("WarehouseIsValid", "model.warehouse.is_valid.email.app_error", nil, "please provide valid email", http.StatusBadRequest)
	}
	if w.ClickAndCollectOption.IsValid() != nil {
		return NewAppError("WarehouseIsValid", "model.warehouse.is_valid.click_and_collect_option.app_error", nil, "please provide valid click and collect option", http.StatusBadRequest)
	}
	if w.CreatedAt <= 0 {
		return NewAppError("WarehouseIsValid", "model.warehouse.is_valid.created_at.app_error", nil, "please provide valid created at", http.StatusBadRequest)
	}
	return nil
}

type WarehouseFilterOption struct {
	CommonQueryOptions
	ShippingZoneCountries qm.QueryMod
	ShippingZoneId        qm.QueryMod
	Search                string
}
