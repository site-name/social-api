package model

import (
	"strings"

	"github.com/Masterminds/squirrel"
)

// max lengths for some fields
const (
	DIGITAL_CONTENT_CONTENT_TYPE_MAX_LENGTH = 128
)

const (
	FILE = "file"
)

// system supported content type
var ContentTypeString = map[string]string{
	FILE: "Digital product",
}

type DigitalContent struct {
	Id                   string `json:"id"`
	UseDefaultSettings   *bool  `json:"use_defaults_settings"` // default true
	AutomaticFulfillment *bool  `json:"automatic_fulfillment"` // default false
	ContentType          string `json:"content_type"`
	ProductVariantID     string `json:"product_variant_id"`
	ContentFile          string `json:"content_file"`
	MaxDownloads         *int   `json:"max_downloads"`
	UrlValidDays         *int   `json:"url_valid_days"`
	ModelMetadata
}

// DigitalContentFilterOption is used for building sql queries
type DigitalContentFilterOption struct {
	Id               squirrel.Sqlizer
	ProductVariantID squirrel.Sqlizer
}

func (d *DigitalContent) IsValid() *AppError {
	outer := CreateAppErrorForModel(
		"model.digital_content.is_valid.%s.app_error",
		"digital_content_id=",
		"DigitalContent.IsValid",
	)
	if !IsValidId(d.Id) {
		return outer("id", nil)
	}
	if len(d.ContentType) > DIGITAL_CONTENT_CONTENT_TYPE_MAX_LENGTH {
		return outer("content_type", &d.Id)
	}
	if ContentTypeString[strings.ToLower(d.ContentType)] == "" {
		return outer("content_type", &d.Id)
	}

	return nil
}

func (d *DigitalContent) ToJSON() string {
	return ModelToJson(d)
}

func (d *DigitalContent) DeepCopy() *DigitalContent {
	res := *d
	if d.UseDefaultSettings != nil {
		res.UseDefaultSettings = NewPrimitive(*d.UseDefaultSettings)
	}
	if d.AutomaticFulfillment != nil {
		res.AutomaticFulfillment = NewPrimitive(*d.AutomaticFulfillment)
	}
	if d.MaxDownloads != nil {
		res.MaxDownloads = NewPrimitive(*d.MaxDownloads)
	}
	if d.UrlValidDays != nil {
		res.UrlValidDays = NewPrimitive(*d.UrlValidDays)
	}
	return &res
}

func (d *DigitalContent) PreSave() {
	if d.Id == "" {
		d.Id = NewId()
	}
	if d.UseDefaultSettings == nil {
		d.UseDefaultSettings = NewPrimitive(true)
	}
	if d.AutomaticFulfillment == nil {
		d.AutomaticFulfillment = NewPrimitive(false)
	}
	if d.ContentType == "" {
		d.ContentType = FILE
	}
}

// max lengths for some fields of DigitalContentUrl
const (
	DIGITAL_CONTENT_URL_TOKEN_MAX_LENGTH = 36
)

type DigitalContentUrl struct {
	Id          string  `json:"id"`
	Token       string  `json:"token"` // uuid field, not editable, unique
	ContentID   string  `json:"content_id"`
	CreateAt    int64   `json:"create_at"`    // DEFAULT UTC now
	DownloadNum int     `json:"download_num"` //
	LineID      *string `json:"line_id"`      // 1-1 order line, unique
}

type DigitalContentUrlFilterOptions struct {
	Id        squirrel.Sqlizer
	Token     squirrel.Sqlizer
	ContentID squirrel.Sqlizer
	LineID    squirrel.Sqlizer
}

func (d *DigitalContentUrl) IsValid() *AppError {
	outer := CreateAppErrorForModel(
		"digital_content_url.is_valid.%s.app_error",
		"digital_content_url_id=",
		"DigitalContentUrl.IsValid",
	)
	if !IsValidId(d.Id) {
		return outer("id", nil)
	}
	if !IsValidId(d.ContentID) {
		return outer("content_id", &d.Id)
	}
	if d.LineID != nil && !IsValidId(*d.LineID) {
		return outer("line_id", &d.Id)
	}
	if len(d.Token) > DIGITAL_CONTENT_URL_TOKEN_MAX_LENGTH {
		return outer("token", &d.Id)
	}

	return nil
}

func (d *DigitalContentUrl) ToJSON() string {
	return ModelToJson(d)
}

func (d *DigitalContentUrl) PreSave() {
	if d.Id == "" {
		d.Id = NewId()
	}
	if d.CreateAt == 0 {
		d.CreateAt = GetMillis()
	}
	if d.Token == "" {
		d.NewToken(true)
	}
}

func (d *DigitalContentUrl) NewToken(force bool) {
	if (d.Token != "" && force) || d.Token == "" {
		d.Token = strings.ReplaceAll(NewId(), "-", "")
	}
}
