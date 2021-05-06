package product_and_discount

import (
	"io"
	"strings"
	"unicode/utf8"

	"github.com/gosimple/slug"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model/seo"
	"github.com/sitename/sitename/modules/measurement"
	"golang.org/x/text/language"
)

const (
	PRODUCT_NAME_MAX_LENGTH = 250
	PRODUCT_SLUG_MAX_LENGTH = 255
)

// Product contains all fields a product contains
type Product struct {
	Id                   string   `json:"id"`
	ProductTypeID        string   `json:"product_type_id"`
	Name                 string   `json:"name"`
	Slug                 string   `json:"slug"`
	Description          *string  `json:"description"`
	DescriptionPlainText string   `json:"description_plaintext"`
	CategoryID           *string  `json:"category_id"`
	CreateAt             int64    `json:"create_at"`
	UpdateAt             int64    `json:"update_at"`
	ChargeTaxes          *bool    `json:"charge_taxes"`
	Weight               *float32 `json:"weight"`
	WeightUnit           string   `json:"weight_unit"`
	DefaultVariantID     *string  `json:"default_variant_id"`
	Rating               *float32 `json:"rating"`
	*model.ModelMetadata
	*seo.Seo
}

func (p *Product) PlainTextDescription() string {
	panic("not implemented")
}

func SortByAttributeFields() []string {
	return []string{"concatenated_values_order", "concatenated_values", "name"}
}

func (p *Product) ToJson() string {
	return model.ModelToJson(p)
}

func ProductFromJson(data io.Reader) *Product {
	var p Product
	model.ModelFromJson(&p, data)
	return &p
}

func (p *Product) IsValid() *model.AppError {
	outer := model.CreateAppErrorForModel(
		"model.product.is_valid.%s.app_error",
		"product_id=",
		"Product.IsValid",
	)
	if !model.IsValidId(p.Id) {
		return outer("id", nil)
	}
	if !model.IsValidId(p.ProductTypeID) {
		return outer("product_type_id", &p.Id)
	}
	if utf8.RuneCountInString(p.Name) > PRODUCT_NAME_MAX_LENGTH {
		return outer("name", &p.Id)
	}
	if p.CategoryID != nil && *p.CategoryID == "" {
		return outer("category_id", &p.Id)
	}
	if p.CreateAt == 0 {
		return outer("create_at", &p.Id)
	}
	if p.UpdateAt == 0 {
		return outer("update_at", &p.Id)
	}
	if utf8.RuneCountInString(p.Slug) > PRODUCT_SLUG_MAX_LENGTH {
		return outer("slug", &p.Id)
	}
	if p.Weight != nil && *p.Weight == 0 {
		return outer("weight", &p.Id)
	}
	if p.Weight != nil {
		if _, ok := measurement.WEIGHT_UNIT_STRINGS[strings.ToLower(p.WeightUnit)]; !ok {
			return outer("weight_unit", &p.Id)
		}
	}

	return nil
}

func (p *Product) PreSave() {
	if p.Id == "" {
		p.Id = model.NewId()
	}
	p.CreateAt = model.GetMillis()
	p.UpdateAt = p.CreateAt
	p.Name = model.SanitizeUnicode(p.Name)
	p.Slug = slug.Make(p.Name)
	if p.Weight != nil && p.WeightUnit == "" {
		p.WeightUnit = measurement.STANDARD_WEIGHT_UNIT
	}
}

func (p *Product) PreUpdate() {
	p.UpdateAt = model.GetMillis()
	p.Name = model.SanitizeUnicode(p.Name)
	p.Slug = slug.Make(p.Name)
	if p.Weight != nil && p.WeightUnit == "" {
		p.WeightUnit = measurement.STANDARD_WEIGHT_UNIT
	}
}

type ProductTranslation struct {
	Id           string                 `json:"id"`
	LanguageCode string                 `json:"language_code"`
	ProductID    string                 `json:"product_id"`
	Name         string                 `json:"name"`
	Description  *model.StringInterface `json:"description"`
	*seo.SeoTranslation
}

func (p *ProductTranslation) String() string {
	return p.Name
}

func (p *ProductTranslation) IsValid() *model.AppError {
	outer := model.CreateAppErrorForModel(
		"model.product_translation.is_valid.%s.app_error",
		"product_translation_id=",
		"ProductTranslation.IsValid",
	)
	if !model.IsValidId(p.Id) {
		return outer("id", nil)
	}
	if !model.IsValidId(p.ProductID) {
		return outer("product_id", &p.Id)
	}
	if tag, err := language.Parse(p.LanguageCode); err != nil || tag.String() != p.LanguageCode {
		return outer("language_code", &p.Id)
	}
	if utf8.RuneCountInString(p.Name) > PRODUCT_NAME_MAX_LENGTH {
		return outer("name", &p.Id)
	}

	return nil
}

func (p *ProductTranslation) ToJson() string {
	return model.ModelToJson(p)
}

func ProductTranslationFromJson(data io.Reader) *ProductTranslation {
	var p ProductTranslation
	model.ModelFromJson(&p, data)
	return &p
}
