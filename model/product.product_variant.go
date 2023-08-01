package model

import (
	"fmt"
	"net/http"

	"github.com/Masterminds/squirrel"
	"github.com/samber/lo"
	"github.com/sitename/sitename/modules/measurement"
	"gorm.io/gorm"
)

// sort by sku
type ProductVariant struct {
	Id                      string                 `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid();column:Id"`
	Name                    string                 `json:"name" gorm:"type:varchar(255);column:Name"` // varchar(255)
	ProductID               string                 `json:"product_id" gorm:"type:uuid;index:productvariants_productid_index_key;column:ProductID"`
	Sku                     string                 `json:"sku" gorm:"type:varchar(255);column:Sku"` // varchar(255)
	Weight                  *float32               `json:"weight" gorm:"column:Weight"`
	WeightUnit              measurement.WeightUnit `json:"weight_unit" gorm:"column:WeightUnit"`
	TrackInventory          *bool                  `json:"track_inventory" gorm:"column:TrackInventory"` // default *true
	IsPreOrder              bool                   `json:"is_preorder" gorm:"column:IsPreOrder"`
	PreorderEndDate         *int64                 `json:"preorder_end_date" column:"type:bigint;column:PreorderEndDate"`
	PreOrderGlobalThreshold *int                   `json:"preorder_global_threshold" gorm:"type:smallint;column:PreOrderGlobalThreshold"`
	Sortable
	ModelMetadata

	variantChannelListings ProductVariantChannelListings `json:"-" gorm:"-"`
	Stocks                 Stocks                        `json:"-" gorm:"foreignKey:ProductVariantID"`
	DigitalContent         *DigitalContent               `json:"-"` // for storing value returned by prefetching
	Sales                  Sales                         `json:"-" gorm:"many2many:SaleProductVariants"`
	Vouchers               Vouchers                      `json:"-" gorm:"many2many:VoucherVariants"`
	Medias                 ProductMedias                 `json:"-" gorm:"many2many:VariantMedias"`
	Attributes             []*AssignedVariantAttribute   `json:"-" gorm:"foreignKey:VariantID"`
	AttributesRelated      []*AttributeVariant           `json:"-" gorm:"many2many:AssignedVariantAttributes"`
	WishlistItems          WishlistItems                 `json:"-" gorm:"many2many:WishlistItemProductVariants"`
	OrderLines             OrderLines                    `json:"-" gorm:"foreignKey:VariantID"`
	Product                *Product                      `json:"-"`
}

func (c *ProductVariant) BeforeCreate(_ *gorm.DB) error { c.commonPre(); return c.IsValid() }
func (c *ProductVariant) BeforeUpdate(_ *gorm.DB) error { c.commonPre(); return c.IsValid() }
func (c *ProductVariant) TableName() string             { return ProductVariantTableName }

// ProductVariantFilterOption is used to build sql queries
type ProductVariantFilterOption struct {
	Conditions squirrel.Sqlizer

	WishlistItemID squirrel.Sqlizer // INNER JOIN WishlistItemProductVariants ON (...) WHERE WishlistItemProductVariants.WishlistItemID ...
	WishlistID     squirrel.Sqlizer // INNER JOIN WishlistItemProductVariants ON (...) INNER JOIN WishlistItems ON (...) WHERE WishlistItems.WishlistID ...

	ProductVariantChannelListingPriceAmount squirrel.Sqlizer // INNER JOIN `ProductVariantChannelListing` ON ... WHERE ProductVariantChannelListing.PriceAmount ...
	ProductVariantChannelListingChannelSlug squirrel.Sqlizer // INNER JOIN `ProductVariantChannelListing` ON ... INNER JOIN Channels ON ... WHERE Channels.Slug ...
	ProductVariantChannelListingChannelID   squirrel.Sqlizer // INNER JOIN ProductVariantChannelListing ON ... WHERE ProductVariantChannelListing.ChannelID ...

	Distinct bool // if true, use SELECT DISTINCT

	VoucherID squirrel.Sqlizer // INNER JOIN VariantVouchers ON ... WHERE VariantVouchers.VoucherID ...
	SaleID    squirrel.Sqlizer // INNER JOIN VariantSales ON ... WHERE VariantSales.SaleID ...

	SelectRelatedDigitalContent bool // if true, JOIN Digital content table and attach related values to returning values(s)
}

func (p *ProductVariant) WeightString() string {
	if p == nil || p.Weight == nil {
		return ""
	}

	u := p.WeightUnit
	if measurement.WEIGHT_UNIT_STRINGS[u] == "" {
		u = measurement.G
	}
	return fmt.Sprintf("%f %s", *p.Weight, u)
}

type ProductVariants []*ProductVariant

func (p ProductVariants) DeepCopy() ProductVariants {
	return lo.Map(p, func(v *ProductVariant, _ int) *ProductVariant { return v.DeepCopy() })
}

// func (s *ProductVariant) SetStocks(stk Stocks) { s.stocks = stk }
// func (p *ProductVariant) GetStocks() Stocks    { return p.stocks }

// func (s *ProductVariant) SetProduct(prd *Product) { s.product = prd }
// func (p *ProductVariant) GetProduct() *Product    { return p.product }

// func (s *ProductVariant) SetDigitalContent(d *DigitalContent) { s.digitalContent = d }
// func (p *ProductVariant) GetDigitalContent() *DigitalContent  { return p.digitalContent }

