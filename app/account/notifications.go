package account

import (
	"net/http"
	"net/url"

	"github.com/sitename/sitename/app/plugin/interfaces"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model_helper"
	"github.com/sitename/sitename/modules/util"
)

func (s *ServiceAccount) GetDefaultUserPayload(user *model.User) map[string]any {
	if user == nil {
		return nil
	}
	return map[string]any{
		"id":               user.ID,
		"email":            user.Email,
		"first_name":       user.FirstName,
		"last_name":        user.LastName,
		"is_active":        user.IsActive,
		"private_metadata": user.PrivateMetadata,
		"metadata":         user.Metadata,
		"language_code":    user.Locale,
	}
}

// Trigger sending a password reset notification for the given customer/staff.
func (s *ServiceAccount) SendPasswordResetNotification(redirectURL string, user model.User, manager interfaces.PluginManagerInterface, channelID string) *model_helper.AppError {
	var (
		token         = util.DefaultTokenGenerator.MakeToken(model_helper.UserWrapper{User: user})
		params        = url.Values{"email": []string{user.Email}, "token": []string{token}}
		resetURL, err = util.PrepareUrl(params, redirectURL)
	)
	if err != nil {
		return model_helper.NewAppError("SendPasswordResetNotification", model_helper.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "redirectURL"}, err.Error(), http.StatusBadRequest)
	}

	payload := map[string]any{
		"user":            s.GetDefaultUserPayload(&user),
		"recipient_email": user.Email,
		"token":           token,
		"reset_url":       resetURL,
		"channel_id":      channelID,
		"domain":          *s.srv.Config().ServiceSettings.SiteURL,
		"site_name":       *s.srv.Config().ServiceSettings.SiteName,
	}

	_, appErr := manager.Notify(model_helper.ACCOUNT_PASSWORD_RESET, payload, channelID, "")
	return appErr
}

// Trigger sending an account confirmation notification for the given user
func (s *ServiceAccount) SendAccountConfirmation(redirectUrl string, user model.User, manager interfaces.PluginManagerInterface, channelID string) *model_helper.AppError {
	var (
		token           = util.DefaultTokenGenerator.MakeToken(model_helper.UserWrapper{User: user})
		params          = url.Values{"email": []string{user.Email}, "token": []string{token}}
		confirmUrl, err = util.PrepareUrl(params, redirectUrl)
	)
	if err != nil {
		return model_helper.NewAppError("SendAccountConfirmation", model_helper.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "redirectUrl"}, err.Error(), http.StatusBadRequest)
	}

	payload := map[string]any{
		"user":            s.GetDefaultUserPayload(&user),
		"recipient_email": user.Email,
		"token":           token,
		"confirm_url":     confirmUrl,
		"channel_id":      channelID,
		"domain":          *s.srv.Config().ServiceSettings.SiteURL,
		"site_name":       *s.srv.Config().ServiceSettings.SiteName,
	}

	_, appErr := manager.Notify(model_helper.ACCOUNT_CONFIRMATION, payload, channelID, "")
	return appErr
}

// Trigger sending a notification change email for the given user
func (s *ServiceAccount) SendRequestUserChangeEmailNotification(redirectUrl string, user model.User, newEmail string, token string, manager interfaces.PluginManagerInterface, channelID string) *model_helper.AppError {
	var (
		params                 = url.Values{"token": []string{token}}
		parsedRedirectUrl, err = util.PrepareUrl(params, redirectUrl)
	)
	if err != nil {
		return model_helper.NewAppError("SendRequestUserChangeEmailNotification", model_helper.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "redirectUrl"}, err.Error(), http.StatusBadRequest)
	}

	payload := map[string]any{
		"user":            s.GetDefaultUserPayload(&user),
		"recipient_email": newEmail,
		"old_email":       user.Email,
		"new_email":       newEmail,
		"token":           token,
		"redirect_url":    parsedRedirectUrl,
		"channel_id":      channelID,
		"domain":          *s.srv.Config().ServiceSettings.SiteURL,
		"site_name":       *s.srv.Config().ServiceSettings.SiteName,
	}

	_, appErr := manager.Notify(model_helper.ACCOUNT_CHANGE_EMAIL_REQUEST, payload, channelID, "")
	return appErr
}

// Trigger sending a email change notification for the given user
func (s *ServiceAccount) SendUserChangeEmailNotification(recipientEmail string, user model.User, manager interfaces.PluginManagerInterface, channelID string) *model_helper.AppError {
	payload := map[string]any{
		"user":            s.GetDefaultUserPayload(&user),
		"recipient_email": recipientEmail,
		"channel_id":      channelID,
		"domain":          *s.srv.Config().ServiceSettings.SiteURL,
		"site_name":       *s.srv.Config().ServiceSettings.SiteName,
	}
	_, appErr := manager.Notify(model_helper.ACCOUNT_CHANGE_EMAIL_CONFIRM, payload, channelID, "")
	return appErr
}

// SendAccountDeleteConfirmationNotification Trigger sending a account delete notification for the given user
func (s *ServiceAccount) SendAccountDeleteConfirmationNotification(redirectUrl string, user model.User, manager interfaces.PluginManagerInterface, channelID string) *model_helper.AppError {
	var (
		token          = util.DefaultTokenGenerator.MakeToken(model_helper.UserWrapper{User: user})
		params         = url.Values{"token": []string{token}}
		deleteUrl, err = util.PrepareUrl(params, redirectUrl)
	)
	if err != nil {
		return model_helper.NewAppError("SendAccountDeleteConfirmationNotification", model_helper.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "redirectUrl"}, err.Error(), http.StatusBadRequest)
	}

	payload := map[string]any{
		"user":            s.GetDefaultUserPayload(&user),
		"recipient_email": user.Email,
		"token":           token,
		"delete_url":      deleteUrl,
		"channel_id":      channelID,
		"domain":          *s.srv.Config().ServiceSettings.SiteURL,
		"site_name":       *s.srv.Config().ServiceSettings.SiteName,
	}

	_, appErr := manager.Notify(model_helper.ACCOUNT_DELETE, payload, channelID, "")
	return appErr
}

// Trigger sending a set password notification for the given customer/staff.
func (s *ServiceAccount) SendSetPasswordNotification(redirectUrl string, user model.User, manager interfaces.PluginManagerInterface, channelID string) *model_helper.AppError {
	var (
		token               = util.DefaultTokenGenerator.MakeToken(model_helper.UserWrapper{User: user})
		params              = url.Values{"token": []string{token}, "email": []string{user.Email}}
		passwordSetURL, err = util.PrepareUrl(params, redirectUrl)
	)
	if err != nil {
		return model_helper.NewAppError("SendSetPasswordNotification", model_helper.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "redirectUrl"}, err.Error(), http.StatusBadRequest)
	}

	payload := map[string]any{
		"user":             s.GetDefaultUserPayload(&user),
		"token":            token,
		"recipient_email":  user.Email,
		"password_set_url": passwordSetURL,
		"channel_id":       channelID,
		"domain":           *s.srv.Config().ServiceSettings.SiteURL,
		"site_name":        *s.srv.Config().ServiceSettings.SiteName,
	}

	_, appErr := manager.Notify(model_helper.ACCOUNT_SET_CUSTOMER_PASSWORD, payload, channelID, "")
	return appErr
}
