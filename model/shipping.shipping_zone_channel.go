package model

import "github.com/Masterminds/squirrel"

// ShippingZoneChannel represents relationships between shipping zones and channels
type ShippingZoneChannel struct {
	Id             string `json:"id"`
	ShippingZoneID string `json:"shipping_zone_id"` // unique together with channelID
	ChannelID      string `json:"channel_id"`       // unique together with shipping zone id
}

type ShippingZoneChannelFilterOptions struct {
	Conditions squirrel.Sqlizer
}

func (s *ShippingZoneChannel) IsValid() *AppError {
	outer := CreateAppErrorForModel(
		"model.shipping_zone_channel.is_valid.%s.app_error",
		"shipping_zone_channel_id=",
		"ShippingZoneChannel.IsValid",
	)

	if !IsValidId(s.Id) {
		return outer("id", nil)
	}
	if !IsValidId(s.ShippingZoneID) {
		return outer("shipping_zone_id", &s.Id)
	}
	if !IsValidId(s.ChannelID) {
		return outer("channel_id", &s.Id)
	}
	return nil
}

func (s *ShippingZoneChannel) PreSave() {
	if s.Id == "" {
		s.Id = NewId()
	}
}
