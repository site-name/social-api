package web

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"reflect"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	spanlog "github.com/opentracing/opentracing-go/log"
	"github.com/sitename/sitename/app"
	app_opentracing "github.com/sitename/sitename/app/opentracing"
	"github.com/sitename/sitename/app/request"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/modules/i18n"
	"github.com/sitename/sitename/modules/slog"
	"github.com/sitename/sitename/modules/util"
	"github.com/sitename/sitename/modules/util/api"
	"github.com/sitename/sitename/services/tracing"
	"github.com/sitename/sitename/store/opentracinglayer"
)

// GetHandlerName returns name of the given argument
func GetHandlerName(h func(*Context, http.ResponseWriter, *http.Request)) string {
	handlerName := runtime.FuncForPC(reflect.ValueOf(h).Pointer()).Name()
	pos := strings.LastIndex(handlerName, ".")
	if pos != -1 && len(handlerName) > pos {
		handlerName = handlerName[pos+1:]
	}

	return handlerName
}

// public handler used for testing or static routes
func (w *Web) NewHandler(h func(*Context, http.ResponseWriter, *http.Request)) http.Handler {
	return &Handler{
		Srv:            w.srv,
		HandleFunc:     h,
		RequireSession: false,
		TrustRequester: false,
		RequireMfa:     false,
		IsStatic:       false,
		IsLocal:        false,
		HandlerName:    GetHandlerName(h),
	}
}

func (w *Web) NewStaticHandler(h func(*Context, http.ResponseWriter, *http.Request)) http.Handler {
	// Determine the CSP SHA directive needed for subpath support, if any. This value is fixed
	// on server start and intentionally requires a restart to take effect.
	subpath, _ := model.GetSubpathFromConfig(w.srv.Config())

	return &Handler{
		Srv:             w.srv,
		HandleFunc:      h,
		HandlerName:     GetHandlerName(h),
		RequireSession:  false,
		TrustRequester:  false,
		RequireMfa:      false,
		IsStatic:        true,
		cspShaDirective: model.GetSubpathScriptHash(subpath),
	}
}

