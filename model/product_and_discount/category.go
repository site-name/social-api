package product_and_discount

import (
	"unicode/utf8"

	"github.com/Masterminds/squirrel"
	"github.com/gosimple/slug"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model/seo"
)

// max length for some fields
const (
	CATEGORY_NAME_MAX_LENGTH         = 250
	CATEGORY_SLUG_MAX_LENGTH         = 255
	CATEGORY_BG_IMAGE_ALT_MAX_LENGTH = 128
)

type Category struct {
	Id                 string  `json:"id"`
	Name               string  `json:"name"` // unique
	Slug               string  `json:"slug"` // unique
	Description        *string `json:"description"`
	ParentID           *string `json:"parent_id"`
	BackgroundImage    *string `json:"background_image"`
	BackgroundImageAlt string  `json:"background_image_alt"`
	seo.Seo
	model.ModelMetadata

	Children Categories `db:"-"`
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
	res := []string{}
	for _, item := range c {
		if item != nil {
			res = append(res, item.Id)
		}
	}
	return res
}

func (c *Category) String() string {
	return c.Name
}

func (c *Category) IsValid() *model.AppError {
	outer := model.CreateAppErrorForModel(
		"model.category.is_valid.%s.app_error",
		"category_id=",
		"Category.IsValid",
	)

	if !model.IsValidId(c.Id) {
		return outer("id", nil)
	}
	if c.ParentID != nil && !model.IsValidId(*c.ParentID) {
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

func (c *Category) PreSave() {
	if c.Id == "" {
		c.Id = model.NewId()
	}
	c.Name = model.SanitizeUnicode(c.Name)
	c.Slug = slug.Make(c.Slug)
}

func (c *Category) PreUpdate() {
	c.Name = model.SanitizeUnicode(c.Name)
}

func (c *Category) ToJSON() string {
	return model.ModelToJson(c)
}

type CategoryTranslation struct {
	Id           string  `json:"id"`
	LanguageCode string  `json:"language_code"`
	CategoryID   string  `json:"category_id"`
	Name         string  `json:"name"`
	Description  *string `json:"description"`
	seo.SeoTranslation
}

func (c *CategoryTranslation) String() string {
	return c.Name
}

func (c *CategoryTranslation) ToJSON() string {
	return model.ModelToJson(c)
}

func (c *CategoryTranslation) IsValid() *model.AppError {
	outer := model.CreateAppErrorForModel(
		"model.category_translation.is_valid.%s.app_error",
		"category_translation_id=",
		"CategoryTranslation.IsValid")

	if !model.IsValidId(c.Id) {
		return outer("id", nil)
	}
	if !model.IsValidId(c.CategoryID) {
		return outer("category_id", &c.Id)
	}
	if utf8.RuneCountInString(c.Name) > CATEGORY_NAME_MAX_LENGTH {
		return outer("name", &c.Id)
	}

	return nil
}

func (c *CategoryTranslation) PreSave() {
	if c.Id == "" {
		c.Id = model.NewId()
	}
	c.Name = model.SanitizeUnicode(c.Name)
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
