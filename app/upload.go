package app

import (
	"errors"
	"io"
	"net/http"

	"github.com/sitename/sitename/app/request"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/store"
)

func (a *App) GetUploadSessionsForUser(userID string) ([]*model.UploadSession, *model.AppError) {
	uss, err := a.Srv().Store.UploadSession().GetForUser(userID)
	if err != nil {
		return nil, model.NewAppError(
			"GetUploadsForUser",
			"app.upload.get_for_user.app_error",
			nil,
			err.Error(),
			http.StatusInternalServerError,
		)
	}
	return uss, nil
}

func (a *App) GetUploadSession(uploadId string) (*model.UploadSession, *model.AppError) {
	us, err := a.Srv().Store.UploadSession().Get(uploadId)
	if err != nil {
		var nfErr *store.ErrNotFound
		switch {
		case errors.As(err, &nfErr):
			return nil, model.NewAppError("GetUpload", "app.upload.get.app_error",
				nil, nfErr.Error(), http.StatusNotFound)
		default:
			return nil, model.NewAppError("GetUpload", "app.upload.get.app_error",
				nil, err.Error(), http.StatusInternalServerError)
		}
	}
	return us, nil
}

func (a *App) CreateUploadSession(us *model.UploadSession) (*model.UploadSession, *model.AppError) {
	// if us.FileSize > *a.Config().FileSettings.MaxFileSize {
	// 	return nil, model.NewAppError("CreateUploadSession", "app.upload.create.upload_too_large.app_error",
	// 		nil, "", http.StatusRequestEntityTooLarge)
	// }

	// us.FileOffset = 0
	// now := time.Now()
	// us.CreateAt = model.GetMillisForTime(now)
	// if us.Type == model.UploadTypeAttachment {
	// 	us.Path = now.Format("20060102") + "/teams/noteam/channels/" + us.ChannelId + "/users/" + us.UserId + "/" + us.Id + "/" + filepath.Base(us.Filename)
	// } else if us.Type == model.UploadTypeImport {
	// 	us.Path = filepath.Clean(*a.Config().ImportSettings.Directory) + "/" + us.Id + "_" + filepath.Base(us.Filename)
	// }
	// if err := us.IsValid(); err != nil {
	// 	return nil, err
	// }

	// if us.Type == model.UploadTypeAttachment {
	// 	channel, err := a.GetChannel(us.ChannelId)
	// 	if err != nil {
	// 		return nil, model.NewAppError("CreateUploadSession", "app.upload.create.incorrect_channel_id.app_error",
	// 			map[string]interface{}{"channelId": us.ChannelId}, "", http.StatusBadRequest)
	// 	}
	// 	if channel.DeleteAt != 0 {
	// 		return nil, model.NewAppError("CreateUploadSession", "app.upload.create.cannot_upload_to_deleted_channel.app_error",
	// 			map[string]interface{}{"channelId": us.ChannelId}, "", http.StatusBadRequest)
	// 	}
	// }

	// us, storeErr := a.Srv().Store.UploadSession().Save(us)
	// if storeErr != nil {
	// 	return nil, model.NewAppError("CreateUploadSession", "app.upload.create.save.app_error", nil, storeErr.Error(), http.StatusInternalServerError)
	// }

	// return us, nil
	panic("not implemented") // TODO: fixme
}

func (a *App) UploadData(c *request.Context, us *model.UploadSession, rd io.Reader) (*model.FileInfo, *model.AppError) {
	// prevent more than one caller to upload data at the same time for a given upload session.
	// This is to avoid possible inconsistencies.
	a.Srv().uploadLockMapMut.Lock()
	locked := a.Srv().uploadLockMap[us.Id]
	if locked {
		// session lock is already taken, return error.
		a.Srv().uploadLockMapMut.Unlock()
		return nil, model.NewAppError("UploadData", "app.upload.upload_data.concurrent.app_error",
			nil, "", http.StatusBadRequest)
	}
	// grab the session lock.
	a.Srv().uploadLockMap[us.Id] = true
	a.Srv().uploadLockMapMut.Unlock()

	// reset the session lock on exit.
	defer func() {
		a.Srv().uploadLockMapMut.Lock()
		delete(a.Srv().uploadLockMap, us.Id)
		a.Srv().uploadLockMapMut.Unlock()
	}()
}