type Handler struct {
	Srv                       *app.Server
	HandleFunc                func(*Context, http.ResponseWriter, *http.Request)
	HandlerName               string
	RequireSession            bool
	RequireCloudKey           bool
	RequireRemoteClusterToken bool
	TrustRequester            bool
	RequireMfa                bool
	IsStatic                  bool
	IsLocal                   bool
	DisableWhenBusy           bool
	cspShaDirective           string
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w = newWrappedWriter(w)

	var (
		now        = time.Now()
		requestID  = model.NewId()
		statusCode string
	)

	defer func() {
		responseLogFields := []slog.Field{
			slog.String("method", r.Method),
			slog.String("url", r.URL.Path),
			slog.String("request_id", requestID),
		}
		// Websockets are returning status code 0 to requests after closing the socket
		if statusCode != "0" {
			responseLogFields = append(responseLogFields, slog.String("status_code", statusCode))
		}
		slog.Debug("Received HTTP request", responseLogFields...)
	}()

	c := &Context{
		AppContext: new(request.Context),
		App:        app.New(app.ServerConnector(h.Srv)),
	}

	t, _ := i18n.GetTranslationsAndLocaleFromRequest(r)
	c.AppContext.SetT(t)
	c.AppContext.SetRequestId(requestID)
	c.AppContext.SetIpAddress(util.GetIPAddress(r, c.App.Config().ServiceSettings.TrustedProxyIPHeader))
	c.AppContext.SetUserAgent(r.UserAgent())
	c.AppContext.SetAcceptLanguage(r.Header.Get("Accept-Language"))
	c.AppContext.SetPath(r.URL.Path)
	c.Logger = c.App.Log()

	// check if open tracing is enabled
	if *c.App.Config().ServiceSettings.EnableOpenTracing {
		span, ctx := tracing.StartRootSpanByContext(context.Background(), "web:ServeHTTP")
		carrier := opentracing.HTTPHeadersCarrier(r.Header)
		_ = opentracing.GlobalTracer().Inject(span.Context(), opentracing.HTTPHeaders, carrier)
		ext.HTTPMethod.Set(span, r.Method)
		ext.HTTPUrl.Set(span, c.AppContext.Path())
		ext.PeerAddress.Set(span, c.AppContext.IpAddress())
		span.SetTag("request_id", c.AppContext.RequestId())
		span.SetTag("user_agent", c.AppContext.UserAgent())

		defer func() {
			if c.Err != nil {
				span.LogFields(spanlog.Error(c.Err))
				ext.HTTPStatusCode.Set(span, uint16(c.Err.StatusCode))
				ext.Error.Set(span, true)
			}
			span.Finish()
		}()

		c.AppContext.SetContext(ctx)

		tmpSrv := *c.App.Srv()
		tmpSrv.Store = opentracinglayer.New(c.App.Srv().Store, ctx)
		c.App.SetServer(&tmpSrv)
		c.App = app_opentracing.NewOpenTracingAppLayer(c.App, ctx)
	}

	// Set the max request body size to be equal to MaxFileSize.
	// Ideally, non-file request bodies should be smaller than file request bodies,
	// but we don't have a clean way to identify all file upload handlers.
	// So to keep it simple, we clamp it to the max file size.
	// We add a buffer of bytes.MinRead so that file sizes close to max file size
	// do not get cut off.
	r.Body = http.MaxBytesReader(w, r.Body, *c.App.Config().FileSettings.MaxFileSize+bytes.MinRead)

	subpath, _ := model.GetSubpathFromConfig(c.App.Config())
	c.SetSiteURLHeader(app.GetProtocol(r) + "://" + r.Host + subpath)

	w.Header().Set(model.HeaderRequestId, c.AppContext.RequestId())
	w.Header().Set(model.HeaderVersionId, fmt.Sprintf("%v.%v.%v", model.CurrentVersion, model.BuildNumber, c.App.ClientConfigHash()))

	if *c.App.Config().ServiceSettings.TLSStrictTransport {
		w.Header().Set("Strict-Transport-Security", fmt.Sprintf("max-age=%d", *c.App.Config().ServiceSettings.TLSStrictTransportMaxAge))
	}

	if h.IsStatic {
		// Instruct the browser not to display us in an iframe unless is the same origin for anti-clickjacking
		w.Header().Set("X-Frame-Options", "SAMEORIGIN")

		// Add unsafe-eval to the content security policy for faster source maps in development mode
		devCSP := ""
		if model.BuildNumber == "dev" {
			devCSP += " 'unsafe-eval'"
		}

		// Add unsafe-inline to unlock extensions like React & Redux DevTools in Firefox
		// see https://github.com/reduxjs/redux-devtools/issues/380
		if model.BuildNumber == "dev" {
			devCSP += " 'unsafe-inline'"
		}

		// Set content security policy. This is also specified in the root.html of the webapp in a meta tag.
		w.Header().Set("Content-Security-Policy", fmt.Sprintf(
			"frame-ancestors 'self'; script-src 'self' cdn.rudderlabs.com%s%s%s",
			"",
			h.cspShaDirective,
			devCSP,
		))
	} else {
		// All api response bodies will be JSON formatted by default
		w.Header().Set("Content-Type", "application/json")

		if r.Method == "GET" {
			w.Header().Set("Expires", "0")
		}
	}

	token, tokenLocation := app.ParseAuthTokenFromRequest(r)

	if token != "" {
		session, err := c.App.Srv().AccountService().GetSession(token)
		defer c.App.Srv().AccountService().ReturnSessionToPool(session)

		if err != nil {
			c.Logger.Info("Invalid session", slog.Err(err))
			if err.StatusCode == http.StatusInternalServerError {
				c.Err = err
			} else if h.RequireSession {
				c.RemoveSessionCookie(w, r)
				c.Err = model.NewAppError("ServeHTTP", "api.context.session_expired.app_error", nil, "token="+token, http.StatusUnauthorized)
			}
		} else {
			c.AppContext.SetSession(session)
		}

		// Rate limit by UserID
		if c.App.Srv().RateLimiter != nil && c.App.Srv().RateLimiter.UserIdRateLimit(c.AppContext.Session().UserId, w) {
			return
		}

		h.checkCSRFToken(c, r, token, tokenLocation, session)
	}

	c.Logger = c.App.Log().With(
		slog.String("path", c.AppContext.Path()),
		slog.String("request_id", c.AppContext.RequestId()),
		slog.String("ip_addr", c.AppContext.IpAddress()),
		slog.String("user_id", c.AppContext.Session().UserId),
		slog.String("method", r.Method),
	)

	if c.Err == nil && h.RequireSession {
		c.SessionRequired()
	}

	if c.Err == nil && h.RequireMfa {
		c.MfaRequired()
	}

	if c.Err == nil && h.DisableWhenBusy && c.App.Srv().Busy.IsBusy() {
		c.SetServerBusyError()
	}

	if c.Err == nil {
		h.HandleFunc(c, w, r)
	}

	if c.Err != nil {
		c.Err.Translate(c.AppContext.T)
		c.Err.RequestId = c.AppContext.RequestId()
		c.LogErrorByCode(c.Err)

		c.Err.Where = r.URL.Path

		// Block out detailed error when not in developer mode
		if !*c.App.Config().ServiceSettings.EnableDeveloper {
			c.Err.DetailedError = ""
		}

		// Sanitize all 5xx error messages in hardened mode
		if *c.App.Config().ServiceSettings.ExperimentalEnableHardenedMode && c.Err.StatusCode >= http.StatusInternalServerError {
			c.Err.Id = ""
			c.Err.Message = "Internal Server Error"
			c.Err.DetailedError = ""
			c.Err.StatusCode = http.StatusInternalServerError
			c.Err.Where = ""
			c.Err.IsOAuth = false
		}

		if IsApiCall(c.App, r) ||
			// IsWebhookCall(c.App, r) ||
			// IsOAuthApiCall(c.App, r) ||
			r.Header.Get("X-Mobile-App") != "" {
			w.WriteHeader(c.Err.StatusCode)
			w.Write([]byte(c.Err.ToJSON()))
		} else {
			api.RenderWebAppError(c.App.Config(), w, r, c.Err, c.App.AsymmetricSigningKey())
		}

		if c.App.Metrics() != nil {
			c.App.Metrics().IncrementHttpError()
		}
	}

	statusCode = strconv.Itoa(w.(*responseWriterWrapper).StatusCode())
	if c.App.Metrics() != nil {
		c.App.Metrics().IncrementHttpRequest()

		if r.URL.Path != /*model.API_URL_SUFFIX+*/ "/websocket" {
			elapsed := float64(time.Since(now)) / float64(time.Second)
			c.App.Metrics().ObserveApiEndpointDuration(h.HandlerName, r.Method, statusCode, elapsed)
		}
	}
}

