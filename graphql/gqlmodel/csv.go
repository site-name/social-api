package gqlmodel

import (
	"strings"
	"time"

	"github.com/sitename/sitename/model/csv"
	"github.com/sitename/sitename/modules/util"
)

// type ExportEvent struct {
// 	ID      string           `json:"id"`
// 	Date    time.Time        `json:"date"`
// 	Type    ExportEventsEnum `json:"type"`
// 	User    *User            `json:"user"`
// 	Message string           `json:"message"`
// }

func (ExportEvent) IsNode() {}

// type ExportFile struct {
// 	ID        string        `json:"id"`
// 	User      *User         `json:"user"`
// 	Status    JobStatusEnum `json:"status"`
// 	CreatedAt time.Time     `json:"createdAt"`
// 	UpdatedAt time.Time     `json:"updatedAt"`
// 	Message   *string       `json:"message"`
// 	URL       *string       `json:"url"`
// 	Events    []*ExportEvent `json:"events"`
// }

func (ExportFile) IsNode() {}
func (ExportFile) IsJob()  {}

type ExportEvent struct {
	ID      string           `json:"id"`
	Date    time.Time        `json:"date"`
	Type    ExportEventsEnum `json:"type"`
	UserID  *string          `json:"user"`
	Message string           `json:"message"`
}

// SystemExportEventToGraphqlExportEvent converts given system export event to graphql export event
func SystemExportEventToGraphqlExportEvent(event *csv.ExportEvent) *ExportEvent {
	res := &ExportEvent{
		ID:     event.Id,
		Date:   util.TimeFromMillis(event.Date),
		Type:   ExportEventsEnum(strings.ToUpper(event.Type)),
		UserID: event.UserID,
	}

	if event.Parameters != nil && (*event.Parameters)["message"] != "" {
		res.Message = (*event.Parameters)["message"]
	}

	return res
}

// SystemExportEventsToGraphqlExportEvents converts given system export events to graphql export events
func SystemExportEventsToGraphqlExportEvents(events []*csv.ExportEvent) []*ExportEvent {
	res := []*ExportEvent{}
	for _, event := range events {
		res = append(res, SystemExportEventToGraphqlExportEvent(event))
	}

	return res
}

type ExportFile struct {
	ID        string                `json:"id"`
	UserID    *string               `json:"user"`
	Status    JobStatusEnum         `json:"status"`
	CreatedAt time.Time             `json:"createdAt"`
	UpdatedAt time.Time             `json:"updatedAt"`
	Message   *string               `json:"message"`
	URL       func() *string        `json:"url"`
	Events    func() []*ExportEvent `json:"events"`
}

// SystemExportFileToGraphqlExportFile converts given system export file to graphql export file
func SystemExportFileToGraphqlExportFile(file *csv.ExportFile) *ExportFile {
	res := &ExportFile{
		ID:        file.Id,
		UserID:    file.UserID,
		CreatedAt: util.TimeFromMillis(file.CreateAt),
		UpdatedAt: util.TimeFromMillis(file.UpdateAt),
	}

	return res
}

// SystemExportFilesToGraphqlExportFiles converts given system export files to graphql export files
func SystemExportFilesToGraphqlExportFiles(files []*csv.ExportFile) []*ExportFile {
	res := []*ExportFile{}
	for _, file := range files {
		res = append(res, SystemExportFileToGraphqlExportFile(file))
	}

	return res
}
