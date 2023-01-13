package model

import (
	"strings"

	"github.com/Masterminds/squirrel"
	"github.com/site-name/decimal"
	"golang.org/x/text/currency"
)

type SaleChannelListing struct {
	Id            string           `json:"id"`
	SaleID        string           `json:"sale_id"`
	ChannelID     string           `json:"channel_id"`
	DiscountValue *decimal.Decimal `json:"discount_value"` // default decimal(0)
	Currency      string           `json:"currency"`
	CreateAt      int64            `json:"create_at"`

	channel *Channel `db:"-"` // this field gets populated in queries that ask for select related channel
}

func (s *SaleChannelListing) GetChannel() *Channel {
	return s.channel
}

func (s *SaleChannelListing) SetChannel(c *Channel) {
	s.channel = c
}

func (s *SaleChannelListing) DeepCopy() *SaleChannelListing {
	if s == nil {
		return new(SaleChannelListing)
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
	Id        squirrel.Sqlizer
	SaleID    squirrel.Sqlizer
	ChannelID squirrel.Sqlizer

	SelectRelatedChannel bool
}

func (s *SaleChannelListing) PreSave() {
	if s.Id == "" {
		s.Id = NewId()
	}
	s.CreateAt = GetMillis()
	if s.DiscountValue == nil || s.DiscountValue.LessThan(decimal.Zero) {
		s.DiscountValue = &decimal.Zero
	}
	s.Currency = strings.ToUpper(s.Currency)
}

func (s *SaleChannelListing) IsValid() *AppError {
	outer := CreateAppErrorForModel(
		"sale_channel_listing.is_valid.%s.app_error",
		"sale_channel_listing_id=",
		"SaleChannelListing.IsValid",
	)
	if !IsValidId(s.Id) {
		return outer("id", nil)
	}
	if s.CreateAt == 0 {
		return outer("create_at", &s.Id)
	}
	if !IsValidId(s.SaleID) {
		return outer("sale_id", &s.Id)
	}
	if !IsValidId(s.ChannelID) {
		return outer("channel_id", &s.Id)
	}
	if unit, err := currency.ParseISO(s.Currency); err == nil || !strings.EqualFold(unit.String(), s.Currency) {
		return outer("currency", &s.Id)
	}

	return nil
}
