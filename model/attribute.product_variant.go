package model

import (
	"github.com/Masterminds/squirrel"
)

// ValueID unique together with AssignmentID
type AssignedVariantAttributeValue struct {
	Id           string `json:"id"`
	ValueID      string `json:"value_id"`      // unique together
	AssignmentID string `json:"assignment_id"` // unique together
	Sortable
}

type AssignedVariantAttributeValueFilterOptions struct {
	ValueID      squirrel.Sqlizer
	AssignmentID squirrel.Sqlizer
}

func (a *AssignedVariantAttributeValue) IsValid() *AppError {
	outer := CreateAppErrorForModel(
		"assigned_variant_attribute_value.is_valid.%s.app_error",
		"assigned_variant_attribute_value_id=",
		"AssignedVariantAttributeValue.IsValid",
	)

	if !IsValidId(a.Id) {
		return outer("id", nil)
	}
	if !IsValidId(a.ValueID) {
		return outer("value_id", &a.Id)
	}
	if !IsValidId(a.AssignmentID) {
		return outer("assignment_id", &a.Id)
	}

	return nil
}

func (a *AssignedVariantAttributeValue) PreSave() {
	if a.Id == "" {
		a.Id = NewId()
	}
}

func (a *AssignedVariantAttributeValue) ToJSON() string {
	return ModelToJson(a)
}

func (a *AssignedVariantAttributeValue) DeepCopy() *AssignedVariantAttributeValue {
	res := *a
	if a.SortOrder != nil {
		res.SortOrder = NewPrimitive(*a.SortOrder)
	}
	return &res
}

// Associate a product type attribute and selected values to a given variant.
type AssignedVariantAttribute struct {
	Id           string `json:"id"`
	VariantID    string `json:"variant_id"`    // to product.ProductVariant
	AssignmentID string `json:"assignment_id"` // to attribute.AttributeVariant
}

// AssignedVariantAttributeFilterOption is used for lookup, if cannot found, creating new instance
type AssignedVariantAttributeFilterOption struct {
	VariantID    squirrel.Sqlizer
	AssignmentID squirrel.Sqlizer

	Assignment_Attribute_VisibleInStoreFront *bool // INNER JOIN AttributeVariant ON ... INNER JOIN Attributes ON ... WHERE Attributes.VisibleInStoreFront ...

	AssignmentAttributeInputType squirrel.Sqlizer // INNER JOIN AttributeVariants ON () INNER JOIN Attributes ON () WHERE Attributes.InputType
	AssignmentAttributeType      squirrel.Sqlizer // INNER JOIN AttributeVariants ON () INNER JOIN Attributes ON () WHERE Attributes.Type
}

func (a *AssignedVariantAttribute) IsValid() *AppError {
	outer := CreateAppErrorForModel(
		"assigned_variant_attribute.is_valid.%s.app_error",
		"assigned_variant_attribute_id=",
		"AssignedVariantAttribute.IsValid",
	)
	if !IsValidId(a.Id) {
		return outer("id", nil)
	}
	if !IsValidId(a.VariantID) {
		return outer("variant_id", &a.Id)
	}
	if !IsValidId(a.AssignmentID) {
		return outer("assignment_id", &a.Id)
	}

	return nil
}

func (a *AssignedVariantAttribute) ToJSON() string {
	return ModelToJson(a)
}

func (a *AssignedVariantAttribute) PreSave() {
	if a.Id == "" {
		a.Id = NewId()
	}
}
