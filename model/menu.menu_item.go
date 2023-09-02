package model

import (
	"net/http"

	"github.com/Masterminds/squirrel"
	"gorm.io/gorm"
)

type MenuItem struct {
	Id           UUID    `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid();column:Id"`
	MenuID       UUID    `json:"menu_id" gorm:"type:uuid;column:MenuID"`
	Name         string  `json:"name" gorm:"type:varchar(128);column:Name"`
	ParentID     *UUID   `json:"parent_id" gorm:"type:uuid;column:ParentID"` // foreign key menu item
	Url          *string `json:"url" gorm:"type:varchar(256);column:Url"`
	CategoryID   *UUID   `json:"category_id" gorm:"type:uuid;column:CategoryID"`     // to category
	CollectionID *UUID   `json:"collection_id" gorm:"type:uuid;column:CollectionID"` // to collection
	PageID       *UUID   `json:"page_id" gorm:"type:uuid;column:PageID"`
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
	if !IsValidId(m.MenuID) {
		return NewAppError("MenuItem.IsValid", "model.menu_item.is_valid.menu_id.app_error", nil, "please provide valid menu item menu id", http.StatusBadRequest)
	}
	if m.ParentID != nil && !IsValidId(*m.ParentID) {
		return NewAppError("MenuItem.IsValid", "model.menu_item.is_valid.parent_id.app_error", nil, "please provide valid menu item parent id", http.StatusBadRequest)
	}
	if m.CategoryID != nil && !IsValidId(*m.CategoryID) {
		return NewAppError("MenuItem.IsValid", "model.menu_item.is_valid.category_id.app_error", nil, "please provide valid menu item category id", http.StatusBadRequest)
	}
	if m.PageID != nil && !IsValidId(*m.PageID) {
		return NewAppError("MenuItem.IsValid", "model.menu_item.is_valid.page_id.app_error", nil, "please provide valid menu item page id", http.StatusBadRequest)
	}
	if m.CollectionID != nil && !IsValidId(*m.CollectionID) {
		return NewAppError("MenuItem.IsValid", "model.menu_item.is_valid.collection_id.app_error", nil, "please provide valid menu item collection id", http.StatusBadRequest)
	}

	return nil
}

func (m *MenuItem) commonPre() {
	m.Name = SanitizeUnicode(m.Name)
}
