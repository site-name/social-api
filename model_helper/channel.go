package model_helper

import (
	"net/http"

	"github.com/gosimple/slug"
	"github.com/sitename/sitename/model"
)

func ChannelPreSave(channel *model.Channel) {
	if channel.ID == "" {
		channel.ID = NewId()
	}
	ChannelCommonPre(channel)
}

func ChannelCommonPre(channel *model.Channel) {
	channel.Name = SanitizeUnicode(channel.Name)
	channel.Slug = slug.Make(channel.Name)
	if channel.Currency.IsValid() != nil {
		channel.Currency = DEFAULT_CURRENCY
	}
	if channel.DefaultCountry.IsValid() != nil {
		channel.DefaultCountry = DEFAULT_COUNTRY
	}
	channel.Annotations = nil
}

func ChannelIsValid(channel model.Channel) *AppError {
	if !IsValidId(channel.ID) {
		return NewAppError("ChannelIsValid", "model.channel.is_valid.id.app_error", nil, "", http.StatusBadRequest)
	}
	if channel.Name == "" {
		return NewAppError("ChannelIsValid", "model.channel.is_valid.name.app_error", nil, "", http.StatusBadRequest)
	}
	if !slug.IsSlug(channel.Slug) {
		return NewAppError("ChannelIsValid", "model.channel.is_valid.slug.app_error", nil, "", http.StatusBadRequest)
	}
	if channel.Currency.IsValid() != nil {
		return NewAppError("ChannelIsValid", "model.channel.is_valid.currency.app_error", nil, "", http.StatusBadRequest)
	}
	if channel.DefaultCountry.IsValid() != nil {
		return NewAppError("ChannelIsValid", "model.channel.is_valid.default_country.app_error", nil, "", http.StatusBadRequest)
	}
	return nil
}

var ChannelAnnotationKeys = struct {
	HasOrders string
}{
	HasOrders: "has_orders",
}
