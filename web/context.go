package web

import (
	"context"
	"net/http"
	"path"
	"strings"

	"github.com/samber/lo"
	"github.com/sitename/sitename/app"
	"github.com/sitename/sitename/app/request"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/modules/i18n"
	"github.com/sitename/sitename/modules/slog"
)

type Context struct {
	App        app.AppIface
	AppContext *request.Context // AppContext holds information of an http request. It is created when an http request is made
	Logger     *slog.Logger
	Err        *model.AppError
	// This is used to track the graphQL query that's being executed,
	// so that we can monitor the timings in Grafana.
	GraphQLOperationName string
	siteURLHeader        string

	CurrentChannelID string
}

// set session missing error for c
func (c *Context) SessionRequired() {
	if !*c.App.Config().ServiceSettings.EnableUserAccessTokens &&
		c.AppContext.Session().Props[model.SESSION_PROP_TYPE] == model.SESSION_TYPE_USER_ACCESS_TOKEN {

		c.Err = model.NewAppError("", "api.context.session_expired.app_error", nil, "UserAccessToken", http.StatusUnauthorized)
		return
	}

	if c.AppContext.Session().UserId == "" {
		c.Err = model.NewAppError("", "api.context.session_expired.app_error", nil, "UserRequired", http.StatusUnauthorized)
		return
	}
}

// CheckAuthenticatedAndHasPermissionToAll checks if user is authenticated, then check if user has all given permission(s)
func (c *Context) CheckAuthenticatedAndHasPermissionToAll(perms ...*model.Permission) {
	c.SessionRequired()
	if c.Err != nil {
		return
	}
	if !c.App.Srv().AccountService().SessionHasPermissionToAll(c.AppContext.Session(), perms...) {
		c.SetPermissionError(perms...)
	}
}

// CheckAuthenticatedAndHasPermissionToAny check user authenticated, then check if user has any of given permission(s)
func (c *Context) CheckAuthenticatedAndHasPermissionToAny(perms ...*model.Permission) {
	c.SessionRequired()
	if c.Err != nil {
		return
	}
	if !c.App.Srv().AccountService().SessionHasPermissionToAny(c.AppContext.Session(), perms...) {
		c.SetPermissionError(perms...)
	}
}

func (c *Context) CheckAuthenticatedAndHasRoles(apiName string, roleIDs ...string) {
	c.SessionRequired()
	if c.Err != nil {
		return
	}
	commonRoles := lo.Intersect(c.AppContext.Session().GetUserRoles(), roleIDs)
	if len(commonRoles) == 0 {
		c.Err = model.NewAppError(apiName, "api.unauthorized.app_error", nil, "you are not allowed to perform this action", http.StatusUnauthorized)
	}
}

