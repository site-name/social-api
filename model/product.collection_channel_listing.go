package model

import (
	"net/http"

	"github.com/Masterminds/squirrel"
	"gorm.io/gorm"
)

type CollectionChannelListing struct {
	Id           string `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid();column:Id"`
	CreateAt     int64  `json:"create_at" gorm:"type:bigint;column:CreateAt;autoCreateTime:milli"`
	CollectionID string `json:"collection_id" gorm:"type:uuid;column:CollectionID"`
	ChannelID    string `json:"channel_id" gorm:"type:uuid;column:ChannelID"`
	Publishable
}

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
