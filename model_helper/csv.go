package model_helper

import (
	"net/http"

	"github.com/sitename/sitename/model"
)

func ExportEventPreSave(event *model.ExportEvent) {
	if event.ID == "" {
		event.ID = NewId()
	}
	if event.Date == 0 {
		event.Date = GetMillis()
	}
}

func ExportEventIsValid(event model.ExportEvent) *AppError {
	if event.Type.IsValid() != nil {
		return NewAppError("ExportEventIsValid", "model.export_event.is_valid.type.app_error", nil, "export event type is invalid", http.StatusBadRequest)
	}
	if event.Date <= 0 {
		return NewAppError("ExportEventIsValid", "model.export_event.is_valid.date.app_error", nil, "export event date is invalid", http.StatusBadRequest)
	}
	if !IsValidId(event.ExportFileID) {
		return NewAppError("ExportEventIsValid", "model.export_event.is_valid.export_file_id.app_error", nil, "export event export file id is invalid", http.StatusBadRequest)
	}
	if !event.UserID.IsNil() && !IsValidId(*event.UserID.String) {
		return NewAppError("ExportEventIsValid", "model.export_event.is_valid.user_id.app_error", nil, "export event user id is invalid", http.StatusBadRequest)
	}
	return nil
}
