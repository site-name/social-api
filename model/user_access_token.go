package model

import (
	"io"
	"net/http"

	"github.com/google/uuid"
	"github.com/sitename/sitename/modules/json"
)

type UserAccessToken struct {
	Id          string `json:"id"`
	Token       string `json:"token,omitempty"`
	UserId      string `json:"user_id"`
	Description string `json:"description"`
	IsActive    bool   `json:"is_active"`
}

const (
	USER_ACCESS_TOKEN_MAX_LENGTH             = 26
	USER_ACCESS_TOKEN_DESCRIPTION_MAX_LENGTH = 255
)

func (t *UserAccessToken) IsValid() *AppError {
	if !IsValidId(t.Id) {
		return NewAppError("UserAccessToken.IsValid", "model.user_access_token.is_valid.id.app_error", nil, "", http.StatusBadRequest)
	}

	if len(t.Token) != USER_ACCESS_TOKEN_MAX_LENGTH {
		return NewAppError("UserAccessToken.IsValid", "model.user_access_token.is_valid.token.app_error", nil, "", http.StatusBadRequest)
	}

	if !IsValidId(t.UserId) {
		return NewAppError("UserAccessToken.IsValid", "model.user_access_token.is_valid.user_id.app_error", nil, "", http.StatusBadRequest)
	}

	if len(t.Description) > USER_ACCESS_TOKEN_DESCRIPTION_MAX_LENGTH {
		return NewAppError("UserAccessToken.IsValid", "model.user_access_token.is_valid.description.app_error", nil, "", http.StatusBadRequest)
	}

	return nil
}

func (t *UserAccessToken) PreSave() {
	t.Id = uuid.NewString()
	t.IsActive = true
}

func (t *UserAccessToken) ToJson() string {
	b, _ := json.JSON.Marshal(t)
	return string(b)
}

func UserAccessTokenFromJson(data io.Reader) *UserAccessToken {
	var t *UserAccessToken
	json.JSON.NewDecoder(data).Decode(&t)
	return t
}

func UserAccessTokenListToJson(t []*UserAccessToken) string {
	b, _ := json.JSON.Marshal(t)
	return string(b)
}

func UserAccessTokenListFromJson(data io.Reader) []*UserAccessToken {
	var t []*UserAccessToken
	json.JSON.NewDecoder(data).Decode(&t)
	return t
}
