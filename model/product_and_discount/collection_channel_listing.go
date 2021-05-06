package product_and_discount

import (
	"io"

	"github.com/sitename/sitename/model"
)

type CollectionChannelListing struct {
	Id           string `json:"id"`
	CollectionID string `json:"collection_id"`
	ChannelID    string `json:"channel_id"`
}

func (c *CollectionChannelListing) IsValid() *model.AppError {
	outer := model.CreateAppErrorForModel(
		"model.collection_channel_listing.is_valid.%s.app_error",
		"collection_channel_listing_id=",
		"CollectionChannelListing.IsValid")

	if c.Id == "" {
		return outer("id", nil)
	}
	if c.CollectionID == "" {
		return outer("collection_id", &c.Id)
	}
	if c.ChannelID == "" {
		return outer("channel_id", &c.Id)
	}

	return nil
}

func (c *CollectionChannelListing) ToJson() string {
	return model.ModelToJson(c)
}

func CollectionChannelListingFromJson(data io.Reader) *CollectionChannelListing {
	var c CollectionChannelListing
	model.ModelFromJson(&c, data)
	return &c
}
