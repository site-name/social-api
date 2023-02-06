package model

import (
	"github.com/Masterminds/squirrel"
)

// AttributeID unique together with ProductTypeID
type AttributeProduct struct {
	Id            string `json:"id"`
	AttributeID   string `json:"attribute_id"`    // to attribute.Attribute
	ProductTypeID string `json:"product_type_id"` // to product.ProductType
	Sortable

	attribute *Attribute `db:"-"`
}

func (a *AttributeProduct) GetAttribute() *Attribute {
	return a.attribute
}

func (a *AttributeProduct) SetAttribute(at *Attribute) {
	a.attribute = at
}

// AttributeProductFilterOption is used when finding an attributeProduct.
type AttributeProductFilterOption struct {
	AttributeID                  squirrel.Sqlizer
	ProductTypeID                squirrel.Sqlizer
	AttributeVisibleInStoreFront *bool
}

func (a *AttributeProduct) IsValid() *AppError {
	outer := CreateAppErrorForModel(
		"attribute_product.is_valid.%s.app_error",
		"attribute_product_id=",
		"AttributeProduct.IsValid",
	)
	if !IsValidId(a.Id) {
		return outer("id", nil)
	}
	if !IsValidId(a.AttributeID) {
		return outer("attribute_id", &a.Id)
	}
	if !IsValidId(a.ProductTypeID) {
		return outer("product_type_id", &a.Id)
	}

	return nil
}

func (a *AttributeProduct) ToJSON() string {
	return ModelToJson(a)
}

func (a *AttributeProduct) PreSave() {
	if a.Id == "" {
		a.Id = NewId()
	}
}

func (a *AttributeProduct) DeepCopy() *AttributeProduct {
	if a == nil {
		return nil
	}

	res := *a
	if a.SortOrder != nil {
		res.SortOrder = NewPrimitive(*a.SortOrder)
	}
	if a.attribute != nil {
		res.attribute = a.attribute.DeepCopy()
	}
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

	AttributeProduct_Attribute_VisibleInStoreFront *bool // INNER JOIN AttributeProduct ON ... INNER JOIN Attribute ON ... WHERE Attribute.VisibleInStoreFront ...
}

func (a *AssignedProductAttribute) IsValid() *AppError {
	outer := CreateAppErrorForModel(
		"assigned_product_attribute.is_valid.%s.app_error",
		"assigned_product_attribute_id=",
		"AssignedProductAttribute.IsValid",
	)
	if !IsValidId(a.Id) {
		return outer("id", nil)
	}
	if !IsValidId(a.ProductID) {
		return outer("product_id", &a.Id)
	}
	if !IsValidId(a.AssignmentID) {
		return outer("assignment_id", &a.Id)
	}

	return nil
}

func (a *AssignedProductAttribute) ToJSON() string {
	return ModelToJson(a)
}

func (a *AssignedProductAttribute) PreSave() {
	if a.Id == "" {
		a.Id = NewId()
	}
}

func (a *AssignedProductAttribute) DeepCopy() *AssignedProductAttribute {
	if a == nil {
		return nil
	}

	res := *a
	res.AttributeValues = a.AttributeValues.DeepCopy()
	if a.AttributeProduct != nil {
		res.AttributeProduct = a.AttributeProduct.DeepCopy()
	}
	return &res
}

func (a AssignedProductAttributes) DeepCopy() AssignedProductAttributes {
	res := make(AssignedProductAttributes, len(a))
	for idx, item := range a {
		res[idx] = item.DeepCopy()
	}

	return res
}

// ValueID unique together AssignmentID
type AssignedProductAttributeValue struct {
	Id           string `json:"id"`
	ValueID      string `json:"value_id"`      // to AttributeValue
	AssignmentID string `json:"assignment_id"` // to AssignedProductAttribute
	Sortable
}

type AssignedProductAttributeValueFilterOptions struct {
	ValueID      squirrel.Sqlizer
	AssignmentID squirrel.Sqlizer
}

func (a *AssignedProductAttributeValue) IsValid() *AppError {
	outer := CreateAppErrorForModel(
		"assigned_product_attribute.is_valid.%s.app_error",
		"assigned_product_attribute_id=",
		"AssignedProductAttributeValue.IsValid",
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

func (a *AssignedProductAttributeValue) ToJSON() string {
	return ModelToJson(a)
}

func (a *AssignedProductAttributeValue) PreSave() {
	if a.Id == "" {
		a.Id = NewId()
	}
}

func (a *AssignedProductAttributeValue) DeepCopy() *AssignedProductAttributeValue {
	res := *a

	if a.SortOrder != nil {
		res.SortOrder = NewPrimitive(*a.SortOrder)
	}
	return &res
}
