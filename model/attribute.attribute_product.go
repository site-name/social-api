package model

import (
	"net/http"

	"github.com/Masterminds/squirrel"
	"gorm.io/gorm"
)

// AttributeID unique together with ProductTypeID
type AttributeProduct struct {
	Id            UUID `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid();column:Id"`
	AttributeID   UUID `json:"attribute_id" gorm:"type:uuid;column:AttributeID;index:attributeid_producttypeid_key"`      // to Attribute
	ProductTypeID UUID `json:"product_type_id" gorm:"type:uuid;column:ProductTypeID;index:attributeid_producttypeid_key"` // to ProductType
	Sortable

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
	if !IsValidId(a.AttributeID) {
		return NewAppError("AttributeProduct.IsValid", "model.attribute_product.is_valid.attribute_id.app_error", nil, "", http.StatusBadRequest)
	}
	if !IsValidId(a.ProductTypeID) {
		return NewAppError("AttributeProduct.IsValid", "model.attribute_product.is_valid.product_type_id.app_error", nil, "", http.StatusBadRequest)
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
	return &res
}

// Associate a product type attribute and selected values to a given product
// ProductID unique with AssignmentID
type AssignedProductAttribute struct {
	Id           UUID `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid();column:Id"`
	ProductID    UUID `json:"product_id" gorm:"type:uuid;index:productid_assignmentid_key;column:ProductID"`       // to Product
	AssignmentID UUID `json:"assignment_id" gorm:"type:uuid;index:productid_assignmentid_key;column:AssignmentID"` // to AttributeProduct

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
	if !IsValidId(a.ProductID) {
		return NewAppError("AssignedProductAttribute.IsValid", "model.assigned_product_attribute.is_valid.product_id.app_error", nil, "", http.StatusBadRequest)
	}
	if !IsValidId(a.AssignmentID) {
		return NewAppError("AssignedProductAttribute.IsValid", "model.assigned_product_attribute.is_valid.assignment_id.app_error", nil, "", http.StatusBadRequest)
	}

	return nil
}

func (a *AssignedProductAttribute) DeepCopy() *AssignedProductAttribute {
	if a == nil {
		return nil
	}

	res := *a
	res.Values = a.Values.DeepCopy()
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
	Id           UUID `json:"id" gorm:"type:uuid;primaryKey;column:Id;default:gen_random_uuid()"`
	ValueID      UUID `json:"value_id" gorm:"type:uuid;index:valueid_assignmentid_key;column:ValueID"`           // to AttributeValue
	AssignmentID UUID `json:"assignment_id" gorm:"type:uuid;index:valueid_assignmentid_key;column:AssignmentID"` // to AssignedProductAttribute
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
	if !IsValidId(a.ValueID) {
		return NewAppError("AssignedProductAttributeValue.IsValid", "model.assigned_product_attribute.is_valid.value_id.app_error", nil, "", http.StatusBadRequest)
	}
	if !IsValidId(a.AssignmentID) {
		return NewAppError("AssignedProductAttributeValue.IsValid", "model.assigned_product_attribute.is_valid.value_id.app_error", nil, "", http.StatusBadRequest)
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
