package model

import (
	"unicode/utf8"

	"gorm.io/gorm"
)

const (
	TERMS_OF_SERVICE_CACHE_SIZE = 1
	POST_MESSAGE_MAX_BYTES_V2   = 65535                         // Maximum size of a TEXT column in MySQL
	POST_MESSAGE_MAX_RUNES_V2   = POST_MESSAGE_MAX_BYTES_V2 / 4 // Assume a worst-case representation
)

type TermsOfService struct {
	Id       string `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid();column:Id"`
	CreateAt int64  `json:"create_at" gorm:"autoCreateTime:milli;type:bigint;column:CreateAt"`
	UserId   string `json:"user_id" gorm:"type:uuid;index:termofservices_userid_key;column:UserId"`
	Text     string `json:"text" gorm:"column:Text;type:varchar(16383)"`
}

func (t *TermsOfService) BeforeCreate(_ *gorm.DB) error {
	return t.IsValid()
}

func (t *TermsOfService) BeforeUpdate(_ *gorm.DB) error {
	return t.IsValid()
}

func (t *TermsOfService) TableName() string {
	return TermsOfServiceTableName
}

func (t *TermsOfService) IsValid() *AppError {
	outer := CreateAppErrorForModel(
		"model.terms_of_service.is_valid.%s.app_error",
		"terms_of_service_id=",
		"TermsOfService.IsValid",
	)
	if !IsValidId(t.UserId) {
		return outer("user_id", &t.Id)
	}
	if utf8.RuneCountInString(t.Text) > POST_MESSAGE_MAX_RUNES_V2 {
		return outer("text", &t.Id)
	}

	return nil
}
