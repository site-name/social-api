package app

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"path"
	"path/filepath"
	"strings"

	"github.com/gorilla/mux"
	"github.com/sitename/sitename/model_helper"
	"github.com/sitename/sitename/modules/plugin"
	"github.com/sitename/sitename/modules/slog"
	"github.com/sitename/sitename/modules/util"
)

func (s *Server) ServePluginRequest(w http.ResponseWriter, r *http.Request) {
	pluginsEnvironment, appErr := s.Plugin.GetPluginsEnvironment()
	if appErr != nil {
		s.Log.Error(appErr.Error())
		w.WriteHeader(appErr.StatusCode)
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(appErr.ToJSON()))
		return
	}

	params := mux.Vars(r)
	hooks, err := pluginsEnvironment.HooksForPlugin(params["plugin_id"])
	if err != nil {
		s.Log.Error("Access to route for non-existent plugin",
			slog.String("missing_plugin_id", params["plugin_id"]),
			slog.String("url", r.URL.String()),
			slog.Err(err))
		http.NotFound(w, r)
		return
	}

	s.servePluginRequest(w, r, hooks.ServeHTTP)
}

// ServePluginPublicRequest serves public plugin files
// at the URL http(s)://$SITE_URL/plugins/$PLUGIN_ID/public/{anything}
func (s *Server) ServePluginPublicRequest(w http.ResponseWriter, r *http.Request) {
	if strings.HasSuffix(r.URL.Path, "/") {
		http.NotFound(w, r)
		return
	}

	// Should be in the form of /$PLUGIN_ID/public/{anything} by the time we get here
	vars := mux.Vars(r)
	pluginID := vars["plugin_id"]

	pluginsEnv, _ := s.Plugin.GetPluginsEnvironment()

	// Check if someone has nullified the pluginsEnv in the meantime
	if pluginsEnv == nil {
		http.NotFound(w, r)
		return
	}

	publicFilesPath, err := pluginsEnv.PublicFilesPath(pluginID)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	publicFilePath := path.Clean(r.URL.Path)
	prefix := fmt.Sprintf("/plugins/%s/public/", pluginID)
	if !strings.HasPrefix(publicFilePath, prefix) {
		http.NotFound(w, r)
		return
	}
	publicFile := filepath.Join(publicFilesPath, strings.TrimPrefix(publicFilePath, prefix))
	http.ServeFile(w, r, publicFile)
}

func (s *Server) servePluginRequest(w http.ResponseWriter, r *http.Request, handler func(*plugin.Context, http.ResponseWriter, *http.Request)) {
	token := ""
	context := &plugin.Context{
		RequestId:      model_helper.NewId(),
		IpAddress:      util.GetIPAddress(r, s.Config().ServiceSettings.TrustedProxyIPHeader),
		AcceptLanguage: r.Header.Get("Accept-Language"),
		UserAgent:      r.UserAgent(),
	}
	cookieAuth := false

	authHeader := r.Header.Get(model_helper.HeaderAuth)
	if strings.HasPrefix(strings.ToUpper(authHeader), model_helper.HeaderBearer+" ") {
		token = authHeader[len(model_helper.HeaderBearer)+1:]
	} else if strings.HasPrefix(strings.ToLower(authHeader), model_helper.HeaderToken+" ") {
		token = authHeader[len(model_helper.HeaderToken)+1:]
	} else if cookie, _ := r.Cookie(model_helper.SESSION_COOKIE_TOKEN); cookie != nil {
		token = cookie.Value
		cookieAuth = true
	} else {
		token = r.URL.Query().Get("access_token")
	}

	// Mattermost-Plugin-ID can only be set by inter-plugin requests
	r.Header.Del("Mattermost-Plugin-ID")

	r.Header.Del("Mattermost-User-Id")
	if token != "" {
		session, err := s.Account.GetSession(token)
		defer s.Account.ReturnSessionToPool(session)

		csrfCheckPassed := false

		if session != nil && err == nil && cookieAuth && r.Method != "GET" {
			sentToken := ""

			if r.Header.Get(model_helper.HeaderCsrfToken) == "" {
				bodyBytes, _ := io.ReadAll(r.Body)
				r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
				r.ParseForm()
				sentToken = r.FormValue(model_helper.SESSION_CSRF_KEY)
				r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
			} else {
				sentToken = r.Header.Get(model_helper.HeaderCsrfToken)
			}

			expectedToken := model_helper.SessionGetCSRF(session)

			if sentToken == expectedToken {
				csrfCheckPassed = true
			}

			// ToDo(DSchalla) 2019/01/04: Remove after deprecation period and only allow CSRF Header (MM-13657)
			if r.Header.Get(model_helper.HeaderRequestedWith) == model_helper.HeaderRequestedWith_XML && !csrfCheckPassed {
				csrfErrorMessage := "CSRF Check failed for request - Please migrate your plugin to either send a CSRF Header or Form Field, XMLHttpRequest is deprecated"
				sid := ""
				userID := ""

				if session.ID != "" {
					sid = session.ID
					userID = session.UserID
				}

				fields := []slog.Field{
					slog.String("path", r.URL.Path),
					slog.String("ip", r.RemoteAddr),
					slog.String("session_id", sid),
					slog.String("user_id", userID),
				}

				if *s.Config().ServiceSettings.ExperimentalStrictCSRFEnforcement {
					s.Log.Warn(csrfErrorMessage, fields...)
				} else {
					s.Log.Debug(csrfErrorMessage, fields...)
					csrfCheckPassed = true
				}
			}
		} else {
			csrfCheckPassed = true
		}

		if (session != nil && session.ID != "") && err == nil && csrfCheckPassed {
			r.Header.Set("Mattermost-User-Id", session.UserID)
			context.SessionId = session.ID
		}
	}

	cookies := r.Cookies()
	r.Header.Del("Cookie")
	for _, c := range cookies {
		if c.Name != model_helper.SESSION_COOKIE_TOKEN {
			r.AddCookie(c)
		}
	}
	r.Header.Del(model_helper.HeaderAuth)
	r.Header.Del("Referer")

	params := mux.Vars(r)

	subpath, _ := model_helper.GetSubpathFromConfig(s.Config())

	newQuery := r.URL.Query()
	newQuery.Del("access_token")
	r.URL.RawQuery = newQuery.Encode()
	r.URL.Path = strings.TrimPrefix(r.URL.Path, path.Join(subpath, "plugins", params["plugin_id"]))

	handler(context, w, r)
}
