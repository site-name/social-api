package model

import (
	"strings"

	"github.com/Masterminds/squirrel"
)

// valid values for product media's Type
const (
	VIDEO = "VIDEO"
	IMAGE = "IMAGE"
)

var ProductMediaTypeChoices = map[string]string{
	VIDEO: "A URL to an external video",
	IMAGE: "An uploaded image or an URL to an image",
}

// max lengths limits for some fields
const (
	PRODUCT_MEDIA_TYPE_MAX_LENGTH         = 32
	PRODUCT_MEDIA_EXTERNAL_URL_MAX_LENGTH = 256
	PRODUCT_MEDIA_ALT_MAX_LENGTH          = 128
	PRODUCT_MEDIA_PPOI_MAX_LENGTH         = 20
	PRODUCT_MEDIA_IMAGE_LINK_MAX_LENGTH   = 100
)

type ProductMedia struct {
	Id          string          `json:"id"`
	CreateAt    int64           `json:"create_at"`
	ProductID   string          `json:"product_id"`
	Ppoi        string          `json:"ppoi"` // holds resolution for images
	Image       string          `json:"image"`
	Alt         string          `json:"alt"`
	Type        string          `json:"type"`
	ExternalUrl *string         `json:"external_url"`
	OembedData  StringInterface `json:"oembed_data"`
	Sortable
}

// ProductMediaFilterOption is used for building squirrel sql queries
type ProductMediaFilterOption struct {
	Id        squirrel.Sqlizer
	ProductID squirrel.Sqlizer
	Type      squirrel.Sqlizer
}

func (p *ProductMedia) IsValid() *AppError {
	outer := CreateAppErrorForModel(
		"product_media.is_valid.%s.app_error",
		"product_media_id=",
		"ProductMedia.IsValid",
	)
	if !IsValidId(p.Id) {
		return outer("id", nil)
	}
	if p.CreateAt == 0 {
		return outer("create_at", &p.Id)
	}
	if !IsValidId(p.ProductID) {
		return outer("product_id", &p.Id)
	}
	if len(p.Ppoi) > PRODUCT_MEDIA_PPOI_MAX_LENGTH {
		return outer("ppoi", &p.Id)
	}
	if len(p.Image) > PRODUCT_MEDIA_IMAGE_LINK_MAX_LENGTH {
		return outer("image", &p.Id)
	}
	if len(p.Alt) > PRODUCT_MEDIA_ALT_MAX_LENGTH {
		return outer("alt", &p.Id)
	}
	if ProductMediaTypeChoices[strings.ToUpper(p.Type)] == "" || len(p.Type) > PRODUCT_MEDIA_TYPE_MAX_LENGTH {
		return outer("type", &p.Id)
	}
	if p.ExternalUrl != nil && len(*p.ExternalUrl) > PRODUCT_MEDIA_EXTERNAL_URL_MAX_LENGTH {
		return outer("external_url", &p.Id)
	}

	return nil
}

func (p *ProductMedia) PreSave() {
	if p.Id == "" {
		p.Id = NewId()
	}
	p.CreateAt = GetMillis()
	p.commonPre()
}

func (p *ProductMedia) commonPre() {
	p.Alt = SanitizeUnicode(p.Alt)
}

func (p *ProductMedia) PreUpdate() {
	p.commonPre()
}