// MfaRequired must be placed after c's SessionRequired() method
func (c *Context) MfaRequired() {
	// OAuth integrations are excepted
	if c.AppContext.Session().IsOAuth {
		return
	}

	user, err := c.App.Srv().AccountService().UserById(context.Background(), c.AppContext.Session().UserId)
	if err != nil {
		c.Err = model.NewAppError("MfaRequired", "api.context.get_user.app_error", nil, err.Error(), http.StatusUnauthorized)
		return
	}

	if !*c.App.Config().GuestAccountsSettings.EnforceMultifactorAuthentication {
		return
	}
	// Only required for email and ldap accounts
	if user.AuthService != "" &&
		user.AuthService != model.USER_AUTH_SERVICE_EMAIL &&
		user.AuthService != model.USER_AUTH_SERVICE_LDAP {
		return
	}

	// Special case to let user get themself
	subpath, _ := model.GetSubpathFromConfig(c.App.Config())
	if c.AppContext.Path() == path.Join(subpath, "/api/v4/users/me") {
		return
	}

	if !user.MfaActive {
		c.Err = model.NewAppError("MfaRequired", "api.context.mfa_required.app_error", nil, "", http.StatusForbidden)
		return
	}
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

// ExtendSessionExpiryIfNeeded will update Session.ExpiresAt based on session lengths in config.
// Session cookies will be resent to the client with updated max age.
func (c *Context) ExtendSessionExpiryIfNeeded(w http.ResponseWriter, r *http.Request) {
	if ok := c.App.Srv().AccountService().ExtendSessionExpiryIfNeeded(c.AppContext.Session()); ok {
		c.App.Srv().AccountService().AttachSessionCookies(c.AppContext, w, r)
	}
}

// RemoveSessionCookie deletes cookie from subpath route
func (c *Context) RemoveSessionCookie(w http.ResponseWriter, r *http.Request) {
	subpath, _ := model.GetSubpathFromConfig(c.App.Config())

	cookie := &http.Cookie{
		Name:     model.SESSION_COOKIE_TOKEN,
		Value:    "",
		Path:     subpath,
		MaxAge:   -1,
		HttpOnly: true,
	}

	http.SetCookie(w, cookie)
}

func (c *Context) SetInvalidParam(parameter string) {
	c.Err = NewInvalidParamError(parameter)
}

func (c *Context) SetInvalidUrlParam(parameter string) {
	c.Err = NewInvalidUrlParamError(parameter)
}

// SetServerBusyError set c's Err property to a non-nil AppError
func (c *Context) SetServerBusyError() {
	c.Err = NewServerBusyError()
}

func (c *Context) SetInvalidRemoteIdError(id string) {
	c.Err = NewInvalidRemoteIdError(id)
}

func (c *Context) SetInvalidRemoteClusterTokenError() {
	c.Err = NewInvalidRemoteClusterTokenError()
}

func (c *Context) SetJSONEncodingError() {
	c.Err = NewJSONEncodingError()
}

func (c *Context) SetCommandNotFoundError() {
	c.Err = model.NewAppError("GetCommand", "store.sql_command.save.get.app_error", nil, "", http.StatusNotFound)
}

func (c *Context) HandleEtag(etag string, routeName string, w http.ResponseWriter, r *http.Request) bool {
	metrics := c.App.Metrics()
	if et := r.Header.Get(model.HEADER_ETAG_CLIENT); etag != "" {
		if et == etag {
			w.Header().Set(model.HEADER_ETAG_SERVER, etag)
			w.WriteHeader(http.StatusNotModified)
			if metrics != nil {
				metrics.IncrementEtagHitCounter(routeName)
			}
			return true
		}
	}

	if metrics != nil {
		metrics.IncrementEtagMissCounter(routeName)
	}

	return false
}

// IsSystemAdmin checks if given session contains info of system's administrator.
func (c *Context) IsSystemAdmin() bool {
	c.SessionRequired()
	return c.Err == nil && c.AppContext.Session().GetUserRoles().Contains(model.SystemAdminRoleId)
}

func NewInvalidParamError(parameter string) *model.AppError {
	err := model.NewAppError("Context", "api.context.invalid_body_param.app_error", map[string]interface{}{"Name": parameter}, "", http.StatusBadRequest)
	return err
}
func NewInvalidUrlParamError(parameter string) *model.AppError {
	err := model.NewAppError("Context", "api.context.invalid_url_param.app_error", map[string]interface{}{"Name": parameter}, "", http.StatusBadRequest)
	return err
}
func NewServerBusyError() *model.AppError {
	err := model.NewAppError("Context", "api.context.server_busy.app_error", nil, "", http.StatusServiceUnavailable)
	return err
}

func NewInvalidRemoteIdError(parameter string) *model.AppError {
	err := model.NewAppError("Context", "api.context.remote_id_invalid.app_error", map[string]interface{}{"RemoteId": parameter}, "", http.StatusBadRequest)
	return err
}

func NewInvalidRemoteClusterTokenError() *model.AppError {
	err := model.NewAppError("Context", "api.context.remote_id_invalid.app_error", nil, "", http.StatusUnauthorized)
	return err
}

func NewJSONEncodingError() *model.AppError {
	err := model.NewAppError("Context", "api.context.json_encoding.app_error", nil, "", http.StatusInternalServerError)
	return err
}

func (c *Context) SetPermissionError(permissions ...*model.Permission) {
	c.Err = c.App.Srv().AccountService().MakePermissionError(c.AppContext.Session(), permissions...)
}

func (c *Context) SetSiteURLHeader(url string) {
	c.siteURLHeader = strings.TrimRight(url, "/")
}

func (c *Context) GetSiteURLHeader() string {
	return c.siteURLHeader
}

func (c *Context) GetRemoteID(r *http.Request) string {
	return r.Header.Get(model.HEADER_REMOTECLUSTER_ID)
}
