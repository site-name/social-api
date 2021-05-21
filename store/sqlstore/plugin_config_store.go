package sqlstore

import (
	"github.com/sitename/sitename/model/plugins"
	"github.com/sitename/sitename/store"
)

type SqlPluginConfigurationStore struct {
	*SqlStore
}

func newSqlPluginConfigurationStore(s *SqlStore) store.PluginConfigurationStore {
	ps := &SqlPluginConfigurationStore{s}

	for _, db := range s.GetAllConns() {
		table := db.AddTableWithName(plugins.PluginConfiguration{}, "PluginConfigurations").SetKeys(false, "Id")
		table.ColMap("Id").SetMaxSize(UUID_MAX_LENGTH)
		table.ColMap("ChannelID").SetMaxSize(UUID_MAX_LENGTH)
		table.ColMap("Name").SetMaxSize(plugins.PLUGIN_CONFIGURATION_COMMON_MAX_LENGHT)
		table.ColMap("Identifier").SetMaxSize(plugins.PLUGIN_CONFIGURATION_COMMON_MAX_LENGHT)
		table.ColMap("Description").SetMaxSize(plugins.PLUGIN_CONFIGURATION_DESCRIPTION_MAX_LENGHT)

		table.SetUniqueTogether("Identifier", "ChannelID")
	}
	return ps
}

func (p *SqlPluginConfigurationStore) createIndexesIfNotExists() {
	p.CreateIndexIfNotExists("idx_plugin_configurations_identifier", "PluginConfigurations", "Identifier")
	p.CreateIndexIfNotExists("idx_plugin_configurations_name", "PluginConfigurations", "Name")
	p.CreateIndexIfNotExists("idx_plugin_configurations_lower_textpattern_name", "PluginConfigurations", "lower(Name) text_pattern_ops")
}
