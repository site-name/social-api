package model

import (
	"github.com/Masterminds/squirrel"
	"github.com/samber/lo"
)

type PreorderAllocation struct {
	Id                             string `json:"id"`
	OrderLineID                    string `json:"order_line_id"`
	Quantity                       int    `json:"quantity"`
	ProductVariantChannelListingID string `json:"product_variant_channel_listing_id"`

	orderLine *OrderLine `json:"-" db:"-"` // related data popularized in some database calls
}

// PreorderAllocationFilterOption is used to build squirrel sql queries
type PreorderAllocationFilterOption struct {
	Id                             squirrel.Sqlizer
	OrderLineID                    squirrel.Sqlizer
	Quantity                       squirrel.Sqlizer
	ProductVariantChannelListingID squirrel.Sqlizer

	SelectRelated_OrderLine       bool // INNER JOIN OrderLines ON ...
	SelectRelated_OrderLine_Order bool // INNER JOIN Orders ON ...
}

type PreorderAllocations []*PreorderAllocation

func (p PreorderAllocations) IDs() []string {
	return lo.Map(p, func(pr *PreorderAllocation, _ int) string { return pr.Id })
}

func (p *PreorderAllocation) GetOrderLine() *OrderLine {
	return p.orderLine
}

func (p *PreorderAllocation) SetOrderLine(l *OrderLine) {
	p.orderLine = l
}

func (p *PreorderAllocation) PreSave() {
	if !IsValidId(p.Id) {
		p.Id = NewId()
	}
}

func (p *PreorderAllocation) IsValid() *AppError {
	outer := CreateAppErrorForModel(
		"model.preorder_allocation.is_valid.%s.app_error",
		"oreorder_allocation_id=",
		"PreorderAllocation.IsValid",
	)

	if !IsValidId(p.Id) {
		return outer("id", nil)
	}
	if !IsValidId(p.OrderLineID) {
		return outer("order_line_id", &p.Id)
	}
	if !IsValidId(p.Id) {
		return outer("product_variant_channel_listing_id", &p.Id)
	}

	return nil
}

func (p *PreorderAllocation) DeepCopy() *PreorderAllocation {
	res := *p

	if p.orderLine != nil {
		res.orderLine = p.orderLine.DeepCopy()
	}
	return &res
}
