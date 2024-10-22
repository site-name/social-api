package model

import (
	"io"

	"gorm.io/gorm"
)

const (
	STATUS_OUT_OF_OFFICE   = "ooo"
	STATUS_OFFLINE         = "offline"
	STATUS_ONLINE          = "online"
	STATUS_CACHE_SIZE      = SESSION_CACHE_SIZE
	STATUS_CHANNEL_TIMEOUT = 20000  // 20 seconds
	STATUS_MIN_UPDATE_TIME = 120000 // 2 minutes
)

type Status struct {
	UserId         string `json:"user_id" gorm:"primaryKey;type:uuid;index:statuses_userid_key;column:UserId"`
	Status         string `json:"status" gorm:"type:varchar(10);column:Status"`
	Manual         bool   `json:"manual" gorm:"column:Manual"`
	LastActivityAt int64  `json:"last_activity_at" gorm:"type:bigint;column:LastActivityAt"`
}

func (*Status) TableName() string           { return StatusTableName }
func (*Status) BeforeCreate(*gorm.DB) error { return nil }
func (*Status) BeforeUpdate(*gorm.DB) error { return nil }

func (o *Status) ToClusterJson() string {
	oCopy := *o
	return ModelToJson(&oCopy)
}

func StatusFromJson(data io.Reader) *Status {
	var o *Status
	ModelFromJson(&o, data)
	return o
}
