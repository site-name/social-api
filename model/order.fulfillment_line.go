package model

import (
	"net/http"

	"github.com/Masterminds/squirrel"
	"github.com/samber/lo"
	"gorm.io/gorm"
)

type FulfillmentLine struct {
	Id            string  `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid();column:Id"`
	OrderLineID   string  `json:"order_line_id" gorm:"type:uuid;column:OrderLineID"`
	FulfillmentID string  `json:"fulfillment_id" gorm:"type:uuid;column:FulfillmentID"`
	Quantity      int     `json:"quantity" gorm:"type:integer;column:Quantity"`
	StockID       *string `json:"stock_id" gorm:"type:uuid;column:StockID"`

	OrderLine *OrderLine `json:"-" db:"-"`
}

func (c *FulfillmentLine) BeforeCreate(_ *gorm.DB) error { return c.IsValid() }
func (c *FulfillmentLine) BeforeUpdate(_ *gorm.DB) error { return c.IsValid() }
func (c *FulfillmentLine) TableName() string             { return FulfillmentLineTableName }

// FulfillmentLineFilterOption is used to build sql queries
type FulfillmentLineFilterOption struct {
	Conditions squirrel.Sqlizer

	FulfillmentOrderID                      squirrel.Sqlizer // INNER JOIN 'Fulfillments' WHERE Fulfillments.OrderID...
	FulfillmentStatus                       squirrel.Sqlizer // INNER JOIN 'Fulfillments' WHERE Fulfillments.Status...
	PrefetchRelatedOrderLine                bool             // this asks to prefetch related order lines of returning fulfillment lines
	PrefetchRelatedOrderLine_ProductVariant bool             // this asks to prefetch related product variants of associated order lines of returning fulfillment lines

	PrefetchRelatedStock bool
}

type FulfillmentLines []*FulfillmentLine

func (f FulfillmentLines) IDs() []string {
	return lo.Map(f, func(item *FulfillmentLine, _ int) string { return item.Id })
}

func (f FulfillmentLines) OrderLineIDs() []string {
	return lo.Map(f, func(item *FulfillmentLine, _ int) string { return item.OrderLineID })

}

// OrderLines returns a slice of order lines attached to every items in f.
//
// NOTE: Make sure the fields `OrderLine` are populated before calling this. If not, the returned slice contains only nil values
func (f FulfillmentLines) OrderLines() OrderLines {
	res := OrderLines{}
	for _, item := range f {
		if item.OrderLine != nil {
			res = append(res, item.OrderLine)
		}
	}

	return res
}

func (f FulfillmentLines) StockIDs() []string {
	res := []string{}
	for _, item := range f {
		if item != nil && item.StockID != nil {
			res = append(res, *item.StockID)
		}
	}

	return res
}

func (f *FulfillmentLine) IsValid() *AppError {
	if !IsValidId(f.OrderLineID) {
		return NewAppError("FulfillmentLine.IsValid", "model.fulfillment_line.is_valid.order_line_id.app_error", nil, "please provide valid order line id", http.StatusBadRequest)
	}
	if !IsValidId(f.FulfillmentID) {
		return NewAppError("FulfillmentLine.IsValid", "model.fulfillment_line.is_valid.fulfillment_id.app_error", nil, "please provide valid fullfillment id", http.StatusBadRequest)
	}
	if f.StockID != nil && !IsValidId(*f.StockID) {
		return NewAppError("FulfillmentLine.IsValid", "model.fulfillment_line.is_valid.stock_id.app_error", nil, "please provide valid stock id", http.StatusBadRequest)
	}
	if f.Quantity <= 0 {
		return NewAppError("FulfillmentLine.IsValid", "model.fulfillment_line.is_valid.quantity.app_error", nil, "quantity must be greater than 0", http.StatusBadRequest)
	}

	return nil
}
