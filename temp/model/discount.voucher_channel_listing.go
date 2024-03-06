package model

import (
	"net/http"
	"strings"

	"github.com/mattermost/squirrel"
	"github.com/site-name/decimal"
	goprices "github.com/site-name/go-prices"
	"golang.org/x/text/currency"
	"gorm.io/gorm"
)

type VoucherChannelListing struct {
	Id             string           `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid();column:Id"`
	CreateAt       int64            `json:"create_at" gorm:"type:bigint;column:CreateAt;autoCreateTime:milli"` // this field is for ordering
	VoucherID      string           `json:"voucher_id" gorm:"type:uuid;column:VoucherID;index:voucherid_channelid_key"`
	ChannelID      string           `json:"channel_id" gorm:"type:uuid;column:ChannelID;index:voucherid_channelid_key"`
	DiscountValue  *decimal.Decimal `json:"discount_value" gorm:"default:0;column:DiscountValue"` // default decimal.Zero
	Currency       string           `json:"currency" gorm:"type:varchar(5);column:Currency"`
	MinSpentAmount *decimal.Decimal `json:"min_spent_amount" gorm:"default:0;column:MinSpentAmount"` // default decimal.Zero

	Discount *goprices.Money `json:"discount,omitempty" gorm:"-"`
	MinSpent *goprices.Money `json:"min_spent,omitempty" gorm:"-"`
}

func (c *VoucherChannelListing) BeforeCreate(_ *gorm.DB) error { c.commonPre(); return c.IsValid() }
func (c *VoucherChannelListing) BeforeUpdate(_ *gorm.DB) error { c.commonPre(); return c.IsValid() }
func (c *VoucherChannelListing) TableName() string             { return VoucherChannelListingTableName }

// VoucherChannelListingFilterOption is mainly used to build sql queries to filter voucher channel listing relationship instances
type VoucherChannelListingFilterOption struct {
	Conditions squirrel.Sqlizer
}

func (v *VoucherChannelListing) commonPre() {
	if v.Discount != nil {
		v.DiscountValue = &v.Discount.Amount
	} else {
		v.DiscountValue = GetPointerOfValue(decimal.Zero)
	}
	if v.MinSpent == nil {
		v.MinSpentAmount = &v.MinSpent.Amount
	} else {
		v.MinSpentAmount = GetPointerOfValue(decimal.Zero)
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

type VoucherChannelListingList []*VoucherChannelListing

func (vs VoucherChannelListingList) PopulateNonDbFields() {
	for _, v := range vs {
		v.PopulateNonDbFields()
	}
}

func (v *VoucherChannelListing) IsValid() *AppError {
	if !IsValidId(v.VoucherID) {
		return NewAppError("VoucherChannelListing.IsValid", "model.voucher_channel_listing.is_valid.voucher_id.app_error", nil, "please provide valid voucher id", http.StatusBadRequest)
	}
	if !IsValidId(v.ChannelID) {
		return NewAppError("VoucherChannelListing.IsValid", "model.voucher_channel_listing.is_valid.channel_id.app_error", nil, "please provide valid channel id", http.StatusBadRequest)
	}
	if unit, err := currency.ParseISO(v.Currency); err != nil || !strings.EqualFold(unit.String(), v.Currency) {
		return NewAppError("VoucherChannelListing.IsValid", "model.voucher_channel_listing.is_valid.currency.app_error", nil, "please provide valid currency", http.StatusBadRequest)
	}

	return nil
}
