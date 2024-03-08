package model

import (
	"github.com/mattermost/squirrel"
	"gorm.io/gorm"
)

type MenuItem struct {
	Id           string  `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid();column:Id"`
	MenuID       string  `json:"menu_id" gorm:"type:uuid;column:MenuID"`
	Name         string  `json:"name" gorm:"type:varchar(128);column:Name"`
	ParentID     *string `json:"parent_id" gorm:"type:uuid;column:ParentID"` // foreign key menu item
	Url          *string `json:"url" gorm:"type:varchar(256);column:Url"`
	CategoryID   *string `json:"category_id" gorm:"type:uuid;column:CategoryID"`     // to category
	CollectionID *string `json:"collection_id" gorm:"type:uuid;column:CollectionID"` // to collection
	PageID       *string `json:"page_id" gorm:"type:uuid;column:PageID"`
	ModelMetadata
	Sortable
}

func (c *MenuItem) BeforeCreate(_ *gorm.DB) error { c.commonPre(); return c.IsValid() }
func (c *MenuItem) BeforeUpdate(_ *gorm.DB) error { c.commonPre(); return c.IsValid() }
func (c *MenuItem) TableName() string             { return MenuItemTableName }

type MenuItemFilterOptions struct {
	Conditions squirrel.Sqlizer
}

func (m *MenuItem) IsValid() *AppError {
	outer := CreateAppErrorForModel(
		"model.menu_item.is_valid.%s.app_error",
		"menu_item_id=",
		"MenuItem.IsValid",
	)
	if !IsValidId(m.MenuID) {
		return outer("menu_id", &m.Id)
	}
	if m.ParentID != nil && !IsValidId(*m.ParentID) {
		return outer("parent_id", &m.Id)
	}
	if m.CategoryID != nil && !IsValidId(*m.CategoryID) {
		return outer("category_id", &m.Id)
	}
	if m.PageID != nil && !IsValidId(*m.PageID) {
		return outer("page_id", &m.Id)
	}
	if m.CollectionID != nil && !IsValidId(*m.CollectionID) {
		return outer("collection_id", &m.Id)
	}

	return nil
}

func (m *MenuItem) commonPre() {
	m.Name = SanitizeUnicode(m.Name)
}
