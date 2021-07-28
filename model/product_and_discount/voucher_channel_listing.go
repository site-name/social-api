package product_and_discount

import (
	"strings"

	"github.com/shopspring/decimal"
	goprices "github.com/site-name/go-prices"
	"github.com/sitename/sitename/model"
	"golang.org/x/text/currency"
)

type VoucherChannelListing struct {
	Id            string           `json:"id"`
	CreateAt      int64            `json:"create_at"` // this field is for ordering
	VoucherID     string           `json:"voucher_id"`
	ChannelID     string           `json:"channel_id"`
	DiscountValue *decimal.Decimal `json:"discount_value"`
	Discount      *goprices.Money  `json:"discount,omitempty" db:"-"`
	MinSpent      *goprices.Money  `json:"min_spent,omitempty" db:"-"`
	Currency      string           `json:"currency"`
	MinSpenAmount *decimal.Decimal `json:"min_spent_amount"`
}

func (v *VoucherChannelListing) PreSave() {
	if v.Id == "" {
		v.Id = model.NewId()
	}
	if v.CreateAt == 0 {
		v.CreateAt = model.GetMillis()
	}
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
	if v.CreateAt == 0 {
		return outer("create_at", &v.Id)
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
