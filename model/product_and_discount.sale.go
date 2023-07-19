package model

import (
	"strings"
	"time"
	"unicode/utf8"

	"github.com/Masterminds/squirrel"
	"github.com/samber/lo"
	"golang.org/x/text/language"
	"gorm.io/gorm"
)

// max lengths for some fields
const (
	SALE_NAME_MAX_LENGTH = 255
	SALE_TYPE_MAX_LENGTH = 10
)

type Sale struct {
	Id        string       `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	Name      string       `json:"name" gorm:"type:varchar(255)"`
	Type      DiscountType `json:"type" gorm:"type:varchar(10)"` // DEFAULT `fixed`
	StartDate time.Time    `json:"start_date"`
	EndDate   *time.Time   `json:"end_date"`
	CreateAt  int64        `json:"create_at" gorm:"autoCreateTime:milli"`
	UpdateAt  int64        `json:"update_at" gorm:"autoUpdateTime:milli"`
	ModelMetadata

	Categories      Categories      `json:"-" gorm:"many2many:SaleCategories"`
	Products        Products        `json:"-" gorm:"many2many:SaleProducts"`
	ProductVariants ProductVariants `json:"-" gorm:"many2many:SaleProductVariants"`
	Collections     Collections     `json:"-" gorm:"many2many:SaleCollections"`
}

type SaleCollection struct {
	SaleID       string
	CollectionID string
}

type SaleProduct struct {
	SaleID    string
	ProductID string
}

type SaleCategory struct {
	SaleID     string
	CategoryID string
}

type SaleProductVariant struct {
	SaleID           string
	ProductVariantID string
}

// BeforeCreate is gorm hook
func (s *Sale) BeforeCreate(_ *gorm.DB) error {
	s.commonPre()
	return nil
}

func (s *Sale) BeforeUpdate() error {
	s.commonPre()
	return nil
}

// SaleFilterOption can be used to
type SaleFilterOption struct {
	StartDate squirrel.Sqlizer
	EndDate   squirrel.Sqlizer
}

type Sales []*Sale

func (s Sales) IDs() []string {
	return lo.Map(s, func(sa *Sale, _ int) string { return sa.Id })
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
	if len(s.Type) > SALE_TYPE_MAX_LENGTH || !s.Type.IsValid() {
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
	s.commonPre()
}

func (s *Sale) commonPre() {
	if s.StartDate.IsZero() {
		s.StartDate = time.Now()
	}
	if !s.Type.IsValid() {
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
		"model.sale_translation.is_valid.%s.app_error",
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
