package model

import (
	"net/http"

	"gorm.io/gorm"
)

type Audit struct {
	Id        string `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid();column:Id"`
	CreateAt  int64  `json:"create_at" gorm:"type:bigint;autoCreateTime:milli;column:CreateAt"`
	UserId    string `json:"user_id" gorm:"type:uuid;column:UserId"`
	Action    string `json:"action" gorm:"type:varchar(512);column:Action"`         // varchar(512)
	ExtraInfo string `json:"extra_info" gorm:"type:varchar(1024);column:ExtraInfo"` // varchar(1024)
	IpAddress string `json:"ip_address" gorm:"type:varchar(64);column:IpAddress"`   // varchar(64)
	SessionId string `json:"session_id" gorm:"type:uuid;column:SessionId"`
}

func (a *Audit) TableName() string             { return AuditTableName }
func (a *Audit) BeforeCreate(_ *gorm.DB) error { return a.IsValid() }
func (a *Audit) BeforeUpdate(_ *gorm.DB) error { a.CreateAt = 0; return a.IsValid() }

func (a *Audit) IsValid() *AppError {
	if !IsValidId(a.UserId) {
		return NewAppError("Audit.IsValid", "model.audit.is_valid.user_id.app_error", nil, "please provide valid user id", http.StatusBadRequest)
	}
	if !IsValidId(a.SessionId) {
		return NewAppError("Audit.IsValid", "model.audit.is_valid.session_id.app_error", nil, "please provide valid session id", http.StatusBadRequest)
	}
	return nil
}

type Audits []Audit

func (o Audits) Etag() string {
	if len(o) > 0 {
		// the first in the list is always the most current
		return Etag(o[0].CreateAt)
	}
	return ""
}

func (o Audits) ToJSON() string {
	return ModelToJson(&o)
}
