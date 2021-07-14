package product_and_discount

import (
	"io"
	"unicode/utf8"

	"github.com/gosimple/slug"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model/seo"
	"github.com/sitename/sitename/modules/measurement"
)

const (
	PRODUCT_NAME_MAX_LENGTH = 250
	PRODUCT_SLUG_MAX_LENGTH = 255
)

// Product contains all fields a product contains
type Product struct {
	Id                   string                 `json:"id"`
	ProductTypeID        string                 `json:"product_type_id"`
	Name                 string                 `json:"name"`
	Slug                 string                 `json:"slug"`
	Description          *string                `json:"description"`
	DescriptionPlainText string                 `json:"description_plaintext"`
	CategoryID           *string                `json:"category_id"`
	CreateAt             int64                  `json:"create_at"`
	UpdateAt             int64                  `json:"update_at"`
	ChargeTaxes          *bool                  `json:"charge_taxes"`
	Weight               *float32               `json:"weight"`
	WeightUnit           measurement.WeightUnit `json:"weight_unit"`
	DefaultVariantID     *string                `json:"default_variant_id"`
	Rating               *float32               `json:"rating"`
	Medias               []*ProductMedia        `json:"medias" db:"-"`
	ProductType          *ProductType           `db:"-"`
	Variants             []*ProductVariant      `db:"-"`
	model.ModelMetadata
	seo.Seo
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
	if p.CategoryID != nil && !model.IsValidId(*p.CategoryID) {
		return outer("category_id", &p.Id)
	}
	if utf8.RuneCountInString(p.Name) > PRODUCT_NAME_MAX_LENGTH {
		return outer("name", &p.Id)
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
	if p.Weight != nil {
		if _, ok := measurement.WEIGHT_UNIT_STRINGS[p.WeightUnit]; !ok {
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
	if p.ChargeTaxes == nil {
		p.ChargeTaxes = model.NewBool(true)
	}
}

func (p *Product) PreUpdate() {
	p.UpdateAt = model.GetMillis()
	p.Name = model.SanitizeUnicode(p.Name)
	if p.Weight != nil && p.WeightUnit == "" {
		p.WeightUnit = measurement.STANDARD_WEIGHT_UNIT
	}
	if p.ChargeTaxes == nil {
		p.ChargeTaxes = model.NewBool(true)
	}
}

func (p *Product) String() string {
	return p.Name
}

// returns 1 ProductMedia if one of them is image type
//
// else return nil
func (p *Product) GetFirstImage() *ProductMedia {
	for _, media := range p.Medias {
		if media.Type == IMAGE {
			return media
		}
	}

	return nil
}
