package account

import (
	"errors"
	"net/http"

	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/store"
)

func (a *ServiceAccount) GetPreferencesForUser(userID string) (model.Preferences, *model.AppError) {
	preferences, err := a.srv.Store.Preference().GetAll(userID)
	if err != nil {
		return nil, model.NewAppError("GetPreferencesForUser", "app.account.preferences_for_user.app_error", nil, err.Error(), http.StatusInternalServerError)
	}
	return preferences, nil
}

func (a *ServiceAccount) GetPreferenceByCategoryForUser(userID string, category string) (model.Preferences, *model.AppError) {
	preferences, err := a.srv.Store.Preference().GetCategory(userID, category)
	if err != nil {
		return nil, model.NewAppError("GetPreferenceByCategoryForUser", "app.account.get_category.app_error", nil, err.Error(), http.StatusInternalServerError)
	}

	return preferences, nil
}

func (a *ServiceAccount) GetPreferenceByCategoryAndNameForUser(userID string, category string, preferenceName string) (*model.Preference, *model.AppError) {
	preference, err := a.srv.Store.Preference().Get(userID, category, preferenceName)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if _, ok := err.(*store.ErrNotFound); ok {
			statusCode = http.StatusNotFound
		}
		return nil, model.NewAppError("GetPreferenceByCategoryAndNameForUser", "app.account.get_preference.app_error", nil, err.Error(), statusCode)
	}
	return preference, nil
}

func (a *ServiceAccount) UpdatePreferences(userID string, preferences model.Preferences) *model.AppError {
	for _, preference := range preferences {
		if userID != preference.UserId {
			return model.NewAppError("savePreferences", "api.preference.update_preferences.set.app_error", nil,
				"userId="+userID+", preference.UserId="+preference.UserId, http.StatusForbidden)
		}
	}

	if err := a.srv.Store.Preference().Save(preferences); err != nil {
		var appErr *model.AppError
		switch {
		case errors.As(err, &appErr):
			return appErr
		default:
			return model.NewAppError("UpdatePreferences", "app.preference.save.updating.app_error", nil, err.Error(), http.StatusBadRequest)
		}
	}

	return nil
}

func (a *ServiceAccount) DeletePreferences(userID string, preferences model.Preferences) *model.AppError {
	for _, preference := range preferences {
		if userID != preference.UserId {
			err := model.NewAppError(
				"DeletePreferences",
				"api.preference.delete_preferences.delete.app_error",
				nil, "userId="+userID+", preference.UserId="+preference.UserId,
				http.StatusForbidden,
			)
			return err
		}
	}

	for _, preference := range preferences {
		if err := a.srv.Store.Preference().Delete(userID, preference.Category, preference.Name); err != nil {
			return model.NewAppError("DeletePreferences", "app.preference.delete.app_error", nil, err.Error(), http.StatusBadRequest)
		}
	}

	return nil
}
