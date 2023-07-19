package model

import (
	"net/http"

	"github.com/Masterminds/squirrel"
	"gorm.io/gorm"
)

// AttributeID unique with PageTypeID
type AttributePage struct {
	Id          string `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid();column:Id"`
	AttributeID string `json:"attribute_id" gorm:"type:uuid;uniqueIndex:,composite:attributeid_pagetypeid_key;column:AttributeID"` // to Attribute
	PageTypeID  string `json:"page_type_id" gorm:"type:uuid;uniqueIndex:,composite:attributeid_pagetypeid_key;column:PageTypeID"`  // to PageType
	Sortable

	PageAssignments []*AssignedPageAttribute `json:"-" gorm:"foreignKey:AssignmentID;constraint:OnDelete:CASCADE;"`
	AssignedPages   []*Page                  `json:"-" gorm:"many2many:AssignedPageAttributes"`
}

func (*AttributePage) TableName() string               { return AttributePageTableName }
func (a *AttributePage) BeforeCreate(_ *gorm.DB) error { return a.IsValid() }
func (a *AttributePage) BeforeUpdate(_ *gorm.DB) error { return a.IsValid() }

// AttributePageFilterOption is used for lookup AttributePage
type AttributePageFilterOption struct {
	Conditions squirrel.Sqlizer
}

func (a *AttributePage) IsValid() *AppError {
	if !IsValidId(a.AttributeID) {
		return NewAppError("AttributePage.IsValid", "model.attribute_page.is_valid.attribute_id.app_error", nil, "please provide valid attribute id", http.StatusBadRequest)
	}
	if !IsValidId(a.PageTypeID) {
		return NewAppError("AttributePage.IsValid", "model.attribute_page.is_valid.page_type_id.app_error", nil, "please provide valid page type id", http.StatusBadRequest)
	}
	return nil
}

// ValueID unique together with AssignmentID
type AssignedPageAttributeValue struct {
	ValueID      string `json:"value_id" gorm:"primeryKey;type:uuid;column:ValueID;uniqueIndex:,composite:valueid_assignmentid_key"`           // AttributeValue
	AssignmentID string `json:"assignment_id" gorm:"primeryKey;type:uuid;column:AssignmentID;uniqueIndex:,composite:valueid_assignmentid_key"` // AssignedPageAttribute
	Sortable
}

func (a *AssignedPageAttributeValue) TableName() string             { return AssignedPageAttributeValueTableName }
func (a *AssignedPageAttributeValue) BeforeCreate(_ *gorm.DB) error { return a.IsValid() }
func (a *AssignedPageAttributeValue) BeforeUpdate(_ *gorm.DB) error { return a.IsValid() }

func (a *AssignedPageAttributeValue) IsValid() *AppError {
	if !IsValidId(a.ValueID) {
		return NewAppError("AssignedPageAttributeValue.IsValid", "model.assigned_page_attribute_value.is_valid.value_id.app_error", nil, "please provide valid value id", http.StatusBadRequest)
	}
	if !IsValidId(a.AssignmentID) {
		return NewAppError("AssignedPageAttributeValue.IsValid", "model.assigned_page_attribute_value.is_valid.assignment_id.app_error", nil, "please provide valid assignment id", http.StatusBadRequest)
	}
	return nil
}

func (a *AssignedPageAttributeValue) DeepCopy() *AssignedPageAttributeValue {
	res := *a

	if a.SortOrder != nil {
		res.SortOrder = NewPrimitive(*a.SortOrder)
	}
	return &res
}

// Associate a page type attribute and selected values to a given page.
// PageID unique together with AssignmentID
type AssignedPageAttribute struct {
	Id           string `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid();column:Id"`
	PageID       string `json:"page_id" gorm:"type:uuid;column:PageID;uniqueIndex:,composite:pageid_assignmentid_key"`             // Page
	AssignmentID string `json:"assignment_id" gorm:"type:uuid;column:AssignmentID;uniqueIndex:,composite:pageid_assignmentid_key"` // AttributePage

	PageValueAssignments []*AssignedPageAttributeValue `json:"-" gorm:"foreignKey:AssignmentID;constraint:OnDelete:CASCADE;"`
	Values               AttributeValues               `json:"-" gorm:"many2many:AssignedPageAttributeValues"`
}

func (a *AssignedPageAttribute) BeforeCreate(_ *gorm.DB) error { return a.IsValid() }
func (a *AssignedPageAttribute) BeforeUpdate(_ *gorm.DB) error { return a.IsValid() }
func (a *AssignedPageAttribute) TableName() string             { return AssignedPageAttributeTableName }

// AssignedPageAttributeFilterOption is used to find or creat new AssignedPageAttribute
type AssignedPageAttributeFilterOption struct {
	Conditions squirrel.Sqlizer
}

func (a *AssignedPageAttribute) IsValid() *AppError {
	outer := CreateAppErrorForModel(
		"model.assigned_page_attribute.is_valid.%s.app_error",
		"assigned_page_attribute_id=",
		"AssignedPageAttribute.IsValid",
	)
	if !IsValidId(a.PageID) {
		return outer("page_id", &a.Id)
	}
	if !IsValidId(a.AssignmentID) {
		return outer("assignment_id", &a.Id)
	}

	return nil
}

func (a *AssignedPageAttribute) DeepCopy() *AssignedPageAttribute {
	res := *a
	return &res
}
