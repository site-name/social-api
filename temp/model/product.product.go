package model

import (
	"fmt"
	"net/http"

	"github.com/gosimple/slug"
	"github.com/mattermost/squirrel"
	"github.com/samber/lo"
	"github.com/sitename/sitename/modules/measurement"
	"gorm.io/gorm"
)

// ordering slug
type Product struct {
	Id                   string                 `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid();column:Id"`
	ProductTypeID        string                 `json:"product_type_id" gorm:"type:uuid;index:producttypeid_key;column:ProductTypeID"`
	Name                 string                 `json:"name" gorm:"type:varchar(250);column:Name"`
	Slug                 string                 `json:"slug" gorm:"type:varchar(255);uniqueIndex:product_slug_unique_key;column:Slug"`
	Description          StringInterface        `json:"description" gorm:"type:jsonb;column:Description"`
	DescriptionPlainText string                 `json:"description_plaintext" gorm:"column:DescriptionPlainText"`
	CategoryID           *string                `json:"category_id" gorm:"type:uuid;index:categoryid_key;column:CategoryID"`
	CreateAt             int64                  `json:"create_at" gorm:"type:bigint;autoCreateTime:milli;column:CreateAt"`
	UpdateAt             int64                  `json:"update_at" gorm:"type:bigint;autoCreateTime:milli;autoUpdateTime:milli;column:UpdateAt"`
	ChargeTaxes          *bool                  `json:"charge_taxes" gorm:"default:true;column:ChargeTaxes"` // default true
	Weight               *float32               `json:"weight" gorm:"column:Weight"`
	WeightUnit           measurement.WeightUnit `json:"weight_unit" gorm:"type:varchar(5);column:WeightUnit"`
	DefaultVariantID     *string                `json:"default_variant_id" gorm:"type:uuid;index:defaultvariantid_key;column:DefaultVariantID"`
	Rating               *float32               `json:"rating" gorm:"column:Rating"`
	ModelMetadata
	Seo

	// medias FileInfos `json:"-" gorm:"-"`

	ProductMedias           ProductMedias             `json:"-" gorm:"foreignKey:ProductID"`
	ProductType             *ProductType              `json:"-"`
	Category                *Category                 `json:"-"`
	ProductChannelListings  ProductChannelListings    `json:"-" gorm:"foreignKey:ProductID"`
	ProductVariants         ProductVariants           `json:"-" gorm:"foreignKey:ProductID"`
	Collections             Collections               `json:"-" gorm:"many2many:ProductCollections"`
	Sales                   Sales                     `json:"-" gorm:"many2many:SaleProducts"`
	Vouchers                Vouchers                  `json:"-" gorm:"many2many:VoucherProducts"`
	Attributes              AssignedProductAttributes `json:"-" gorm:"foreignKey:ProductID;constraint:OnDelete:CASCADE;"`
	AttributesRelated       []*AttributeProduct       `json:"-" gorm:"many2many:AssignedProductAttributes"`
	ShippingMethodsExcluded ShippingMethods           `json:"-" gorm:"many2many:ShippingMethodExcludedProducts"`
}

// column names for product table
const (
	ProductColumnId                   = "Id"
	ProductColumnProductTypeID        = "ProductTypeID"
	ProductColumnName                 = "Name"
	ProductColumnSlug                 = "Slug"
	ProductColumnDescription          = "Description"
	ProductColumnDescriptionPlainText = "DescriptionPlainText"
	ProductColumnCategoryID           = "CategoryID"
	ProductColumnCreateAt             = "CreateAt"
	ProductColumnUpdateAt             = "UpdateAt"
	ProductColumnChargeTaxes          = "ChargeTaxes"
	ProductColumnWeight               = "Weight"
	ProductColumnWeightUnit           = "WeightUnit"
	ProductColumnDefaultVariantID     = "DefaultVariantID"
	ProductColumnRating               = "Rating"
)

// func (p *Product) GetProductType() *ProductType                        { return p.productType }
// func (p *Product) SetProductType(pt *ProductType)                      { p.productType = pt }
// func (p *Product) GetProductVariants() ProductVariants                 { return p.productVariants }
// func (p *Product) SetProductVariants(pvs ProductVariants)              { p.productVariants = pvs }
// func (p *Product) GetCategory() *Category                              { return p.category }
// func (p *Product) SetCategory(c *Category)                             { p.category = c }
// func (p *Product) SetProductChannelListings(pc ProductChannelListings) { p.productChannelListings = pc }
// func (p *Product) GetProductChannelListings() ProductChannelListings   { return p.productChannelListings }
// func (p *Product) GetMedias() FileInfos          { return p.medias }
// func (p *Product) SetMedias(ms FileInfos)        { p.medias = ms }
func (c *Product) BeforeCreate(_ *gorm.DB) error { c.PreSave(); return c.IsValid() }
func (c *Product) BeforeUpdate(_ *gorm.DB) error { c.PreUpdate(); return c.IsValid() }
func (c *Product) TableName() string             { return ProductTableName }

type ProductCountByCategoryID struct {
	CategoryID   string `json:"category_id"`
	ProductCount uint64 `json:"product_count"`
}

