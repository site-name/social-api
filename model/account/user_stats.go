package account

import (
	"io"

	"github.com/sitename/sitename/model"
)

type UsersStats struct {
	TotalUsersCount int64 `json:"total_users_count"`
}

func (o *UsersStats) ToJSON() string {
	return model.ModelToJson(o)
}

func UsersStatsFromJson(data io.Reader) *UsersStats {
	var o *UsersStats
	model.ModelFromJson(&o, data)
	return o
}
