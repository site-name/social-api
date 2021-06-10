package product_and_discount

import (
	"io"
	"strings"

	"github.com/shopspring/decimal"
	goprices "github.com/site-name/go-prices"
	"github.com/sitename/sitename/model"
	"golang.org/x/text/currency"
)

type VoucherChannelListing struct {
	Id            string               `json:"id"`
	VoucherID     string               `json:"voucher_id"`
	ChannelID     string               `json:"channel_id"`
	DiscountValue *decimal.NullDecimal `json:"discount_value"`
	Discount      *goprices.Money      `json:"discount,omitempty" db:"-"`
	MinSpent      *goprices.Money      `json:"min_spent,omitempty" db:"-"`
	Currency      string               `json:"currency"`
	MinSpenAmount *decimal.NullDecimal `json:"min_spent_amount"`
}

func (v *VoucherChannelListing) ToJson() string {
	v.Discount = &goprices.Money{
		Amount:   &v.DiscountValue.Decimal,
		Currency: v.Currency,
	}
	v.MinSpent = &goprices.Money{
		Amount:   &v.MinSpenAmount.Decimal,
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
