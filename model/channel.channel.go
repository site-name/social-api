package model

import (
	"net/http"
	"strings"

	"github.com/Masterminds/squirrel"
	"github.com/gosimple/slug"
	"github.com/samber/lo"
	"golang.org/x/text/currency"
	"gorm.io/gorm"
)

type Channel struct {
	Id             UUID        `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid();column:Id"`
	Name           string      `json:"name" gorm:"type:varchar(250);column:Name"`
	IsActive       bool        `json:"is_active" gorm:"column:IsActive"`
	Slug           string      `json:"slug" gorm:"type:varchar(255);column:Slug;uniqueIndex:slug_unique_key"` // unique
	Currency       string      `json:"currency" gorm:"column:Currency;type:varchar(3)"`
	DefaultCountry CountryCode `json:"default_country" gorm:"column:DefaultCountry;type:varchar(10)"` // default "US"

	ShippingZones ShippingZones `json:"-" gorm:"many2many:ShippingZoneChannels"`
	Orders        Orders        `json:"-" gorm:"foreignKey:ChannelID"`

	hasOrders bool `gorm:"-"`
}

func (c *Channel) GetHasOrders() bool            { return c.hasOrders }
func (c *Channel) SetHasOrders(b bool)           { c.hasOrders = b }
func (c *Channel) BeforeCreate(_ *gorm.DB) error { c.PreSave(); return c.IsValid() }
func (c *Channel) BeforeUpdate(_ *gorm.DB) error { c.PreUpdate(); return c.IsValid() }
func (c *Channel) TableName() string             { return ChannelTableName }

// ChannelFilterOption is used for building sql queries
type ChannelFilterOption struct {
	ShippingZoneChannels_ShippingZoneID squirrel.Sqlizer // INNER JOIN ShippingZoneChannels ON ... WHERE ChannelShippingZones.ShippingZoneID ...
	AnnotateHasOrders                   bool             // to check if there are at least 1 order associated to this channel

	VoucherChannelListing_VoucherID squirrel.Sqlizer // INNER JOIN VoucherChannelListings ON ... WHERE VoucherChannelListings.VoucherID ...
	Conditions                      squirrel.Sqlizer
	Limit                           int
}

type Channels []*Channel

func (c Channels) IDs() []UUID {
	return lo.Map(c, func(ch *Channel, _ int) UUID { return ch.Id })
}

func (c Channels) Currencies() []string {
	return lo.Map(c, func(ch *Channel, _ int) string { return ch.Currency })
}

func (c Channels) Len() int { return len(c) }

func (c *Channel) String() string {
	return c.Name
}

func (c *Channel) IsValid() *AppError {
	if un, err := currency.ParseISO(c.Currency); err != nil || !strings.EqualFold(un.String(), c.Currency) {
		return NewAppError("Channel.IsValid", "model.channel.is_valid.currency.app_error", nil, "please provide valid currency", http.StatusBadRequest)
	}
	if !c.DefaultCountry.IsValid() {
		return NewAppError("Channel.IsValid", "model.channel.is_valid.default_country.app_error", nil, "please provide valid default country", http.StatusBadRequest)
	}

	return nil
}

func (c *Channel) PreSave() {
	c.commonPre()
	c.Slug = slug.Make(c.Name)
}

func (c *Channel) commonPre() {
	c.Name = SanitizeUnicode(c.Name)
	c.Currency = strings.ToUpper(c.Currency)
	if !c.DefaultCountry.IsValid() {
		c.DefaultCountry = DEFAULT_COUNTRY
	}
}

func (c *Channel) PreUpdate() {
	c.commonPre()
}

func (c *Channel) DeepCopy() *Channel {
	res := *c
	return &res
}
