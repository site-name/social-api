package model

import (
	"net/http"

	"github.com/mattermost/squirrel"
	"gorm.io/gorm"
)

// AttributeID unique with PageTypeID
type AttributePage struct {
	Id          string `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid();column:Id"`
	AttributeID string `json:"attribute_id" gorm:"type:uuid;index:attributeid_pagetypeid_key;column:AttributeID"` // to Attribute
	PageTypeID  string `json:"page_type_id" gorm:"type:uuid;index:attributeid_pagetypeid_key;column:PageTypeID"`  // to PageType
	Sortable

	PageAssignments []*AssignedPageAttribute `json:"-" gorm:"foreignKey:AssignmentID;constraint:OnDelete:CASCADE;"`
	AssignedPages   []*Page                  `json:"-" gorm:"many2many:AssignedPageAttributes"` // through AssignedPageAttribute
	Attribute       *Attribute               `json:"-" gorm:"foreignKey:AttributeID;constraint:OnDelete:CASCADE;"`
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
	Id           string `json:"id" gorm:"type:uuid;primaryKey;column:Id;default:gen_random_uuid()"`
	ValueID      string `json:"value_id" gorm:"primeryKey;type:uuid;column:ValueID;index:valueid_assignmentid_key"`           // AttributeValue
	AssignmentID string `json:"assignment_id" gorm:"primeryKey;type:uuid;column:AssignmentID;index:valueid_assignmentid_key"` // AssignedPageAttribute
	Sortable

	AssignedPageAttribute *AssignedPageAttribute `json:"-" gorm:"foreignKey:AssignmentID;constraint:OnDelete:CASCADE;"`
	AttributeValue        *AttributeValue        `json:"-" gorm:"foreignKey:ValueID;constraint:OnDelete:CASCADE;"`
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
		res.SortOrder = GetPointerOfValue(*a.SortOrder)
	}
	return &res
}

// Associate a page type attribute and selected values to a given page.
// PageID unique together with AssignmentID
type AssignedPageAttribute struct {
	Id           string `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid();column:Id"`
	PageID       string `json:"page_id" gorm:"type:uuid;column:PageID;index:pageid_assignmentid_key"`             // Page
	AssignmentID string `json:"assignment_id" gorm:"type:uuid;column:AssignmentID;index:pageid_assignmentid_key"` // AttributePage

	PageValueAssignments []*AssignedPageAttributeValue `json:"-" gorm:"foreignKey:AssignmentID;constraint:OnDelete:CASCADE;"`
	Values               AttributeValues               `json:"-" gorm:"many2many:AssignedPageAttributeValues"`
	AttributePage        *AttributePage                `json:"-" gorm:"foreignKey:AssignmentID;constraint:OnDelete:CASCADE;"`
}

func (a *AssignedPageAttribute) BeforeCreate(_ *gorm.DB) error { return a.IsValid() }
func (a *AssignedPageAttribute) BeforeUpdate(_ *gorm.DB) error { return a.IsValid() }
func (a *AssignedPageAttribute) TableName() string             { return AssignedPageAttributeTableName }

// AssignedPageAttributeFilterOption is used to find or creat new AssignedPageAttribute
type AssignedPageAttributeFilterOption struct {
	Conditions squirrel.Sqlizer
}

func (a *AssignedPageAttribute) IsValid() *AppError {
	if !IsValidId(a.PageID) {
		return NewAppError("AssignedPageAttribute.IsValid", "model.assigned_page_attribute.is_valid.page_id.app_error", nil, "", http.StatusBadRequest)
	}
	if !IsValidId(a.AssignmentID) {
		return NewAppError("AssignedPageAttribute.IsValid", "model.assigned_page_attribute.is_valid.assignment_id.app_error", nil, "", http.StatusBadRequest)
	}

	return nil
}

func (a *AssignedPageAttribute) DeepCopy() *AssignedPageAttribute {
	res := *a
	return &res
}
