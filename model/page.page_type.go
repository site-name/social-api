package model

import (
	"github.com/gosimple/slug"
	"gorm.io/gorm"
)

type PageType struct {
	Id   UUID   `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid();column:Id"`
	Name string `json:"name" gorm:"type:varchar(250);column:Name"`
	Slug string `json:"alug" gorm:"uniqueIndex:slug_key;type:varchar(255);column:Slug"`
	ModelMetadata

	AttributePages []*AttributePage `json:"-" gorm:"foreignKey:PageTypeID;constraint:OnDelete:CASCADE;"`
}

func (c *PageType) BeforeCreate(_ *gorm.DB) error { c.PreSave(); return c.IsValid() }
func (c *PageType) BeforeUpdate(_ *gorm.DB) error { c.PreUpdate(); return c.IsValid() }
func (c *PageType) TableName() string             { return PageTypeTableName }

func (pt *PageType) IsValid() *AppError {
	return nil
}

func (pt *PageType) PreSave() {
	pt.Name = SanitizeUnicode(pt.Name)
	pt.Slug = slug.Make(pt.Name)
}

func (pt *PageType) PreUpdate() {
	pt.Name = SanitizeUnicode(pt.Name)
}
