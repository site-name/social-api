// this plugin config is borrowed from saleor
package plugin

import (
	"database/sql"

	"github.com/Masterminds/squirrel"
	"github.com/pkg/errors"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model/channel"
	"github.com/sitename/sitename/model/plugins"
	"github.com/sitename/sitename/store"
)

type SqlPluginConfigurationStore struct {
	store.Store
}

func NewSqlPluginConfigurationStore(s store.Store) store.PluginConfigurationStore {
	return &SqlPluginConfigurationStore{s}
}

func (s *SqlPluginConfigurationStore) ModelFields(prefix string) model.StringArray {
	res := model.StringArray{
		"Id",
		"Identifier",
		"Name",
		"ChannelID",
		"Description",
		"Active",
		"Configuration",
	}
	if prefix == "" {
		return res
	}

	return res.Map(func(_ int, s string) string {
		return prefix + s
	})
}

// Upsert inserts or updates given plugin configuration and returns it
func (p *SqlPluginConfigurationStore) Upsert(config *plugins.PluginConfiguration) (*plugins.PluginConfiguration, error) {
	var isSaving bool

	if config.Id == "" {
		isSaving = true
		config.PreSave()
	} else {
		config.PreUpdate()
	}

	if err := config.IsValid(); err != nil {
		return nil, err
	}

	var (
		err        error
		numUpdated int64
	)
	if isSaving {
		query := "INSERT INTO " + store.PluginConfigurationTableName + "(" + p.ModelFields("").Join(",") + ") VALUES (" + p.ModelFields(":").Join(",") + ")"
		_, err = p.GetMasterX().NamedExec(query, config)

	} else {
		query := "UPDATE " + store.PluginConfigurationTableName + " SET " + p.
			ModelFields("").
			Map(func(_ int, s string) string {
				return s + "=:" + s
			}).
			Join(",") + " WHERE Id=:Id"

		var result sql.Result
		result, err = p.GetMasterX().NamedExec(query, config)
		if err == nil && result != nil {
			numUpdated, _ = result.RowsAffected()
		}
	}

	if err != nil {
		if p.IsUniqueConstraintError(err, []string{"Identifier", "ChannelID", "pluginconfigurations_identifier_channelid_key"}) {
			return nil, store.NewErrInvalidInput(store.PluginConfigurationTableName, "Identifier/ChannelID", "duplicate")
		}
		return nil, errors.Wrapf(err, "failed to upsert plugin configuration with id=%s", config.Id)
	}
	if numUpdated != 1 {
		return nil, errors.Errorf("%d configuration(s) were/was updated instewad of 1", numUpdated)
	}

	return config, nil
}

// Get finds a plugin configuration with given id then returns it
func (p *SqlPluginConfigurationStore) Get(id string) (*plugins.PluginConfiguration, error) {
	var res plugins.PluginConfiguration
	err := p.GetReplicaX().Get(&res, "SELECT * FROM "+store.PluginConfigurationTableName+" WHERE Id = ?", id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, store.NewErrNotFound(store.PluginConfigurationTableName, id)
		}
		return nil, errors.Wrapf(err, "failed to find plugon configuration with id=%s", id)
	}

	return &res, nil
}

func (p *SqlPluginConfigurationStore) optionsParse(options *plugins.PluginConfigurationFilterOptions) (string, []interface{}, error) {
	query := p.GetQueryBuilder().
		Select("*").
		From(store.PluginConfigurationTableName)

	// parse options
	if options.Id != nil {
		query = query.Where(options.Id)
	}
	if options.Identifier != nil {
		query = query.Where(options.Identifier)
	}
	if options.ChannelID != nil {
		query = query.Where(options.ChannelID)
	}

	return query.ToSql()
}

// FilterPluginConfigurations finds and returns a list of configs with given options then returns them
func (p *SqlPluginConfigurationStore) FilterPluginConfigurations(options plugins.PluginConfigurationFilterOptions) ([]*plugins.PluginConfiguration, error) {
	queryStr, args, err := p.optionsParse(&options)
	if err != nil {
		return nil, errors.Wrap(err, "FilterPluginConfigurations_ToSql")
	}

	var configs plugins.PluginConfigurations
	err = p.GetReplicaX().Select(&configs, queryStr, args...)
	if err != nil {
		return nil, errors.Wrap(err, "failed to find plugin configurations with given options")
	}

	// check if we need to prefetch
	if options.PrefetchRelatedChannel && len(configs) != 0 {
		channels, err := p.Channel().FilterByOption(&channel.ChannelFilterOption{
			Id: squirrel.Eq{store.PluginConfigurationTableName + ".Id": configs.ChannelIDs()},
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
func (p *SqlPluginConfigurationStore) GetByOptions(options *plugins.PluginConfigurationFilterOptions) (*plugins.PluginConfiguration, error) {
	queryStr, args, err := p.optionsParse(options)
	if err != nil {
		return nil, errors.Wrap(err, "GetByOptions_ToSql")
	}

	var res plugins.PluginConfiguration
	err = p.GetReplicaX().Get(&res, queryStr, args...)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, store.NewErrNotFound(store.PluginConfigurationTableName, "options")
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
