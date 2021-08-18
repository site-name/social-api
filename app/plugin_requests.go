package app

import (
	"fmt"
	"net/http"
	"path"
	"path/filepath"
	"strings"

	"github.com/gorilla/mux"
	"github.com/sitename/sitename/modules/slog"
)

func (s *Server) ServePluginRequest(w http.ResponseWriter, r *http.Request) {
	pluginsEnvironment, appErr := s.GetPluginsEnvironment()
	if appErr != nil {
		s.Log.Error(appErr.Error())
		w.WriteHeader(appErr.StatusCode)
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(appErr.ToJson()))
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

	pluginsEnv, _ := s.GetPluginsEnvironment()

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
