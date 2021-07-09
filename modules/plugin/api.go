package plugin

import (
	"github.com/sitename/sitename/model"
)

type API interface {
	LoadPluginConfiguration(dest interface{}) error  // LoadPluginConfiguration loads the plugin's configuration. dest should be a pointer to a struct that the configuration JSON can be unmarshalled to.
	GetConfig() *model.Config                        // GetConfig fetches the currently persisted config
	SaveConfig(config *model.Config) *model.AppError // SaveConfig sets the given config and persists the changes
	GetPluginConfig() map[string]interface{}         // GetPluginConfig fetches the currently persisted config of plugin
}
