package model

import (
	"io"

	"github.com/sitename/sitename/modules/json"
)

type SecurityBulletin struct {
	Id               string `json:"id"`
	AppliesToVersion string `json:"applies_to_version"`
}

type SecurityBulletins []SecurityBulletin

func (sb *SecurityBulletin) ToJson() string {
	return ModelToJson(sb)
}

func SecurityBulletinFromJson(data io.Reader) *SecurityBulletin {
	var o *SecurityBulletin
	ModelFromJson(&o, data)
	return o
}

func (sb *SecurityBulletins) ToJson() string {
	b, err := json.JSON.Marshal(sb)
	if err != nil {
		return "[]"
	}
	return string(b)
}

func SecurityBulletinsFromJson(data io.Reader) SecurityBulletins {
	var o SecurityBulletins
	json.JSON.NewDecoder(data).Decode(&o)
	return o
}
