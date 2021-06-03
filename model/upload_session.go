package model

import (
	"fmt"
	"io"
)

// UploadType defines the type of an upload.
type UploadType string

const (
	UploadTypeAttachment UploadType = "attachment"
	UploadTypeImport     UploadType = "import"
)

// UploadNoUserID is a "fake" user id used by the API layer when in local mode.
const UploadNoUserID = "nouser"

// UploadSession contains information used to keep track of a file upload.
type UploadSession struct {
	Id         string     `json:"id"`
	Type       UploadType `json:"type"`
	CreateAt   int64      `json:"create_at"`
	UserID     string     `json:"user_id"`
	FileName   string     `json:"filename"`
	Path       string     `json:"-"`
	FileSize   int64      `json:"file_size"`
	FileOffset int64      `json:"file_offset"`
}

func (us *UploadSession) ToJson() string {
	return ModelToJson(us)
}

func UploadSessionFromJson(data io.Reader) *UploadSession {
	var us *UploadSession
	ModelFromJson(&us, data)
	return us
}

// UploadSessionsToJson serializes a list of UploadSession into JSON and
// returns it as string.
func UploadSessionsToJson(uss []*UploadSession) string {
	return ModelToJson(uss)
}

// UploadSessionsFromJson deserializes a list of UploadSession from JSON data.
func UploadSessionsFromJson(data io.Reader) []*UploadSession {
	var uss []*UploadSession
	if err := ModelFromJson(&uss, data); err != nil {
		return nil
	}
	return uss
}

// PreSave is a utility function used to fill required information.
func (us *UploadSession) PreSave() {
	if us.Id == "" {
		us.Id = NewId()
	}

	if us.CreateAt == 0 {
		us.CreateAt = GetMillis()
	}
}

// IsValid validates an UploadType. It returns an error in case of
// failure.
func (t UploadType) IsValid() error {
	switch t {
	case UploadTypeAttachment:
		return nil
	case UploadTypeImport:
		return nil
	default:
	}
	return fmt.Errorf("invalid UploadType %s", t)
}

func (us *UploadSession) IsValid() *AppError {
	outer := CreateAppErrorForModel(
		"model.upload_session.is_valid.%s.app_error",
		"upload_session_id=",
		"UploadSession.IsValid",
	)
	if !IsValidId(us.Id) {
		return outer("id", nil)
	}
	if err := us.Type.IsValid(); err != nil {
		return outer("type", &us.Id)
	}
	if !IsValidId(us.UserID) && us.UserID != UploadNoUserID {
		return outer("user_id", &us.Id)
	}
	if us.CreateAt == 0 {
		return outer("create_at", &us.Id)
	}
	if us.FileName == "" {
		return outer("file_name", &us.Id)
	}
	if us.FileSize <= 0 {
		return outer("file_size", &us.Id)
	}
	if us.FileOffset < 0 || us.FileOffset > us.FileSize {
		return outer("file_offset", &us.Id)
	}
	if us.Path == "" {
		return outer("path", &us.Id)
	}

	return nil
}
