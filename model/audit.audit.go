package model

import (
	"io"
)

type Audit struct {
	Id        string `json:"id"`
	CreateAt  int64  `json:"create_at"`
	UserId    string `json:"user_id"`
	Action    string `json:"action"`
	ExtraInfo string `json:"extra_info"`
	IpAddress string `json:"ip_address"`
	SessionId string `json:"session_id"`
}

func (a *Audit) IsValid() *AppError {
	outer := CreateAppErrorForModel(
		"audit.is_valid.%s.app_error",
		"audit_id=",
		"Audit.IsValid",
	)
	if !IsValidId(a.Id) {
		return outer("id", nil)
	}
	if !IsValidId(a.UserId) {
		return outer("user_id", &a.Id)
	}
	if !IsValidId(a.SessionId) {
		return outer("session_id", &a.Id)
	}
	if a.CreateAt == 0 {
		return outer("create_at", &a.Id)
	}

	return nil
}

func (a *Audit) PreSave() {
	if a.Id == "" {
		a.Id = NewId()
	}
	a.CreateAt = GetMillis()
}

func (a *Audit) ToJSON() string {
	return ModelToJson(a)
}

func AuditFromJson(data io.Reader) *Audit {
	var a *Audit
	ModelFromJson(&a, data)
	return a
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
