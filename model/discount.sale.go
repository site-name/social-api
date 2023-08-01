package model

import (
	"net/http"
	"strings"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/samber/lo"
	"golang.org/x/text/language"
	"gorm.io/gorm"
)

type Sale struct {
	Id        string       `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid();column:Id"`
	Name      string       `json:"name" gorm:"type:varchar(255);column:Name"`
	Type      DiscountType `json:"type" gorm:"type:varchar(10);column:Type"` // DEFAULT `fixed`
	StartDate time.Time    `json:"start_date" gorm:"column:StartDate"`
	EndDate   *time.Time   `json:"end_date" gorm:"column:EndDate"`
	CreateAt  int64        `json:"create_at" gorm:"autoCreateTime:milli;column:CreateAt"`
	UpdateAt  int64        `json:"update_at" gorm:"autoUpdateTime:milli;autoCreateTime:milli;column:UpdateAt"`
	ModelMetadata

	Categories      Categories      `json:"-" gorm:"many2many:SaleCategories"`
	Products        Products        `json:"-" gorm:"many2many:SaleProducts"`
	ProductVariants ProductVariants `json:"-" gorm:"many2many:SaleProductVariants"`
	Collections     Collections     `json:"-" gorm:"many2many:SaleCollections"`
}

func (c *Sale) BeforeCreate(_ *gorm.DB) error { c.commonPre(); return c.IsValid() }
func (c *Sale) BeforeUpdate(_ *gorm.DB) error { c.commonPre(); return c.IsValid() }
func (c *Sale) TableName() string             { return SaleTableName }

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

// SaleFilterOption can be used to
type SaleFilterOption struct {
	Conditions                   squirrel.Sqlizer
	SaleChannelListing_ChannelID squirrel.Sqlizer // INNER JOIN SaleChannelListings ON ... WHERE SaleChannelListings.ChannelID ...
}

type Sales []*Sale

func (s Sales) IDs() []string {
	return lo.Map(s, func(sa *Sale, _ int) string { return sa.Id })
}

func (s *Sale) String() string {
	return s.Name
}

func (s *Sale) IsValid() *AppError {
	if !s.Type.IsValid() {
		return NewAppError("Sale.IsValid", "model.sale.is_valid.type.app_error", nil, "please provide valid type", http.StatusBadRequest)
	}
	if s.StartDate.IsZero() {
		return NewAppError("Sale.IsValid", "model.sale.is_valid.start_date.app_error", nil, "please provide valid start date", http.StatusBadRequest)
	}
	if s.EndDate != nil && s.EndDate.IsZero() {
		return NewAppError("Sale.IsValid", "model.sale.is_valid.end_date.app_error", nil, "please provide valid end date", http.StatusBadRequest)
	}
	if s.EndDate != nil && s.EndDate.Before(s.StartDate) {
		return NewAppError("Sale.IsValid", "model.sale.is_valid.dates.app_error", nil, "start date must be before end date", http.StatusBadRequest)
	}

	return nil
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

type SaleTranslation struct {
	Id           string `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid();column:Id"`
	LanguageCode string `json:"language_code" gorm:"type:varchar(5);column:LanguageCode"`
	Name         string `json:"name" gorm:"type:varchar(255);column:Name"`
	SaleID       string `json:"sale_id" gorm:"type:uuid;column:SaleID"`
}

func (s *SaleTranslation) commonPre() {
	s.Name = SanitizeUnicode(s.Name)
}

func (s *SaleTranslation) IsValid() *AppError {
	if !IsValidId(s.SaleID) {
		return NewAppError("SaleTranslation.IsValid", "model.sale_transalation.is_valid.sale_id.app_error", nil, "please provide valid sale id", http.StatusBadRequest)
	}
	if tag, err := language.Parse(s.LanguageCode); err != nil || !strings.EqualFold(tag.String(), s.LanguageCode) {
		return NewAppError("SaleTranslation.IsValid", "model.sale_transalation.is_valid.language_code.app_error", nil, "please provide valid language code", http.StatusBadRequest)
	}

	return nil
}

func (c *SaleTranslation) BeforeCreate(_ *gorm.DB) error { c.commonPre(); return c.IsValid() }
func (c *SaleTranslation) BeforeUpdate(_ *gorm.DB) error { c.commonPre(); return c.IsValid() }
func (c *SaleTranslation) TableName() string             { return SaleTranslationTableName }
