package model

import (
	"github.com/Masterminds/squirrel"
	"github.com/gosimple/slug"
	"gorm.io/gorm"
)

type Menu struct {
	Id       UUID   `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid();column:Id"`
	Name     string `json:"name" gorm:"type:varchar(250);column:Name;uniqueIndex:name_unique_key"`
	Slug     string `json:"slug" gorm:"type:varchar(255);column:Slug;uniqueIndex:slug_unique_key"` // unique, index
	CreateAt int64  `json:"create_at" gorm:"type:bigint;column:CreateAt;autoCreateTime:milli"`     // this field can be used for ordering
	ModelMetadata
}

func (c *Menu) BeforeCreate(_ *gorm.DB) error { c.PreSave(); return c.IsValid() }
func (c *Menu) BeforeUpdate(_ *gorm.DB) error { c.PreUpdate(); return c.IsValid() }
func (c *Menu) TableName() string             { return MenuTableName }

type MenuFilterOptions struct {
	Conditions squirrel.Sqlizer
}

func (m *Menu) IsValid() *AppError {
	return nil
}

func (m *Menu) PreSave() {
	m.Name = SanitizeUnicode(m.Name)
	if m.Slug == "" {
		m.Slug = slug.Make(m.Name)
	}
}

func (m *Menu) PreUpdate() {
	m.Name = SanitizeUnicode(m.Name)
}
