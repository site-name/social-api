package model

import (
	"strings"
	"time"
	"unicode/utf8"

	"github.com/Masterminds/squirrel"
	"golang.org/x/text/language"
)

// max lengths for some fields
const (
	SALE_NAME_MAX_LENGTH = 255
	SALE_TYPE_MAX_LENGTH = 10
)

type Sale struct {
	Id        string     `json:"id"`
	Name      string     `json:"name"`
	Type      string     `json:"type"` // DEFAULT `fixed`
	StartDate time.Time  `json:"start_date"`
	EndDate   *time.Time `json:"end_date"`
	CreateAt  int64      `json:"create_at"`
	UpdateAt  int64      `json:"update_at"`
	ModelMetadata
}

// SaleFilterOption can be used to
type SaleFilterOption struct {
	StartDate squirrel.Sqlizer
	EndDate   squirrel.Sqlizer
}

type Sales []*Sale

func (s Sales) IDs() []string {
	res := []string{}
	for _, item := range s {
		if item != nil {
			res = append(res, item.Id)
		}
	}

	return res
}

func (s *Sale) String() string {
	return s.Name
}

func (s *Sale) IsValid() *AppError {
	outer := CreateAppErrorForModel(
		"model.sale.is_valid.%s.app_error",
		"sale_id=",
		"Sale.IsValid",
	)
	if !IsValidId(s.Id) {
		return outer("id", nil)
	}
	if utf8.RuneCountInString(s.Name) > SALE_NAME_MAX_LENGTH {
		return outer("name", &s.Id)
	}
	if len(s.Type) > SALE_TYPE_MAX_LENGTH || !SALE_TYPES.Contains(s.Type) {
		return outer("type", &s.Id)
	}
	if s.StartDate.IsZero() {
		return outer("start_date", &s.Id)
	}
	if s.EndDate != nil && s.EndDate.IsZero() {
		return outer("end_date", &s.Id)
	}
	if s.CreateAt == 0 {
		return outer("create_at", &s.Id)
	}
	if s.UpdateAt == 0 {
		return outer("update_at", &s.Id)
	}

	return nil
}

func (s *Sale) PreSave() {
	if s.Id == "" {
		s.Id = NewId()
	}
	s.CreateAt = GetMillis()
	s.UpdateAt = s.CreateAt

	if s.StartDate.IsZero() {
		s.StartDate = time.Now()
	}
	s.commonPre()
}

func (s *Sale) commonPre() {
	if s.Type == "" || !SALE_TYPES.Contains(s.Type) {
		s.Type = FIXED
	}
	s.Name = SanitizeUnicode(s.Name)
}

func (s *Sale) PreUpdate() {
	s.UpdateAt = GetMillis()
	s.commonPre()
}

type SaleTranslation struct {
	Id           string `json:"id"`
	LanguageCode string `json:"language_code"`
	Name         string `json:"name"`
	SaleID       string `json:"sale_id"`
}

func (s *SaleTranslation) IsValid() *AppError {
	outer := CreateAppErrorForModel(
		"sale_translation.is_valid.%s.app_error",
		"sale_translation_id=",
		"SaleTranslation.IsValid",
	)
	if !IsValidId(s.Id) {
		return outer("id", nil)
	}
	if !IsValidId(s.SaleID) {
		return outer("sale_id", &s.Id)
	}
	if tag, err := language.Parse(s.LanguageCode); err != nil || !strings.EqualFold(tag.String(), s.LanguageCode) {
		return outer("language_code", &s.Id)
	}
	if utf8.RuneCountInString(s.Name) > SALE_NAME_MAX_LENGTH {
		return outer("name", &s.Id)
	}

	return nil
}
