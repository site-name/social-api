package warehouse

import (
	"github.com/Masterminds/squirrel"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model/order"
)

type PreorderAllocation struct {
	Id                             string `json:"id"`
	OrderLineID                    string `json:"order_line_id"`
	Quantity                       int    `json:"quantity"`
	ProductVariantChannelListingID string `json:"product_variant_channel_listing_id"`

	OrderLine *order.OrderLine `json:"-" db:"-"` // related data popularized in some database calls
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
	res := []string{}
	for _, item := range p {
		if item != nil {
			res = append(res, item.Id)
		}
	}

	return res
}

func (p *PreorderAllocation) PreSave() {
	if !model.IsValidId(p.Id) {
		p.Id = model.NewId()
	}
}

func (p *PreorderAllocation) IsValid() *model.AppError {
	outer := model.CreateAppErrorForModel(
		"model.preorder_allocation.is_valid.%s.app_error",
		"oreorder_allocation_id=",
		"PreorderAllocation.IsValid",
	)

	if !model.IsValidId(p.Id) {
		return outer("id", nil)
	}
	if !model.IsValidId(p.OrderLineID) {
		return outer("order_line_id", &p.Id)
	}
	if !model.IsValidId(p.Id) {
		return outer("product_variant_channel_listing_id", &p.Id)
	}

	return nil
}

func (p *PreorderAllocation) DeepCopy() *PreorderAllocation {
	res := *p
	return &res
}