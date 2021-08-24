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
	res := make([]*OrderLineData, len(a))
	for i := range a {
		res[i] = &(*a[i])
	}

	return res
}

func (a OrderLineDatas) Variants() []*product_and_discount.ProductVariant {
	res := []*product_and_discount.ProductVariant{}
	for _, item := range a {
		if item != nil && item.Variant != nil {
			res = append(res, item.Variant)
		}
	}

	return res
}

type FulfillmentLineData struct {
	Line     FulfillmentLine
	Quantity int
	Replace  bool // default false
}
