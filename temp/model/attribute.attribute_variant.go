package model

import (
	"net/http"

	"github.com/Masterminds/squirrel"
	"gorm.io/gorm"
)

// AttributeID unique together with ProductTypeID
type AttributeVariant struct {
	Id            string `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid();column:Id"`
	AttributeID   string `json:"attribute_id" gorm:"type:uuid;index:attributeid_producttypeid_key;column:AttributeID"`
	ProductTypeID string `json:"product_type_id" gorm:"type:uuid;index:attributeid_producttypeid_key;column:ProductTypeID"`
	Sortable

	AssignedVariants   ProductVariants             `json:"-" gorm:"many2many:AssignedVariantAttributes"`
	VariantAssignments []*AssignedVariantAttribute `json:"-" gorm:"foreignKey:AssignmentID"`
}

// column names for table attribute variants
const (
	AttributeVariantColumnId          = "Id"
	AttributeVariantColumnAttributeID = "AttributeID"
	AttributeVariantProductTypeID     = "ProductTypeID"
)

func (a *AttributeVariant) BeforeCreate(_ *gorm.DB) error { return a.IsValid() }
func (a *AttributeVariant) BeforeUpdate(_ *gorm.DB) error { return a.IsValid() }
func (a *AttributeVariant) TableName() string             { return AttributeVariantTableName }

// AttributeVariantFilterOption is used to find `AttributeVariant`.
type AttributeVariantFilterOption struct {
	Conditions                   squirrel.Sqlizer
	AttributeVisibleInStoreFront *bool // INNER JOIN Attributes ON ... WHERE Attributes.VisibleInStoreFront ...
}

func (a *AttributeVariant) IsValid() *AppError {
	if !IsValidId(a.AttributeID) {
		return NewAppError("AttributeVariant.IsValid", "model.attribute_variant.is_valid.attribute_id.app_error", nil, "", http.StatusBadRequest)
	}
	if !IsValidId(a.ProductTypeID) {
		return NewAppError("AttributeVariant.IsValid", "model.attribute_variant.is_valid.product_type_id.app_error", nil, "", http.StatusBadRequest)
	}

	return nil
}

// Associate a product type attribute and selected values to a given variant.
type AssignedVariantAttribute struct {
	Id           string `json:"id" gorm:"column:Id;type:uuid;primaryKey;default:gen_random_uuid()"`
	VariantID    string `json:"variant_id" gorm:"column:VariantID;index:variantid_assignmentid_key"`       // to ProductVariant
	AssignmentID string `json:"assignment_id" gorm:"column:AssignmentID;index:variantid_assignmentid_key"` // to AttributeVariant

	Values                 AttributeValues                  `json:"-" gorm:"many2many:AssignedVariantAttributeValues"`
	VariantValueAssignment []*AssignedVariantAttributeValue `json:"-" gorm:"foreignKey:AssignmentID"`
}

// column names for table AssignedVariantAttribute
const (
	AssignedVariantAttributeColumnId           = "Id"
	AssignedVariantAttributeColumnVariantID    = "VariantID"
	AssignedVariantAttributeColumnAssignmentID = "AssignmentID"
)

func (a *AssignedVariantAttribute) BeforeCreate(_ *gorm.DB) error { return a.IsValid() }
func (a *AssignedVariantAttribute) BeforeUpdate(_ *gorm.DB) error { return a.IsValid() }
func (a *AssignedVariantAttribute) TableName() string             { return AssignedVariantAttributeTableName }

// AssignedVariantAttributeFilterOption is used for lookup, if cannot found, creating new instance
type AssignedVariantAttributeFilterOption struct {
	Conditions squirrel.Sqlizer

	Assignment_Conditions           squirrel.Sqlizer // INNER JOIN AttributeVariants ON ... WHERE AttributeVariants...
	Assignment_Attribute_Conditions squirrel.Sqlizer // INNER JOIN AttributeVariants ON ... INNER JOIN Attributes ON ... WHERE Attributes...

	// E.g:
	//  "Values"
	Preloads []string
}

func (a *AssignedVariantAttribute) IsValid() *AppError {
	if !IsValidId(a.VariantID) {
		return NewAppError("AssignedVariantAttribute.IsValid", "model.assigned_variant_attribute.is_valid.variant_id.app_error", nil, "", http.StatusBadRequest)
	}
	if !IsValidId(a.AssignmentID) {
		return NewAppError("AssignedVariantAttribute.IsValid", "model.assigned_variant_attribute.is_valid.assignment_id.app_error", nil, "", http.StatusBadRequest)
	}

	return nil
}

// ValueID unique together with AssignmentID
type AssignedVariantAttributeValue struct {
	Id           string `json:"id" gorm:"column:Id;type:uuid;primaryKey;default:gen_random_uuid()"`
	ValueID      string `json:"value_id" gorm:"column:ValueID;primaryKey;index:valueid_assignmentid_key;type:uuid"`
	AssignmentID string `json:"assignment_id" gorm:"column:AssignmentID;primaryKey;index:valueid_assignmentid_key;type:uuid"`
	Sortable

	AttributeValue *AttributeValue `json:"-" gorm:"foreignKey:ValueID;constraint:OnDelete:CASCADE;"`
}

// column names for table AssignedVariantAttributeValue
const (
	AssignedVariantAttributeValueColumnId           = "Id"
	AssignedVariantAttributeValueColumnValueID      = "ValueID"
	AssignedVariantAttributeValueColumnAssignmentID = "AssignmentID"
)

func (a *AssignedVariantAttributeValue) BeforeCreate(_ *gorm.DB) error { return a.IsValid() }
func (a *AssignedVariantAttributeValue) BeforeUpdate(_ *gorm.DB) error { return a.IsValid() }
func (*AssignedVariantAttributeValue) TableName() string {
	return AssignedVariantAttributeValueTableName
}

type AssignedVariantAttributeValueFilterOptions struct {
	Conditions squirrel.Sqlizer
}

func (a *AssignedVariantAttributeValue) IsValid() *AppError {
	if !IsValidId(a.ValueID) {
		return NewAppError("AssignedVariantAttributeValue.IsValid", "model.assigned_variant_attribute_value.is_valid.value_id.app_error", nil, "", http.StatusBadRequest)
	}
	if !IsValidId(a.AssignmentID) {
		return NewAppError("AssignedVariantAttributeValue.IsValid", "model.assigned_variant_attribute_value.is_valid.assignment_id.app_error", nil, "", http.StatusBadRequest)
	}

	return nil
}

func (a *AssignedVariantAttributeValue) DeepCopy() *AssignedVariantAttributeValue {
	res := *a
	res.SortOrder = CopyPointer(a.SortOrder)
	return &res
}
