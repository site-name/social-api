package model

import (
	"io"
	"net/http"

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
	UserId         UUID   `json:"user_id" gorm:"primaryKey;type:uuid;index:statuses_userid_key;column:UserId"`
	Status         string `json:"status" gorm:"type:varchar(10);column:Status"`
	Manual         bool   `json:"manual" gorm:"column:Manual"`
	LastActivityAt int64  `json:"last_activity_at" gorm:"type:bigint;column:LastActivityAt"`
}

func (*Status) TableName() string             { return StatusTableName }
func (s *Status) BeforeCreate(*gorm.DB) error { return s.IsValid() }
func (s *Status) BeforeUpdate(*gorm.DB) error { return s.IsValid() }

func (s *Status) IsValid() *AppError {
	if !IsValidId(s.UserId) {
		return NewAppError("Status.IsValid", "model.account.status.is_valid.user_id.app_error", nil, "pleaseprovide valid user id", http.StatusBadRequest)
	}
	return nil
}

func (o *Status) ToClusterJson() string {
	oCopy := *o
	return ModelToJson(&oCopy)
}

func StatusFromJson(data io.Reader) *Status {
	var o *Status
	ModelFromJson(&o, data)
	return o
}
