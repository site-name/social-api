package product_and_discount

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"unicode/utf8"

	"github.com/gosimple/slug"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/modules/json"
	"github.com/sitename/sitename/modules/measurement"
)

const (
	PRODUCT_TYPE_NAME_MAX_LENGTH = 250
	PRODUCT_TYPE_SLUG_MAX_LENGTH = 255
)

type ProductType struct {
	Id                 string   `json:"id"`
	Name               string   `json:"name"`
	Slug               string   `json:"slug"`
	HasVariants        *bool    `json:"has_variants"`
	IsShippingRequired *bool    `json:"is_shipping_required"`
	IsDigital          *bool    `json:"is_digital"`
	Weight             *float32 `json:"weight"`
	WeightUnit         string   `json:"weight_unit"`
	*model.ModelMetadata
}

func (p *ProductType) String() string {
	return p.Name
}

func (p *ProductType) createAppError(fieldName string) *model.AppError {
	id := fmt.Sprintf("model.product_type.is_valid.%s.app_error", fieldName)
	var details string
	if !strings.EqualFold(fieldName, "id") {
		details = "product_type_id=" + p.Id
	}

	return model.NewAppError("ProductType.IsValid", id, nil, details, http.StatusBadRequest)
}

func (p *ProductType) IsValid() *model.AppError {
	if !model.IsValidId(p.Id) {
		return p.createAppError("id")
	}
	if utf8.RuneCountInString(p.Name) > PRODUCT_TYPE_NAME_MAX_LENGTH {
		return p.createAppError("name")
	}
	if utf8.RuneCountInString(p.Slug) > PRODUCT_TYPE_SLUG_MAX_LENGTH {
		return p.createAppError("slug")
	}
	if p.Weight != nil && *p.Weight < 0 {
		return p.createAppError("weight")
	}
	if _, ok := measurement.WEIGHT_UNIT_STRINGS[strings.ToLower(p.WeightUnit)]; !ok {
		return p.createAppError("weight_unit")
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
	if p.Weight != nil && p.WeightUnit == "" {
		// p.Weight = model.NewFloat32(0)
		p.WeightUnit = measurement.STANDARD_WEIGHT_UNIT
	}
}

func (p *ProductType) PreUpdate() {
	p.Name = model.SanitizeUnicode(p.Name)
	p.Slug = slug.Make(p.Name)

	if p.Weight != nil && p.WeightUnit == "" {
		// p.Weight = model.NewFloat32(0)
		p.WeightUnit = measurement.STANDARD_WEIGHT_UNIT
	}
}

func (p *ProductType) ToJson() string {
	b, _ := json.JSON.Marshal(p)
	return string(b)
}

func ProductTypeFromJson(data io.Reader) *ProductType {
	var pt ProductType
	err := json.JSON.NewDecoder(data).Decode(&pt)
	if err != nil {
		return nil
	}
	return &pt
}
