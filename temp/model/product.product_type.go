package model

import (
	"net/http"

	"github.com/gosimple/slug"
	"github.com/mattermost/squirrel"
	"github.com/sitename/sitename/modules/measurement"
	"gorm.io/gorm"
)

type ProductTypeKind string

func (p ProductTypeKind) IsValid() bool {
	return ProductTypeKindStrings[p] != ""
}

// some valid value for product type kind
const (
	NORMAL    ProductTypeKind = "normal"
	GIFT_CARD ProductTypeKind = "gift_card"
)

var ProductTypeKindStrings = map[ProductTypeKind]string{
	NORMAL:    "A standard product type.",
	GIFT_CARD: "A gift card product type.",
}

// Orderby Slug
type ProductType struct {
	Id                 string                 `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid();column:Id"`
	Name               string                 `json:"name" gorm:"type:varchar(250);column:Name"`
	Slug               string                 `json:"slug" gorm:"type:varchar(255);column:Slug;uniqueIndex:slug_key"`
	Kind               ProductTypeKind        `json:"kind" gorm:"type:varchar(32);column:Kind"`
	HasVariants        *bool                  `json:"has_variants" gorm:"column:HasVariants;default:true"`                // default true
	IsShippingRequired *bool                  `json:"is_shipping_required" gorm:"column:IsShippingRequired;default:true"` // default true
	IsDigital          *bool                  `json:"is_digital" gorm:"column:IsDigital;default:false"`                   // default false
	Weight             *float32               `json:"weight" gorm:"column:Weight;default:0"`
	WeightUnit         measurement.WeightUnit `json:"weight_unit" gorm:"column:WeightUnit;type:varchar(5)"`
	ModelMetadata

	ProductAttributes Attributes `json:"-" gorm:"many2many:AttributeProducts"`
	VariantAttributes Attributes `json:"-" gorm:"many2many:AttributeVariants"`
}

// column names of product type's fields
const (
	ProductTypeColumnId                 = "Id"
	ProductTypeColumnName               = "Name"
	ProductTypeColumnSlug               = "Slug"
	ProductTypeColumnKind               = "Kind"
	ProductTypeColumnHasVariants        = "HasVariants"
	ProductTypeColumnIsShippingRequired = "IsShippingRequired"
	ProductTypeColumnIsDigital          = "IsDigital"
	ProductTypeColumnWeight             = "Weight"
	ProductTypeColumnWeightUnit         = "WeightUnit"
)

func (c *ProductType) BeforeCreate(_ *gorm.DB) error { c.PreSave(); return c.IsValid() }
func (c *ProductType) BeforeUpdate(_ *gorm.DB) error { c.commonPre(); return c.IsValid() }
func (c *ProductType) TableName() string             { return ProductTypeTableName }

// ProductTypeFilterOption is used to build squirrel sql queries
type ProductTypeFilterOption struct {
	Conditions squirrel.Sqlizer

	AttributeProducts_AttributeID squirrel.Sqlizer // INNER JOIN AttributeProducts ON (...) WHERE AttributeProducts.AttributeID ...
	AttributeVariants_AttributeID squirrel.Sqlizer // INNER JOIN AttributeVariants ON (...) WHERE AttributeVariants.AttributeID ...

	CountTotal              bool
	GraphqlPaginationValues GraphqlPaginationValues
}

func (p *ProductType) DeepCopy() *ProductType {
	if p == nil {
		return nil
	}

	res := *p
	res.HasVariants = CopyPointer(p.HasVariants)
	res.IsShippingRequired = CopyPointer(p.IsShippingRequired)
	res.IsDigital = CopyPointer(p.IsDigital)
	res.Weight = CopyPointer(p.Weight)
	res.ModelMetadata = p.ModelMetadata.DeepCopy()

	return &res
}

func (p *ProductType) String() string {
	return p.Name
}

func (p *ProductType) IsValid() *AppError {
	if !p.Kind.IsValid() {
		return NewAppError("ProductType.IsValid", "model.product_type.is_valid.kind.app_error", nil, "please provide valid product kind", http.StatusBadRequest)
	}
	if p.Weight != nil && *p.Weight < 0 {
		return NewAppError("ProductType.IsValid", "model.product_type.is_valid.weight.app_error", nil, "please provide valid product weight", http.StatusBadRequest)
	}
	if _, ok := measurement.WEIGHT_UNIT_STRINGS[p.WeightUnit]; !ok {
		return NewAppError("ProductType.IsValid", "model.product_type.is_valid.weight_unit.app_error", nil, "please provide valid weight unit", http.StatusBadRequest)
	}

	return nil
}

func (p *ProductType) PreSave() {
	p.commonPre()
	p.Slug = slug.Make(p.Name)
}

func (p *ProductType) commonPre() {
	p.Name = SanitizeUnicode(p.Name)

	if p.HasVariants == nil {
		p.HasVariants = GetPointerOfValue(true)
	}
	if p.IsShippingRequired == nil {
		p.IsShippingRequired = GetPointerOfValue(true)
	}
	if p.IsDigital == nil {
		p.IsDigital = GetPointerOfValue(false)
	}
	if p.Weight == nil {
		p.Weight = GetPointerOfValue[float32](0)
	}
	if p.Weight != nil && p.WeightUnit == "" {
		p.WeightUnit = measurement.STANDARD_WEIGHT_UNIT
	}
}

// IsGiftcard checks if current product type has kind of "gift_card"
func (p *ProductType) IsGiftcard() bool {
	return p.Kind == GIFT_CARD
}
