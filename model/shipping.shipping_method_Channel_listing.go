package model

import (
	"strings"

	"github.com/Masterminds/squirrel"
	"github.com/samber/lo"
	"github.com/site-name/decimal"
	goprices "github.com/site-name/go-prices"
	"golang.org/x/text/currency"
)

type ShippingMethodChannelListing struct {
	Id                      string           `json:"id"`
	ShippingMethodID        string           `json:"shipping_method_id"`
	ChannelID               string           `json:"channel_id"`
	MinimumOrderPriceAmount *decimal.Decimal `json:"minimum_order_price_amount"`
	MinimumOrderPrice       *goprices.Money  `json:"minimum_order_price" db:"-"`
	Currency                string           `json:"currency"`
	MaximumOrderPriceAmount *decimal.Decimal `json:"maximum_order_price_amount"`
	MaximumOrderPrice       *goprices.Money  `json:"maximum_order_price" db:"-"`
	Price                   *goprices.Money  `json:"price" db:"-"`
	PriceAmount             *decimal.Decimal `json:"price_amount"`
	CreateAt                int64            `json:"create_at"`
}

// ShippingMethodChannelListingFilterOption is used to build sql queries
type ShippingMethodChannelListingFilterOption struct {
	ShippingMethodID squirrel.Sqlizer
	ChannelID        squirrel.Sqlizer

	ChannelSlug squirrel.Sqlizer // INNER JOIN Channels ON ... WHERE Channels.Slug ...

	ShippingMethod_ShippingZoneID_Inner squirrel.Sqlizer // INNER JOIN ShippingMethods ON ... INNER JOIN ShippingZones ON ... WHERE ShippingZones.Id ...
}

type ShippingMethodChannelListings []*ShippingMethodChannelListing

func (ss ShippingMethodChannelListings) IDs() []string {
	return lo.Map(ss, func(s *ShippingMethodChannelListing, _ int) string { return s.Id })
}

func (ss ShippingMethodChannelListings) ShippingMethodIDs() []string {
	return lo.Map(ss, func(s *ShippingMethodChannelListing, _ int) string { return s.ShippingMethodID })
}

func (s *ShippingMethodChannelListing) IsValid() *AppError {
	outer := CreateAppErrorForModel(
		"model.shipping_method_channel_listing.is_valid.%s.app_error",
		"shipping_method_channel_listing_id=",
		"ShippingMethodChannelListing.IsValid",
	)
	if !IsValidId(s.Id) {
		return outer("id", nil)
	}
	if s.CreateAt == 0 {
		return outer("create_at", &s.Id)
	}
	if !IsValidId(s.ShippingMethodID) {
		return outer("shipping_method_id", &s.Id)
	}
	if !IsValidId(s.ChannelID) {
		return outer("channel_id", &s.Id)
	}
	if unit, err := currency.ParseISO(s.Currency); err != nil || !strings.EqualFold(unit.String(), s.Currency) {
		return outer("currency", &s.Id)
	}

	return nil
}

// PopulateNonDbFields populates non db fields of shipping method channel listing
func (s *ShippingMethodChannelListing) PopulateNonDbFields() {
	if s.MinimumOrderPriceAmount != nil {
		s.MinimumOrderPrice = &goprices.Money{
			Amount:   *s.MinimumOrderPriceAmount,
			Currency: s.Currency,
		}
	}
	if s.MaximumOrderPriceAmount != nil {
		s.MaximumOrderPrice = &goprices.Money{
			Amount:   *s.MaximumOrderPriceAmount,
			Currency: s.Currency,
		}
	}
	if s.PriceAmount != nil {
		s.Price = &goprices.Money{
			Amount:   *s.PriceAmount,
			Currency: s.Currency,
		}
	}
}

func (s *ShippingMethodChannelListing) PreSave() {
	if s.Id == "" {
		s.Id = NewId()
	}
	s.CreateAt = GetMillis()
	s.commonPre()
}

func (s *ShippingMethodChannelListing) commonPre() {
	if s.MinimumOrderPriceAmount == nil {
		s.MinimumOrderPriceAmount = &decimal.Zero
	}
	if s.PriceAmount == nil {
		s.PriceAmount = &decimal.Zero
	}
}

func (s *ShippingMethodChannelListing) PreUpdate() {
	s.commonPre()
}

// GetTotal retuns current ShippingMethodChannelListing's Price fields
func (s *ShippingMethodChannelListing) GetTotal() *goprices.Money {
	s.PopulateNonDbFields()
	return s.Price
}

func (s *ShippingMethodChannelListing) ToJSON() string {
	s.PopulateNonDbFields()
	return ModelToJson(s)
}
