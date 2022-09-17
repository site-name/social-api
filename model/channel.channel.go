package model

import (
	"strings"
	"unicode/utf8"

	"github.com/Masterminds/squirrel"
	"github.com/gosimple/slug"
	"golang.org/x/text/currency"
)

// max lengths for some channel's fields
const (
	CHANNEL_NAME_MAX_LENGTH = 250
	CHANNEL_SLUG_MAX_LENGTH = 255
)

type Channel struct {
	Id             string `json:"id"`
	ShopID         string `json:"shop_id"`
	Name           string `json:"name"`
	IsActive       bool   `json:"is_active"`
	Slug           string `json:"slug"`            // unique
	Currency       string `json:"currency"`        //
	DefaultCountry string `json:"default_country"` // default "US"
}

// ChannelFilterOption is used for building sql queries
type ChannelFilterOption struct {
	Id       squirrel.Sqlizer
	ShopID   squirrel.Sqlizer
	Name     squirrel.Sqlizer
	IsActive *bool
	Slug     squirrel.Sqlizer
	Currency squirrel.Sqlizer
}

type Channels []*Channel

func (c Channels) IDs() []string {
	res := []string{}
	for _, item := range c {
		if item != nil {
			res = append(res, item.Id)
		}
	}

	return res
}

func (c *Channel) String() string {
	return c.Name
}

func (c *Channel) IsValid() *AppError {
	outer := CreateAppErrorForModel(
		"channel.is_valid.%s.app_error",
		"channel_id=",
		"Channel.IsValid",
	)
	if !IsValidId(c.Id) {
		return outer("id", nil)
	}
	if !IsValidId(c.ShopID) {
		return outer("shop_id", &c.Id)
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
	if _, exist := Countries[c.DefaultCountry]; !exist {
		c.DefaultCountry = DEFAULT_COUNTRY
	}
	c.DefaultCountry = strings.ToUpper(c.DefaultCountry)
}

func (c *Channel) PreUpdate() {
	c.commonPre()
}

func (c *Channel) DeepCopy() *Channel {
	res := *c
	return &res
}
