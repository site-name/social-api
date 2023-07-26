// this plugin config is borrowed from saleor
package plugin

import (
	"github.com/Masterminds/squirrel"
	"github.com/pkg/errors"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/store"
	"gorm.io/gorm"
)

type SqlPluginConfigurationStore struct {
	store.Store
}

func NewSqlPluginConfigurationStore(s store.Store) store.PluginConfigurationStore {
	return &SqlPluginConfigurationStore{s}
}

// Upsert inserts or updates given plugin configuration and returns it
func (p *SqlPluginConfigurationStore) Upsert(config *model.PluginConfiguration) (*model.PluginConfiguration, error) {
	var err error

	if config.Id == "" {
		err = p.GetMaster().Create(config).Error
	} else {
		err = p.GetMaster().Model(config).Updates(config).Error
	}

	if err != nil {
		if p.IsUniqueConstraintError(err, []string{"Identifier", "ChannelID", "pluginconfigurations_identifier_channelid_key"}) {
			return nil, store.NewErrInvalidInput(model.PluginConfigurationTableName, "Identifier/ChannelID", "duplicate")
		}
		return nil, errors.Wrapf(err, "failed to upsert plugin configuration with id=%s", config.Id)
	}

	return config, nil
}

// Get finds a plugin configuration with given id then returns it
func (p *SqlPluginConfigurationStore) Get(id string) (*model.PluginConfiguration, error) {
	var res model.PluginConfiguration
	err := p.GetReplica().First(&res, "Id = ?", id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, store.NewErrNotFound(model.PluginConfigurationTableName, id)
		}
		return nil, errors.Wrapf(err, "failed to find plugon configuration with id=%s", id)
	}

	return &res, nil
}

// FilterPluginConfigurations finds and returns a list of configs with given options then returns them
func (p *SqlPluginConfigurationStore) FilterPluginConfigurations(options model.PluginConfigurationFilterOptions) ([]*model.PluginConfiguration, error) {
	var configs model.PluginConfigurations
	err := p.GetReplica().Find(&configs, store.BuildSqlizer(options.Conditions)...).Error
	if err != nil {
		return nil, errors.Wrap(err, "failed to find plugin configurations with given options")
	}

	// check if we need to prefetch
	if options.PrefetchRelatedChannel && len(configs) != 0 {
		channels, err := p.Channel().FilterByOption(&model.ChannelFilterOption{
			Conditions: squirrel.Eq{model.PluginConfigurationTableName + ".Id": configs.ChannelIDs()},
		})

		if err != nil {
			return nil, errors.Wrap(err, "failed to find related channels of plugin configs")
		}

		for _, channel := range channels {
			for _, config := range configs {
				if channel.Id == config.ChannelID {
					config.SetRelatedChannel(channel)
				}
			}
		}
	}

	return configs, nil
}

// GetByOptions finds and returns 1 plugin configuration with given options
func (p *SqlPluginConfigurationStore) GetByOptions(options *model.PluginConfigurationFilterOptions) (*model.PluginConfiguration, error) {
	var res model.PluginConfiguration
	err := p.GetReplica().First(&res, store.BuildSqlizer(options.Conditions)...).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, store.NewErrNotFound(model.PluginConfigurationTableName, "options")
		}
		return nil, errors.Wrap(err, "failed to find plugin configuration with given options")
	}

	// check if we need to prefetch
	if options.PrefetchRelatedChannel && model.IsValidId(res.Id) {
		channel, err := p.Channel().Get(res.ChannelID)

		if err != nil {
			return nil, errors.Wrap(err, "failed to find related channels of plugin configs")
		}

		res.SetRelatedChannel(channel)
	}

	return &res, nil
}
