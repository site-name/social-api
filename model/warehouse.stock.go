package model

import (
	"github.com/Masterminds/squirrel"
	"github.com/samber/lo"
)

type Stock struct {
	Id               string `json:"id"`
	CreateAt         int64  `json:"create_at"`
	WarehouseID      string `json:"warehouse_id"`       // NOT NULL
	ProductVariantID string `json:"product_variant_id"` // NOT NULL
	Quantity         int    `json:"quantity"`           // DEFAULT 0

	AvailableQuantity int             `json:"-" db:"-"` // this field will be populated in same queries
	warehouse         *WareHouse      `db:"-"`          // this foreign field is populated with select related data
	productVariant    *ProductVariant `db:"-"`          // this foreign field is populated with select related data
}

// StockFilterForChannelOption is used by a filter function at store/sqlstore/channel/channel_store.go
type StockFilterForChannelOption struct {
	ChannelID string

	Id               squirrel.Sqlizer // WHERE Id ...
	WarehouseID      squirrel.Sqlizer // WHERE WarehouseID ...
	ProductVariantID squirrel.Sqlizer // WHERE ProductVariantID ...

	SelectRelatedProductVariant bool // inner join ProductVariants and attachs them to returning stocks

	ReturnQueryOnly bool // if true, only the squirrel query will be returned, no execution will be performed
}

// StockFilterOption is used for build squirrel sql queries
type StockFilterOption struct {
	Id               squirrel.Sqlizer //
	WarehouseID      squirrel.Sqlizer //
	ProductVariantID squirrel.Sqlizer //

	Warehouse_ShippingZone_countries squirrel.Sqlizer // INNER JOIN Warehouses ON ... INNER JOIN WarehouseShippingZones ON ... INNER JOIN ShippingZones ON ... WHERE ShippingZones.Countries ...
	Warehouse_ShippingZone_ChannelID squirrel.Sqlizer // INNER JOIN Warehouses ON ... INNER JOIN WarehouseShippingZones ON ... INNER JOIN ShippingZones ON ... INNER JOIN ShippingZoneChannels WHERE ShippingZoneChannels.ChannelID ...

	SelectRelatedProductVariant bool // inner join ProductVariants and attachs them to returning stocks
	SelectRelatedWarehouse      bool // inner join Warehouses and attachs them to returning stocks

	AnnotateAvailabeQuantity bool // if true, store selects another column: `Stocks.Quantity - COALESCE(SUM(Allocations.QuantityAllocated), 0) AS AvailableQuantity`

	// set this to true if you want to lock selected rows for update.
	// This add `FOR UPDATE` to the end of sql queries
	LockForUpdate bool
	// adds something after `FOR UPDATE` to the end of sql queries.
	// It tells the database to lock accesses to specific rows instead of both selecting rows and relative rows (foreign key rows)
	//
	// E.g:  ForUpdateOf: "Warehouses" results in `FOR UPDATE OF Warehouses`.
	//
	// NOTE: Remember to set `LockForUpdate` to true before setting this.
	ForUpdateOf string
}

type StockFilterForCountryAndChannel struct {
	CountryCode      string
	ChannelSlug      string
	WarehouseID      string
	ProductVariantID string
	ProductID        string

	// additional fields
	Id                     squirrel.Sqlizer
	WarehouseIDFilter      squirrel.Sqlizer
	ProductVariantIDFilter squirrel.Sqlizer

	AnnotateAvailabeQuantity bool // if true, store selects another column: `Stocks.Quantity - COALESCE(SUM(Allocations.QuantityAllocated), 0) AS AvailableQuantity`

	// set this to true if you want to lock selected rows for update.
	// This add `FOR UPDATE` to the end of sql queries
	LockForUpdate bool
	// adds something after `FOR UPDATE` to the end of sql queries.
	// It tells the database to lock accesses to specific rows instead of both selecting rows and relative rows (foreign key rows)
	//
	// E.g:  ForUpdateOf: "Warehouses" results in `FOR UPDATE OF Warehouses`.
	//
	// NOTE: Remember to set `LockForUpdate` to true before setting this.
	ForUpdateOf string
}

type Stocks []*Stock

// IDs returns a slice of ids of stocks contained in current `Stocks`
func (s Stocks) IDs() []string {
	return lo.Map(s, func(st *Stock, _ int) string { return st.Id })
}

func (s Stocks) DeepCopy() Stocks {
	return lo.Map(s, func(st *Stock, _ int) *Stock { return st.DeepCopy() })
}

func (s *Stock) IsValid() *AppError {
	outer := CreateAppErrorForModel(
		"model.stock.is_valid.%s.app_error",
		"stock_id=",
		"Stock.IsValid",
	)
	if !IsValidId(s.Id) {
		return outer("id", nil)
	}
	if s.CreateAt == 0 {
		return outer("create_at", &s.Id)
	}
	if !IsValidId(s.WarehouseID) {
		return outer("warehouse_id", &s.Id)
	}
	if !IsValidId(s.ProductVariantID) {
		return outer("product_variant_id", &s.Id)
	}

	return nil
}

func (s *Stock) ToJSON() string {
	return ModelToJson(s)
}

func (s *Stock) PreSave() {
	if s.Id == "" {
		s.Id = NewId()
	}
	s.CreateAt = GetMillis()
	s.commonPre()
}

func (s *Stock) commonPre() {
	if s.Quantity < 0 {
		s.Quantity = 0
	}
}

func (s *Stock) PreUpdate() {
	s.commonPre()
}

func (s *Stock) DeepCopy() *Stock {
	res := *s

	if s.warehouse != nil {
		res.warehouse = s.warehouse.DeepCopy()
	}

	if s.productVariant != nil {
		res.productVariant = s.productVariant.DeepCopy()
	}

	return &res
}

func (s *Stock) GetWarehouse() *WareHouse {
	return s.warehouse
}

func (s *Stock) SetWarehouse(w *WareHouse) {
	s.warehouse = w
}

func (s *Stock) GetProductVariant() *ProductVariant {
	return s.productVariant
}

func (s *Stock) SetProductVariant(p *ProductVariant) {
	s.productVariant = p
}
