package model

import (
	"io"
	"unicode/utf8"
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
	outer := CreateAppErrorForModel(
		"model.terms_of_service.is_valid.%s.app_error",
		"terms_of_service_id=",
		"TermsOfService.IsValid",
	)
	if !IsValidId(t.Id) {
		return outer("id", nil)
	}
	if t.CreateAt == 0 {
		return outer("create_at", &t.Id)
	}
	if !IsValidId(t.UserId) {
		return outer("user_id", &t.Id)
	}
	if utf8.RuneCountInString(t.Text) > POST_MESSAGE_MAX_RUNES_V2 {
		return outer("text", &t.Id)
	}

	return nil
}

func (t *TermsOfService) ToJSON() string {
	return ModelToJson(t)
}

func TermsOfServiceFromJson(data io.Reader) *TermsOfService {
	var termsOfService TermsOfService
	ModelFromJson(&termsOfService, data)
	return &termsOfService
}

func (t *TermsOfService) PreSave() {
	if t.Id == "" {
		t.Id = NewId()
	}

	t.CreateAt = GetMillis()
}
