package product_and_discount

import (
	"io"
	"strings"

	"github.com/sitename/sitename/model"
)

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
	// PRODUCT_MEDIA_TYPE_MAX_LENGTH         = 32
	PRODUCT_MEDIA_EXTERNAL_URL_MAX_LENGTH = 256
	PRODUCT_MEDIA_ALT_MAX_LENGTH          = 128
	PRODUCT_MEDIA_PPOI_MAX_LENGTH         = 20
)

// TODO: not done yet
type ProductMedia struct {
	Id          string                `json:"id"`
	ProductID   string                `json:"product_id"`
	Ppoi        string                `json:"ppoi"` // NOTE: need investigation
	Image       *model.FileInfo       `db:"-"`
	Alt         string                `json:"alt"`
	Type        string                `json:"type"`
	ExternalUrl *string               `json:"external_url"`
	Product     *Product              `json:"product" db:"-"`
	OembedData  model.StringInterface `json:"oembed_data"`
	*model.Sortable
}

// TODO: not done yet
func (p *ProductMedia) IsValid() *model.AppError {
	outer := model.CreateAppErrorForModel(
		"model.product_media.is_valid.%s.app_error",
		"product_media_id=",
		"ProductMedia.IsValid",
	)
	if !model.IsValidId(p.Id) {
		return outer("id", nil)
	}
	if !model.IsValidId(p.ProductID) {
		return outer("product_id", &p.Id)
	}
	if len(p.Ppoi) > PRODUCT_MEDIA_PPOI_MAX_LENGTH {
		return outer("ppoi", &p.Id)
	}
	if len(p.Alt) > PRODUCT_MEDIA_ALT_MAX_LENGTH {
		return outer("alt", &p.Id)
	}
	if ProductMediaTypeChoices[strings.ToUpper(p.Type)] == "" {
		return outer("type", &p.Id)
	}
	if p.ExternalUrl != nil && len(*p.ExternalUrl) > PRODUCT_MEDIA_EXTERNAL_URL_MAX_LENGTH {
		return outer("external_url", &p.Id)
	}

	return nil
}

func (p *ProductMedia) ToJson() string {
	return model.ModelToJson(p)
}

func ProductMediaFromJson(data io.Reader) *ProductMedia {
	var prd ProductMedia
	model.ModelFromJson(&prd, data)
	return &prd
}

func (p *ProductMedia) PreSave() {
	if p.Id == "" {
		p.Id = model.NewId()
	}
	if p.Ppoi == "" {
		p.Ppoi = "0.5x0.5"
	}
}

func (p *ProductMedia) GetOrderingQueryset() []*ProductMedia {
	return p.Product.Medias
}
