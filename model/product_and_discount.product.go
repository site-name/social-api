package model

import (
	"fmt"
	"unicode/utf8"

	"github.com/Masterminds/squirrel"
	"github.com/gosimple/slug"
	"github.com/samber/lo"
	"github.com/sitename/sitename/modules/measurement"
	"github.com/sitename/sitename/modules/util"
)

// max lengths for some fields of products
const (
	PRODUCT_NAME_MAX_LENGTH = 250
	PRODUCT_SLUG_MAX_LENGTH = 255
)

// ordering slug
type Product struct {
	Id                   string                 `json:"id"`
	ProductTypeID        string                 `json:"product_type_id"`
	Name                 string                 `json:"name"`
	Slug                 string                 `json:"slug"`
	Description          StringInterface        `json:"description"`
	DescriptionPlainText string                 `json:"description_plaintext"`
	CategoryID           *string                `json:"category_id"`
	CreateAt             int64                  `json:"create_at"`
	UpdateAt             int64                  `json:"update_at"`
	ChargeTaxes          *bool                  `json:"charge_taxes"` // default true
	Weight               *float32               `json:"weight"`
	WeightUnit           measurement.WeightUnit `json:"weight_unit"`
	DefaultVariantID     *string                `json:"default_variant_id"`
	Rating               *float32               `json:"rating"`
	ModelMetadata
	Seo

	Collections               Collections               `json:"-" db:"-"`
	ProductType               *ProductType              `json:"-" db:"-"`
	AssignedProductAttributes AssignedProductAttributes `json:"-" db:"-"`
	ProductVariants           ProductVariants           `json:"-" db:"-"`
	Category                  *Category                 `json:"-" db:"-"`
	Medias                    FileInfos                 `json:"-" db:"-"`
	ProductChannelListings    ProductChannelListings    `json:"-" db:"-"`
}

// ProductFilterOption is used to compose squirrel sql queries
type ProductFilterOption struct {
	Id squirrel.Sqlizer

	// LEFT/INNER JOIN ProductVariants ON (...) WHERE ProductVariants.Id ...
	//
	// LEFT JOIN when squirrel.Eq{...: nil}, INNER JOIN otherwise
	ProductVariantID squirrel.Sqlizer
	VoucherID        squirrel.Sqlizer // SELECT * FROM Products INNER JOIN ProductVouchers ON (...) WHERE ProductVouchers.VoucherID ...
	SaleID           squirrel.Sqlizer // SELECT * FROM Products INNER JOIN ProductSales ON (...) WHERE ProductSales.SaleID ...
	CreateAt         squirrel.Sqlizer

	Limit *uint64

	PrefetchRelatedAssignedProductAttributes bool
	PrefetchRelatedVariants                  bool
	PrefetchRelatedCollections               bool
	PrefetchRelatedMedia                     bool
	PrefetchRelatedProductType               bool
	PrefetchRelatedCategory                  bool

	Prefetch_Related_AssignedProductAttribute_AttributeValues                                 bool
	Prefetch_Related_AssignedProductAttribute_AttributeProduct_Attribute                      bool
	Prefetch_Related_ProductChannelListings                                                   bool
	Prefetch_Related_ProductChannelListings_Channel                                           bool
	Prefetch_Related_ProductVariants_Stocks                                                   bool
	Prefetch_Related_ProductVariants_Stocks_Warehouses                                        bool
	Prefetch_Related_ProductVariants_AssignedVariantAttributeValue_AttributeValues            bool
	Prefetch_Related_ProductVariants_AssignedVariantAttributeValue_AttributeVariant_Attribute bool
	Prefetch_Related_ProductVariants_ProductVariantChannelListing_Channel                     bool
	Prefetch_Related_ProductVariants_ProductVariantChannelListings                            bool
}

type Products []*Product

func (ps Products) IDs() []string {
	return lo.Map(ps, func(p *Product, _ int) string { return p.Id })
}

func (p Products) CategoryIDs() []string {
	res := []string{}
	for _, product := range p {
		if product.CategoryID != nil {
			res = append(res, *product.CategoryID)
		}
	}

	return res
}

func (p *Product) WeightString() string {
	if p == nil || p.Weight == nil {
		return ""
	}

	u := p.WeightUnit
	if measurement.WEIGHT_UNIT_STRINGS[u] == "" {
		u = measurement.G
	}

	return fmt.Sprintf("%f %s", *p.Weight, u)
}

