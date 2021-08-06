package product_and_discount

import (
	"github.com/sitename/sitename/model"
)

type CollectionChannelListing struct {
	Id           string `json:"id"`
	CreateAt     int64  `json:"create_at"`
	CollectionID string `json:"collection_id"`
	ChannelID    string `json:"channel_id"`
}

func (c *CollectionChannelListing) IsValid() *model.AppError {
	outer := model.CreateAppErrorForModel(
		"model.collection_channel_listing.is_valid.%s.app_error",
		"collection_channel_listing_id=",
		"CollectionChannelListing.IsValid")

	if !model.IsValidId(c.Id) {
		return outer("id", nil)
	}
	if c.CreateAt == 0 {
		return outer("create_at", &c.Id)
	}
	if !model.IsValidId(c.CollectionID) {
		return outer("collection_id", &c.Id)
	}
	if !model.IsValidId(c.ChannelID) {
		return outer("channel_id", &c.Id)
	}

	return nil
}

func (c *CollectionChannelListing) PreSave() {
	if c.Id == "" {
		c.Id = model.NewId()
	}
	if c.CreateAt == 0 {
		c.CreateAt = model.GetMillis()
	}
}
