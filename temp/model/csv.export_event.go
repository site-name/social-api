package model

import (
	"net/http"

	"github.com/Masterminds/squirrel"
	"gorm.io/gorm"
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
	Id           string          `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid();column:Id"`
	Date         int64           `json:"date" gorm:"type:bigint;column:Date;autoCreateTime:milli"` // NOTE editable
	Type         ExportEventType `json:"type" gorm:"type:varchar(255);column:Type"`
	Parameters   *StringMap      `json:"parameters" gorm:"type:jsonb;column:Parameters"`
	ExportFileID string          `json:"export_file_id" gorm:"type:uuid;column:ExportFileID"` // delete CASCADE
	UserID       *string         `json:"user_id" gorm:"type:uuid;column:UserID"`              // delete CASCADE
}

func (c *ExportEvent) BeforeCreate(_ *gorm.DB) error { return c.IsValid() }
func (c *ExportEvent) BeforeUpdate(_ *gorm.DB) error { return c.IsValid() }
func (c *ExportEvent) TableName() string             { return CsvExportEventTableName }

// ExportEventFilterOption is used to build squirrel queries
type ExportEventFilterOption struct {
	Conditions squirrel.Sqlizer
}

func (e *ExportEvent) IsValid() *AppError {
	if !e.Type.IsValid() {
		return NewAppError("ExportEvent.IsValid", "model.export_event.is_valid.type.app_error", nil, "export event type is invalid", http.StatusBadRequest)
	}
	return nil
}
