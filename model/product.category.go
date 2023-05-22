package model

import (
	"unicode/utf8"

	"github.com/Masterminds/squirrel"
	"github.com/gosimple/slug"
	"github.com/samber/lo"
	"github.com/sitename/sitename/modules/util"
)

// max length for some fields
const (
	CATEGORY_NAME_MAX_LENGTH         = 250
	CATEGORY_SLUG_MAX_LENGTH         = 255
	CATEGORY_BG_IMAGE_ALT_MAX_LENGTH = 128

	CATEGORY_MIN_LEVEL = 0
	CATEGORY_MAX_LEVEL = 4

	CATEGORY_IMAGES_MAX_LENGTH = 1000
)

type Category struct {
	Id                 string          `json:"id"`
	Name               string          `json:"name"` // unique, English
	Slug               string          `json:"slug"` // unique
	Description        StringInterface `json:"description,omitempty"`
	ParentID           *string         `json:"parent_id,omitempty"`
	Level              uint8           `json:"level"` // 0, 1, 2, 3, 4
	BackgroundImage    *string         `json:"background_image,omitempty"`
	BackgroundImageAlt string          `json:"background_image_alt"`
	Images             string          `json:"images"` // space-seperated urls
	Seo
	NameTranslation StringMAP `json:"name_translation,omitempty"` // e.g {"vi": "Xin Chao"}
	ModelMetadata

	NumOfProducts uint64 `json:"num_of_products" db:"-"` // this field gets fulfilled in some db quesries
	// Children      Categories `json:"children,omitempty" db:"-"` // this field gets populated sometimes
	NumOfChildren int `json:"num_of_children" db:"-"`
}

// CategoryFilterOption is used for building sql queries
type CategoryFilterOption struct {
	All   bool // if true, select all categories
	Id    squirrel.Sqlizer
	Name  squirrel.Sqlizer
	Level squirrel.Sqlizer
	Slug  squirrel.Sqlizer
	Extra squirrel.Sqlizer

	SaleID    squirrel.Sqlizer // SELECT * FROM Categories INNER JOIN SaleCategories ON (Categories.Id = SaleCategories.CategoryID) WHERE SaleCategories.SaleID ...
	VoucherID squirrel.Sqlizer // SELECT * FROM Categories INNER JOIN VoucherCategories ON (VoucherCategories.CategoryID = Categories.Id) WHERE VoucherCategories.VoucherID ...
	ProductID squirrel.Sqlizer // SELECT * FROM Categories INNER JOIN Products (ON ...) WHERE ProductID IN (...)

	LockForUpdate bool // set this to true if you want to add "FOR UPDATE" suffix to the end of queries

	OrderBy string
	Limit   uint64
}

type Categories []*Category

// set flat to true to recursively get all ids of child categories to
func (cs Categories) IDs(flat bool) util.AnyArray[string] {
	// if !flat {
	return lo.Map(cs, func(g *Category, _ int) string { return g.Id })
	// }

	// var res util.AnyArray[string]
	// for _, cate := range cs {
	// 	res = append(res, cate.Id)
	// 	if cate.Children.Len() > 0 {
	// 		res = append(res, cate.Children.IDs(flat)...)
	// 	}
	// }

	// return res
}

func (c Categories) Len() int {
	return len(c)
}

func (cs Categories) DeepCopy() Categories {
	return lo.Map(cs, func(g *Category, _ int) *Category { return g.DeepCopy() })
}

func (c *Category) String() string {
	return c.Name
}

func (c *Category) IsValid() *AppError {
	outer := CreateAppErrorForModel(
		"model.category.is_valid.%s.app_error",
		"category_id=",
		"Category.IsValid",
	)

	if !IsValidId(c.Id) {
		return outer("id", nil)
	}
	if c.ParentID != nil && !IsValidId(*c.ParentID) {
		return outer("id", &c.Id)
	}
	if len(c.BackgroundImageAlt) > CATEGORY_BG_IMAGE_ALT_MAX_LENGTH {
		return outer("background_image_alt", &c.Id)
	}
	if utf8.RuneCountInString(c.Name) > CATEGORY_NAME_MAX_LENGTH {
		return outer("name", &c.Id)
	}
	if utf8.RuneCountInString(c.Slug) > CATEGORY_SLUG_MAX_LENGTH {
		return outer("slug", &c.Id)
	}
	if c.Level < CATEGORY_MIN_LEVEL || c.Level > CATEGORY_MAX_LEVEL {
		return outer("level", &c.Id)
	}
	if c.Images != "" && len(c.Images) > CATEGORY_IMAGES_MAX_LENGTH {
		return outer("images", &c.Id)
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
		res.ParentID = NewPrimitive(*s.ParentID)
	}
	if s.BackgroundImage != nil {
		res.BackgroundImage = NewPrimitive(*s.BackgroundImage)
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
	if c.Id == "" {
		c.Id = NewId()
	}
	c.Name = SanitizeUnicode(c.Name)
	if c.Slug == "" {
		c.Slug = slug.Make(c.Name)
	}
}

func (c *Category) PreUpdate() {
	c.Name = SanitizeUnicode(c.Name)
}

func (c *Category) ToJSON() string {
	return ModelToJson(c)
}

type CategoryTranslation struct {
	Id           string  `json:"id"`
	LanguageCode string  `json:"language_code"`
	CategoryID   string  `json:"category_id"`
	Name         string  `json:"name"`
	Description  *string `json:"description"`
	SeoTranslation
}

func (c *CategoryTranslation) String() string {
	return c.Name
}

func (c *CategoryTranslation) ToJSON() string {
	return ModelToJson(c)
}

func (c *CategoryTranslation) IsValid() *AppError {
	outer := CreateAppErrorForModel(
		"model.category_translation.is_valid.%s.app_error",
		"category_translation_id=",
		"CategoryTranslation.IsValid")

	if !IsValidId(c.CategoryID) {
		return outer("category_id", &c.Id)
	}
	if utf8.RuneCountInString(c.Name) > CATEGORY_NAME_MAX_LENGTH {
		return outer("name", &c.Id)
	}

	return nil
}

func (c *CategoryTranslation) PreSave() {
	if c.Id == "" {
		c.Id = NewId()
	}
	c.Name = SanitizeUnicode(c.Name)
}
