package model

import (
	"io"
	"unicode/utf8"

	"github.com/gosimple/slug"
)

const (
	PAGE_TYPE_NAME_MAX_LENGTH = 250
	PAGE_TYPE_SLUG_MAX_LENGTH = 255
)

type PageType struct {
	Id   string `json:"id"`
	Name string `json:"name"`
	Slug string `json:"alug"`
	ModelMetadata
}

func (pt *PageType) IsValid() *AppError {
	outer := CreateAppErrorForModel(
		"page_type.is_valid.%s.app_error",
		"page_type_id=",
		"PageType.IsValid",
	)
	if !IsValidId(pt.Id) {
		return outer("id", nil)
	}
	if utf8.RuneCountInString(pt.Name) > PAGE_TYPE_NAME_MAX_LENGTH {
		return outer("name", &pt.Id)
	}
	if len(pt.Slug) > PAGE_TYPE_SLUG_MAX_LENGTH {
		return outer("slug", &pt.Id)
	}

	return nil
}

func (pt *PageType) PreSave() {
	pt.Name = SanitizeUnicode(pt.Name)
	pt.Slug = slug.Make(pt.Name)
}

func (pt *PageType) PreUpdate() {
	pt.Name = SanitizeUnicode(pt.Name)
	// slug should be kept unchanged
}

func (p *PageType) ToJSON() string {
	return ModelToJson(p)
}

func PageTypeFromJson(data io.Reader) *PageType {
	var p PageType
	ModelFromJson(&p, data)
	return &p
}
