package model_helper

import (
	"net/http"

	"github.com/gosimple/slug"
	"github.com/mattermost/squirrel"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/modules/model_types"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

func AllocationPreSave(allocation *model.Allocation) {
	if allocation.ID == "" {
		allocation.ID = NewId()
	}
	if allocation.CreatedAt == 0 {
		allocation.CreatedAt = GetMillis()
	}
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
	OrderLineOrderID               qm.QueryMod
	Preloads                       []string
	AnnotateStockAvailableQuantity bool
}

var AllocationAnnotationKeys = struct {
	AvailableQuantity string
}{
	AvailableQuantity: "available_quantity",
}

func WarehousePreSave(w *model.Warehouse) {
	if w.ID == "" {
		w.ID = NewId()
	}
	if w.CreatedAt == 0 {
		w.CreatedAt = GetMillis()
	}
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
	ShippingZoneCountries qm.QueryMod // INNER JOIN WarehouseShippingZones ON ... INNER JOIN ShippingZones ON ... WHERE ShippingZones.Countries ...
	ShippingZoneId        qm.QueryMod // INNER JOIN WarehouseShippingZones ON ... WHERE WarehouseShippingZones.ShippingZoneID ...
	Search                string

	Preloads []string
}

type CustomAllocation struct {
	model.Allocation
}

type StockFilterOptionsForCountryAndChannel struct {
	CountryCode model.CountryCode
	ChannelSlug string
}

type StockFilterVariantStocksForCountryFilterOptions struct {
	StockFilterOptionsForCountryAndChannel
	ProductVariantID string
}

type StockFilterProductStocksForCountryAndChannelFilterOptions struct {
	StockFilterOptionsForCountryAndChannel
	ProductID string
}

type StockFilterOption struct {
	CommonQueryOptions

	Warehouse_ShippingZone_countries qm.QueryMod // INNER JOIN Warehouses ON ... INNER JOIN WarehouseShippingZones ON ... INNER JOIN ShippingZones ON ... WHERE ShippingZones.Countries ...
	Warehouse_ShippingZone_ChannelID qm.QueryMod // INNER JOIN Warehouses ON ... INNER JOIN WarehouseShippingZones ON ... INNER JOIN ShippingZones ON ... INNER JOIN ShippingZoneChannels WHERE ShippingZoneChannels.ChannelID ...

	AnnotateAvailableQuantity bool // if true, store selects another column: `Stocks.Quantity - COALESCE(SUM(Allocations.QuantityAllocated), 0) AS AvailableQuantity`

	// NOTE: If Set, store use OR ILIKEs to check this value against:
	//
	// relevant product of this stock's name (INNER JOIN ProductVariants ON ... INNER JOIN Products ON ... WHERE Products.Name ...),
	//
	// relevant product variant's name (INNER JOIN ProductVariants ON ... WHERE ProductVariants.Name ...),
	//
	// relevent warehouse's name (INNER JOIN Warehouses ON ... WHERE Warehouses.Name ...),
	//
	// company name of relevent address of relevent warehouse of this stock (INNER JOIN Warehouses ON ... INNER JOIN Addresses ON ... WHERE Addresses.CompanyName ...)
	Search string

	Preloads []string
}

type StockFilterForChannelOption struct {
	Conditions squirrel.Sqlizer
	ChannelID  string
}

type PreorderAllocationFilterOption struct {
	CommonQueryOptions
	Preloads []string
}

var StockAnnotationKeys = struct {
	AvailableQuantity string
}{
	AvailableQuantity: "available_quantity",
}

func StockPreSave(s *model.Stock) {
	if s.ID == "" {
		s.ID = NewId()
	}
	if s.CreatedAt == 0 {
		s.CreatedAt = GetMillis()
	}
}

func StockIsValid(s model.Stock) *AppError {
	if !IsValidId(s.ID) {
		return NewAppError("StockIsValid", "model.stock.is_valid.id.app_error", nil, "please provide valid id", http.StatusBadRequest)
	}
	if !IsValidId(s.ProductVariantID) {
		return NewAppError("StockIsValid", "model.stock.is_valid.product_variant_id.app_error", nil, "please provide valid product variant id", http.StatusBadRequest)
	}
	if !IsValidId(s.WarehouseID) {
		return NewAppError("StockIsValid", "model.stock.is_valid.warehouse_id.app_error", nil, "please provide valid warehouse id", http.StatusBadRequest)
	}
	if s.Quantity < 0 {
		return NewAppError("StockIsValid", "model.stock.is_valid.quantity.app_error", nil, "please provide valid quantity", http.StatusBadRequest)
	}
	if s.CreatedAt <= 0 {
		return NewAppError("StockIsValid", "model.stock.is_valid.created_at.app_error", nil, "please provide valid created at", http.StatusBadRequest)
	}
	return nil
}
