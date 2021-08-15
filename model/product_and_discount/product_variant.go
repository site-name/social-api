package product_and_discount

import (
	"io"
	"strings"
	"unicode/utf8"

	goprices "github.com/site-name/go-prices"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model/channel"
	"golang.org/x/text/language"

	"github.com/sitename/sitename/modules/measurement"
)

const (
	PRODUCT_VARIANT_NAME_MAX_LENGTH = 255
	PRODUCT_VARIANT_SKU_MAX_LENGTH  = 255
)

type ProductVariant struct {
	Id             string                 `json:"id"`
	Name           string                 `json:"name"`
	ProductID      string                 `json:"product_id"`
	Sku            string                 `json:"sku"`
	Weight         *float32               `json:"weight"`
	WeightUnit     measurement.WeightUnit `json:"weight_unit"`
	TrackInventory *bool                  `json:"track_inventory"` // default *true
	model.Sortable
	model.ModelMetadata

	DigitalContent *DigitalContent `json:"-" db:"-"` // for storing value returned by prefetching
}

// ProductVariantFilterOption is used to build sql queries
type ProductVariantFilterOption struct {
	Id   *model.StringFilter
	Name *model.StringFilter
}

func (p *ProductVariant) IsValid() *model.AppError {
	outer := model.CreateAppErrorForModel(
		"model.product_variant.is_valid.%s.app_error",
		"product_variant_id=",
		"ProductVariant.IsValid",
	)
	if !model.IsValidId(p.Id) {
		return outer("id", nil)
	}
	if !model.IsValidId(p.ProductID) {
		return outer("product_id", &p.Id)
	}
	if len(p.Sku) > PRODUCT_VARIANT_SKU_MAX_LENGTH {
		return outer("sku", &p.Id)
	}
	if utf8.RuneCountInString(p.Name) > PRODUCT_VARIANT_NAME_MAX_LENGTH {
		return outer("name", &p.Id)
	}
	if p.Weight != nil && *p.Weight <= 0 {
		return outer("weight", &p.Id)
	}
	if p.WeightUnit != "" {
		if _, ok := measurement.WEIGHT_UNIT_CONVERSION[p.WeightUnit]; !ok {
			return outer("weight_unit", &p.Id)
		}
	}

	return nil
}

func (p *ProductVariant) String() string {
	if p.Name != "" {
		return p.Name
	}
	return p.Sku
}

func (p *ProductVariant) ToJson() string {
	return model.ModelToJson(p)
}

func (p *ProductVariant) PreSave() {
	if p.Id == "" {
		p.Id = model.NewId()
	}
	p.Name = model.SanitizeUnicode(p.Name)
	if p.TrackInventory == nil {
		p.TrackInventory = model.NewBool(true)
	}
	if p.Weight != nil && p.WeightUnit == "" {
		p.WeightUnit = measurement.STANDARD_WEIGHT_UNIT
	}
	p.ModelMetadata.PreSave()
}

func (p *ProductVariant) PreUpdate() {
	p.Name = model.SanitizeUnicode(p.Name)
	if p.Weight != nil && p.WeightUnit == "" {
		p.WeightUnit = measurement.STANDARD_WEIGHT_UNIT
	}
	p.ModelMetadata.PreUpdate()
}

func ProductVariantFromJson(data io.Reader) *ProductVariant {
	var prd ProductVariant
	model.ModelFromJson(&prd, data)
	return &prd
}

// TODO: fixme
func (p *ProductVariant) GetPrice(product *Product, collections []*Collection, channel *channel.Channel, channelListing *ProductChannelListing, discounts []*DiscountInfo) *goprices.Money {
	panic("not impl")
}

// TODO: fixme
func (p *ProductVariant) DisplayProduct() {
	panic("not implemented")
}

// --------------------
type ProductVariantTranslation struct {
	Id               string `json:"id"`
	LanguageCode     string `json:"language_code"`
	ProductVariantID string `json:"product_variant_id"`
	Name             string `json:"name"`
}

func (p *ProductVariantTranslation) String() string {
	return p.Name
}

func (p *ProductVariantTranslation) IsValid() *model.AppError {
	outer := model.CreateAppErrorForModel(
		"model.product_variant_translation.is_valid.%s.app_error",
		"product_variant_translation_id=",
		"ProductVariantTranslation.IsValid",
	)
	if !model.IsValidId(p.Id) {
		return outer("id", nil)
	}
	if !model.IsValidId(p.ProductVariantID) {
		return outer("product_variant_id", &p.Id)
	}
	tag, err := language.Parse(p.LanguageCode)
	if err != nil || !strings.EqualFold(tag.String(), p.LanguageCode) || model.Languages[strings.ToLower(p.LanguageCode)] == "" {
		return outer("language_code", &p.Id)
	}

	return nil
}

func (p *ProductVariantTranslation) ToJson() string {
	return model.ModelToJson(p)
}

func ProductVariantTranslationFromJson(data io.Reader) *ProductVariantTranslation {
	var p ProductVariantTranslation
	model.ModelFromJson(&p, data)
	return &p
}
