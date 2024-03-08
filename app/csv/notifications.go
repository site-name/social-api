package csv

import (
	"context"

	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model_helper"
)

// GetDefaultExportPayload returns a map for mapping
func (a *ServiceCsv) GetDefaultExportPayload(exportFile model.ExportFile) (map[string]any, *model_helper.AppError) {
	var (
		user   *model.User
		appErr *model_helper.AppError
	)

	if exportFile.UserID != nil {
		user, appErr = a.srv.AccountService().UserById(context.Background(), *exportFile.UserID)
	}
	if appErr != nil {
		return nil, appErr
	}

	return map[string]any{
		"user_id":    user.Id,
		"user_email": user.Email,
		"id":         exportFile.Id,
		"status":     nil,
		"message":    nil,
		"created_at": exportFile.CreateAt,
		"updated_at": exportFile.UpdateAt,
	}, nil
}
