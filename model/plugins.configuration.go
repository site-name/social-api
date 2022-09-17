package model

import (
	"unicode/utf8"

	"github.com/Masterminds/squirrel"
)

// max length for some fields
const (
	PLUGIN_CONFIGURATION_COMMON_MAX_LENGHT      = 128
	PLUGIN_CONFIGURATION_DESCRIPTION_MAX_LENGHT = 1000
)

type PluginConfiguration struct {
	Id            string            `json:"id"`
	Identifier    string            `json:"identifier"`
	Name          string            `json:"name"`
	ChannelID     string            `json:"channel_id"`
	Description   string            `json:"description"`
	Active        bool              `json:"active"`
	Configuration []StringInterface `json:"configuration"` // default [{}]

	relatedChannel *Channel `json:"-" db:"-"` // this field is populated in some sql queries
}

type PluginConfigurations []*PluginConfiguration

func (p PluginConfigurations) IDs() []string {
	var res []string
	for _, item := range p {
		if item != nil {
			res = append(res, item.Id)
		}
	}

	return res
}

func (p PluginConfigurations) ChannelIDs() []string {
	var res []string
	for _, item := range p {
		if item != nil {
			res = append(res, item.ChannelID)
		}
	}

	return res
}

// PluginConfigurationFilterOptions is used to build sql queries
type PluginConfigurationFilterOptions struct {
	Id         squirrel.Sqlizer
	Identifier squirrel.Sqlizer
	ChannelID  squirrel.Sqlizer

	PrefetchRelatedChannel bool // this tells store to prefetch related channel also
}

func (p *PluginConfiguration) SetRelatedChannel(ch *Channel) {
	p.relatedChannel = ch
}

func (p *PluginConfiguration) GetRelatedChannel() *Channel {
	return p.relatedChannel
}

func (p *PluginConfiguration) IsValid() *AppError {
	outer := CreateAppErrorForModel(
		"plugin.is_valid.%s.app_error",
		"plugin_id=",
		"Plugin.IsValid",
	)
	if !IsValidId(p.Id) {
		return outer("id", nil)
	}
	if !IsValidId(p.ChannelID) {
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
	return ModelToJson(p)
}

func (p *PluginConfiguration) PreSave() {
	if p.Id == "" {
		p.Id = NewId()
	}
	p.commonPre()
}

func (p *PluginConfiguration) commonPre() {
	p.Identifier = SanitizeUnicode(p.Identifier)
	p.Name = SanitizeUnicode(p.Name)
	p.Description = SanitizeUnicode(p.Description)

	if p.Configuration == nil {
		p.Configuration = []StringInterface{}
	}
}

func (p *PluginConfiguration) PreUpdate() {
	p.commonPre()
}