// Flat returns a slice of map[string]interface{} items
// each item has keys are values, values of values of attributes of app/csv.ProductExportFields
func (ps Products) Flat() []StringInterface {
	var res = []StringInterface{}

	for _, prd := range ps {
		maxLength := util.Max(
			len(prd.Collections),
			len(prd.Medias),
			len(prd.AssignedProductAttributes),
			len(prd.ProductVariants),
		)

		var categorySlug string
		var productTypeName string

		if prd.Category != nil {
			categorySlug = prd.Category.Slug
		}
		if prd.ProductType != nil {
			productTypeName = prd.ProductType.Name
		}

		for i := 0; i < maxLength; i++ {
			data := StringInterface{
				"id":                 prd.Id,
				"name":               prd.Name,
				"description_as_str": prd.Description,
				"category__slug":     categorySlug,
				"product_type__name": productTypeName,
				"charge_taxes":       *prd.ChargeTaxes,
				"product_weight":     prd.WeightString(),
			}

			if i < len(prd.Collections) {
				data["collections__slug"] = prd.Collections[i].Slug
			}
			if i < len(prd.Medias) {
				data["media__image"] = prd.Medias[i].Path
			}
			if i < len(prd.AssignedProductAttributes) {
				panic("not implemented")
			}
			if i < len(prd.ProductVariants) {
				data["variant_weight"] = prd.ProductVariants[i].WeightString()
				data["variants__id"] = prd.ProductVariants[i].Id
				data["variants__sku"] = prd.ProductVariants[i].Sku // can be nil
				data["variants__is_preorder"] = prd.ProductVariants[i].IsPreOrder
				data["variants__preorder_global_threshold"] = prd.ProductVariants[i].PreOrderGlobalThreshold // can be nil
				data["variants__preorder_end_date"] = prd.ProductVariants[i].PreorderEndDate                 // can be nil
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

func (p *Product) IsValid() *AppError {
	outer := CreateAppErrorForModel(
		"product.is_valid.%s.app_error",
		"product_id=",
		"Product.IsValid",
	)

	if !IsValidId(p.Id) {
		return outer("id", nil)
	}
	if !IsValidId(p.ProductTypeID) {
		return outer("product_type_id", &p.Id)
	}
	if p.CategoryID != nil && !IsValidId(*p.CategoryID) {
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
		p.Id = NewId()
	}
	p.CreateAt = GetMillis()
	p.UpdateAt = p.CreateAt
	p.Slug = slug.Make(p.Name)
}

func (p *Product) PreUpdate() {
	p.commonPre()
	p.UpdateAt = GetMillis()
}

func (p *Product) commonPre() {
	p.Name = SanitizeUnicode(p.Name)
	if p.WeightUnit == "" {
		p.WeightUnit = measurement.STANDARD_WEIGHT_UNIT
	}
	if p.ChargeTaxes == nil {
		p.ChargeTaxes = NewPrimitive(true)
	}
}

// String returns exact product's name
func (p *Product) String() string {
	return p.Name
}

func (p *Product) DeepCopy() *Product {
	res := *p

	if p.CategoryID != nil {
		res.CategoryID = NewPrimitive(*p.CategoryID)
	}
	if p.DefaultVariantID != nil {
		res.DefaultVariantID = NewPrimitive(*p.DefaultVariantID)
	}
	if p.Weight != nil {
		res.Weight = NewPrimitive(*p.Weight)
	}
	if p.Rating != nil {
		res.Rating = NewPrimitive(*p.Rating)
	}

	if p.Collections != nil {
		res.Collections = p.Collections.DeepCopy()
	}
	if p.ProductType != nil {
		res.ProductType = p.ProductType.DeepCopy()
	}
	if p.AssignedProductAttributes != nil {
		res.AssignedProductAttributes = p.AssignedProductAttributes.DeepCopy()
	}
	if p.ProductVariants != nil {
		res.ProductVariants = p.ProductVariants.DeepCopy()
	}
	if p.Category != nil {
		res.Category = p.Category.DeepCopy()
	}
	if p.Medias != nil {
		res.Medias = p.Medias.DeepCopy()
	}
	if p.ProductChannelListings != nil {
		res.ProductChannelListings = p.ProductChannelListings.DeepCopy()
	}

	return &res
}
