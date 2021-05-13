package product_and_discount

import (
	"io"
	"unicode/utf8"

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
	Name               string  `json:"name"`
	Slug               string  `json:"slug"`
	Description        *string `json:"description"`
	ParentID           string  `json:"parent_id"`
	BackgroundImage    *string `json:"background_image"`
	BackgroundImageAlt string  `json:"background_image_alt"`
	seo.Seo
	model.ModelMetadata
}

func (c *Category) String() string {
	return c.Name
}

func (c *Category) IsValid() *model.AppError {
	outer := model.CreateAppErrorForModel(
		"model.category.is_valid.%s.app_error",
		"category_id=",
		"Category.IsValid")

	if !model.IsValidId(c.Id) {
		return outer("id", nil)
	}
	if !model.IsValidId(c.ParentID) {
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
	c.Slug = slug.Make(c.Slug)
}

func (c *Category) ToJson() string {
	return model.ModelToJson(c)
}

func CategoryFromJSON(data io.Reader) *Category {
	var c Category
	model.ModelFromJson(&c, data)
	return &c
}

type CategoryTranslation struct {
	Id           string  `json:"id"`
	LanguageCode string  `json:"language_code"`
	CategoryID   string  `json:"category_id"`
	Name         string  `json:"name"`
	Description  *string `json:"description"`
	*seo.SeoTranslation
}

func (c *CategoryTranslation) String() string {
	return c.Name
}

func (c *CategoryTranslation) ToJson() string {
	return model.ModelToJson(c)
}

func CategoryTranslationFromJSON(data io.Reader) *Category {
	var c Category
	model.ModelFromJson(&c, data)
	return &c
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
