package audit

import (
	"io"

	"github.com/sitename/sitename/model"
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

func (a *Audit) IsValid() *model.AppError {
	outer := model.CreateAppErrorForModel(
		"model.audit.is_valid.%s.app_error",
		"audit_id=",
		"Audit.IsValid",
	)
	if !model.IsValidId(a.Id) {
		return outer("id", nil)
	}
	if !model.IsValidId(a.UserId) {
		return outer("user_id", &a.Id)
	}
	if !model.IsValidId(a.SessionId) {
		return outer("session_id", &a.Id)
	}

	return nil
}

func (a *Audit) ToJson() string {
	return model.ModelToJson(a)
}

func AuditFromJson(data io.Reader) *Audit {
	var a *Audit
	model.ModelFromJson(&a, data)
	return a
}

type Audits []Audit

func (o Audits) Etag() string {
	if len(o) > 0 {
		// the first in the list is always the most current
		return model.Etag(o[0].CreateAt)
	}
	return ""
}

func (o Audits) ToJson() string {
	return model.ModelToJson(&o)
}

func AuditsFromJson(data io.Reader) Audits {
	var o Audits
	model.ModelFromJson(&o, data)
	return o
}
