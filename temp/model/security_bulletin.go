package model

type SecurityBulletin struct {
	Id               string `json:"id"`
	AppliesToVersion string `json:"applies_to_version"`
}

type SecurityBulletins []SecurityBulletin
