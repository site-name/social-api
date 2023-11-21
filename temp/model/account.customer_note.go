package model

import (
	"net/http"

	"gorm.io/gorm"
)

type CustomerNote struct {
	Id         string  `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid();column:Id"`
	UserID     *string `json:"user_id" gorm:"type:uuid;index:customernotes_userid_key;column:UserID"`
	Date       int64   `json:"date" gorm:"autoCreateTime:milli;column:Date"` // default now()
	Content    string  `json:"content" gorm:"type:text;column:Content"`
	IsPublic   *bool   `json:"is_public" gorm:"default:true;column:IsPublic"`
	CustomerID string  `json:"customer_id" gorm:"type:uuid;index:customernotes_customerid_key;column:CustomerID"`
}

func (c *CustomerNote) BeforeCreate(_ *gorm.DB) error {
	return c.IsValid()
}

func (c *CustomerNote) BeforeUpdate(_ *gorm.DB) error {
	return c.IsValid()
}

func (*CustomerNote) TableName() string {
	return CustomerNoteTableName
}

func (c *CustomerNote) IsValid() *AppError {
	if c.UserID != nil && !IsValidId(*c.UserID) {
		return NewAppError("CustomerNote.IsValid", "model.customer_note.is_valid.user_id.app_error", nil, "please provide valid user id", http.StatusBadRequest)
	}
	if !IsValidId(c.CustomerID) {
		return NewAppError("CustomerNote.IsValid", "model.customer_note.is_valid.customer_id.app_error", nil, "please provide valid customer id", http.StatusBadRequest)
	}
	return nil
}
