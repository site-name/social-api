package model

import "encoding/json"

type SecurityBulletin struct {
	Id               string `json:"id"`
	AppliesToVersion string `json:"applies_to_version"`
}

type SecurityBulletins []SecurityBulletin

func (sb *SecurityBulletin) ToJSON() string {
	return ModelToJson(sb)
}

func (sb *SecurityBulletins) ToJSON() string {
	b, err := json.Marshal(sb)
	if err != nil {
		return "[]"
	}
	return string(b)
}
