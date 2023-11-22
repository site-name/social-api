package model_helper

import (
	"net/http"

	"github.com/sitename/sitename/model"
)

func AttributePageIsValid(a *model.AttributePage) *AppError {
	if !IsValidId(a.AttributeID) {
		return NewAppError("AttributePage.IsValid", "model.attribute_page.is_valid.attribute_id.app_error", nil, "please provide valid attribute id", http.StatusBadRequest)
	}
	if !IsValidId(a.PageTypeID) {
		return NewAppError("AttributePage.IsValid", "model.attribute_page.is_valid.page_type_id.app_error", nil, "please provide valid page type id", http.StatusBadRequest)
	}
	return nil
}

func AssignedPageAttributeValueIsValid(a *model.AssignedPageAttributeValue) *AppError {
	if !IsValidId(a.ValueID) {
		return NewAppError("AssignedPageAttributeValue.IsValid", "model.assigned_page_attribute_value.is_valid.value_id.app_error", nil, "please provide valid value id", http.StatusBadRequest)
	}
	if !IsValidId(a.AssignmentID) {
		return NewAppError("AssignedPageAttributeValue.IsValid", "model.assigned_page_attribute_value.is_valid.assignment_id.app_error", nil, "please provide valid assignment id", http.StatusBadRequest)
	}
	return nil
}

func AssignedPageAttributeIsValid(a *model.AssignedPageAttribute) *AppError {
	if !IsValidId(a.PageID) {
		return NewAppError("AssignedPageAttribute.IsValid", "model.assigned_page_attribute.is_valid.page_id.app_error", nil, "", http.StatusBadRequest)
	}
	if !IsValidId(a.AssignmentID) {
		return NewAppError("AssignedPageAttribute.IsValid", "model.assigned_page_attribute.is_valid.assignment_id.app_error", nil, "", http.StatusBadRequest)
	}

	return nil
}
