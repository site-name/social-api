package model

import (
	"github.com/Masterminds/squirrel"
	"gorm.io/gorm"
)

// AttributeID unique together with ProductTypeID
type AttributeVariant struct {
	Id            string `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid();column:Id"`
	AttributeID   string `json:"attribute_id" gorm:"type:uuid;uniqueIndex:,composite:attributeid_producttypeid_key;column:AttributeID"`
	ProductTypeID string `json:"product_type_id" gorm:"type:uuid;uniqueIndex:,composite:attributeid_producttypeid_key;column:ProductTypeID"`
	Sortable

	AssignedVariants   ProductVariants             `json:"-" gorm:"many2many:AssignedVariantAttributes"`
	VariantAssignments []*AssignedVariantAttribute `json:"-" gorm:"foreignKey:AssignmentID"`
}

func (a *AttributeVariant) BeforeCreate(_ *gorm.DB) error { return a.IsValid() }
func (a *AttributeVariant) BeforeUpdate(_ *gorm.DB) error { return a.IsValid() }
func (a *AttributeVariant) TableName() string             { return AttributeVariantTableName }

// AttributeVariantFilterOption is used to find `AttributeVariant`.
type AttributeVariantFilterOption struct {
	Conditions                   squirrel.Sqlizer
	AttributeVisibleInStoreFront *bool
}

func (a *AttributeVariant) IsValid() *AppError {
	outer := CreateAppErrorForModel(
		"model.attribute_variant.is_valid.%s.app_error",
		"attribute_variant_id=",
		"AttributeVariant.IsValid",
	)
	if !IsValidId(a.AttributeID) {
		return outer("attribute_id", &a.Id)
	}
	if !IsValidId(a.ProductTypeID) {
		return outer("product_type_id", &a.Id)
	}

	return nil
}

// Associate a product type attribute and selected values to a given variant.
type AssignedVariantAttribute struct {
	VariantID    string `json:"variant_id" gorm:"primaryKey;column:VariantID;uniqueIndex:,composite:variantid_assignmentid_key"`       // to ProductVariant
	AssignmentID string `json:"assignment_id" gorm:"primaryKey;column:AssignmentID;uniqueIndex:,composite:variantid_assignmentid_key"` // to AttributeVariant

	Values                 AttributeValues                  `json:"-" gorm:"many2many:AssignedVariantAttributeValues"`
	VariantValueAssignment []*AssignedVariantAttributeValue `json:"-" gorm:"foreignKey:AssignmentID"`
}

func (a *AssignedVariantAttribute) BeforeCreate(_ *gorm.DB) error { return a.IsValid() }
func (a *AssignedVariantAttribute) BeforeUpdate(_ *gorm.DB) error { return a.IsValid() }
func (a *AssignedVariantAttribute) TableName() string             { return AssignedVariantAttributeTableName }

// AssignedVariantAttributeFilterOption is used for lookup, if cannot found, creating new instance
type AssignedVariantAttributeFilterOption struct {
	Conditions squirrel.Sqlizer

	Assignment_Attribute_VisibleInStoreFront *bool // INNER JOIN AttributeVariant ON ... INNER JOIN Attributes ON ... WHERE Attributes.VisibleInStoreFront ...

	AssignmentAttributeInputType squirrel.Sqlizer // INNER JOIN AttributeVariants ON () INNER JOIN Attributes ON () WHERE Attributes.InputType
	AssignmentAttributeType      squirrel.Sqlizer // INNER JOIN AttributeVariants ON () INNER JOIN Attributes ON () WHERE Attributes.Type
}

func (a *AssignedVariantAttribute) IsValid() *AppError {
	outer := CreateAppErrorForModel(
		"model.assigned_variant_attribute.is_valid.%s.app_error",
		"assigned_variant_attribute_id=",
		"AssignedVariantAttribute.IsValid",
	)
	if !IsValidId(a.VariantID) {
		return outer("variant_id", nil)
	}
	if !IsValidId(a.AssignmentID) {
		return outer("assignment_id", nil)
	}

	return nil
}

// ValueID unique together with AssignmentID
type AssignedVariantAttributeValue struct {
	ValueID      string `json:"value_id" gorm:"primaryKey;uniqueIndex:,composite:valueid_assignmentid_key;type:uuid;column:ValueID"`
	AssignmentID string `json:"assignment_id" gorm:"primaryKey;uniqueIndex:,composite:valueid_assignmentid_key;type:uuid;column:AssignmentID"`
	Sortable
}

func (a *AssignedVariantAttributeValue) BeforeCreate(_ *gorm.DB) error { return a.IsValid() }
func (a *AssignedVariantAttributeValue) BeforeUpdate(_ *gorm.DB) error { return a.IsValid() }
func (a *AssignedVariantAttributeValue) TableName() string {
	return AssignedVariantAttributeValueTableName
}

type AssignedVariantAttributeValueFilterOptions struct {
	Conditions squirrel.Sqlizer
}

func (a *AssignedVariantAttributeValue) IsValid() *AppError {
	outer := CreateAppErrorForModel(
		"model.assigned_variant_attribute_value.is_valid.%s.app_error",
		"assigned_variant_attribute_value_id=",
		"AssignedVariantAttributeValue.IsValid",
	)

	if !IsValidId(a.ValueID) {
		return outer("value_id", nil)
	}
	if !IsValidId(a.AssignmentID) {
		return outer("assignment_id", nil)
	}

	return nil
}

func (a *AssignedVariantAttributeValue) DeepCopy() *AssignedVariantAttributeValue {
	res := *a
	if a.SortOrder != nil {
		res.SortOrder = NewPrimitive(*a.SortOrder)
	}
	return &res
}
