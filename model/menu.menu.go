package model

import (
	"github.com/Masterminds/squirrel"
	"github.com/gosimple/slug"
	"gorm.io/gorm"
)

// max lengths for some manu's fields
const (
	MENU_NAME_MAX_LENGTH = 250
	MENU_SLUG_MAX_LENGTH = 255
)

type Menu struct {
	Id       string `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid();column:Id"`
	Name     string `json:"name" gorm:"type:varchar(250);column:Name"`
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
	// outer := CreateAppErrorForModel(
	// 	"model.menu.is_valid.%s.app_error",
	// 	"menu_id=",
	// 	"Menu.IsValid",
	// )
	// if !IsValidId(m.Id) {
	// 	return outer("id", nil)
	// }
	// if m.CreateAt <= 0 {
	// 	return outer("create_at", &m.Id)
	// }
	// if utf8.RuneCountInString(m.Name) > MENU_NAME_MAX_LENGTH {
	// 	return outer("Name", &m.Id)
	// }
	// if utf8.RuneCountInString(m.Slug) > MENU_SLUG_MAX_LENGTH {
	// 	return outer("Slug", &m.Id)
	// }

	return nil
}

func (m *Menu) PreSave() {
	m.Name = SanitizeUnicode(m.Name)
	m.Slug = slug.Make(m.Name)
}

func (m *Menu) PreUpdate() {
	m.Name = SanitizeUnicode(m.Name)
}
