package model

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"unicode/utf8"

	"github.com/sitename/sitename/modules/json"
)

// standard units for weight
const (
	G     = "g"
	LB    = "lb"
	OZ    = "oz"
	KG    = "kg"
	TONNE = "tonne"
)

const (
	PRODUCT_TYPE_NAME_MAX_LENGTH = 250
	PRODUCT_TYPE_SLUG_MAX_LENGTH = 255
)

var WeightUnitString = map[string]string{
	G:     "Gram",
	LB:    "Pound",
	OZ:    "Ounce",
	KG:    "kg",
	TONNE: "Tonne",
}

type ProductType struct {
	Id                 string   `json:"id"`
	Name               string   `json:"name"`
	Slug               string   `json:"slug"`
	HasVariants        *bool    `json:"has_variants"`
	IsShippingRequired *bool    `json:"is_shipping_required"`
	IsDigital          *bool    `json:"is_digital"`
	Weight             *float32 `json:"weight"`
	WeightUnit         string   `json:"weight_unit"`
}

func (p *ProductType) createAppError(fieldName string) *AppError {
	id := fmt.Sprintf("model.product_type.is_valid.%s.app_error", fieldName)
	var details string
	if !strings.EqualFold(fieldName, "id") {
		details = "product_type_id=" + p.Id
	}

	return NewAppError("ProductType.IsValid", id, nil, details, http.StatusBadRequest)
}

func (p *ProductType) IsValid() *AppError {
	if !IsValidId(p.Id) {
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
	if _, ok := WeightUnitString[p.WeightUnit]; !ok {
		return p.createAppError("weight_unit")
	}

	return nil
}

func (p *ProductType) PreSave() {
	if p.Id == "" {
		p.Id = NewId()
	}
	p.Name = SanitizeUnicode(p.Name)

	if p.HasVariants == nil {
		p.HasVariants = NewBool(true)
	}
	if p.IsShippingRequired == nil {
		p.IsShippingRequired = NewBool(true)
	}
	if p.IsDigital == nil {
		p.IsDigital = NewBool(false)
	}
	if p.Weight == nil {
		p.Weight = NewFloat32(0)
		p.WeightUnit = KG
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
