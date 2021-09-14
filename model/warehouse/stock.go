package warehouse

import (
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model/product_and_discount"
)

type Stock struct {
	Id               string `json:"id"`
	CreateAt         int64  `json:"create_at"`
	WarehouseID      string `json:"warehouse_id"`       // NOT NULL
	ProductVariantID string `json:"product_variant_id"` // NOT NULL
	Quantity         int    `json:"quantity"`           // DEFAULT 0

	AvailableQuantity int                                  `json:"-" db:"-"` // this field will be populated in same queries
	Warehouse         *WareHouse                           `json:"-" db:"-"` // this foreign field is populated with select related data
	ProductVariant    *product_and_discount.ProductVariant `json:"-" db:"-"` // this foreign field is populated with select related data
}

// StockFilterForChannelOption is used by a filter function at store/sqlstore/channel/channel_store.go
type StockFilterForChannelOption struct {
	ChannelSlug string

	Id               *model.StringFilter // WHERE Id ...
	WarehouseID      *model.StringFilter // WHERE WarehouseID ...
	ProductVariantID *model.StringFilter // WHERE ProductVariantID ...

	SelectRelatedProductVariant bool // inner join ProductVariants and attachs them to returning stocks
	// SelectRelatedWarehouse      bool // inner join Warehouses and attachs them to returning stocks
}

// StockFilterOption is used for build squirrel sql queries
type StockFilterOption struct {
	Id               *model.StringFilter //
	WarehouseID      *model.StringFilter //
	ProductVariantID *model.StringFilter //

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

// StockFilterForCountryAndChannel is used in specific filter function located at store/sqlstore/channel/channel_store.go
type StockFilterForCountryAndChannel struct {
	CountryCode      string
	ChannelSlug      string
	WarehouseID      string
	ProductVariantID string
	ProductID        string

	// additional fields
	Id                     *model.StringFilter //
	WarehouseIDFilter      *model.StringFilter //
	ProductVariantIDFilter *model.StringFilter //

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
	res := []string{}
	for _, item := range s {
		if item != nil {
			res = append(res, item.Id)
		}
	}

	return res
}

func (s *Stock) IsValid() *model.AppError {
	outer := model.CreateAppErrorForModel(
		"model.stock.is_valid.%s.app_error",
		"stock_id=",
		"Stock.IsValid",
	)
	if !model.IsValidId(s.Id) {
		return outer("id", nil)
	}
	if s.CreateAt == 0 {
		return outer("create_at", &s.Id)
	}
	if !model.IsValidId(s.WarehouseID) {
		return outer("warehouse_id", &s.Id)
	}
	if !model.IsValidId(s.ProductVariantID) {
		return outer("product_variant_id", &s.Id)
	}

	return nil
}

func (s *Stock) ToJson() string {
	return model.ModelToJson(s)
}

func (s *Stock) PreSave() {
	if s.Id == "" {
		s.Id = model.NewId()
	}
	s.CreateAt = model.GetMillis()
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
