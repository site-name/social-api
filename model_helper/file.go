package model_helper

import (
	"image"
	"image/gif"
	"io"
	"mime"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/modules/model_types"
)

func UploadSessionPreSave(u *model.UploadSession) {
	if u.ID == "" {
		u.ID = NewId()
	}
	if u.CreatedAt == 0 {
		u.CreatedAt = GetMillis()
	}
	UploadSessionCommonPre(u)
}

func UploadSessionCommonPre(u *model.UploadSession) {
	u.FileName = SanitizeUnicode(u.FileName)
}

type UploadSessionFilterOption struct {
	CommonQueryOptions
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

func FileInfoPreSave(f *model.FileInfo) {
	if f.ID == "" {
		f.ID = NewId()
	}
	if f.CreatedAt == 0 {
		f.CreatedAt = GetMillis()
	}
	f.UpdatedAt = f.CreatedAt
	fileInfoCommonPre(f)
}

func fileInfoCommonPre(f *model.FileInfo) {
	f.Name = SanitizeUnicode(f.Name)
}

func FileInfoPreUpdate(f *model.FileInfo) {
	f.UpdatedAt = GetMillis()
	fileInfoCommonPre(f)
}

func FileInfoIsValid(f model.FileInfo) *AppError {
	if !IsValidId(f.ID) {
		return NewAppError("FileInfoIsValid", "model.file_info.is_valid.id.app_error", nil, "", http.StatusBadRequest)
	}
	if !IsValidId(f.CreatorID) && f.CreatorID != "nouser" {
		return NewAppError("FileInfoIsValid", "model.file_info.is_valid.creator_id.app_error", nil, "", http.StatusBadRequest)
	}
	if f.ParentID != "" && !IsValidId(f.ParentID) {
		return NewAppError("FileInfoIsValid", "model.file_info.is_valid.parent_id.app_error", nil, "", http.StatusBadRequest)
	}
	if f.Path == "" {
		return NewAppError("FileInfoIsValid", "model.file_info.is_valid.path.app_error", nil, "", http.StatusBadRequest)
	}
	if f.CreatedAt <= 0 {
		return NewAppError("FileInfoIsValid", "model.file_info.is_valid.created_at.app_error", nil, "", http.StatusBadRequest)
	}
	if f.UpdatedAt <= 0 {
		return NewAppError("FileInfoIsValid", "model.file_info.is_valid.updated_at.app_error", nil, "", http.StatusBadRequest)
	}
	return nil
}

type FileInfoFilterOption struct {
	CommonQueryOptions
}

func FileInfoIsImage(f model.FileInfo) bool {
	return strings.HasPrefix(f.MimeType, "image")
}

func NewFileInfo(name string) *model.FileInfo {
	info := &model.FileInfo{
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

func GetInfoForBytes(name string, data io.ReadSeeker, size int) (*model.FileInfo, *AppError) {
	info := &model.FileInfo{
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

	if FileInfoIsImage(*info) {
		if config, _, err := image.DecodeConfig(data); err == nil {
			info.Width = model_types.NewNullInt(config.Width)
			info.Height = model_types.NewNullInt(config.Height)

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

func GetEtagForFileInfos(infos model.FileInfoSlice) string {
	if len(infos) == 0 {
		return Etag()
	}

	var maxUpdateAt int64

	for _, info := range infos {
		if info.UpdatedAt > maxUpdateAt {
			maxUpdateAt = info.UpdatedAt
		}
	}

	return Etag(infos[0].ParentID, maxUpdateAt)
}

func ExportFilePreSave(f *model.ExportFile) {
	if f.ID == "" {
		f.ID = NewId()
	}
	if f.CreatedAt == 0 {
		f.CreatedAt = GetMillis()
	}
	f.CreatedAt = GetMillis()
	f.UpdatedAt = f.CreatedAt
}

func ExportFileIsValid(f model.ExportFile) *AppError {
	if !IsValidId(f.ID) {
		return NewAppError("ExportFileIsValid", "model.export_file.is_valid.id.app_error", nil, "", http.StatusBadRequest)
	}
	if f.CreatedAt <= 0 {
		return NewAppError("ExportFileIsValid", "model.export_file.is_valid.created_at.app_error", nil, "", http.StatusBadRequest)
	}
	if f.UpdatedAt <= 0 {
		return NewAppError("ExportFileIsValid", "model.export_file.is_valid.updated_at.app_error", nil, "", http.StatusBadRequest)
	}
	if !f.UserID.IsNil() && !IsValidId(*f.UserID.String) {
		return NewAppError("ExportFileIsValid", "model.export_file.is_valid.user_id.app_error", nil, "", http.StatusBadRequest)
	}
	return nil
}
