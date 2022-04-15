package attribute

import (
	"strings"
	"unicode/utf8"

	"github.com/Masterminds/squirrel"
	"github.com/gosimple/slug"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/modules/measurement"
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
	PRODUCT_TYPE = "product_type"
	PAGE_TYPE    = "page_type"
)

// choices for attribute's entity type
const (
	PAGE    = "page"
	PRODUCT = "product"
)

// choices for attribute's input type
const (
	DROPDOWN    string = "dropdown"
	MULTISELECT string = "multiselect"
	FILE        string = "file"
	REFERENCE   string = "reference"
	NUMERIC     string = "numeric"
	RICH_TEXT   string = "rich_text"
	SWATCH      string = "swatch"
	BOOLEAN     string = "boolean"
	DATE        string = "date"
	DATE_TIME   string = "date_time"
)

var (
	ALLOWED_IN_VARIANT_SELECTION = model.StringArray{DROPDOWN, BOOLEAN, SWATCH, NUMERIC}
	TYPES_WITH_CHOICES           = model.StringArray{DROPDOWN, MULTISELECT, SWATCH}
	TYPES_WITH_UNIQUE_VALUES     = model.StringArray{FILE, REFERENCE, RICH_TEXT, NUMERIC, DATE, DATE_TIME} // list of the translatable attributes, excluding attributes with choices.
	TRANSLATABLE_ATTRIBUTES      = model.StringArray{RICH_TEXT}
)
var AttributeInputTypeStrings = map[string]string{
	DROPDOWN:    "Dropdown",
	MULTISELECT: "Multi Select",
	FILE:        "File",
	REFERENCE:   "Reference",
	NUMERIC:     "Numeric",
	RICH_TEXT:   "Rich Text",
	BOOLEAN:     "Boolean",
	DATE:        "Date",
	DATE_TIME:   "Date Time",
}

var AttributeTypeStrings = map[string]string{
	PRODUCT_TYPE: "Product type",
	PAGE_TYPE:    "Page type",
}

var AttributeEntityTypeStrings = map[string]string{
	PAGE:    "Page",
	PRODUCT: "Product",
}

// Attribute
type Attribute struct {
	Id                       string  `json:"id"`
	Slug                     string  `json:"slug"` // unique
	Name                     string  `json:"name"`
	Type                     string  `json:"type"`
	InputType                string  `json:"input_type"`
	EntityType               *string `json:"entity_type"`
	Unit                     *string `json:"unit"` // lower cased
	ValueRequired            bool    `json:"value_required"`
	IsVariantOnly            bool    `json:"is_variant_only"`
	VisibleInStoreFront      bool    `json:"visible_in_storefront"`
	FilterableInStorefront   bool    `json:"filterable_in_storefront"`
	FilterableInDashboard    bool    `json:"filterable_in_dashboard"`
	StorefrontSearchPosition int     `json:"storefront_search_position"`
	AvailableInGrid          bool    `json:"available_in_grid"`
	model.ModelMetadata

	AttributeValues AttributeValues `json:"-" db:"-"`
}

type AttributeFilterOption struct {
	Id                  squirrel.Sqlizer
	Slug                squirrel.Sqlizer
	InputType           squirrel.Sqlizer
	ProductTypes        squirrel.Sqlizer // INNER JOIN AttributeProducts ON ... WHERE AttributeProducts.ProductTypeID ...
	ProductVariantTypes squirrel.Sqlizer // INNER JOIN AttributeVariants ON ... WHERE AttributeVariants.ProductTypeID ...
	Distinct            bool
	VisibleInStoreFront *bool

	ValueRequired          *bool
	IsVariantOnly          *bool
	VisibleInStorefront    *bool
	FilterableInStorefront *bool
	FilterableInDashboard  *bool
	AvailableInGrid        *bool
	Type                   string

	OrderBy string

	PrefetchRelatedAttributeValues bool
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
	if len(a.Type) > ATTRIBUTE_TYPE_MAX_LENGTH || AttributeTypeStrings[a.Type] == "" {
		return outer("type", &a.Id)
	}
	if len(a.InputType) > ATTRIBUTE_TYPE_MAX_LENGTH || AttributeInputTypeStrings[a.InputType] == "" {
		return outer("input_type", &a.Id)
	}
	if (a.EntityType != nil && len(*a.EntityType) > ATTRIBUTE_ENTITY_TYPE_MAX_LENGTH) ||
		(a.EntityType != nil && AttributeEntityTypeStrings[*a.EntityType] == "") {
		return outer("entity_type", &a.Id)
	}
	if (a.Unit != nil && len(*a.Unit) > ATTRIBUTE_UNIT_MAX_LENGTH) ||
		(a.Unit != nil && measurement.MeasurementUnitMap[strings.ToUpper(*a.Unit)] == "") {
		return outer("unit", &a.Id)
	}

	return nil
}

type Attributes []*Attribute

func (as Attributes) IDs() []string {
	res := []string{}
	for _, item := range as {
		if item != nil {
			res = append(res, item.Id)
		}
	}

	return res
}

func (a *Attribute) PreSave() {
	if a.Id == "" {
		a.Id = model.NewId()
	}
	if a.InputType == "" {
		a.InputType = DROPDOWN
	}
	a.Name = model.SanitizeUnicode(a.Name)
	a.Slug = slug.Make(a.Name)
}

func (a *Attribute) PreUpdate() {
	a.Name = model.SanitizeUnicode(a.Name)
}

func (a *Attribute) String() string {
	return a.Name
}

func (a *Attribute) ToJSON() string {
	return model.ModelToJson(a)
}

func (a *Attribute) DeepCopy() *Attribute {
	if a == nil {
		return nil
	}
	res := *a

	res.ModelMetadata = *(a.ModelMetadata.DeepCopy())
	res.AttributeValues = a.AttributeValues.DeepCopy()

	return &res
}

// max lengths for attribute translation's fields
const (
	ATTRIBUTE_TRANSLATION_NAME_MAX_LENGTH = 100
)

type AttributeTranslation struct {
	Id           string `json:"id"`
	AttributeID  string `json:"attribute_id"`
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
	if !model.IsValidId(a.AttributeID) {
		return outer("attribute_id", nil)
	}
	if utf8.RuneCountInString(a.Name) > ATTRIBUTE_TRANSLATION_NAME_MAX_LENGTH {
		return outer("name", &a.Id)
	}
	if tag, err := language.Parse(a.LanguageCode); err != nil || !strings.EqualFold(tag.String(), a.LanguageCode) {
		return outer("language_code", &a.Id)
	}
	if model.Languages[strings.ToLower(a.LanguageCode)] == "" {
		return outer("language_coce", &a.Id)
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
