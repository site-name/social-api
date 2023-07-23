package model

import (
	"strings"

	"github.com/Masterminds/squirrel"
	"github.com/gosimple/slug"
	"github.com/samber/lo"
	"github.com/sitename/sitename/modules/measurement"
	"github.com/sitename/sitename/modules/util"
	"gorm.io/gorm"
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
	Id                       string               `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid();column:Id"`
	Slug                     string               `json:"slug" gorm:"type:varchar(255);uniqueIndex:slug_unique_key;column:Slug"` // varchar(255); unique
	Name                     string               `json:"name" gorm:"type:varchar(250);column:Name"`                             // varchar(250)
	Type                     AttributeType        `json:"type" gorm:"type:varchar(50);column:Type"`                              // varchar(50)
	InputType                AttributeInputType   `json:"input_type" gorm:"type:varchar(50);column:InputType"`                   // default "dropdown"
	EntityType               *AttributeEntityType `json:"entity_type" gorm:"type:varchar(50);column:EntityType"`
	Unit                     *string              `json:"unit" gorm:"type:varchar(100);column:Unit"` // lower cased
	ValueRequired            bool                 `json:"value_required" gorm:"column:ValueRequired"`
	IsVariantOnly            bool                 `json:"is_variant_only" gorm:"column:IsVariantOnly"`
	VisibleInStoreFront      bool                 `json:"visible_in_storefront" gorm:"column:VisibleInStoreFront"`
	FilterableInStorefront   bool                 `json:"filterable_in_storefront" gorm:"column:FilterableInStorefront"`
	FilterableInDashboard    bool                 `json:"filterable_in_dashboard" gorm:"column:FilterableInDashboard"`
	StorefrontSearchPosition int                  `json:"storefront_search_position" gorm:"column:StorefrontSearchPosition"`
	AvailableInGrid          bool                 `json:"available_in_grid" gorm:"column:AvailableInGrid"`
	ModelMetadata

	AttributeValues AttributeValues  `json:"-" gorm:"foreignKey:AttributeID;constraint:OnDelete:CASCADE;"`
	AttributePages  []*AttributePage `json:"-" gorm:"foreignKey:AttributeID;constraint:OnDelete:CASCADE;"`
}

func (a *Attribute) BeforeCreate(_ *gorm.DB) error { a.PreSave(); return a.IsValid() }
func (a *Attribute) BeforeUpdate(_ *gorm.DB) error { a.PreUpdate(); return a.IsValid() }
func (a *Attribute) TableName() string             { return AttributeTableName }
func (a *Attribute) PreSave() {
	a.commonPre()
	if a.Slug == "" {
		a.Slug = slug.Make(a.Name)
	}
}

func (a *Attribute) commonPre() {
	a.Name = SanitizeUnicode(a.Name)
	if !a.InputType.IsValid() {
		a.InputType = AttributeInputTypeDropDown
	}
}

func (a *Attribute) PreUpdate()     { a.commonPre() }
func (a *Attribute) String() string { return a.Name }

type AttributeFilterOption struct {
	Conditions squirrel.Sqlizer

	AttributeProduct_ProductTypeID squirrel.Sqlizer // INNER JOIN AttributeProducts ON ... WHERE AttributeProducts.ProductTypeID ...
	AttributeVariant_ProductTypeID squirrel.Sqlizer // INNER JOIN AttributeVariants ON ... WHERE AttributeVariants.ProductTypeID ...

	Metadata     StringMAP
	Search       string // Slug or Name ILIKE ...
	InCollection *string
	InCategory   *string
	Channel      *string // channel id or slug in which attributes reside

	OrderBy         string
	Distinct        bool
	UserIsShopStaff bool // user has role 'shop_admin' or 'shop_staff'
	Limit           int

	PrefetchRelatedAttributeValues bool
}

func (a *Attribute) IsValid() *AppError {
	outer := CreateAppErrorForModel(
		"model.attribute.is_valid.%s.app_error",
		"attribute_id=",
		"Attribute.IsValid",
	)
	if !a.Type.IsValid() {
		return outer("type", &a.Id)
	}
	if !a.InputType.IsValid() {
		return outer("input_type", &a.Id)
	}
	if a.EntityType != nil && !(*a.EntityType).IsValid() {
		return outer("entity_type", &a.Id)
	}
	if a.Unit != nil && measurement.MeasurementUnitMap[strings.ToUpper(*a.Unit)] == "" {
		return outer("unit", &a.Id)
	}

	return nil
}

type Attributes []*Attribute

func (as Attributes) IDs() []string {
	return lo.Map(as, func(a *Attribute, _ int) string { return a.Id })
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
	res.AttributeValues = a.AttributeValues.DeepCopy()

	return &res
}

type AttributeTranslation struct {
	Id           string           `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid();column:Id"`
	AttributeID  string           `json:"attribute_id" gorm:"type:uuid;column:AttributeID"`
	LanguageCode LanguageCodeEnum `json:"language_code" gorm:"type:varchar(35);column:LanguageCode"`
	Name         string           `json:"name" gorm:"type:varchar(100);column:Name"`
}

func (a *AttributeTranslation) BeforeCreate(_ *gorm.DB) error { a.commonPre(); return a.IsValid() }
func (a *AttributeTranslation) BeforeUpdate(_ *gorm.DB) error { a.commonPre(); return a.IsValid() }
func (a *AttributeTranslation) TableName() string             { return AttributeTranslationTableName }

func (a *AttributeTranslation) IsValid() *AppError {
	outer := CreateAppErrorForModel(
		"model.attribute_translation.is_valid.%s.app_error",
		"attribute_translation_id=",
		"AttributeTranslation.IsValid",
	)
	if !IsValidId(a.AttributeID) {
		return outer("attribute_id", nil)
	}
	if !a.LanguageCode.IsValid() {
		return outer("language_code", &a.Id)
	}

	return nil
}

func (a *AttributeTranslation) commonPre() {
	a.Name = SanitizeUnicode(a.Name)
}

func (a *AttributeTranslation) String() string {
	return a.Name
}
