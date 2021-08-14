package order

import "github.com/sitename/sitename/model/product_and_discount"

type OrderLineData struct {
	Line        OrderLine
	Quantity    uint
	Variant     *product_and_discount.ProductVariant // can be nil
	Replace     bool                                 // default false
	WarehouseID string                               // can be empty
}

type FulfillmentLineData struct {
	Line     FulfillmentLine
	Quantity uint
	Replace  bool // default false
}
