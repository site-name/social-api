package order

import "github.com/sitename/sitename/model/product_and_discount"

type OrderLineData struct {
	Line        OrderLine
	Quantity    int
	Variant     *product_and_discount.ProductVariant // can be nil
	Replace     bool                                 // default false
	WarehouseID *string                              // can be nil
}

type OrderLineDatas []*OrderLineData

func (a OrderLineDatas) DeepCopy() []*OrderLineData {
	res := []*OrderLineData{}
	for _, orderLineData := range a {
		if orderLineData != nil {
			newItem := *orderLineData
			res = append(res, &newItem)
		}
	}

	return res
}

func (a OrderLineDatas) Variants() product_and_discount.ProductVariants {
	res := []*product_and_discount.ProductVariant{}
	for _, item := range a {
		if item != nil && item.Variant != nil {
			res = append(res, item.Variant)
		}
	}

	return res
}

func (a OrderLineDatas) OrderLines() OrderLines {
	res := []*OrderLine{}
	for _, item := range a {
		if item != nil {
			res = append(res, &item.Line)
		}
	}

	return res
}

func (a OrderLineDatas) WarehouseIDs() []string {
	res := []string{}
	for _, item := range a {
		if item != nil && item.WarehouseID != nil {
			res = append(res, *item.WarehouseID)
		}
	}

	return res
}

type FulfillmentLineData struct {
	Line     FulfillmentLine
	Quantity int
	Replace  bool // default false
}

// QuantityOrderLine
type QuantityOrderLine struct {
	Quantity  int
	OrderLine *OrderLine
}

type QuantityOrderLines []*QuantityOrderLine

func (q QuantityOrderLines) OrderLines() OrderLines {
	res := []*OrderLine{}
	for _, item := range q {
		if item != nil {
			res = append(res, item.OrderLine)
		}
	}

	return res
}

// NOTE: ReplicateWarehouseAllocation is identical to warehouse.Allocation
// We re-define this struct sine cycle import the package model/warehouse is not allowed
// You should update this struct definition whenever a change to warehouse.Allocation is made
type ReplicateWarehouseAllocation struct {
	Id                string `json:"id"`
	CreateAt          int64  `json:"create_at"`
	OrderLineID       string `json:"order_ldine_id"`     // NOT NULL
	StockID           string `json:"stock_id"`           // NOT NULL
	QuantityAllocated int    `json:"quantity_allocated"` // default 0

	stock *ReplicateWarehouseStock // this field get populated
}

func (a *ReplicateWarehouseAllocation) GetStock() *ReplicateWarehouseStock {
	return a.stock
}

func (a *ReplicateWarehouseAllocation) SetStock(stock *ReplicateWarehouseStock) {
	a.stock = stock
}

func (a *ReplicateWarehouseAllocation) DeepCopy() *ReplicateWarehouseAllocation {
	res := *a

	if st := a.GetStock(); st != nil {
		res.SetStock(st.DeepCopy())
	}

	return &res
}

// NOTE: ReplicateWarehouseStock is identical to warehouse.Stock
// We re-define this struct sine cycle import the package model/warehouse is not allowed
// You should update this struct definition whenever a change to warehouse.Stock is made
type ReplicateWarehouseStock struct {
	Id               string `json:"id"`
	CreateAt         int64  `json:"create_at"`
	WarehouseID      string `json:"warehouse_id"`       // NOT NULL
	ProductVariantID string `json:"product_variant_id"` // NOT NULL
	Quantity         int    `json:"quantity"`           // DEFAULT 0
}

func (s *ReplicateWarehouseStock) DeepCopy() *ReplicateWarehouseStock {
	res := *s
	return &res
}
