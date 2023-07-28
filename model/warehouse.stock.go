package model

import (
	"net/http"

	"github.com/Masterminds/squirrel"
	"github.com/samber/lo"
	"gorm.io/gorm"
)

type StockAvailability string

func (s StockAvailability) IsValid() bool {
	return s == StockAvailabilityInStock || s == StockAvailabilityOutOfStock
}

const (
	StockAvailabilityInStock    StockAvailability = "in_stock"
	StockAvailabilityOutOfStock StockAvailability = "out_of_stock"
)

type Stock struct {
	Id               string `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid();column:Id"`
	CreateAt         int64  `json:"create_at" gorm:"type:bigint;column:CreateAt;autoCreateTime:milli"`
	WarehouseID      string `json:"warehouse_id" gorm:"type:uuid;column:WarehouseID;index:warehouseid_productvariantid_key"`            // NOT NULL
	ProductVariantID string `json:"product_variant_id" gorm:"type:uuid;column:ProductVariantID;index:warehouseid_productvariantid_key"` // NOT NULL
	Quantity         int    `json:"quantity" gorm:"check:Quantity >= 0;column:Quantity"`                                                // DEFAULT 0

	AvailableQuantity int             `json:"-" gorm:"-"` // this field will be populated in same queries
	warehouse         *WareHouse      `gorm:"-"`          // this foreign field is populated with select related data
	productVariant    *ProductVariant `gorm:"-"`          // this foreign field is populated with select related data
}

func (s *Stock) GetWarehouse() *WareHouse            { return s.warehouse }
func (s *Stock) SetWarehouse(w *WareHouse)           { s.warehouse = w }
func (s *Stock) GetProductVariant() *ProductVariant  { return s.productVariant }
func (s *Stock) SetProductVariant(p *ProductVariant) { s.productVariant = p }
func (c *Stock) BeforeCreate(_ *gorm.DB) error       { c.commonPre(); return c.IsValid() }
func (c *Stock) BeforeUpdate(_ *gorm.DB) error {
	c.commonPre()
	c.CreateAt = 0 // prevent updating
	return c.IsValid()
}
func (c *Stock) TableName() string { return StockTableName }

func (s *Stock) IsValid() *AppError {
	if !IsValidId(s.WarehouseID) {
		return NewAppError("Stock.IsValid", "model.stock.is_valid.warehouse_id.app_error", nil, "please provide valid warehouse id", http.StatusBadRequest)
	}
	if !IsValidId(s.ProductVariantID) {
		return NewAppError("Stock.IsValid", "model.stock.is_valid.product_variant_id.app_error", nil, "please provide valid product variant id", http.StatusBadRequest)
	}

	return nil
}

// StockFilterForChannelOption is used by a filter function at store/sqlstore/channel/channel_store.go
type StockFilterForChannelOption struct {
	ChannelID  string
	Conditions squirrel.Sqlizer

	SelectRelatedProductVariant bool // inner join ProductVariants and attachs them to returning stocks

	ReturnQueryOnly bool // if true, only the squirrel query will be returned, no execution will be performed
}

// StockFilterOption is used for build squirrel sql queries
type StockFilterOption struct {
	Conditions squirrel.Sqlizer // all stock's native field lookup should be put here

	Warehouse_ShippingZone_countries squirrel.Sqlizer // INNER JOIN Warehouses ON ... INNER JOIN WarehouseShippingZones ON ... INNER JOIN ShippingZones ON ... WHERE ShippingZones.Countries ...
	Warehouse_ShippingZone_ChannelID squirrel.Sqlizer // INNER JOIN Warehouses ON ... INNER JOIN WarehouseShippingZones ON ... INNER JOIN ShippingZones ON ... INNER JOIN ShippingZoneChannels WHERE ShippingZoneChannels.ChannelID ...

	SelectRelatedProductVariant bool // inner join ProductVariants and attachs them to returning stocks
	SelectRelatedWarehouse      bool // inner join Warehouses and attachs them to returning stocks

	AnnotateAvailabeQuantity bool // if true, store selects another column: `Stocks.Quantity - COALESCE(SUM(Allocations.QuantityAllocated), 0) AS AvailableQuantity`

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

	// set this to true if you want to lock selected rows for update.
	// This add `FOR UPDATE` to the end of sql queries
	// NOTE: only apply if `Transaction` field is set
	LockForUpdate bool
	Transaction   *gorm.DB
	// adds something after `FOR UPDATE` to the end of sql queries.
	// It tells the database to lock accesses to specific rows instead of both selecting rows and relative rows (foreign key rows)
	//
	// E.g:  ForUpdateOf: "Warehouses" results in `FOR UPDATE OF Warehouses`.
	//
	// NOTE: Remember to set `LockForUpdate` to true before setting this.
	ForUpdateOf string

	PaginationValues PaginationValues
}

type StockFilterForCountryAndChannel struct {
	CountryCode      CountryCode
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
	// NOTE: Ony apply if `Transaction` field is set
	LockForUpdate bool
	Transaction   *gorm.DB
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

func (s *Stock) commonPre() {
	if s.Quantity < 0 {
		s.Quantity = 0
	}
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
