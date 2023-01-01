package model

import "github.com/samber/lo"

type OrderLineData struct {
	Line        OrderLine
	Quantity    int
	Variant     *ProductVariant // can be nil
	Replace     bool            // default false
	WarehouseID *string         // can be nil
}

func (o *OrderLineData) DeepCopy() *OrderLineData {
	if o == nil {
		return &OrderLineData{}
	}

	res := *o
	res.Line = *o.Line.DeepCopy()
	if o.WarehouseID != nil {
		res.WarehouseID = NewString(*o.WarehouseID)
	}
	if o.Variant != nil {
		res.Variant = o.Variant.DeepCopy()
	}

	return &res
}

type OrderLineDatas []*OrderLineData

func (a OrderLineDatas) DeepCopy() []*OrderLineData {
	return lo.Map(a, func(o *OrderLineData, _ int) *OrderLineData { return o.DeepCopy() })
}

func (a OrderLineDatas) Variants() ProductVariants {
	return lo.Map(a, func(o *OrderLineData, _ int) *ProductVariant { return o.Variant })
}

func (a OrderLineDatas) OrderLines() OrderLines {
	return lo.Map(a, func(o *OrderLineData, _ int) *OrderLine { return &o.Line })
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
	return lo.Map(q, func(o *QuantityOrderLine, _ int) *OrderLine { return o.OrderLine })
}
