package attribute

import (
	"io"

	"github.com/sitename/sitename/model"
)

// AttributeID unique together with ProductTypeID
type AttributeProduct struct {
	Id            string `json:"id"`
	AttributeID   string `json:"attribute_id"`    // to attribute.Attribute
	ProductTypeID string `json:"product_type_id"` // to product.ProductType
	model.Sortable
}

// AttributeProductFilterOption is used when finding an attributeProduct.
type AttributeProductFilterOption struct {
	AttributeID   *model.StringFilter
	ProductTypeID *model.StringFilter
}

func (a *AttributeProduct) IsValid() *model.AppError {
	outer := model.CreateAppErrorForModel(
		"model.attribute_product.is_valid.%s.app_error",
		"attribute_product_id=",
		"AttributeProduct.IsValid",
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

func (a *AttributeProduct) ToJson() string {
	return model.ModelToJson(a)
}

func AttributeProductFromJson(data io.Reader) *AttributeProduct {
	var a AttributeProduct
	model.ModelFromJson(&a, data)
	return &a
}

func (a *AttributeProduct) PreSave() {
	if a.Id == "" {
		a.Id = model.NewId()
	}
}

// Associate a product type attribute and selected values to a given product
// ProductID unique with AssignmentID
type AssignedProductAttribute struct {
	Id           string `json:"id"`
	ProductID    string `json:"product_id"`    // to product.Product
	AssignmentID string `json:"assignment_id"` // to attribute.AttributeProduct
	BaseAssignedAttribute
}

// AssignedProductAttributeFilterOption is used to filter or creat new AssignedProductAttribute
type AssignedProductAttributeFilterOption struct {
	ProductID    *model.StringFilter
	AssignmentID *model.StringFilter
}

func (a *AssignedProductAttribute) IsValid() *model.AppError {
	outer := model.CreateAppErrorForModel(
		"model.assigned_product_attribute.is_valid.%s.app_error",
		"assigned_product_attribute_id=",
		"AssignedProductAttribute.IsValid",
	)
	if !model.IsValidId(a.Id) {
		return outer("id", nil)
	}
	if !model.IsValidId(a.ProductID) {
		return outer("product_id", &a.Id)
	}
	if !model.IsValidId(a.AssignmentID) {
		return outer("assignment_id", &a.Id)
	}

	return nil
}

func (a *AssignedProductAttribute) ToJson() string {
	return model.ModelToJson(a)
}

func (a *AssignedProductAttribute) PreSave() {
	if a.Id == "" {
		a.Id = model.NewId()
	}
}

func AssignedProductAttributeFromJson(data io.Reader) *AssignedProductAttribute {
	var a AssignedProductAttribute
	model.ModelFromJson(&a, data)
	return &a
}

// ValueID unique together AssignmentID
type AssignedProductAttributeValue struct {
	Id           string `json:"id"`
	ValueID      string `json:"value_id"`      // to attribute.AttributeValue
	AssignmentID string `json:"assignment_id"` // to attribute.AssignedProductAttribute
	model.Sortable
}

func (a *AssignedProductAttributeValue) IsValid() *model.AppError {
	outer := model.CreateAppErrorForModel(
		"model.assigned_product_attribute.is_valid.%s.app_error",
		"assigned_product_attribute_id=",
		"AssignedProductAttributeValue.IsValid",
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

func (a *AssignedProductAttributeValue) ToJson() string {
	return model.ModelToJson(a)
}

func (a *AssignedProductAttributeValue) PreSave() {
	if a.Id == "" {
		a.Id = model.NewId()
	}
}

func AssignedProductAttributeValueFromJson(data io.Reader) *AssignedProductAttributeValue {
	var a AssignedProductAttributeValue
	model.ModelFromJson(&a, data)
	return &a
}
