package attribute

import (
	"github.com/Masterminds/squirrel"
	"github.com/sitename/sitename/model"
)

// AttributeID unique together with ProductTypeID
type AttributeProduct struct {
	Id            string `json:"id"`
	AttributeID   string `json:"attribute_id"`    // to attribute.Attribute
	ProductTypeID string `json:"product_type_id"` // to product.ProductType
	model.Sortable

	Attribute *Attribute `json:"-" db:"-"`
}

// AttributeProductFilterOption is used when finding an attributeProduct.
type AttributeProductFilterOption struct {
	AttributeID   squirrel.Sqlizer
	ProductTypeID squirrel.Sqlizer
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

func (a *AttributeProduct) ToJSON() string {
	return model.ModelToJson(a)
}

func (a *AttributeProduct) PreSave() {
	if a.Id == "" {
		a.Id = model.NewId()
	}
}

func (a *AttributeProduct) DeepCopy() *AttributeProduct {
	if a == nil {
		return nil
	}

	res := *a
	res.Attribute = a.Attribute.DeepCopy()
	return &res
}

// Associate a product type attribute and selected values to a given product
// ProductID unique with AssignmentID
type AssignedProductAttribute struct {
	Id           string `json:"id"`
	ProductID    string `json:"product_id"`    // to product.Product
	AssignmentID string `json:"assignment_id"` // to attribute.AttributeProduct

	AttributeValues  AttributeValues   `json:"-" db:"-"`
	AttributeProduct *AttributeProduct `json:"-" db:"-"`
}

type AssignedProductAttributes []*AssignedProductAttribute

// AssignedProductAttributeFilterOption is used to filter or creat new AssignedProductAttribute
type AssignedProductAttributeFilterOption struct {
	ProductID    squirrel.Sqlizer
	AssignmentID squirrel.Sqlizer
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

func (a *AssignedProductAttribute) ToJSON() string {
	return model.ModelToJson(a)
}

func (a *AssignedProductAttribute) PreSave() {
	if a.Id == "" {
		a.Id = model.NewId()
	}
}

func (a *AssignedProductAttribute) DeepCopy() *AssignedProductAttribute {
	if a == nil {
		return nil
	}

	res := *a
	res.AttributeValues = a.AttributeValues.DeepCopy()
	res.AttributeProduct = a.AttributeProduct.DeepCopy()
	return &res
}

func (a AssignedProductAttributes) DeepCopy() AssignedProductAttributes {
	res := AssignedProductAttributes{}
	for _, item := range a {
		res = append(res, item.DeepCopy())
	}

	return res
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

func (a *AssignedProductAttributeValue) ToJSON() string {
	return model.ModelToJson(a)
}

func (a *AssignedProductAttributeValue) PreSave() {
	if a.Id == "" {
		a.Id = model.NewId()
	}
}

func (a *AssignedProductAttributeValue) DeepCopy() *AssignedProductAttributeValue {
	res := *a
	return &res
}