func (s *ProductVariant) SetVariantChannelListings(d ProductVariantChannelListings) {
	s.variantChannelListings = d
}
func (s *ProductVariant) AppendVariantChannelListing(d *ProductVariantChannelListing) {
	s.variantChannelListings = append(s.variantChannelListings, d)
}
func (p *ProductVariant) GetVariantChannelListings() ProductVariantChannelListings {
	return p.variantChannelListings
}

// FilterNils returns new ProductVariants contains all non-nil items from current ProductVariants
func (p ProductVariants) FilterNils() ProductVariants {
	return lo.Filter(p, func(v *ProductVariant, _ int) bool { return v != nil })
}

func (p ProductVariants) IDs() []string {
	return lo.Map(p, func(v *ProductVariant, _ int) string { return v.Id })
}

// ProductIDs returns all product ids of current product variants
func (p ProductVariants) ProductIDs() []string {
	return lo.Map(p, func(v *ProductVariant, _ int) string { return v.ProductID })
}

func (p *ProductVariant) IsValid() *AppError {
	if !IsValidId(p.ProductID) {
		return NewAppError("ProductVariant.IsValid", "model.product_variant.is_valid.product_id.app_error", nil, "please provide valid product id", http.StatusBadRequest)
	}
	if p.Weight != nil && *p.Weight <= 0 {
		return NewAppError("ProductVariant.IsValid", "model.product_variant.is_valid.weight.app_error", nil, "please provide valid weight", http.StatusBadRequest)
	}
	if _, ok := measurement.WEIGHT_UNIT_CONVERSION[p.WeightUnit]; !ok {
		return NewAppError("ProductVariant.IsValid", "model.product_variant.is_valid.weight_unit.app_error", nil, "please provide valid weight unit", http.StatusBadRequest)
	}

	return nil
}

// String returns exact product variant name or Sku depends on their truth value
func (p *ProductVariant) String() string {
	if p.Name != "" {
		return p.Name
	}

	return fmt.Sprintf("ID:%s", p.Id)
}

func (p *ProductVariant) IsPreorderActive() bool {
	return p.IsPreOrder && (p.PreorderEndDate == nil || (p.PreorderEndDate != nil && GetMillis() <= *p.PreorderEndDate))
}

func (p *ProductVariant) commonPre() {
	p.Name = SanitizeUnicode(p.Name)
	if p.TrackInventory == nil {
		p.TrackInventory = NewPrimitive(true)
	}
	if p.Weight != nil && p.WeightUnit == "" {
		p.WeightUnit = measurement.G
	}
	p.ModelMetadata.PopulateFields()

}

func (p *ProductVariant) DeepCopy() *ProductVariant {
	if p == nil {
		return nil
	}

	res := *p

	res.Weight = CopyPointer(p.Weight)
	res.TrackInventory = CopyPointer(p.TrackInventory)
	res.PreorderEndDate = CopyPointer(p.PreorderEndDate)
	res.PreOrderGlobalThreshold = CopyPointer(p.PreOrderGlobalThreshold)
	res.SortOrder = CopyPointer(p.SortOrder)

	res.ModelMetadata = p.ModelMetadata.DeepCopy()
	if p.Product != nil {
		res.Product = p.Product.DeepCopy()
	}
	if p.DigitalContent != nil {
		res.DigitalContent = p.DigitalContent.DeepCopy()
	}
	if p.Stocks != nil {
		res.Stocks = p.Stocks.DeepCopy()
	}

	return &res
}

type ProductVariantTranslation struct {
	Id               string           `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid();column:Id"`
	LanguageCode     LanguageCodeEnum `json:"language_code" gorm:"type:varchar(5);column:LanguageCode;index:languagecode_productvariantid_key"`
	ProductVariantID string           `json:"product_variant_id" gorm:"type:uuid;column:ProductVariantID;index:languagecode_productvariantid_key"`
	Name             string           `json:"name" gorm:"type:varchar(255);column:Name"`
}

func (c *ProductVariantTranslation) BeforeCreate(_ *gorm.DB) error { c.commonPre(); return c.IsValid() }
func (c *ProductVariantTranslation) BeforeUpdate(_ *gorm.DB) error { c.commonPre(); return c.IsValid() }
func (c *ProductVariantTranslation) TableName() string             { return ProductVariantTranslationTableName }

// ProductVariantTranslationFilterOption is used to build squirrel sql queries
type ProductVariantTranslationFilterOption struct {
	Conditions squirrel.Sqlizer
}

func (p *ProductVariantTranslation) String() string {
	if p.Name != "" {
		return p.Name
	}

	return p.ProductVariantID
}

func (p *ProductVariantTranslation) commonPre() {
	p.Name = SanitizeUnicode(p.Name)
}

func (p *ProductVariantTranslation) IsValid() *AppError {
	if !IsValidId(p.ProductVariantID) {
		return NewAppError("ProductVariantTranslation.IsValid", "model.product_variant_translation.is_valid.product_variant_id.app_error", nil, "please provide valid product variant id", http.StatusBadRequest)
	}
	if !p.LanguageCode.IsValid() {
		return NewAppError("ProductVariantTranslation.IsValid", "model.product_variant_translation.is_valid.language_code.app_error", nil, "please provide valid language code", http.StatusBadRequest)
	}

	return nil
}
