package csv

import (
	"context"

	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model/account"
	"github.com/sitename/sitename/model/csv"
)

// GetDefaultExportPayload returns a map for mapping
func (a *ServiceCsv) GetDefaultExportPayload(exportFile csv.ExportFile) (map[string]interface{}, *model.AppError) {
	var (
		user   *account.User
		appErr *model.AppError
	)

	if exportFile.UserID != nil {
		user, appErr = a.srv.AccountService().UserById(context.Background(), *exportFile.UserID)
	}
	if appErr != nil {
		return nil, appErr
	}

	return map[string]interface{}{
		"user_id":    user.Id,
		"user_email": user.Email,
		"id":         exportFile.Id,
		"status":     nil,
		"message":    nil,
		"created_at": exportFile.CreateAt,
		"updated_at": exportFile.UpdateAt,
	}, nil
}
