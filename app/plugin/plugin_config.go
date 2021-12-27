package plugin

import (
	"net/http"

	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model/plugins"
	"github.com/sitename/sitename/store"
)

// UpsertPluginConfiguration updates/inserts given configuration into database then returns it
func (s *ServicePlugin) UpsertPluginConfiguration(config *plugins.PluginConfiguration) (*plugins.PluginConfiguration, *model.AppError) {
	config, err := s.srv.Store.PluginConfiguration().Upsert(config)
	if err != nil {
		if appErr, ok := err.(*model.AppError); ok {
			return nil, appErr
		}
		var statusCode int = http.StatusInternalServerError
		if _, ok := err.(*store.ErrInvalidInput); ok {
			statusCode = http.StatusBadRequest
		}

		return nil, model.NewAppError("UpsertPluginConfiguration", "app.plugin.error_upsertig_plugin_configuration.app_error", nil, err.Error(), statusCode)
	}

	return config, nil
}
