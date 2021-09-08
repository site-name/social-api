package channel

import (
	"strings"
	"unicode/utf8"

	"github.com/gosimple/slug"
	"github.com/sitename/sitename/model"
	"golang.org/x/text/currency"
)

// max lengths for some channel's fields
const (
	CHANNEL_NAME_MAX_LENGTH = 250
	CHANNEL_SLUG_MAX_LENGTH = 255
)

type Channel struct {
	Id             string `json:"id"`
	Name           string `json:"name"`
	IsActive       bool   `json:"is_active"`
	Slug           string `json:"slug"` // unique
	Currency       string `json:"currency"`
	DefaultCountry string `json:"default_country"`
}

// ChannelFilterOption is used for building sql queries
type ChannelFilterOption struct {
	Id       *model.StringFilter
	Name     *model.StringFilter
	IsActive *bool
	Slug     *model.StringFilter
	Currency *model.StringFilter
}

func (c *Channel) String() string {
	return c.Name
}

func (c *Channel) IsValid() *model.AppError {
	outer := model.CreateAppErrorForModel(
		"model.channel.is_valid.%s.app_error",
		"channel_id=",
		"Channel.IsValid",
	)
	if !model.IsValidId(c.Id) {
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
	if c.DefaultCountry != "" && model.Countries[c.DefaultCountry] == "" {
		return outer("default_country", &c.Id)
	}

	return nil
}

func (c *Channel) ToJson() string {
	return model.ModelToJson(c)
}

func (c *Channel) PreSave() {
	if c.Id == "" {
		c.Id = model.NewId()
	}
	c.Name = model.SanitizeUnicode(c.Name)
	c.Slug = slug.Make(c.Name)
	c.Currency = strings.ToUpper(c.Currency)
}

func (c *Channel) PreUpdate() {
	c.Name = model.SanitizeUnicode(c.Name)
	c.Currency = strings.ToUpper(c.Currency)
	// c.Slug = slug.Make(c.Name)
}
