package model

import (
	"github.com/Masterminds/squirrel"
)

type ExportEventType string

// choices for export event's type
const (
	EXPORT_PENDING          ExportEventType = "export_pending"
	EXPORT_SUCCESS          ExportEventType = "export_success"
	EXPORT_FAILED           ExportEventType = "export_failed"
	EXPORT_DELETED          ExportEventType = "export_deleted"
	EXPORTED_FILE_SENT      ExportEventType = "exported_file_sent"
	EXPORT_FAILED_INFO_SENT ExportEventType = "export_failed_info_sent"
)

var ExportTypeString = map[ExportEventType]string{
	EXPORT_PENDING:          "Data export was started.",
	EXPORT_SUCCESS:          "Data export was completed successfully.",
	EXPORT_FAILED:           "Data export failed.",
	EXPORT_DELETED:          "Export file was deleted.",
	EXPORTED_FILE_SENT:      "Email with link to download file was sent to the customer.",
	EXPORT_FAILED_INFO_SENT: "Email with info that export failed was sent to the customer.",
}

func (t ExportEventType) IsValid() bool {
	return ExportTypeString[t] != ""
}

// Model used to store events that happened during the export file lifecycle.
type ExportEvent struct {
	Id           string          `json:"id"`
	Date         int64           `json:"date"`
	Type         ExportEventType `json:"type"`
	Parameters   *StringMap      `json:"parameters"`
	ExportFileID string          `json:"export_file_id"`
	UserID       *string         `json:"user_id"`
}

// ExportEventFilterOption is used to build squirrel queries
type ExportEventFilterOption struct {
	Id           squirrel.Sqlizer
	ExportFileID squirrel.Sqlizer
	UserID       squirrel.Sqlizer
}

func (e *ExportEvent) IsValid() *AppError {
	outer := CreateAppErrorForModel(
		"export_event.is_valid.%s.app_error",
		"export_event_id",
		"ExportEvent.IsValid",
	)
	if !IsValidId(e.Id) {
		return outer("id", nil)
	}
	if !IsValidId(e.ExportFileID) {
		return outer("export_file_id", &e.Id)
	}
	if e.UserID != nil && !IsValidId(*e.UserID) {
		return outer("user_id", &e.Id)
	}
	if !e.Type.IsValid() {
		return outer("type", &e.Id)
	}
	if e.Date == 0 {
		return outer("date", &e.Id)
	}

	return nil
}

func (e *ExportEvent) ToJSON() string {
	return ModelToJson(e)
}

func (e *ExportEvent) PreSave() {
	if e.Id == "" {
		e.Id = NewId()
	}
	e.Date = GetMillis()
}
