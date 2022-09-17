package model

import (
	"io"
)

type UsersStats struct {
	TotalUsersCount int64 `json:"total_users_count"`
}

func (o *UsersStats) ToJSON() string {
	return ModelToJson(o)
}

func UsersStatsFromJson(data io.Reader) *UsersStats {
	var o *UsersStats
	ModelFromJson(&o, data)
	return o
}
