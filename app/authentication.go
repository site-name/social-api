package app

import (
	"net/http"
	"strings"

	"github.com/sitename/sitename/app/request"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model/account"
	"github.com/sitename/sitename/modules/mfa"
)

type TokenLocation int

const (
	TokenLocationNotFound TokenLocation = iota
	TokenLocationHeader
	TokenLocationCookie
	TokenLocationQueryString
	TokenLocationCloudHeader
	TokenLocationRemoteClusterHeader
)

func (tl TokenLocation) String() string {
	switch tl {
	case TokenLocationNotFound:
		return "Not Found"
	case TokenLocationHeader:
		return "Header"
	case TokenLocationCookie:
		return "Cookie"
	case TokenLocationQueryString:
		return "QueryString"
	case TokenLocationCloudHeader:
		return "CloudHeader"
	case TokenLocationRemoteClusterHeader:
		return "RemoteClusterHeader"
	default:
		return "Unknown"
	}
}

func (a *App) IsPasswordValid(password string) *model.AppError {

	if *a.Config().ServiceSettings.EnableDeveloper {
		return nil
	}

	return IsPasswordValidWithSettings(password, &a.Config().PasswordSettings)
}

func (a *App) CheckPasswordAndAllCriteria(user *account.User, password string, mfaToken string) *model.AppError {
	if err := a.CheckUserPreflightAuthenticationCriteria(user, mfaToken); err != nil {
		return err
	}

	if err := a.checkUserPassword(user, password); err != nil {
		if passErr := a.Srv().Store.User().UpdateFailedPasswordAttempts(user.Id, user.FailedAttempts+1); passErr != nil {
			return model.NewAppError("CheckPasswordAndAllCriteria", "app.user.update_failed_pwd_attempts.app_error", nil, passErr.Error(), http.StatusInternalServerError)
		}

		a.InvalidateCacheForUser(user.Id)

		return err
	}

	if err := a.CheckUserMfa(user, mfaToken); err != nil {
		// If the mfaToken is not set, we assume the client used this as a pre-flight request to query the server
		// about the MFA state of the user in question
		if mfaToken != "" {
			if passErr := a.Srv().Store.User().UpdateFailedPasswordAttempts(user.Id, user.FailedAttempts+1); passErr != nil {
				return model.NewAppError("CheckPasswordAndAllCriteria", "app.user.update_failed_pwd_attempts.app_error", nil, passErr.Error(), http.StatusInternalServerError)
			}
		}

		a.InvalidateCacheForUser(user.Id)

		return err
	}

	if passErr := a.Srv().Store.User().UpdateFailedPasswordAttempts(user.Id, 0); passErr != nil {
		return model.NewAppError("CheckPasswordAndAllCriteria", "app.user.update_failed_pwd_attempts.app_error", nil, passErr.Error(), http.StatusInternalServerError)
	}

	a.InvalidateCacheForUser(user.Id)

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
func (a *App) DoubleCheckPassword(user *account.User, password string) *model.AppError {
	if err := checkUserLoginAttempts(user, *a.Config().ServiceSettings.MaximumLoginAttempts); err != nil {
		return err
	}

	if err := a.checkUserPassword(user, password); err != nil {
		if passErr := a.Srv().Store.User().UpdateFailedPasswordAttempts(user.Id, user.FailedAttempts+1); passErr != nil {
			return model.NewAppError("DoubleCheckPassword", "app.user.update_failed_pwd_attempts.app_error", nil, passErr.Error(), http.StatusInternalServerError)
		}

		a.InvalidateCacheForUser(user.Id)

		return err
	}

	if passErr := a.Srv().Store.User().UpdateFailedPasswordAttempts(user.Id, 0); passErr != nil {
		return model.NewAppError("DoubleCheckPassword", "app.user.update_failed_pwd_attempts.app_error", nil, passErr.Error(), http.StatusInternalServerError)
	}

	a.InvalidateCacheForUser(user.Id)

	return nil
}

// checkUserPassword compares user's password to given password. If they dont match, return an error
func (a *App) checkUserPassword(user *account.User, password string) *model.AppError {
	if !account.ComparePassword(user.Password, password) {
		return model.NewAppError("checkUserPassword", "api.user.check_user_password.invalid.app_error", nil, "user_id="+user.Id, http.StatusUnauthorized)
	}

	return nil
}

func (a *App) checkLdapUserPasswordAndAllCriteria(ldapId *string, password string, mfaToken string) (*account.User, *model.AppError) {
	if a.Ldap() == nil || ldapId == nil {
		err := model.NewAppError("doLdapAuthentication", "api.user.login_ldap.not_available.app_error", nil, "", http.StatusNotImplemented)
		return nil, err
	}

	ldapUser, err := a.Ldap().DoLogin(*ldapId, password)
	if err != nil {
		err.StatusCode = http.StatusUnauthorized
		return nil, err
	}

	if err := a.CheckUserMfa(ldapUser, mfaToken); err != nil {
		return nil, err
	}

	if err := checkUserNotDisabled(ldapUser); err != nil {
		return nil, err
	}

	// user successfully authenticated
	return ldapUser, nil
}

func (a *App) CheckUserAllAuthenticationCriteria(user *account.User, mfaToken string) *model.AppError {
	if err := a.CheckUserPreflightAuthenticationCriteria(user, mfaToken); err != nil {
		return err
	}

	if err := a.CheckUserPostflightAuthenticationCriteria(user); err != nil {
		return err
	}

	return nil
}

func (a *App) CheckUserPreflightAuthenticationCriteria(user *account.User, mfaToken string) *model.AppError {
	if err := checkUserNotDisabled(user); err != nil {
		return err
	}

	// if err := checkUserNotBot(user); err != nil {
	// 	return err
	// }

	if err := checkUserLoginAttempts(user, *a.Config().ServiceSettings.MaximumLoginAttempts); err != nil {
		return err
	}

	return nil
}

func (a *App) CheckUserPostflightAuthenticationCriteria(user *account.User) *model.AppError {
	if !user.EmailVerified && *a.Config().EmailSettings.RequireEmailVerification {
		return model.NewAppError("Login", "api.user.login.not_verified.app_error", nil, "user_id="+user.Id, http.StatusUnauthorized)
	}

	return nil
}

func (a *App) CheckUserMfa(user *account.User, token string) *model.AppError {
	if !user.MfaActive || !*a.Config().ServiceSettings.EnableMultifactorAuthentication {
		return nil
	}

	if !*a.Config().ServiceSettings.EnableMultifactorAuthentication {
		return model.NewAppError("CheckUserMfa", "mfa.mfa_disabled.app_error", nil, "", http.StatusNotImplemented)
	}

	ok, err := mfa.New(a.Srv().Store.User()).ValidateToken(user.MfaSecret, token)
	if err != nil {
		return model.NewAppError("CheckUserMfa", "mfa.validate_token.authenticate.app_error", nil, err.Error(), http.StatusBadRequest)
	}

	if !ok {
		return model.NewAppError("checkUserMfa", "api.user.check_user_mfa.bad_code.app_error", nil, "", http.StatusUnauthorized)
	}

	return nil
}

func checkUserLoginAttempts(user *account.User, max int) *model.AppError {
	if user.FailedAttempts >= max {
		return model.NewAppError("checkUserLoginAttempts", "api.user.check_user_login_attempts.too_many.app_error", nil, "user_id="+user.Id, http.StatusUnauthorized)
	}

	return nil
}

func checkUserNotDisabled(user *account.User) *model.AppError {
	if user.DeleteAt > 0 {
		return model.NewAppError("Login", "api.user.login.inactive.app_error", nil, "user_id="+user.Id, http.StatusUnauthorized)
	}
	return nil
}

// func checkUserNotBot(user *model.User) *model.AppError {
// 	if user.IsBot {
// 		return model.NewAppError("Login", "api.user.login.bot_login_forbidden.app_error", nil, "user_id="+user.Id, http.StatusUnauthorized)
// 	}
// 	return nil
// }

func (a *App) authenticateUser(c *request.Context, user *account.User, password, mfaToken string) (*account.User, *model.AppError) {
	ldapAvailable := *a.Config().LdapSettings.Enable && a.Ldap() != nil

	if user.AuthService == model.USER_AUTH_SERVICE_LDAP {
		if !ldapAvailable {
			err := model.NewAppError("login", "api.user.login_ldap.not_available.app_error", nil, "", http.StatusNotImplemented)
			return user, err
		}

		ldapUser, err := a.checkLdapUserPasswordAndAllCriteria(user.AuthData, password, mfaToken)
		if err != nil {
			err.StatusCode = http.StatusUnauthorized
			return user, err
		}

		// slightly redundant to get the user again, but we need to get it from the LDAP server
		return ldapUser, nil
	}

	if user.AuthService != "" {
		authService := user.AuthService
		if authService == model.USER_AUTH_SERVICE_SAML {
			authService = strings.ToUpper(authService)
		}
		err := model.NewAppError("login", "api.user.login.use_auth_service.app_error", map[string]interface{}{"AuthService": authService}, "", http.StatusBadRequest)
		return user, err
	}

	if err := a.CheckPasswordAndAllCriteria(user, password, mfaToken); err != nil {
		err.StatusCode = http.StatusUnauthorized
		return user, err
	}

	return user, nil
}

func ParseAuthTokenFromRequest(r *http.Request) (string, TokenLocation) {
	authHeader := r.Header.Get(model.HEADER_AUTH)

	// Attempt to parse the token from the cookie
	if cookie, err := r.Cookie(model.SESSION_COOKIE_TOKEN); err == nil {
		return cookie.Value, TokenLocationCookie
	}

	// Parse the token from the header
	if len(authHeader) > 6 && strings.ToUpper(authHeader[0:6]) == model.HEADER_BEARER {
		// Default session token
		return authHeader[7:], TokenLocationHeader
	}

	if len(authHeader) > 5 && strings.ToLower(authHeader[0:5]) == model.HEADER_TOKEN {
		// OAuth token
		return authHeader[6:], TokenLocationHeader
	}

	// Attempt to parse token out of the query string
	if token := r.URL.Query().Get("access_token"); token != "" {
		return token, TokenLocationQueryString
	}

	if token := r.Header.Get(model.HEADER_CLOUD_TOKEN); token != "" {
		return token, TokenLocationCloudHeader
	}

	if token := r.Header.Get(model.HEADER_REMOTECLUSTER_TOKEN); token != "" {
		return token, TokenLocationRemoteClusterHeader
	}

	return "", TokenLocationNotFound
}

func IsPasswordValidWithSettings(password string, settings *model.PasswordSettings) *model.AppError {
	id := "model.user.is_valid.pwd"
	isError := false

	if len(password) < *settings.MinimumLength || len(password) > model.PASSWORD_MAXIMUM_LENGTH {
		isError = true
	}

	if *settings.Lowercase {
		if !strings.ContainsAny(password, model.LOWERCASE_LETTERS) {
			isError = true
		}

		id = id + "_lowercase"
	}

	if *settings.Uppercase {
		if !strings.ContainsAny(password, model.UPPERCASE_LETTERS) {
			isError = true
		}

		id = id + "_uppercase"
	}

	if *settings.Number {
		if !strings.ContainsAny(password, model.NUMBERS) {
			isError = true
		}

		id = id + "_number"
	}

	if *settings.Symbol {
		if !strings.ContainsAny(password, model.SYMBOLS) {
			isError = true
		}

		id = id + "_symbol"
	}

	if isError {
		return model.NewAppError("User.IsValid", id+".app_error", map[string]interface{}{"Min": *settings.MinimumLength}, "", http.StatusBadRequest)
	}

	return nil
}
