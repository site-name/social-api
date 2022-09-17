package plugin

import (
	"net/http"

	"github.com/sitename/sitename/model"
)

func (a *ServicePlugin) GetPluginStatus(id string) (*model.PluginStatus, *model.AppError) {
	pluginsEnvironment, appErr := a.GetPluginsEnvironment()
	if appErr != nil {
		return nil, appErr
	}

	pluginStatuses, err := pluginsEnvironment.Statuses()
	if err != nil {
		return nil, model.NewAppError("GetPluginStatus", "app.plugin.get_statuses.app_error", nil, err.Error(), http.StatusInternalServerError)
	}

	for _, status := range pluginStatuses {
		if status.PluginId == id {
			// Add our cluster ID
			if a.srv.Cluster != nil {
				status.ClusterId = a.srv.Cluster.GetClusterId()
			}

			return status, nil
		}
	}

	return nil, model.NewAppError("GetPluginStatus", "app.plugin.not_installed.app_error", nil, "", http.StatusNotFound)
}

// GetPluginStatuses returns the status for plugins installed on this server.
func (a *ServicePlugin) GetPluginStatuses() (model.PluginStatuses, *model.AppError) {
	pluginsEnvironment, appErr := a.GetPluginsEnvironment()
	if appErr != nil {
		return nil, appErr
	}

	pluginStatuses, err := pluginsEnvironment.Statuses()
	if err != nil {
		return nil, model.NewAppError("GetPluginStatuses", "app.plugin.get_statuses.app_error", nil, err.Error(), http.StatusInternalServerError)
	}

	// Add our cluster ID
	for _, status := range pluginStatuses {
		if a.srv.Cluster != nil {
			status.ClusterId = a.srv.Cluster.GetClusterId()
		} else {
			status.ClusterId = ""
		}
	}

	return pluginStatuses, nil
}

// GetClusterPluginStatuses returns the status for plugins installed anywhere in the cluster.
func (a *ServicePlugin) GetClusterPluginStatuses() (model.PluginStatuses, *model.AppError) {
	pluginStatuses, err := a.GetPluginStatuses()
	if err != nil {
		return nil, err
	}

	if a.srv.Cluster != nil && *a.srv.Config().ClusterSettings.Enable {
		clusterPluginStatuses, err := a.srv.Cluster.GetPluginStatuses()
		if err != nil {
			return nil, model.NewAppError("GetClusterPluginStatuses", "app.plugin.get_cluster_plugin_statuses.app_error", nil, err.Error(), http.StatusInternalServerError)
		}

		pluginStatuses = append(pluginStatuses, clusterPluginStatuses...)
	}

	return pluginStatuses, nil
}

func (a *ServicePlugin) notifyPluginStatusesChanged() error {
	pluginStatuses, err := a.GetClusterPluginStatuses()
	if err != nil {
		return err
	}

	// Notify any system admins.
	message := model.NewWebSocketEvent(model.WebsocketEventPluginStatusesChanged, "", nil)
	message.Add("plugin_statuses", pluginStatuses)
	message.GetBroadcast().ContainsSensitiveData = true
	a.srv.Publish(message)

	return nil
}
