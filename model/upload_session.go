package model

// UploadType defines the type of an upload.
type UploadType string

const (
	UploadTypeAttachment UploadType = "attachment"
	UploadTypeImport     UploadType = "import"
)

// UploadSession contains information used to keep track of a file upload.
type UploadSession struct {
	Id         string     `json:"id"`
	Type       UploadType `json:"type"`
	CreateAt   int64      `json:"create_at"`
	UserID     string     `json:"user_id"`
	FileName   string     `json:"filename"`
	Path       string     `json:"path"`
	FileSize   int64      `json:"file_size"`
	FileOffset int64      `json:"file_offset"`
}
