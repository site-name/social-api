package model

import (
	"net/http"

	"github.com/Masterminds/squirrel"
	"github.com/samber/lo"
	"gorm.io/gorm"
)

type PluginConfiguration struct {
	Id            string          `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid();column:Id"`
	Identifier    string          `json:"identifier" gorm:"type:varchar(128);column:Identifier;uniqueIndex:identifier_channelid_key"`
	Name          string          `json:"name" gorm:"type:varchar(128);column:Name"`
	ChannelID     string          `json:"channel_id" gorm:"type:uuid;column:ChannelID;uniqueIndex:identifier_channelid_key"`
	Description   string          `json:"description" gorm:"column:Description"`
	Active        bool            `json:"active" gorm:"column:Active"`
	Configuration StringInterface `json:"configuration" gorm:"type:jsonb;column:Configuration"`

	relatedChannel *Channel `db:"-"` // this field is populated in some sql queries
}

func (c *PluginConfiguration) BeforeCreate(_ *gorm.DB) error { c.commonPre(); return c.IsValid() }
func (c *PluginConfiguration) BeforeUpdate(_ *gorm.DB) error { c.commonPre(); return c.IsValid() }
func (c *PluginConfiguration) TableName() string             { return TransactionTableName }

type PluginConfigurations []*PluginConfiguration

func (p PluginConfigurations) IDs() []string {
	return lo.Map(p, func(c *PluginConfiguration, _ int) string { return c.Id })
}

func (p PluginConfigurations) ChannelIDs() []string {
	return lo.Map(p, func(c *PluginConfiguration, _ int) string { return c.ChannelID })
}

// PluginConfigurationFilterOptions is used to build sql queries
type PluginConfigurationFilterOptions struct {
	Conditions squirrel.Sqlizer

	PrefetchRelatedChannel bool // this tells store to prefetch related channel also
}

func (p *PluginConfiguration) SetRelatedChannel(ch *Channel) {
	p.relatedChannel = ch
}

func (p *PluginConfiguration) GetRelatedChannel() *Channel {
	return p.relatedChannel
}

func (p *PluginConfiguration) IsValid() *AppError {
	if !IsValidId(p.ChannelID) {
		return NewAppError("PluginConfiguration.IsValid", "model.plugin_config.is_valid.channel_id.app_error", nil, "please provide valid channel id", http.StatusBadRequest)
	}

	return nil
}

func (p *PluginConfiguration) commonPre() {
	p.Identifier = SanitizeUnicode(p.Identifier)
	p.Name = SanitizeUnicode(p.Name)
	p.Description = SanitizeUnicode(p.Description)

	if p.Configuration == nil {
		p.Configuration = StringInterface{}
	}
}
