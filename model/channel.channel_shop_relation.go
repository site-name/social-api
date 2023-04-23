package model

import (
	"fmt"
	"net/http"

	"github.com/Masterminds/squirrel"
)

type ChannelShopRelation struct {
	Id        string `json:"id"`
	ChannelID string `json:"channel_id"` // unique together with shop_id
	CreateAt  int64  `json:"create_at"`
	EndAt     int64  `json:"end_at"`
}

type ChannelShopRelationFilterOptions struct {
	Id        squirrel.Sqlizer
	ChannelID squirrel.Sqlizer
}

func (c *ChannelShopRelation) IsValid() *AppError {
	if !IsValidId(c.Id) {
		return NewAppError("ChannelShopRelation.IsValid", "model.channel_shop.id.app_error", nil, fmt.Sprintf("%s is invalid id", c.Id), http.StatusBadRequest)
	}

	if !IsValidId(c.ChannelID) {
		return NewAppError("ChannelShopRelation.IsValid", "model.channel_shop.channel_id.app_error", nil, fmt.Sprintf("%s is invalid channelID", c.ChannelID), http.StatusBadRequest)
	}
	if c.CreateAt <= 0 {
		return NewAppError("ChannelShopRelation.IsValid", "model.channel_shop.create_at.app_error", nil, "createAt must not be less than or equal to zero", http.StatusBadRequest)
	}
	if c.EndAt < 0 {
		return NewAppError("ChannelShopRelation.IsValid", "model.channel_shop.end_at.app_error", nil, "endAt must not be negative", http.StatusBadRequest)
	}

	return nil
}

func (c *ChannelShopRelation) PreSave() {
	if !IsValidId(c.Id) {
		c.Id = NewId()
	}
	c.CreateAt = GetMillis()
}
