package model

import (
	"net/http"

	"github.com/Masterminds/squirrel"
	"github.com/samber/lo"
	"gorm.io/gorm"
)

type PreorderAllocation struct {
	Id                             string `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid();column:Id"`
	OrderLineID                    string `json:"order_line_id" gorm:"type:uuid;column:OrderLineID;index:orderlineid_productvariantchannellistingid_key"`
	ProductVariantChannelListingID string `json:"product_variant_channel_listing_id" gorm:"type:uuid;column:ProductVariantChannelListingID;index:orderlineid_productvariantchannellistingid_key"`
	Quantity                       int    `json:"quantity" gorm:"check:Quantity >= 0;column:Quantity"`

	orderLine *OrderLine `json:"-" gorm:"-"` // related data popularized in some database calls
}

func (c *PreorderAllocation) BeforeCreate(_ *gorm.DB) error { return c.IsValid() }
func (c *PreorderAllocation) BeforeUpdate(_ *gorm.DB) error { return c.IsValid() }
func (c *PreorderAllocation) TableName() string             { return PreOrderAllocationTableName }

// PreorderAllocationFilterOption is used to build squirrel sql queries
type PreorderAllocationFilterOption struct {
	Conditions squirrel.Sqlizer

	SelectRelated_OrderLine       bool // INNER JOIN OrderLines ON ...
	SelectRelated_OrderLine_Order bool // INNER JOIN Orders ON ...
}

type PreorderAllocations []*PreorderAllocation

func (p PreorderAllocations) IDs() []string {
	return lo.Map(p, func(pr *PreorderAllocation, _ int) string { return pr.Id })
}

func (p *PreorderAllocation) GetOrderLine() *OrderLine { return p.orderLine }

func (p *PreorderAllocation) SetOrderLine(l *OrderLine) { p.orderLine = l }

func (p *PreorderAllocation) IsValid() *AppError {
	if !IsValidId(p.OrderLineID) {
		return NewAppError("PreorderAllocation.IsValid", "model.preorder_allocation.is_valid.orderline_id.app_error", nil, "please provide valid order line id", http.StatusBadRequest)
	}
	if !IsValidId(p.ProductVariantChannelListingID) {
		return NewAppError("PreorderAllocation.IsValid", "model.preorder_allocation.is_valid.orderline_id.app_error", nil, "please provide valid order line id", http.StatusBadRequest)
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
