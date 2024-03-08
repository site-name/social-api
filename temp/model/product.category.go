package model

import (
	"net/http"

	"github.com/gosimple/slug"
	"github.com/mattermost/squirrel"
	"github.com/samber/lo"
	"github.com/sitename/sitename/modules/util"
	"gorm.io/gorm"
)

const (
	CATEGORY_MIN_LEVEL = 0
)

type Category struct {
	Id                 string          `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid();column:Id"`
	Name               string          `json:"name" gorm:"unique;type:varchar(250);column:Name"`                      // unique, English
	Slug               string          `json:"slug" gorm:"uniqueIndex:slug_unique_key;type:varchar(255);column:Slug"` // unique
	Description        StringInterface `json:"description,omitempty" gorm:"type:jsonb;column:Description"`
	ParentID           *string         `json:"parent_id,omitempty" gorm:"type:uuid;column:ParentID"`
	Level              uint8           `json:"level" gorm:"type:smallint;check:level >= 0;column:Level"` // 0, 1, 2, 3, 4
	BackgroundImage    *string         `json:"background_image,omitempty" gorm:"type:varchar(1000);column:BackgroundImage"`
	BackgroundImageAlt string          `json:"background_image_alt" gorm:"type:varchar(128);column:BackgroundImageAlt"`
	Images             string          `json:"images" gorm:"type:varchar(1000);column:Images"`                      // space-seperated urls
	NameTranslation    StringMAP       `json:"name_translation,omitempty" gorm:"type:jsonb;column:NameTranslation"` // e.g {"vi": "Xin Chao"}
	Seo
	ModelMetadata

	NumOfProducts uint64 `json:"num_of_products" gorm:"-"` // this field gets fulfilled in some db quesries
	NumOfChildren int    `json:"num_of_children" gorm:"-"`

	Sales    Sales    `json:"-" gorm:"many2many:SaleCategories"`
	Vouchers Vouchers `json:"-" gorm:"many2many:VoucherCategories"`

	// Children      Categories `json:"children,omitempty" db:"-"` // this field gets populated sometimes
}

func (c *Category) BeforeCreate(_ *gorm.DB) error { c.PreSave(); return c.IsValid() }
func (c *Category) BeforeUpdate(_ *gorm.DB) error { c.PreUpdate(); return c.IsValid() }
func (c *Category) TableName() string             { return CategoryTableName }
func (c Categories) Len() int                     { return len(c) }
func (c *Category) String() string                { return c.Name }

// CategoryFilterOption is used for building sql queries
type CategoryFilterOption struct {
	Conditions squirrel.Sqlizer

	SaleID    squirrel.Sqlizer // SELECT * FROM Categories INNER JOIN SaleCategories ON (Categories.Id = SaleCategories.CategoryID) WHERE SaleCategories.SaleID ...
	VoucherID squirrel.Sqlizer // SELECT * FROM Categories INNER JOIN VoucherCategories ON (VoucherCategories.CategoryID = Categories.Id) WHERE VoucherCategories.VoucherID ...
	ProductID squirrel.Sqlizer // SELECT * FROM Categories INNER JOIN Products (ON ...) WHERE ProductID IN (...)

	LockForUpdate bool // set this to true if you want to add "FOR UPDATE" suffix to the end of queries. NOTE: only applies when Transaction field is set
	Transaction   *gorm.DB

	OrderBy string
	Limit   uint64
}

type Categories []*Category

// set flat to true to recursively get all ids of child categories to
func (cs Categories) IDs(flat bool) util.AnyArray[string] {
	return lo.Map(cs, func(g *Category, _ int) string { return g.Id })
}

func (ps Categories) Contains(c *Category) bool {
	return c != nil && lo.SomeBy(ps, func(ct *Category) bool { return ct != nil && ct.Id == c.Id })
}

func (cs Categories) DeepCopy() Categories {
	return lo.Map(cs, func(g *Category, _ int) *Category { return g.DeepCopy() })
}

func (c *Category) IsValid() *AppError {
	if c.ParentID != nil && !IsValidId(*c.ParentID) {
		return NewAppError("Category.IsValid", "model.category.is_valid.parent_id.app_error", nil, "please provide valid parent id", http.StatusBadRequest)
	}
	return nil
}

func (s *Category) DeepCopy() *Category {
	if s == nil {
		return nil
	}

	res := *s
	if s.Description != nil {
		res.Description = s.Description.DeepCopy()
	}
	if s.ParentID != nil {
		*res.ParentID = *s.ParentID
	}
	if s.BackgroundImage != nil {
		*res.BackgroundImage = *s.BackgroundImage
	}
	if s.NameTranslation != nil {
		res.NameTranslation = s.NameTranslation.DeepCopy()
	}
	// if len(s.Children) > 0 {
	// 	res.Children = s.Children.DeepCopy()
	// }
	return &res
}

func (c *Category) PreSave() {
	c.commonPre()
	if c.Slug == "" {
		c.Slug = slug.Make(c.Name)
	}
}

func (c *Category) commonPre() {
	c.Seo.commonPre()
	c.Name = SanitizeUnicode(c.Name)
}

func (c *Category) PreUpdate() {
	c.commonPre()
}

type CategoryTranslation struct {
	Id           string           `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid();column:Id"`
	LanguageCode LanguageCodeEnum `json:"language_code" gorm:"type:varchar(5);column:LanguageCode"`
	CategoryID   string           `json:"category_id" gorm:"type:uuid;column:CategoryID"`
	Name         string           `json:"name" gorm:"type:varchar(250);column:Name"`
	Description  StringInterface  `json:"description" gorm:"type:jsonb;column:Description"`
	SeoTranslation
}

func (c *CategoryTranslation) BeforeCreate(_ *gorm.DB) error { c.commonPre(); return c.IsValid() }
func (c *CategoryTranslation) BeforeUpdate(_ *gorm.DB) error { c.commonPre(); return c.IsValid() }
func (c *CategoryTranslation) TableName() string             { return CategoryTranslationTableName }
func (c *CategoryTranslation) String() string                { return c.Name }

func (c *CategoryTranslation) IsValid() *AppError {
	if !IsValidId(c.CategoryID) {
		return NewAppError("CategoryTranslation.IsValid", "model.category_translation.is_valid.category_id.app_error", nil, "please provide valid category id", http.StatusBadRequest)
	}
	if !c.LanguageCode.IsValid() {
		return NewAppError("CategoryTranslation.IsValid", "model.category_translation.is_valid.language_code.app_error", nil, "please provide valid language code", http.StatusBadRequest)
	}

	return nil
}

func (c *CategoryTranslation) commonPre() {
	c.Name = SanitizeUnicode(c.Name)
}
