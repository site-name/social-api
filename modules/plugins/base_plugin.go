package plugins

import (
	"net/http"
	"reflect"
	"strings"

	"github.com/sitename/sitename/app"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model/channel"
	"github.com/sitename/sitename/model/plugins"
	"github.com/sitename/sitename/modules/util"
)

type BasePlugin struct {
	Manifest PluginManifest

	Active        bool
	Channel       *channel.Channel // can be nil
	Configuration PluginConfigurationType
	srv           *app.Server
}

func NewBasePlugin(active bool, chanNel *channel.Channel, configuration PluginConfigurationType, srv *app.Server) *BasePlugin {
	manifest := PluginManifest{
		ConfigStructure:         make(map[string]model.StringInterface),
		ConfigurationPerChannel: true,
		DefaultConfiguration:    []model.StringInterface{},
	}

	return &BasePlugin{
		Manifest:      manifest,
		Active:        active,
		Channel:       chanNel,
		Configuration: configuration,
		srv:           srv,
	}
}

func (b *BasePlugin) String() string {
	return b.Manifest.Name
}

func (b *BasePlugin) ExternalAuthenticationUrl(data model.StringInterface, request *http.Request, previousValue interface{}) (model.StringInterface, *PluginMethodNotImplemented) {
	return nil, new(PluginMethodNotImplemented)
}

func (b *BasePlugin) CheckPluginId(pluginID string) bool {
	return b.Manifest.ID == pluginID
}

func (b *BasePlugin) GetDefaultActive() bool {
	return b.Manifest.DefaultActive
}

func (b *BasePlugin) UpdateConfigurationStructure(config PluginConfigurationType) PluginConfigurationType {
	var updatedConfiguration []model.StringInterface

	configStructure := b.Manifest.ConfigStructure
	if configStructure == nil {
		configStructure = make(map[string]model.StringInterface)
	}

	desiredConfigKeys := []string{}
	for key := range configStructure {
		desiredConfigKeys = append(desiredConfigKeys, key)
	}
	desiredConfigKeys = util.RemoveDuplicatesFromStringArray(desiredConfigKeys)

	for _, configField := range config {
		if name, ok := configField["name"]; ok && !util.StringInSlice(name.(string), desiredConfigKeys) {
			continue
		}

		updatedConfiguration = append(updatedConfiguration, model.CopyStringInterface(configField))
	}

	configuredKeys := []string{}
	for _, cfg := range updatedConfiguration {
		configuredKeys = append(configuredKeys, cfg["name"].(string)) // name should exist
	}
	configuredKeys = util.RemoveDuplicatesFromStringArray(configuredKeys)

	missingKeys := []string{}
	for _, value := range desiredConfigKeys {
		if !util.StringInSlice(value, configuredKeys) {
			missingKeys = append(missingKeys, value)
		}
	}

	if len(missingKeys) == 0 {
		return updatedConfiguration
	}

	if len(b.Manifest.DefaultConfiguration) == 0 {
		return updatedConfiguration
	}

	updatedValues := []model.StringInterface{}
	for _, item := range b.Manifest.DefaultConfiguration {
		if util.StringInSlice(item["name"].(string), missingKeys) {
			updatedValues = append(updatedValues, model.CopyStringInterface(item))
		}
	}

	if len(updatedValues) > 0 {
		updatedConfiguration = append(updatedConfiguration, updatedValues...)
	}

	return updatedConfiguration
}

func (b *BasePlugin) GetPluginConfiguration(config PluginConfigurationType) PluginConfigurationType {
	if config == nil {
		config = PluginConfigurationType{}
	}

	config = b.UpdateConfigurationStructure(config)

	if len(config) > 0 {
		config = b.AppendConfigStructure(config)
	}

	return config
}

// Append configuration structure to config from the database.
//
// Database stores "key: value" pairs, the definition of fields should be declared
// inside of the plugin. Based on this, the plugin will generate a structure of
// configuration with current values and provide access to it via API.
func (b *BasePlugin) AppendConfigStructure(config PluginConfigurationType) PluginConfigurationType {
	configStructure := b.Manifest.ConfigStructure
	if configStructure == nil {
		configStructure = make(map[string]model.StringInterface)
	}

	fieldsWithoutStructure := []model.StringInterface{}

	for _, configurationField := range config {
		structureToAdd, ok := configStructure[configurationField["name"].(string)]
		if ok && structureToAdd != nil {
			for key, value := range structureToAdd {
				configurationField[key] = value
			}
		} else {
			fieldsWithoutStructure = append(fieldsWithoutStructure, configurationField)
		}
	}

	if len(fieldsWithoutStructure) > 0 {
		for _, field := range fieldsWithoutStructure {
			for idx, item := range config {
				if reflect.DeepEqual(field, item) {
					config = append(config[:idx], config[idx+1:]...)
				}
			}
		}
	}

	return config
}

