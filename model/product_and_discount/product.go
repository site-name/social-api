package product_and_discount

import (
	"unicode/utf8"

	"github.com/gosimple/slug"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model/seo"
	"github.com/sitename/sitename/modules/measurement"
)

// max lengths for some fields of products
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
	ChargeTaxes          *bool                  `json:"charge_taxes"` // default true
	Weight               *float32               `json:"weight"`
	WeightUnit           measurement.WeightUnit `json:"weight_unit"`
	DefaultVariantID     *string                `json:"default_variant_id"`
	Rating               *float32               `json:"rating"`
	model.ModelMetadata
	seo.Seo

	Collections Collections  `json:"-" db:"-"`
	ProductType *ProductType `json:"-" db:"-"`
}

// ProductFilterOption is used to compose squirrel sql queries
type ProductFilterOption struct {
	Id               *model.StringFilter
	ProductVariantID *model.StringFilter // LEFT/INNER JOIN ProductVariants ON (...) WHERE ProductVariants.Id ...
	VoucherIDs       []string            // SELECT * FROM Products WHERE Id IN (SELECT ProductID FROM ... WHERE VoucherID IN ?)
	SaleIDs          []string
}

type Products []*Product

func (p Products) IDs() []string {
	res := []string{}
	for _, product := range p {
		if product != nil {
			res = append(res, product.Id)
		}
	}

	return res
}

// PlainTextDescription Convert DraftJS JSON content to plain text
func (p *Product) PlainTextDescription() string {
	panic("not implemented")
}

func SortByAttributeFields() []string {
	return []string{"concatenated_values_order", "concatenated_values", "name"}
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
	p.commonPre()
	if p.Id == "" {
		p.Id = model.NewId()
	}
	p.CreateAt = model.GetMillis()
	p.UpdateAt = p.CreateAt
	p.Slug = slug.Make(p.Name)
}

func (p *Product) PreUpdate() {
	p.commonPre()
	p.UpdateAt = model.GetMillis()
}

func (p *Product) commonPre() {
	p.Name = model.SanitizeUnicode(p.Name)
	if p.WeightUnit == "" {
		p.WeightUnit = measurement.STANDARD_WEIGHT_UNIT
	}
	if p.ChargeTaxes == nil {
		p.ChargeTaxes = model.NewBool(true)
	}
}

// String returns exact product's name
func (p *Product) String() string {
	return p.Name
}
