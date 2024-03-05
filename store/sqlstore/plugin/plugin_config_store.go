// this plugin config is borrowed from saleor
package plugin

import (
	"database/sql"

	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model_helper"
	"github.com/sitename/sitename/store"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

type SqlPluginConfigurationStore struct {
	store.Store
}

func NewSqlPluginConfigurationStore(s store.Store) store.PluginConfigurationStore {
	return &SqlPluginConfigurationStore{s}
}

func (p *SqlPluginConfigurationStore) Upsert(config model.PluginConfiguration) (*model.PluginConfiguration, error) {
	isSaving := config.ID == ""
	if isSaving {
		model_helper.PluginConfigurationPreSave(&config)
	} else {
		model_helper.PluginConfigurationCommonPre(&config)
	}

	if err := model_helper.PluginConfigurationIsValid(config); err != nil {
		return nil, err
	}

	var err error
	if isSaving {
		err = config.Insert(p.GetMaster(), boil.Infer())
	} else {
		_, err = config.Update(p.GetMaster(), boil.Infer())
	}

	if err != nil {
		if p.IsUniqueConstraintError(err, []string{"plugin_configurations_identifier_channel_id_key", model.PluginConfigurationColumns.Identifier, model.PluginConfigurationColumns.ChannelID}) {
			return nil, store.NewErrInvalidInput(model.TableNames.PluginConfigurations, model.PluginConfigurationColumns.Identifier, config.Identifier)
		}
		return nil, err
	}

	return &config, nil
}

func (p *SqlPluginConfigurationStore) Get(id string) (*model.PluginConfiguration, error) {
	record, err := model.FindPluginConfiguration(p.GetReplica(), id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, store.NewErrNotFound(model.TableNames.PluginConfigurations, id)
		}
		return nil, err
	}

	return record, nil
}

func (p *SqlPluginConfigurationStore) FilterPluginConfigurations(options model_helper.PluginConfigurationFilterOptions) (model.PluginConfigurationSlice, error) {
	conds := options.Conditions
	for _, load := range options.Preloads {
		conds = append(conds, qm.Load(load))
	}

	return model.PluginConfigurations(conds...).All(p.GetReplica())
}
