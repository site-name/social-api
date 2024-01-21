package account

import (
	"errors"
	"net/http"

	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model_helper"
	"github.com/sitename/sitename/store"
)

func (a *ServiceAccount) GetPreferencesForUser(userID string) (model.PreferenceSlice, *model_helper.AppError) {
	preferences, err := a.srv.Store.Preference().GetAll(userID)
	if err != nil {
		return nil, model_helper.NewAppError("GetPreferencesForUser", "app.account.preferences_for_user.app_error", nil, err.Error(), http.StatusInternalServerError)
	}
	return preferences, nil
}

func (a *ServiceAccount) GetPreferenceByCategoryForUser(userID string, category string) (model.PreferenceSlice, *model_helper.AppError) {
	preferences, err := a.srv.Store.Preference().GetCategory(userID, category)
	if err != nil {
		return nil, model_helper.NewAppError("GetPreferenceByCategoryForUser", "app.account.get_category.app_error", nil, err.Error(), http.StatusInternalServerError)
	}

	return preferences, nil
}

func (a *ServiceAccount) GetPreferenceByCategoryAndNameForUser(userID string, category string, preferenceName string) (*model.Preference, *model_helper.AppError) {
	preference, err := a.srv.Store.Preference().Get(userID, category, preferenceName)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if _, ok := err.(*store.ErrNotFound); ok {
			statusCode = http.StatusNotFound
		}
		return nil, model_helper.NewAppError("GetPreferenceByCategoryAndNameForUser", "app.account.get_preference.app_error", nil, err.Error(), statusCode)
	}
	return preference, nil
}

func (a *ServiceAccount) UpdatePreferences(userID string, preferences model.PreferenceSlice) *model_helper.AppError {
	for _, preference := range preferences {
		if userID != preference.UserID {
			return model_helper.NewAppError("savePreferences", "api.preference.update_preferences.set.app_error", nil,
				"userId="+userID+", preference.UserId="+preference.UserID, http.StatusForbidden)
		}
	}

	if err := a.srv.Store.Preference().Save(preferences); err != nil {
		var appErr *model_helper.AppError
		switch {
		case errors.As(err, &appErr):
			return appErr
		default:
			return model_helper.NewAppError("UpdatePreferences", "app.preference.save.updating.app_error", nil, err.Error(), http.StatusBadRequest)
		}
	}

	return nil
}

func (a *ServiceAccount) DeletePreferences(userID string, preferences model.PreferenceSlice) *model_helper.AppError {
	for _, preference := range preferences {
		if userID != preference.UserID {
			err := model_helper.NewAppError(
				"DeletePreferences",
				"api.preference.delete_preferences.delete.app_error",
				nil, "userId="+userID+", preference.UserID="+preference.UserID,
				http.StatusForbidden,
			)
			return err
		}
	}

	for _, preference := range preferences {
		if err := a.srv.Store.Preference().Delete(userID, preference.Category, preference.Name); err != nil {
			return model_helper.NewAppError("DeletePreferences", "app.preference.delete.app_error", nil, err.Error(), http.StatusBadRequest)
		}
	}

	return nil
}
