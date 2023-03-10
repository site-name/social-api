package model

import (
	"github.com/Masterminds/squirrel"
	"github.com/samber/lo"
)

type FulfillmentLine struct {
	Id            string  `json:"id"`
	OrderLineID   string  `json:"order_line_id"`
	FulfillmentID string  `json:"fulfillment_id"`
	Quantity      int     `json:"quantity"`
	StockID       *string `json:"stock_id"`

	OrderLine *OrderLine `json:"-" db:"-"`
}

// FulfillmentLineFilterOption is used to build sql queries
type FulfillmentLineFilterOption struct {
	Id                 squirrel.Sqlizer
	OrderLineID        squirrel.Sqlizer
	FulfillmentID      squirrel.Sqlizer
	FulfillmentOrderID squirrel.Sqlizer // INNER JOIN 'Fulfillments' WHERE Fulfillments.OrderID...
	FulfillmentStatus  squirrel.Sqlizer // INNER JOIN 'Fulfillments' WHERE Fulfillments.Status...

	PrefetchRelatedOrderLine                bool // this asks to prefetch related order lines of returning fulfillment lines
	PrefetchRelatedOrderLine_ProductVariant bool // this asks to prefetch related product variants of associated order lines of returning fulfillment lines

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
	outer := CreateAppErrorForModel(
		"model.fulfillment_line.is_valid.%s.app_error",
		"fulfillment_line_id=",
		"FulfillmentLine.IsValid",
	)
	if !IsValidId(f.Id) {
		return outer("id", nil)
	}
	if !IsValidId(f.OrderLineID) {
		return outer("order_id", &f.Id)
	}
	if !IsValidId(f.FulfillmentID) {
		return outer("fulfillment_id", &f.Id)
	}
	if f.StockID != nil && !IsValidId(*f.StockID) {
		return outer("stock_id", &f.Id)
	}

	return nil
}

func (f *FulfillmentLine) ToJSON() string {
	return ModelToJson(f)
}

func (f *FulfillmentLine) PreSave() {
	if f.Id == "" {
		f.Id = NewId()
	}
}
