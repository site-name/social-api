package product_and_discount

import (
	"io"
	"strings"

	"github.com/shopspring/decimal"
	"github.com/sitename/sitename/model"
	"golang.org/x/text/currency"
)

type VoucherChannelListing struct {
	Id            string           `json:"id"`
	VoucherID     string           `json:"voucher_id"`
	ChannelID     string           `json:"channel_id"`
	DiscountValue *decimal.Decimal `json:"discount_value"`
	Discount      *model.Money     `json:"discount,omitempty" db:"_"`
	MinSpent      *model.Money     `json:"min_spent,omitempty" db:"_"`
	Currency      string           `json:"currency"`
	MinSpenAmount *decimal.Decimal `json:"min_spent_amount"`
}

func (v *VoucherChannelListing) ToJson() string {
	v.Discount = &model.Money{
		Amount:   v.DiscountValue,
		Currency: v.Currency,
	}
	v.MinSpent = &model.Money{
		Amount:   v.MinSpenAmount,
		Currency: v.Currency,
	}
	return model.ModelToJson(v)
}

func VoucherChannelListingFromJson(data io.Reader) *VoucherChannelListing {
	var vcl VoucherChannelListing
	model.ModelFromJson(&vcl, data)
	return &vcl
}

func (v *VoucherChannelListing) IsValid() *model.AppError {
	outer := model.CreateAppErrorForModel(
		"model.voucher_channel_listing.is_valid.%s.app_error",
		"voucher_channel_listing_id=",
		"VoucherChannelListing.IsValid",
	)
	if !model.IsValidId(v.Id) {
		return outer("id", nil)
	}
	if !model.IsValidId(v.VoucherID) {
		return outer("voucher_id", &v.Id)
	}
	if !model.IsValidId(v.ChannelID) {
		return outer("channel_id", &v.Id)
	}
	if unit, err := currency.ParseISO(v.Currency); err != nil || !strings.EqualFold(unit.String(), v.Currency) {
		return outer("currency", &v.Id)
	}

	return nil
}
