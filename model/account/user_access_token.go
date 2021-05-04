package account

import (
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/sitename/sitename/model"
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

func (t *UserAccessToken) IsValid() *model.AppError {
	if !model.IsValidId(t.Id) {
		return t.createAppError("id")
	}
	if len(t.Token) != USER_ACCESS_TOKEN_MAX_LENGTH {
		return t.createAppError("token")
	}
	if !model.IsValidId(t.UserId) {
		return t.createAppError("user_id")
	}
	if len(t.Description) > USER_ACCESS_TOKEN_DESCRIPTION_MAX_LENGTH {
		return t.createAppError("description")
	}

	return nil
}

func (u *UserAccessToken) createAppError(fieldName string) *model.AppError {
	id := fmt.Sprintf("model.user_access_token.is_valid.%s.app_error", fieldName)
	var details string
	if !strings.EqualFold(fieldName, "id") {
		details = "user_access_token_id=" + u.Id
	}

	return model.NewAppError("UserAccessToken.IsValid", id, nil, details, http.StatusBadRequest)
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
	var t UserAccessToken
	err := json.JSON.NewDecoder(data).Decode(&t)
	if err != nil {
		return nil
	}
	return &t
}

func UserAccessTokenListToJson(t []*UserAccessToken) string {
	b, _ := json.JSON.Marshal(t)
	return string(b)
}

func UserAccessTokenListFromJson(data io.Reader) []*UserAccessToken {
	var t []*UserAccessToken
	err := json.JSON.NewDecoder(data).Decode(&t)
	if err != nil {
		return nil
	}
	return t
}
