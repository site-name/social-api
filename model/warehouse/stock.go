package warehouse

import (
	"strings"

	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model/order"
	"github.com/sitename/sitename/model/product_and_discount"
)

type Stock struct {
	Id               string `json:"id"`
	CreateAt         int64  `json:"create_at"`
	WarehouseID      string `json:"warehouse_id"`       // NOT NULL
	ProductVariantID string `json:"product_variant_id"` // NOT NULL
	Quantity         int    `json:"quantity"`           // DEFAULT 0

	*WareHouse                           `json:"-" db:"-"` // this foreign field is populated with select related data
	*product_and_discount.ProductVariant `json:"-" db:"-"` // this foreign field is populated with select related data
}

// StockFilterOption is used for build squirrel sql queries
type StockFilterOption struct {
	Id               *model.StringFilter //
	WarehouseID      *model.StringFilter //
	ProductVariantID *model.StringFilter //

	// set this to true if you want to lock selected rows for update.
	// This add `FOR UPDATE` to the end of sql queries
	LockForUpdate bool
	// add something after `FOR UPDATE` to the end of sql queries, to tell the database to lock specific rows instead of both selecting rows and foreign rows
	//
	// E.g:  ForUpdateOf: "Warehouse" => FOR UPDATE OF Warehouse.
	//
	// NOTE: Remember to set `LockForUpdate` property to true before setting this.
	ForUpdateOf string

	// set this if you want to make use of `GetForCountryAndChannel`
	ForCountryAndChannel *StockFilterForCountryAndChannel
}

type Stocks []*Stock

func (s Stocks) IDs() []string {
	res := []string{}
	for _, item := range s {
		if item != nil {
			res = append(res, item.Id)
		}
	}

	return res
}

type StockFilterForCountryAndChannel struct {
	CountryCode      string
	ChannelSlug      string
	WarehouseID      string
	ProductVariantID string
	ProductID        string
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
	if s.Quantity < 0 {
		s.Quantity = 0
	}
}

type InsufficientStockData struct {
	Variant           product_and_discount.ProductVariant // Product variant ID
	OrderLine         *order.OrderLine                    // OrderLine id
	WarehouseID       *string
	AvailableQuantity *int
}

// InsufficientStock is an error indicating stock is insufficient
type InsufficientStock struct {
	Items []*InsufficientStockData
}

func (i *InsufficientStock) Error() string {
	var builder strings.Builder

	builder.WriteString("Insufficient stock for ")
	for idx, item := range i.Items {
		builder.WriteString(item.Variant.String())
		if idx == 0 {
			continue
		}
		if idx == len(i.Items)-1 {
			break
		}
		builder.WriteString(", ")
	}

	return builder.String()
}
