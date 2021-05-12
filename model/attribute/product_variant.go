package attribute

import "github.com/sitename/sitename/model"

type AssignedVariantAttributeValue struct {
	Id           string `json:"id"`
	ValueID      string `json:"value_id"`
	AssignmentID string `json:"assignment_id"`
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

type AssignedVariantAttribute struct {
	Id                    string            `json:"id"`
	VariantID             string            `json:"variant_id"`
	AssignmentID          string            `json:"assignment_id"`
	Values                []*AttributeValue `json:"values"`
	BaseAssignedAttribute `db:"-"`
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

func (a *AssignedVariantAttribute) ToJson() string {
	return model.ModelToJson(a)
}
