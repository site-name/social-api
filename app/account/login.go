package account

import (
	"context"
	"crypto/subtle"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/avct/uasurfer"
	"github.com/sitename/sitename/app"
	"github.com/sitename/sitename/app/request"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model_helper"
	"github.com/sitename/sitename/modules/model_types"
	"github.com/sitename/sitename/modules/slog"
	"github.com/sitename/sitename/store"
)

const cwsTokenEnv = "CWS_CLOUD_TOKEN"

// CheckForClientSideCert checks request's header's `X-SSL-Client-Cert` and `X-SSL-Client-Cert-Subject-DN` keys
func (a *ServiceAccount) CheckForClientSideCert(r *http.Request) (string, string, string) {
	pem := r.Header.Get("X-SSL-Client-Cert")
	subject := r.Header.Get("X-SSL-Client-Cert-Subject-DN")
	email := ""

	if subject != "" {
		for _, v := range strings.Split(subject, "/") {
			kv := strings.Split(v, "=")
			if len(kv) == 2 && kv[0] == "emailAddress" {
				email = kv[1]
			}
		}
	}

	return pem, subject, email
}

// AuthenticateUserForLogin
func (a *ServiceAccount) AuthenticateUserForLogin(c *request.Context, id, loginId, password, mfaToken, cwsToken string, ldapOnly bool) (user *model.User, err *model_helper.AppError) {
	// Do statistics
	defer func() {
		if a.metrics != nil {
			if user == nil || err != nil {
				a.metrics.IncrementLoginFail()
			} else {
				a.metrics.IncrementLogin()
			}
		}
	}()

	if password == "" && !IsCWSLogin(a, cwsToken) {
		return nil, model_helper.NewAppError("AuthenticateUserForLogin", "api.user.login.blank_pwd.app_error", nil, "", http.StatusBadRequest)
	}

	// get the sn user we are trying to login
	if user, err = a.GetUserForLogin(id, loginId); err != nil {
		return nil, err
	}

	// CWS login allow to use the one-time token to login the users when they're redirected to their
	// installation for the first time
	if IsCWSLogin(a, cwsToken) {
		token, err := a.srv.Store.Token().GetByToken(cwsToken)
		if nfErr := new(store.ErrNotFound); err != nil && !errors.As(err, &nfErr) {
			slog.Debug("Error retrieving the cws token from the store", slog.Err(err))
			return nil, model_helper.NewAppError(
				"AuthenticateUserForLogin",
				"api.user.login_by_cws.invalid_token.app_error",
				nil,
				"",
				http.StatusInternalServerError,
			)
		}
		// If token is stored in the database that means it was used
		if token != nil {
			return nil, model_helper.NewAppError(
				"AuthenticateUserForLogin",
				"api.user.login_by_cws.invalid_token.app_error",
				nil,
				"",
				http.StatusBadRequest,
			)
		}
		envToken, ok := os.LookupEnv(cwsTokenEnv)
		if ok && subtle.ConstantTimeCompare([]byte(envToken), []byte(cwsToken)) == 1 {
			token := model.Token{
				Token:     cwsToken,
				CreatedAt: model_helper.GetMillis(),
				Type:      model_helper.TokenTypeCWSAccess.String(),
			}
			_, err := a.srv.Store.Token().Save(token)
			if err != nil {
				slog.Debug("Error storing the cws token in the store", slog.Err(err))
				return nil, model_helper.NewAppError(
					"AuthenticateUserForLogin",
					"api.user.login_by_cws.invalid_token.app_error",
					nil,
					"",
					http.StatusInternalServerError,
				)
			}
			return user, nil
		}
		return nil, model_helper.NewAppError(
			"AuthenticateUserForLogin",
			"api.user.login_by_cws.invalid_token.app_error",
			nil,
			"",
			http.StatusBadRequest,
		)
	}

	// If client side cert is enable and it's checking as a primary source
	// then trust the proxy and cert that the correct user is supplied and allow
	// them access
	if *a.srv.Config().ExperimentalSettings.ClientSideCertEnable && *a.srv.Config().ExperimentalSettings.ClientSideCertCheck == model_helper.CLIENT_SIDE_CERT_CHECK_PRIMARY_AUTH {
		return user, nil
	}

	// and then authenticate them
	if user, err = a.authenticateUser(c, *user, password, mfaToken); err != nil {
		return nil, err
	}
	return user, nil
}

// GetUserForLogin
func (a *ServiceAccount) GetUserForLogin(id, loginId string) (*model.User, *model_helper.AppError) {
	enableUsername := *a.srv.Config().EmailSettings.EnableSignInWithUsername
	enableEmail := *a.srv.Config().EmailSettings.EnableSignInWithEmail

	// if we are given a userID then fail if we can't find a user with that ID
	if id != "" {
		user, err := a.UserById(context.Background(), id)
		if err != nil {
			if err.Id != MissingAccountError {
				err.StatusCode = http.StatusInternalServerError
				return nil, err
			}
			err.StatusCode = http.StatusBadRequest
			return nil, err
		}
		return user, nil
	}

	// Try to get the user by username/email
	user, err := a.srv.Store.User().GetForLogin(loginId, enableUsername, enableEmail)
	if err == nil {
		return user, nil
	}

	// Try to get the user with LDAP if enabled
	if *a.srv.Config().LdapSettings.Enable && a.srv.Ldap != nil {
		if ldapUser, err := a.srv.Ldap.GetUser(loginId); err == nil {
			if !ldapUser.AuthData.IsNil() && *ldapUser.AuthData.String != "" {
				if user, err := a.GetUserByOptions(
					model.UserWhere.AuthData.EQ(ldapUser.AuthData),
					model.UserWhere.AuthService.EQ(model_helper.USER_AUTH_SERVICE_LDAP),
				); err == nil {
					return user, nil
				}
			}

			return ldapUser, nil
		}
	}

	return nil, model_helper.NewAppError("GetUserForLogin", "store.sql_user.get_for_login.app_error", nil, "", http.StatusBadRequest)
}

