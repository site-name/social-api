package plugins

import (
	"unicode/utf8"

	"github.com/sitename/sitename/model"
)

// max length for some fields
const (
	PLUGIN_CONFIGURATION_COMMON_MAX_LENGHT      = 128
	PLUGIN_CONFIGURATION_DESCRIPTION_MAX_LENGHT = 1000
)

type PluginConfiguration struct {
	Id            string                  `json:"id"`
	Identifier    string                  `json:"identifier"`
	Name          string                  `json:"name"`
	ChannelID     string                  `json:"channel_id"`
	Description   string                  `json:"description"`
	Active        bool                    `json:"active"`
	Configuration []model.StringInterface `json:"configuration"` // default [{}]
}

func (p *PluginConfiguration) IsValid() *model.AppError {
	outer := model.CreateAppErrorForModel(
		"model.plugin.is_valid.%s.app_error",
		"plugin_id=",
		"Plugin.IsValid",
	)
	if !model.IsValidId(p.Id) {
		return outer("id", nil)
	}
	if !model.IsValidId(p.ChannelID) {
		return outer("channel_id", &p.Id)
	}
	if utf8.RuneCountInString(p.Identifier) > PLUGIN_CONFIGURATION_COMMON_MAX_LENGHT {
		return outer("identifier", &p.Id)
	}
	if utf8.RuneCountInString(p.Name) > PLUGIN_CONFIGURATION_COMMON_MAX_LENGHT {
		return outer("name", &p.Id)
	}
	if utf8.RuneCountInString(p.Description) > PLUGIN_CONFIGURATION_DESCRIPTION_MAX_LENGHT {
		return outer("description", &p.Id)
	}

	return nil
}

func (p *PluginConfiguration) ToJSON() string {
	return model.ModelToJson(p)
}

func (p *PluginConfiguration) PreSave() {
	if p.Id == "" {
		p.Id = model.NewId()
	}
	p.commonPre()
}

func (p *PluginConfiguration) commonPre() {
	p.Identifier = model.SanitizeUnicode(p.Identifier)
	p.Name = model.SanitizeUnicode(p.Name)
	p.Description = model.SanitizeUnicode(p.Description)

	if p.Configuration == nil {
		p.Configuration = []model.StringInterface{}
	}
}

func (p *PluginConfiguration) PreUpdate() {
	p.commonPre()
}
