package model

import (
	"fmt"
	"io"
	"net/http"
	"unicode/utf8"

	"github.com/sitename/sitename/modules/json"
)

const (
	TERMS_OF_SERVICE_CACHE_SIZE = 1
	POST_MESSAGE_MAX_BYTES_V2   = 65535                         // Maximum size of a TEXT column in MySQL
	POST_MESSAGE_MAX_RUNES_V2   = POST_MESSAGE_MAX_BYTES_V2 / 4 // Assume a worst-case representation
)

type TermsOfService struct {
	Id       string `json:"id"`
	CreateAt int64  `json:"create_at"`
	UserId   string `json:"user_id"`
	Text     string `json:"text"`
}

func (t *TermsOfService) IsValid() *AppError {
	if !IsValidId(t.Id) {
		return InvalidTermsOfServiceError("id", "")
	}
	if t.CreateAt == 0 {
		return InvalidTermsOfServiceError("create_at", t.Id)
	}
	if !IsValidId(t.UserId) {
		return InvalidTermsOfServiceError("user_id", t.Id)
	}
	if utf8.RuneCountInString(t.Text) > POST_MESSAGE_MAX_RUNES_V2 {
		return InvalidTermsOfServiceError("text", t.Id)
	}

	return nil
}

func (t *TermsOfService) ToJson() string {
	b, _ := json.JSON.Marshal(t)
	return string(b)
}

func TermsOfServiceFromJson(data io.Reader) *TermsOfService {
	var termsOfService TermsOfService
	err := json.JSON.NewDecoder(data).Decode(&termsOfService)
	if err != nil {
		return nil
	}
	return &termsOfService
}

func InvalidTermsOfServiceError(fieldName string, termsOfServiceId string) *AppError {
	id := fmt.Sprintf("model.terms_of_service.is_valid.%s.app_error", fieldName)
	details := ""
	if termsOfServiceId != "" {
		details = "terms_of_service_id=" + termsOfServiceId
	}
	return NewAppError("TermsOfService.IsValid", id, map[string]interface{}{"MaxLength": POST_MESSAGE_MAX_RUNES_V2}, details, http.StatusBadRequest)
}

func (t *TermsOfService) PreSave() {
	if t.Id == "" {
		t.Id = NewId()
	}

	t.CreateAt = GetMillis()
}
