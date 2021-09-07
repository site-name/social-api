package gqlmodel

import "time"

// type ExportEvent struct {
// 	ID      string           `json:"id"`
// 	Date    time.Time        `json:"date"`
// 	Type    ExportEventsEnum `json:"type"`
// 	User    *User            `json:"user"`
// 	App     *App             `json:"app"`
// 	Message string           `json:"message"`
// }

// func (ExportEvent) IsNode() {}

type ExportEvent struct {
	ID      string           `json:"id"`
	Date    time.Time        `json:"date"`
	Type    ExportEventsEnum `json:"type"`
	UserID  *string          `json:"user"`
	Message string           `json:"message"`
}

func (ExportEvent) IsNode() {}
