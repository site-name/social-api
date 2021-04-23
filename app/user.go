package app

import (
	"net/http"

	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/modules/slog"
)

const (
	TokenTypePasswordRecovery  = "password_recovery"
	TokenTypeVerifyEmail       = "verify_email"
	TokenTypeTeamInvitation    = "team_invitation"
	TokenTypeGuestInvitation   = "guest_invitation"
	TokenTypeCWSAccess         = "cws_access_token"
	PasswordRecoverExpiryTime  = 1000 * 60 * 60      // 1 hour
	InvitationExpiryTime       = 1000 * 60 * 60 * 48 // 48 hours
	ImageProfilePixelDimension = 128
)

// func (a *App) CreateUserWithToken(user *model.User, token *model.Token) (*model.User, *model.AppError) {
// 	if err := a.IsUserSignupAllowed(); err != nil {
// 		return nil, err
// 	}

// 	if model.GetMillis()-token.CreateAt >= InvitationExpiryTime {
// 		a.DeleteToken(token)
// 		return nil, model.NewAppError("CreateUserWithToken", "api.user.create_user.signup_link_expired.app_error", nil, "", http.StatusBadRequest)
// 	}

// 	tokenData := model.MapFromJson(strings.NewReader(token.Extra))

// 	user.Email = tokenData["email"]
// 	user.EmailVerified = true

// 	var ruser *model.User
// 	var err *model.AppError

// }

func (a *App) CreateUserFromSignup(user *model.User, redirect string) (*model.User, *model.AppError) {
	if err := a.IsUserSignupAllowed(); err != nil {
		return nil, err
	}

	if !a.IsFirstUserAccount() {
		err := model.NewAppError("CreateUserFromSignup", "api.user.create_user.no_open_server", nil, "email="+user.Email, http.StatusForbidden)
		return nil, err
	}

	user.EmailVerified = false

	ruser, err := a.CreateUser(user)
	if err != nil {
		return nil, err
	}

	if err := a.Srv().EmailService.sendWelcomeEmail(ruser.Id, ruser.Email, ruser.EmailVerified, ruser.DisableWelcomeEmail, ruser.Locale, a.GetSiteURL(), redirect); err != nil {
		slog.Warn("Failed to send welcome email on create user from signup", slog.Err(err))
	}

	return ruser, nil
}

// CreateUser creates a user and sets several fields of the returned User struct to
// their zero values.
func (a *App) CreateUser(user *model.User) (*model.User, *model.AppError) {
	// return a.createUserOrGuest(user, false)
}

func (s *Server) IsFirstUserAccount() bool {
	cachedSessions, err := s.sessionCache.Len()
	if err != nil {
		return false
	}
	if cachedSessions == 0 {
		count, err := s.Store.User().Count(model.UserCountOptions{IncludeDeleted: true})
		if err != nil {
			slog.Debug("There was an error fetching if first usder account", slog.Err(err))
			return false
		}
		if count <= 0 {
			return true
		}
	}
	return false
}

func (a *App) IsFirstUserAccount() bool {
	return a.Srv().IsFirstUserAccount()
}

// IsUserSignupAllowed checks email settings if signing up with email is allowed
func (a *App) IsUserSignupAllowed() *model.AppError {
	if !*a.Config().EmailSettings.EnableSignUpWithEmail {
		err := model.NewAppError("IsUserSignupAllowed", "api.user.create_user.signup_email_disabled.app_error", nil, "", http.StatusNotImplemented)
		return err
	}
	return nil
}

func (a *App) DeleteToken(token *model.Token) *model.AppError {
	err := a.Srv().Store.Token().Delete(token.Token)
	if err != nil {
		return model.NewAppError("DeleteToken", "app.recover.delete.app_error", nil, err.Error(), http.StatusInternalServerError)
	}
	return nil
}

// func (a *App) DeactivateGuests() *model.AppError {
// 	userIDs, err := a.Srv().Store.User().DeactivateGuests()
// 	if err != nil {
// 		return model.NewAppError("DeactivateGuests", "app.user.update_active_for_multiple_users.updating.app_error", nil, err.Error(), http.StatusInternalServerError)
// 	}

// 	for _, userID := range userIDs {
// 		if err := a.userDeactivated(userID); err != nil {
// 			return err
// 		}
// 	}

// 	a.Srv().Store.User().ClearCaches()
// }

// func (a *App) userDeactivated(userID string) *model.AppError {
// 	if err := a.Revo
// }
