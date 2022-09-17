package model

import (
	"github.com/Masterminds/squirrel"
)

// ValueID unique together with AssignmentID
type AssignedPageAttributeValue struct {
	Id           string `json:"id"`
	ValueID      string `json:"value_id"`      // AttributeValue
	AssignmentID string `json:"assignment_id"` // AssignedPageAttribute
	Sortable
}

func (a *AssignedPageAttributeValue) IsValid() *AppError {
	outer := CreateAppErrorForModel(
		"assigned_page_attribute_value.is_valid.%s.app_error",
		"assigned_page_sttribute_value_id=",
		"AssignedPageAttributeValue.IsValid",
	)
	if !IsValidId(a.Id) {
		return outer("id", nil)
	}
	if !IsValidId(a.ValueID) {
		return outer("value_id", &a.ValueID)
	}
	if !IsValidId(a.AssignmentID) {
		return outer("assignment_id", &a.Id)
	}

	return nil
}

func (a *AssignedPageAttributeValue) ToJSON() string {
	return ModelToJson(a)
}

func (a *AssignedPageAttributeValue) PreSave() {
	if a.Id == "" {
		a.Id = NewId()
	}
}

func (a *AssignedPageAttributeValue) DeepCopy() *AssignedPageAttributeValue {
	res := *a

	return &res
}

// Associate a page type attribute and selected values to a given page.
// PageID unique together with AssignmentID
type AssignedPageAttribute struct {
	Id           string `json:"id"`
	PageID       string `json:"page_id"`
	AssignmentID string `json:"assignment_id"` // AttributePage
}

// AssignedPageAttributeFilterOption is used to find or creat new AssignedPageAttribute
type AssignedPageAttributeFilterOption struct {
	PageID       squirrel.Sqlizer
	AssignmentID squirrel.Sqlizer
}

func (a *AssignedPageAttribute) IsValid() *AppError {
	outer := CreateAppErrorForModel(
		"assigned_page_attribute.is_valid.%s.app_error",
		"assigned_page_attribute_id=",
		"AssignedPageAttribute.IsValid",
	)

	if !IsValidId(a.Id) {
		return outer("id", nil)
	}
	if !IsValidId(a.PageID) {
		return outer("page_id", &a.Id)
	}
	if !IsValidId(a.AssignmentID) {
		return outer("assignment_id", &a.Id)
	}

	return nil
}

func (a *AssignedPageAttribute) ToJSON() string {
	return ModelToJson(a)
}

func (a *AssignedPageAttribute) PreSave() {
	if a.Id == "" {
		a.Id = NewId()
	}
}

func (a *AssignedPageAttribute) DeepCopy() *AssignedPageAttribute {
	res := *a

	return &res
}
