package attribute

import (
	"io"

	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model/product_and_discount"
)

// AttributeID unique together with ProductTypeID
type AttributeVariant struct {
	Id               string                                 `json:"id"`
	AttributeID      string                                 `json:"attribute_id"`
	ProductTypeID    string                                 `json:"product_type_id"`
	AssignedVariants []*product_and_discount.ProductVariant `json:"assigned_variants" db:"-"` // through attribute.AssignedVariantAttribute
	model.Sortable
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
