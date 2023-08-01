package model

import (
	"strings"

	"github.com/Masterminds/squirrel"
	"github.com/site-name/decimal"
	"golang.org/x/text/currency"
	"gorm.io/gorm"
)

type SaleChannelListing struct {
	Id            string           `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid();column:Id"`
	SaleID        string           `json:"sale_id" gorm:"type:uuid;column:SaleID"`
	ChannelID     string           `json:"channel_id" gorm:"type:uuid;column:ChannelID"`
	DiscountValue *decimal.Decimal `json:"discount_value" gorm:"default:0;column:DiscountValue"` // default decimal(0)
	Currency      string           `json:"currency" gorm:"type:varchar(3);column:Currency"`
	CreateAt      int64            `json:"create_at" gorm:"autoCreateTime:milli;column:CreateAt"`

	channel *Channel `json:"-"` // this field gets populated in queries that ask for select related channel
}

func (c *SaleChannelListing) BeforeCreate(_ *gorm.DB) error { c.commonPre(); return c.IsValid() }
func (c *SaleChannelListing) BeforeUpdate(_ *gorm.DB) error {
	c.commonPre()
	c.CreateAt = 0 // prevent upda
	return c.IsValid()
}
func (c *SaleChannelListing) TableName() string     { return SaleChannelListingTableName }
func (s *SaleChannelListing) GetChannel() *Channel  { return s.channel }
func (s *SaleChannelListing) SetChannel(c *Channel) { s.channel = c }

func (s *SaleChannelListing) DeepCopy() *SaleChannelListing {
	if s == nil {
		return nil
	}

	res := *s
	if s.DiscountValue != nil {
		res.DiscountValue = NewPrimitive(*s.DiscountValue)
	}
	if s.channel != nil {
		res.channel = s.channel.DeepCopy()
	}

	return &res
}

type SaleChannelListingFilterOption struct {
	Conditions           squirrel.Sqlizer
	SelectRelatedChannel bool
}

func (s *SaleChannelListing) commonPre() {
	if s.DiscountValue == nil || s.DiscountValue.LessThan(decimal.Zero) {
		s.DiscountValue = &decimal.Zero
	}
	s.Currency = strings.ToUpper(s.Currency)
}

func (s *SaleChannelListing) IsValid() *AppError {
	outer := CreateAppErrorForModel(
		"model.sale_channel_listing.is_valid.%s.app_error",
		"sale_channel_listing_id=",
		"SaleChannelListing.IsValid",
	)
	if !IsValidId(s.SaleID) {
		return outer("sale_id", &s.Id)
	}
	if !IsValidId(s.ChannelID) {
		return outer("channel_id", &s.Id)
	}
	if unit, err := currency.ParseISO(s.Currency); err == nil || !strings.EqualFold(unit.String(), s.Currency) {
		return outer("currency", &s.Id)
	}
	err := ValidateDecimal("SaleChannelListing.IsValid.DiscountValue", s.DiscountValue, DECIMAL_TOTAL_DIGITS_ALLOWED, DECIMAL_MAX_DECIMAL_PLACES_ALLOWED)
	if err != nil {
		return err
	}

	return nil
}
