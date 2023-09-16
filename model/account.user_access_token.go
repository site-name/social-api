package model

import "gorm.io/gorm"

type UserAccessToken struct {
	Id          string `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid();column:Id"`
	Token       string `json:"token,omitempty" gorm:"type:uuid;column:Token"`
	UserId      string `json:"user_id" gorm:"type:uuid;index:useraccesstokens_userid_key;column:UserId"`
	Description string `json:"description" gorm:"type:varchar(255);column:Description"`
	IsActive    *bool  `json:"is_active" gorm:"default:true;column:IsActive"` // defaut true
}

func (c *UserAccessToken) BeforeCreate(_ *gorm.DB) error { c.commonPre(); return c.IsValid() }
func (c *UserAccessToken) BeforeUpdate(_ *gorm.DB) error { c.commonPre(); return c.IsValid() }
func (*UserAccessToken) TableName() string               { return UserAccessTokenTableName }

func (t *UserAccessToken) IsValid() *AppError {
	outer := CreateAppErrorForModel(
		"model.user_access_token.is_valid.%s.app_error",
		"user_access_token_id=",
		"UserAccessToken.IsValid",
	)
	if !IsValidId(t.Token) {
		return outer("token", &t.Id)
	}
	if !IsValidId(t.UserId) {
		return outer("user_id", &t.Id)
	}
	return nil
}

func (t *UserAccessToken) commonPre() {
	if t.Token == "" {
		t.Token = NewId()
	}
	if t.IsActive == nil {
		t.IsActive = GetPointerOfValue(true)
	}
}
