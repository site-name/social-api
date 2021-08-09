package attribute

import (
	"io"

	"github.com/sitename/sitename/model"
)

// AttributeID unique together with ProductTypeID
type AttributeVariant struct {
	Id            string `json:"id"`
	AttributeID   string `json:"attribute_id"`
	ProductTypeID string `json:"product_type_id"`
	model.Sortable
}

// AttributeVariantFilterOption is used to find `AttributeVariant`.
//
// properties can be provided partially or fully
type AttributeVariantFilterOption struct {
	AttributeID string `json:"attribute_id"`
	ProductID   string `json:"product_id"` // required
}

func (a *AttributeVariant) IsValid() *model.AppError {
	outer := model.CreateAppErrorForModel(
		"model.attribute_variant.is_valid.%s.app_error",
		"attribute_variant_id=",
		"AttributeVariant.IsValid",
	)
	if !model.IsValidId(a.Id) {
		return outer("id", nil)
	}
	if !model.IsValidId(a.AttributeID) {
		return outer("attribute_id", &a.Id)
	}
	if !model.IsValidId(a.ProductTypeID) {
		return outer("product_type_id", &a.Id)
	}

	return nil
}

func (a *AttributeVariant) ToJson() string {
	return model.ModelToJson(a)
}

func AttributeVariantFromJson(data io.Reader) *AttributeVariant {
	var a AttributeVariant
	model.ModelFromJson(&a, data)
	return &a
}

func (a *AttributeVariant) PreSave() {
	if a.Id == "" {
		a.Id = model.NewId()
	}
}
