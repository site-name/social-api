package app

import (
	"net/http"

	"github.com/sitename/sitename/model"
)

func(a *App) DeactivateGuests() *model.AppError {
	userIDs, err := a.Srv().Store.User().DeactivateGuests()
	if err != nil {
		return model.NewAppError("DeactivateGuests", "app.user.update_active_for_multiple_users.updating.app_error", nil, err.Error(), http.StatusInternalServerError)
	}

	for _, userID := range userIDs {
		if err := a.user
	}
}

func (a *App) userDeactivated(userID string) *model.AppError {
	if err := a.Revo
}
