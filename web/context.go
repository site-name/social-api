package web

import (
	"net/http"
	"path"
	"strings"

	"github.com/sitename/sitename/app"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/modules/i18n"
	"github.com/sitename/sitename/modules/slog"
	"github.com/sitename/sitename/modules/util"
)

type Context struct {
	App           app.AppIface
	Logger        *slog.Logger
	Params        *Params
	Err           *model.AppError
	siteURLHeader string
}

func (c *Context) SessionRequired() {
	if !*c.App.Config().ServiceSettings.EnableUserAccessTokens &&
		c.App.Session().Props[model.SESSION_PROP_TYPE] == model.SESSION_TYPE_USER_ACCESS_TOKEN &&
		c.App.Session().Props[model.SESSION_PROP_IS_BOT] != model.SESSION_PROP_IS_BOT_VALUE {

		c.Err = model.NewAppError("", "api.context.session_expired.app_error", nil, "UserAccessToken", http.StatusUnauthorized)
		return
	}

	if c.App.Session().UserId == "" {
		c.Err = model.NewAppError("", "api.context.session_expired.app_error", nil, "UserRequired", http.StatusUnauthorized)
		return
	}
}

func (c *Context) SetSiteURLHeader(url string) {
	c.siteURLHeader = strings.TrimRight(url, "/")
}

func (c *Context) RemoveSessionCookie(w http.ResponseWriter, r *http.Request) {
	subpath, _ := util.GetSubpathFromConfig(c.App.Config())

	cookie := &http.Cookie{
		Name:     model.SESSION_COOKIE_TOKEN,
		Value:    "",
		Path:     subpath,
		MaxAge:   -1,
		HttpOnly: true,
	}

	http.SetCookie(w, cookie)
}

func (c *Context) MfaRequired() {
	// Must be licensed for MFA and have it configured for enforcement
	if !*c.App.Config().ServiceSettings.EnableMultifactorAuthentication || !*c.App.Config().ServiceSettings.EnforceMultifactorAuthentication {
		return
	}

	// OAuth integrations are excepted
	if c.App.Session().IsOAuth {
		return
	}

	user, err := c.App.GetUser(c.App.Session().UserId)
	if err != nil {
		c.Err = model.NewAppError("MfaRequired", "api.context.get_user.app_error", nil, err.Error(), http.StatusUnauthorized)
		return
	}

	if user.IsGuest() && !*c.App.Config().GuestAccountsSettings.EnforceMultifactorAuthentication {
		return
	}
	// Only required for email and ldap accounts
	if user.AuthService != "" &&
		user.AuthService != model.USER_AUTH_SERVICE_EMAIL &&
		user.AuthService != model.USER_AUTH_SERVICE_LDAP {
		return
	}

	// Special case to let user get themself
	subpath, _ := util.GetSubpathFromConfig(c.App.Config())
	if c.App.Path() == path.Join(subpath, "/api/v4/users/me") {
		return
	}

	// Bots are exempt
	if user.IsBot {
		return
	}

	if !user.MfaActive {
		c.Err = model.NewAppError("MfaRequired", "api.context.mfa_required.app_error", nil, "", http.StatusForbidden)
		return
	}
}

func (c *Context) SetServerBusyError() {
	c.Err = NewServerBusyError()
}

func NewServerBusyError() *model.AppError {
	err := model.NewAppError("Context", "api.context.server_busy.app_error", nil, "", http.StatusServiceUnavailable)
	return err
}

func (c *Context) LogErrorByCode(err *model.AppError) {
	code := err.StatusCode
	msg := err.SystemMessage(i18n.TDefault)
	fields := []slog.Field{
		slog.String("err_where", err.Where),
		slog.Int("http_code", err.StatusCode),
		slog.String("err_details", err.DetailedError),
	}
	switch {
	case (code >= http.StatusBadRequest && code < http.StatusInternalServerError) ||
		err.Id == "web.check_browser_compatibility.app_error":
		c.Logger.Debug(msg, fields...)
	case code == http.StatusNotImplemented:
		c.Logger.Info(msg, fields...)
	default:
		c.Logger.Error(msg, fields...)
	}
}
