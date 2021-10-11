package warehouse

import (
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
	Id                             *model.StringFilter
	OrderLineID                    *model.StringFilter
	Quantity                       *model.NumberFilter
	ProductVariantChannelListingID *model.StringFilter

	SelectRelated_OrderLine       bool // INNER JOIN OrderLines ON ...
	SelectRelated_OrderLine_Order bool // INNER JOIN Orders ON ...
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
