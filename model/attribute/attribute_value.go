package attribute

import (
	"io"
	"strings"
	"unicode/utf8"

	"github.com/gosimple/slug"
	"github.com/sitename/sitename/model"
	"golang.org/x/text/language"
)

// max lengths for some fields
const (
	ATTRIBUTE_VALUE_NAME_MAX_LENGTH         = 250
	ATTRIBUTE_VALUE_VALUE_MAX_LENGTH        = 100
	ATTRIBUTE_VALUE_SLUG_MAX_LENGTH         = 255
	ATTRIBUTE_VALUE_CONTENT_TYPE_MAX_LENGTH = 50
)

type AttributeValue struct {
	Id          string                 `json:"id"`
	Name        string                 `json:"name"`
	Value       string                 `json:"value"`
	Slug        string                 `json:"slug"`
	FileUrl     *string                `json:"file_url"`
	ContentType *string                `json:"content_file"`
	AttributeID string                 `json:"attribute_id"`
	Attribute   *Attribute             `json:"attribute" db:"-"`
	RichText    *model.StringInterface `json:"rich_text"`
}

func (a *AttributeValue) String() string {
	return a.Name
}

func (a *AttributeValue) InputType() string {
	return a.Attribute.InputType
}

func (a *AttributeValue) IsValid() *model.AppError {
	outer := model.CreateAppErrorForModel(
		"model.attribute_value.is_valid.%s.app_error",
		"attribute_value_id=",
		"AttributeValue.IsValid",
	)
	if !model.IsValidId(a.Id) {
		return outer("id", nil)
	}
	if !model.IsValidId(a.AttributeID) {
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

	return nil
}

func (a *AttributeValue) PreSave() {
	if a.Id == "" {
		a.Id = model.NewId()
	}
	a.Name = model.SanitizeUnicode(a.Name)
	a.Slug = slug.Make(a.Name)
}

func (a *AttributeValue) PreUpdate() {
	a.Name = model.SanitizeUnicode(a.Name)
}

func (a *AttributeValue) ToJson() string {
	return model.ModelToJson(a)
}

func AttributeValueFromJson(data io.Reader) *AttributeValue {
	var a AttributeValue
	model.ModelFromJson(&a, data)
	return &a
}

// ---------------------------

const (
	ATTRIBUTE_VALUE_TRANSLATION_NAME_MAX_LENGTH = 100
)

type AttributeValueTranslation struct {
	Id               string                 `json:"id"`
	LanguageCode     string                 `json:"language_code"`
	AttributeValueID string                 `json:"attribute_value"`
	Name             string                 `json:"name"`
	RichText         *model.StringInterface `json:"rich_text"`
}

func (a *AttributeValueTranslation) String() string {
	return a.Name
}

func (a *AttributeValueTranslation) IsValid() *model.AppError {
	outer := model.CreateAppErrorForModel(
		"model.attribute_value_translation.is_valid.%s.app_error",
		"attribute_value_translation_id=",
		"AttributeValueTranslation.IsValid",
	)
	if !model.IsValidId(a.Id) {
		return outer("id", nil)
	}
	if !model.IsValidId(a.AttributeValueID) {
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
		a.Id = model.NewId()
	}
	a.Name = model.SanitizeUnicode(a.Name)
}

func (a *AttributeValueTranslation) PreUpdate() {
	a.Name = model.SanitizeUnicode(a.Name)
}

func (a *AttributeValueTranslation) ToJson() string {
	return model.ModelToJson(a)
}

func AttributeValueTranslationFromJson(data io.Reader) *AttributeValueTranslation {
	var a *AttributeValueTranslation
	model.ModelFromJson(&a, data)
	return a
}
