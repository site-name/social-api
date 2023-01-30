package model

import "github.com/Masterminds/squirrel"

type CollectionChannelListing struct {
	Id           string `json:"id"`
	CreateAt     int64  `json:"create_at"`
	CollectionID string `json:"collection_id"`
	ChannelID    string `json:"channel_id"`
	Publishable
}

type CollectionChannelListingFilterOptions struct {
	Id           squirrel.Sqlizer
	CollectionID squirrel.Sqlizer
	ChannelID    squirrel.Sqlizer
}

func (c *CollectionChannelListing) IsValid() *AppError {
	outer := CreateAppErrorForModel(
		"collection_channel_listing.is_valid.%s.app_error",
		"collection_channel_listing_id=",
		"CollectionChannelListing.IsValid")

	if !IsValidId(c.Id) {
		return outer("id", nil)
	}
	if c.CreateAt == 0 {
		return outer("create_at", &c.Id)
	}
	if !IsValidId(c.CollectionID) {
		return outer("collection_id", &c.Id)
	}
	if !IsValidId(c.ChannelID) {
		return outer("channel_id", &c.Id)
	}

	return nil
}

func (c *CollectionChannelListing) PreSave() {
	if c.Id == "" {
		c.Id = NewId()
	}
	c.CreateAt = GetMillis()

}
