package model

import (
	"net/http"
	"strings"

	"github.com/Masterminds/squirrel"
	"gorm.io/gorm"
)

const (
	FILE = "file"
)

// system supported content type
var ContentTypeString = map[string]string{
	FILE: "Digital product",
}

type DigitalContent struct {
	Id                   string `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid();column:Id"`
	UseDefaultSettings   *bool  `json:"use_defaults_settings" gorm:"column:UseDefaultSettings"`   // default true
	AutomaticFulfillment *bool  `json:"automatic_fulfillment" gorm:"column:AutomaticFulfillment"` // default false
	ContentType          string `json:"content_type" gorm:"type:varchar(128);column:ContentType"`
	ProductVariantID     string `json:"product_variant_id" gorm:"type:uuid;column:ProductVariantID"`
	ContentFile          string `json:"content_file" gorm:"type:varchar(300);column:ContentFile"`
	MaxDownloads         *int   `json:"max_downloads" gorm:"column:MaxDownloads"`
	UrlValidDays         *int   `json:"url_valid_days" gorm:"column:UrlValidDays"`
	ModelMetadata
}

func (c *DigitalContent) BeforeCreate(_ *gorm.DB) error { c.commonPre(); return c.IsValid() }
func (c *DigitalContent) BeforeUpdate(_ *gorm.DB) error { c.commonPre(); return c.IsValid() }
func (c *DigitalContent) TableName() string             { return DigitalContentTableName }

// DigitalContentFilterOption is used for building sql queries
type DigitalContentFilterOption struct {
	Conditions squirrel.Sqlizer

	PaginationValues GraphqlPaginationValues
	CountTotal       bool
}

func (d *DigitalContent) IsValid() *AppError {
	if ContentTypeString[strings.ToLower(d.ContentType)] == "" {
		return NewAppError("DigitalContent.IsValid", "model.digital_content.is_valid.content_type.app_error", nil, "please provide valid content type", http.StatusBadRequest)
	}

	return nil
}

func (d *DigitalContent) DeepCopy() *DigitalContent {
	res := *d
	if d.UseDefaultSettings != nil {
		res.UseDefaultSettings = GetPointerOfValue(*d.UseDefaultSettings)
	}
	if d.AutomaticFulfillment != nil {
		res.AutomaticFulfillment = GetPointerOfValue(*d.AutomaticFulfillment)
	}
	if d.MaxDownloads != nil {
		res.MaxDownloads = GetPointerOfValue(*d.MaxDownloads)
	}
	if d.UrlValidDays != nil {
		res.UrlValidDays = GetPointerOfValue(*d.UrlValidDays)
	}
	return &res
}

func (d *DigitalContent) commonPre() {
	if d.UseDefaultSettings == nil {
		d.UseDefaultSettings = GetPointerOfValue(true)
	}
	if d.AutomaticFulfillment == nil {
		d.AutomaticFulfillment = GetPointerOfValue(false)
	}
	if d.ContentType != FILE {
		d.ContentType = FILE
	}
}

type DigitalContentUrl struct {
	Id          string  `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid();column:Id"`
	Token       string  `json:"token" gorm:"type:uuid;uniqueIndex:token_key;column:Token"` // uuid field, not editable, unique
	ContentID   string  `json:"content_id" gorm:"type:uuid;column:ContentID"`
	CreateAt    int64   `json:"create_at" gorm:"type:bigint;column:CreateAt"` // DEFAULT UTC now
	DownloadNum int     `json:"download_num" gorm:"column:DownloadNum"`       //
	LineID      *string `json:"line_id" gorm:"type:uuid;column:LineID"`       // 1-1 order line, unique
}

func (c *DigitalContentUrl) BeforeCreate(_ *gorm.DB) error { c.commonPre(); return c.IsValid() }
func (c *DigitalContentUrl) BeforeUpdate(_ *gorm.DB) error {
	c.commonPre()
	c.CreateAt = 0 // prevent update
	c.Token = ""
	return c.IsValid()
}
func (c *DigitalContentUrl) TableName() string { return DigitalContentURLTableName }

type DigitalContentUrlFilterOptions struct {
	Conditions squirrel.Sqlizer
}

func (d *DigitalContentUrl) IsValid() *AppError {
	if !IsValidId(d.ContentID) {
		return NewAppError("DigitalContentUrl.IsValid", "model.digital_content_url.is_valid.content_id.app_error", nil, "please provide valid content id", http.StatusBadRequest)
	}
	if d.LineID != nil && !IsValidId(*d.LineID) {
		return NewAppError("DigitalContentUrl.IsValid", "model.digital_content_url.is_valid.line_id.app_error", nil, "please provide valid line id", http.StatusBadRequest)
	}
	return nil
}

func (d *DigitalContentUrl) commonPre() {
	if d.Token == "" {
		d.NewToken(true)
	}
}

func (d *DigitalContentUrl) NewToken(force bool) {
	if (d.Token != "" && force) || d.Token == "" {
		d.Token = strings.ReplaceAll(NewId(), "-", "")
	}
}
