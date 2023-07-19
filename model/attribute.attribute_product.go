package model

import (
	"github.com/Masterminds/squirrel"
	"gorm.io/gorm"
)

// AttributeID unique together with ProductTypeID
type AttributeProduct struct {
	Id            string `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid();column:Id"`
	AttributeID   string `json:"attribute_id" gorm:"type:uuid;column:AttributeID;uniqueIndex:,composite:attributeid_producttypeid_key"`      // to Attribute
	ProductTypeID string `json:"product_type_id" gorm:"type:uuid;column:ProductTypeID;uniqueIndex:,composite:attributeid_producttypeid_key"` // to ProductType
	Sortable

	Attribute          *Attribute                  `json:"-"`
	AssignedProducts   Products                    `json:"-" gorm:"many2many:AssignedProductAttributes"`
	ProductAssignments []*AssignedProductAttribute `json:"-" gorm:"foreignKey:AssignmentID;constraint:OnDelete:CASCADE;"`
}

func (*AttributeProduct) TableName() string               { return AttributeProductTableName }
func (a *AttributeProduct) BeforeCreate(_ *gorm.DB) error { return a.IsValid() }
func (a *AttributeProduct) BeforeUpdate(_ *gorm.DB) error { return a.IsValid() }

// AttributeProductFilterOption is used when finding an attributeProduct.
type AttributeProductFilterOption struct {
	Conditions                   squirrel.Sqlizer
	AttributeVisibleInStoreFront *bool
}

func (a *AttributeProduct) IsValid() *AppError {
	outer := CreateAppErrorForModel(
		"model.attribute_product.is_valid.%s.app_error",
		"attribute_product_id=",
		"AttributeProduct.IsValid",
	)
	if !IsValidId(a.AttributeID) {
		return outer("attribute_id", &a.Id)
	}
	if !IsValidId(a.ProductTypeID) {
		return outer("product_type_id", &a.Id)
	}

	return nil
}

func (a *AttributeProduct) DeepCopy() *AttributeProduct {
	if a == nil {
		return nil
	}
	res := *a
	if a.SortOrder != nil {
		res.SortOrder = NewPrimitive(*a.SortOrder)
	}
	if a.Attribute != nil {
		res.Attribute = a.Attribute.DeepCopy()
	}
	return &res
}

// Associate a product type attribute and selected values to a given product
// ProductID unique with AssignmentID
type AssignedProductAttribute struct {
	Id           string `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid();column:Id"`
	ProductID    string `json:"product_id" gorm:"type:uuid;uniqueIndex:,composite:productid_assignmentid_key;column:ProductID"`       // to Product
	AssignmentID string `json:"assignment_id" gorm:"type:uuid;uniqueIndex:,composite:productid_assignmentid_key;column:AssignmentID"` // to AttributeProduct

	Values                  AttributeValues                  `json:"-" gorm:"many2many:AssignedProductAttributeValues"`
	ProductValueAssignments []*AssignedProductAttributeValue `json:"-" gorm:"foreignKey:AssignmentID;constraint:OnDelete:CASCADE;"`
}

func (*AssignedProductAttribute) TableName() string               { return AssignedProductAttributeTableName }
func (a *AssignedProductAttribute) BeforeCreate(_ *gorm.DB) error { return a.IsValid() }
func (a *AssignedProductAttribute) BeforeUpdate(_ *gorm.DB) error { return a.IsValid() }

type AssignedProductAttributes []*AssignedProductAttribute

// AssignedProductAttributeFilterOption is used to filter or creat new AssignedProductAttribute
type AssignedProductAttributeFilterOption struct {
	Conditions squirrel.Sqlizer

	AttributeProduct_Attribute_VisibleInStoreFront *bool // INNER JOIN AttributeProduct ON ... INNER JOIN Attribute ON ... WHERE Attribute.VisibleInStoreFront ...
}

func (a *AssignedProductAttribute) IsValid() *AppError {
	outer := CreateAppErrorForModel(
		"model.assigned_product_attribute.is_valid.%s.app_error",
		"assigned_product_attribute_id=",
		"AssignedProductAttribute.IsValid",
	)
	if !IsValidId(a.ProductID) {
		return outer("product_id", &a.Id)
	}
	if !IsValidId(a.AssignmentID) {
		return outer("assignment_id", &a.Id)
	}

	return nil
}

func (a *AssignedProductAttribute) DeepCopy() *AssignedProductAttribute {
	if a == nil {
		return nil
	}

	res := *a
	res.Values = a.Values.DeepCopy()
	// if a.attributeProduct != nil {
	// 	res.attributeProduct = a.attributeProduct.DeepCopy()
	// }
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
	ValueID      string `json:"value_id" gorm:"type:uuid;uniqueIndex:,composite:valueid_assignmentid_key;column:ValueID"`           // to AttributeValue
	AssignmentID string `json:"assignment_id" gorm:"type:uuid;uniqueIndex:,composite:valueid_assignmentid_key;column:AssignmentID"` // to AssignedProductAttribute
	Sortable
}

func (*AssignedProductAttributeValue) TableName() string {
	return AssignedProductAttributeValueTableName
}
func (a *AssignedProductAttributeValue) BeforeCreate(_ *gorm.DB) error { return a.IsValid() }
func (a *AssignedProductAttributeValue) BeforeUpdate(_ *gorm.DB) error { return a.IsValid() }

type AssignedProductAttributeValueFilterOptions struct {
	Conditions squirrel.Sqlizer
}

func (a *AssignedProductAttributeValue) IsValid() *AppError {
	outer := CreateAppErrorForModel(
		"model.assigned_product_attribute.is_valid.%s.app_error",
		"assigned_product_attribute_id=",
		"AssignedProductAttributeValue.IsValid",
	)
	if !IsValidId(a.ValueID) {
		return outer("value_id", nil)
	}
	if !IsValidId(a.AssignmentID) {
		return outer("assignment_id", nil)
	}

	return nil
}

func (a *AssignedProductAttributeValue) DeepCopy() *AssignedProductAttributeValue {
	res := *a

	if a.SortOrder != nil {
		res.SortOrder = NewPrimitive(*a.SortOrder)
	}
	return &res
}
