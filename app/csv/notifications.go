package csv

import (
	"context"

	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model/account"
	"github.com/sitename/sitename/model/csv"
)

func (a *AppCsv) GetDefaultExportPayload(exportFile *csv.ExportFile) (map[string]interface{}, *model.AppError) {
	var user *account.User
	var appErr *model.AppError

	if exportFile.UserID != nil {
		user, appErr = a.AccountApp().UserById(context.Background(), *exportFile.UserID)
	}
	if appErr != nil {
		return nil, appErr
	}

	userID := ""
	userEmail := ""
	if user != nil {
		userID = user.Id
		userEmail = user.Email
	}

	return map[string]interface{}{
		"user_id":    userID,
		"user_email": userEmail,
		"id":         exportFile.Id,
		"status":     "",
	}, nil
}