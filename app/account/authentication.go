package account

import (
	"net/http"
	"strings"

	"github.com/sitename/sitename/app/request"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model_helper"
	"github.com/sitename/sitename/modules/mfa"
)

// CheckPasswordAndAllCriteria
func (a *ServiceAccount) CheckPasswordAndAllCriteria(user model.User, password string, mfaToken string) *model_helper.AppError {
	if err := a.CheckUserPreflightAuthenticationCriteria(user, mfaToken); err != nil {
		return err
	}

	if err := a.CheckUserPassword(user, password); err != nil {
		if passErr := a.srv.Store.User().UpdateFailedPasswordAttempts(user.ID, user.FailedAttempts+1); passErr != nil {
			return model_helper.NewAppError("CheckPasswordAndAllCriteria", "app.user.update_failed_pwd_attempts.app_error", nil, passErr.Error(), http.StatusInternalServerError)
		}

		a.InvalidateCacheForUser(user.ID)

		return err
	}

	if err := a.CheckUserMfa(user, mfaToken); err != nil {
		// If the mfaToken is not set, we assume the client used this as a pre-flight request to query the server
		// about the MFA state of the user in question
		if mfaToken != "" {
			if passErr := a.srv.Store.User().UpdateFailedPasswordAttempts(user.ID, user.FailedAttempts+1); passErr != nil {
				return model_helper.NewAppError("CheckPasswordAndAllCriteria", "app.user.update_failed_pwd_attempts.app_error", nil, passErr.Error(), http.StatusInternalServerError)
			}
		}

		a.InvalidateCacheForUser(user.ID)

		return err
	}

	if passErr := a.srv.Store.User().UpdateFailedPasswordAttempts(user.ID, 0); passErr != nil {
		return model_helper.NewAppError("CheckPasswordAndAllCriteria", "app.user.update_failed_pwd_attempts.app_error", nil, passErr.Error(), http.StatusInternalServerError)
	}

	a.InvalidateCacheForUser(user.ID)

	if err := a.CheckUserPostflightAuthenticationCriteria(user); err != nil {
		return err
	}

	return nil
}

// DoubleCheckPassword performs:
//
// 1) check if number of failed login is not exceed the limit. If yes returns an error
//
// 2) check if user's password and given password don't match, update number of attempts failed in database, return an error
//
// otherwise: set number of failed attempts to 0
func (a *ServiceAccount) DoubleCheckPassword(user model.User, password string) *model_helper.AppError {
	if err := checkUserLoginAttempts(user, *a.srv.Config().ServiceSettings.MaximumLoginAttempts); err != nil {
		return err
	}

	if err := a.CheckUserPassword(user, password); err != nil {
		if passErr := a.srv.Store.User().UpdateFailedPasswordAttempts(user.ID, user.FailedAttempts+1); passErr != nil {
			return model_helper.NewAppError("DoubleCheckPassword", "app.user.update_failed_pwd_attempts.app_error", nil, passErr.Error(), http.StatusInternalServerError)
		}

		a.InvalidateCacheForUser(user.ID)

		return err
	}

	if passErr := a.srv.Store.User().UpdateFailedPasswordAttempts(user.ID, 0); passErr != nil {
		return model_helper.NewAppError("DoubleCheckPassword", "app.user.update_failed_pwd_attempts.app_error", nil, passErr.Error(), http.StatusInternalServerError)
	}

	a.InvalidateCacheForUser(user.ID)

	return nil
}

// CheckUserPassword compares user's password to given password. If they dont match, return an error
func (a *ServiceAccount) CheckUserPassword(user model.User, password string) *model_helper.AppError {
	if err := ComparePassword(user.Password, password); err != nil {
		return model_helper.NewAppError("CheckUserPassword", "api.user.check_user_password.invalid.app_error", nil, "user_id="+user.ID, http.StatusUnauthorized)
	}

	return nil
}

// checkLdapUserPasswordAndAllCriteria
func (a *ServiceAccount) checkLdapUserPasswordAndAllCriteria(ldapId *string, password string, mfaToken string) (*model.User, *model_helper.AppError) {
	if a.srv.Ldap == nil || ldapId == nil {
		err := model_helper.NewAppError("doLdapAuthentication", "api.user.login_ldap.not_available.app_error", nil, "", http.StatusNotImplemented)
		return nil, err
	}

	ldapUser, err := a.srv.Ldap.DoLogin(*ldapId, password)
	if err != nil {
		err.StatusCode = http.StatusUnauthorized
		return nil, err
	}

	if err := a.CheckUserMfa(*ldapUser, mfaToken); err != nil {
		return nil, err
	}

	if err := checkUserNotDisabled(*ldapUser); err != nil {
		return nil, err
	}

	// user successfully authenticated
	return ldapUser, nil
}

