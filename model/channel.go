package model

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"unicode/utf8"

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
	return c.Slug
}

func (c *Channel) invalidChannelErr(field string) *AppError {
	id := fmt.Sprintf("model.channel.is_valid.%s.app_error", field)
	var details string
	if strings.ToLower(field) != "id" {
		details = "channel_id=" + c.Id
	}
	return NewAppError("Channel.IsValid", id, nil, details, http.StatusBadRequest)
}

func (c *Channel) IsValid() *AppError {
	if c.Id == "" {
		return c.invalidChannelErr("id")
	}
	if utf8.RuneCountInString(c.Name) > CHANNEL_NAME_MAX_LENGTH {
		c.invalidChannelErr("name")
	}
	if utf8.RuneCountInString(c.Slug) > CHANNEL_SLUG_MAX_LENGTH {
		c.invalidChannelErr("slug")
	}
	if len(c.Currency) > MAX_LENGTH_CURRENCY_CODE || c.Currency == "" {
		c.invalidChannelErr("currency")
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
