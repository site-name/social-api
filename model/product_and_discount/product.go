package product_and_discount

import (
	"unicode/utf8"

	"github.com/Masterminds/squirrel"
	"github.com/gosimple/slug"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model/attribute"
	"github.com/sitename/sitename/model/file"
	"github.com/sitename/sitename/model/seo"
	"github.com/sitename/sitename/modules/measurement"
	"github.com/sitename/sitename/modules/util"
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
	Description          model.StringInterface  `json:"description"`
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

	Collections               Collections                         `json:"-" db:"-"`
	ProductType               *ProductType                        `json:"-" db:"-"`
	AssignedProductAttributes attribute.AssignedProductAttributes `json:"-" db:"-"`
	ProductVariants           ProductVariants                     `json:"-" db:"-"`
	Category                  *Category                           `json:"-" db:"-"`
	Medias                    file.FileInfos                      `json:"-" db:"-"`
}

// ProductFilterOption is used to compose squirrel sql queries
type ProductFilterOption struct {
	Id squirrel.Sqlizer

	// LEFT/INNER JOIN ProductVariants ON (...) WHERE ProductVariants.Id ...
	//
	// LEFT JOIN when squirrel.Eq{"": nil}, INNER JOIN otherwise
	ProductVariantID squirrel.Sqlizer
	VoucherID        squirrel.Sqlizer // SELECT * FROM Products INNER JOIN ProductVouchers ON (...) WHERE ProductVouchers.VoucherID ...
	SaleID           squirrel.Sqlizer // SELECT * FROM Products INNER JOIN ProductSales ON (...) WHERE ProductSales.SaleID ...
}

type ProductFilterByQueryOptions struct {
	CreateAt squirrel.Sqlizer
	Limit    *uint64

	PrefetchRelatedAssignedProductAttributes bool
	PrefetchRelatedVariants                  bool
	PrefetchRelatedCollections               bool
	PrefetchRelatedMedia                     bool
	PrefetchRelatedProductType               bool
	PrefetchRelatedCategory                  bool
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

func (p Products) CategoryIDs() []string {
	res := []string{}
	for _, product := range p {
		if product != nil && product.CategoryID != nil {
			res = append(res, *product.CategoryID)
		}
	}

	return res
}

func (ps Products) Flat() []model.StringInterface {
	var res = []model.StringInterface{}

	for _, prd := range ps {

		maxLength := util.Max(
			len(prd.Collections),
			len(prd.Medias),
			len(prd.AssignedProductAttributes),
			len(prd.ProductVariants),
		)

		for i := 0; i < maxLength; i++ {
			data := model.StringInterface{
				"id":                 prd.Id,
				"name":               prd.Name,
				"description_as_str": prd.DescriptionPlainText,
				"charge_taxes":       *prd.ChargeTaxes,
			}

			if i < len(prd.Collections) {
				data["collections__slug"] = prd.Collections[i].Slug
			}

			if i < len(prd.Medias) {
				data["media__image"] = prd.Medias[i].Path
			}

			if i < len(prd.AssignedProductAttributes) {

			}

			if i < len(prd.ProductVariants) {

			}

			res = append(res, data)
		}
	}

	return res
}

// PlainTextDescription Convert DraftJS JSON content to plain text
func (p *Product) PlainTextDescription() string {
	return p.Name
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

func (p *Product) DeepCopy() *Product {
	res := *p

	if p.Collections != nil {
		res.Collections = p.Collections.DeepCopy()
	}
	if p.ProductType != nil {
		res.ProductType = p.ProductType.DeepCopy()
	}
	if len(p.AssignedProductAttributes) > 0 {
		res.AssignedProductAttributes = p.AssignedProductAttributes.DeepCopy()
	}
	if len(p.ProductVariants) > 0 {
		res.ProductVariants = p.ProductVariants.DeepCopy()
	}
	if p.Category != nil {
		res.Category = p.Category.DeepCopy()
	}
	if len(p.Medias) > 0 {
		res.Medias = p.Medias.DeepCopy()
	}

	return &res
}
