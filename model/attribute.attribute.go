package model

import (
	"strings"
	"unicode/utf8"

	"github.com/Masterminds/squirrel"
	"github.com/gosimple/slug"
	"github.com/samber/lo"
	"github.com/sitename/sitename/modules/measurement"
	"github.com/sitename/sitename/modules/util"
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
type AttributeType string

const (
	PRODUCT_TYPE AttributeType = "product_type"
	PAGE_TYPE    AttributeType = "page_type"
)

func (e AttributeType) IsValid() bool {
	return e == PRODUCT_TYPE || e == PAGE_TYPE
}

// choices for attribute's entity type
type AttributeEntityType string

func (t AttributeEntityType) IsValid() bool {
	return t == PAGE || t == PRODUCT
}

const (
	PAGE    AttributeEntityType = "page"
	PRODUCT AttributeEntityType = "product"
)

type AttributeInputType string

func (e AttributeInputType) IsValid() bool {
	switch e {
	case AttributeInputTypeDropDown, AttributeInputTypeMultiSelect, AttributeInputTypeFile, AttributeInputTypeReference,
		AttributeInputTypeNumeric, AttributeInputTypeRichText, AttributeInputTypeSwatch, AttributeInputTypeBoolean, AttributeInputTypeDate, AttributeInputTypeDateTime:
		return true
	default:
		return false
	}
}

const (
	AttributeInputTypeDropDown    AttributeInputType = "dropdown"    //
	AttributeInputTypeMultiSelect AttributeInputType = "multiselect" //
	AttributeInputTypeFile        AttributeInputType = "file"        //
	AttributeInputTypeReference   AttributeInputType = "reference"   //
	AttributeInputTypeNumeric     AttributeInputType = "numeric"     //
	AttributeInputTypeRichText    AttributeInputType = "rich_text"   //
	AttributeInputTypeSwatch      AttributeInputType = "swatch"      //
	AttributeInputTypeBoolean     AttributeInputType = "boolean"     //
	AttributeInputTypeDate        AttributeInputType = "date"        //
	AttributeInputTypeDateTime    AttributeInputType = "date_time"   //
)

var (
	ALLOWED_IN_VARIANT_SELECTION = util.AnyArray[AttributeInputType]{AttributeInputTypeDropDown, AttributeInputTypeBoolean, AttributeInputTypeSwatch, AttributeInputTypeNumeric}
	TYPES_WITH_CHOICES           = util.AnyArray[AttributeInputType]{AttributeInputTypeDropDown, AttributeInputTypeMultiSelect, AttributeInputTypeSwatch}
	TYPES_WITH_UNIQUE_VALUES     = util.AnyArray[AttributeInputType]{AttributeInputTypeFile, AttributeInputTypeReference, AttributeInputTypeRichText, AttributeInputTypeNumeric, AttributeInputTypeDate, AttributeInputTypeDateTime} // list of the translatable attributes, excluding attributes with choices.
	TRANSLATABLE_ATTRIBUTES      = util.AnyArray[AttributeInputType]{AttributeInputTypeRichText}
)

var ATTRIBUTE_PROPERTIES_CONFIGURATION = map[string][]AttributeInputType{
	"filterable_in_storefront": {
		AttributeInputTypeDropDown,
		AttributeInputTypeMultiSelect,
		AttributeInputTypeNumeric,
		AttributeInputTypeSwatch,
		AttributeInputTypeBoolean,
		AttributeInputTypeDate,
		AttributeInputTypeDateTime,
	},
	"filterable_in_dashboard": {
		AttributeInputTypeDropDown,
		AttributeInputTypeMultiSelect,
		AttributeInputTypeNumeric,
		AttributeInputTypeSwatch,
		AttributeInputTypeBoolean,
		AttributeInputTypeDate,
		AttributeInputTypeDateTime,
	},
	"available_in_grid": {
		AttributeInputTypeDropDown,
		AttributeInputTypeMultiSelect,
		AttributeInputTypeNumeric,
		AttributeInputTypeSwatch,
		AttributeInputTypeBoolean,
		AttributeInputTypeDate,
		AttributeInputTypeDateTime,
	},
	"storefront_search_position": {
		AttributeInputTypeDropDown,
		AttributeInputTypeMultiSelect,
		AttributeInputTypeBoolean,
		AttributeInputTypeDate,
		AttributeInputTypeDateTime,
	},
}

// ORDER BY Slug
type Attribute struct {
	Id                       string               `json:"id"`
	Slug                     string               `json:"slug"` // unique
	Name                     string               `json:"name"`
	Type                     AttributeType        `json:"type"`
	InputType                AttributeInputType   `json:"input_type"`
	EntityType               *AttributeEntityType `json:"entity_type"`
	Unit                     *string              `json:"unit"` // lower cased
	ValueRequired            bool                 `json:"value_required"`
	IsVariantOnly            bool                 `json:"is_variant_only"`
	VisibleInStoreFront      bool                 `json:"visible_in_storefront"`
	FilterableInStorefront   bool                 `json:"filterable_in_storefront"`
	FilterableInDashboard    bool                 `json:"filterable_in_dashboard"`
	StorefrontSearchPosition int                  `json:"storefront_search_position"`
	AvailableInGrid          bool                 `json:"available_in_grid"`
	ModelMetadata

	attributeValues AttributeValues `db:"-"`
}

func (a *Attribute) SetAttributeValues(v AttributeValues) {
	a.attributeValues = v
}

func (a *Attribute) GetAttributeValues() AttributeValues {
	return a.attributeValues
}

type AttributeFilterOption struct {
	Id                     squirrel.Sqlizer
	Slug                   squirrel.Sqlizer
	InputType              squirrel.Sqlizer
	ProductTypes           squirrel.Sqlizer // INNER JOIN AttributeProducts ON ... WHERE AttributeProducts.ProductTypeID ...
	ProductVariantTypes    squirrel.Sqlizer // INNER JOIN AttributeVariants ON ... WHERE AttributeVariants.ProductTypeID ...
	Type                   squirrel.Sqlizer
	VisibleInStoreFront    *bool
	ValueRequired          *bool
	IsVariantOnly          *bool
	FilterableInStorefront *bool
	FilterableInDashboard  *bool
	AvailableInGrid        *bool
	Metadata               StringMAP
	Search                 *string // Slug or Name ILIKE ...
	InCollection           *string
	InCategory             *string
	Channel                *string // channel id or slug in which attributes reside

	OrderBy                        string
	Distinct                       bool
	UserHasOneOfProductPermissions *bool

	Extra squirrel.Sqlizer

	PrefetchRelatedAttributeValues bool
}

func (a *Attribute) IsValid() *AppError {
	outer := CreateAppErrorForModel(
		"model.attribute.is_valid.%s.app_error",
		"attribute_id=",
		"Attribute.IsValid",
	)
	if !IsValidId(a.Id) {
		return outer("id", &a.Id)
	}
	if utf8.RuneCountInString(a.Name) > ATTRIBUTE_NAME_MAX_LENGTH {
		return outer("name", &a.Id)
	}
	if len(a.Slug) > ATTRIBUTE_SLUG_MAX_LENGTH {
		return outer("slug", &a.Id)
	}
	if len(a.Type) > ATTRIBUTE_TYPE_MAX_LENGTH || !a.Type.IsValid() {
		return outer("type", &a.Id)
	}
	if len(a.InputType) > ATTRIBUTE_TYPE_MAX_LENGTH || !a.InputType.IsValid() {
		return outer("input_type", &a.Id)
	}
	if a.EntityType != nil &&
		(len(*a.EntityType) > ATTRIBUTE_ENTITY_TYPE_MAX_LENGTH || !(*a.EntityType).IsValid()) {
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
	return lo.Map(as, func(a *Attribute, _ int) string { return a.Id })
}

func (a *Attribute) PreSave() {
	if a.Id == "" {
		a.Id = NewId()
	}
	if a.InputType == "" {
		a.InputType = AttributeInputTypeDropDown
	}
	a.Name = SanitizeUnicode(a.Name)
	if a.Slug == "" {
		a.Slug = slug.Make(a.Name)
	}
}

func (a *Attribute) PreUpdate() {
	a.Name = SanitizeUnicode(a.Name)
}

func (a *Attribute) String() string {
	return a.Name
}

func (a *Attribute) ToJSON() string {
	return ModelToJson(a)
}

func (a *Attribute) DeepCopy() *Attribute {
	if a == nil {
		return nil
	}
	res := *a

	if a.EntityType != nil {
		res.EntityType = NewPrimitive(*a.EntityType)
	}
	if a.Unit != nil {
		res.Unit = NewPrimitive(*a.Unit)
	}
	res.ModelMetadata = a.ModelMetadata.DeepCopy()
	res.attributeValues = a.attributeValues.DeepCopy()

	return &res
}

// max lengths for attribute translation's fields
const (
	ATTRIBUTE_TRANSLATION_NAME_MAX_LENGTH = 100
)

type AttributeTranslation struct {
	Id           string           `json:"id"`
	AttributeID  string           `json:"attribute_id"`
	LanguageCode LanguageCodeEnum `json:"language_code"`
	Name         string           `json:"name"`
}

func (a *AttributeTranslation) IsValid() *AppError {
	outer := CreateAppErrorForModel(
		"model.attribute_translation.is_valid.%s.app_error",
		"attribute_translation_id=",
		"AttributeTranslation.IsValid",
	)
	if !IsValidId(a.Id) {
		return outer("id", nil)
	}
	if !IsValidId(a.AttributeID) {
		return outer("attribute_id", nil)
	}
	if utf8.RuneCountInString(a.Name) > ATTRIBUTE_TRANSLATION_NAME_MAX_LENGTH {
		return outer("name", &a.Id)
	}
	if !a.LanguageCode.IsValid() {
		return outer("language_code", &a.Id)
	}

	return nil
}

func (a *AttributeTranslation) PreSave() {
	if a.Id == "" {
		a.Id = NewId()
	}
	a.Name = SanitizeUnicode(a.Name)
}

func (a *AttributeTranslation) PreUpdate() {
	a.Name = SanitizeUnicode(a.Name)
}

func (a *AttributeTranslation) String() string {
	return a.Name
}
