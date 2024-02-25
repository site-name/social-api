package model_helper

import (
	"net/http"

	"github.com/gosimple/slug"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/modules/measurement"
	"github.com/sitename/sitename/modules/model_types"
)

func AttributePageIsValid(a model.AttributePage) *AppError {
	if !IsValidId(a.ID) {
		return NewAppError("AttributePageIsValid", "model.attribute_page.is_valid.id.app_error", nil, "invalid id", http.StatusBadRequest)
	}
	if !IsValidId(a.AttributeID) {
		return NewAppError("AttributePageIsValid", "model.attribute_page.is_valid.attribute_id.app_error", nil, "please provide valid attribute id", http.StatusBadRequest)
	}
	if !IsValidId(a.PageTypeID) {
		return NewAppError("AttributePageIsValid", "model.attribute_page.is_valid.page_type_id.app_error", nil, "please provide valid page type id", http.StatusBadRequest)
	}
	return nil
}

func AttributePagePreSave(a *model.AttributePage) {
	if a.ID == "" {
		a.ID = NewId()
	}
}

type AttributePageFilterOption struct {
	CommonQueryOptions
}

func AssignedPageAttributeValueIsValid(a model.AssignedPageAttributeValue) *AppError {
	if !IsValidId(a.ID) {
		return NewAppError("AssignedPageAttributeValueIsValid", "model.assigned_page_attribute_value.is_valid.id.app_error", nil, "invalid id", http.StatusBadRequest)
	}
	if !IsValidId(a.ValueID) {
		return NewAppError("AssignedPageAttributeValueIsValid", "model.assigned_page_attribute_value.is_valid.value_id.app_error", nil, "please provide valid value id", http.StatusBadRequest)
	}
	if !IsValidId(a.AssignmentID) {
		return NewAppError("AssignedPageAttributeValueIsValid", "model.assigned_page_attribute_value.is_valid.assignment_id.app_error", nil, "please provide valid assignment id", http.StatusBadRequest)
	}
	return nil
}

func AssignedPageAttributeValuePreSave(a *model.AssignedPageAttributeValue) {
	if a.ID == "" {
		a.ID = NewId()
	}
}

func AssignedPageAttributeIsValid(a model.AssignedPageAttribute) *AppError {
	if !IsValidId(a.ID) {
		return NewAppError("AssignedPageAttributeIsValid", "model.assigned_page_attribute.is_valid.id.app_error", nil, "", http.StatusBadRequest)
	}
	if !IsValidId(a.PageID) {
		return NewAppError("AssignedPageAttributeIsValid", "model.assigned_page_attribute.is_valid.page_id.app_error", nil, "", http.StatusBadRequest)
	}
	if !IsValidId(a.AssignmentID) {
		return NewAppError("AssignedPageAttributeIsValid", "model.assigned_page_attribute.is_valid.assignment_id.app_error", nil, "", http.StatusBadRequest)
	}

	return nil
}

func AttributeProductIsValid(a model.CategoryAttribute) *AppError {
	if !IsValidId(a.ID) {
		return NewAppError("AttributeProductIsValid", "model.attribute_product.is_valid.id.app_error", nil, "", http.StatusBadRequest)
	}
	if !IsValidId(a.AttributeID) {
		return NewAppError("AttributeProductIsValid", "model.attribute_product.is_valid.attribute_id.app_error", nil, "", http.StatusBadRequest)
	}
	if !IsValidId(a.CategoryID) {
		return NewAppError("AttributeProductIsValid", "model.attribute_product.is_valid.product_type_id.app_error", nil, "", http.StatusBadRequest)
	}
	return nil
}

func AssignedProductAttributeIsValid(a model.AssignedProductAttribute) *AppError {
	if !IsValidId(a.ID) {
		return NewAppError("AssignedProductAttributeIsValid", "model.assigned_product_attribute.is_valid.id.app_error", nil, "", http.StatusBadRequest)
	}
	if !IsValidId(a.ProductID) {
		return NewAppError("AssignedProductAttributeIsValid", "model.assigned_product_attribute.is_valid.product_id.app_error", nil, "", http.StatusBadRequest)
	}
	if !IsValidId(a.AssignmentID) {
		return NewAppError("AssignedProductAttributeIsValid", "model.assigned_product_attribute.is_valid.assignment_id.app_error", nil, "", http.StatusBadRequest)
	}

	return nil
}

func AssignedProductAttributeValueIsValid(a model.AssignedProductAttributeValue) *AppError {
	if !IsValidId(a.ID) {
		return NewAppError("AssignedProductAttributeValueIsValid", "model.assigned_product_attribute_value.is_valid.id.app_error", nil, "", http.StatusBadRequest)
	}
	if !IsValidId(a.AssignmentID) {
		return NewAppError("AssignedProductAttributeValueIsValid", "model.assigned_product_attribute_value.is_valid.assignment_id.app_error", nil, "", http.StatusBadRequest)
	}
	if !IsValidId(a.ValueID) {
		return NewAppError("AssignedProductAttributeValueIsValid", "model.assigned_product_attribute_value.is_valid.value_id.app_error", nil, "", http.StatusBadRequest)
	}

	return nil
}

func AssignedProductAttributeValuePreSave(a *model.AssignedProductAttributeValue) {
	if a.ID == "" {
		a.ID = NewId()
	}
}

func AttributeValueCommonPre(a *model.AttributeValue) {
	a.Name = SanitizeUnicode(a.Name)
}

func AttributeValuePreSave(a *model.AttributeValue) {
	if a.ID == "" {
		a.ID = NewId()
	}
	AttributeValueCommonPre(a)
	a.Slug = slug.Make(a.Name)
}

