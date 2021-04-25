package model

import (
	"io"

	"github.com/sitename/sitename/modules/json"
)

const (
	MaxImageSize = int64(6048 * 4032) // 24 megapixels, roughly 36MB as a raw image
)

type FileUploadResponse struct {
	FileInfos []*FileInfo `json:"file_infos"`
	ClientIds []string    `json:"client_ids"`
}

func FileUploadResponseFromJson(data io.Reader) *FileUploadResponse {
	var o *FileUploadResponse
	json.JSON.NewDecoder(data).Decode(&o)
	return o
}

func (o *FileUploadResponse) ToJson() string {
	b, _ := json.JSON.Marshal(o)
	return string(b)
}
