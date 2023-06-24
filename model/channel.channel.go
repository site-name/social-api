package model

import (
	"strings"
	"unicode/utf8"

	"github.com/Masterminds/squirrel"
	"github.com/gosimple/slug"
	"github.com/samber/lo"
	"golang.org/x/text/currency"
)

// max lengths for some channel's fields
const (
	CHANNEL_NAME_MAX_LENGTH = 250
	CHANNEL_SLUG_MAX_LENGTH = 255
)

type Channel struct {
	Id             string      `json:"id"`
	Name           string      `json:"name"`
	IsActive       bool        `json:"is_active"`
	Slug           string      `json:"slug"`            // unique
	Currency       string      `json:"currency"`        //
	DefaultCountry CountryCode `json:"default_country"` // default "US"

	hasOrders     bool          `db:"-"`
	shippingZones ShippingZones `db:"-"` // get populated in some queries that require selected related shipping zones
}

func (c *Channel) GetShippingZones() ShippingZones {
	return c.shippingZones
}

func (c *Channel) SetShippingZones(s ShippingZones) {
	c.shippingZones = s
}

func (c *Channel) GetHasOrders() bool {
	return c.hasOrders
}

func (c *Channel) SetHasOrders(b bool) {
	c.hasOrders = b
}

// ChannelFilterOption is used for building sql queries
type ChannelFilterOption struct {
	Id       squirrel.Sqlizer
	Name     squirrel.Sqlizer
	IsActive *bool
	Slug     squirrel.Sqlizer
	Currency squirrel.Sqlizer

	ShippingZoneChannels_ShippingZoneID squirrel.Sqlizer // INNER JOIN ShippingZoneChannels ON ... WHERE ChannelShippingZones.ShippingZoneID ...
	AnnotateHasOrders                   bool             // to check if there are at least 1 order associated to this channel

	Extra squirrel.Sqlizer
}

type Channels []*Channel

func (c Channels) IDs() []string {
	return lo.Map(c, func(ch *Channel, _ int) string { return ch.Id })
}

func (c Channels) Currencies() []string {
	return lo.Map(c, func(ch *Channel, _ int) string { return ch.Currency })
}

func (c *Channel) String() string {
	return c.Name
}

func (c *Channel) IsValid() *AppError {
	outer := CreateAppErrorForModel(
		"model.channel.is_valid.%s.app_error",
		"channel_id=",
		"Channel.IsValid",
	)
	if !IsValidId(c.Id) {
		return outer("id", nil)
	}
	if utf8.RuneCountInString(c.Name) > CHANNEL_NAME_MAX_LENGTH {
		outer("name", &c.Id)
	}
	if utf8.RuneCountInString(c.Slug) > CHANNEL_SLUG_MAX_LENGTH {
		outer("slug", &c.Id)
	}
	if un, err := currency.ParseISO(c.Currency); err != nil || !strings.EqualFold(un.String(), c.Currency) {
		return outer("currency", &c.Id)
	}
	if c.DefaultCountry != "" && Countries[c.DefaultCountry] == "" {
		return outer("default_country", &c.Id)
	}

	return nil
}

func (c *Channel) ToJSON() string {
	return ModelToJson(c)
}

func (c *Channel) PreSave() {
	if c.Id == "" {
		c.Id = NewId()
	}
	c.commonPre()
	c.Slug = slug.Make(c.Name)
}

func (c *Channel) commonPre() {
	c.Name = SanitizeUnicode(c.Name)
	c.Currency = strings.ToUpper(c.Currency)
	if !c.DefaultCountry.IsValid() {
		c.DefaultCountry = DEFAULT_COUNTRY
	}
	// c.DefaultCountry = strings.ToUpper(c.DefaultCountry)
}

func (c *Channel) PreUpdate() {
	c.commonPre()
}

func (c *Channel) DeepCopy() *Channel {
	res := *c
	res.shippingZones = c.shippingZones.DeepCopy()
	return &res
}
