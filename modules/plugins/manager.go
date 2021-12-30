package plugins

import (
	"github.com/Masterminds/squirrel"
	"github.com/sitename/sitename/app"
	"github.com/sitename/sitename/model/plugins"
)

type PluginManager struct {
	srv        *app.Server
	AllPlugins []BasePluginInterface // keys are channel id
}

func NewPluginManager(srv *app.Server, channelID string) *PluginManager {
	m := &PluginManager{
		srv: srv,
	}

	configs, appErr := m.srv.PluginService().FilterPluginConfigurations(&plugins.PluginConfigurationFilterOptions{
		ChannelID:              squirrel.Eq{m.srv.Store.PluginConfiguration().TableName("ChannelID"): channelID},
		PrefetchRelatedChannel: true, // note this.
	})
}

// getAllPluginConfigs returns:
//
// map with keys are channel ids, map values have keys are plugin identifier
// func (m *PluginManager) getAllPluginConfigs(channelID string) (map[string]map[string]*plugins.PluginConfiguration, *model.AppError) {
// 	configs, appErr := m.srv.PluginService().FilterPluginConfigurations(&plugins.PluginConfigurationFilterOptions{
// 		ChannelID:              squirrel.Eq{m.srv.Store.PluginConfiguration().TableName("ChannelID"): channelID},
// 		PrefetchRelatedChannel: true, // note this.
// 	})

// 	if appErr != nil {
// 		return nil, appErr
// 	}

// 	var (
// 		configsPerChannel = make(map[string]map[string]*plugins.PluginConfiguration)
// 	)

// 	for _, config := range configs {
// 		configsPerChannel[config.GetRelatedChannel().Id] = make(map[string]*plugins.PluginConfiguration)

// 	}

// 	return configsPerChannel, nil
// }
