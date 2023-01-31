package model

import (
	"unicode/utf8"

	"github.com/Masterminds/squirrel"
	"github.com/gosimple/slug"
	"github.com/samber/lo"
)

// max length for some fields
const (
	CATEGORY_NAME_MAX_LENGTH         = 250
	CATEGORY_SLUG_MAX_LENGTH         = 255
	CATEGORY_BG_IMAGE_ALT_MAX_LENGTH = 128
)

type Category struct {
	Id                 string          `json:"id"`
	Name               string          `json:"name"` // unique
	Slug               string          `json:"slug"` // unique
	NameEn             string          `json:"name_en"`
	Description        StringInterface `json:"description"`
	ParentID           *string         `json:"parent_id"`
	BackgroundImage    *string         `json:"background_image"`
	BackgroundImageAlt string          `json:"background_image_alt"`
	Seo
	ModelMetadata

	Children Categories `json:"children"`
	Images   []string   `json:"images"`
}

// CategoryFilterOption is used for building sql queries
type CategoryFilterOption struct {
	All  bool // if true, select all categories
	Id   squirrel.Sqlizer
	Name squirrel.Sqlizer
	Slug squirrel.Sqlizer

	SaleID    squirrel.Sqlizer // SELECT * FROM Categories INNER JOIN SaleCategories ON (Categories.Id = SaleCategories.CategoryID) WHERE SaleCategories.SaleID ...
	VoucherID squirrel.Sqlizer // SELECT * FROM Categories INNER JOIN VoucherCategories ON (VoucherCategories.CategoryID = Categories.Id) WHERE VoucherCategories.VoucherID ...
	ProductID squirrel.Sqlizer // SELECT * FROM Categories INNER JOIN Products (ON ...) WHERE ProductID IN (...)

	LockForUpdate bool // set this to true if you want to add "FOR UPDATE" suffix to the end of queries
}

type Categories []*Category

func (c Categories) IDs() []string {
	return lo.Map(c, func(g *Category, _ int) string { return g.Id })
}

func (cs Categories) DeepCopy() Categories {
	return lo.Map(cs, func(g *Category, _ int) *Category { return g.DeepCopy() })

}

func (c *Category) String() string {
	return c.Name
}

func (c *Category) IsValid() *AppError {
	outer := CreateAppErrorForModel(
		"category.is_valid.%s.app_error",
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
	if len(s.Children) > 0 {
		res.Children = s.Children.DeepCopy()
	}
	return &res
}

func (c *Category) PreSave() {
	if c.Id == "" {
		c.Id = NewId()
	}
	c.Name = SanitizeUnicode(c.Name)
	c.Slug = slug.Make(c.Slug)
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
		"category_translation.is_valid.%s.app_error",
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

// ClassifyCategories takes a slice of single categories.
// Returns a slice of category families
func ClassifyCategories(categories Categories) Categories {
	if len(categories) <= 1 {
		return categories
	}

	var res Categories

	// trackMap has keys are category ids
	var trackMap = map[string]*Category{}

	for _, cate := range categories {
		if cate != nil {
			trackMap[cate.Id] = cate
		}
	}

	for _, cate := range categories {
		if cate != nil {
			if cate.ParentID == nil { // category has no child category
				res = append(res, cate)
				continue
			}

			trackMap[*cate.ParentID].Children = append(trackMap[*cate.ParentID].Children, cate)
		}
	}

	return res
}
