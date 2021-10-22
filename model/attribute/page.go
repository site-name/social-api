package attribute

import (
	"github.com/sitename/sitename/model"
)

// ValueID unique together with AssignmentID
type AssignedPageAttributeValue struct {
	Id           string `json:"id"`
	ValueID      string `json:"value_id"`      // AttributeValue
	AssignmentID string `json:"assignment_id"` // AssignedPageAttribute
	model.Sortable
}

func (a *AssignedPageAttributeValue) IsValid() *model.AppError {
	outer := model.CreateAppErrorForModel(
		"model.assigned_page_attribute_value.is_valid.%s.app_error",
		"assigned_page_sttribute_value_id=",
		"AssignedPageAttributeValue.IsValid",
	)
	if !model.IsValidId(a.Id) {
		return outer("id", nil)
	}
	if !model.IsValidId(a.ValueID) {
		return outer("value_id", &a.ValueID)
	}
	if !model.IsValidId(a.AssignmentID) {
		return outer("assignment_id", &a.Id)
	}

	return nil
}

func (a *AssignedPageAttributeValue) ToJSON() string {
	return model.ModelToJson(a)
}

func (a *AssignedPageAttributeValue) PreSave() {
	if a.Id == "" {
		a.Id = model.NewId()
	}
}

func (a *AssignedPageAttributeValue) DeepCopy() *AssignedPageAttributeValue {
	res := *a

	return &res
}

// Associate a page type attribute and selected values to a given page.
// PageID unique together with AssignmentID
type AssignedPageAttribute struct {
	Id                    string `json:"id"`
	PageID                string `json:"page_id"`
	AssignmentID          string `json:"assignment_id"` // AttributePage
	BaseAssignedAttribute `db:"-"`
}

// AssignedPageAttributeFilterOption is used to find or creat new AssignedPageAttribute
type AssignedPageAttributeFilterOption struct {
	PageID       *model.StringFilter
	AssignmentID *model.StringFilter
}

func (a *AssignedPageAttribute) IsValid() *model.AppError {
	outer := model.CreateAppErrorForModel(
		"model.assigned_page_attribute.is_valid.%s.app_error",
		"assigned_page_attribute_id=",
		"AssignedPageAttribute.IsValid",
	)

	if !model.IsValidId(a.Id) {
		return outer("id", nil)
	}
	if !model.IsValidId(a.PageID) {
		return outer("page_id", &a.Id)
	}
	if !model.IsValidId(a.AssignmentID) {
		return outer("assignment_id", &a.Id)
	}

	return nil
}

func (a *AssignedPageAttribute) ToJSON() string {
	return model.ModelToJson(a)
}

func (a *AssignedPageAttribute) PreSave() {
	if a.Id == "" {
		a.Id = model.NewId()
	}
}

func (a *AssignedPageAttribute) DeepCopy() *AssignedPageAttribute {
	res := *a

	return &res
}
