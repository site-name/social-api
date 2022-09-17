package model

import (
	"github.com/Masterminds/squirrel"
)

// AttributeID unique with PageTypeID
type AttributePage struct {
	Id          string `json:"id"`
	AttributeID string `json:"attribute_id"` // to attribute.Attribute
	PageTypeID  string `json:"page_type_id"` // to page.PageType
	Sortable
}

// AttributePageFilterOption is used for lookup AttributePage
type AttributePageFilterOption struct {
	PageTypeID  squirrel.Sqlizer
	AttributeID squirrel.Sqlizer
}

func (a *AttributePage) IsValid() *AppError {
	outer := CreateAppErrorForModel(
		"attribute_page.is_valid.%s.app_error",
		"attribute_page_id=",
		"AttributePage.IsValid",
	)
	if !IsValidId(a.Id) {
		return outer("id", nil)
	}
	if !IsValidId(a.AttributeID) {
		return outer("attribute_id", &a.Id)
	}
	if !IsValidId(a.PageTypeID) {
		return outer("page_type_id", &a.Id)
	}

	return nil
}

func (a *AttributePage) PreSave() {
	if a.Id == "" {
		a.Id = NewId()
	}
}

func (a *AttributePage) ToJSON() string {
	return ModelToJson(a)
}
