package model

import (
	"net/http"

	"gorm.io/gorm"
)

type TermsOfService struct {
	Id       UUID   `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid();column:Id"`
	CreateAt int64  `json:"create_at" gorm:"autoCreateTime:milli;type:bigint;column:CreateAt"`
	UserId   UUID   `json:"user_id" gorm:"type:uuid;index:termofservices_userid_key;column:UserId"`
	Text     string `json:"text" gorm:"column:Text;type:varchar(16383)"`
}

func (t *TermsOfService) BeforeCreate(_ *gorm.DB) error { return t.IsValid() }
func (t *TermsOfService) BeforeUpdate(_ *gorm.DB) error { return t.IsValid() }
func (t *TermsOfService) TableName() string             { return TermsOfServiceTableName }

func (t *TermsOfService) IsValid() *AppError {
	if !IsValidId(t.UserId) {
		return NewAppError("TermOfService.IsValid", "model.terms_of_service.is_valid.user_id.app_error", nil, "please provide valid user id", http.StatusBadRequest)
	}

	return nil
}
