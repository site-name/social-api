package model

import (
	"strings"
	"time"
	"unicode/utf8"

	"github.com/Masterminds/squirrel"
	"github.com/gosimple/slug"
	"github.com/sitename/sitename/store/store_iface"
	"golang.org/x/text/language"
)

// max lengths for some fields
const (
	ATTRIBUTE_VALUE_NAME_MAX_LENGTH         = 250
	ATTRIBUTE_VALUE_VALUE_MAX_LENGTH        = 9
	ATTRIBUTE_VALUE_SLUG_MAX_LENGTH         = 255
	ATTRIBUTE_VALUE_CONTENT_TYPE_MAX_LENGTH = 50
)

type AttributeValue struct {
	Id          string          `json:"id"`
	Name        string          `json:"name"`
	Value       string          `json:"value"`
	Slug        string          `json:"slug"` // unique
	FileUrl     *string         `json:"file_url"`
	ContentType *string         `json:"content_file"`
	AttributeID string          `json:"attribute_id"`
	RichText    StringInterface `json:"rich_text"`
	Boolean     *bool           `json:"boolean"`
	Datetime    *time.Time      `json:"date_time"`
	Sortable

	Attribute *Attribute `db:"-" json:"-"`
}

type AttributeValueFilterOptions struct {
	Id          squirrel.Sqlizer
	AttributeID squirrel.Sqlizer

	Extra                  squirrel.Sqlizer
	SelectRelatedAttribute bool

	Transaction     store_iface.SqlxTxExecutor
	OrderBy         string
	SelectForUpdate bool // is true, add `FOR UPDATE` suffic to the end of sql query

	Limit int
}

type AttributeValues []*AttributeValue

func (a AttributeValues) IDs() []string {
	var res []string
	meetMap := map[string]bool{}
	for _, item := range a {
		if _, met := meetMap[item.Id]; !met {
			res = append(res, item.Id)
			meetMap[item.Id] = true
		}
	}

	return res
}

func (a AttributeValues) DeepCopy() AttributeValues {
	if a == nil {
		return nil
	}

	res := AttributeValues{}
	for _, item := range a {
		res = append(res, item.DeepCopy())
	}

	return res
}

func (a *AttributeValue) String() string {
	return a.Name
}

func (a *AttributeValue) IsValid() *AppError {
	outer := CreateAppErrorForModel(
		"attribute_value.is_valid.%s.app_error",
		"attribute_value_id=",
		"AttributeValue.IsValid",
	)
	if !IsValidId(a.Id) {
		return outer("id", nil)
	}
	if !IsValidId(a.AttributeID) {
		return outer("attribute_id", &a.Id)
	}
	if utf8.RuneCountInString(a.Name) > ATTRIBUTE_VALUE_NAME_MAX_LENGTH {
		return outer("name", &a.Id)
	}
	if utf8.RuneCountInString(a.Value) > ATTRIBUTE_VALUE_VALUE_MAX_LENGTH {
		return outer("value", &a.Id)
	}
	if len(a.Slug) > ATTRIBUTE_VALUE_SLUG_MAX_LENGTH {
		return outer("slug", &a.Id)
	}
	if a.ContentType != nil && len(*a.ContentType) > ATTRIBUTE_VALUE_CONTENT_TYPE_MAX_LENGTH {
		return outer("content_type", &a.Id)
	}
	if a.Datetime != nil && a.Datetime.IsZero() {
		return outer("date_time", &a.Id)
	}

	return nil
}

func (a *AttributeValue) PreSave() {
	if a.Id == "" {
		a.Id = NewId()
	}
	a.Name = SanitizeUnicode(a.Name)
	a.Slug = slug.Make(a.Name)
}

func (a *AttributeValue) PreUpdate() {
	a.Name = SanitizeUnicode(a.Name)
}

func (a *AttributeValue) ToJSON() string {
	return ModelToJson(a)
}

func (a *AttributeValue) DeepCopy() *AttributeValue {
	res := *a

	if a.RichText != nil {
		res.RichText = a.RichText.DeepCopy()
	}

	if a.Attribute != nil {
		res.Attribute = a.Attribute.DeepCopy()
	}

	return &res
}

// max lengths for some fields of attribute value translation
const (
	ATTRIBUTE_VALUE_TRANSLATION_NAME_MAX_LENGTH = 100
)

// LanguageCode unique together AttributeValueID
type AttributeValueTranslation struct {
	Id               string          `json:"id"`
	LanguageCode     string          `json:"language_code"`
	AttributeValueID string          `json:"attribute_value"`
	Name             string          `json:"name"`
	RichText         StringInterface `json:"rich_text"`
}

func (a *AttributeValueTranslation) String() string {
	return a.Name
}

func (a *AttributeValueTranslation) IsValid() *AppError {
	outer := CreateAppErrorForModel(
		"attribute_value_translation.is_valid.%s.app_error",
		"attribute_value_translation_id=",
		"AttributeValueTranslation.IsValid",
	)
	if !IsValidId(a.Id) {
		return outer("id", nil)
	}
	if !IsValidId(a.AttributeValueID) {
		return outer("attribute_value_id", &a.Id)
	}
	if utf8.RuneCountInString(a.Name) > ATTRIBUTE_VALUE_TRANSLATION_NAME_MAX_LENGTH {
		return outer("name", &a.Id)
	}
	if tag, err := language.Parse(a.LanguageCode); err != nil || !strings.EqualFold(tag.String(), a.LanguageCode) {
		return outer("language_code", &a.Id)
	}

	return nil
}

func (a *AttributeValueTranslation) PreSave() {
	if a.Id == "" {
		a.Id = NewId()
	}
	a.Name = SanitizeUnicode(a.Name)
}

func (a *AttributeValueTranslation) PreUpdate() {
	a.Name = SanitizeUnicode(a.Name)
}

func (a *AttributeValueTranslation) ToJSON() string {
	return ModelToJson(a)
}