func (a *ServiceAccount) DoLogin(c *request.Context, w http.ResponseWriter, r *http.Request, user model.User, deviceID string, isMobile, isOAuthUser, isSaml bool) *model_helper.AppError {
	session := model.Session{
		UserID:   user.ID,
		Roles:    model_helper.UserGetRawRoles(user),
		DeviceID: deviceID,
		IsOauth:  false,
		Props: model_types.JsonMap{
			model_helper.USER_AUTH_SERVICE_IS_MOBILE: strconv.FormatBool(isMobile),
			model_helper.USER_AUTH_SERVICE_IS_SAML:   strconv.FormatBool(isSaml),
			model_helper.USER_AUTH_SERVICE_IS_OAUTH:  strconv.FormatBool(isOAuthUser),
		},
	}
	model_helper.SessionGenerateCSRF(&session)

	if deviceID != "" {
		a.SetSessionExpireInDays(&session, *a.srv.Config().ServiceSettings.SessionLengthMobileInDays)

		// A special case where we log out of all other sessions with the same Id
		if err := a.RevokeSessionsForDeviceId(user.ID, deviceID, ""); err != nil {
			err.StatusCode = http.StatusInternalServerError
			return err
		}
	} else if isMobile {
		a.SetSessionExpireInDays(&session, *a.srv.Config().ServiceSettings.SessionLengthMobileInDays)
	} else if isOAuthUser || isSaml {
		a.SetSessionExpireInDays(&session, *a.srv.Config().ServiceSettings.SessionLengthSSOInDays)
	} else {
		a.SetSessionExpireInDays(&session, *a.srv.Config().ServiceSettings.SessionLengthWebInDays)
	}

	ua := uasurfer.Parse(r.UserAgent())

	model_helper.SessionAddProp(&session, model_helper.SESSION_PROP_PLATFORM, app.GetPlatformName(ua))
	model_helper.SessionAddProp(&session, model_helper.SESSION_PROP_OS, app.GetOSName(ua))
	model_helper.SessionAddProp(&session, model_helper.SESSION_PROP_BROWSER, fmt.Sprintf("%s/%s", app.GetBrowserName(ua, r.UserAgent()), app.GetBrowserVersion(ua, r.UserAgent())))

	savedSession, err := a.CreateSession(session)
	if err != nil {
		return err
	}

	w.Header().Set(model_helper.HeaderToken, session.Token)

	c.SetSession(*savedSession)
	if a.srv.Ldap != nil {
		userVal := user
		sessionVal := session

		a.srv.Go(func() {
			a.srv.Ldap.UpdateProfilePictureIfNecessary(userVal, sessionVal)
		})
	}

	return nil
}

// AttachSessionCookies sets:
//
// 1) session cookie with value of given s's session's token to given w
//
// 2) user cookie with value of user id
//
// 3) csrf cookie with value of csrf in session
func (a *ServiceAccount) AttachSessionCookies(c *request.Context, w http.ResponseWriter, r *http.Request) {
	secure := app.GetProtocol(r) == "https"

	maxAge := *a.srv.Config().ServiceSettings.SessionLengthWebInDays * 60 * 60 * 24
	domain := a.srv.GetCookieDomain()
	subpath, _ := model_helper.GetSubpathFromConfig(a.srv.Config())

	expiresAt := time.Unix(model_helper.GetMillis()/1000+int64(maxAge), 0)
	sessionCookie := &http.Cookie{
		Name:     model_helper.SESSION_COOKIE_TOKEN,
		Value:    c.Session().Token,
		Path:     subpath,
		MaxAge:   maxAge,
		Expires:  expiresAt,
		HttpOnly: true,
		Domain:   domain,
		Secure:   secure,
	}

	userCookie := &http.Cookie{
		Name:    model_helper.SESSION_COOKIE_USER,
		Value:   c.Session().UserID,
		Path:    subpath,
		MaxAge:  maxAge,
		Expires: expiresAt,
		Domain:  domain,
		Secure:  secure,
	}

	csrfCookie := &http.Cookie{
		Name:    model_helper.SESSION_COOKIE_CSRF,
		Value:   model_helper.SessionGetCSRF(*c.Session()),
		Path:    subpath,
		MaxAge:  maxAge,
		Expires: expiresAt,
		Domain:  domain,
		Secure:  secure,
	}

	http.SetCookie(w, sessionCookie)
	http.SetCookie(w, userCookie)
	http.SetCookie(w, csrfCookie)
}

// IsCWSLogin returns true if token != "" else false
func IsCWSLogin(a *ServiceAccount, token string) bool {
	return token != ""
}
