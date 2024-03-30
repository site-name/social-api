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

	if !exportFile.UserID.IsNil() {
		user, appErr = a.srv.Account.UserById(context.Background(), *exportFile.UserID.String)
	}
	if appErr != nil {
		return nil, appErr
	}

	return map[string]any{
		"user_id":    user.ID,
		"user_email": user.Email,
		"id":         exportFile.ID,
		"status":     nil,
		"message":    nil,
		"created_at": exportFile.CreatedAt,
		"updated_at": exportFile.UpdatedAt,
	}, nil
}
