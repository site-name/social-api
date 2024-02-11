package model_helper

import (
	"net/http"

	"github.com/sitename/sitename/model"
)

func UploadSessionPreSave(u *model.UploadSession) {
	u.CreatedAt = GetMillis()
	u.FileName = SanitizeUnicode(u.FileName)
}

func UploadSessionIsValid(u model.UploadSession) *AppError {
	if !IsValidId(u.ID) {
		return NewAppError("UploadSessionIsValid", "model.upload_session.is_valid.id.app_error", nil, "", http.StatusBadRequest)
	}
	if err := u.Type.IsValid(); err != nil {
		return NewAppError("UploadSessionIsValid", "model.upload_session.is_valid.type.app_error", nil, "", http.StatusBadRequest)
	}
	if !IsValidId(u.UserID) {
		return NewAppError("UploadSessionIsValid", "model.upload_session.is_valid.user_id.app_error", nil, "", http.StatusBadRequest)
	}
	if u.FileName == "" {
		return NewAppError("UploadSessionIsValid", "model.upload_session.is_valid.filename.app_error", nil, "", http.StatusBadRequest)
	}
	if u.FileSize <= 0 {
		return NewAppError("UploadSessionIsValid", "model.upload_session.is_valid.file_size.app_error", nil, "", http.StatusBadRequest)
	}
	if u.FileOffset < 0 || u.FileOffset > u.FileSize {
		return NewAppError("UploadSessionIsValid", "model.upload_session.is_valid.file_offset.app_error", nil, "", http.StatusBadRequest)
	}
	return nil
}
