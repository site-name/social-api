package product

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"unicode/utf8"

	"github.com/gosimple/slug"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model/seo"
	"github.com/sitename/sitename/modules/json"
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
	b, _ := json.JSON.Marshal(p)
	return string(b)
}

func ProductFromJson(data io.Reader) *Product {
	var p Product
	err := json.JSON.NewDecoder(data).Decode(&p)
	if err != nil {
		return nil
	}
	return &p
}

func (p *Product) createAppError(fieldName string) *model.AppError {
	id := fmt.Sprintf("model.product.is_valid.%s.app_error", fieldName)
	var details string
	if !strings.EqualFold(fieldName, "id") {
		details = "product_id=" + p.Id
	}

	return model.NewAppError("Product.IsValid", id, nil, details, http.StatusBadRequest)
}

func (p *Product) IsValid() *model.AppError {
	if !model.IsValidId(p.Id) {
		return p.createAppError("id")
	}
	if !model.IsValidId(p.ProductTypeID) {
		return p.createAppError("product_type_id")
	}
	if utf8.RuneCountInString(p.Name) > PRODUCT_NAME_MAX_LENGTH {
		return p.createAppError("name")
	}
	if p.CategoryID != nil && *p.CategoryID == "" {
		return p.createAppError("category_id")
	}
	if p.CreateAt == 0 {
		return p.createAppError("create_at")
	}
	if p.UpdateAt == 0 {
		return p.createAppError("update_at")
	}
	if utf8.RuneCountInString(p.Slug) > PRODUCT_SLUG_MAX_LENGTH {
		return p.createAppError("slug")
	}
	if p.Weight != nil && *p.Weight == 0 {
		return p.createAppError("weight")
	}
	if p.Weight != nil {
		if _, ok := WeightUnitString[p.WeightUnit]; !ok {
			return p.createAppError("weight_unit")
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
	if p.Weight != nil {
		if p.WeightUnit == "" {
			p.WeightUnit = KG
		}
	}
}

func (p *Product) PreUpdate() {
	p.UpdateAt = model.GetMillis()
	p.Name = model.SanitizeUnicode(p.Name)
	p.Slug = slug.Make(p.Name)
	if p.Weight != nil {
		if p.WeightUnit == "" {
			p.WeightUnit = KG
		}
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

func (p *ProductTranslation) createAppError(fieldName string) *model.AppError {
	id := fmt.Sprintf("model.product_translation.is_valid.%s.app_error", fieldName)
	var details string
	if !strings.EqualFold(fieldName, "id") {
		details = "product_translation_id=" + p.Id
	}

	return model.NewAppError("ProductTranslation.IsValid", id, nil, details, http.StatusBadRequest)
}

func (p *ProductTranslation) IsValid() *model.AppError {
	if !model.IsValidId(p.Id) {
		return p.createAppError("id")
	}
	if !model.IsValidId(p.ProductID) {
		return p.createAppError("product_id")
	}
	if tag, err := language.Parse(p.LanguageCode); err != nil || tag.String() != p.LanguageCode {
		return p.createAppError("language_code")
	}
	if utf8.RuneCountInString(p.Name) > PRODUCT_NAME_MAX_LENGTH {
		return p.createAppError("name")
	}

	return nil
}

func (p *ProductTranslation) ToJson() string {
	b, _ := json.JSON.Marshal(p)
	return string(b)
}

func ProductTranslationFromJson(data io.Reader) *ProductTranslation {
	var p ProductTranslation
	err := json.JSON.NewDecoder(data).Decode(&p)
	if err != nil {
		return nil
	}
	return &p
}
