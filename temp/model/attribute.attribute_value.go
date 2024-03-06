package model

import (
	"net/http"
	"time"

	"github.com/gosimple/slug"
	"github.com/mattermost/squirrel"
	"github.com/samber/lo"
	"gorm.io/gorm"
)

// constants related to attribute value model
const (
	AttributeValueNameMaxLength = 250
)

type AttributeValue struct {
	Id          string          `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid();column:Id"`
	Name        string          `json:"name" gorm:"type:varchar(250);column:Name"`                                      // varchar(250)
	Value       string          `json:"value" gorm:"type:varchar(9);column:Value"`                                      // varchar(9)
	Slug        string          `json:"slug" gorm:"uniqueIndex:idx_slug_attributeid_key;type:varchar(255);column:Slug"` // unique with attribute_id; varchar(255)
	AttributeID string          `json:"attribute_id" gorm:"uniqueIndex:idx_slug_attributeid_key;type:uuid;column:AttributeID"`
	FileUrl     *string         `json:"file_url" gorm:"type:varchar(500);column:FileUrl"`        // varchar(500)
	ContentType *string         `json:"content_file" gorm:"type:varchar(50);column:ContentType"` // varchar(50)
	RichText    StringInterface `json:"rich_text" gorm:"column:RichText"`
	Boolean     *bool           `json:"boolean" gorm:"column:Boolean"`
	Datetime    *time.Time      `json:"date_time" gorm:"column:Datetime"`
	Sortable

	VariantsAssignments     []*AssignedVariantAttribute      `json:"-" gorm:"many2many:AssignedVariantAttributeValues"`
	PageAssignments         []*AssignedPageAttribute         `json:"-" gorm:"many2many:AssignedPageAttributeValues"`
	ProductAssignments      []*AssignedProductAttribute      `json:"-" gorm:"many2many:AssignedProductAttributeValues"`
	Attribute               *Attribute                       `json:"-" gorm:"foreignKey:AttributeID;constraint:OnDelete:CASCADE;"`
	VariantValueAssignment  []*AssignedVariantAttributeValue `json:"-" gorm:"foreignKey:ValueID;constraint:OnDelete:CASCADE;"`
	PageValueAssignment     []*AssignedPageAttributeValue    `json:"-" gorm:"foreignKey:ValueID;constraint:OnDelete:CASCADE;"`
	ProductValueAssignments []*AssignedProductAttributeValue `json:"-" gorm:"foreignKey:ValueID;constraint:OnDelete:CASCADE;"`
}

// column names of attribute value table
const (
	AttributeValueColumnId          = "Id"
	AttributeValueColumnName        = "Name"
	AttributeValueColumnValue       = "Value"
	AttributeValueColumnSlug        = "Slug"
	AttributeValueColumnAttributeID = "AttributeID"
	AttributeValueColumnFileUrl     = "FileUrl"
	AttributeValueColumnContentType = "ContentType"
	AttributeValueColumnRichText    = "RichText"
	AttributeValueColumnBoolean     = "Boolean"
	AttributeValueColumnDatetime    = "Datetime"
)

func (a *AttributeValue) BeforeCreate(_ *gorm.DB) error { a.PreSave(); return a.IsValid() }
func (a *AttributeValue) BeforeUpdate(_ *gorm.DB) error { a.PreUpdate(); return a.IsValid() }
func (a *AttributeValue) TableName() string             { return AttributeValueTableName }
func (a *AttributeValue) PreSave()                      { a.commonPre(); a.Slug = slug.Make(a.Name) }
func (a *AttributeValue) commonPre()                    { a.Name = SanitizeUnicode(a.Name) }
func (a *AttributeValue) PreUpdate()                    { a.commonPre() }

type AttributeValueFilterOptions struct {
	Conditions             squirrel.Sqlizer
	SelectRelatedAttribute bool

	Transaction *gorm.DB // if provided, this will be responsible for perform queries
	// set to true to add `FOR UPDATE` suffix to the end of sql queries.
	//
	// NOTE: only apply when Transaction field is provided
	SelectForUpdate bool

	Ordering string
}

type AttributeValues []*AttributeValue

func (a AttributeValues) IDs() []string {
	return lo.Map(a, func(v *AttributeValue, _ int) string { return v.Id })
}

func (a AttributeValues) DeepCopy() AttributeValues {
	return lo.Map(a, func(v *AttributeValue, _ int) *AttributeValue { return v.DeepCopy() })
}

func (a *AttributeValue) String() string { return a.Name }
func (a *AttributeValue) IsValid() *AppError {
	if !IsValidId(a.AttributeID) {
		return NewAppError("AttributeValue.IsValid", "model.attribute_value.is_valid.attribute_id.app_error", nil, "", http.StatusBadRequest)
	}
	if a.Datetime != nil && a.Datetime.IsZero() {
		return NewAppError("AttributeValue.IsValid", "model.attribute_value.is_valid.date_time.app_error", nil, "", http.StatusBadRequest)
	}

	return nil
}

func (a *AttributeValue) DeepCopy() *AttributeValue {
	res := *a

	if a.RichText != nil {
		res.RichText = a.RichText.DeepCopy()
	}
	if a.FileUrl != nil {
		res.FileUrl = GetPointerOfValue(*a.FileUrl)
	}
	if a.Boolean != nil {
		res.Boolean = GetPointerOfValue(*a.Boolean)
	}
	if a.Datetime != nil {
		res.Datetime = GetPointerOfValue(*a.Datetime)
	}

	return &res
}

// LanguageCode unique together AttributeValueID
type AttributeValueTranslation struct {
	Id               string           `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid();column:Id"`
	LanguageCode     LanguageCodeEnum `json:"language_code" gorm:"type:varchar(35);index::languagecode_attributevalueid_key;column:LanguageCode"` // varchar(35); unique together with attributeid
	AttributeValueID string           `json:"attribute_value" gorm:"type:uuid;index::languagecode_attributevalueid_key;column:AttributeValueID"`
	Name             string           `json:"name" gorm:"type:varchar(100);column:Name"` // varchar(100)
	RichText         StringInterface  `json:"rich_text" gorm:"column:RichText"`
}

func (a *AttributeValueTranslation) BeforeCreate(_ *gorm.DB) error { a.commonPre(); return a.IsValid() }
func (a *AttributeValueTranslation) BeforeUpdate(_ *gorm.DB) error { a.commonPre(); return a.IsValid() }
func (a *AttributeValueTranslation) TableName() string             { return AttributeValueTranslationTableName }
func (a *AttributeValueTranslation) String() string                { return a.Name }
func (a *AttributeValueTranslation) commonPre()                    { a.Name = SanitizeUnicode(a.Name) }

func (a *AttributeValueTranslation) IsValid() *AppError {
	if !IsValidId(a.AttributeValueID) {
		return NewAppError("AttributeValueTranslation.IsValid", "model.attribute_value_translation.is_valid.attribute_value_id.app_error", nil, "", http.StatusBadRequest)
	}
	if !a.LanguageCode.IsValid() {
		return NewAppError("AttributeValueTranslation.IsValid", "model.attribute_value_translation.is_valid.language_code.app_error", nil, "", http.StatusBadRequest)
	}
	return nil
}
