package model

import (
	""
)

const (
	FILEINFO_SORT_BY_CREATED = "CreateAt"
	FILEINFO_SORT_BY_SIZE    = "Size"
)

// GetFileInfosOptions contains options for getting FileInfos
type GetFileInfoOptions struct {
	UserIds []string `json:"user_ids"`
	// ChannelIds []string `json:"channel_ids"`
	Since int64 `json:"since"`
	IncludeDeleted bool `json:"include_deleted"`
	SortBy string `json:"sort_by"`
	SortDescending bool `json:"sort_descending"`
}

type FileInfo struct {
	Id              string  `json:"id"`
	CreatorId       string  `json:"user_id"`
	PostId          string  `json:"post_id,omitempty"`
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

func (fi *FileInfo) ToJson() string {
	// b, _ := json.JSON.New
}