// checkCSRFToken performs a CSRF check on the provided request with the given CSRF token. Returns whether or not
// a CSRF check occurred and whether or not it succeeded.
func (h *Handler) checkCSRFToken(c *Context, r *http.Request, token string, tokenLocation app.TokenLocation, session *model.Session) (checked bool, passed bool) {
	csrfCheckNeeded := session != nil && c.Err == nil && tokenLocation == app.TokenLocationCookie && !h.TrustRequester && r.Method != "GET"
	csrfCheckPassed := false

	if csrfCheckNeeded {
		csrfHeader := r.Header.Get(model.HEADER_CSRF_TOKEN)

		if csrfHeader == session.GetCSRF() {
			csrfCheckPassed = true
		} else if r.Header.Get(model.HEADER_REQUESTED_WITH) == model.HEADER_REQUESTED_WITH_XML {
			// ToDo(DSchalla) 2019/01/04: Remove after deprecation period and only allow CSRF Header (MM-13657)
			csrfErrorMessage := "CSRF Header check failed for request - Please upgrade your web application or custom app to set a CSRF Header"

			sid := ""
			userId := ""

			if session != nil {
				sid = session.Id
				userId = session.UserId
			}

			fields := []slog.Field{
				slog.String("path", r.URL.Path),
				slog.String("ip", r.RemoteAddr),
				slog.String("session_id", sid),
				slog.String("user_id", userId),
			}

			if *c.App.Config().ServiceSettings.ExperimentalStrictCSRFEnforcement {
				c.Logger.Warn(csrfErrorMessage, fields...)
			} else {
				c.Logger.Debug(csrfErrorMessage, fields...)
				csrfCheckPassed = true
			}
		}

		if !csrfCheckPassed {
			c.AppContext.SetSession(&model.Session{})
			c.Err = model.NewAppError("ServeHTTP", "api.context.session_expired.app_error", nil, "token="+token+" Appears to be a CSRF attempt", http.StatusUnauthorized)
		}
	}

	return csrfCheckNeeded, csrfCheckPassed
}

// ApiSessionRequired provides a handler for API endpoints which require the user to be logged in in order for access to
// be granted.
// func (w *Web) ApiSessionRequired(h func(*Context, http.ResponseWriter, *http.Request)) http.Handler {
// 	handler := &Handler{
// 		Srv:            w.srv,
// 		HandleFunc:     h,
// 		HandlerName:    GetHandlerName(h),
// 		RequireSession: true,
// 		TrustRequester: false,
// 		RequireMfa:     true,
// 		IsStatic:       false,
// 		IsLocal:        false,
// 	}
// 	if *w.srv.Config().ServiceSettings.WebserverMode == "gzip" {
// 		return gziphandler.GzipHandler(handler)
// 	}
// 	return handler
// }
