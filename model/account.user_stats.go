package model

type UsersStats struct {
	TotalUsersCount int64 `json:"total_users_count"`
}

func (o *UsersStats) ToJSON() string {
	return ModelToJson(o)
}
