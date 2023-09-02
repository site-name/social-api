package model

import (
	"net/http"

	"github.com/Masterminds/squirrel"
	"github.com/gosimple/slug"
	"gorm.io/gorm"
)

type Page struct {
	Id         UUID            `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid();column:Id"`
	Title      string          `json:"title" gorm:"type:varchar(250);column:Title"`                           //
	Slug       string          `json:"slug" gorm:"type:varchar(255);column:Slug;uniqueIndex:slug_unique_key"` // unique
	PageTypeID UUID            `json:"page_type_id" gorm:"type:uuid;column:PageTypeID"`
	Content    StringInterface `json:"content" gorm:"column:Content"`
	CreateAt   int64           `json:"create_at" gorm:"type:bigint;column:CreateAt"`
	ModelMetadata
	Publishable
	Seo

	Attributes       []*AssignedPageAttribute `json:"-" gorm:"foreignKey:AssignmentID;constraint:OnDelete:CASCADE;"`
	AttributesRelate []*AttributePage         `json:"-" gorm:"many2many:AssignedPageAttributes"`
}

func (c *Page) BeforeCreate(_ *gorm.DB) error { c.PreSave(); return c.IsValid() }
func (c *Page) BeforeUpdate(_ *gorm.DB) error { c.PreUpdate(); return c.IsValid() }
func (c *Page) TableName() string             { return PageTableName }

type PageFilterOptions struct {
	Conditions squirrel.Sqlizer
}

func (p *Page) IsValid() *AppError {
	if !IsValidId(p.PageTypeID) {
		return NewAppError("Page.IsValid", "model.page.is_valid.page_type_id.app_error", nil, "please provide valid page type id", http.StatusBadRequest)
	}

	return nil
}

func (p *Page) PreSave() {
	p.Title = SanitizeUnicode(p.Title)
	p.Slug = slug.Make(p.Title)
}

func (p *Page) PreUpdate() {
	p.Title = SanitizeUnicode(p.Title)
}

func (p *Page) String() string {
	return p.Title
}
