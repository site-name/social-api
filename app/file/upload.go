package file

import (
	"io"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/sitename/sitename/app/request"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model_helper"
	"github.com/sitename/sitename/modules/slog"
	"github.com/sitename/sitename/store"
)

const minFirstPartSize = 5 * 1024 * 1024 // 5MB
const IncompleteUploadSuffix = ".tmp"

func (a *ServiceFile) GetUploadSessionsForUser(userID string) ([]*model.UploadSession, *model_helper.AppError) {
	uss, err := a.srv.Store.UploadSession().GetForUser(userID)
	var (
		statusCode int
		errMsg     string
	)
	if err != nil {
		statusCode = http.StatusInternalServerError
		errMsg = err.Error()
	} else if len(uss) == 0 {
		statusCode = http.StatusNotFound
	}

	if statusCode != 0 {
		return nil, model_helper.NewAppError("GetUploadSessionsForUser", "app.file_error_finding_upload_sessions_for_user.app_error", nil, errMsg, statusCode)
	}

	return uss, nil
}

func (a *ServiceFile) GetUploadSession(uploadId string) (*model.UploadSession, *model_helper.AppError) {
	us, err := a.srv.Store.UploadSession().Get(uploadId)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if _, ok := err.(*store.ErrNotFound); ok {
			statusCode = http.StatusNotFound
		}

		return nil, model_helper.NewAppError("GetUploadSession", "app.file.error_finding_upload_session_by_id.app_error", nil, err.Error(), statusCode)
	}

	return us, nil
}

// func (a *ServiceFile) CreateUploadSession(us *model.UploadSession) (*model.UploadSession, *model_helper.AppError) {
// 	if us.FileSize > *a.srv.Config().FileSettings.MaxFileSize {
// 		return nil, model_helper.NewAppError(
// 			"CreateUploadSession",
// 			"app.upload.create.upload_too_large.app_error",
// 			nil, "", http.StatusRequestEntityTooLarge,
// 		)
// 	}

// 	us.FileOffset = 0
// 	now := time.Now()
// 	us.CreateAt = model.GetMillisForTime(now)
// 	if us.Type == model.UploadTypeAttachment {
// 		us.Path = now.Format("20060102") + "/teams/noteam/channels/" + us.ChannelId + "/users/" + us.UserID + "/" + us.Id + "/" + filepath.Base(us.FileName)
// 	} else if us.Type == model.UploadTypeImport {
// 		us.Path = filepath.Clean(*a.srv.Config().ImportSettings.Directory) + "/" + us.Id + "_" + filepath.Base(us.FileName)
// 	}
// 	if err := us.IsValid(); err != nil {
// 		return nil, err
// 	}

// 	if us.Type == model.UploadTypeAttachment {
// 		channel, err := a.GetChannel(us.ChannelId)
// 		if err != nil {
// 			return nil, model_helper.NewAppError("CreateUploadSession", "app.upload.create.incorrect_channel_id.app_error",
// 				map[string]interface{}{"channelId": us.ChannelId}, "", http.StatusBadRequest)
// 		}
// 		if channel.DeleteAt != 0 {
// 			return nil, model_helper.NewAppError("CreateUploadSession", "app.upload.create.cannot_upload_to_deleted_channel.app_error",
// 				map[string]interface{}{"channelId": us.ChannelId}, "", http.StatusBadRequest)
// 		}
// 	}

// 	us, storeErr := a.srv.Store.UploadSession().Save(us)
// 	if storeErr != nil {
// 		return nil, model_helper.NewAppError("CreateUploadSession", "app.upload.create.save.app_error", nil, storeErr.Error(), http.StatusInternalServerError)
// 	}

// 	return us, nil
// }

