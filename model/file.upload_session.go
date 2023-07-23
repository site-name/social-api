package model

import (
	"fmt"

	"gorm.io/gorm"
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
	Id         string     `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid();column:Id"` // The unique identifier for the session.
	Type       UploadType `json:"type" gorm:"type:varchar(32);column:Type"`                           // The type of the upload.
	CreateAt   int64      `json:"create_at" gorm:"type:bigint;autoCreateTime:milli;column:CreateAt"`  // The timestamp of creation.
	UserID     string     `json:"user_id" gorm:"type:uuid;column:UserID"`                             // The id of the user performing the upload.
	FileName   string     `json:"filename" gorm:"type:varchar(256);column:FileName"`                  // The name of the file to upload.
	Path       string     `json:"-" gorm:"type:varchar(512);column:Path"`                             // The path where the file is stored.
	FileSize   int64      `json:"file_size" gorm:"type:bigint;column:FileSize"`                       // The size of the file to upload.
	FileOffset int64      `json:"file_offset" gorm:"type:bigint;column:FileOffset"`                   // The amount of received data in bytes. If equal to FileSize it means the upload has finished.
}

func (c *UploadSession) BeforeCreate(_ *gorm.DB) error { return c.IsValid() }
func (c *UploadSession) BeforeUpdate(_ *gorm.DB) error { return c.IsValid() }
func (c *UploadSession) TableName() string             { return UploadSessionTableName }

func (us *UploadSession) ToJSON() string {
	return ModelToJson(us)
}

// UploadSessionsToJson serializes a list of UploadSession into JSON and
// returns it as string.
func UploadSessionsToJson(uss []*UploadSession) string {
	return ModelToJson(uss)
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

	if err := us.Type.IsValid(); err != nil {
		return outer("type", &us.Id)
	}
	if !IsValidId(us.UserID) && us.UserID != UploadNoUserID {
		return outer("user_id", &us.Id)
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
