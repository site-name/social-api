package plugin

import (
	"net/http"

	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model_helper"
	"github.com/sitename/sitename/store"
)

// UpsertPluginConfiguration updates/inserts given configuration into database then returns it
func (s *ServicePlugin) UpsertPluginConfiguration(config *model.PluginConfiguration) (*model.PluginConfiguration, *model_helper.AppError) {
	config, err := s.srv.Store.PluginConfiguration().Upsert(config)
	if err != nil {
		if appErr, ok := err.(*model_helper.AppError); ok {
			return nil, appErr
		}
		var statusCode int = http.StatusInternalServerError
		if _, ok := err.(*store.ErrInvalidInput); ok {
			statusCode = http.StatusBadRequest
		}

		return nil, model_helper.NewAppError("UpsertPluginConfiguration", "app.plugin.error_upsertig_plugin_configuration.app_error", nil, err.Error(), statusCode)
	}

	return config, nil
}

// FilterPluginConfigurations returns a list of plugin configurations filtered using given options
func (s *ServicePlugin) FilterPluginConfigurations(options *model.PluginConfigurationFilterOptions) (model.PluginConfigurations, *model_helper.AppError) {
	if options == nil {
		return nil, model_helper.NewAppError("FilterPluginConfigurations", model_helper.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "options"}, "", http.StatusBadRequest)
	}

	configs, err := s.srv.Store.PluginConfiguration().FilterPluginConfigurations(*options)
	if err != nil {
		return nil, model_helper.NewAppError("FilterPluginConfigurations", "app.plugins.error_finding_plugins.app_error", nil, err.Error(), http.StatusInternalServerError)
	}

	return configs, nil
}

// GetPluginConfiguration finds and returns a plugin configuration based on given options
func (s *ServicePlugin) GetPluginConfiguration(options *model.PluginConfigurationFilterOptions) (*model.PluginConfiguration, *model_helper.AppError) {
	config, err := s.srv.Store.PluginConfiguration().GetByOptions(options)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if _, ok := err.(*store.ErrNotFound); ok {
			statusCode = http.StatusInternalServerError
		}

		return nil, model_helper.NewAppError("GetPluginConfiguration", "app.plugin.error_finding_plugin_configuration_by_options.app_error", nil, err.Error(), statusCode)
	}

	return config, nil
}
