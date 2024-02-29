package model_helper

import (
	"net/http"

	"github.com/sitename/sitename/model"
)

func PluginConfigurationPreSave(p *model.PluginConfiguration) {
	if p.ID == "" {
		p.ID = NewId()
	}
	PluginConfigurationCommonPre(p)
}

func PluginConfigurationCommonPre(p *model.PluginConfiguration) {
	p.Identifier = SanitizeUnicode(p.Identifier)
	p.ChannelID = SanitizeUnicode(p.ChannelID)
	p.Name = SanitizeUnicode(p.Name)
	p.Description = SanitizeUnicode(p.Description)
}

func PluginConfigurationIsValid(p model.PluginConfiguration) *AppError {
	if !IsValidId(p.ID) {
		return NewAppError("PluginConfiguration.IsValid", "model.plugin_configuration.is_valid.id.app_error", nil, "please provide valid id", http.StatusBadRequest)
	}
	if p.Identifier == "" {
		return NewAppError("PluginConfiguration.IsValid", "model.plugin_configuration.is_valid.identifier.app_error", nil, "please provide valid identifier", http.StatusBadRequest)
	}
	if !IsValidId(p.ChannelID) {
		return NewAppError("PluginConfiguration.IsValid", "model.plugin_configuration.is_valid.channel_id.app_error", nil, "please provide valid channel id", http.StatusBadRequest)
	}
	if p.Name == "" {
		return NewAppError("PluginConfiguration.IsValid", "model.plugin_configuration.is_valid.name.app_error", nil, "please provide valid name", http.StatusBadRequest)
	}
	return nil
}
