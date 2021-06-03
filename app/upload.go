package app

import (
	"net/http"

	"github.com/sitename/sitename/model"
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
