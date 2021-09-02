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
		return nil, store.AppErrorFromDatabaseLookupError("GetPreferencesForUser", "app.account.get_all.app_error", err)
	}
	return preferences, nil
}

func (a *ServiceAccount) GetPreferenceByCategoryForUser(userID string, category string) (model.Preferences, *model.AppError) {
	preferences, err := a.srv.Store.Preference().GetCategory(userID, category)
	if err != nil {
		return nil, store.AppErrorFromDatabaseLookupError("GetPreferenceByCategoryForUser", "app.account.get_category.app_error", err)
	}
	if len(preferences) == 0 {
		err := model.NewAppError("GetPreferenceByCategoryForUser", "api.account.preferences_category.get.app_error", nil, "", http.StatusNotFound)
		return nil, err
	}
	return preferences, nil
}

func (a *ServiceAccount) GetPreferenceByCategoryAndNameForUser(userID string, category string, preferenceName string) (*model.Preference, *model.AppError) {
	res, err := a.srv.Store.Preference().Get(userID, category, preferenceName)
	if err != nil {
		return nil, store.AppErrorFromDatabaseLookupError("GetPreferenceByCategoryAndNameForUser", "app.account.get_preference.app_error", err)
	}
	return res, nil
}

func (a *ServiceAccount) UpdatePreferences(userID string, preferences model.Preferences) *model.AppError {
	for _, preference := range preferences {
		if userID != preference.UserId {
			return model.NewAppError("savePreferences", "api.preference.update_preferences.set.app_error", nil,
				"userId="+userID+", preference.UserId="+preference.UserId, http.StatusForbidden)
		}
	}

	if err := a.srv.Store.Preference().Save(&preferences); err != nil {
		var appErr *model.AppError
		switch {
		case errors.As(err, &appErr):
			return appErr
		default:
			return model.NewAppError("UpdatePreferences", "app.preference.save.updating.app_error", nil, err.Error(), http.StatusBadRequest)
		}
	}

	// if err := a.srv.Store.Channel().UpdateSidebarChannelsByPreferences(&preferences); err != nil {
	// 	return model.NewAppError("UpdatePreferences", "api.preference.update_preferences.update_sidebar.app_error", nil, err.Error(), http.StatusInternalServerError)
	// }

	// message := model.NewWebSocketEvent(model.WEBSOCKET_EVENT_SIDEBAR_CATEGORY_UPDATED, "", "", userID, nil)
	// // TODO this needs to be updated to include information on which categories changed
	// a.Publish(message)

	// message = model.NewWebSocketEvent(model.WEBSOCKET_EVENT_PREFERENCES_CHANGED, "", "", userID, nil)
	// message.Add("preferences", preferences.ToJson())
	// a.Publish(message)

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

	// if err := a.srv.Store.Channel().DeleteSidebarChannelsByPreferences(&preferences); err != nil {
	// 	return model.NewAppError("DeletePreferences", "api.preference.delete_preferences.update_sidebar.app_error", nil, err.Error(), http.StatusInternalServerError)
	// }

	// message := model.NewWebSocketEvent(model.WEBSOCKET_EVENT_SIDEBAR_CATEGORY_UPDATED, "", "", userID, nil)
	// // TODO this needs to be updated to include information on which categories changed
	// a.Publish(message)

	// message = model.NewWebSocketEvent(model.WEBSOCKET_EVENT_PREFERENCES_DELETED, "", "", userID, nil)
	// message.Add("preferences", preferences.ToJson())
	// a.Publish(message)

	return nil
}
