package model

import (
	"net/http"

	"github.com/Masterminds/squirrel"
	"gorm.io/gorm"
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

type ProductMedia struct {
	Id          string          `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid();column:Id"`
	CreateAt    int64           `json:"create_at" gorm:"type:bigint;column:CreateAt;autoCreateTime:milli"`
	ProductID   string          `json:"product_id" gorm:"type:uuid;column:ProductID"`
	Ppoi        string          `json:"ppoi" gorm:"type:varchar(20);column:Ppoi"` // holds resolution for images, not editable
	Image       string          `json:"image" gorm:"type:varchar(500);column:Image"`
	Alt         string          `json:"alt" gorm:"type:varchar(128);column:Alt"`
	Type        string          `json:"type" gorm:"type:varchar(32);column:Type"` // default to "IMAGE"
	ExternalUrl *string         `json:"external_url" gorm:"type:varchar(256);column:ExternalUrl"`
	OembedData  StringInterface `json:"oembed_data" gorm:"type:jsonb;column:OembedData"`
	Sortable

	ProductVariants ProductVariants `json:"-" gorm:"many2many:VariantMedias"`
}

func (c *ProductMedia) BeforeCreate(_ *gorm.DB) error { c.commonPre(); return c.IsValid() }
func (c *ProductMedia) BeforeUpdate(_ *gorm.DB) error { c.commonPre(); return c.IsValid() }
func (c *ProductMedia) TableName() string             { return ProductMediaTableName }

type ProductMedias []*ProductMedia

// ProductMediaFilterOption is used for building squirrel sql queries
type ProductMediaFilterOption struct {
	Conditions squirrel.Sqlizer

	VariantID squirrel.Sqlizer // INNER JOIN VariantMedias ON VariantMedias.MediaID = ProductMedias.Id Where VariantMedias.VariantID ...
}

func (p *ProductMedia) IsValid() *AppError {
	if !IsValidId(p.ProductID) {
		return NewAppError("ProductMedia.IsValid", "model.product_media.is_valid.product_id.app_error", nil, "please provide valid product id", http.StatusBadRequest)
	}
	if ProductMediaTypeChoices[p.Type] == "" {
		return NewAppError("ProductMedia.IsValid", "model.product_media.is_valid.type.app_error", nil, "please provide valid type", http.StatusBadRequest)
	}

	return nil
}

func (p *ProductMedia) commonPre() {
	p.Alt = SanitizeUnicode(p.Alt)
	if ProductMediaTypeChoices[p.Type] == "" {
		p.Type = IMAGE
	}
}
