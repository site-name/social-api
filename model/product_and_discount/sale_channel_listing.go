package product_and_discount

import "github.com/site-name/decimal"

type SaleChannelListing struct {
	Id            string           `json:"id"`
	SaleID        string           `json:"sale_id"`
	ChannelID     string           `json:"channel_id"`
	DiscountValue *decimal.Decimal `json:"discount_value"`
	Currency      string           `json:"currency"`
}