func (a *ServiceFile) UploadData(c *request.Context, us *model.UploadSession, rd io.Reader) (*model.FileInfo, *model_helper.AppError) {
	// prevent more than one caller to upload data at the same time for a given upload session.
	// This is to avoid possible inconsistencies.
	a.uploadLockMapMut.Lock()
	locked := a.uploadLockMap[us.Id]
	if locked {
		// session lock is already taken, return error.
		a.uploadLockMapMut.Unlock()
		return nil, model_helper.NewAppError("UploadData", "app.upload.upload_data.concurrent.app_error", nil, "", http.StatusBadRequest)
	}
	// grab the session lock.
	a.uploadLockMap[us.Id] = true
	a.uploadLockMapMut.Unlock()

	// reset the session lock on exit.
	defer func() {
		a.uploadLockMapMut.Lock()
		delete(a.uploadLockMap, us.Id)
		a.uploadLockMapMut.Unlock()
	}()

	// fetch the session from store to check for inconsistencies.
	if storedSession, err := a.GetUploadSession(us.Id); err != nil {
		return nil, err
	} else if us.FileOffset != storedSession.FileOffset {
		return nil, model_helper.NewAppError("UploadData", "app.upload.upload_data.concurrent.app_error", nil, "FileOffset mismatch", http.StatusBadRequest)
	}

	uploadPath := us.Path
	if us.Type == model.UploadTypeImport {
		uploadPath += IncompleteUploadSuffix
	}

	// make sure it's not possible to upload more data than what is expected.
	lr := &io.LimitedReader{
		R: rd,
		N: us.FileSize - us.FileOffset,
	}
	var err *model_helper.AppError
	var written int64
	if us.FileOffset == 0 {
		// new upload
		written, err = a.WriteFile(lr, uploadPath)
		if err != nil && written == 0 {
			return nil, err
		}
		if written < minFirstPartSize && written != us.FileSize {
			a.RemoveFile(uploadPath)
			var errStr string
			if err != nil {
				errStr = err.Error()
			}
			return nil, model_helper.NewAppError("UploadData", "app.upload.upload_data.first_part_too_small.app_error", map[string]interface{}{"Size": minFirstPartSize}, errStr, http.StatusBadRequest)
		}
	} else if us.FileOffset < us.FileSize {
		// resume upload
		written, err = a.AppendFile(lr, uploadPath)
	}
	if written > 0 {
		us.FileOffset += written
		if storeErr := a.srv.Store.UploadSession().Update(us); storeErr != nil {
			return nil, model_helper.NewAppError("UploadData", "app.upload.upload_data.update.app_error", nil, storeErr.Error(), http.StatusInternalServerError)
		}
	}
	if err != nil {
		return nil, err
	}

	// upload is incomplete
	if us.FileOffset != us.FileSize {
		return nil, nil
	}

	// upload is done, create FileInfo
	f, err := a.FileReader(uploadPath)
	if err != nil {
		return nil, model_helper.NewAppError("UploadData", "app.upload.upload_data.read_file.app_error", nil, err.Error(), http.StatusInternalServerError)
	}

	info, err := model.GetInfoForBytes(us.FileName, f, int(us.FileSize))
	f.Close()
	if err != nil {
		return nil, err
	}

	info.CreatorId = us.UserID
	info.Path = us.Path

	// info.RemoteId = model.GetPointerOfValue(us.RemoteId)
	// if us.ReqFileId != "" {
	// 	info.Id = us.ReqFileId
	// }

	// run plugins upload hook
	// if err := a.runPluginsHook(c, info, file); err != nil {
	// 	return nil, err
	// }

	// image post-processing
	if info.IsImage() {
		if limitErr := checkImageResolutionLimit(info.Width, info.Height, *a.srv.Config().FileSettings.MaxImageResolution); limitErr != nil {
			return nil, model_helper.NewAppError(
				"uploadData",
				"app.upload.upload_data.large_image.app_error",
				map[string]interface{}{
					"Filename": us.FileName,
					"Width":    info.Width,
					"Height":   info.Height,
				}, "", http.StatusBadRequest)
		}

		nameWithoutExtension := info.Name[:strings.LastIndex(info.Name, ".")]
		info.PreviewPath = filepath.Dir(info.Path) + "/" + nameWithoutExtension + "_preview.jpg"
		info.ThumbnailPath = filepath.Dir(info.Path) + "/" + nameWithoutExtension + "_thumb.jpg"
		imgData, fileErr := a.ReadFile(uploadPath)
		if fileErr != nil {
			return nil, fileErr
		}
		a.HandleImages([]string{info.PreviewPath}, []string{info.ThumbnailPath}, [][]byte{imgData})
	}

	if us.Type == model.UploadTypeImport {
		if err := a.MoveFile(uploadPath, us.Path); err != nil {
			return nil, model_helper.NewAppError("UploadData", "app.upload.upload_data.move_file.app_error", nil, err.Error(), http.StatusInternalServerError)
		}
	}

	var storeErr error
	if info, storeErr = a.srv.Store.FileInfo().Upsert(info); storeErr != nil {
		if appErr, ok := storeErr.(*model_helper.AppError); ok {
			return nil, appErr
		}
		return nil, model_helper.NewAppError("uploadData", "app.upload.upload_data.save.app_error", nil, storeErr.Error(), http.StatusInternalServerError)
	}

	if *a.srv.Config().FileSettings.ExtractContent {
		infoCopy := *info
		a.srv.Go(func() {
			err := a.ExtractContentFromFileInfo(&infoCopy)
			if err != nil {
				slog.Error("Failed to extract file content", slog.Err(err), slog.String("fileInfoId", infoCopy.Id))
			}
		})
	}

	// delete upload session
	if storeErr := a.srv.Store.UploadSession().Delete(us.Id); storeErr != nil {
		slog.Warn("Failed to delete UploadSession", slog.Err(storeErr))
	}

	return info, nil
}
