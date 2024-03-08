package model

import (
	"net/http"

	"github.com/mattermost/squirrel"
	"gorm.io/gorm"
)

type CollectionChannelListing struct {
	Id           string `json:"id" gorm:"column:Id;type:uuid;primaryKey;default:gen_random_uuid()"`
	CreateAt     int64  `json:"create_at" gorm:"column:CreateAt;type:bigint;autoCreateTime:milli"`
	CollectionID string `json:"collection_id" gorm:"column:CollectionID;type:uuid"`
	ChannelID    string `json:"channel_id" gorm:"column:ChannelID;type:uuid"`
	Publishable
}

// column names for collection-channel-listing table
const (
	CollectionChannelListingColumnId           = "Id"
	CollectionChannelListingColumnCreateAt     = "CreateAt"
	CollectionChannelListingColumnCollectionID = "CollectionID"
	CollectionChannelListingColumnChannelID    = "ChannelID"
)

func (c *CollectionChannelListing) BeforeCreate(_ *gorm.DB) error { return c.IsValid() }
func (c *CollectionChannelListing) BeforeUpdate(_ *gorm.DB) error { return c.IsValid() }
func (c *CollectionChannelListing) TableName() string             { return CollectionChannelListingTableName }

type CollectionChannelListingFilterOptions struct {
	Conditions squirrel.Sqlizer
}

func (c *CollectionChannelListing) IsValid() *AppError {
	if !IsValidId(c.CollectionID) {
		return NewAppError("CollectionChannelListing.IsValid", "model.collection_channel_listing.is_valid.collection_id.app_error", nil, "please provide valid collection id", http.StatusBadRequest)
	}
	if !IsValidId(c.ChannelID) {
		return NewAppError("CollectionChannelListing.IsValid", "model.collection_channel_listing.is_valid.channel_id.app_error", nil, "please provide valid channel id", http.StatusBadRequest)
	}

	return nil
}
