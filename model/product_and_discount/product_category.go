package product_and_discount

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"unicode/utf8"

	"github.com/gosimple/slug"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model/seo"
	"github.com/sitename/sitename/modules/json"
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
	*seo.Seo
	*model.ModelMetadata
}

func (c *Category) String() string {
	return c.Name
}

func (c *Category) createAppError(fieldName string) *model.AppError {
	id := fmt.Sprintf("model.category.is_valid.%s.app_error", fieldName)
	var details string
	if !strings.EqualFold(fieldName, "id") {
		details = "category_id=" + c.Id
	}

	return model.NewAppError("Category.IsValid", id, nil, details, http.StatusBadRequest)
}

func (c *Category) IsValid() *model.AppError {
	if !model.IsValidId(c.Id) {
		return c.createAppError("id")
	}
	if !model.IsValidId(c.ParentID) {
		return c.createAppError("id")
	}
	if len(c.BackgroundImageAlt) > CATEGORY_BG_IMAGE_ALT_MAX_LENGTH {
		return c.createAppError("background_image_alt")
	}
	if utf8.RuneCountInString(c.Name) > CATEGORY_NAME_MAX_LENGTH {
		return c.createAppError("name")
	}
	if utf8.RuneCountInString(c.Slug) > CATEGORY_SLUG_MAX_LENGTH {
		return c.createAppError("slug")
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
	b, _ := json.JSON.Marshal(c)
	return string(b)
}

func CategoryFromJSON(data io.Reader) *Category {
	var c Category
	err := json.JSON.NewDecoder(data).Decode(&c)
	if err != nil {
		return nil
	}
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
	b, _ := json.JSON.Marshal(c)
	return string(b)
}

func CategoryTranslationFromJSON(data io.Reader) *Category {
	var c Category
	err := json.JSON.NewDecoder(data).Decode(&c)
	if err != nil {
		return nil
	}
	return &c
}

func (c *CategoryTranslation) createAppError(fieldName string) *model.AppError {
	id := fmt.Sprintf("model.category_translation.is_valid.%s.app_error", fieldName)
	var details string
	if !strings.EqualFold(fieldName, "id") {
		details = "category_translation_id=" + c.Id
	}

	return model.NewAppError("CategoryTranslation.IsValid", id, nil, details, http.StatusBadRequest)
}

func (c *CategoryTranslation) IsValid() *model.AppError {
	if !model.IsValidId(c.Id) {
		return c.createAppError("id")
	}
	if !model.IsValidId(c.CategoryID) {
		return c.createAppError("category_id")
	}
	if utf8.RuneCountInString(c.Name) > CATEGORY_NAME_MAX_LENGTH {
		return c.createAppError("name")
	}

	return nil
}

func (c *CategoryTranslation) PreSave() {
	if c.Id == "" {
		c.Id = model.NewId()
	}
	c.Name = model.SanitizeUnicode(c.Name)
}
