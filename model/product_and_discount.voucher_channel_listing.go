package model

import (
	"strings"

	"github.com/Masterminds/squirrel"
	"github.com/site-name/decimal"
	goprices "github.com/site-name/go-prices"
	"golang.org/x/text/currency"
)

type VoucherChannelListing struct {
	Id             string           `json:"id"`
	CreateAt       int64            `json:"create_at"` // this field is for ordering
	VoucherID      string           `json:"voucher_id"`
	ChannelID      string           `json:"channel_id"`
	DiscountValue  *decimal.Decimal `json:"discount_value"` // default decimal.Zero
	Currency       string           `json:"currency"`
	MinSpentAmount *decimal.Decimal `json:"min_spent_amount"` // default decimal.Zero
	Discount       *goprices.Money  `json:"discount,omitempty" db:"-"`
	MinSpent       *goprices.Money  `json:"min_spent,omitempty" db:"-"`
}

// VoucherChannelListingFilterOption is mainly used to build sql queries to filter voucher channel listing relationship instances
type VoucherChannelListingFilterOption struct {
	Id        squirrel.Sqlizer
	VoucherID squirrel.Sqlizer
	ChannelID squirrel.Sqlizer
}

func (v *VoucherChannelListing) PreSave() {
	if v.Id == "" {
		v.Id = NewId()
	}
	v.CreateAt = GetMillis()

	v.commonPre()
}

func (v *VoucherChannelListing) commonPre() {
	if v.Discount != nil {
		v.DiscountValue = &v.Discount.Amount
	} else {
		v.DiscountValue = &decimal.Zero
	}
	if v.MinSpent == nil {
		v.MinSpentAmount = &v.MinSpent.Amount
	} else {
		v.MinSpentAmount = &decimal.Zero
	}
	if v.Currency != "" {
		v.Currency = strings.ToUpper(v.Currency)
	} else {
		v.Currency = DEFAULT_CURRENCY
	}
}

func (v *VoucherChannelListing) PopulateNonDbFields() {
	v.MinSpent = &goprices.Money{
		Currency: v.Currency,
		Amount:   *v.MinSpentAmount,
	}
	v.Discount = &goprices.Money{
		Amount:   *v.DiscountValue,
		Currency: v.Currency,
	}
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

func (v *VoucherChannelListing) IsValid() *AppError {
	outer := CreateAppErrorForModel(
		"model.voucher_channel_listing.is_valid.%s.app_error",
		"voucher_channel_listing_id=",
		"VoucherChannelListing.IsValid",
	)
	if !IsValidId(v.Id) {
		return outer("id", nil)
	}
	if v.CreateAt == 0 {
		return outer("create_at", &v.Id)
	}
	if !IsValidId(v.VoucherID) {
		return outer("voucher_id", &v.Id)
	}
	if !IsValidId(v.ChannelID) {
		return outer("channel_id", &v.Id)
	}
	if unit, err := currency.ParseISO(v.Currency); err != nil || !strings.EqualFold(unit.String(), v.Currency) {
		return outer("currency", &v.Id)
	}

	return nil
}
