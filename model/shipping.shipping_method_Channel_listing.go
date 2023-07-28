package model

import (
	"strings"

	"github.com/Masterminds/squirrel"
	"github.com/samber/lo"
	"github.com/site-name/decimal"
	goprices "github.com/site-name/go-prices"
	"github.com/sitename/sitename/modules/util"
	"golang.org/x/text/currency"
	"gorm.io/gorm"
)

type ShippingMethodChannelListing struct {
	Id                      string           `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid();column:Id"`
	ShippingMethodID        string           `json:"shipping_method_id" gorm:"type:uuid;column:ShippingMethodID;index:shippingmethodid_channelid_key"`
	ChannelID               string           `json:"channel_id" gorm:"type:uuid;column:ChannelID;index:shippingmethodid_channelid_key"`
	MinimumOrderPriceAmount *decimal.Decimal `json:"minimum_order_price_amount" gorm:"default:0;column:MinimumOrderPriceAmount"` // default 0
	Currency                string           `json:"currency" gorm:"type:varchar(5);column:Currency"`
	MaximumOrderPriceAmount *decimal.Decimal `json:"maximum_order_price_amount" gorm:"column:MaximumOrderPriceAmount"`
	PriceAmount             *decimal.Decimal `json:"price_amount" gorm:"default:0;column:PriceAmount"`
	CreateAt                int64            `json:"create_at" gorm:"type:bigint;column:CreateAt;autoCreateTime:milli"`

	MaximumOrderPrice *goprices.Money `json:"maximum_order_price" gorm:"-"`
	Price             *goprices.Money `json:"price" gorm:"-"`
	MinimumOrderPrice *goprices.Money `json:"minimum_order_price" gorm:"-"`
}

func (c *ShippingMethodChannelListing) BeforeCreate(_ *gorm.DB) error {
	c.commonPre()
	return c.IsValid()
}
func (c *ShippingMethodChannelListing) BeforeUpdate(_ *gorm.DB) error {
	c.commonPre()
	c.CreateAt = 0 // prevent update
	return c.IsValid()
}
func (c *ShippingMethodChannelListing) TableName() string {
	return ShippingMethodChannelListingTableName
}

// ShippingMethodChannelListingFilterOption is used to build sql queries
type ShippingMethodChannelListingFilterOption struct {
	Conditions squirrel.Sqlizer

	ChannelSlug                         squirrel.Sqlizer // INNER JOIN Channels ON ... WHERE Channels.Slug ...
	ShippingMethod_ShippingZoneID_Inner squirrel.Sqlizer // INNER JOIN ShippingMethods ON ... INNER JOIN ShippingZones ON ... WHERE ShippingZones.Id ...
}

type ShippingMethodChannelListings []*ShippingMethodChannelListing

func (ss ShippingMethodChannelListings) IDs() util.AnyArray[string] {
	return lo.Map(ss, func(s *ShippingMethodChannelListing, _ int) string { return s.Id })
}

func (ss ShippingMethodChannelListings) ShippingMethodIDs() util.AnyArray[string] {
	return lo.Map(ss, func(s *ShippingMethodChannelListing, _ int) string { return s.ShippingMethodID })
}

func (ss ShippingMethodChannelListings) ChannelIDs() util.AnyArray[string] {
	return lo.Map(ss, func(s *ShippingMethodChannelListing, _ int) string { return s.ChannelID })
}

func (s *ShippingMethodChannelListing) IsValid() *AppError {
	outer := CreateAppErrorForModel(
		"model.shipping_method_channel_listing.is_valid.%s.app_error",
		"shipping_method_channel_listing_id=",
		"ShippingMethodChannelListing.IsValid",
	)
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

func (s *ShippingMethodChannelListing) commonPre() {
	if s.MinimumOrderPriceAmount == nil {
		s.MinimumOrderPriceAmount = &decimal.Zero
	} else if s.MinimumOrderPrice != nil {
		s.MinimumOrderPriceAmount = &s.MinimumOrderPrice.Amount
	}

	if s.PriceAmount == nil {
		s.PriceAmount = &decimal.Zero
	} else if s.Price != nil {
		s.PriceAmount = &s.Price.Amount
	}

	if s.MaximumOrderPrice != nil {
		s.MaximumOrderPriceAmount = &s.MaximumOrderPrice.Amount
	}
}

// GetTotal retuns current ShippingMethodChannelListing's Price fields
func (s *ShippingMethodChannelListing) GetTotal() *goprices.Money {
	s.PopulateNonDbFields()
	return s.Price
}
