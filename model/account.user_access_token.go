package model

import (
	"io"
)

type UserAccessToken struct {
	Id          string `json:"id"`
	Token       string `json:"token,omitempty"`
	UserId      string `json:"user_id"`
	Description string `json:"description"`
	IsActive    bool   `json:"is_active"`
}

const (
	USER_ACCESS_TOKEN_DESCRIPTION_MAX_LENGTH = 255
)

func (t *UserAccessToken) IsValid() *AppError {
	outer := CreateAppErrorForModel(
		"user_access_token.is_valid.%s.app_error",
		"user_access_token_id=",
		"UserAccessToken.IsValid",
	)
	if !IsValidId(t.Id) {
		return outer("id", nil)
	}
	if !IsValidId(t.Token) {
		return outer("token", &t.Id)
	}
	if !IsValidId(t.UserId) {
		return outer("user_id", &t.Id)
	}
	if len(t.Description) > USER_ACCESS_TOKEN_DESCRIPTION_MAX_LENGTH {
		return outer("description", &t.Id)
	}

	return nil
}

func (t *UserAccessToken) PreSave() {
	if t.Id == "" {
		t.Id = NewId()
	}
	if t.Token == "" {
		t.Token = NewId()
	}
	t.IsActive = true
}

func (t *UserAccessToken) ToJSON() string {
	return ModelToJson(t)
}

func UserAccessTokenFromJson(data io.Reader) *UserAccessToken {
	var t UserAccessToken
	ModelFromJson(&t, data)
	return &t
}

func UserAccessTokenListToJson(t []*UserAccessToken) string {
	return ModelToJson(&t)
}

func UserAccessTokenListFromJson(data io.Reader) []*UserAccessToken {
	var t []*UserAccessToken
	ModelFromJson(&t, data)
	return t
}
