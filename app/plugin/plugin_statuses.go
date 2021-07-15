package plugin

import (
	"net/http"

	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model/plugins"
)

func (a *AppPlugin) GetPluginStatus(id string) (*plugins.PluginStatus, *model.AppError) {
	pluginsEnvironment := a.GetPluginsEnvironment()
	if pluginsEnvironment == nil {
		return nil, model.NewAppError("GetPluginStatus", "app.plugin.disabled.app_error", nil, "", http.StatusNotImplemented)
	}

	pluginStatuses, err := pluginsEnvironment.Statuses()
	if err != nil {
		return nil, model.NewAppError("GetPluginStatus", "app.plugin.get_statuses.app_error", nil, err.Error(), http.StatusInternalServerError)
	}

	for _, status := range pluginStatuses {
		if status.PluginId == id {
			// Add our cluster ID
			if a.Srv().Cluster != nil {
				status.ClusterId = a.Srv().Cluster.GetClusterId()
			}

			return status, nil
		}
	}

	return nil, model.NewAppError("GetPluginStatus", "app.plugin.not_installed.app_error", nil, "", http.StatusNotFound)
}

// GetPluginStatuses returns the status for plugins installed on this server.
func (a *AppPlugin) GetPluginStatuses() (plugins.PluginStatuses, *model.AppError) {
	pluginsEnvironment := a.GetPluginsEnvironment()
	if pluginsEnvironment == nil {
		return nil, model.NewAppError("GetPluginStatuses", "app.plugin.disabled.app_error", nil, "", http.StatusNotImplemented)
	}

	pluginStatuses, err := pluginsEnvironment.Statuses()
	if err != nil {
		return nil, model.NewAppError("GetPluginStatuses", "app.plugin.get_statuses.app_error", nil, err.Error(), http.StatusInternalServerError)
	}

	// Add our cluster ID
	for _, status := range pluginStatuses {
		if a.Srv().Cluster != nil {
			status.ClusterId = a.Srv().Cluster.GetClusterId()
		} else {
			status.ClusterId = ""
		}
	}

	return pluginStatuses, nil
}

// GetClusterPluginStatuses returns the status for plugins installed anywhere in the cluster.
func (a *AppPlugin) GetClusterPluginStatuses() (plugins.PluginStatuses, *model.AppError) {
	pluginStatuses, err := a.GetPluginStatuses()
	if err != nil {
		return nil, err
	}

	if a.Srv().Cluster != nil && *a.Config().ClusterSettings.Enable {
		clusterPluginStatuses, err := a.Srv().Cluster.GetPluginStatuses()
		if err != nil {
			return nil, model.NewAppError("GetClusterPluginStatuses", "app.plugin.get_cluster_plugin_statuses.app_error", nil, err.Error(), http.StatusInternalServerError)
		}

		pluginStatuses = append(pluginStatuses, clusterPluginStatuses...)
	}

	return pluginStatuses, nil
}

func (a *AppPlugin) notifyPluginStatusesChanged() error {
	pluginStatuses, err := a.GetClusterPluginStatuses()
	if err != nil {
		return err
	}

	// Notify any system admins.
	message := model.NewWebSocketEvent(model.WEBSOCKET_EVENT_PLUGIN_STATUSES_CHANGED, "", nil)
	message.Add("plugin_statuses", pluginStatuses)
	message.GetBroadcast().ContainsSensitiveData = true
	a.Srv().Publish(message)

	return nil
}
