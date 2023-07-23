package model

import (
	"io"
)

type Audit struct {
	Id        string `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid();column:Id"`
	CreateAt  int64  `json:"create_at" gorm:"type:bigint;default:autoCreateTime:milli;column:CreateAt"`
	UserId    string `json:"user_id" gorm:"type:uuid;column:UserId"`
	Action    string `json:"action" gorm:"type:varchar(512);column:Action"`         // varchar(512)
	ExtraInfo string `json:"extra_info" gorm:"type:varchar(1024);column:ExtraInfo"` // varchar(1024)
	IpAddress string `json:"ip_address" gorm:"type:varchar(64);column:IpAddress"`   // varchar(64)
	SessionId string `json:"session_id" gorm:"type:uuid;column:SessionId"`
}

func (*Audit) TableName() string { return AuditTableName }

func (a *Audit) IsValid() *AppError {
	outer := CreateAppErrorForModel(
		"model.audit.is_valid.%s.app_error",
		"audit_id=",
		"Audit.IsValid",
	)
	if !IsValidId(a.UserId) {
		return outer("user_id", &a.Id)
	}
	if !IsValidId(a.SessionId) {
		return outer("session_id", &a.Id)
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

func AuditsFromJson(data io.Reader) Audits {
	var o Audits
	ModelFromJson(&o, data)
	return o
}
