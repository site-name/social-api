package model

import (
	"image"
	"image/gif"
	"io"
	"mime"
	"net/http"
	"path/filepath"
	"strings"
)

const (
	FILEINFO_SORT_BY_CREATED = "CreateAt"
	FILEINFO_SORT_BY_SIZE    = "Size"
)

// GetFileInfosOptions contains options for getting FileInfos
type GetFileInfosOptions struct {
	UserIds        []string `json:"user_ids"`
	Since          int64    `json:"since"`
	IncludeDeleted bool     `json:"include_deleted"`
	SortBy         string   `json:"sort_by"`
	SortDescending bool     `json:"sort_descending"`
	ParentID       []string `json:"parent_id"`
}

type FileForIndexing struct {
	FileInfo
	Content string `json:"content"`
}

type FileInfo struct {
	Id              string  `json:"id"`
	CreatorId       string  `json:"user_id"`
	ParentID        string  `json:"parent_id,omitempty"` // can be a product's id, comment's id
	CreateAt        int64   `json:"create_at"`
	UpdateAt        int64   `json:"update_at"`
	DeleteAt        int64   `json:"delete_at"`
	Path            string  `json:"-"` // not sent back to the client
	ThumbnailPath   string  `json:"-"` // not sent back to the client
	PreviewPath     string  `json:"-"` // not sent back to the client
	Name            string  `json:"name"`
	Extension       string  `json:"extension"`
	Size            int64   `json:"size"`
	MimeType        string  `json:"mime_type"`
	Width           int     `json:"width,omitempty"`
	Height          int     `json:"height,omitempty"`
	HasPreviewImage bool    `json:"has_preview_image,omitempty"`
	MiniPreview     *[]byte `json:"mini_preview"` // declared as *[]byte to avoid postgres/mysql differences in deserialization
	Content         string  `json:"-"`
	RemoteId        *string `json:"remote_id"`
}

type FileInfos []*FileInfo

func (fi *FileInfo) ToJSON() string {
	return ModelToJson(fi)
}

func (fi *FileInfo) PreSave() {
	if fi.Id == "" {
		fi.Id = NewId()
	}

	if fi.CreateAt == 0 {
		fi.CreateAt = GetMillis()
	}
	if fi.UpdateAt < fi.CreateAt {
		fi.UpdateAt = fi.CreateAt
	}

	if fi.RemoteId == nil {
		fi.RemoteId = NewString("")
	}
}

func (fi *FileInfo) DeepCopy() *FileInfo {
	if fi == nil {
		return nil
	}

	res := *fi
	return &res
}

func (f FileInfos) DeepCopy() FileInfos {
	res := FileInfos{}
	for _, item := range f {
		res = append(res, item.DeepCopy())
	}

	return res
}

func (fi *FileInfo) IsValid() *AppError {
	outer := CreateAppErrorForModel(
		"file_info.is_valid.%s.app_error",
		"file_info_id=",
		"FileInfo.IsValid",
	)
	if !IsValidId(fi.Id) {
		return outer("id", nil)
	}
	if !IsValidId(fi.CreatorId) && fi.CreatorId != "nouser" {
		return outer("creator_id", &fi.Id)
	}
	if fi.ParentID != "" && !IsValidId(fi.ParentID) {
		return outer("product_id", &fi.Id)
	}
	if fi.CreateAt == 0 {
		return outer("create_at", &fi.Id)
	}
	if fi.UpdateAt == 0 {
		return outer("update_at", &fi.Id)
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