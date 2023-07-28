package model

import (
	"image"
	"image/gif"
	"io"
	"mime"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/Masterminds/squirrel"
	"github.com/samber/lo"
	"gorm.io/gorm"
)

const (
	FILEINFO_SORT_BY_CREATED = "CreateAt"
	FILEINFO_SORT_BY_SIZE    = "Size"
)

// GetFileInfosOptions contains options for getting FileInfos
type GetFileInfosOptions struct {
	Conditions squirrel.Sqlizer
	OrderBy    string // E.g "CreateAt ASC"

	Limit  uint64 // if 0, no limit
	Offset uint64 // if 0, no offset
}

type FileForIndexing struct {
	FileInfo
	Content string `json:"content"`
}

type FileInfo struct {
	Id              string  `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid();column:Id"`
	CreatorId       string  `json:"user_id" gorm:"type:uuid;column:CreatorId"`
	ParentID        string  `json:"parent_id,omitempty" gorm:"type:uuid;column:ParentID"` // can be a product's id, comment's id
	CreateAt        int64   `json:"create_at" gorm:"type:bigint;column:CreateAt;autoCreateTime:milli"`
	UpdateAt        int64   `json:"update_at" gorm:"type:bigint;column:UpdateAt;autoUpdateTime:milli"`
	DeleteAt        int64   `json:"delete_at" gorm:"type:bigint;column:DeleteAt"`
	Path            string  `json:"-" gorm:"type:varchar(512);column:Path"`          // not sent back to the client
	ThumbnailPath   string  `json:"-" gorm:"type:varchar(512);column:ThumbnailPath"` // not sent back to the client
	PreviewPath     string  `json:"-" gorm:"type:varchar(512);column:PreviewPath"`   // not sent back to the client
	Name            string  `json:"name" gorm:"type:varchar(256);column:Name"`
	Extension       string  `json:"extension" gorm:"type:varchar(64);column:Extension"`
	Size            int64   `json:"size" gorm:"type:bigint;column:Size"`
	MimeType        string  `json:"mime_type" gorm:"type:varchar(256);column:MimeType"`
	Width           int     `json:"width,omitempty" gorm:"column:Width"`
	Height          int     `json:"height,omitempty" gorm:"column:Height"`
	HasPreviewImage bool    `json:"has_preview_image,omitempty" gorm:"column:HasPreviewImage"`
	MiniPreview     *[]byte `json:"mini_preview" gorm:"column:MiniPreview;type:bytea"` // declared as *[]byte to avoid postgres/mysql differences in deserialization
	Content         string  `json:"-" gorm:"column:Content"`
	RemoteId        *string `json:"remote_id" gorm:"column:RemoteId;type:uuid"`
}

func (c *FileInfo) BeforeCreate(_ *gorm.DB) error { c.PreSave(); return c.IsValid() }
func (c *FileInfo) BeforeUpdate(_ *gorm.DB) error { c.CreateAt = 0; return c.IsValid() }
func (c *FileInfo) TableName() string             { return FileInfoTableName }

type FileInfos []*FileInfo

func (fi *FileInfo) ToJSON() string {
	return ModelToJson(fi)
}

func (fi *FileInfo) PreSave() {
	if fi.UpdateAt < fi.CreateAt {
		fi.UpdateAt = fi.CreateAt
	}

	if fi.RemoteId == nil {
		fi.RemoteId = NewPrimitive("")
	}
}

func (fi *FileInfo) DeepCopy() *FileInfo {
	if fi == nil {
		return nil
	}

	res := *fi
	if fi.RemoteId != nil {
		res.RemoteId = NewPrimitive(*fi.RemoteId)
	}
	if fi.MiniPreview != nil {
		m := append([]byte{}, *fi.MiniPreview...)
		res.MiniPreview = &m
	}
	return &res
}

func (f FileInfos) DeepCopy() FileInfos {
	return lo.Map(f, func(fi *FileInfo, _ int) *FileInfo { return fi.DeepCopy() })
}

func (fi *FileInfo) IsValid() *AppError {
	outer := CreateAppErrorForModel(
		"model.file_info.is_valid.%s.app_error",
		"file_info_id=",
		"FileInfo.IsValid",
	)

	if !IsValidId(fi.CreatorId) && fi.CreatorId != "nouser" {
		return outer("creator_id", &fi.Id)
	}
	if fi.ParentID != "" && !IsValidId(fi.ParentID) {
		return outer("product_id", &fi.Id)
	}
	if fi.Path == "" {
		return outer("path", &fi.Id)
	}

	return nil
}

// IsImage check if fileInfo's MimeType is prefixed with "image"
func (fi *FileInfo) IsImage() bool {
	return strings.HasPrefix(fi.MimeType, "image")
}

// NewInfo create new FileInfo, attributes 'Name', 'MimeType' and 'Extension' are created
func NewInfo(name string) *FileInfo {
	info := &FileInfo{
		Name: name,
	}

	extension := strings.ToLower(filepath.Ext(name))
	info.MimeType = mime.TypeByExtension(extension)

	if extension != "" && extension[0] == '.' {
		// The client expects a file extension without the leading period
		info.Extension = extension[1:]
	} else {
		info.Extension = extension
	}

	return info
}

func GetInfoForBytes(name string, data io.ReadSeeker, size int) (*FileInfo, *AppError) {
	info := &FileInfo{
		Name: name,
		Size: int64(size),
	}

	var err *AppError

	extension := strings.ToLower(filepath.Ext(name))
	info.MimeType = mime.TypeByExtension(extension)

	if extension != "" && extension[0] == '.' {
		info.Extension = extension[1:]
	} else {
		info.Extension = extension
	}

	if info.IsImage() {
		if config, _, err := image.DecodeConfig(data); err == nil {
			info.Width = config.Width
			info.Height = config.Height

			if info.MimeType == "image/gif" {
				data.Seek(0, io.SeekStart)
				gifConfig, err := gif.DecodeAll(data)
				if err != nil {
					// Still return the rest of the info even though it doesn't appear to be an actual gif
					info.HasPreviewImage = true
					return info, NewAppError("GetInfoForBytes", "file_info.get.gif.app_error", nil, err.Error(), http.StatusBadRequest)
				}
				info.HasPreviewImage = len(gifConfig.Image) == 1
			} else {
				info.HasPreviewImage = true
			}
		}
	}

	return info, err
}

func GetEtagForFileInfos(infos []*FileInfo) string {
	if len(infos) == 0 {
		return Etag()
	}

	var maxUpdateAt int64

	for _, info := range infos {
		if info.UpdateAt > maxUpdateAt {
			maxUpdateAt = info.UpdateAt
		}
	}

	return Etag(infos[0].ParentID, maxUpdateAt)
}
