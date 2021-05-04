package discount

import (
	"io"

	"github.com/shopspring/decimal"
	"github.com/sitename/sitename/model/checkout"
	"github.com/sitename/sitename/modules/json"
)

type VoucherChannelListing struct {
	Id            string           `json:"id"`
	VoucherID     string           `json:"voucher_id"`
	ChannelID     string           `json:"channel_id"`
	DiscountValue *decimal.Decimal `json:"discount_value"`
	Discount      *checkout.Money  `json:"discount" db:"_"`
	MinSpent      *checkout.Money  `json:"min_spent" db:"_"`
	Currency      string           `json:"currency"`
	MinSpenAmount *decimal.Decimal `json:"min_spent_amount"`
}

func (v *VoucherChannelListing) ToJson() string {
	b, _ := json.JSON.Marshal(v)
	return string(b)
}

func VoucherChannelListingFromJson(data io.Reader) *VoucherChannelListing {
	var vcl VoucherChannelListing
	err := json.JSON.NewDecoder(data).Decode(&vcl)
	if err != nil {
		return nil
	}
	return &vcl
}
