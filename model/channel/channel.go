package channel

import (
	"io"
	"strings"
	"unicode/utf8"

	"github.com/gosimple/slug"
	"github.com/sitename/sitename/model"
	"golang.org/x/text/currency"
)

const (
	CHANNEL_NAME_MAX_LENGTH = 250
	CHANNEL_SLUG_MAX_LENGTH = 255
)

type Channel struct {
	Id       string `json:"id"`
	Name     string `json:"name"`
	IsActive bool   `json:"is_active"`
	Slug     string `json:"slug"`
	Currency string `json:"currency"`
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

	return nil
}

func (c *Channel) ToJson() string {
	return model.ModelToJson(c)
}

func ChannelFromJson(data io.Reader) *Channel {
	var channel Channel
	model.ModelFromJson(&channel, data)
	return &channel
}

func (c *Channel) PreSave() {
	if c.Id == "" {
		c.Id = model.NewId()
	}
	c.Name = model.SanitizeUnicode(c.Name)
	c.Slug = slug.Make(c.Name)
}

func (c *Channel) PreUpdate() {
	c.Name = model.SanitizeUnicode(c.Name)
	c.Slug = slug.Make(c.Name)
}