func (a *ServiceAccount) CheckUserAllAuthenticationCriteria(user model.User, mfaToken string) *model_helper.AppError {
	if err := a.CheckUserPreflightAuthenticationCriteria(user, mfaToken); err != nil {
		return err
	}

	if err := a.CheckUserPostflightAuthenticationCriteria(user); err != nil {
		return err
	}

	return nil
}

// CheckUserPreflightAuthenticationCriteria checks:
//
// 1) user is not disabled
//
// 2) numbers of failed logins is not exceed the limit
func (a *ServiceAccount) CheckUserPreflightAuthenticationCriteria(user model.User, mfaToken string) *model_helper.AppError {
	if err := checkUserNotDisabled(user); err != nil {
		return err
	}

	if err := checkUserLoginAttempts(user, *a.srv.Config().ServiceSettings.MaximumLoginAttempts); err != nil {
		return err
	}

	return nil
}

// checkUserLoginAttempts checks if user's FailedAttempts >= max, then returns error
func checkUserLoginAttempts(user model.User, max int) *model_helper.AppError {
	if user.FailedAttempts >= max {
		return model_helper.NewAppError("checkUserLoginAttempts", "api.user.check_user_login_attempts.too_many.app_error", nil, "user_id="+user.ID, http.StatusUnauthorized)
	}

	return nil
}

// checkUserNotDisabled checks if user's DeleteAt > 0, then returns error
func checkUserNotDisabled(user model.User) *model_helper.AppError {
	if user.DeleteAt > 0 {
		return model_helper.NewAppError("Login", "api.user.login.inactive.app_error", nil, "user_id="+user.ID, http.StatusUnauthorized)
	}
	return nil
}

// CheckUserPostflightAuthenticationCriteria checks if:
//
// Given user's `EmailVerified` attribute is false && email verification is required,
// Then it return an error.
func (a *ServiceAccount) CheckUserPostflightAuthenticationCriteria(user model.User) *model_helper.AppError {
	if !user.EmailVerified && *a.srv.Config().EmailSettings.RequireEmailVerification {
		return model_helper.NewAppError("Login", "api.user.login.not_verified.app_error", nil, "user_id="+user.ID, http.StatusUnauthorized)
	}

	return nil
}

// CheckUserMfa checks
//
// 1) if given user's `MfaActive` is false || multi factor authentication is not enabled => return nil
//
// 2) multi factor authentication is not enabled => return non-nil error
//
// 3) validates user's `MfaSecret` and given token, if error occur or not valid => return concret error
func (a *ServiceAccount) CheckUserMfa(user model.User, token string) *model_helper.AppError {
	if !user.MfaActive || !*a.srv.Config().ServiceSettings.EnableMultifactorAuthentication {
		return nil
	}

	if !*a.srv.Config().ServiceSettings.EnableMultifactorAuthentication {
		return model_helper.NewAppError("CheckUserMfa", "mfa.mfa_disabled.app_error", nil, "", http.StatusNotImplemented)
	}

	ok, err := mfa.New(a.srv.Store.User()).ValidateToken(user.MfaSecret, token)
	if err != nil {
		return model_helper.NewAppError("CheckUserMfa", "mfa.validate_token.authenticate.app_error", nil, err.Error(), http.StatusBadRequest)
	}

	if !ok {
		return model_helper.NewAppError("checkUserMfa", "api.user.check_user_mfa.bad_code.app_error", nil, "", http.StatusUnauthorized)
	}

	return nil
}

func (a *ServiceAccount) authenticateUser(c *request.Context, user model.User, password, mfaToken string) (*model.User, *model_helper.AppError) {
	ldapAvailable := *a.srv.Config().LdapSettings.Enable && a.srv.Ldap != nil

	if model_helper.UserIsLDAP(user) {
		if !ldapAvailable {
			return nil, model_helper.NewAppError("login", "api.user.login_ldap.not_available.app_error", nil, "", http.StatusNotImplemented)
		}

		ldapUser, err := a.checkLdapUserPasswordAndAllCriteria(user.AuthData.String, password, mfaToken)
		if err != nil {
			err.StatusCode = http.StatusUnauthorized
			return nil, err
		}

		// slightly redundant to get the user again, but we need to get it from the LDAP server
		return ldapUser, nil
	}

	if model_helper.UserIsLDAP(user) {
		return nil, model_helper.NewAppError(
			"login",
			"api.user.login.use_auth_service.app_error",
			map[string]any{
				"AuthService": strings.ToUpper(user.AuthService),
			},
			"", http.StatusBadRequest,
		)
	}

	if err := a.CheckPasswordAndAllCriteria(user, password, mfaToken); err != nil {
		err.StatusCode = http.StatusUnauthorized
		return nil, err
	}

	return &user, nil
}
