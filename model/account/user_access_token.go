package account

import (
	"io"

	"github.com/sitename/sitename/model"
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

func (t *UserAccessToken) IsValid() *model.AppError {
	outer := model.CreateAppErrorForModel(
		"model.user_access_token.is_valid.%s.app_error",
		"user_access_token_id=",
		"UserAccessToken.IsValid",
	)
	if !model.IsValidId(t.Id) {
		return outer("id", nil)
	}
	if !model.IsValidId(t.Token) {
		return outer("token", &t.Id)
	}
	if !model.IsValidId(t.UserId) {
		return outer("user_id", &t.Id)
	}
	if len(t.Description) > USER_ACCESS_TOKEN_DESCRIPTION_MAX_LENGTH {
		return outer("description", &t.Id)
	}

	return nil
}

func (t *UserAccessToken) PreSave() {
	if t.Id == "" {
		t.Id = model.NewId()
	}
	if t.Token == "" {
		t.Token = model.NewId()
	}
	t.IsActive = true
}

func (t *UserAccessToken) ToJson() string {
	return model.ModelToJson(t)
}

func UserAccessTokenFromJson(data io.Reader) *UserAccessToken {
	var t UserAccessToken
	model.ModelFromJson(&t, data)
	return &t
}

func UserAccessTokenListToJson(t []*UserAccessToken) string {
	return model.ModelToJson(&t)
}

func UserAccessTokenListFromJson(data io.Reader) []*UserAccessToken {
	var t []*UserAccessToken
	model.ModelFromJson(&t, data)
	return t
}
