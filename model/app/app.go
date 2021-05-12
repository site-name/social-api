package app

type App struct {
	Id         string  `json:"id"`
	Name       string  `json:"name"`
	CreateAt   int64   `json:"create_at"`
	IsActive   bool    `json:"is_active"`
	Type       string  `json:"type"`
	Identifier *string `json:"identifier"`
	// Permissions []*model.Permission `json:"permissions"`
}
