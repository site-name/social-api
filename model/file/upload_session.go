package file

import (
	"fmt"
	"io"

	"github.com/sitename/sitename/model"
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
	Id         string     `json:"id"`          // The unique identifier for the session.
	Type       UploadType `json:"type"`        // The type of the upload.
	CreateAt   int64      `json:"create_at"`   // The timestamp of creation.
	UserID     string     `json:"user_id"`     // The id of the user performing the upload.
	FileName   string     `json:"filename"`    // The name of the file to upload.
	Path       string     `json:"-"`           // The path where the file is stored.
	FileSize   int64      `json:"file_size"`   // The size of the file to upload.
	FileOffset int64      `json:"file_offset"` // The amount of received data in bytes. If equal to FileSize it means the upload has finished.
}

func (us *UploadSession) ToJson() string {
	return model.ModelToJson(us)
}

func UploadSessionFromJson(data io.Reader) *UploadSession {
	var us *UploadSession
	model.ModelFromJson(&us, data)
	return us
}

// UploadSessionsToJson serializes a list of UploadSession into JSON and
// returns it as string.
func UploadSessionsToJson(uss []*UploadSession) string {
	return model.ModelToJson(uss)
}

// UploadSessionsFromJson deserializes a list of UploadSession from JSON data.
func UploadSessionsFromJson(data io.Reader) []*UploadSession {
	var uss []*UploadSession
	if err := model.ModelFromJson(&uss, data); err != nil {
		return nil
	}
	return uss
}

// PreSave is a utility function used to fill required information.
func (us *UploadSession) PreSave() {
	if us.Id == "" {
		us.Id = model.NewId()
	}

	us.CreateAt = model.GetMillis()
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

func (us *UploadSession) IsValid() *model.AppError {
	outer := model.CreateAppErrorForModel(
		"model.upload_session.is_valid.%s.app_error",
		"upload_session_id=",
		"UploadSession.IsValid",
	)
	if !model.IsValidId(us.Id) {
		return outer("id", nil)
	}
	if err := us.Type.IsValid(); err != nil {
		return outer("type", &us.Id)
	}
	if !model.IsValidId(us.UserID) && us.UserID != UploadNoUserID {
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