func (b *BasePlugin) UpdateConfigItems(configurationToUpdate []model.StringInterface, currentConfig []model.StringInterface) []model.StringInterface {
	configStructure := b.Manifest.ConfigStructure
	if configStructure == nil {
		configStructure = make(map[string]model.StringInterface)
	}

	for _, configItem := range currentConfig {
		for _, configItemToUpdate := range configurationToUpdate {

			configItemToUpdateName, ok1 := configItemToUpdate["name"]
			configItemName, ok2 := configItem["name"]

			if ok1 && ok2 && configItemToUpdateName == configItemName {

				newValue, ok3 := configItemToUpdate["value"]

				newValueIsNotNullNorBoolean := ok3 && newValue != nil
				if newValueIsNotNullNorBoolean {
					_, newValueIsBoolean := newValue.(bool)
					newValueIsNotNullNorBoolean = newValueIsNotNullNorBoolean && !newValueIsBoolean
				}

				configStructureValue, ok4 := configStructure[configItemToUpdateName.(string)]

				if !ok4 || configStructureValue == nil {
					configStructureValue = make(model.StringInterface)
				}

				itemType, ok5 := configStructureValue["type"]

				if ok5 &&
					itemType != nil &&
					itemType.(ConfigurationTypeField) == BOOLEAN &&
					newValueIsNotNullNorBoolean {
					newValue = strings.ToLower(newValue.(string)) == "true"
				}

				if val, ok := itemType.(ConfigurationTypeField); ok && val == OUTPUT {
					// OUTPUT field is read only. No need to update it
					continue
				}

				configItem["value"] = newValue
			}
		}
	}

	// Get new keys that don't exist in currentConfig and extend it:
	currentConfigKeys := []string{}
	for _, cField := range currentConfig {
		currentConfigKeys = append(currentConfigKeys, cField["name"].(string))
	}
	currentConfigKeys = util.RemoveDuplicatesFromStringArray(currentConfigKeys)

	configurationToUpdateDict := make(model.StringInterface)
	configurationToUpdateDictKeys := []string{}

	for _, item := range configurationToUpdate {
		configurationToUpdateDict[item["name"].(string)] = item["value"]
		configurationToUpdateDictKeys = append(configurationToUpdateDictKeys, item["name"].(string))
	}
	configurationToUpdateDictKeys = util.RemoveDuplicatesFromStringArray(configurationToUpdateDictKeys)

	for _, item := range configurationToUpdateDictKeys {
		if !util.StringInSlice(item, currentConfigKeys) {
			if val, ok := configStructure[item]; !ok || val == nil {
				continue
			}

			currentConfig = append(currentConfig, model.StringInterface{
				"name":  item,
				"value": configurationToUpdateDict[item],
			})
		}
	}

	return currentConfig
}

func (b *BasePlugin) SavePluginConfiguration(pluginConfiguration *plugins.PluginConfiguration, cleanedData model.StringInterface) (*plugins.PluginConfiguration, *model.AppError, *PluginMethodNotImplemented) {
	currentConfig := pluginConfiguration.Configuration
	configurationToUpdate, ok := cleanedData["configuration"]

	if ok && configurationToUpdate != nil {
		pluginConfiguration.Configuration = b.UpdateConfigItems(configurationToUpdate.([]model.StringInterface), currentConfig)
	}

	if active, ok := cleanedData["active"]; ok && active != nil {
		pluginConfiguration.Active = active.(bool)
	}

	appErr, notImplt := b.ValidatePluginConfiguration(pluginConfiguration)
	if notImplt != nil {
		return nil, nil, notImplt
	}
	if appErr != nil {
		return nil, appErr, nil
	}
	appErr, notImplt = b.PreSavePluginConfiguration(pluginConfiguration)
	if notImplt != nil {
		return nil, nil, notImplt
	}
	if appErr != nil {
		return nil, appErr, nil
	}

	pluginConfiguration, appErr = b.srv.PluginService().UpsertPluginConfiguration(pluginConfiguration)
	if appErr != nil {
		return nil, appErr, nil
	}

	if len(pluginConfiguration.Configuration) > 0 {
		pluginConfiguration.Configuration = b.AppendConfigStructure(pluginConfiguration.Configuration)
	}

	return pluginConfiguration, nil, nil
}

func (b *BasePlugin) ValidatePluginConfiguration(pluginConfiguration *plugins.PluginConfiguration) (*model.AppError, *PluginMethodNotImplemented) {
	return nil, new(PluginMethodNotImplemented)
}

func (b *BasePlugin) PreSavePluginConfiguration(pluginConfiguration *plugins.PluginConfiguration) (*model.AppError, *PluginMethodNotImplemented) {
	return nil, new(PluginMethodNotImplemented)
}
