package model

const (
	MaxImageSize = int64(6048 * 4032) // 24 megapixels, roughly 36MB as a raw image
)

type FileUploadResponse struct {
	FileInfos []*FileInfo `json:"file_infos"`
	ClientIds []string    `json:"client_ids"`
}

func (o *FileUploadResponse) ToJSON() string {
	return ModelToJson(o)
}
