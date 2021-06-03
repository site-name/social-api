package app

import (
	"net/http"
	"strings"
	"time"

	"github.com/sitename/sitename/app/request"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model/account"
	"github.com/sitename/sitename/modules/util"
)

func (a *App) CheckForClientSideCert(r *http.Request) (string, string, string) {
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

func (a *App) AuthenticateUserForLogin(c *request.Context, id, loginId, password, mfaToken, cwsToken string, ldapOnly bool) (user *account.User, err *model.AppError) {
	// Do statistics
	defer func() {
		if a.Metrics() != nil {
			if user == nil || err != nil {
				a.Metrics().IncrementLoginFail()
			} else {
				a.Metrics().IncrementLogin()
			}
		}
	}()

	if password == "" && !IsCWSLogin(a, cwsToken) {
		return nil, model.NewAppError("AuthenticateUserForLogin", "api.user.login.blank_pwd.app_error", nil, "", http.StatusBadRequest)
	}

	// get the sn user we are trying to login
	if user, err = a.GetUserForLogin(id, cwsToken); err != nil {
		return nil, err
	}
}

func (a *App) GetUserForLogin(id, loginId string) (*account.User, *model.AppError) {
	enableUsername := *a.Config().EmailSettings.EnableSignInWithUsername
	enableEmail := *a.Config().EmailSettings.EnableSignInWithEmail

	// if we are given a userID then fail if we can't find a user with that ID
	if id != "" {
		user, err := a.GetUser(id)
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
	if user, err := a.Srv().Store.User().GetForLogin(loginId, enableUsername, enableEmail); err == nil {
		return user, nil
	}

	// Try to get the user with LDAP if enabled
	if *a.Config().LdapSettings.Enable && a.Ldap() != nil {
		if ldapUser, err := a.Ldap().GetUser(loginId); err == nil {
			if user, err := a.GetUserByAuth(ldapUser.AuthData, model.USER_AUTH_SERVICE_LDAP); err == nil {
				return user, nil
			}
			return ldapUser, nil
		}
	}

	return nil, model.NewAppError("GetUserForLogin", "store.sql_user.get_for_login.app_error", nil, "", http.StatusBadRequest)
}

func GetProtocol(r *http.Request) string {
	if r.Header.Get(model.HEADER_FORWARDED_PROTO) == "https" || r.TLS != nil {
		return "https"
	}
	return "http"
}

func (a *App) AttachSessionCookies(c *request.Context, w http.ResponseWriter, r *http.Request) {
	secure := false
	if GetProtocol(r) == "https" {
		secure = true
	}

	maxAge := *a.Config().ServiceSettings.SessionLengthWebInDays * 60 * 60 * 24
	domain := a.GetCookieDomain()
	subpath, _ := util.GetSubpathFromConfig(a.Config())

	expiresAt := time.Unix(model.GetMillis()/1000+int64(maxAge), 0)
	sessionCookie := &http.Cookie{
		Name:     model.SESSION_COOKIE_TOKEN,
		Value:    c.Session().Token,
		Path:     subpath,
		MaxAge:   maxAge,
		Expires:  expiresAt,
		HttpOnly: true,
		Domain:   domain,
		Secure:   secure,
	}

	userCookie := &http.Cookie{
		Name:    model.SESSION_COOKIE_USER,
		Value:   c.Session().UserId,
		Path:    subpath,
		MaxAge:  maxAge,
		Expires: expiresAt,
		Domain:  domain,
		Secure:  secure,
	}

	csrfCookie := &http.Cookie{
		Name:    model.SESSION_COOKIE_CSRF,
		Value:   c.Session().GetCSRF(),
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

func IsCWSLogin(a *App, token string) bool {
	return token != ""
}
