package gqlmodel

import (
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model/channel"
)

// ORIGINAL IMPLEMENTEATION

// type Channel struct {
// 	ID             string          `json:"id"`
// 	Name           string          `json:"name"`
// 	IsActive       bool            `json:"isActive"`
// 	Slug           string          `json:"slug"`
// 	CurrencyCode   string          `json:"currencyCode"`
// 	HasOrders      bool            `json:"hasOrders"`
// 	DefaultCountry *CountryDisplay `json:"defaultCountry"`
// }

type Channel struct {
	ID             string          `json:"id"`
	Name           string          `json:"name"`
	IsActive       bool            `json:"isActive"`
	Slug           string          `json:"slug"`
	CurrencyCode   string          `json:"currencyCode"`
	HasOrders      bool            `json:"hasOrders"`
	DefaultCountry *CountryDisplay `json:"defaultCountry"`
}

// SystemChannelToGraphqlChannel converts given system channel to graphql channel
func SystemChannelToGraphqlChannel(c *channel.Channel) *Channel {
	return &Channel{
		ID:           c.Id,
		Name:         c.Name,
		Slug:         c.Slug,
		CurrencyCode: c.Currency,
		IsActive:     c.IsActive,
		DefaultCountry: &CountryDisplay{
			Code:    c.DefaultCountry,
			Country: model.Countries[c.DefaultCountry],
		},
	}
}

func (Channel) IsNode() {}
