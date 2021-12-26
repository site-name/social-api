// this plugin config is borrowed from saleor
package plugin

import (
	"database/sql"

	"github.com/pkg/errors"
	"github.com/sitename/sitename/model/plugins"
	"github.com/sitename/sitename/store"
)

type SqlPluginConfigurationStore struct {
	store.Store
}

func NewSqlPluginConfigurationStore(s store.Store) store.PluginConfigurationStore {
	ps := &SqlPluginConfigurationStore{s}

	for _, db := range s.GetAllConns() {
		table := db.AddTableWithName(plugins.PluginConfiguration{}, ps.TableName("")).SetKeys(false, "Id")
		table.ColMap("Id").SetMaxSize(store.UUID_MAX_LENGTH)
		table.ColMap("ChannelID").SetMaxSize(store.UUID_MAX_LENGTH)
		table.ColMap("Name").SetMaxSize(plugins.PLUGIN_CONFIGURATION_COMMON_MAX_LENGHT)
		table.ColMap("Identifier").SetMaxSize(plugins.PLUGIN_CONFIGURATION_COMMON_MAX_LENGHT)
		table.ColMap("Description").SetMaxSize(plugins.PLUGIN_CONFIGURATION_DESCRIPTION_MAX_LENGHT)

		table.SetUniqueTogether("Identifier", "ChannelID")
	}
	return ps
}

func (p *SqlPluginConfigurationStore) CreateIndexesIfNotExists() {
	p.CreateIndexIfNotExists("idx_plugin_configurations_identifier", p.TableName(""), "Identifier")
	p.CreateIndexIfNotExists("idx_plugin_configurations_name", p.TableName(""), "Name")
	p.CreateIndexIfNotExists("idx_plugin_configurations_lower_textpattern_name", p.TableName(""), "lower(Name) text_pattern_ops")
}

func (p *SqlPluginConfigurationStore) TableName(withField string) string {
	name := "PluginConfigurations"

	if withField != "" {
		return name + "." + withField
	}

	return name
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
		err = p.GetMaster().Insert(config)
	} else {
		_, err = p.Get(config.Id)
		if err != nil {
			return nil, err
		}

		numUpdated, err = p.GetMaster().Update(config)
	}

	if err != nil {
		if p.IsUniqueConstraintError(err, []string{}) {
			return nil, store.NewErrInvalidInput(p.TableName(""), "Identifier/ChannelID", "duplicate")
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
	err := p.GetReplica().SelectOne(&res, "SELECT * FROM "+p.TableName("")+" WHERE Id = :ID", map[string]interface{}{"ID": id})
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, store.NewErrNotFound(p.TableName(""), id)
		}
		return nil, errors.Wrapf(err, "failed to find plugon configuration with id=%s", id)
	}

	return &res, nil
}
