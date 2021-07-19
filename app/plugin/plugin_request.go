package plugin

import (
	"net/http"

	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/modules/plugin"
	"github.com/sitename/sitename/modules/slog"
)

func (a *AppPlugin) ServeInterPluginRequest(w http.ResponseWriter, r *http.Request, sourcePluginId, destinationPluginId string) {
	pluginEnvironment, appErr := a.GetPluginsEnvironment()
	if appErr != nil {
		a.Log().Error(appErr.Error())
		w.WriteHeader(appErr.StatusCode)
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(appErr.ToJson()))
		return
	}

	hooks, err := pluginEnvironment.HooksForPlugin(destinationPluginId)
	if err != nil {
		a.Log().Error("Access to route for non-existent plugin in inter plugin request",
			slog.String("source_plugin_id", sourcePluginId),
			slog.String("destination_plugin_id", destinationPluginId),
			slog.String("url", r.URL.String()),
			slog.Err(err),
		)
		http.NotFound(w, r)
		return
	}

	context := &plugin.Context{
		RequestId:      model.NewId(),
		UserAgent:      r.UserAgent(),
		SourcePluginId: sourcePluginId,
	}

	r.Header.Set("Sitename-Plugin-ID", sourcePluginId)

	hooks.ServeHTTP(context, w, r)
}
