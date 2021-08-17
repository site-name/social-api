package product_and_discount

import (
	"strings"

	"github.com/site-name/decimal"
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

// VoucherChannelListingFilterOption is mainly used to build sql queries to filter voucher channel listing relationship instances
type VoucherChannelListingFilterOption struct {
	Id        *model.StringFilter
	VoucherID *model.StringFilter
	ChannelID *model.StringFilter
}

func (v *VoucherChannelListing) PreSave() {
	if v.Id == "" {
		v.Id = model.NewId()
	}
	v.CreateAt = model.GetMillis()

	v.commonPre()
}

func (v *VoucherChannelListing) commonPre() {
	if v.Discount != nil {
		v.DiscountValue = v.Discount.Amount
	} else {
		v.DiscountValue = &decimal.Zero
	}
	if v.MinSpent == nil {
		v.MinSpenAmount = v.MinSpent.Amount
	} else {
		v.MinSpenAmount = &decimal.Zero
	}
	if v.Currency != "" {
		v.Currency = strings.ToUpper(v.Currency)
	} else {
		v.Currency = model.DEFAULT_CURRENCY
	}
}

func (v *VoucherChannelListing) PopulateNonDbFields() {
	v.MinSpent, _ = goprices.NewMoney(v.MinSpenAmount, v.Currency)
	v.Discount, _ = goprices.NewMoney(v.DiscountValue, v.Currency)
}

func (v *VoucherChannelListing) PreUpdate() {
	v.commonPre()
}

type VoucherChannelListingList []*VoucherChannelListing

func (vs VoucherChannelListingList) PopulateNonDbFields() {
	for _, v := range vs {
		v.PopulateNonDbFields()
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
