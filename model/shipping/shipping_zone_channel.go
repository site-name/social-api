package shipping

import "github.com/sitename/sitename/model"

// ShippingZoneChannel represents relationships between shipping zones and channels
type ShippingZoneChannel struct {
	Id             string `json:"id"`
	ShippingZoneID string `json:"shipping_zone_id"` // unique together with channelID
	ChannelID      string `json:"channel_id"`       // unique together with shipping zone id
}

func (s *ShippingZoneChannel) IsValid() *model.AppError {
	outer := model.CreateAppErrorForModel(
		"model.shipping_zone_channel.is_valid.%s.app_error",
		"shipping_zone_channel_id=",
		"ShippingZoneChannel.IsValid",
	)

	if !model.IsValidId(s.Id) {
		return outer("id", nil)
	}
	if !model.IsValidId(s.ShippingZoneID) {
		return outer("shipping_zone_id", &s.Id)
	}
	if !model.IsValidId(s.ChannelID) {
		return outer("channel_id", &s.Id)
	}
	return nil
}

func (s *ShippingZoneChannel) PreSave() {
	if s.Id == "" {
		s.Id = model.NewId()
	}
}
