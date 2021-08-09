package product_and_discount

import (
	"unicode/utf8"

	"github.com/gosimple/slug"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/modules/measurement"
)

// max lengths for some product type's fields
const (
	PRODUCT_TYPE_NAME_MAX_LENGTH = 250
	PRODUCT_TYPE_SLUG_MAX_LENGTH = 255
)

type ProductType struct {
	Id                 string                 `json:"id"`
	Name               string                 `json:"name"`
	Slug               string                 `json:"slug"`
	HasVariants        *bool                  `json:"has_variants"`
	IsShippingRequired *bool                  `json:"is_shipping_required"` // default true
	IsDigital          *bool                  `json:"is_digital"`           // default false
	Weight             *float32               `json:"weight"`
	WeightUnit         measurement.WeightUnit `json:"weight_unit"`
	model.ModelMetadata
}

func (p *ProductType) String() string {
	return p.Name
}

func (p *ProductType) IsValid() *model.AppError {
	outer := model.CreateAppErrorForModel(
		"model.product_type.is_valid.%s.app_error",
		"product_type_id=",
		"ProductType.IsValid")

	if !model.IsValidId(p.Id) {
		return outer("id", nil)
	}
	if utf8.RuneCountInString(p.Name) > PRODUCT_TYPE_NAME_MAX_LENGTH {
		return outer("name", &p.Id)
	}
	if utf8.RuneCountInString(p.Slug) > PRODUCT_TYPE_SLUG_MAX_LENGTH {
		return outer("slug", &p.Id)
	}
	if p.Weight != nil && *p.Weight < 0 {
		return outer("weight", &p.Id)
	}
	if _, ok := measurement.WEIGHT_UNIT_STRINGS[p.WeightUnit]; !ok {
		return outer("weight_unit", &p.Id)
	}

	return nil
}

func (p *ProductType) PreSave() {
	if p.Id == "" {
		p.Id = model.NewId()
	}
	p.Name = model.SanitizeUnicode(p.Name)
	p.Slug = slug.Make(p.Name)

	if p.HasVariants == nil {
		p.HasVariants = model.NewBool(true)
	}
	if p.IsShippingRequired == nil {
		p.IsShippingRequired = model.NewBool(true)
	}
	if p.IsDigital == nil {
		p.IsDigital = model.NewBool(false)
	}
	if p.Weight == nil {
		p.Weight = model.NewFloat32(0)
	}
	if p.Weight != nil && p.WeightUnit == "" {
		p.WeightUnit = measurement.STANDARD_WEIGHT_UNIT
	}
}

func (p *ProductType) PreUpdate() {
	p.Name = model.SanitizeUnicode(p.Name)
	p.Slug = slug.Make(p.Name)

	if p.Weight != nil && p.WeightUnit == "" {
		p.WeightUnit = measurement.STANDARD_WEIGHT_UNIT
	}
	if p.IsShippingRequired == nil {
		p.IsShippingRequired = model.NewBool(true)
	}
}

func (p *ProductType) ToJson() string {
	return model.ModelToJson(p)
}
