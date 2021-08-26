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
