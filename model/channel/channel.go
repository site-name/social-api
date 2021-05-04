package channel

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"unicode/utf8"

	"github.com/gosimple/slug"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/modules/json"
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

func (c *Channel) invalidChannelErr(field string) *model.AppError {
	id := fmt.Sprintf("model.channel.is_valid.%s.app_error", field)
	var details string
	if strings.ToLower(field) != "id" {
		details = "channel_id=" + c.Id
	}
	return model.NewAppError("Channel.IsValid", id, nil, details, http.StatusBadRequest)
}

func (c *Channel) IsValid() *model.AppError {
	if c.Id == "" {
		return c.invalidChannelErr("id")
	}
	if utf8.RuneCountInString(c.Name) > CHANNEL_NAME_MAX_LENGTH {
		c.invalidChannelErr("name")
	}
	if utf8.RuneCountInString(c.Slug) > CHANNEL_SLUG_MAX_LENGTH {
		c.invalidChannelErr("slug")
	}
	if un, err := currency.ParseISO(c.Currency); err != nil || !strings.EqualFold(un.String(), c.Currency) {
		return c.invalidChannelErr("currency")
	}

	return nil
}

func (c *Channel) ToJson() string {
	b, _ := json.JSON.Marshal(c)
	return string(b)
}

func ChannelFromJson(data io.Reader) *Channel {
	var channel Channel
	err := json.JSON.NewDecoder(data).Decode(&channel)
	if err != nil {
		return nil
	}
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
	if c.Id == "" {
		c.Id = model.NewId()
	}
	c.Name = model.SanitizeUnicode(c.Name)
	c.Slug = slug.Make(c.Name)
}
