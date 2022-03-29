package attribute

import (
	"github.com/Masterminds/squirrel"
	"github.com/sitename/sitename/model"
)

// ValueID unique together with AssignmentID
type AssignedVariantAttributeValue struct {
	Id           string `json:"id"`
	ValueID      string `json:"value_id"`      // unique together
	AssignmentID string `json:"assignment_id"` // unique together
	model.Sortable
}

func (a *AssignedVariantAttributeValue) IsValid() *model.AppError {
	outer := model.CreateAppErrorForModel(
		"model.assigned_variant_attribute_value.is_valid.%s.app_error",
		"assigned_variant_attribute_value_id=",
		"AssignedVariantAttributeValue.IsValid",
	)

	if !model.IsValidId(a.Id) {
		return outer("id", nil)
	}
	if !model.IsValidId(a.ValueID) {
		return outer("value_id", &a.Id)
	}
	if !model.IsValidId(a.AssignmentID) {
		return outer("assignment_id", &a.Id)
	}

	return nil
}

func (a *AssignedVariantAttributeValue) PreSave() {
	if a.Id == "" {
		a.Id = model.NewId()
	}
}

func (a *AssignedVariantAttributeValue) ToJSON() string {
	return model.ModelToJson(a)
}

func (a *AssignedVariantAttributeValue) DeepCopy() *AssignedVariantAttributeValue {
	res := *a
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

	AssignmentAttributeInputType squirrel.Sqlizer // INNER JOIN AttributeVariants ON () INNER JOIN Attributes ON () WHERE Attributes.InputType
	AssignmentAttributeType      squirrel.Sqlizer // INNER JOIN AttributeVariants ON () INNER JOIN Attributes ON () WHERE Attributes.Type
}

func (a *AssignedVariantAttribute) IsValid() *model.AppError {
	outer := model.CreateAppErrorForModel(
		"model.assigned_variant_attribute.is_valid.%s.app_error",
		"assigned_variant_attribute_id=",
		"AssignedVariantAttribute.IsValid",
	)
	if !model.IsValidId(a.Id) {
		return outer("id", nil)
	}
	if !model.IsValidId(a.VariantID) {
		return outer("variant_id", &a.Id)
	}
	if !model.IsValidId(a.AssignmentID) {
		return outer("assignment_id", &a.Id)
	}

	return nil
}

func (a *AssignedVariantAttribute) ToJSON() string {
	return model.ModelToJson(a)
}

func (a *AssignedVariantAttribute) PreSave() {
	if a.Id == "" {
		a.Id = model.NewId()
	}
}
