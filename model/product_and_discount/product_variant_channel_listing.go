package product_and_discount

import (
	"github.com/shopspring/decimal"
	"github.com/sitename/sitename/model"
)

type ProductVariantChannelListing struct {
	Id              string           `json:"id"`
	VariantID       string           `json:"variant_id"`
	ChannelID       string           `json:"channel_id"`
	Currency        string           `json:"currency"`
	PriceAmount     *decimal.Decimal `json:"price_amount,omitempty"`
	Price           *model.Money     `json:"price" db:"-"`
	CostPriceAmount *decimal.Decimal `json:"cost_price_amount"`
	CostPrice       *model.Money     `json:"cost_price,omitempty" db:"-"`
}