// ProductFilterOption is used to compose squirrel sql queries
type ProductFilterOption struct {
	// native fields
	Conditions squirrel.Sqlizer

	Limit uint64

	HasNoProductVariants bool             // LEFT JOIN ProductVariants ON ... WHERE ProductVariants.ProductID IS NULL
	ProductVariantID     squirrel.Sqlizer // INNER JOIN ProductVariants ON ... WHERE ProductVariants.Id ...
	VoucherID            squirrel.Sqlizer // INNER JOIN ProductVouchers ON (...) WHERE ProductVouchers.VoucherID ...
	SaleID               squirrel.Sqlizer // INNER JOIN ProductSales ON (...) WHERE ProductSales.SaleID ...
	CollectionID         squirrel.Sqlizer // INNER JOIN ProductCollections ON ... WHERE ProductCollections.CollectionID ...

	// PrefetchRelatedAssignedProductAttributes bool
	// PrefetchRelatedVariants                  bool
	// PrefetchRelatedCollections               bool
	// PrefetchRelatedMedia                     bool
	// PrefetchRelatedProductType               bool
	// PrefetchRelatedCategory                  bool

	Preloads []string

	// Prefetch_Related_AssignedProductAttribute_AttributeValues                                 bool
	// Prefetch_Related_AssignedProductAttribute_AttributeProduct_Attribute                      bool
	// Prefetch_Related_ProductChannelListings                                                   bool
	// Prefetch_Related_ProductChannelListings_Channel                                           bool
	// Prefetch_Related_ProductVariants_Stocks                                                   bool
	// Prefetch_Related_ProductVariants_Stocks_Warehouses                                        bool
	// Prefetch_Related_ProductVariants_AssignedVariantAttributeValue_AttributeValues            bool
	// Prefetch_Related_ProductVariants_AssignedVariantAttributeValue_AttributeVariant_Attribute bool
	// Prefetch_Related_ProductVariants_ProductVariantChannelListing_Channel                     bool
	// Prefetch_Related_ProductVariants_ProductVariantChannelListings                            bool
}

type Products []*Product

func (ps Products) IDs() []string {
	return lo.Map(ps, func(p *Product, _ int) string {
		return p.Id
	})
}

func (ps Products) Contains(p *Product) bool {
	return p != nil && lo.SomeBy(ps, func(prd *Product) bool { return prd != nil && prd.Id == p.Id })
}

func (ps Products) ProductTypeIDs() []string {
	return lo.Map(ps, func(p *Product, _ int) string {
		return p.ProductTypeID
	})
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

func (p *Product) WeightString() string {
	if p.Weight == nil {
		return ""
	}

	u := p.WeightUnit
	if measurement.WEIGHT_UNIT_STRINGS[u] == "" {
		u = measurement.G
	}

	return fmt.Sprintf("%f %s", *p.Weight, u)
}

// Flat returns a slice of map[string]any items
// each item has keys are values, values of values of attributes of app/csv.ProductExportFields
func (ps Products) Flat() []StringInterface {
	var res = []StringInterface{}

	for _, prd := range ps {
		maxLength := max(
			len(prd.Collections),
			len(prd.ProductMedias),
			len(prd.Attributes),
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
			if i < len(prd.ProductMedias) {
				data["media__image"] = prd.ProductMedias[i].Image
			}
			if i < len(prd.Attributes) {
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
	if !IsValidId(p.ProductTypeID) {
		return NewAppError("Product.IsValid", "model.product.is_valid.product_type_id.app_error", nil, "please provide valid product type id", http.StatusBadRequest)
	}
	if p.CategoryID != nil && !IsValidId(*p.CategoryID) {
		return NewAppError("Product.IsValid", "model.product.is_valid.category_id.app_error", nil, "please provide valid category id", http.StatusBadRequest)
	}
	if p.Weight != nil {
		if _, ok := measurement.WEIGHT_UNIT_STRINGS[p.WeightUnit]; !ok {
			return NewAppError("Product.IsValid", "model.product.is_valid.weight_unit.app_error", nil, "please provide valid weight unit", http.StatusBadRequest)
		}
	}

	return nil
}

func (p *Product) PreSave() {
	p.commonPre()
	p.Slug = slug.Make(p.Name)
}

func (p *Product) PreUpdate() {
	p.commonPre()
}

func (p *Product) commonPre() {
	p.Name = SanitizeUnicode(p.Name)
	if p.WeightUnit == "" {
		p.WeightUnit = measurement.G
	}
	if p.ChargeTaxes == nil {
		p.ChargeTaxes = GetPointerOfValue(true)
	}
	p.Seo.commonPre()
}

// String returns exact product's name
func (p *Product) String() string {
	return p.Name
}

func (p *Product) DeepCopy() *Product {
	res := *p

	res.CategoryID = CopyPointer(p.CategoryID)
	res.DefaultVariantID = CopyPointer(p.DefaultVariantID)
	res.Weight = CopyPointer(p.Weight)
	res.Rating = CopyPointer(p.Rating)

	if p.Collections != nil {
		res.Collections = p.Collections.DeepCopy()
	}
	if p.ProductType != nil {
		res.ProductType = p.ProductType.DeepCopy()
	}
	if p.Attributes != nil {
		res.Attributes = p.Attributes.DeepCopy()
	}
	if p.ProductVariants != nil {
		res.ProductVariants = p.ProductVariants.DeepCopy()
	}
	if p.Category != nil {
		res.Category = p.Category.DeepCopy()
	}
	if p.ProductMedias != nil {
		res.ProductMedias = p.ProductMedias.DeepCopy()
	}
	if p.ProductChannelListings != nil {
		res.ProductChannelListings = p.ProductChannelListings.DeepCopy()
	}
	res.ModelMetadata = p.ModelMetadata.DeepCopy()

	return &res
}