func AttributeValueIsValid(a model.AttributeValue) *AppError {
	if !IsValidId(a.ID) {
		return NewAppError("AttributeValueIsValid", "model.attribute_value.is_valid.id.app_error", nil, "", http.StatusBadRequest)
	}
	if !IsValidId(a.AttributeID) {
		return NewAppError("AttributeValue.IsValid", "model.attribute_value.is_valid.attribute_id.app_error", nil, "", http.StatusBadRequest)
	}
	if a.Datetime.IsNil() || a.Datetime.Time.IsZero() {
		return NewAppError("AttributeValue.IsValid", "model.attribute_value.is_valid.date_time.app_error", nil, "", http.StatusBadRequest)
	}
	if !slug.IsSlug(a.Slug) {
		return NewAppError("AttributeValue.IsValid", "model.attribute_value.is_valid.slug.app_error", nil, "", http.StatusBadRequest)
	}

	return nil
}

func AttributePreSave(a *model.Attribute) {
	if a.ID == "" {
		a.ID = NewId()
	}
	attributeCommonPre(a)
	if a.Slug == "" {
		a.Slug = slug.Make(a.Name)
	}
}

func attributeCommonPre(a *model.Attribute) {
	a.Name = SanitizeUnicode(a.Name)
	if a.InputType.IsValid() != nil {
		a.InputType = model.AttributeInputTypeDropdown
	}
}

func AttributePreUpdate(a *model.Attribute) {
	attributeCommonPre(a)
}

type AttributeFilterOption struct {
	CommonQueryOptions
	Search   string // WHERE Attributes.Name ILIKE ... OR Attributes.Slug ILIKE ...
	Metadata model_types.JsonMap
	Preload  []string
}

func AttributeIsValid(a model.Attribute) *AppError {
	if !IsValidId(a.ID) {
		return NewAppError("Attribute.IsValid", "model.attribute.is_valid.id.app_error", nil, "please provide valid attribute id", http.StatusBadRequest)
	}
	if a.Type.IsValid() != nil {
		return NewAppError("Attribute.IsValid", "model.attribute.is_valid.type.app_error", nil, "please provide valid attribute type", http.StatusBadRequest)
	}
	if a.InputType.IsValid() != nil {
		return NewAppError("Attribute.IsValid", "model.attribute.is_valid.input_type.app_error", nil, "please provide valid attribute input type", http.StatusBadRequest)
	}
	if a.EntityType.Valid && a.EntityType.Val.IsValid() != nil {
		return NewAppError("Attribute.IsValid", "model.attribute.is_valid.entity_type.app_error", nil, "please provide valid attribute entity type", http.StatusBadRequest)
	}
	if !a.Unit.IsNil() && measurement.MeasurementUnitMap[*a.Unit.String] == "" {
		return NewAppError("Attribute.IsValid", "model.attribute.is_valid.unit.app_error", nil, "please provide valid attribute unit", http.StatusBadRequest)
	}
	if !slug.IsSlug(a.Slug) {
		return NewAppError("Attribute.IsValid", "model.attribute.is_valid.slug.app_error", nil, "please provide valid attribute slug", http.StatusBadRequest)
	}

	return nil
}

type AttributeValueFilterOptions struct {
	CommonQueryOptions
}

type AttributeCategoryFilterOptions struct {
	CommonQueryOptions
}

func CategoryAttributePreSave(a *model.CategoryAttribute) {
	if a.ID == "" {
		a.ID = NewId()
	}
}

func CategoryAttributeIsValid(a model.CategoryAttribute) *AppError {
	if !IsValidId(a.ID) {
		return NewAppError("CategoryAttributeIsValid", "model.category_attribute.is_valid.id.app_error", nil, "", http.StatusBadRequest)
	}
	if !IsValidId(a.CategoryID) {
		return NewAppError("CategoryAttributeIsValid", "model.category_attribute.is_valid.category_id.app_error", nil, "", http.StatusBadRequest)
	}
	if !IsValidId(a.AttributeID) {
		return NewAppError("CategoryAttributeIsValid", "model.category_attribute.is_valid.attribute_id.app_error", nil, "", http.StatusBadRequest)
	}
	return nil
}

type CustomProductAttributeFilterOptions struct {
	CommonQueryOptions
}

func CustomProductAttributePreSave(a *model.CustomProductAttribute) {
	if a.ID == "" {
		a.ID = NewId()
	}
	CustomProductAttributeCommonPre(a)
}

func CustomProductAttributeCommonPre(a *model.CustomProductAttribute) {
	a.Name = SanitizeUnicode(a.Name)
	a.Slug = slug.Make(a.Name)
}

func CustomProductAttributeIsValid(a model.CustomProductAttribute) *AppError {
	if !IsValidId(a.ID) {
		return NewAppError("CustomProductAttributeIsValid", "model.custom_product_attribute.is_valid.id.app_error", nil, "", http.StatusBadRequest)
	}
	if !IsValidId(a.ProductID) {
		return NewAppError("CustomProductAttributeIsValid", "model.custom_product_attribute.is_valid.product_id.app_error", nil, "", http.StatusBadRequest)
	}
	if !slug.IsSlug(a.Slug) {
		return NewAppError("CustomProductAttributeIsValid", "model.custom_product_attribute.is_valid.slug.app_error", nil, "", http.StatusBadRequest)
	}
	return nil
}

type AssignedProductAttributeValueFilterOptions struct {
	CommonQueryOptions
}

type AssignedProductAttributeFilterOption struct {
	CommonQueryOptions
	AttributeProduct_Attribute_VisibleInStoreFront *bool // INNER JOIN AttributeCategory ON ... INNER JOIN Attributes ON ... WHERE Attributes.VisibleInStoreFront ...
}

type AssignedPageAttributeFilterOption struct {
	CommonQueryOptions
}
