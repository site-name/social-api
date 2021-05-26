package attribute

import (
	"io"
	"strings"
	"unicode/utf8"

	"github.com/gosimple/slug"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model/page"
	"github.com/sitename/sitename/model/product_and_discount"
	"golang.org/x/text/language"
)

// max lengths for some fields
const (
	ATTRIBUTE_SLUG_MAX_LENGTH        = 255
	ATTRIBUTE_NAME_MAX_LENGTH        = 250
	ATTRIBUTE_TYPE_MAX_LENGTH        = 50
	ATTRIBUTE_UNIT_MAX_LENGTH        = 100
	ATTRIBUTE_INPUT_TYPE_MAX_LENGTH  = 50
	ATTRIBUTE_ENTITY_TYPE_MAX_LENGTH = 50
)

// choices for attribute's type
const (
	PRODUCT_TYPE = "product-type"
	PAGE_TYPE    = "page-type"
)

// choices for attribute's input type
const (
	DROPDOWN    = "dropdown"
	MULTISELECT = "multiselect"
	FILE        = "file"
	REFERENCE   = "reference"
	NUMERIC     = "numeric"
	RICH_TEXT   = "rich-text"
)

var AttributeInputTypeStrings = map[string]string{
	DROPDOWN:    "Dropdown",
	MULTISELECT: "Multi Select",
	FILE:        "File",
	REFERENCE:   "Reference",
	NUMERIC:     "Numeric",
	RICH_TEXT:   "Rich Text",
}

var AttributeTypeStrings = map[string]string{
	PRODUCT_TYPE: "Product type",
	PAGE_TYPE:    "Page type",
}

type Attribute struct {
	Id                       string                              `json:"id"`
	Slug                     string                              `json:"slug"` // unique
	Name                     string                              `json:"name"`
	Type                     string                              `json:"type"`
	InputType                string                              `json:"input_type"`
	EntityType               *string                             `json:"entity_type"`
	ProductTypes             []*product_and_discount.ProductType `json:"product_types" db:"-"`
	ProductVariantTypes      []*product_and_discount.ProductType `json:"product_variant_types" db:"-"`
	PageTypes                []*page.PageType                    `json:"page_types" db:"-"`
	Unit                     *string                             `json:"unit"`
	ValueRequired            *bool                               `json:"value_required"`
	IsVariantOnly            *bool                               `json:"is_variant_only"`
	VisibleInStoreFront      *bool                               `json:"visible_in_storefront"`
	FilterableInStorefront   *bool                               `json:"filterable_in_storefront"`
	FilterableInDashboard    *bool                               `json:"filterable_in_dashboard"`
	StorefrontSearchPosition int                                 `json:"storefront_search_position"`
	AvailableInGrid          *bool                               `json:"available_in_grid"`
	model.ModelMetadata
}

func (a *Attribute) IsValid() *model.AppError {
	outer := model.CreateAppErrorForModel(
		"model.attribute.is_valid.%s.app_error",
		"attribute_id=",
		"Attribute.IsValid",
	)
	if !model.IsValidId(a.Id) {
		return outer("id", &a.Id)
	}
	if utf8.RuneCountInString(a.Name) > ATTRIBUTE_NAME_MAX_LENGTH {
		return outer("name", &a.Id)
	}
	if len(a.Slug) > ATTRIBUTE_SLUG_MAX_LENGTH {
		return outer("slug", &a.Id)
	}
	if len(a.Type) > ATTRIBUTE_TYPE_MAX_LENGTH {
		return outer("type", &a.Id)
	}
	if len(a.InputType) > ATTRIBUTE_TYPE_MAX_LENGTH || AttributeInputTypeStrings[strings.ToLower(a.InputType)] == "" {
		return outer("input_type", &a.Id)
	}
	if a.EntityType != nil && len(*a.EntityType) > ATTRIBUTE_ENTITY_TYPE_MAX_LENGTH {
		return outer("entity_type", &a.Id)
	}
	if a.Unit != nil && len(*a.Unit) > ATTRIBUTE_UNIT_MAX_LENGTH {
		return outer("unit", &a.Id)
	}

	return nil
}

func (a *Attribute) PreSave() {
	if a.Id == "" {
		a.Id = model.NewId()
	}
	if a.InputType == "" {
		a.InputType = DROPDOWN
	}
	if a.ValueRequired == nil {
		a.ValueRequired = model.NewBool(false)
	}
	if a.IsVariantOnly == nil {
		a.IsVariantOnly = model.NewBool(false)
	}
	if a.VisibleInStoreFront == nil {
		a.VisibleInStoreFront = model.NewBool(true)
	}
	if a.FilterableInStorefront == nil {
		a.FilterableInStorefront = model.NewBool(false)
	}
	if a.FilterableInDashboard == nil {
		a.FilterableInDashboard = model.NewBool(false)
	}
	if a.AvailableInGrid == nil {
		a.AvailableInGrid = model.NewBool(false)
	}
	a.Name = model.SanitizeUnicode(a.Name)
	a.Slug = slug.Make(a.Name)
}

func (a *Attribute) PreUpdate() {
	a.Name = model.SanitizeUnicode(a.Name)
	a.Slug = slug.Make(a.Name)
}

func (a *Attribute) String() string {
	return a.Name
}

func (a *Attribute) ToJson() string {
	return model.ModelToJson(a)
}

func AttributeFromJson(data io.Reader) *Attribute {
	var a Attribute
	model.ModelFromJson(&a, data)
	return &a
}

// -----------------

const (
	ATTRIBUTE_TRANSLATION_NAME_MAX_LENGTH = 100
)

type AttributeTranslation struct {
	Id           string `json:"id"`
	AttributeId  string `json:"attribute_id"`
	LanguageCode string `json:"language_code"`
	Name         string `json:"name"`
}

func (a *AttributeTranslation) IsValid() *model.AppError {
	outer := model.CreateAppErrorForModel(
		"model.attribute_translation.is_valid.%s.app_error",
		"attribute_translation_id=",
		"AttributeTranslation.IsValid",
	)
	if !model.IsValidId(a.Id) {
		return outer("id", nil)
	}
	if !model.IsValidId(a.AttributeId) {
		return outer("attribute_id", nil)
	}
	if utf8.RuneCountInString(a.Name) > ATTRIBUTE_TRANSLATION_NAME_MAX_LENGTH {
		return outer("name", &a.Id)
	}
	if tag, err := language.Parse(a.LanguageCode); err != nil || !strings.EqualFold(tag.String(), a.LanguageCode) {
		return outer("language_code", &a.Id)
	}

	return nil
}

func (a *AttributeTranslation) PreSave() {
	if a.Id == "" {
		a.Id = model.NewId()
	}
	a.Name = model.SanitizeUnicode(a.Name)
}

func (a *AttributeTranslation) PreUpdate() {
	a.Name = model.SanitizeUnicode(a.Name)
}

func (a *AttributeTranslation) String() string {
	return a.Name
}
